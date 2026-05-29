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

// ClaimedReward lưu thông tin vật phẩm đã nhận (dạng nguyên thủy, độc lập với pve package)
type ClaimedReward struct {
	Type     string `bson:"type" json:"type"`
	RefID    string `bson:"refId" json:"refId"`
	Quantity int64  `bson:"quantity" json:"quantity"`
	IsBonus  bool   `bson:"isBonus" json:"isBonus"`
}

// CombatSession lưu toàn bộ state của trận đấu (Lưu MongoDB).
type CombatSession struct {
	ID                     string           `bson:"_id" json:"id"`
	UserID                 string           `bson:"userId" json:"userId"`   // Chủ phòng
	GuildID                string           `bson:"guildId" json:"guildId"` // Nơi xảy ra trận chiến
	AreaID                 string           `bson:"areaId" json:"areaId"`
	ActivityType           string           `bson:"activityType" json:"activityType"`
	Stage                  int              `bson:"stage" json:"stage"`
	State                  SessionState     `bson:"state" json:"state"`
	Turn                   int              `bson:"turn" json:"turn"`
	Player                 CombatActor      `bson:"player" json:"player"`
	Enemies                []CombatActor    `bson:"enemies" json:"enemies"`
	CurrentActorID         string           `bson:"currentActorId" json:"currentActorId"`
	TurnOrder              []TurnOrderEntry `bson:"turnOrder" json:"turnOrder"`
	GuaranteedRewardPoolID string           `bson:"guaranteedRewardPoolId" json:"guaranteedRewardPoolId"`
	BonusRewardPoolID      string           `bson:"bonusRewardPoolId" json:"bonusRewardPoolId"`
	Logs                   []CombatLogEntry `bson:"logs" json:"logs"`
	IdempotencyKeys        []string         `bson:"idempotencyKeys" json:"idempotencyKeys"` // Chống double-click
	AutoBattle             AutoBattlePolicy `bson:"autoBattle" json:"autoBattle"`           // Chính sách đánh tự động
	RewardClaimed          bool             `bson:"rewardClaimed" json:"rewardClaimed"`
	RewardClaimStatus      string           `bson:"rewardClaimStatus,omitempty" json:"rewardClaimStatus,omitempty"`
	RewardClaimID          string           `bson:"rewardClaimId,omitempty" json:"rewardClaimId,omitempty"`
	RewardClaimError       string           `bson:"rewardClaimError,omitempty" json:"rewardClaimError,omitempty"`
	RewardClaimStartedAt   time.Time        `bson:"rewardClaimStartedAt,omitempty" json:"rewardClaimStartedAt,omitempty"`
	RewardClaimedAt        time.Time        `bson:"rewardClaimedAt" json:"rewardClaimedAt"`
	ClaimedRewards         []ClaimedReward  `bson:"claimedRewards" json:"claimedRewards"`
	CreatedAt              time.Time        `bson:"createdAt" json:"createdAt"`
	UpdatedAt              time.Time        `bson:"updatedAt" json:"updatedAt"`
	ExpiresAt              time.Time        `bson:"expiresAt" json:"expiresAt"` // MongoDB TTL index sẽ dọn dẹp
}
