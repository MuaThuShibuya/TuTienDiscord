// File: internal/discord/menu/session_model.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa struct Session — phiên menu của người chơi.
// Bảo mật: SessionID được tạo bằng crypto/rand để không thể đoán được.
//           Mọi interaction phải xác thực OwnedBy() trước khi xử lý.
// Ghi chú: MongoDB TTL index trên expiresAt tự động xóa session hết hạn.

package menu

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session một phiên menu đang hoạt động của người chơi.
type Session struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"   json:"id"`
	SessionID       string             `bson:"sessionId"       json:"sessionId"` // ID ngẫu nhiên mạnh — nhúng vào custom_id
	UserID          string             `bson:"userId"          json:"userId"`
	GuildID         string             `bson:"guildId"         json:"guildId"`
	ChannelID       string             `bson:"channelId"       json:"channelId"`
	MessageID       string             `bson:"messageId"       json:"messageId"` // Discord message ID để edit lại
	CurrentPage     Page               `bson:"currentPage"     json:"currentPage"`
	CurrentCategory string             `bson:"currentCategory" json:"currentCategory"` // Sub-navigation trong trang
	ExpiresAt       time.Time          `bson:"expiresAt"       json:"expiresAt"`
	CreatedAt       time.Time          `bson:"createdAt"       json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"       json:"updatedAt"`
}

// IsExpired trả về true nếu phiên đã hết hạn.
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// OwnedBy trả về true nếu phiên thuộc về userID đã cho.
func (s *Session) OwnedBy(userID string) bool {
	return s.UserID == userID
}
