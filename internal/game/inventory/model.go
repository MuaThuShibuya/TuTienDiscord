// File: internal/game/inventory/model.go
package inventory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	UserID         string             `bson:"userId"`
	GuildID        string             `bson:"guildId"`
	SlotLimit      int                `bson:"slotLimit"`
	StarterGranted bool               `bson:"starterGranted"`
	CreatedAt      time.Time          `bson:"createdAt"`
	UpdatedAt      time.Time          `bson:"updatedAt"`
}
