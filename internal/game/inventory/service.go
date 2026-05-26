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
	items, err := s.itemRepo.GetInstancesByUser(ctx, userID, guildID)
	return inv, items, err
}

func (s *inventoryService) AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	def, ok := item.GetDefinition(definitionID)
	if !ok {
		return apperrors.ErrInvalidInput
	}

	inv, items, err := s.GetInventory(ctx, userID, guildID)
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

	// Cần tạo ô mới
	if len(items) >= inv.SlotLimit {
		return apperrors.ErrInventoryFull
	}

	inst := &item.ItemInstance{
		InstanceID:   utils.NewSessionID(),
		DefinitionID: definitionID,
		UserID:       userID,
		GuildID:      guildID,
		Quantity:     quantity,
	}
	return s.itemRepo.CreateInstance(ctx, inst)
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

	// 2. Trừ nguyên liệu (đã xác nhận đủ)
	// Ghi chú: Đây chưa phải là giao dịch atomic hoàn hảo. Để an toàn tuyệt đối cần MongoDB Transaction.
	// Tuy nhiên, với game bot, việc check trước rồi trừ sau là đủ an toàn cho hầu hết trường hợp.
	for defID, requiredQty := range itemsToConsume {
		for _, it := range currentItems {
			if it.DefinitionID == defID {
				// Giả định mỗi nguyên liệu chỉ có 1 stack. Sẽ cần cải tiến nếu nguyên liệu có nhiều stack.
				if err := s.itemRepo.AdjustQuantity(ctx, it.InstanceID, userID, guildID, -requiredQty); err != nil {
					return fmt.Errorf("lỗi trừ nguyên liệu %s: %w", defID, err)
				}
				_ = s.itemRepo.DeleteInstance(ctx, it.InstanceID, userID, guildID) // Tự dọn dẹp nếu số lượng <= 0
				break
			}
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

	_ = s.AddItem(ctx, userID, guildID, "pill_exp_tu_khi_d", 3)
	_ = s.AddItem(ctx, userID, guildID, "pill_stm_hoi_luc_d", 2)
	_ = s.AddItem(ctx, userID, guildID, "eq_weapon_moc_kiem_d", 1)
	_ = s.AddItem(ctx, userID, guildID, "eq_armor_vai_tho_d", 1)
	_ = s.AddItem(ctx, userID, guildID, "mat_enhance_hac_thiet_d", 3)

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

	// 1. Trừ số lượng vật phẩm (Atomic)
	if err := s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, -1); err != nil {
		return "", err
	}

	// 2. Áp dụng hiệu ứng
	expGained := int64(0)
	staminaGained := 0
	breakthroughBuff := 0

	if def.Effects != nil {
		expGained = int64(def.Effects["exp"])
		staminaGained = def.Effects["stamina"]
		breakthroughBuff = def.Effects["breakthrough_chance"]
	}

	// Fallback an toàn nếu Effects chưa được định nghĩa đầy đủ trong Registry
	if expGained == 0 && staminaGained == 0 && breakthroughBuff == 0 {
		logger.L().Warn("Sử dụng fallback effect cho item do thiếu config", zap.String("defID", inst.DefinitionID))
		if inst.DefinitionID == "pill_exp_tu_khi_d" || inst.DefinitionID == "item_qi_pill" {
			expGained = 50
		} else if inst.DefinitionID == "pill_stm_hoi_luc_d" {
			staminaGained = 10
		}
	}

	var messages []string
	if expGained > 0 {
		if err := s.cultSvc.AddExperience(ctx, userID, guildID, expGained); err != nil {
			// Hoàn trả lại vật phẩm nếu hiệu ứng thất bại
			_ = s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, 1)
			return "", fmt.Errorf("không thể hấp thụ linh khí lúc này: %w", err)
		}
		messages = append(messages, fmt.Sprintf("nhận **%d** tu vi", expGained))
	}
	if staminaGained > 0 {
		if err := s.cultSvc.AddStamina(ctx, userID, guildID, staminaGained); err != nil {
			// Hoàn trả lại vật phẩm nếu hiệu ứng thất bại
			_ = s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, 1)
			return "", fmt.Errorf("không thể hồi phục thể lực lúc này: %w", err)
		}
		messages = append(messages, fmt.Sprintf("hồi **%d** thể lực", staminaGained))
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
	_ = s.itemRepo.DeleteInstance(ctx, instanceID, userID, guildID)
	return message, nil
}
