// File: internal/game/inventory/service.go
package inventory

import (
	"context"
	"fmt"
	"strings"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
	"go.uber.org/zap"
)

type Service interface {
	GetInventory(ctx context.Context, userID, guildID string) (*Inventory, []*item.ItemInstance, error)
	AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error
	GrantStarterItems(ctx context.Context, userID, guildID string) error
	ConsumeItems(ctx context.Context, userID, guildID string, itemsToConsume map[string]int64) error
	UseItem(ctx context.Context, userID, guildID, instanceID string) (string, error)
	DismantleItem(ctx context.Context, userID, guildID, instanceID string) (string, error)
}

type inventoryService struct {
	invRepo  Repository
	itemRepo item.Repository
	cultSvc  cultivation.Service
}

func NewService(invRepo Repository, itemRepo item.Repository, cultSvc cultivation.Service) Service {
	return &inventoryService{invRepo: invRepo, itemRepo: itemRepo, cultSvc: cultSvc}
}

func (s *inventoryService) GetInventory(ctx context.Context, userID, guildID string) (*Inventory, []*item.ItemInstance, error) {
	inv, err := s.invRepo.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		return nil, nil, err
	}
	logger.L().Debug("InventoryService.GetInventory", zap.String("userId", userID), zap.String("guildId", guildID))
	items, err := s.itemRepo.GetInstancesByUser(ctx, userID, guildID)
	logger.L().Debug("InventoryRepository result",
		zap.String("userId", userID),
		zap.Int("loadedCount", len(items)),
		zap.Error(err),
	)
	return inv, items, err
}

func (s *inventoryService) AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	def, ok := item.GetDefinition(definitionID)
	if !ok {
		return apperrors.ErrInvalidInput
	}

	_, items, err := s.GetInventory(ctx, userID, guildID)
	if err != nil {
		return err
	}

	if def.Stackable {
		for _, it := range items {
			if it.DefinitionID == definitionID {
				return s.itemRepo.AdjustQuantity(ctx, it.InstanceID, userID, guildID, quantity)
			}
		}
	}

	// Cần tạo ô mới -> Xin cấp phát 1 ô an toàn ở tầng Database (Atomic)
	if err := s.invRepo.AcquireSlot(ctx, userID, guildID); err != nil {
		return apperrors.ErrInventoryFull
	}

	inst := &item.ItemInstance{
		InstanceID:   utils.NewSessionID(),
		DefinitionID: definitionID,
		UserID:       userID,
		GuildID:      guildID,
		Quantity:     quantity,
	}
	err = s.itemRepo.CreateInstance(ctx, inst)
	if err != nil {
		_ = s.invRepo.ReleaseSlot(ctx, userID, guildID) // Trả lại ô nếu lỗi tạo item
	}
	return err
}

func (s *inventoryService) ConsumeItems(ctx context.Context, userID, guildID string, itemsToConsume map[string]int64) error {
	_, currentItems, err := s.GetInventory(ctx, userID, guildID)
	if err != nil {
		return err
	}

	// 1. Kiểm tra xem có đủ nguyên liệu không
	playerHas := make(map[string]int64)
	for _, it := range currentItems {
		playerHas[it.DefinitionID] += it.Quantity
	}

	for defID, requiredQty := range itemsToConsume {
		if playerHas[defID] < requiredQty {
			def, ok := item.GetDefinition(defID)
			name := defID // fallback nếu chưa có dữ liệu item
			if ok && def.Name != "" {
				name = def.Name
			}
			return fmt.Errorf("không đủ nguyên liệu: thiếu **%s**", name)
		}
	}

	// 2. Trừ nguyên liệu an toàn với Rollback (Saga Pattern)
	deductedHistory := make(map[string]int64)

	for defID, requiredQty := range itemsToConsume {
		rem := requiredQty
		for _, it := range currentItems {
			if it.DefinitionID == defID && it.Quantity > 0 {
				toDeduct := it.Quantity
				if toDeduct > rem {
					toDeduct = rem
				}

				if err := s.itemRepo.AdjustQuantity(ctx, it.InstanceID, userID, guildID, -toDeduct); err != nil {
					// Rollback toàn bộ quá trình nếu trừ lỗi (vd: bị xài hết bởi thread khác)
					for instID, amount := range deductedHistory {
						_ = s.itemRepo.AdjustQuantity(ctx, instID, userID, guildID, amount)
					}
					return fmt.Errorf("tài nguyên %s đang biến động, vui lòng thử lại", defID)
				}

				deductedHistory[it.InstanceID] += toDeduct
				rem -= toDeduct

				if rem <= 0 {
					break
				}
			}
		}

		if rem > 0 {
			// Lỗi race condition khiến kho không đủ như tính toán ban đầu
			for instID, amount := range deductedHistory {
				_ = s.itemRepo.AdjustQuantity(ctx, instID, userID, guildID, amount)
			}
			return fmt.Errorf("nguyên liệu %s bị thiếu hụt đột ngột", defID)
		}
	}

	// 3. Dọn dẹp túi đồ nếu nguyên liệu đã dùng hết
	for instID := range deductedHistory {
		// DeleteInstance có sẵn "$lte: 0" ở Tầng Repo, gọi hoàn toàn an toàn
		if err := s.itemRepo.DeleteInstance(ctx, instID, userID, guildID); err == nil {
			_ = s.invRepo.ReleaseSlot(ctx, userID, guildID)
		}
	}
	return nil
}

func (s *inventoryService) GrantStarterItems(ctx context.Context, userID, guildID string) error {
	inv, err := s.invRepo.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		return err
	}
	if inv.StarterGranted {
		return nil
	}

	var hasError bool
	grant := func(defID string, qty int64) {
		if err := s.AddItem(ctx, userID, guildID, defID, qty); err != nil {
			logger.L().Error("Lỗi cấp vật phẩm tân thủ", zap.String("userId", userID), zap.String("defID", defID), zap.Error(err))
			hasError = true
		}
	}

	grant("pill_exp_tu_khi_d", 3)
	grant("eq_weapon_moc_kiem_d", 1)
	grant("eq_armor_vai_tho_d", 1)
	grant("mat_enhance_hac_thiet_d", 3)

	if hasError {
		return fmt.Errorf("không thể khởi tạo đủ vật phẩm tân thủ")
	}

	return s.invRepo.MarkStarterGranted(ctx, userID, guildID)
}

func (s *inventoryService) UseItem(ctx context.Context, userID, guildID, instanceID string) (string, error) {
	inst, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return "", err
	}

	def, ok := item.GetDefinition(inst.DefinitionID)
	if !ok {
		return "", fmt.Errorf("định nghĩa vật phẩm không tồn tại")
	}
	if !def.Usable || def.Type != item.TypePill {
		return "", apperrors.ErrItemNotUsable
	}

	// 0. Chặn đan thể lực cũ: nếu item CHỈ CÓ effect stamina thì không cho xài, không trừ đồ
	if def.Effects != nil {
		if s, ok := def.Effects["stamina"]; ok && s > 0 {
			if def.Effects["exp"] == 0 && def.Effects["breakthrough_chance"] == 0 {
				return "Cơ chế Thể Lực đã được gỡ bỏ. Vật phẩm này hiện không còn tác dụng và sẽ không bị tiêu hao.", nil
			}
		}
	}

	// 1. Trừ số lượng vật phẩm (Atomic)
	if err := s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, -1); err != nil {
		return "", err
	}

	// 2. Áp dụng hiệu ứng
	expGained := int64(0)
	breakthroughBuff := 0

	if def.Effects != nil {
		expGained = int64(def.Effects["exp"])
		breakthroughBuff = def.Effects["breakthrough_chance"]
	}

	// Fallback an toàn nếu Effects chưa được định nghĩa đầy đủ trong Registry
	if expGained == 0 && breakthroughBuff == 0 {
		logger.L().Error("Vật phẩm thiếu cấu hình effect", zap.String("defID", inst.DefinitionID))
		return "", fmt.Errorf("vật phẩm **%s** bị lỗi cấu hình hệ thống, vui lòng báo Admin", def.Name)
	}

	var messages []string
	if expGained > 0 {
		if err := s.cultSvc.AddExperience(ctx, userID, guildID, expGained); err != nil {
			// Hoàn trả lại vật phẩm nếu hiệu ứng thất bại
			if rbErr := s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, 1); rbErr != nil {
				logger.L().Error("CRITICAL: Mất vật phẩm do rollback thất bại", zap.String("instanceId", instanceID), zap.Error(rbErr))
			}
			return "", fmt.Errorf("không thể hấp thụ linh khí lúc này: %w", err)
		}
		messages = append(messages, fmt.Sprintf("nhận **%d** tu vi", expGained))
	}
	if breakthroughBuff > 0 {
		// Hiện tại CultivationProfile chưa có field lưu Buff đột phá. Tạm thời báo text lên UI.
		messages = append(messages, fmt.Sprintf("tăng **%d%%** tỉ lệ đột phá", breakthroughBuff))
	}

	var message string
	if len(messages) > 0 {
		message = fmt.Sprintf("Sử dụng **%s** thành công! %s.", def.Name, strings.Join(messages, ", "))
	} else {
		message = fmt.Sprintf("Sử dụng **%s** thành công nhưng không có tác dụng gì!", def.Name)
	}

	// 3. Gọi hàm xóa (Hàm này đã được bảo vệ Atomic $lte: 0 ở Tầng Repository)
	if err := s.itemRepo.DeleteInstance(ctx, instanceID, userID, guildID); err == nil {
		_ = s.invRepo.ReleaseSlot(ctx, userID, guildID)
	}
	return message, nil
}

func (s *inventoryService) DismantleItem(ctx context.Context, userID, guildID, instanceID string) (string, error) {
	inst, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return "", err
	}
	def, ok := item.GetDefinition(inst.DefinitionID)
	if !ok || def.Type != item.TypeEquipment {
		return "", fmt.Errorf("chỉ có thể phân giải trang bị")
	}

	matDefID := "mat_enhance_hac_thiet_d"
	qty := int64(1)
	switch def.Rarity {
	case item.RarityC:
		qty = 2
	case item.RarityB:
		qty = 4
	case item.RarityA:
		qty = 10
	case item.RarityS:
		qty = 30
	}

	if err := s.itemRepo.DeleteInstance(ctx, instanceID, userID, guildID); err != nil {
		return "", err
	}
	_ = s.invRepo.ReleaseSlot(ctx, userID, guildID)

	if err := s.AddItem(ctx, userID, guildID, matDefID, qty); err != nil {
		return "", fmt.Errorf("phân giải thành công nhưng túi đồ đầy, không thể chứa nguyên liệu")
	}

	return fmt.Sprintf("Phân giải **%s** thành công! Nhận được **%d** Hắc Thiết.", def.Name, qty), nil
}
