// File: internal/game/combat/repository.go
package combat

import "context"

type Repository interface {
	CreateSession(ctx context.Context, session *CombatSession) error
	GetSession(ctx context.Context, sessionID string) (*CombatSession, error)
	GetActiveSessionByUser(ctx context.Context, userID string) (*CombatSession, error)
	UpdateSession(ctx context.Context, session *CombatSession) error
	MarkSessionState(ctx context.Context, sessionID string, state SessionState) error
}
