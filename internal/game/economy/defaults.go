// File: internal/game/economy/defaults.go
// Phiên bản: v0.1.1
// Mục đích: Giá trị tài nguyên khởi đầu cho người chơi mới.
// Ghi chú: Chỉnh số liệu ở đây để điều chỉnh balance game.

package economy

import "go.mongodb.org/mongo-driver/bson/primitive"

// Tài nguyên khởi đầu khi người chơi đăng ký lần đầu.
const (
	DefaultSpiritStones = int64(500) // 500 Linh Thạch để bắt đầu
	DefaultSpiritJades  = int64(0)
	DefaultFateTickets  = 3 // 3 vé gacha miễn phí lúc đầu
)

// NewWallet tạo ví mới với tài nguyên khởi đầu.
func NewWallet(userID, guildID string) *Wallet {
	return &Wallet{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		GuildID:      guildID,
		SpiritStones: DefaultSpiritStones,
		SpiritJades:  DefaultSpiritJades,
		FateTickets:  DefaultFateTickets,
	}
}
