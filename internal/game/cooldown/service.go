// File: internal/game/cooldown/service.go
// Version: v0.1
// Purpose: Business logic for cooldown — check if active, set, and clear cooldowns.
// Security: Action names are always server-side constants. Never pass raw user input as action.
// Notes: IsOnCooldown returns false if DB lookup fails, to avoid blocking players due to DB errors.

package cooldown

import (
	"context"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// Service defines cooldown business operations.
type Service interface {
	// IsOnCooldown returns (true, remaining) if the action is on cooldown.
	IsOnCooldown(ctx context.Context, userID, guildID string, action Action) (bool, time.Duration)
	// SetCooldown starts a cooldown for the given duration.
	SetCooldown(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error
	// ClearCooldown removes a cooldown before it expires (e.g., admin reset).
	ClearCooldown(ctx context.Context, userID, guildID string, action Action) error
}

type service struct {
	repo Repository
	log  *zap.Logger
}

// NewService creates a new cooldown service.
func NewService(repo Repository) Service {
	return &service{repo: repo, log: logger.L().Named("cooldown.service")}
}

func (s *service) IsOnCooldown(ctx context.Context, userID, guildID string, action Action) (bool, time.Duration) {
	cd, err := s.repo.Get(ctx, userID, guildID, action)
	if err != nil {
		if !apperrors.IsNotFound(err) {
			s.log.Warn("IsOnCooldown DB error (allowing action)",
				zap.String("userId", userID),
				zap.String("action", string(action)),
				zap.Error(err),
			)
		}
		return false, 0
	}
	remaining := cd.RemainingDuration()
	if remaining <= 0 {
		return false, 0
	}
	return true, remaining
}

func (s *service) SetCooldown(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error {
	if err := s.repo.Set(ctx, userID, guildID, action, duration); err != nil {
		s.log.Error("SetCooldown failed",
			zap.String("userId", userID),
			zap.String("action", string(action)),
			zap.Duration("duration", duration),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (s *service) ClearCooldown(ctx context.Context, userID, guildID string, action Action) error {
	return s.repo.Delete(ctx, userID, guildID, action)
}
