// File: internal/game/pve/encounter.go
// Chức năng: Sinh tổ đội quái vật cho một trận đánh (Encounter).

package pve

import (
	"fmt"
	"math/rand"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

type EncounterDefinition struct {
	AreaID                 string
	Stage                  int
	ActivityType           ActivityType
	Enemies                []combat.CombatActor
	RecommendedCombatPower int64
	GuaranteedRewardPoolID string
	BonusRewardPoolID      string
}

// GenerateEncounter sinh ra mảng kẻ địch cho trận chiến. Chặn tối đa 9 quái.
func GenerateEncounter(area PvEAreaDefinition, stage int, rng *rand.Rand) (*EncounterDefinition, error) {
	stageDef, err := GetStageDefinition(area, stage)
	if err != nil {
		return nil, err
	}

	pool, ok := MonsterPoolRegistry[stageDef.MonsterPoolID]
	if !ok || len(pool.MonsterIDs) == 0 {
		return nil, fmt.Errorf("không tìm thấy monster pool %s", stageDef.MonsterPoolID)
	}

	// Tính toán số lượng quái
	enemyCount := stageDef.EnemyCountMin
	if stageDef.EnemyCountMax > stageDef.EnemyCountMin {
		enemyCount += rng.Intn(stageDef.EnemyCountMax - stageDef.EnemyCountMin + 1)
	}
	if enemyCount > MaxEnemiesPerEncounter {
		enemyCount = MaxEnemiesPerEncounter
	}
	// Boss stage thường ít quái rác hơn
	if stageDef.IsBossStage && enemyCount > 3 {
		enemyCount = 3
	}

	var enemies []combat.CombatActor
	for i := 0; i < enemyCount; i++ {
		var mDef MonsterDefinition
		role := MonsterRoleNormal

		// Quyết định quái thường / Elite / Boss
		if stageDef.IsBossStage && i == 0 && len(pool.BossMonsterIDs) > 0 {
			mDef = MonsterRegistry[pool.BossMonsterIDs[rng.Intn(len(pool.BossMonsterIDs))]]
			role = MonsterRoleBoss
		} else if rng.Float64() < 0.15 && len(pool.EliteMonsterIDs) > 0 { // 15% ra elite
			mDef = MonsterRegistry[pool.EliteMonsterIDs[rng.Intn(len(pool.EliteMonsterIDs))]]
			role = MonsterRoleElite
		} else {
			mDef = MonsterRegistry[pool.MonsterIDs[rng.Intn(len(pool.MonsterIDs))]]
		}

		scaledStats := ScaleMonsterStats(mDef.BaseStats, stage, area.ActivityType, role, DefaultScalingConfig)
		enemies = append(enemies, combat.CombatActor{
			ID:        fmt.Sprintf("e_%d_%s", i, mDef.ID),
			Type:      combat.ActorTypeMonster,
			Name:      mDef.Name,
			Level:     mDef.BaseLevel + stage,
			Stats:     scaledStats,
			CurrentHP: scaledStats.MaxHP,
		})
	}

	return &EncounterDefinition{AreaID: area.ID, Stage: stage, ActivityType: area.ActivityType, Enemies: enemies, GuaranteedRewardPoolID: stageDef.GuaranteedRewardPoolID, BonusRewardPoolID: stageDef.BonusRewardPoolID}, nil
}
