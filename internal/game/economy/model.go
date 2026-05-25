// File: internal/game/economy/model.go
// Version: v0.1
// Purpose: Define the Wallet model — in-game currencies (linh thạch, linh ngọc, vé cơ duyên).
// Security: Currency values must only be modified through the service using atomic DB operations.
//           Never accept currency amounts directly from Discord user input.
// Notes: All currencies are in-game only. No real-money transactions. Anti-race-condition
//        updates will use MongoDB $inc with findOneAndUpdate in v0.2+.

package economy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallet holds all in-game currency for a player.
type Wallet struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"  json:"id"`
	UserID       string             `bson:"userId"         json:"userId"`
	GuildID      string             `bson:"guildId"        json:"guildId"`
	SpiritStones int64              `bson:"spiritStones"   json:"spiritStones"`  // Linh Thạch — common currency
	SpiritJades  int64              `bson:"spiritJades"    json:"spiritJades"`   // Linh Ngọc — premium in-game currency
	FateTickets  int                `bson:"fateTickets"    json:"fateTickets"`   // Vé Cơ Duyên — gacha tickets
	CreatedAt    time.Time          `bson:"createdAt"      json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"      json:"updatedAt"`
}

// DefaultWallet returns a wallet with starting currencies for a new player.
func DefaultWallet(userID, guildID string) *Wallet {
	return &Wallet{
		UserID:       userID,
		GuildID:      guildID,
		SpiritStones: 500,   // Starting linh thạch
		SpiritJades:  0,
		FateTickets:  3,     // 3 free pulls to start
	}
}
