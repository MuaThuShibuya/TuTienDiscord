// File: internal/game/cooldown/repository.go
// Phiên bản: v0.1.1
// Mục đích: Interface Repository cho Cooldown.
//           MongoDB implementation nằm trong mongo_repository.go.

package cooldown

import (
	"context"
	"time"
)

// Repository định nghĩa thao tác lưu trữ cho Cooldown.
type Repository interface {
	// Get lấy cooldown đang hoạt động (chưa hết hạn).
	// Trả về ErrNotFound nếu không có hoặc đã hết hạn.
	Get(ctx context.Context, userID, guildID string, action Action) (*Cooldown, error)

	// Set thiết lập hoặc cập nhật cooldown cho action với thời gian đã cho.
	Set(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error

	// Delete xóa cooldown trước khi hết hạn (ví dụ: admin reset).
	Delete(ctx context.Context, userID, guildID string, action Action) error
}
