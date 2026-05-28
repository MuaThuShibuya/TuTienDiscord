// File: internal/game/characterstats/service.go
package characterstats

import (
	"context"
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/game/aptitude"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"go.uber.org/zap"
)

type pipelineSvc struct {
	aptSvc   aptitude.Service
	cultSvc  cultivation.Service
	equipSvc equipment.Service
	log      *zap.Logger
}

func NewPipelineService(aptSvc aptitude.Service, cultSvc cultivation.Service, equipSvc equipment.Service) Provider {
	return &pipelineSvc{
		aptSvc:   aptSvc,
		cultSvc:  cultSvc,
		equipSvc: equipSvc,
		log:      logger.L().Named("characterstats"),
	}
}

func (s *pipelineSvc) GetEffectiveStatsBreakdown(ctx context.Context, userID string) (*StatsBreakdown, error) {
	bd := &StatsBreakdown{UserID: userID, Warnings: []string{}}

	cultProfile, err := s.cultSvc.GetProfile(ctx, userID, "")
	if err != nil {
		bd.Warnings = append(bd.Warnings, fmt.Sprintf("Missing cultivation: %v", err))
		return bd, err
	}

	bd.RealmID = cultivation.NormalizeRealmID(string(cultProfile.Realm))
	bd.RealmLevel = cultProfile.RealmLevel
	if realmDef, ok := cultivation.RealmRegistry[bd.RealmID]; ok {
		bd.RealmName = realmDef.Name
	}
	bd.BaseFromRealm = cultivation.GetRealmBaseStats(bd.RealmID, bd.RealmLevel)

	_, aptDef, err := s.aptSvc.GetProfile(ctx, userID)
	if err != nil || aptDef == nil {
		bd.Warnings = append(bd.Warnings, "Missing aptitude, using fallback")
		fallback := aptitude.Registry["apt_pham_tu"]
		aptDef = &fallback
	}
	bd.AptitudeID = aptDef.ID
	bd.AptitudeName = aptDef.Name
	bd.BaseFromAptitude = combat.CombatStats{MaxHP: aptDef.BaseStats.MaxHP, ATK: aptDef.BaseStats.ATK, DEF: aptDef.BaseStats.DEF, Speed: aptDef.BaseStats.Speed}
	bd.GrowthFromAptitude = aptDef.GrowthStats.MaxHPMultiplier

	eqStats, err := s.equipSvc.GetEffectiveStats(ctx, userID, "")
	if err != nil {
		bd.Warnings = append(bd.Warnings, fmt.Sprintf("Equipment error: %v", err))
	} else {
		bd.EquipmentStats = combat.CombatStats{MaxHP: eqStats.MaxHP, ATK: eqStats.ATK, DEF: eqStats.DEF, CritRate: eqStats.CritRate}
	}

	final, _ := s.GetEffectiveStats(ctx, userID)
	bd.FinalStats = final
	return bd, nil
}

func (s *pipelineSvc) GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error) {
	s.log.Debug("GetEffectiveStats: start", zap.String("userId", userID))
	final := combat.CombatStats{}

	// 1. Lấy thông tin Cảnh Giới
	cultProfile, err := s.cultSvc.GetProfile(ctx, userID, "")
	if err != nil {
		s.log.Error("GetEffectiveStats: missing cultivation", zap.String("userId", userID), zap.Error(err))
		return final, fmt.Errorf("missing cultivation: user=%s err=%w", userID, err)
	}

	realmID := cultivation.NormalizeRealmID(string(cultProfile.Realm))
	if _, ok := cultivation.RealmRegistry[realmID]; !ok {
		s.log.Error("GetEffectiveStats: missing realm definition", zap.String("userId", userID), zap.String("realmID", realmID))
		return final, fmt.Errorf("missing realm definition: user=%s realmID=%s", userID, realmID)
	}
	realmBase := cultivation.GetRealmBaseStats(realmID, cultProfile.RealmLevel)

	// 2. Lấy Modifier từ Tư Chất
	_, aptDef, err := s.aptSvc.GetProfile(ctx, userID)
	if err != nil || aptDef == nil {
		s.log.Error("GetEffectiveStats: missing aptitude", zap.String("userId", userID), zap.Error(err))
		return final, fmt.Errorf("missing aptitude: user=%s err=%w", userID, err)
	}

	s.log.Debug("GetEffectiveStats: loaded profiles",
		zap.String("userId", userID),
		zap.String("aptitudeID", aptDef.ID),
		zap.String("realmID", realmID),
		zap.Int("realmLevel", cultProfile.RealmLevel),
	)

	// 3. TÍNH TOÁN BASE STATS GỐC (Realm * Aptitude Growth + Aptitude Base)
	final.MaxHP = int64(float64(realmBase.MaxHP)*aptDef.GrowthStats.MaxHPMultiplier) + aptDef.BaseStats.MaxHP
	final.ATK = int64(float64(realmBase.ATK)*aptDef.GrowthStats.ATKMultiplier) + aptDef.BaseStats.ATK
	final.DEF = int64(float64(realmBase.DEF)*aptDef.GrowthStats.DEFMultiplier) + aptDef.BaseStats.DEF
	final.Speed = int64(float64(realmBase.Speed)*aptDef.GrowthStats.SpeedMultiplier) + aptDef.BaseStats.Speed
	final.CritRate = aptDef.BaseStats.CritRate
	final.CritDamage = aptDef.BaseStats.CritDamage

	// 4. Cộng dồn Đạo Lộ (DaoPath)
	switch cultProfile.Path {
	case cultivation.PathSword:
		final.ATK = int64(float64(final.ATK) * 1.05)
	case cultivation.PathBody:
		final.MaxHP = int64(float64(final.MaxHP) * 1.10)
		final.DEF = int64(float64(final.DEF) * 1.05)
	}

	// 5. Cộng dồn Equipment
	eqStats, err := s.equipSvc.GetEffectiveStats(ctx, userID, "")
	if err != nil {
		s.log.Error("GetEffectiveStats: equipment error", zap.String("userId", userID), zap.Error(err))
		return final, fmt.Errorf("equipment error: user=%s err=%w", userID, err)
	}
	final.MaxHP += eqStats.MaxHP
	final.ATK += eqStats.ATK
	final.DEF += eqStats.DEF
	final.CritRate += eqStats.CritRate

	// TODO_SAFE: add skill passive stats
	// TODO_SAFE: add pet stats
	// TODO_SAFE: add puppet stats

	// 6. Validate & Normalize
	if final.MaxHP <= 0 || final.ATK <= 0 || final.Speed <= 0 {
		return final, fmt.Errorf("invalid stats: userID=%s aptitudeID=%s realmID=%s realmLevel=%d equipment summary=applied MaxHP=%d ATK=%d DEF=%d Speed=%d CombatPower=%d",
			userID, aptDef.ID, realmID, cultProfile.RealmLevel, final.MaxHP, final.ATK, final.DEF, final.Speed, final.CombatPower)
	}
	if final.CritRate > 1.0 {
		final.CritRate = 1.0
	}

	// 7. Cập nhật CombatPower
	final.CombatPower = final.ATK*2 + final.DEF*2 + final.MaxHP/10

	s.log.Debug("GetEffectiveStats: final", zap.Int64("MaxHP", final.MaxHP), zap.Int64("ATK", final.ATK), zap.Int64("DEF", final.DEF), zap.Int64("Speed", final.Speed), zap.Int64("CombatPower", final.CombatPower))
	return final, nil
}
