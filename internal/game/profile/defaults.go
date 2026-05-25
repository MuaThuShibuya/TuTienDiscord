// File: internal/game/profile/defaults.go
// Phiên bản: v0.1.1
// Mục đích: Giá trị mặc định khi khởi tạo Player mới và hằng số giới hạn đạo hiệu.
// Ghi chú: Tách ra khỏi model.go để dễ chỉnh sửa game balance mà không ảnh hưởng struct.

package profile

import "go.mongodb.org/mongo-driver/bson/primitive"

// MaxDaoNameLength là số ký tự (rune) tối đa cho đạo hiệu.
// Giới hạn này tránh đạo hiệu quá dài làm vỡ layout embed Discord.
const MaxDaoNameLength = 24

// MinDaoNameLength là số ký tự (rune) tối thiểu cho đạo hiệu.
// Đạo hiệu quá ngắn (1 ký tự) trông thiếu ý nghĩa và dễ gây nhầm lẫn.
const MinDaoNameLength = 2

// NewPlayer tạo Player mới với giá trị khởi đầu cho người chơi mới.
// daoName là đạo hiệu mặc định, có thể thay đổi sau.
func NewPlayer(userID, guildID, username, displayName, daoName string) *Player {
	return &Player{
		ID:          primitive.NewObjectID(),
		UserID:      userID,
		GuildID:     guildID,
		Username:    username,
		DisplayName: displayName,
		DaoName:     daoName,
		Status:      StatusActive,
	}
}
