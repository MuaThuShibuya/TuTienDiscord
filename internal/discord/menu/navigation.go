// File: internal/discord/menu/navigation.go
// Version: v0.1
// Purpose: Session service — business logic for creating, validating, and navigating menu sessions.
// Security: ValidateOwner must be called before every button/select interaction that modifies state.
// Notes: NavigateTo updates the session's current page and triggers a UI re-render via the page router.

package menu

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
	"github.com/yourname/tu-tien-bot/pkg/utils"
)

// Service manages menu session lifecycle and navigation.
type Service interface {
	// OpenMenu creates (or replaces) a session for the user and returns it.
	OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error)
	// ValidateOwner checks that the given sessionID exists, is not expired, and belongs to userID.
	ValidateOwner(ctx context.Context, sessionID, userID string) (*Session, error)
	// NavigateTo updates the session's current page.
	NavigateTo(ctx context.Context, sessionID string, page Page) error
	// Close deletes a session immediately (user pressed "Đóng").
	Close(ctx context.Context, sessionID string) error
	// SetMessageID persists the Discord message ID after the menu is first sent.
	SetMessageID(ctx context.Context, sessionID, messageID string) error
	// Refresh extends the session TTL (called on every user interaction).
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error
}

type sessionService struct {
	repo Repository
	log  *zap.Logger
}

// NewService creates a new menu session service.
func NewService(repo Repository) Service {
	return &sessionService{repo: repo, log: logger.L().Named("menu.service")}
}

// OpenMenu closes any old session for the user and opens a fresh one.
func (s *sessionService) OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error) {
	// Clean up expired sessions for this user first.
	if err := s.repo.DeleteExpiredByUser(ctx, userID, guildID); err != nil {
		s.log.Warn("OpenMenu: failed to clean expired sessions",
			zap.String("userId", userID), zap.Error(err))
	}

	session := &Session{
		SessionID:   utils.NewSessionID(),
		UserID:      userID,
		GuildID:     guildID,
		ChannelID:   channelID,
		CurrentPage: PageMain,
		ExpiresAt:   time.Now().UTC().Add(ttl),
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("OpenMenu: %w", err)
	}

	s.log.Debug("Menu session opened",
		zap.String("userId", userID),
		zap.String("sessionId", session.SessionID),
	)

	return session, nil
}

// ValidateOwner verifies the session exists, is not expired, and belongs to the given user.
func (s *sessionService) ValidateOwner(ctx context.Context, sessionID, userID string) (*Session, error) {
	session, err := s.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			return nil, apperrors.ErrSessionExpired
		}
		return nil, fmt.Errorf("ValidateOwner: %w", err)
	}

	if session.IsExpired() {
		_ = s.repo.Delete(ctx, sessionID)
		return nil, apperrors.ErrSessionExpired
	}

	if !session.OwnedBy(userID) {
		return nil, apperrors.ErrSessionNotOwner
	}

	return session, nil
}

// NavigateTo sets the current page and clears sub-category navigation.
func (s *sessionService) NavigateTo(ctx context.Context, sessionID string, page Page) error {
	return s.repo.UpdatePage(ctx, sessionID, page, "")
}

// Close deletes the session, disabling all its buttons.
func (s *sessionService) Close(ctx context.Context, sessionID string) error {
	return s.repo.Delete(ctx, sessionID)
}

// SetMessageID stores the Discord message ID after first render.
func (s *sessionService) SetMessageID(ctx context.Context, sessionID, messageID string) error {
	return s.repo.UpdateMessageID(ctx, sessionID, messageID)
}

// Refresh extends the session TTL after any user interaction.
func (s *sessionService) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	return s.repo.Refresh(ctx, sessionID, ttl)
}
