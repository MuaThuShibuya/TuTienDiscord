// File: internal/game/economy/service.go
// Version: v0.1
// Purpose: Business logic for player wallet — get or create, check balance, spend/earn currency.
// Security: Amount validation happens here before any DB call. No negative input allowed from outside.
// Notes: All currency operations log the action for future audit trail (v0.8+).

package economy

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// Service defines wallet business operations.
type Service interface {
	GetOrCreate(ctx context.Context, userID, guildID string) (*Wallet, error)
	GetWallet(ctx context.Context, userID, guildID string) (*Wallet, error)
	EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)
	SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)
	EarnSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)
	SpendSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error)
	EarnFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error)
	SpendFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error)
}

type service struct {
	repo Repository
	log  *zap.Logger
}

// NewService creates a new economy service.
func NewService(repo Repository) Service {
	return &service{repo: repo, log: logger.L().Named("economy.service")}
}

func (s *service) GetOrCreate(ctx context.Context, userID, guildID string) (*Wallet, error) {
	wallet, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return wallet, nil
	}
	if !apperrors.IsNotFound(err) {
		return nil, err
	}

	newWallet := DefaultWallet(userID, guildID)
	if err := s.repo.Upsert(ctx, newWallet); err != nil {
		s.log.Error("GetOrCreate wallet failed",
			zap.String("userId", userID),
			zap.String("guildId", guildID),
			zap.Error(err),
		)
		return nil, err
	}
	return newWallet, nil
}

func (s *service) GetWallet(ctx context.Context, userID, guildID string) (*Wallet, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}

func (s *service) EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng linh thạch phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustSpiritStones(ctx, userID, guildID, amount)
	if err != nil {
		return nil, fmt.Errorf("EarnSpiritStones: %w", err)
	}
	s.log.Debug("SpiritStones earned",
		zap.String("userId", userID), zap.Int64("amount", amount), zap.String("reason", reason))
	return wallet, nil
}

func (s *service) SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng linh thạch phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustSpiritStones(ctx, userID, guildID, -amount)
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return nil, apperrors.New("INSUFFICIENT_SPIRIT_STONES",
				"Đạo hữu không đủ linh thạch.", err)
		}
		return nil, fmt.Errorf("SpendSpiritStones: %w", err)
	}
	s.log.Debug("SpiritStones spent",
		zap.String("userId", userID), zap.Int64("amount", amount), zap.String("reason", reason))
	return wallet, nil
}

func (s *service) EarnSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng linh ngọc phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustSpiritJades(ctx, userID, guildID, amount)
	if err != nil {
		return nil, fmt.Errorf("EarnSpiritJades: %w", err)
	}
	s.log.Debug("SpiritJades earned",
		zap.String("userId", userID), zap.Int64("amount", amount), zap.String("reason", reason))
	return wallet, nil
}

func (s *service) SpendSpiritJades(ctx context.Context, userID, guildID string, amount int64, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng linh ngọc phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustSpiritJades(ctx, userID, guildID, -amount)
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return nil, apperrors.New("INSUFFICIENT_SPIRIT_JADES",
				"Đạo hữu không đủ linh ngọc.", err)
		}
		return nil, fmt.Errorf("SpendSpiritJades: %w", err)
	}
	s.log.Debug("SpiritJades spent",
		zap.String("userId", userID), zap.Int64("amount", amount), zap.String("reason", reason))
	return wallet, nil
}

func (s *service) EarnFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng vé cơ duyên phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustFateTickets(ctx, userID, guildID, amount)
	if err != nil {
		return nil, fmt.Errorf("EarnFateTickets: %w", err)
	}
	s.log.Debug("FateTickets earned",
		zap.String("userId", userID), zap.Int("amount", amount), zap.String("reason", reason))
	return wallet, nil
}

func (s *service) SpendFateTickets(ctx context.Context, userID, guildID string, amount int, reason string) (*Wallet, error) {
	if amount <= 0 {
		return nil, apperrors.New("INVALID_AMOUNT", "Số lượng vé cơ duyên phải lớn hơn 0.", nil)
	}
	wallet, err := s.repo.AdjustFateTickets(ctx, userID, guildID, -amount)
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return nil, apperrors.New("INSUFFICIENT_FATE_TICKETS",
				"Đạo hữu không đủ vé cơ duyên.", err)
		}
		return nil, fmt.Errorf("SpendFateTickets: %w", err)
	}
	s.log.Debug("FateTickets spent",
		zap.String("userId", userID), zap.Int("amount", amount), zap.String("reason", reason))
	return wallet, nil
}
