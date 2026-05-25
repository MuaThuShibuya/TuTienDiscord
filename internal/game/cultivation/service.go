// File: internal/game/cultivation/service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic cho hệ thống tu luyện — lấy hoặc khởi tạo hồ sơ.
// Ghi chú: Logic tĩnh tu, đột phá, chọn đạo lộ sẽ thêm vào v0.2.
//          Service không gọi Discord API, không build embed.

package cultivation

import (
	"context"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Service định nghĩa các nghiệp vụ tu luyện.
type Service interface {
	// GetOrCreate lấy hồ sơ tu luyện hoặc tạo mới với giá trị mặc định.
	GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error)

	// GetProfile lấy hồ sơ tu luyện, trả về ErrNotFound nếu chưa có.
	GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error)
	// TODO v0.2: Cultivate, Breakthrough, ChoosePath, ConsumePill
}

type cultivationService struct {
	repo Repository
	log  *zap.Logger
}

// NewService tạo cultivation service.
func NewService(repo Repository) Service {
	return &cultivationService{repo: repo, log: logger.L().Named("cultivation.service")}
}

func (s *cultivationService) GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return profile, nil
	}
	if !apperrors.IsNotFound(err) {
		s.log.Error("GetOrCreate: lỗi DB", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	// Tạo hồ sơ mới với giá trị khởi đầu
	newProfile := NewCultivationProfile(userID, guildID)
	if err := s.repo.Upsert(ctx, newProfile); err != nil {
		s.log.Error("GetOrCreate: không tạo được hồ sơ tu luyện",
			zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	s.log.Info("Hồ sơ tu luyện mới được tạo",
		zap.String("userId", userID), zap.String("guildId", guildID))
	return newProfile, nil
}

func (s *cultivationService) GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}
