// File: internal/game/economy/service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic cho ví tiền — lấy hoặc tạo ví, kiếm/tiêu tài nguyên.
// Bảo mật: Validate amount > 0 trước mọi thao tác. Không chấp nhận số âm từ caller.
// Ghi chú: Không gọi Discord API. Không build embed.

package economy

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Service định nghĩa nghiệp vụ quản lý ví tiền.
type Service interface {
	// GetOrCreate lấy ví hoặc tạo mới với tài nguyên khởi đầu.
	GetOrCreate(ctx context.Context, userID, guildID string) (*Wallet, error)

	// GetWallet lấy ví, trả về ErrNotFound nếu chưa có.
	GetWallet(ctx context.Context, userID, guildID string) (*Wallet, error)

	// EarnSpiritStones cộng linh thạch (amount phải > 0).
	EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)

	// SpendSpiritStones trừ linh thạch. Trả về ErrInsufficientFunds nếu không đủ.
	SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)

	// EarnSpiritJades cộng linh ngọc.
	EarnSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)

	// SpendSpiritJades trừ linh ngọc.
	SpendSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)

	// EarnFateTickets cộng vé cơ duyên.
	EarnFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error)

	// SpendFateTickets trừ vé cơ duyên.
	SpendFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error)
}

type economyService struct {
	repo Repository
	log  *zap.Logger
}

// NewService tạo economy service.
func NewService(repo Repository) Service {
	return &economyService{repo: repo, log: logger.L().Named("economy.service")}
}

func (s *economyService) GetOrCreate(ctx context.Context, userID, guildID string) (*Wallet, error) {
	wallet, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return wallet, nil
	}
	if !apperrors.IsNotFound(err) {
		return nil, err
	}
	newWallet := NewWallet(userID, guildID)
	if err := s.repo.Upsert(ctx, newWallet); err != nil {
		s.log.Error("GetOrCreate wallet thất bại", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}
	return newWallet, nil
}

func (s *economyService) GetWallet(ctx context.Context, userID, guildID string) (*Wallet, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}

func (s *economyService) EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số linh thạch phải lớn hơn 0.", nil)
	}
	w, err := s.repo.AdjustSpiritStones(ctx, userID, guildID, amount)
	if err != nil {
		return nil, fmt.Errorf("EarnSpiritStones: %w", err)
	}
	s.log.Debug("Cộng linh thạch", zap.String("userId", userID), zap.Int64("amount", amount), zap.String("lý do", reason))
	return w, nil
}

func (s *economyService) SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số linh thạch phải lớn hơn 0.", nil)
	}
	w, err := s.repo.AdjustSpiritStones(ctx, userID, guildID, -amount)
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return nil, apperrors.New("INSUFFICIENT_SPIRIT_STONES", "Đạo hữu không đủ linh thạch.", err)
		}
		return nil, fmt.Errorf("SpendSpiritStones: %w", err)
	}
	s.log.Debug("Trừ linh thạch", zap.String("userId", userID), zap.Int64("amount", amount), zap.String("lý do", reason))
	return w, nil
}

func (s *economyService) EarnSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số linh ngọc phải lớn hơn 0.", nil)
	}
	return s.repo.AdjustSpiritJades(ctx, userID, guildID, amount)
}

func (s *economyService) SpendSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số linh ngọc phải lớn hơn 0.", nil)
	}
	w, err := s.repo.AdjustSpiritJades(ctx, userID, guildID, -amount)
	if err != nil && apperrors.IsInsufficientFunds(err) {
		return nil, apperrors.New("INSUFFICIENT_SPIRIT_JADES", "Đạo hữu không đủ linh ngọc.", err)
	}
	return w, err
}

func (s *economyService) EarnFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số vé phải lớn hơn 0.", nil)
	}
	return s.repo.AdjustFateTickets(ctx, userID, guildID, amount)
}

func (s *economyService) SpendFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số vé phải lớn hơn 0.", nil)
	}
	w, err := s.repo.AdjustFateTickets(ctx, userID, guildID, -amount)
	if err != nil && apperrors.IsInsufficientFunds(err) {
		return nil, apperrors.New("INSUFFICIENT_FATE_TICKETS", "Đạo hữu không đủ vé cơ duyên.", err)
	}
	return w, err
}
