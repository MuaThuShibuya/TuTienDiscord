// File: internal/game/characterstats/provider.go
package characterstats

import (
	"context"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

type StatsBreakdown struct {
	UserID             string
	AptitudeID         string
	AptitudeName       string
	RealmID            string
	RealmName          string
	RealmLevel         int
	BaseFromRealm      combat.CombatStats
	BaseFromAptitude   combat.CombatStats
	GrowthFromAptitude float64
	EquipmentStats     combat.CombatStats
	SkillPassiveStats  combat.CombatStats
	PetStats           combat.CombatStats
	PuppetStats        combat.CombatStats
	BuffStats          combat.CombatStats
	FinalStats         combat.CombatStats
	Warnings           []string
}

type Provider interface {
	GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error)
	GetEffectiveStatsBreakdown(ctx context.Context, userID string) (*StatsBreakdown, error)
}
