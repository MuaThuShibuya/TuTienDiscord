// File: internal/game/pve/service.go

package pve

import (
	"context"
	"errors"
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
)

type ProgressService interface {
	CanEnterArea(ctx context.Context, userID string, area PvEAreaDefinition, combatPower int64, realm string) error
	MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error
	GetNextStage(ctx context.Context, userID, areaID string) (int, error)
}

type progressService struct {
	repo ProgressRepository
}

func NewProgressService(repo ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) CanEnterArea(ctx context.Context, userID string, area PvEAreaDefinition, combatPower int64, realm string) error {
	if area.RequiredCombatPower > 0 && combatPower < area.RequiredCombatPower {
		return fmt.Errorf("cần tối thiểu %d lực chiến để vào khu vực này", area.RequiredCombatPower)
	}

	if area.RequiredRealm != "" && realm != area.RequiredRealm {
		return fmt.Errorf("cảnh giới không phù hợp, yêu cầu: %s", area.RequiredRealm)
	}

	_, err := s.repo.GetAreaProgress(ctx, userID, area.ID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return err
	}
	// TODO: Kiểm tra daily limit attempts ở đây nếu thiết kế yêu cầu giới hạn lượt đánh.
	return nil
}

func (s *progressService) MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error {
	prog, err := s.repo.GetAreaProgress(ctx, userID, areaID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			def, ok := AreaRegistry[areaID]
			if !ok {
				return fmt.Errorf("khu vực không tồn tại")
			}
			newProg := AreaProgress{AreaID: areaID, ActivityType: def.ActivityType, HighestStageCleared: stage}
			return s.repo.UpsertAreaProgress(ctx, userID, newProg)
		}
		return err
	}
	if stage > prog.HighestStageCleared {
		return s.repo.MarkStageCleared(ctx, userID, areaID, stage)
	}
	return nil
}

func (s *progressService) GetNextStage(ctx context.Context, userID, areaID string) (int, error) {
	prog, err := s.repo.GetAreaProgress(ctx, userID, areaID)
	def, ok := AreaRegistry[areaID]
	if !ok {
		return 1, fmt.Errorf("khu vực không tồn tại")
	}
	if err != nil {
		return def.MinStage, nil
	}
	next := prog.HighestStageCleared + 1
	if next > def.MaxStage {
		return def.MaxStage, nil
	}
	return next, nil
}
