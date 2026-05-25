// File: internal/game/cooldown/model.go
// Version: v0.1
// Purpose: Define the Cooldown model for tracking per-user, per-action cooldowns.
// Security: action name comes from server-side constants, never raw user input.
// Notes: MongoDB TTL index on expiresAt auto-deletes expired cooldowns. See database/indexes.go.

package cooldown

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Action defines a named cooldown type. Use constants to avoid typos.
type Action string

const (
	ActionCultivate  Action = "cultivate"   // Tĩnh tu
	ActionDungeon    Action = "dungeon"      // Phó bản
	ActionDaily      Action = "daily"        // Điểm danh hàng ngày
	ActionGacha      Action = "gacha"        // Quay cơ duyên — rate limited
	ActionPvP        Action = "pvp"          // PvP
	ActionBoss       Action = "boss"         // Boss server
	// TODO v0.2+: add more actions as systems are built
)

// Cooldown is one record per (userId, guildId, action) triple.
type Cooldown struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId"        json:"userId"`
	GuildID   string             `bson:"guildId"       json:"guildId"`
	Action    Action             `bson:"action"        json:"action"`
	ExpiresAt time.Time          `bson:"expiresAt"     json:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
}

// IsExpired returns true if the cooldown has passed.
func (c *Cooldown) IsExpired() bool {
	return time.Now().UTC().After(c.ExpiresAt)
}

// RemainingDuration returns how long until the cooldown expires.
func (c *Cooldown) RemainingDuration() time.Duration {
	remaining := time.Until(c.ExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}
