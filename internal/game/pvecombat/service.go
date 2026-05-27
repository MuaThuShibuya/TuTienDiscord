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
}

type Service struct {
	repo          combat.Repository
	statsProvider StatsProvider
	pveProvider   PvEProvider
	grantService  RewardGrantService
	turnOrder     *combat.TurnOrderService
	rng           *rand.Rand
	now           func() time.Time
}

func NewService(
	repo combat.Repository,
	statsProvider StatsProvider,
	pveProvider PvEProvider,
	grantService RewardGrantService,
	turnOrder *combat.TurnOrderService,
	rng *rand.Rand,
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

	return &Service{
		repo:          repo,
		statsProvider: statsProvider,
		pveProvider:   pveProvider,
		grantService:  grantService,
		turnOrder:     turnOrder,
		rng:           rng,
		now:           time.Now,
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
		return activeSession, nil
	}

	// 2. Lấy chỉ số chiến đấu
	stats, err := s.statsProvider.GetEffectiveStats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy chỉ số chiến đấu: %w", err)
	}
	if stats.MaxHP <= 0 {
		return nil, combat.ErrInvalidCombatStats
	}
	if stats.Speed <= 0 {
		stats.Speed = 100
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
	sessionID := fmt.Sprintf("ss_%s_%d", userID, now.UnixNano())

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

	return session, nil
}

// ClaimReward giải quyết việc trao phần thưởng khi trận chiến kết thúc (State = Won).
func (s *Service) ClaimReward(ctx context.Context, userID, sessionID, idempotencyKey string) ([]combat.ClaimedReward, error) {
	if userID == "" || sessionID == "" || idempotencyKey == "" {
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

	// 1. Idempotency Check: Nếu đã nhận rồi, trả về kết quả cũ
	if session.RewardClaimed {
		if session.RewardIdempotencyKey == idempotencyKey {
			return session.ClaimedRewards, nil
		}
		return nil, combat.ErrRewardAlreadyClaimed
	}

	// 2. Resolve (Đổ xúc xắc mảng thưởng)
	rollResult := pve.ResolveStageRewards(session.GuaranteedRewardPoolID, session.BonusRewardPoolID, s.rng)

	var claimed []combat.ClaimedReward

	// Helper function để duyệt và grant
	processRewards := func(rewards []pve.ResolvedReward, isBonus bool) error {
		for _, r := range rewards {
			switch r.Type {
			case "exp":
				if err := s.grantService.GrantExp(ctx, userID, r.Quantity); err != nil {
					return err
				}
			case "stones":
				if err := s.grantService.GrantStones(ctx, userID, r.Quantity); err != nil {
					return err
				}
			default:
				if err := s.grantService.GrantItem(ctx, userID, r.RefID, r.Quantity); err != nil {
					return err
				}
			}
			claimed = append(claimed, combat.ClaimedReward{
				Type: r.Type, RefID: r.RefID, Quantity: r.Quantity, IsBonus: isBonus,
			})
		}
		return nil
	}

	// 3. Thực thi Grant
	if err := processRewards(rollResult.Guaranteed, false); err != nil {
		return nil, fmt.Errorf("%w: %v", combat.ErrRewardGrantFailed, err)
	}
	if err := processRewards(rollResult.Bonus, true); err != nil {
		return nil, fmt.Errorf("%w: %v", combat.ErrRewardGrantFailed, err)
	}

	// 4. Update Progression
	_ = s.pveProvider.MarkStageCleared(ctx, userID, session.AreaID, session.Stage)

	// 5. Lưu trạng thái Session
	session.RewardClaimed = true
	session.RewardClaimedAt = s.now().UTC()
	session.RewardIdempotencyKey = idempotencyKey
	session.ClaimedRewards = claimed

	if err := s.repo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	return claimed, nil
}
