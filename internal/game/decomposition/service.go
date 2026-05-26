package decomposition

import (
	"context"
	"errors"
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"go.uber.org/zap"
)

type Service interface {
	Decompose(ctx context.Context, userID, guildID, instanceID string) ([]string, error)
}

type service struct {
	itemRepo item.Repository
	invSvc   inventory.Service
	equipSvc equipment.Service
	log      *zap.Logger
}

func NewService(itemRepo item.Repository, invSvc inventory.Service, equipSvc equipment.Service) Service {
	return &service{
		itemRepo: itemRepo,
		invSvc:   invSvc,
		equipSvc: equipSvc,
		log:      logger.L().Named("game.decomposition"),
	}
}

// Bảng Rarity Map
var decompMaterials = map[string][]string{
	"D":    {"mat_scrap_linh_tran_d", "mat_scrap_thiet_vun_d"},
	"C":    {"mat_scrap_thanh_dong_c", "mat_scrap_moc_linh_c"},
	"B":    {"mat_scrap_huyen_thiet_b", "mat_scrap_tu_van_b"},
	"A":    {"mat_scrap_bach_ngoc_a", "mat_scrap_linh_hoa_a"},
	"S":    {"mat_scrap_long_van_s", "mat_scrap_phuong_hoa_s"},
	"SS":   {"mat_scrap_tinh_ha_ss", "mat_scrap_van_phap_ss"},
	"SSS":  {"mat_scrap_hu_khong_sss"},
	"SSS+": {"mat_scrap_luan_hoi_sssp", "mat_scrap_nghich_menh_sssp"},
}

func (s *service) Decompose(ctx context.Context, userID, guildID, instanceID string) ([]string, error) {
	inst, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return nil, errors.New("vật phẩm không tồn tại hoặc không thuộc về bạn")
	}

	if locked, ok := inst.Metadata["locked"].(bool); ok && locked {
		return nil, errors.New("không thể phân giải vật phẩm đang bị khóa")
	}

	eqSet, err := s.equipSvc.GetEquipment(ctx, userID, guildID)
	if err == nil && eqSet != nil && eqSet.Slots != nil {
		// Sửa lỗi duyệt map slots thay vì struct properties
		for _, equippedID := range eqSet.Slots {
			if equippedID == instanceID {
				return nil, errors.New("không thể phân giải trang bị đang mặc, vui lòng tháo ra trước")
			}
		}
	}

	// Sửa lỗi lấy Rarity thông qua Definition
	def, ok := item.GetDefinition(inst.DefinitionID)
	if !ok {
		return nil, errors.New("không tìm thấy định nghĩa vật phẩm")
	}
	rarity := string(def.Rarity)

	materials, ok := decompMaterials[rarity]
	if !ok || len(materials) == 0 {
		return nil, fmt.Errorf("không tìm thấy công thức phân giải cho phẩm chất %s", rarity)
	}

	// Sửa lỗi Xóa vật phẩm: AdjustQuantity trước để đảm bảo quantity <= 0 rồi mới Delete
	if err := s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, -1); err != nil {
		s.log.Error("Trừ item phân giải thất bại", zap.Error(err), zap.String("instanceID", instanceID))
		return nil, errors.New("lỗi hệ thống khi phân giải vật phẩm")
	}
	_ = s.itemRepo.DeleteInstance(ctx, instanceID, userID, guildID)

	var addedMats []string
	for _, matID := range materials {
		if err := s.invSvc.AddItem(ctx, userID, guildID, matID, 1); err == nil {
			addedMats = append(addedMats, matID)
		}
	}

	s.log.Info("Phân giải thành công", zap.String("userID", userID), zap.String("instanceID", instanceID))
	return addedMats, nil
}
