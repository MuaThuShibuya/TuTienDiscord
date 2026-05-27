package equipment

import (
	"context"
	"errors"
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

type CombatStats struct {
	MaxHP       int64
	ATK         int64
	DEF         int64
	CritRate    float64
	CombatPower int64
}

type Service interface {
	GetEquipment(ctx context.Context, userID, guildID string) (*EquipmentSet, error)
	Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error
	Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
	Enhance(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
	GetEffectiveStats(ctx context.Context, userID, guildID string) (CombatStats, error)
}

const DefaultMaxEnhanceLevel = 10

type equipmentService struct {
	repo     Repository
	itemRepo item.Repository
	invSvc   inventory.Service
}

func NewService(repo Repository, itemRepo item.Repository, invSvc inventory.Service) Service {
	return &equipmentService{repo: repo, itemRepo: itemRepo, invSvc: invSvc}
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
	eq, err := s.repo.Get(ctx, userID, guildID)
	if err != nil {
		return err
	}

	instID, ok := eq.Slots[string(slot)]
	if !ok || instID == "" {
		return errors.New("không có trang bị ở vị trí này để cường hóa")
	}

	inst, err := s.itemRepo.GetInstanceByID(ctx, instID, userID, guildID)
	if err != nil {
		return errors.New("không tìm thấy dữ liệu trang bị")
	}

	def, ok := item.GetDefinition(inst.DefinitionID)
	if !ok {
		return errors.New("không tìm thấy định nghĩa trang bị")
	}

	level := 0
	if inst.Metadata != nil && inst.Metadata["level"] != nil {
		if l, ok := inst.Metadata["level"].(int32); ok {
			level = int(l)
		} else if l, ok := inst.Metadata["level"].(float64); ok {
			level = int(l)
		}
		if l, ok := inst.Metadata["level"].(int); ok {
			level = l
		}
	} else {
		inst.Metadata = make(map[string]interface{})
	}

	maxLevel := def.MaxEnhanceLevel
	if maxLevel <= 0 {
		maxLevel = DefaultMaxEnhanceLevel
	}

	if level >= maxLevel {
		return fmt.Errorf("trang bị đã đạt cấp cường hóa tối đa")
	}

	cost := int64(level + 1)
	if err := s.invSvc.ConsumeItems(ctx, userID, guildID, map[string]int64{"mat_enhance_hac_thiet_d": cost}); err != nil {
		return fmt.Errorf("không đủ Hắc Thiết. Cần **%d** để thăng cấp", cost)
	}

	inst.Metadata["level"] = level + 1
	return s.itemRepo.UpdateMetadata(ctx, inst.InstanceID, userID, guildID, inst.Metadata)
}

func (s *equipmentService) GetEffectiveStats(ctx context.Context, userID, guildID string) (CombatStats, error) {
	var stats CombatStats
	eq, err := s.repo.Get(ctx, userID, guildID)
	if err != nil {
		return stats, err
	}

	for _, instID := range eq.Slots {
		inst, err := s.itemRepo.GetInstanceByID(ctx, instID, userID, guildID)
		if err != nil {
			continue
		}
		def, ok := item.GetDefinition(inst.DefinitionID)
		if !ok || def.Stats == nil {
			continue
		}

		level := 0
		if inst.Metadata != nil && inst.Metadata["level"] != nil {
			if l, ok := inst.Metadata["level"].(int32); ok {
				level = int(l)
			} else if l, ok := inst.Metadata["level"].(float64); ok {
				level = int(l)
			}
		}

		// Cường hóa: buff 10% chỉ số mỗi cấp
		multiplier := 1.0 + (float64(level) * 0.1)

		stats.MaxHP += int64(float64(def.Stats["hp"]) * multiplier)
		stats.ATK += int64(float64(def.Stats["atk"]) * multiplier)
		stats.DEF += int64(float64(def.Stats["def"]) * multiplier)

		// Crit thường cố định không tăng theo % multiplier, chỉ cộng dồn thẳng
		if critRate, exists := def.Stats["crit"]; exists {
			// Giả định config ghi crit là 15 -> Tương đương 15% (0.15)
			stats.CritRate += float64(critRate) / 100.0
		}
	}

	stats.CombatPower = stats.ATK*2 + stats.DEF*2 + stats.MaxHP/10

	return stats, nil
}
