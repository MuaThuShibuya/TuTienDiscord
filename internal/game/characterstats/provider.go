// File: internal/game/characterstats/provider.go
package characterstats

import (
	"context"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

type Provider interface {
	GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error)
}
