// File: internal/game/cultivation/model.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa struct CultivationProfile — dữ liệu tu luyện của người chơi.
// Ghi chú: Giá trị khởi đầu nằm trong defaults.go. Thêm cảnh giới mới trong Realm constants.

package cultivation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Realm cảnh giới tu tiên.
type Realm string

const (
	RealmQiRefining Realm = "qi_refining" // Luyện Khí
	RealmFoundation Realm = "foundation"  // Trúc Cơ
	RealmGoldCore   Realm = "gold_core"   // Kim Đan
	// TODO v0.2+: thêm cảnh giới cao hơn
)

// CultivationPath đạo lộ tu luyện.
type CultivationPath string

const (
	PathNone   CultivationPath = ""       // Chưa chọn
	PathSword  CultivationPath = "sword"  // Kiếm Tu
	PathBody   CultivationPath = "body"   // Thể Tu
	PathSpirit CultivationPath = "spirit" // Linh Tu
	PathPoison CultivationPath = "poison" // Độc Tu
)

// CultivationProfile hồ sơ tu luyện của người chơi.
type CultivationProfile struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"          json:"id"`
	UserID                 string             `bson:"userId"                 json:"userId"`
	GuildID                string             `bson:"guildId"                json:"guildId"`
	Realm                  Realm              `bson:"realm"                  json:"realm"`
	RealmLevel             int                `bson:"realmLevel"             json:"realmLevel"` // Tầng trong cảnh giới (1-9)
	CultivationExp         int64              `bson:"cultivationExp"         json:"cultivationExp"`
	CultivationExpRequired int64              `bson:"cultivationExpRequired" json:"cultivationExpRequired"`
	CombatPower            int64              `bson:"combatPower"            json:"combatPower"`
	Stamina                int                `bson:"stamina"                json:"stamina"` // Thể lực hiện tại
	MaxStamina             int                `bson:"maxStamina"             json:"maxStamina"`
	MindState              int                `bson:"mindState"              json:"mindState"` // 0 - 100
	Path                   CultivationPath    `bson:"path"                   json:"path"`
	CreatedAt              time.Time          `bson:"createdAt"              json:"createdAt"`
	UpdatedAt              time.Time          `bson:"updatedAt"              json:"updatedAt"`
}

// RealmDisplayName trả về tên tiếng Việt của cảnh giới.
func (r Realm) DisplayName() string {
	names := map[Realm]string{
		RealmQiRefining: "Luyện Khí",
		RealmFoundation: "Trúc Cơ",
		RealmGoldCore:   "Kim Đan",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return string(r)
}

// MindStateDisplayName trả về tên tiếng Việt của tâm cảnh.
func (c *CultivationProfile) MindStateDisplayName() string {
	switch {
	case c.MindState >= 80:
		return "Ngộ Đạo"
	case c.MindState >= 50:
		return "Bình Tĩnh"
	case c.MindState >= 20:
		return "Bất Ổn"
	default:
		return "Tâm Ma Xâm Nhập"
	}
}

// CanBreakthrough kiểm tra xem có thể đột phá không (exp đủ).
func (c *CultivationProfile) CanBreakthrough() bool {
	return c.CultivationExp >= c.CultivationExpRequired
}

// --- DTOs cho Service ---

type CultivationActionInput struct {
	UserID  string
	GuildID string
	Now     time.Time
}

type BreakthroughInput struct {
	UserID  string
	GuildID string
	Now     time.Time
	Rand    RandomSource
}

type RandomSource interface {
	Float64() float64
}

type CultivationActionResult struct {
	Action              string
	ExpGained           int64
	CombatPowerGained   int64
	StaminaSpent        int
	NewCultivationExp   int64
	CultivationRequired int64
	NewCombatPower      int64
	NewStamina          int
	NewMindState        int
	CooldownExpiresAt   time.Time
	Message             string
}

type BreakthroughResult struct {
	Success                bool
	OldRealm               string
	OldRealmLevel          int
	NewRealm               string
	NewRealmLevel          int
	AdvancedRealm          bool
	CostPaid               int64
	ExpChanged             int64
	CombatPowerGained      int64
	NewCultivationExp      int64
	NewCultivationRequired int64
	NewMindState           int
	CooldownExpiresAt      time.Time
	Rate                   float64
	Roll                   float64
	Message                string
}
