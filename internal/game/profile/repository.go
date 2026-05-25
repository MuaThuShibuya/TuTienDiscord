// File: internal/game/profile/repository.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa interface Repository cho Player.
//           Implementation MongoDB nằm trong mongo_repository.go.
//           Implementation giả (cho test) nằm trong service_test.go.
// Ghi chú: Service chỉ phụ thuộc vào interface này — không biết MongoDB là gì.

package profile

import "context"

// Repository định nghĩa các thao tác lưu trữ cho Player.
// Mọi implementation phải tuân theo interface này.
type Repository interface {
	// FindByUserID tìm player theo cặp (userId, guildId).
	// Trả về apperrors.ErrNotFound nếu không tồn tại.
	FindByUserID(ctx context.Context, userID, guildID string) (*Player, error)

	// Create lưu player mới vào database.
	// Trả về apperrors.ErrAlreadyExists nếu đã tồn tại.
	Create(ctx context.Context, player *Player) error

	// UpdateLastActive cập nhật thời điểm hoạt động gần nhất.
	UpdateLastActive(ctx context.Context, userID, guildID string) error

	// UpdateDaoName thay đổi đạo hiệu của player.
	UpdateDaoName(ctx context.Context, userID, guildID, daoName string) error
}
