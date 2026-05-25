// File: internal/game/cultivation/defaults.go
// Phiên bản: v0.1.1
// Mục đích: Giá trị khởi đầu cho hồ sơ tu luyện của người chơi mới.
// Ghi chú: Chỉnh sửa các hằng số ở đây để điều chỉnh game balance mà không ảnh hưởng struct.

package cultivation

import "go.mongodb.org/mongo-driver/bson/primitive"

// Giá trị khởi đầu cho người chơi mới tạo tài khoản.
const (
	DefaultRealm          = RealmQiRefining // Cảnh giới bắt đầu: Luyện Khí
	DefaultRealmLevel     = 1
	DefaultCultivationExp = int64(0)
	DefaultExpRequired    = int64(1000) // Tu vi cần để lên tầng tiếp theo
	DefaultCombatPower    = int64(100)
	DefaultStamina        = 100
	DefaultMaxStamina     = 100
	DefaultMindState      = MindStateCalm
	DefaultPath           = PathNone
)

// NewCultivationProfile tạo hồ sơ tu luyện mới với giá trị khởi đầu.
func NewCultivationProfile(userID, guildID string) *CultivationProfile {
	return &CultivationProfile{
		ID:                     primitive.NewObjectID(),
		UserID:                 userID,
		GuildID:                guildID,
		Realm:                  DefaultRealm,
		RealmLevel:             DefaultRealmLevel,
		CultivationExp:         DefaultCultivationExp,
		CultivationExpRequired: DefaultExpRequired,
		CombatPower:            DefaultCombatPower,
		Stamina:                DefaultStamina,
		MaxStamina:             DefaultMaxStamina,
		MindState:              DefaultMindState,
		Path:                   DefaultPath,
	}
}
