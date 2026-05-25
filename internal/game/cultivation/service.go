// File: internal/game/cultivation/service.go
// Phiên bản: v0.1.1
// Mục đích: Business logic cho hệ thống tu luyện — lấy hoặc khởi tạo hồ sơ.
// Ghi chú: Logic tĩnh tu, đột phá, chọn đạo lộ sẽ thêm vào v0.2.
//          Service không gọi Discord API, không build embed.

package cultivation

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cooldown"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

// Service định nghĩa các nghiệp vụ tu luyện.
type Service interface {
	// GetOrCreate lấy hồ sơ tu luyện hoặc tạo mới với giá trị mặc định.
	GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error)

	// GetProfile lấy hồ sơ tu luyện, trả về ErrNotFound nếu chưa có.
	GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error)

	Meditate(ctx context.Context, input CultivationActionInput) (*CultivationActionResult, error)
	Seclusion(ctx context.Context, input CultivationActionInput) (*CultivationActionResult, error)
	BodyTraining(ctx context.Context, input CultivationActionInput) (*CultivationActionResult, error)
	Breakthrough(ctx context.Context, input BreakthroughInput) (*BreakthroughResult, error)

	// ChoosePath cho phép người chơi chọn con đường tu tiên (chỉ được chọn 1 lần).
	ChoosePath(ctx context.Context, userID, guildID string, path CultivationPath) error
}

type cultivationService struct {
	repo        Repository
	cooldownSvc cooldown.Service
	economySvc  economy.Service
	log         *zap.Logger
}

// NewService tạo cultivation service.
func NewService(repo Repository, cdSvc cooldown.Service, ecoSvc economy.Service) Service {
	return &cultivationService{
		repo:        repo,
		cooldownSvc: cdSvc,
		economySvc:  ecoSvc,
		log:         logger.L().Named("cultivation.service"),
	}
}

func (s *cultivationService) GetOrCreate(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID, guildID)
	if err == nil {
		return profile, nil
	}
	if !apperrors.IsNotFound(err) {
		// Log rõ khả năng sai lệch schema (ví dụ string vs int)
		s.log.Error("GetOrCreate: lỗi DB (có thể do schema mismatch)", zap.String("userId", userID), zap.Error(err))
		return nil, fmt.Errorf("dữ liệu hồ sơ đang gặp sự cố, vui lòng liên hệ admin: %w", err)
	}

	// Tạo hồ sơ mới với giá trị khởi đầu
	newProfile := NewCultivationProfile(userID, guildID)
	if err := s.repo.Upsert(ctx, newProfile); err != nil {
		s.log.Error("GetOrCreate: không tạo được hồ sơ tu luyện",
			zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	s.log.Info("Hồ sơ tu luyện mới được tạo",
		zap.String("userId", userID), zap.String("guildId", guildID))
	return newProfile, nil
}

func (s *cultivationService) GetProfile(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	return s.repo.FindByUserID(ctx, userID, guildID)
}

// clamp giới hạn giá trị trong khoảng min-max.
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func clampF(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func (s *cultivationService) checkCooldown(ctx context.Context, userID, guildID string, act cooldown.Action) error {
	onCD, remaining := s.cooldownSvc.IsOnCooldown(ctx, userID, guildID, act)
	if onCD {
		return &apperrors.CooldownError{Action: string(act), Remaining: utils.FormatDuration(remaining)}
	}
	return nil
}

func (s *cultivationService) Meditate(ctx context.Context, in CultivationActionInput) (*CultivationActionResult, error) {
	if err := s.checkCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionMeditate); err != nil {
		return nil, err
	}

	prof, err := s.GetProfile(ctx, in.UserID, in.GuildID)
	if err != nil {
		return nil, err
	}

	cost := ApplyPathStaminaCostBonus(prof.Path, "meditate", 5)
	if prof.Stamina < cost {
		return nil, apperrors.ErrInsufficientStamina
	}

	expGained := int64(10 + prof.RealmLevel*5)
	if prof.MindState >= 70 {
		expGained = expGained * 110 / 100
	}
	expGained = ApplyPathExpBonus(prof.Path, "meditate", expGained)

	prof.Stamina -= cost
	prof.CultivationExp += expGained
	prof.MindState = clamp(prof.MindState+1, 0, 100)

	if err := s.repo.UpdateStats(ctx, prof); err != nil {
		return nil, err
	}
	_ = s.cooldownSvc.SetCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionMeditate, 5*time.Minute)

	return &CultivationActionResult{
		Action: "Tĩnh Tu", ExpGained: expGained, StaminaSpent: cost,
		NewCultivationExp: prof.CultivationExp, CultivationRequired: prof.CultivationExpRequired,
		NewStamina: prof.Stamina, NewMindState: prof.MindState,
		CooldownExpiresAt: in.Now.Add(5 * time.Minute),
		Message:           fmt.Sprintf("Tĩnh tu thành công, nhận %d tu vi.", expGained),
	}, nil
}

func (s *cultivationService) Seclusion(ctx context.Context, in CultivationActionInput) (*CultivationActionResult, error) {
	if err := s.checkCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionSeclusion); err != nil {
		return nil, err
	}

	prof, err := s.GetProfile(ctx, in.UserID, in.GuildID)
	if err != nil {
		return nil, err
	}

	cost := ApplyPathStaminaCostBonus(prof.Path, "seclusion", 20)
	if prof.Stamina < cost {
		return nil, apperrors.ErrInsufficientStamina
	}

	expGained := int64(50 + prof.RealmLevel*15)
	if prof.MindState >= 70 {
		expGained = expGained * 115 / 100
	}
	expGained = ApplyPathExpBonus(prof.Path, "seclusion", expGained)

	prof.Stamina -= cost
	prof.CultivationExp += expGained
	prof.MindState = clamp(prof.MindState+2, 0, 100)

	if err := s.repo.UpdateStats(ctx, prof); err != nil {
		return nil, err
	}
	_ = s.cooldownSvc.SetCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionSeclusion, 60*time.Minute)

	return &CultivationActionResult{
		Action: "Bế Quan", ExpGained: expGained, StaminaSpent: cost,
		NewCultivationExp: prof.CultivationExp, CultivationRequired: prof.CultivationExpRequired,
		NewStamina: prof.Stamina, NewMindState: prof.MindState,
		CooldownExpiresAt: in.Now.Add(60 * time.Minute),
		Message:           fmt.Sprintf("Bế quan thành công, nhận %d tu vi.", expGained),
	}, nil
}

func (s *cultivationService) BodyTraining(ctx context.Context, in CultivationActionInput) (*CultivationActionResult, error) {
	if err := s.checkCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionBodyTraining); err != nil {
		return nil, err
	}

	prof, err := s.GetProfile(ctx, in.UserID, in.GuildID)
	if err != nil {
		return nil, err
	}

	cost := ApplyPathStaminaCostBonus(prof.Path, "body_training", 10)
	if prof.Stamina < cost {
		return nil, apperrors.ErrInsufficientStamina
	}

	expGained := int64(5 + prof.RealmLevel*3)
	cpGained := int64(3 + prof.RealmLevel*2)
	cpGained = ApplyPathCombatPowerBonus(prof.Path, "body_training", cpGained)

	prof.Stamina -= cost
	prof.CultivationExp += expGained
	prof.CombatPower += cpGained

	if err := s.repo.UpdateStats(ctx, prof); err != nil {
		return nil, err
	}
	_ = s.cooldownSvc.SetCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionBodyTraining, 15*time.Minute)

	return &CultivationActionResult{
		Action: "Luyện Thể", ExpGained: expGained, CombatPowerGained: cpGained, StaminaSpent: cost,
		NewCultivationExp: prof.CultivationExp, CultivationRequired: prof.CultivationExpRequired,
		NewCombatPower: prof.CombatPower, NewStamina: prof.Stamina, NewMindState: prof.MindState,
		CooldownExpiresAt: in.Now.Add(15 * time.Minute),
		Message:           fmt.Sprintf("Luyện thể thành công, nhận %d tu vi và %d chiến lực.", expGained, cpGained),
	}, nil
}

func (s *cultivationService) Breakthrough(ctx context.Context, in BreakthroughInput) (*BreakthroughResult, error) {
	if err := s.checkCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionBreakthrough); err != nil {
		return nil, err
	}

	prof, err := s.GetProfile(ctx, in.UserID, in.GuildID)
	if err != nil {
		return nil, err
	}

	if prof.CultivationExp < prof.CultivationExpRequired {
		return nil, apperrors.ErrInsufficientCultivationExp
	}
	if prof.MindState < 50 {
		return nil, apperrors.ErrInsufficientMindState
	}

	cost := CalculateBreakthroughCost(prof.RealmLevel)
	baseRate := 0.60
	mindBonus := float64(prof.MindState-50) * 0.005
	realmPenalty := float64(prof.RealmLevel) * 0.01
	finalRate := baseRate + mindBonus - realmPenalty
	finalRate = ApplyPathBreakthroughRateBonus(prof.Path, finalRate)
	finalRate = clampF(finalRate, 0.20, 0.95) // Clamp lần cuối

	roll := in.Rand.Float64()
	success := roll <= finalRate

	res := &BreakthroughResult{
		Success: success, Rate: finalRate, Roll: roll, OldRealm: string(prof.Realm), OldRealmLevel: prof.RealmLevel,
	}

	if success {
		// 1. Trừ tiền atomic trước khi thăng cấp
		if _, err := s.economySvc.SpendSpiritStones(ctx, in.UserID, in.GuildID, cost, "breakthrough_success"); err != nil {
			return nil, err
		}

		nextRealm, nextLevel, adv, err := NextRealm(string(prof.Realm), prof.RealmLevel)
		if err != nil {
			return nil, err
		} // ErrMaxRealmReached

		prof.CultivationExp -= prof.CultivationExpRequired
		prof.Realm = Realm(nextRealm)
		prof.RealmLevel = nextLevel
		prof.CultivationExpRequired = CalculateNextExpRequired(nextRealm, nextLevel)
		prof.CombatPower += int64(100 * nextLevel)
		prof.MindState = clamp(prof.MindState-3, 0, 100)

		_ = s.cooldownSvc.SetCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionBreakthrough, 30*time.Minute)

		res.NewRealm = string(prof.Realm)
		res.NewRealmLevel = prof.RealmLevel
		res.AdvancedRealm = adv
		res.CostPaid = cost
		res.Message = "Đột phá thành công! Cảnh giới thăng cấp."
	} else {
		penaltyCost := cost / 2
		if _, err := s.economySvc.SpendSpiritStones(ctx, in.UserID, in.GuildID, penaltyCost, "breakthrough_fail"); err != nil {
			return nil, err
		}

		penalty := ApplyPathBreakthroughFailurePenalty(prof.Path, 5)
		prof.MindState = clamp(prof.MindState-penalty, 0, 100)
		_ = s.cooldownSvc.SetCooldown(ctx, in.UserID, in.GuildID, cooldown.ActionBreakthrough, 15*time.Minute)

		res.NewRealm = string(prof.Realm)
		res.NewRealmLevel = prof.RealmLevel
		res.CostPaid = penaltyCost
		res.Message = "Đột phá thất bại. Tiêu hao linh thạch và tổn thương tâm cảnh."
	}

	if err := s.repo.UpdateStats(ctx, prof); err != nil {
		return nil, err
	}

	res.NewCultivationExp = prof.CultivationExp
	res.NewCultivationRequired = prof.CultivationExpRequired
	res.NewMindState = prof.MindState

	return res, nil
}

func (s *cultivationService) ChoosePath(ctx context.Context, userID, guildID string, path CultivationPath) error {
	if !path.IsValid() {
		return apperrors.ErrInvalidInput
	}

	prof, err := s.GetProfile(ctx, userID, guildID)
	if err != nil {
		return err
	}
	if prof.Path != PathNone {
		return apperrors.ErrPathAlreadyChosen
	}

	prof.Path = path
	return s.repo.UpdateStats(ctx, prof)
}
