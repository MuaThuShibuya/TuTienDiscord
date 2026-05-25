// File: internal/game/economy/repository.go
// Phiên bản: v0.1.1
// Mục đích: Interface Repository cho Wallet.
//           MongoDB implementation nằm trong mongo_repository.go.
// Ghi chú: AdjustCurrency dùng atomic $inc — tránh race condition khi nhiều người giao dịch.

package economy

import "context"

// Repository định nghĩa thao tác lưu trữ cho Wallet.
type Repository interface {
	// FindByUserID tìm ví theo (userId, guildId).
	FindByUserID(ctx context.Context, userID, guildID string) (*Wallet, error)

	// Upsert tạo mới hoặc cập nhật ví.
	Upsert(ctx context.Context, wallet *Wallet) error

	// AdjustSpiritStones thay đổi linh thạch một cách atomic.
	// amount âm = trừ tiền. Trả về ErrInsufficientFunds nếu số dư < |amount|.
	AdjustSpiritStones(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error)

	// AdjustSpiritJades thay đổi linh ngọc một cách atomic.
	AdjustSpiritJades(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error)

	// AdjustFateTickets thay đổi vé cơ duyên một cách atomic.
	AdjustFateTickets(ctx context.Context, userID, guildID string, amount int) (*Wallet, error)
}
