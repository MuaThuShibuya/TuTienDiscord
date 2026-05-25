// File: internal/discord/menu/session_service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic quản lý vòng đời phiên menu — mở, xác thực chủ sở hữu,
//           điều hướng trang, gia hạn TTL, và đóng phiên.
// Bảo mật: ValidateOwner PHẢI được gọi trước mọi thao tác button/select menu.
//           Người khác bấm menu sẽ bị từ chối với thông báo ephemeral.
// Ghi chú: SessionService không chứa logic game — chỉ quản lý metadata phiên.

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
	// OpenMenu tạo phiên menu mới cho người dùng (xóa phiên cũ đã hết hạn trước).
	OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error)

	// ValidateOwner kiểm tra phiên tồn tại, chưa hết hạn, và thuộc về userID.
	// Trả về ErrSessionExpired hoặc ErrSessionNotOwner nếu không hợp lệ.
	ValidateOwner(ctx context.Context, sessionID, userID string) (*Session, error)

	// NavigateTo cập nhật trang hiện tại của phiên.
	NavigateTo(ctx context.Context, sessionID string, page Page) error

	// SetMessageID lưu Discord message ID sau khi gửi menu lần đầu.
	SetMessageID(ctx context.Context, sessionID, messageID string) error

	// Refresh gia hạn TTL của phiên sau mỗi lần tương tác.
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error

	// Close xóa phiên ngay lập tức (người dùng nhấn "Đóng").
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

// OpenMenu dọn phiên hết hạn cũ rồi tạo phiên mới.
func (s *sessionSvc) OpenMenu(ctx context.Context, userID, guildID, channelID string, ttl time.Duration) (*Session, error) {
	// Dọn phiên cũ đã hết hạn của người dùng này
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

// ValidateOwner xác thực phiên tồn tại, chưa hết hạn, và thuộc về đúng người.
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

// NavigateTo cập nhật trang hiện tại và xóa sub-category.
func (s *sessionSvc) NavigateTo(ctx context.Context, sessionID string, page Page) error {
	return s.repo.UpdatePage(ctx, sessionID, page, "")
}

// SetMessageID lưu message ID để edit lại về sau.
func (s *sessionSvc) SetMessageID(ctx context.Context, sessionID, messageID string) error {
	return s.repo.UpdateMessageID(ctx, sessionID, messageID)
}

// Refresh gia hạn TTL phiên sau mỗi tương tác.
func (s *sessionSvc) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	return s.repo.Refresh(ctx, sessionID, ttl)
}

// Close xóa phiên ngay lập tức.
func (s *sessionSvc) Close(ctx context.Context, sessionID string) error {
	return s.repo.Delete(ctx, sessionID)
}
