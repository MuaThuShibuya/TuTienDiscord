// File: internal/game/combat/repository.go
package combat

import (
	"context"
	"time"
)

type Repository interface {
	CreateSession(ctx context.Context, session *CombatSession) error
	GetSession(ctx context.Context, sessionID string) (*CombatSession, error)
	GetActiveSessionByUser(ctx context.Context, userID string) (*CombatSession, error)
	UpdateSession(ctx context.Context, session *CombatSession) error
	MarkSessionState(ctx context.Context, sessionID string, state SessionState) error
	TryStartRewardClaim(ctx context.Context, sessionID string, claimID string, now time.Time) (*CombatSession, error)
	CompleteRewardClaim(ctx context.Context, sessionID string, claimID string, details []ClaimedReward, now time.Time) error
	FailRewardClaim(ctx context.Context, sessionID string, claimID string, reason string, now time.Time) error
}
