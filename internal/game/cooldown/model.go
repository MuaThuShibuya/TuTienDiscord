// File: internal/game/cooldown/model.go
// Phiên bản: v0.1.1
// Mục đích: Struct Cooldown và các hằng Action cho hệ thống cooldown.
// Ghi chú: MongoDB TTL index trên expiresAt tự động xóa cooldown đã hết hạn.

package cooldown

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Action tên hành động được cooldown — luôn là hằng server-side, không nhận từ user.
type Action string

const (
	ActionCultivate    Action = "cultivate"     // Tĩnh tu
	ActionMeditate     Action = "meditate"      // Tĩnh tu
	ActionSeclusion    Action = "seclusion"     // Bế quan
	ActionBodyTraining Action = "body_training" // Luyện thể
	ActionBreakthrough Action = "breakthrough"  // Đột phá
	ActionDungeon      Action = "dungeon"       // Phó bản
	ActionDaily        Action = "daily"         // Điểm danh hàng ngày
	ActionGacha        Action = "gacha"         // Quay cơ duyên
	ActionPvP          Action = "pvp"           // PvP
	ActionBoss         Action = "boss"          // Boss server
	// TODO v0.2+: thêm action mới khi xây dựng các hệ thống tương ứng
)

// Cooldown một bản ghi cooldown per (userId, guildId, action).
type Cooldown struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId"        json:"userId"`
	GuildID   string             `bson:"guildId"       json:"guildId"`
	Action    Action             `bson:"action"        json:"action"`
	ExpiresAt time.Time          `bson:"expiresAt"     json:"expiresAt"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
}

// IsExpired trả về true nếu cooldown đã hết hạn.
func (c *Cooldown) IsExpired() bool {
	return time.Now().UTC().After(c.ExpiresAt)
}

// Remaining trả về thời gian còn lại của cooldown. 0 nếu đã hết.
func (c *Cooldown) Remaining() time.Duration {
	r := time.Until(c.ExpiresAt)
	if r < 0 {
		return 0
	}
	return r
}
