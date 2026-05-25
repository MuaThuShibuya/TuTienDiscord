// File: internal/game/item/model.go
package item

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemType string

const (
	TypeMaterial  ItemType = "material"
	TypePill      ItemType = "pill"
	TypeEquipment ItemType = "equipment"
)

type ItemRarity string

const (
	RarityD ItemRarity = "D"
	RarityC ItemRarity = "C"
	RarityB ItemRarity = "B"
	RarityA ItemRarity = "A"
	RarityS ItemRarity = "S"
)

type ItemDefinition struct {
	ID        string
	Name      string
	Type      ItemType
	Rarity    ItemRarity
	Stackable bool
	MaxStack  int
	Usable    bool
}

type ItemInstance struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty"`
	InstanceID   string                 `bson:"instanceId"`
	DefinitionID string                 `bson:"definitionId"`
	UserID       string                 `bson:"userId"`
	GuildID      string                 `bson:"guildId"`
	Quantity     int64                  `bson:"quantity"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty"`
	CreatedAt    time.Time              `bson:"createdAt"`
	UpdatedAt    time.Time              `bson:"updatedAt"`
}
