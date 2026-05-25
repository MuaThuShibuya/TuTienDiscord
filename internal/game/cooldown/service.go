// File: internal/game/cooldown/service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic cooldown — kiểm tra, thiết lập, xóa cooldown.
// Bảo mật: Action luôn là hằng server-side — không nhận tên action từ user input.
// Ghi chú: IsOnCooldown trả về false nếu lỗi DB để không chặn người chơi oan.

package cooldown

import (
	"context"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Service định nghĩa nghiệp vụ quản lý cooldown.
type Service interface {
	// IsOnCooldown trả về (true, thời_gian_còn_lại) nếu đang trong cooldown.
	IsOnCooldown(ctx context.Context, userID, guildID string, action Action) (bool, time.Duration)

	// SetCooldown bắt đầu cooldown trong khoảng thời gian duration.
	SetCooldown(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error

	// ClearCooldown xóa cooldown trước khi hết hạn (dùng cho admin reset).
	ClearCooldown(ctx context.Context, userID, guildID string, action Action) error
}

type cooldownService struct {
	repo Repository
	log  *zap.Logger
}

// NewService tạo cooldown service.
func NewService(repo Repository) Service {
	return &cooldownService{repo: repo, log: logger.L().Named("cooldown.service")}
}

func (s *cooldownService) IsOnCooldown(ctx context.Context, userID, guildID string, action Action) (bool, time.Duration) {
	cd, err := s.repo.Get(ctx, userID, guildID, action)
	if err != nil {
		if !apperrors.IsNotFound(err) {
			// Lỗi DB không mong đợi — log nhưng không chặn người chơi
			s.log.Warn("IsOnCooldown: lỗi DB (cho phép hành động)",
				zap.String("userId", userID), zap.String("action", string(action)), zap.Error(err))
		}
		return false, 0
	}
	remaining := cd.Remaining()
	if remaining <= 0 {
		return false, 0
	}
	return true, remaining
}

func (s *cooldownService) SetCooldown(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error {
	if err := s.repo.Set(ctx, userID, guildID, action, duration); err != nil {
		s.log.Error("SetCooldown thất bại",
			zap.String("userId", userID), zap.String("action", string(action)), zap.Error(err))
		return err
	}
	return nil
}

func (s *cooldownService) ClearCooldown(ctx context.Context, userID, guildID string, action Action) error {
	return s.repo.Delete(ctx, userID, guildID, action)
}
