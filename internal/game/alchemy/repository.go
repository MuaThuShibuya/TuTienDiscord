// File: internal/game/alchemy/repository.go
// Chức năng: Interface định nghĩa các thao tác lưu trữ hồ sơ luyện đan.

package alchemy

import "context"

type Repository interface {
	// Get lấy hồ sơ luyện đan của người chơi.
	Get(ctx context.Context, userID, guildID string) (*AlchemyProfile, error)

	// Upsert tạo mới nếu chưa có hoặc cập nhật hồ sơ.
	Upsert(ctx context.Context, profile *AlchemyProfile) error
}
