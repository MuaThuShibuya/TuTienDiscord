// File: internal/game/profile/model.go
// Version: v0.1
// Purpose: Define the Player data model stored in the "players" MongoDB collection.
// Security: userId and guildId are always validated before any DB operation.
// Notes: daoName (đạo hiệu) is user-chosen and must be sanitized before storage.

package profile

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlayerStatus represents the account status of a player.
type PlayerStatus string

const (
	StatusActive  PlayerStatus = "active"
	StatusBanned  PlayerStatus = "banned"
	StatusInactive PlayerStatus = "inactive"
)

// Player is the core user record. One record per (userId, guildId) pair.
type Player struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"    json:"id"`
	UserID       string             `bson:"userId"           json:"userId"`
	GuildID      string             `bson:"guildId"          json:"guildId"`
	Username     string             `bson:"username"         json:"username"`     // Discord username
	DisplayName  string             `bson:"displayName"      json:"displayName"`  // Discord display name
	DaoName      string             `bson:"daoName"          json:"daoName"`      // đạo hiệu (player-chosen)
	Status       PlayerStatus       `bson:"status"           json:"status"`
	CreatedAt    time.Time          `bson:"createdAt"        json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"        json:"updatedAt"`
	LastActiveAt time.Time          `bson:"lastActiveAt"     json:"lastActiveAt"`
}

// IsActive returns true if the player account is not banned or deactivated.
func (p *Player) IsActive() bool {
	return p.Status == StatusActive
}
