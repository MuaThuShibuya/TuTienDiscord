// File: internal/game/profile/service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic cho profile người chơi — đăng ký, tra cứu, đổi đạo hiệu.
// Bảo mật: Validate input ở đây trước khi gọi repository. DaoName được sanitize.
// Ghi chú: Service không gọi Discord API và không biết embed/component là gì.

package profile

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Service định nghĩa các nghiệp vụ liên quan đến profile người chơi.
type Service interface {
	// GetOrCreate tìm player hoặc đăng ký mới nếu chưa có.
	GetOrCreate(ctx context.Context, userID, guildID, username, displayName string) (*Player, error)

	// GetPlayer lấy thông tin player. Trả về ErrNotFound nếu chưa đăng ký.
	GetPlayer(ctx context.Context, userID, guildID string) (*Player, error)

	// SetDaoName validate và thay đổi đạo hiệu của player.
	SetDaoName(ctx context.Context, userID, guildID, daoName string) error

	// TouchLastActive cập nhật thời điểm online gần nhất (background, lỗi chỉ log).
	TouchLastActive(ctx context.Context, userID, guildID string)
}

type profileService struct {
	repo Repository
	log  *zap.Logger
}

// NewService tạo service profile với repository đã cho.
func NewService(repo Repository) Service {
	return &profileService{repo: repo, log: logger.L().Named("profile.service")}
}

// GetOrCreate tìm player hoặc tạo mới với giá trị khởi đầu.
func (s *profileService) GetOrCreate(ctx context.Context, userID, guildID, username, displayName string) (*Player, error) {
	// Tìm player hiện có
	player, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return player, nil // Đã tồn tại
	}
	if !apperrors.IsNotFound(err) {
		s.log.Error("GetOrCreate: lỗi DB khi tìm player",
			zap.String("userId", userID), zap.String("guildId", guildID), zap.Error(err))
		return nil, err
	}

	// Tạo đạo hiệu mặc định từ tên Discord
	defaultDaoName := sanitizeDaoName(displayName)
	if defaultDaoName == "" {
		defaultDaoName = sanitizeDaoName(username)
	}
	if defaultDaoName == "" {
		defaultDaoName = "Vô Danh"
	}

	newPlayer := NewPlayer(userID, guildID, username, displayName, defaultDaoName)
	if err := s.repo.Create(ctx, newPlayer); err != nil {
		s.log.Error("GetOrCreate: không tạo được player",
			zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	s.log.Info("Người chơi mới đã đăng ký",
		zap.String("userId", userID), zap.String("guildId", guildID), zap.String("daoName", newPlayer.DaoName))
	return newPlayer, nil
}

// GetPlayer lấy thông tin player, trả về lỗi nếu chưa đăng ký.
func (s *profileService) GetPlayer(ctx context.Context, userID, guildID string) (*Player, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}

// SetDaoName validate và cập nhật đạo hiệu người chơi tự chọn.
// Khác với GetOrCreate — ở đây KHÔNG truncate, phải báo lỗi nếu vượt giới hạn.
// Dùng utf8.RuneCountInString thay vì len() vì len() đếm byte, sai với tiếng Việt có dấu.
func (s *profileService) SetDaoName(ctx context.Context, userID, guildID, daoName string) error {
	// Chỉ trim khoảng trắng — không truncate để giữ nguyên ý định của người dùng
	daoName = strings.TrimSpace(daoName)
	runeLen := utf8.RuneCountInString(daoName)

	if runeLen < MinDaoNameLength {
		return fmt.Errorf("%w: đạo hiệu phải có ít nhất %d ký tự", apperrors.ErrInvalidDaoName, MinDaoNameLength)
	}
	if runeLen > MaxDaoNameLength {
		return fmt.Errorf("%w: đạo hiệu không được quá %d ký tự", apperrors.ErrInvalidDaoName, MaxDaoNameLength)
	}
	return s.repo.UpdateDaoName(ctx, userID, guildID, daoName)
}

// TouchLastActive cập nhật lastActiveAt, lỗi chỉ được log không trả về caller.
func (s *profileService) TouchLastActive(ctx context.Context, userID, guildID string) {
	if err := s.repo.UpdateLastActive(ctx, userID, guildID); err != nil {
		s.log.Warn("TouchLastActive thất bại",
			zap.String("userId", userID), zap.Error(err))
	}
}

// sanitizeDaoName dùng khi GetOrCreate tạo đạo hiệu mặc định từ tên Discord.
// Truncate về MaxDaoNameLength để đảm bảo đạo hiệu mặc định luôn hợp lệ.
// KHÔNG dùng hàm này trong SetDaoName — khi user tự chọn phải báo lỗi, không tự cắt.
func sanitizeDaoName(name string) string {
	name = strings.TrimSpace(name)
	runes := []rune(name)
	if len(runes) > MaxDaoNameLength {
		runes = runes[:MaxDaoNameLength]
	}
	return string(runes)
}
