// File: internal/game/alchemy/service.go
// Chức năng: Nghiệp vụ (Business logic) luyện chế đan dược.

package alchemy

import (
	"context"
	"fmt"
	"math/rand"

	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

type CraftResult struct {
	Success bool
	Message string
}

type Service interface {
	GetProfile(ctx context.Context, userID, guildID string) (*AlchemyProfile, error)
	Craft(ctx context.Context, userID, guildID, recipeID string, rnd *rand.Rand) (*CraftResult, error)
}

type alchemyService struct {
	repo         Repository
	inventorySvc inventory.Service
	log          *zap.Logger
}

func NewService(repo Repository, invSvc inventory.Service) Service {
	return &alchemyService{
		repo:         repo,
		inventorySvc: invSvc,
		log:          logger.L().Named("alchemy.service"),
	}
}

func (s *alchemyService) GetProfile(ctx context.Context, userID, guildID string) (*AlchemyProfile, error) {
	return s.repo.Get(ctx, userID, guildID)
}

func (s *alchemyService) Craft(ctx context.Context, userID, guildID, recipeID string, rnd *rand.Rand) (*CraftResult, error) {
	recipe, ok := Recipes[recipeID]
	if !ok {
		return nil, apperrors.ErrInvalidInput
	}

	profile, err := s.repo.Get(ctx, userID, guildID)
	if err != nil {
		return nil, err
	}

	if profile.Level < recipe.LevelRequired {
		return nil, fmt.Errorf("cần cấp luyện đan %d để luyện chế %s", recipe.LevelRequired, recipe.Name)
	}

	// Kiểm tra và trừ nguyên liệu từ inventorySvc
	if err := s.inventorySvc.ConsumeItems(ctx, userID, guildID, recipe.RequiredItems); err != nil {
		return nil, err // Trả về lỗi không đủ nguyên liệu
	}

	// Roll random tính tỉ lệ thành công
	roll := rnd.Float64()
	success := roll <= recipe.SuccessRate

	if !success {
		// Thất bại: Lò nổ, mất nguyên liệu (đã trừ ở trên)
		return &CraftResult{
			Success: false,
			Message: fmt.Sprintf("Luyện chế **%s** thất bại. Linh thảo đã hóa thành tro bụi!", recipe.Name),
		}, nil
	}

	// Thành công: Cộng đan dược vào túi đồ. Nếu túi đầy, gửi trả lại nguyên liệu.
	if err := s.inventorySvc.AddItem(ctx, userID, guildID, recipe.OutputItem, recipe.OutputQuantity); err != nil {
		if apperrors.IsInventoryFull(err) {
			// Hoàn trả nguyên liệu nếu túi đầy
			for defID, qty := range recipe.RequiredItems {
				_ = s.inventorySvc.AddItem(ctx, userID, guildID, defID, qty)
			}
			return nil, fmt.Errorf("túi đồ đã đầy, không thể nhận thêm đan dược")
		}
		return nil, err
	}

	// Cộng kinh nghiệm luyện đan
	profile.Exp += recipe.ExpReward
	// Logic tính toán lên cấp đơn giản (Ví dụ: Cấp * 100 Exp thì lên cấp)
	if profile.Exp >= int64(profile.Level*100) {
		profile.Exp -= int64(profile.Level * 100)
		profile.Level++
	}

	if err := s.repo.Upsert(ctx, profile); err != nil {
		s.log.Error("Không thể cập nhật hồ sơ luyện đan", zap.Error(err))
	}

	return &CraftResult{
		Success: true,
		Message: fmt.Sprintf("Luyện chế thành công! Nhận được **%dx %s**.", recipe.OutputQuantity, recipe.Name),
	}, nil
}
