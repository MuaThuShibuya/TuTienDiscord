// File: internal/game/pvecombat/service_test.go
package pvecombat

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"go.uber.org/zap"
)

// --- Fakes ---

type fakeRepo struct {
	sessions map[string]*combat.CombatSession
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{sessions: make(map[string]*combat.CombatSession)}
}

func (r *fakeRepo) CreateSession(ctx context.Context, session *combat.CombatSession) error {
	for _, s := range r.sessions {
		if s.UserID == session.UserID && s.State == combat.StateActive && s.ExpiresAt.After(time.Now().UTC()) {
			return combat.ErrCombatSessionAlreadyActive
		}
	}
	r.sessions[session.ID] = session
	return nil
}

func (r *fakeRepo) GetSession(ctx context.Context, sessionID string) (*combat.CombatSession, error) {
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return s, nil
}

func (r *fakeRepo) GetActiveSessionByUser(ctx context.Context, userID string) (*combat.CombatSession, error) {
	for _, s := range r.sessions {
		if s.UserID == userID && s.State == combat.StateActive && s.ExpiresAt.After(time.Now().UTC()) {
			return s, nil
		}
	}
	return nil, apperrors.ErrNotFound
}

func (r *fakeRepo) UpdateSession(ctx context.Context, session *combat.CombatSession) error {
	r.sessions[session.ID] = session
	return nil
}
func (r *fakeRepo) MarkSessionState(ctx context.Context, sessionID string, state combat.SessionState) error {
	return nil
}

func (r *fakeRepo) TryStartRewardClaim(ctx context.Context, sessionID string, claimID string, now time.Time) (*combat.CombatSession, error) {
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	if s.RewardClaimed || s.RewardClaimStatus == "claimed" {
		return nil, combat.ErrRewardAlreadyClaimed
	}
	if s.RewardClaimStatus == "claiming" {
		return nil, combat.ErrRewardClaimInProgress
	}
	if s.RewardClaimStatus == "claim_failed" {
		return nil, combat.ErrRewardClaimFailedNeedsAdmin
	}
	s.RewardClaimStatus = "claiming"
	s.RewardClaimID = claimID
	s.UpdatedAt = now
	return s, nil
}
func (r *fakeRepo) CompleteRewardClaim(ctx context.Context, sessionID string, claimID string, details []combat.ClaimedReward, now time.Time) error {
	s, _ := r.sessions[sessionID]
	s.RewardClaimed = true
	s.RewardClaimStatus = "claimed"
	s.ClaimedRewards = details
	return nil
}
func (r *fakeRepo) FailRewardClaim(ctx context.Context, sessionID string, claimID string, reason string, now time.Time) error {
	s, _ := r.sessions[sessionID]
	s.RewardClaimStatus = "claim_failed"
	s.RewardClaimError = reason
	return nil
}

type fakeStatsProvider struct {
	stats combat.CombatStats
	err   error
}

func (p *fakeStatsProvider) GetEffectiveStats(ctx context.Context, userID string) (combat.CombatStats, error) {
	return p.stats, p.err
}

type fakePvEProvider struct {
	area      pve.PvEAreaDefinition
	nextStage int
	encounter *pve.EncounterDefinition
	errArea   error
	errEnter  error
}

func (p *fakePvEProvider) GetArea(areaID string) (pve.PvEAreaDefinition, error) {
	if p.errArea != nil {
		return pve.PvEAreaDefinition{}, p.errArea
	}
	return p.area, nil
}
func (p *fakePvEProvider) GetNextStage(ctx context.Context, userID, areaID string) (int, error) {
	return p.nextStage, nil
}
func (p *fakePvEProvider) CanEnterArea(ctx context.Context, userID string, area pve.PvEAreaDefinition, combatPower int64, realm string) error {
	return p.errEnter
}
func (p *fakePvEProvider) GenerateEncounter(area pve.PvEAreaDefinition, stage int, rng *rand.Rand) (*pve.EncounterDefinition, error) {
	if p.encounter == nil {
		return &pve.EncounterDefinition{}, nil
	}
	return p.encounter, nil
}

func (p *fakePvEProvider) MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error {
	return nil
}

type fakeRewardGrantService struct {
	err error
}

func (f *fakeRewardGrantService) PreflightInventoryCapacity(ctx context.Context, userID string, items []RewardItemPlan) error {
	return f.err
}

func (f *fakeRewardGrantService) GrantExp(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid")
	}
	return f.err
}

func (f *fakeRewardGrantService) GrantStones(ctx context.Context, userID string, amount int64) error {
	return f.err
}

func (f *fakeRewardGrantService) GrantItem(ctx context.Context, userID, defID string, quantity int64) error {
	return f.err
}

func newTestService(
	repo combat.Repository,
	stats StatsProvider,
	pveProv PvEProvider,
	grantSvc RewardGrantService,
) *Service {
	if repo == nil {
		repo = newFakeRepo()
	}
	if stats == nil {
		stats = &fakeStatsProvider{}
	}
	if pveProv == nil {
		pveProv = &fakePvEProvider{}
	}
	if grantSvc == nil {
		grantSvc = &fakeRewardGrantService{}
	}

	svc, err := NewService(repo, stats, pveProv, grantSvc, combat.NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())
	if err != nil {
		panic(err)
	}
	return svc
}

// --- Tests ---

func TestStartPvECombat_Success(t *testing.T) {
	repo := newFakeRepo()
	stats := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 1000, ATK: 50, Speed: 120}}
	pveProv := &fakePvEProvider{
		area:      pve.PvEAreaDefinition{ID: "area_1", Name: "Rừng Trúc", ActivityType: pve.ActivityDuNgoan},
		nextStage: 1,
		encounter: &pve.EncounterDefinition{
			GuaranteedRewardPoolID: "pool_1",
			Enemies: []pve.Enemy{
				{ID: "e_1", Stats: pve.MonsterStats{MaxHP: 100}},
			},
		},
	}

	svc := newTestService(repo, stats, pveProv, nil)

	session, err := svc.StartPvECombat(context.Background(), "u1", "area_1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi, nhận %v", err)
	}

	if session.State != combat.StateActive {
		t.Errorf("Mong đợi StateActive, nhận %v", session.State)
	}
	if session.Player.CurrentHP != 1000 {
		t.Errorf("Máu player khởi tạo sai")
	}
	if session.ActivityType != string(pve.ActivityDuNgoan) {
		t.Errorf("Activity type lưu sai")
	}
	if session.GuaranteedRewardPoolID != "pool_1" {
		t.Errorf("Không copy được Reward Pool")
	}
}

func TestStartPvECombat_ReturnsExistingActiveSession(t *testing.T) {
	repo := newFakeRepo()
	repo.sessions["ss_old"] = &combat.CombatSession{
		ID:        "ss_old",
		UserID:    "u1",
		State:     combat.StateActive,
		ExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}

	svc := newTestService(repo, &fakeStatsProvider{}, &fakePvEProvider{}, nil)
	session, err := svc.StartPvECombat(context.Background(), "u1", "area_1")

	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if session.ID != "ss_old" {
		t.Errorf("Phải trả về phiên cũ thay vì tạo mới. Trả về: %v", session.ID)
	}
}

func TestStartPvECombat_RejectInvalidStats(t *testing.T) {
	repo := newFakeRepo()
	stats := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 0}}

	svc := newTestService(repo, stats, &fakePvEProvider{}, nil)
	_, err := svc.StartPvECombat(context.Background(), "u1", "area_1")

	if err == nil {
		t.Fatalf("Mong đợi lỗi invalid stats, nhưng không có lỗi")
	}

	if !errors.Is(err, combat.ErrInvalidCombatStats) {
		t.Fatalf("Mong đợi ErrInvalidCombatStats, nhận %v", err)
	}

	msg := err.Error()
	requiredParts := []string{"user=u1", "hp=0", "atk=0", "def=0", "speed=0", "cp=0"}
	for _, part := range requiredParts {
		if !strings.Contains(msg, part) {
			t.Errorf("Error thiếu debug context %q: %v", part, msg)
		}
	}
}

func TestStartPvECombat_RejectEmptyEncounter(t *testing.T) {
	repo := newFakeRepo()
	stats := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 100, ATK: 10, Speed: 90}}
	pveProv := &fakePvEProvider{
		encounter: &pve.EncounterDefinition{Enemies: []pve.Enemy{}},
	}

	svc := newTestService(repo, stats, pveProv, nil)
	_, err := svc.StartPvECombat(context.Background(), "u1", "area_1")

	if err != combat.ErrEncounterEmpty {
		t.Errorf("Mong đợi lỗi ErrEncounterEmpty, nhận %v", err)
	}
}

func TestStartPvECombat_EnemyLimit(t *testing.T) {
	repo := newFakeRepo()
	stats := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 100, ATK: 10, Speed: 90}}

	enemies := make([]pve.Enemy, 10)
	for i := 0; i < 10; i++ {
		enemies[i] = pve.Enemy{ID: "e"}
	}

	pveProv := &fakePvEProvider{
		encounter: &pve.EncounterDefinition{Enemies: enemies},
	}

	svc := newTestService(repo, stats, pveProv, nil)
	_, err := svc.StartPvECombat(context.Background(), "u1", "area_1")

	if err != combat.ErrEnemyLimitExceeded {
		t.Errorf("Mong đợi lỗi vượt quá số lượng địch, nhận %v", err)
	}
}

func TestStartPvECombat_StatsProviderErrorIncludesContext(t *testing.T) {
	repo := newFakeRepo()
	stats := &fakeStatsProvider{err: errors.New("missing aptitude")}
	svc := newTestService(repo, stats, &fakePvEProvider{}, nil)

	_, err := svc.StartPvECombat(context.Background(), "u1", "area_1")
	if err == nil {
		t.Fatal("Mong đợi lỗi nhưng nhận nil")
	}
	if !strings.Contains(err.Error(), "user=u1") || !strings.Contains(err.Error(), "missing aptitude") {
		t.Errorf("Error thiếu context, nhận: %v", err)
	}
}
