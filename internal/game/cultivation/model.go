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

// MindState tâm cảnh tu luyện.
type MindState string

const (
	MindStateCalm      MindState = "calm"        // Bình Tĩnh
	MindStateFocused   MindState = "focused"     // Chuyên Tâm
	MindStateUnstable  MindState = "unstable"    // Bất Ổn
	MindStateEnlighten MindState = "enlightened" // Ngộ Đạo
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
	MindState              MindState          `bson:"mindState"              json:"mindState"`
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
	return "Không Rõ"
}

// MindStateDisplayName trả về tên tiếng Việt của tâm cảnh.
func (m MindState) DisplayName() string {
	names := map[MindState]string{
		MindStateCalm:      "Bình Tĩnh",
		MindStateFocused:   "Chuyên Tâm",
		MindStateUnstable:  "Bất Ổn",
		MindStateEnlighten: "Ngộ Đạo",
	}
	if name, ok := names[m]; ok {
		return name
	}
	return "Không Rõ"
}

// CanBreakthrough kiểm tra xem có thể đột phá không (exp đủ).
func (c *CultivationProfile) CanBreakthrough() bool {
	return c.CultivationExp >= c.CultivationExpRequired
}
