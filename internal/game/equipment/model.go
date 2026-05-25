// File: internal/game/equipment/model.go
package equipment

import "time"

type EquipmentSlot string

const (
	SlotWeapon   EquipmentSlot = "weapon"
	SlotArmor    EquipmentSlot = "armor"
	SlotArtifact EquipmentSlot = "artifact"
	SlotTreasure EquipmentSlot = "treasure"
)

type EquipmentSet struct {
	UserID    string            `bson:"userId"`
	GuildID   string            `bson:"guildId"`
	Slots     map[string]string `bson:"slots"`
	CreatedAt time.Time         `bson:"createdAt"`
	UpdatedAt time.Time         `bson:"updatedAt"`
}

func GetSlotForDefinition(defID string) EquipmentSlot {
	switch defID {
	case "eq_wood_sword", "eq_iron_sword":
		return SlotWeapon
	case "eq_cloth_robe", "eq_iron_armor":
		return SlotArmor
	case "eq_spirit_bell":
		return SlotArtifact
	case "eq_guardian_jade":
		return SlotTreasure
	default:
		return ""
	}
}
