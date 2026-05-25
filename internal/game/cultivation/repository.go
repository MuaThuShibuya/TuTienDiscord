// File: internal/game/cultivation/repository.go
// Phiên bản: v0.1.1
// Mục đích: Interface Repository cho CultivationProfile.
//           MongoDB implementation nằm trong mongo_repository.go.
// Ghi chú: Service chỉ phụ thuộc vào interface — có thể mock trong test.

package cultivation

import "context"

// Repository định nghĩa thao tác lưu trữ cho CultivationProfile.
type Repository interface {
	// FindByUserID tìm hồ sơ tu luyện theo (userId, guildId).
	FindByUserID(ctx context.Context, userID, guildID string) (*CultivationProfile, error)

	// Upsert tạo mới hoặc cập nhật hồ sơ tu luyện.
	Upsert(ctx context.Context, profile *CultivationProfile) error
	// TODO v0.2: thêm AddExp, SetRealm, SetMindState với atomic update
}
