// File: internal/discord/menu/session_service.go
// Chức năng: Business logic quản lý vòng đời phiên menu.
// Bảo mật: ValidateOwner PHẢI được gọi trước mọi thao tác button/select menu.

package menu

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

// SessionService quản lý vòng đời phiên menu.
type SessionService interface {
	OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error)
	ValidateOwner(ctx context.Context, sessionID, userID string) (*Session, error)
	NavigateTo(ctx context.Context, sessionID string, page Page, category string) error
	SetMessageID(ctx context.Context, sessionID, messageID string) error
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error
	Close(ctx context.Context, sessionID string) error
}

type sessionSvc struct {
	repo SessionRepository
	log  *zap.Logger
}

// NewSessionService tạo session service với repository đã cho.
func NewSessionService(repo SessionRepository) SessionService {
	return &sessionSvc{repo: repo, log: logger.L().Named("menu.session")}
}

func (s *sessionSvc) OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error) {
	if err := s.repo.DeleteExpiredByUser(ctx, userID, guildID); err != nil {
		s.log.Warn("OpenMenu: không dọn được phiên cũ",
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

	s.log.Debug("Phiên menu mới được tạo",
		zap.String("userId", userID),
		zap.String("sessionId", session.SessionID),
	)
	return session, nil
}

func (s *sessionSvc) ValidateOwner(ctx context.Context, sessionID, userID string) (*Session, error) {
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

func (s *sessionSvc) NavigateTo(ctx context.Context, sessionID string, page Page, category string) error {
	return s.repo.UpdatePage(ctx, sessionID, page, category)
}

func (s *sessionSvc) SetMessageID(ctx context.Context, sessionID, messageID string) error {
	return s.repo.UpdateMessageID(ctx, sessionID, messageID)
}

func (s *sessionSvc) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	return s.repo.Refresh(ctx, sessionID, ttl)
}

func (s *sessionSvc) Close(ctx context.Context, sessionID string) error {
	return s.repo.Delete(ctx, sessionID)
}
