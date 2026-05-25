// File: internal/game/equipment/repository.go
package equipment

import "context"

type Repository interface {
	Get(ctx context.Context, userID, guildID string) (*EquipmentSet, error)
	Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error
	Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error
}
