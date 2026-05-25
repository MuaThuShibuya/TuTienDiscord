// File: internal/game/economy/model.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa struct Wallet — ví tiền trong game.
// Bảo mật: Chỉ thay đổi currency qua service dùng atomic DB operation.
//          Tuyệt đối không nhận số lượng tiền trực tiếp từ Discord user input.
// Ghi chú: Không có tiền thật. Gacha chỉ dùng vé/linh ngọc trong game.

package economy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Wallet chứa toàn bộ tài nguyên tiền tệ trong game của một người chơi.
type Wallet struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"userId"        json:"userId"`
	GuildID      string             `bson:"guildId"       json:"guildId"`
	SpiritStones int64              `bson:"spiritStones"  json:"spiritStones"` // Linh Thạch — tiền phổ thông
	SpiritJades  int64              `bson:"spiritJades"   json:"spiritJades"`  // Linh Ngọc — tiền cao cấp trong game
	FateTickets  int                `bson:"fateTickets"   json:"fateTickets"`  // Vé Cơ Duyên — dùng để gacha
	CreatedAt    time.Time          `bson:"createdAt"     json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"     json:"updatedAt"`
}
