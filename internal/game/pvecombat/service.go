// File: internal/game/pvecombat/service.go
// Chức năng: Điều phối bắt đầu trận đấu PvE, kết nối giữa PvE world và Combat engine.

package pvecombat

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"go.uber.org/zap"
)

// StatsProvider cung cấp chỉ số người chơi từ Equipment hoặc Profile service.
type StatsProvider interface {
	GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error)
}

// PvEProvider cung cấp dữ liệu thế giới (Area, Encounter, Limits).
type PvEProvider interface {
	GetArea(areaID string) (pve.PvEAreaDefinition, error)
	GetNextStage(ctx context.Context, userID, areaID string) (int, error)
	CanEnterArea(ctx context.Context, userID string, area pve.PvEAreaDefinition, combatPower int64, realm string) error
	GenerateEncounter(area pve.PvEAreaDefinition, stage int, rng *rand.Rand) (*pve.EncounterDefinition, error)
	MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error
}

// RewardGrantService thao tác với Inventory, Cultivation, Economy để cấp đồ cho user.
type RewardGrantService interface {
	GrantExp(ctx context.Context, userID string, amount int64) error
	GrantStones(ctx context.Context, userID string, amount int64) error
	GrantItem(ctx context.Context, userID, defID string, quantity int64) error
	PreflightInventoryCapacity(ctx context.Context, userID string, items []RewardItemPlan) error
}

type RewardPlan struct {
	Exp     int64
	Stones  int64
	Items   []RewardItemPlan
	Summary []combat.ClaimedReward
}

type RewardItemPlan struct {
	ItemID   string
	Quantity int
}

type Service struct {
	repo          combat.Repository
	statsProvider StatsProvider
	pveProvider   PvEProvider
	grantService  RewardGrantService
	turnOrder     *combat.TurnOrderService
	rng           *rand.Rand
	now           func() time.Time
	log           *zap.Logger
}

func NewService(
	repo combat.Repository,
	statsProvider StatsProvider,
	pveProvider PvEProvider,
	grantService RewardGrantService,
	turnOrder *combat.TurnOrderService,
	rng *rand.Rand,
	log *zap.Logger,
) (*Service, error) {
	if repo == nil {
		return nil, errors.New("combat repo cannot be nil")
	}
	if statsProvider == nil {
		return nil, errors.New("statsProvider cannot be nil")
	}
	if pveProvider == nil {
		return nil, errors.New("pveProvider cannot be nil")
	}
	if grantService == nil {
		return nil, errors.New("grantService cannot be nil")
	}
	if turnOrder == nil {
		return nil, errors.New("turnOrder cannot be nil")
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if log == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Service{
		repo:          repo,
		statsProvider: statsProvider,
		pveProvider:   pveProvider,
		grantService:  grantService,
		turnOrder:     turnOrder,
		rng:           rng,
		now:           time.Now,
		log:           log.Named("game.pvecombat"),
	}, nil
}

func (s *Service) StartPvECombat(ctx context.Context, userID, areaID string) (*combat.CombatSession, error) {
	if userID == "" || areaID == "" {
		return nil, errors.New("userID và areaID không được rỗng")
	}

	// 1. Kiểm tra session đang hoạt động
	activeSession, err := s.repo.GetActiveSessionByUser(ctx, userID)
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, fmt.Errorf("lỗi kiểm tra trận đấu cũ: %w", err)
	}
	if activeSession != nil {
		// Tự động thanh tẩy Ghost Session từ phiên bản cũ (ID quá dài gây crash Discord UI)
		if len(activeSession.ID) > 25 {
			s.log.Warn("Phát hiện tàn trận có ID quá dài, tiến hành hủy bỏ để tạo trận mới", zap.String("old_id", activeSession.ID))
			activeSession.State = combat.StateLost
			_ = s.repo.UpdateSession(ctx, activeSession)
		} else {
			return activeSession, nil
		}
	}

	// 2. Lấy chỉ số chiến đấu
	s.log.Debug("StartPvECombat: gọi GetEffectiveStats", zap.String("userId", userID))
	stats, err := s.statsProvider.GetEffectiveStats(ctx, userID)
	s.log.Debug("StartPvECombat: GetEffectiveStats xong", zap.String("userId", userID), zap.Error(err))
	if err != nil {
		return nil, fmt.Errorf("không thể lấy chỉ số chiến đấu: user=%s reason=%w", userID, err)
	}
	if stats.MaxHP <= 0 || stats.ATK <= 0 || stats.Speed <= 0 {
		debugInfo := fmt.Sprintf("invalid combat stats: user=%s hp=%d atk=%d def=%d speed=%d cp=%d", userID, stats.MaxHP, stats.ATK, stats.DEF, stats.Speed, stats.CombatPower)
		s.log.Warn("Từ chối tạo trận do chỉ số lỗi", zap.String("userId", userID), zap.Any("stats", stats))
		return nil, fmt.Errorf("%w: %s", combat.ErrInvalidCombatStats, debugInfo)
	}

	// 3. Lấy thông tin Area và điều kiện
	area, err := s.pveProvider.GetArea(areaID)
	if err != nil {
		return nil, combat.ErrAreaNotFound
	}

	if err := s.pveProvider.CanEnterArea(ctx, userID, area, 1000, ""); err != nil {
		return nil, err
	}

	// 4. Khởi tạo ải và kẻ địch
	stage, err := s.pveProvider.GetNextStage(ctx, userID, areaID)
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo ải: %w", err)
	}
	encounter, err := s.pveProvider.GenerateEncounter(area, stage, s.rng)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo kẻ địch: %w", err)
	}
	if len(encounter.Enemies) == 0 {
		return nil, combat.ErrEncounterEmpty
	}
	if len(encounter.Enemies) > pve.MaxEnemiesPerEncounter {
		return nil, combat.ErrEnemyLimitExceeded
	}

	// 5. Khởi tạo tác nhân người chơi
	player := combat.CombatActor{
		ID:        userID,
		Type:      combat.ActorTypePlayer,
		Name:      "Đạo Hữu",
		Level:     1,
		Stats:     stats,
		CurrentHP: stats.MaxHP,
	}

	var enemies []combat.CombatActor
	for _, e := range encounter.Enemies {
		enemies = append(enemies, combat.CombatActor{
			ID:    e.ID,
			Type:  combat.ActorTypeMonster,
			Name:  e.Name,
			Level: e.Level,
			Stats: combat.CombatStats{
				MaxHP: e.Stats.MaxHP,
				ATK:   e.Stats.ATK,
				DEF:   e.Stats.DEF,
				Speed: e.Stats.Speed,
			},
			CurrentHP: e.Stats.MaxHP,
		})
	}

	// 6. Tính toán thứ tự lượt đánh
	actors := make([]combat.CombatActor, 0, 1+len(enemies))
	actors = append(actors, player)
	actors = append(actors, enemies...)

	turnOrder := s.turnOrder.BuildInitialOrder(actors)
	var currentActorID string
	if len(turnOrder) > 0 {
		currentActorID = turnOrder[0].ActorID
	}

	// 7. Tạo Session
	now := s.now().UTC()
	// Rút ngắn ID bằng Hex (từ 43 -> 19 ký tự), tránh lỗi vượt quá 100 ký tự CustomID của Discord UI
	sessionID := fmt.Sprintf("ss_%x", now.UnixNano())

	session := &combat.CombatSession{
		ID:                     sessionID,
		UserID:                 userID,
		AreaID:                 area.ID,
		ActivityType:           string(area.ActivityType),
		Stage:                  stage,
		State:                  combat.StateActive,
		Turn:                   1,
		Player:                 player,
		Enemies:                enemies,
		CurrentActorID:         currentActorID,
		TurnOrder:              turnOrder,
		GuaranteedRewardPoolID: encounter.GuaranteedRewardPoolID,
		BonusRewardPoolID:      encounter.BonusRewardPoolID,
		Logs: []combat.CombatLogEntry{
			{
				Turn:      1,
				ActorID:   "system",
				Action:    "start",
				Message:   fmt.Sprintf("Đạo hữu bước vào %s, linh khí dao động dữ dội...", area.Name),
				CreatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(30 * time.Minute),
	}

	// 8. Lưu DB
	if err := s.repo.CreateSession(ctx, session); err != nil {
		if errors.Is(err, combat.ErrCombatSessionAlreadyActive) {
			return s.repo.GetActiveSessionByUser(ctx, userID)
		}
		return nil, err
	}

	s.log.Info("Tạo trận PvE thành công",
		zap.String("userId", userID),
		zap.String("areaId", areaID),
		zap.String("activityType", string(area.ActivityType)),
		zap.Int("stage", stage),
		zap.String("sessionId", sessionID),
		zap.Int("enemyCount", len(enemies)),
	)
	return session, nil
}

func (s *Service) buildRewardPlan(session *combat.CombatSession) (*RewardPlan, error) {
	rollResult := pve.ResolveStageRewards(session.GuaranteedRewardPoolID, session.BonusRewardPoolID, s.rng)
	plan := &RewardPlan{}

	process := func(rewards []pve.ResolvedReward, isBonus bool) {
		for _, r := range rewards {
			if r.Quantity <= 0 {
				continue
			}
			switch r.Type {
			case "exp":
				plan.Exp += r.Quantity
			case "stones":
				plan.Stones += r.Quantity
			default:
				plan.Items = append(plan.Items, RewardItemPlan{ItemID: r.RefID, Quantity: int(r.Quantity)})
			}
			plan.Summary = append(plan.Summary, combat.ClaimedReward{Type: r.Type, RefID: r.RefID, Quantity: r.Quantity, IsBonus: isBonus})
		}
	}
	process(rollResult.Guaranteed, false)
	process(rollResult.Bonus, true)
	return plan, nil
}

func (s *Service) validateRewardPlan(plan *RewardPlan) error {
	if plan.Exp < 0 || plan.Stones < 0 {
		return errors.New("config error: reward quantity cannot be negative")
	}
	for _, item := range plan.Items {
		if item.Quantity <= 0 {
			return errors.New("config error: item quantity must be positive")
		}
		if item.ItemID == "" {
			return errors.New("config error: itemID cannot be empty")
		}
	}
	return nil
}

func (s *Service) grantRewardPlan(ctx context.Context, userID string, plan *RewardPlan) error {
	if plan.Exp > 0 {
		if err := s.grantService.GrantExp(ctx, userID, plan.Exp); err != nil {
			return err
		}
	}
	if plan.Stones > 0 {
		if err := s.grantService.GrantStones(ctx, userID, plan.Stones); err != nil {
			return err
		}
	}
	for _, it := range plan.Items {
		if err := s.grantService.GrantItem(ctx, userID, it.ItemID, int64(it.Quantity)); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ClaimReward(ctx context.Context, userID, sessionID string) ([]combat.ClaimedReward, error) {
	if userID == "" || sessionID == "" {
		return nil, errors.New("thiếu tham số bắt buộc")
	}

	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, combat.ErrCombatSessionNotFound
		}
		return nil, err
	}

	if session.UserID != userID {
		return nil, combat.ErrCombatSessionForbidden
	}

	if session.State != combat.StateWon {
		return nil, combat.ErrRewardSessionNotWon
	}

	if session.RewardClaimed || session.RewardClaimStatus == "claimed" {
		return nil, combat.ErrRewardAlreadyClaimed
	}
	if session.RewardClaimStatus == "claiming" {
		return nil, combat.ErrRewardClaimInProgress
	}
	if session.RewardClaimStatus == "claim_failed" {
		return nil, combat.ErrRewardClaimFailedNeedsAdmin
	}

	plan, err := s.buildRewardPlan(session)
	if err != nil {
		return nil, err
	}
	if err := s.validateRewardPlan(plan); err != nil {
		return nil, err
	}

	if len(plan.Items) > 0 {
		if err := s.grantService.PreflightInventoryCapacity(ctx, userID, plan.Items); err != nil {
			return nil, err
		}
	}

	claimID := "pve:claim:" + sessionID
	now := s.now().UTC()
	lockedSession, err := s.repo.TryStartRewardClaim(ctx, sessionID, claimID, now)
	if err != nil {
		return nil, err
	}

	if err := s.grantRewardPlan(ctx, userID, plan); err != nil {
		_ = s.repo.FailRewardClaim(ctx, sessionID, claimID, err.Error(), time.Now().UTC())
		return nil, err
	}

	// 4. Update Progression
	_ = s.pveProvider.MarkStageCleared(ctx, userID, lockedSession.AreaID, lockedSession.Stage)

	if err := s.repo.CompleteRewardClaim(ctx, sessionID, claimID, plan.Summary, time.Now().UTC()); err != nil {
		s.log.Error("CRITICAL: reward granted but claim finalize failed", zap.String("sessionId", sessionID), zap.Error(err))
		return nil, fmt.Errorf("reward granted but claim finalize failed: %w", err)
	}

	s.log.Info("Nhận thưởng PvE thành công",
		zap.String("userId", userID),
		zap.String("sessionId", sessionID),
		zap.Int("stage", session.Stage),
		zap.Int("totalRewards", len(plan.Summary)),
	)
	return plan.Summary, nil
}
