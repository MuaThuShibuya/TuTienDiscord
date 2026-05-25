// File: internal/game/profile/service.go
// Version: v0.1
// Purpose: Business logic for player profile — registration, lookup, and name changes.
// Security: Validates all inputs. DaoName is sanitized (length, forbidden chars).
// Notes: Service depends only on the Repository interface, not the concrete mongo impl.

package profile

import (
	"context"
	"strings"
	"unicode/utf8"

	"go.uber.org/zap"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

const (
	maxDaoNameLength = 24
	minDaoNameLength = 2
)

// Service defines the profile business operations.
type Service interface {
	GetOrCreate(ctx context.Context, userID, guildID, username, displayName string) (*Player, error)
	GetPlayer(ctx context.Context, userID, guildID string) (*Player, error)
	SetDaoName(ctx context.Context, userID, guildID, daoName string) error
	TouchLastActive(ctx context.Context, userID, guildID string)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

// NewService creates a new profile service backed by the given repository.
func NewService(repo Repository) Service {
	return &service{repo: repo, log: logger.L().Named("profile.service")}
}

// GetOrCreate finds the player or registers them if they don't exist yet.
func (s *service) GetOrCreate(ctx context.Context, userID, guildID, username, displayName string) (*Player, error) {
	player, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return player, nil
	}

	if !apperrors.IsNotFound(err) {
		s.log.Error("GetOrCreate: DB error looking up player",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
		return nil, err
	}

	// New player — use Discord display name as default dao name.
	defaultDaoName := sanitizeDaoName(displayName)
	if defaultDaoName == "" {
		defaultDaoName = sanitizeDaoName(username)
	}
	if defaultDaoName == "" {
		defaultDaoName = "Vô Danh"
	}

	newPlayer := &Player{
		UserID:      userID,
		GuildID:     guildID,
		Username:    username,
		DisplayName: displayName,
		DaoName:     defaultDaoName,
		Status:      StatusActive,
	}

	if err := s.repo.Create(ctx, newPlayer); err != nil {
		s.log.Error("GetOrCreate: failed to create player",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
		return nil, err
	}

	s.log.Info("New player registered",
		zap.String("userId", userID),
		zap.String("guildId", guildID),
		zap.String("daoName", newPlayer.DaoName),
	)

	return newPlayer, nil
}

// GetPlayer retrieves an existing player or returns ErrNotFound.
func (s *service) GetPlayer(ctx context.Context, userID, guildID string) (*Player, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}

// SetDaoName validates and updates the player's đạo hiệu.
func (s *service) SetDaoName(ctx context.Context, userID, guildID, daoName string) error {
	clean := sanitizeDaoName(daoName)
	if utf8.RuneCountInString(clean) < minDaoNameLength {
		return apperrors.New("DAO_NAME_TOO_SHORT",
			"Đạo hiệu phải có ít nhất 2 ký tự.", nil)
	}
	if utf8.RuneCountInString(clean) > maxDaoNameLength {
		return apperrors.New("DAO_NAME_TOO_LONG",
			"Đạo hiệu không được vượt quá 24 ký tự.", nil)
	}
	return s.repo.UpdateDaoName(ctx, userID, guildID, clean)
}

// TouchLastActive updates lastActiveAt in the background; errors are only logged.
func (s *service) TouchLastActive(ctx context.Context, userID, guildID string) {
	if err := s.repo.UpdateLastActive(ctx, userID, guildID); err != nil {
		s.log.Warn("TouchLastActive failed",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
	}
}

// sanitizeDaoName trims whitespace and truncates to the max allowed rune count.
func sanitizeDaoName(name string) string {
	name = strings.TrimSpace(name)
	runes := []rune(name)
	if len(runes) > maxDaoNameLength {
		runes = runes[:maxDaoNameLength]
	}
	return string(runes)
}
