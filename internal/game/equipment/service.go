// File: internal/game/equipment/service.go
// Phiên bản: v0.3
// Mục đích: Quản lý trang bị của người chơi. Tách biệt khỏi inventory để dễ tính chỉ số.

package equipment

import (
	"context"

	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

type Service interface {
	GetEquipment(ctx context.Context, userID, guildID string) (*EquipmentSet, error)
	Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error
	Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
}

type equipmentService struct {
	repo     Repository
	itemRepo item.Repository
}

func NewService(repo Repository, itemRepo item.Repository) Service {
	return &equipmentService{repo: repo, itemRepo: itemRepo}
}

func (s *equipmentService) GetEquipment(ctx context.Context, userID, guildID string) (*EquipmentSet, error) {
	return s.repo.Get(ctx, userID, guildID)
}

func (s *equipmentService) Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error {
	// BẢO MẬT: Xác thực vật phẩm tồn tại và thuộc về user này trước khi cho phép mặc
	// Ngăn chặn hacker thay đổi payload custom_id để mặc trang bị của người khác.
	_, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return err // Sẽ trả về lỗi không tìm thấy (ErrNotFound)
	}

	return s.repo.Equip(ctx, userID, guildID, slot, instanceID)
}

func (s *equipmentService) Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error {
	return s.repo.Unequip(ctx, userID, guildID, slot)
}
