// File: internal/game/equipment/model.go
package equipment

import (
	"strings"
	"time"
)

type EquipmentSlot string

const (
	SlotWeapon   EquipmentSlot = "weapon"
	SlotArmor    EquipmentSlot = "armor"
	SlotArtifact EquipmentSlot = "artifact"
	SlotTreasure EquipmentSlot = "treasure"
	SlotBoots    EquipmentSlot = "boots"
)

type EquipmentSet struct {
	UserID    string            `bson:"userId"`
	GuildID   string            `bson:"guildId"`
	Slots     map[string]string `bson:"slots"`
	CreatedAt time.Time         `bson:"createdAt"`
	UpdatedAt time.Time         `bson:"updatedAt"`
}

func GetSlotForDefinition(defID string) EquipmentSlot {
	if strings.HasPrefix(defID, "eq_weapon_") {
		return SlotWeapon
	}
	if strings.HasPrefix(defID, "eq_armor_") {
		return SlotArmor
	}
	if strings.HasPrefix(defID, "eq_boots_") {
		return SlotBoots
	}
	if strings.HasPrefix(defID, "eq_artifact_") {
		return SlotArtifact
	}
	if strings.HasPrefix(defID, "eq_treasure_") {
		return SlotTreasure
	}
	return ""
}
