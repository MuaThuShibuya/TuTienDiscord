// File: internal/game/skill/equipped.go
// Chức năng: Quản lý 3 Slot trang bị kỹ năng của một thực thể (Player/Pet).

package skill

import (
	"context"
	"fmt"
	"time"
)

type EquippedSkillSlot struct {
	SlotIndex int    `bson:"slotIndex" json:"slotIndex"`
	SkillID   string `bson:"skillId" json:"skillId"`
	Unlocked  bool   `bson:"unlocked" json:"unlocked"`
}

type EquippedSkillSet struct {
	OwnerID   string               `bson:"ownerId" json:"ownerId"`
	OwnerType string               `bson:"ownerType" json:"ownerType"`
	Slots     [3]EquippedSkillSlot `bson:"slots" json:"slots"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

// NewEquippedSkillSet khởi tạo bộ kỹ năng mặc định với 3 slot đã mở khóa.
func NewEquippedSkillSet(ownerType, ownerID string) *EquippedSkillSet {
	return &EquippedSkillSet{
		OwnerID:   ownerID,
		OwnerType: ownerType,
		Slots: [3]EquippedSkillSlot{
			{SlotIndex: 0, Unlocked: true},
			{SlotIndex: 1, Unlocked: true},
			{SlotIndex: 2, Unlocked: true},
		},
		UpdatedAt: time.Now().UTC(),
	}
}

// SetSlot gài kỹ năng vào slot hợp lệ.
func (e *EquippedSkillSet) SetSlot(index int, skillID string) error {
	if index < 0 || index > 2 {
		return fmt.Errorf("slot không hợp lệ, chỉ từ 0 đến 2")
	}
	if !e.Slots[index].Unlocked {
		return fmt.Errorf("slot %d chưa được mở khóa", index)
	}
	e.Slots[index].SlotIndex = index
	e.Slots[index].SkillID = skillID
	e.UpdatedAt = time.Now().UTC()
	return nil
}

type EquippedSkillService interface {
	GetEquippedSkills(ctx context.Context, ownerType, ownerID string) (*EquippedSkillSet, error)
	EquipSkill(ctx context.Context, ownerType, ownerID string, slot int, skillID string) error
	UnequipSkill(ctx context.Context, ownerType, ownerID string, slot int) error
}
