// File: internal/game/characterstats/service.go
package characterstats

import (
	"context"

	"github.com/whiskey/tu-tien-bot/internal/game/aptitude"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
)

type pipelineSvc struct {
	aptSvc   aptitude.Service
	cultSvc  cultivation.Service
	equipSvc equipment.Service
}

func NewPipelineService(aptSvc aptitude.Service, cultSvc cultivation.Service, equipSvc equipment.Service) Provider {
	return &pipelineSvc{aptSvc: aptSvc, cultSvc: cultSvc, equipSvc: equipSvc}
}

func (s *pipelineSvc) GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error) {
	final := combat.CombatStats{}

	// 1. Lấy thông tin Cảnh Giới
	cultProfile, err := s.cultSvc.GetProfile(ctx, userID, "")
	if err != nil {
		return final, err
	}
	realmBase := cultivation.GetRealmBaseStats(string(cultProfile.Realm), cultProfile.RealmLevel)

	// 2. Lấy Modifier từ Tư Chất
	_, aptDef, err := s.aptSvc.GetProfile(ctx, userID)
	if err != nil || aptDef == nil {
		fallback := aptitude.Registry["apt_pham_tu"]
		aptDef = &fallback
	}

	// 3. TÍNH TOÁN BASE STATS GỐC (Realm * Aptitude Growth + Aptitude Base)
	final.MaxHP = int64(float64(realmBase.MaxHP)*aptDef.GrowthStats.MaxHPMultiplier) + aptDef.BaseStats.MaxHP
	final.ATK = int64(float64(realmBase.ATK)*aptDef.GrowthStats.ATKMultiplier) + aptDef.BaseStats.ATK
	final.DEF = int64(float64(realmBase.DEF)*aptDef.GrowthStats.DEFMultiplier) + aptDef.BaseStats.DEF
	final.Speed = int64(float64(realmBase.Speed)*aptDef.GrowthStats.SpeedMultiplier) + aptDef.BaseStats.Speed
	final.CritRate = aptDef.BaseStats.CritRate
	final.CritDamage = aptDef.BaseStats.CritDamage

	// 4. Cộng dồn Đạo Lộ (DaoPath)
	if cultProfile.Path == cultivation.PathSword {
		final.ATK = int64(float64(final.ATK) * 1.05)
	} else if cultProfile.Path == cultivation.PathBody {
		final.MaxHP = int64(float64(final.MaxHP) * 1.10)
		final.DEF = int64(float64(final.DEF) * 1.05)
	}

	// 5. Cộng dồn Equipment
	eqStats, err := s.equipSvc.GetEffectiveStats(ctx, userID, "")
	if err == nil {
		final.MaxHP += eqStats.MaxHP
		final.ATK += eqStats.ATK
		final.DEF += eqStats.DEF
		final.CritRate += eqStats.CritRate
	}

	// 6. Normalize an toàn
	if final.MaxHP <= 0 {
		final.MaxHP = 1
	}
	if final.ATK < 0 {
		final.ATK = 0
	}
	if final.Speed <= 0 {
		final.Speed = 90
	}
	if final.CritRate > 1.0 {
		final.CritRate = 1.0
	}

	// 7. Cập nhật CombatPower
	final.CombatPower = final.ATK*2 + final.DEF*2 + final.MaxHP/10
	return final, nil
}
