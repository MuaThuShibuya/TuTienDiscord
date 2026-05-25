// File: internal/game/cultivation/service.go
// Version: v0.1
// Purpose: Business logic for cultivation profiles — get or initialize.
// Security: No direct user input touches DB without validation.
// Notes: Cultivation logic (tĩnh tu, đột phá, etc.) will be added in v0.2.

package cultivation

import (
	"context"

	"go.uber.org/zap"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// Service defines cultivation business operations.
type Service interface {
	GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error)
	GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error)
	// TODO v0.2: Cultivate, Breakthrough, SetPath, ConsumePill
}

type service struct {
	repo Repository
	log  *zap.Logger
}

// NewService creates a new cultivation service.
func NewService(repo Repository) Service {
	return &service{repo: repo, log: logger.L().Named("cultivation.service")}
}

// GetOrCreate retrieves the cultivation profile or seeds a new default one.
func (s *service) GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return profile, nil
	}

	if !apperrors.IsNotFound(err) {
		s.log.Error("GetOrCreate: DB error",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
		return nil, err
	}

	newProfile := DefaultCultivationProfile(userID, guildID)
	if err := s.repo.Upsert(ctx, newProfile); err != nil {
		s.log.Error("GetOrCreate: failed to upsert cultivation profile",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
		return nil, err
	}

	s.log.Info("New cultivation profile created",
		zap.String("userId", userID),
		zap.String("guildId", guildID),
	)

	return newProfile, nil
}

// GetProfile retrieves an existing cultivation profile or returns ErrNotFound.
func (s *service) GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}
