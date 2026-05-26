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
	TypeFurnace   ItemType = "furnace"
)

type ItemRarity string

const (
	RarityD    ItemRarity = "D"
	RarityC    ItemRarity = "C"
	RarityB    ItemRarity = "B"
	RarityA    ItemRarity = "A"
	RarityS    ItemRarity = "S"
	RaritySS   ItemRarity = "SS"
	RaritySSS  ItemRarity = "SSS"
	RaritySSSP ItemRarity = "SSS+"
)

type ItemDefinition struct {
	ID            string
	Name          string
	Type          ItemType
	Rarity        ItemRarity
	Stackable     bool
	MaxStack      int
	Usable        bool
	Description   string
	Effects       map[string]int // Dùng cho đan dược (buff exp, stamina...)
	Stats         map[string]int // Dùng cho trang bị & lò đan (attack, defense...)
	RequiredRealm string         // Cảnh giới yêu cầu (nếu có)
	SellPrice     int64          // Giá bán ra linh thạch
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
