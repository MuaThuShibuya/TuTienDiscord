// File: internal/game/cultivation/model.go
// Version: v0.1
// Purpose: Define the CultivationProfile model stored in "cultivation_profiles" collection.
// Security: userId and guildId are always paired in queries to prevent cross-guild leaks.
// Notes: Realm and RealmLevel drive progression; combatPower is derived, not raw user input.

package cultivation

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Realm represents a cultivation stage (cảnh giới).
type Realm string

const (
	RealmMortual    Realm = "mortal"       // Phàm Nhân
	RealmQiRefining Realm = "qi_refining"  // Luyện Khí
	RealmFoundation Realm = "foundation"   // Trúc Cơ
	RealmGoldCore   Realm = "gold_core"    // Kim Đan
	// TODO v0.2+: add higher realms (Nguyên Anh, Hóa Thần, ...)
)

// MindState represents the cultivator's mental state (tâm cảnh).
type MindState string

const (
	MindStateCalm      MindState = "calm"       // Bình Tĩnh
	MindStateFocused   MindState = "focused"    // Chuyên Tâm
	MindStateUnstable  MindState = "unstable"   // Bất Ổn
	MindStateEnlighten MindState = "enlightened" // Ngộ Đạo
)

// CultivationPath represents the chosen path (đạo lộ).
type CultivationPath string

const (
	PathNone   CultivationPath = ""        // Chưa chọn
	PathSword  CultivationPath = "sword"   // Kiếm Tu
	PathBody   CultivationPath = "body"    // Thể Tu
	PathSpirit CultivationPath = "spirit"  // Linh Tu
	PathPoison CultivationPath = "poison"  // Độc Tu
	// TODO v0.2+: add more paths
)

// CultivationProfile holds all cultivation-related data for a player.
type CultivationProfile struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"           json:"id"`
	UserID                string             `bson:"userId"                  json:"userId"`
	GuildID               string             `bson:"guildId"                 json:"guildId"`
	Realm                 Realm              `bson:"realm"                   json:"realm"`
	RealmLevel            int                `bson:"realmLevel"              json:"realmLevel"`            // Level within current realm (1-9)
	CultivationExp        int64              `bson:"cultivationExp"          json:"cultivationExp"`
	CultivationExpRequired int64             `bson:"cultivationExpRequired"  json:"cultivationExpRequired"`
	CombatPower           int64              `bson:"combatPower"             json:"combatPower"`
	Stamina               int                `bson:"stamina"                 json:"stamina"`               // Thể lực hiện tại
	MaxStamina            int                `bson:"maxStamina"              json:"maxStamina"`
	MindState             MindState          `bson:"mindState"               json:"mindState"`
	Path                  CultivationPath    `bson:"path"                    json:"path"`
	CreatedAt             time.Time          `bson:"createdAt"               json:"createdAt"`
	UpdatedAt             time.Time          `bson:"updatedAt"               json:"updatedAt"`
}

// RealmDisplayName returns the Vietnamese name for the realm.
func (r Realm) DisplayName() string {
	names := map[Realm]string{
		RealmMortual:    "Phàm Nhân",
		RealmQiRefining: "Luyện Khí",
		RealmFoundation: "Trúc Cơ",
		RealmGoldCore:   "Kim Đan",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return "Không Rõ"
}

// MindStateDisplayName returns the Vietnamese name for the mind state.
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

// DefaultCultivationProfile returns a new profile with starting values.
func DefaultCultivationProfile(userID, guildID string) *CultivationProfile {
	return &CultivationProfile{
		UserID:                 userID,
		GuildID:                guildID,
		Realm:                  RealmQiRefining,
		RealmLevel:             1,
		CultivationExp:         0,
		CultivationExpRequired: 1000,
		CombatPower:            100,
		Stamina:                100,
		MaxStamina:             100,
		MindState:              MindStateCalm,
		Path:                   PathNone,
	}
}
