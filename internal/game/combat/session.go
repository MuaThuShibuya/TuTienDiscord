// File: internal/game/combat/session.go
// Chức năng: Quản lý trạng thái và dữ liệu của một phiên chiến đấu PvE (CombatSession), lưu trữ trên MongoDB và chống spam double-click.

package combat

import "time"

type SessionState string

const (
	StateActive  SessionState = "active"
	StateWon     SessionState = "won"
	StateLost    SessionState = "lost"
	StateEscaped SessionState = "escaped"
	StateExpired SessionState = "expired"
)

// CombatLogEntry lưu lịch sử hành động để render lên Discord Embed.
type CombatLogEntry struct {
	Turn      int       `bson:"turn" json:"turn"`
	ActorID   string    `bson:"actorId" json:"actorId"`
	Action    string    `bson:"action" json:"action"`
	Message   string    `bson:"message" json:"message"`
	Damage    int64     `bson:"damage" json:"damage"`
	IsCrit    bool      `bson:"isCrit" json:"isCrit"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// CombatSession lưu toàn bộ state của trận đấu (Lưu MongoDB).
type CombatSession struct {
	ID              string           `bson:"_id" json:"id"`
	UserID          string           `bson:"userId" json:"userId"` // Chủ phòng
	ZoneID          string           `bson:"zoneId" json:"zoneId"`
	Stage           int              `bson:"stage" json:"stage"`
	State           SessionState     `bson:"state" json:"state"`
	Turn            int              `bson:"turn" json:"turn"`
	Player          CombatActor      `bson:"player" json:"player"`
	Enemies         []CombatActor    `bson:"enemies" json:"enemies"`
	Logs            []CombatLogEntry `bson:"logs" json:"logs"`
	IdempotencyKeys []string         `bson:"idempotencyKeys" json:"idempotencyKeys"` // Chống double-click
	AutoBattle      AutoBattlePolicy `bson:"autoBattle" json:"autoBattle"`           // Chính sách đánh tự động
	CreatedAt       time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time        `bson:"updatedAt" json:"updatedAt"`
	ExpiresAt       time.Time        `bson:"expiresAt" json:"expiresAt"` // MongoDB TTL index sẽ dọn dẹp
}
