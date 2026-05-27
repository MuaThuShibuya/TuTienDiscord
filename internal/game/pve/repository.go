// File: internal/game/pve/repository.go

package pve

import "context"

type ProgressRepository interface {
	GetProgress(ctx context.Context, userID string) (*UserPvEProgress, error)
	GetAreaProgress(ctx context.Context, userID, areaID string) (*AreaProgress, error)
	UpsertAreaProgress(ctx context.Context, userID string, progress AreaProgress) error
	MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error
	IncrementAttempt(ctx context.Context, userID, areaID string) error
}
