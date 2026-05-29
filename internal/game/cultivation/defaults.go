// File: internal/game/cultivation/defaults.go
// Phiên bản: v0.1.1
// Mục đích: Giá trị khởi đầu cho hồ sơ tu luyện của người chơi mới.
// Ghi chú: Chỉnh sửa các hằng số ở đây để điều chỉnh game balance mà không ảnh hưởng struct.

package cultivation

import "go.mongodb.org/mongo-driver/bson/primitive"

// Giá trị khởi đầu cho người chơi mới tạo tài khoản.
const (
	DefaultRealm          = Realm("pham_nhan") // Cảnh giới bắt đầu ID chuẩn
	DefaultRealmLevel     = 1
	DefaultCultivationExp = int64(0)
	DefaultExpRequired    = int64(200) // Khớp với CalculateNextExpRequired
	DefaultCombatPower    = int64(100)
	DefaultMindState      = 50 // Khởi đầu v0.2: Tâm cảnh Bình tĩnh (50/100)
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
		MindState:              DefaultMindState,
		Path:                   DefaultPath,
	}
}
