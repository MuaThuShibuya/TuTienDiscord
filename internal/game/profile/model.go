// File: internal/game/profile/model.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa struct Player — dữ liệu cốt lõi của người chơi.
// Bảo mật: userId và guildId luôn đi cặp khi truy vấn để tránh rò rỉ dữ liệu chéo server.
// Ghi chú: Giá trị mặc định khi tạo player mới xem trong defaults.go.

package profile

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PlayerStatus trạng thái tài khoản người chơi.
type PlayerStatus string

const (
	StatusActive   PlayerStatus = "active"   // Đang hoạt động
	StatusBanned   PlayerStatus = "banned"   // Đã bị ban
	StatusInactive PlayerStatus = "inactive" // Không hoạt động
)

// Player là bản ghi người chơi chính. Mỗi cặp (userId, guildId) có một Player duy nhất.
type Player struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"userId"        json:"userId"`
	GuildID      string             `bson:"guildId"       json:"guildId"`
	Username     string             `bson:"username"      json:"username"`    // Tên Discord username
	DisplayName  string             `bson:"displayName"   json:"displayName"` // Tên hiển thị Discord
	DaoName      string             `bson:"daoName"       json:"daoName"`     // Đạo hiệu (người chơi chọn)
	Status       PlayerStatus       `bson:"status"        json:"status"`
	CreatedAt    time.Time          `bson:"createdAt"     json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"     json:"updatedAt"`
	LastActiveAt time.Time          `bson:"lastActiveAt"  json:"lastActiveAt"`
}

// IsActive trả về true nếu tài khoản không bị cấm hoặc vô hiệu hóa.
func (p *Player) IsActive() bool {
	return p.Status == StatusActive
}
