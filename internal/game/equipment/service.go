package equipment

import (
	"context"
	"errors"

	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

type Service interface {
	GetEquipment(ctx context.Context, userID, guildID string) (*EquipmentSet, error)
	Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error
	Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
	Enhance(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
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
	inst, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return errors.New("không tìm thấy trang bị hoặc không thuộc sở hữu")
	}

	defSlot := GetSlotForDefinition(inst.DefinitionID)
	if defSlot != slot {
		return errors.New("trang bị không phù hợp với vị trí này")
	}

	return s.repo.Equip(ctx, userID, guildID, slot, instanceID)
}

func (s *equipmentService) Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error {
	return s.repo.Unequip(ctx, userID, guildID, slot)
}

func (s *equipmentService) Enhance(ctx context.Context, userID, guildID string, slot EquipmentSlot) error {
	return errors.New("tính năng cường hóa đang được hoàn thiện")
}
