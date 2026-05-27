// File: internal/game/pvecombat/integration_test.go
// Chức năng: Kiểm tra toàn bộ vòng đời Combat PvE từ lúc Start -> Won -> Claim Reward -> Tiến trình.

package pvecombat

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"go.uber.org/zap"
)

// --- Fake Services Dành Cho Integration ---

type fakeGrantService struct {
	failNext   bool
	totalExp   int64
	totalStone int64
	items      map[string]int64
}

func newFakeGrant() *fakeGrantService {
	return &fakeGrantService{items: make(map[string]int64)}
}

func (f *fakeGrantService) GrantExp(ctx context.Context, userID string, amount int64) error {
	if f.failNext {
		return errors.New("mock db error")
	}
	f.totalExp += amount
	return nil
}
func (f *fakeGrantService) GrantStones(ctx context.Context, userID string, amount int64) error {
	if f.failNext {
		return errors.New("mock db error")
	}
	f.totalStone += amount
	return nil
}
func (f *fakeGrantService) GrantItem(ctx context.Context, userID, defID string, quantity int64) error {
	if f.failNext {
		return errors.New("mock inventory full")
	}
	f.items[defID] += quantity
	return nil
}

type integrationPvEProvider struct {
	clearedStage int
}

func (p *integrationPvEProvider) GetArea(areaID string) (pve.PvEAreaDefinition, error) {
	return pve.AreaRegistry["area_du_ngoan_rung_truc"], nil
}
func (p *integrationPvEProvider) GetNextStage(ctx context.Context, userID, areaID string) (int, error) {
	return 1, nil
}
func (p *integrationPvEProvider) CanEnterArea(ctx context.Context, userID string, area pve.PvEAreaDefinition, combatPower int64, realm string) error {
	return nil
}
func (p *integrationPvEProvider) GenerateEncounter(area pve.PvEAreaDefinition, stage int, rng *rand.Rand) (*pve.EncounterDefinition, error) {
	return &pve.EncounterDefinition{
		AreaID: area.ID, Stage: stage, ActivityType: area.ActivityType,
		GuaranteedRewardPoolID: "reward_du_ngoan_basic",
		Enemies:                []pve.Enemy{{ID: "e1", Name: "Thỏ", Level: 1, Stats: pve.MonsterStats{MaxHP: 100}}},
	}, nil
}
func (p *integrationPvEProvider) MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error {
	if stage > p.clearedStage {
		p.clearedStage = stage
	}
	return nil
}

// --- Setup Môi trường ---

func setupIntegrationEnv() (*Service, *fakeGrantService, *fakeRepo, *integrationPvEProvider) {
	repo := newFakeRepo()
	grantSvc := newFakeGrant()
	pveProv := &integrationPvEProvider{}
	statsProv := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 1000, Speed: 100}}

	svc, _ := NewService(repo, statsProv, pveProv, grantSvc, combat.NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())
	return svc, grantSvc, repo, pveProv
}

// --- T E S T S ---

func TestPvECombat_EndToEnd_WinClaimRewardProgress(t *testing.T) {
	svc, grantSvc, _, pveProv := setupIntegrationEnv()
	ctx := context.Background()

	// 1. Start Combat
	session, err := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	if err != nil {
		t.Fatalf("StartPvECombat failed: %v", err)
	}

	// 2. Mock thắng trận (Mô phỏng hàm PlayerBasicAttack của combat engine)
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	// 3. Claim Reward
	rewards, err := svc.ClaimReward(ctx, "u1", session.ID, "idempotency_1")
	if err != nil {
		t.Fatalf("ClaimReward failed: %v", err)
	}

	// 4. Xác minh
	if len(rewards) == 0 {
		t.Errorf("Phải có phần thưởng rơi ra")
	}
	if !session.RewardClaimed {
		t.Errorf("Trạng thái Session phải được đánh dấu là đã claim")
	}
	if grantSvc.totalExp == 0 && grantSvc.totalStone == 0 {
		t.Errorf("Grant service không nhận được tài nguyên")
	}
	if pveProv.clearedStage != 1 {
		t.Errorf("Tiến trình (Progress) chưa được cập nhật. Mong đợi 1, nhận %d", pveProv.clearedStage)
	}
}

func TestPvECombat_ClaimRewardTwice_NoDoubleGrant(t *testing.T) {
	svc, grantSvc, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	// Claim lần 1
	_, _ = svc.ClaimReward(ctx, "u1", session.ID, "idem_abc")
	expFirstClaim := grantSvc.totalExp

	// Claim lần 2 cùng key (Double click)
	rewards2, err := svc.ClaimReward(ctx, "u1", session.ID, "idem_abc")
	if err != nil {
		t.Fatalf("Idempotency sai, đáng lẽ phải trả kết quả cũ, nhưng lại lỗi: %v", err)
	}
	if len(rewards2) == 0 {
		t.Errorf("Phải trả về đúng mảng reward cũ")
	}
	if grantSvc.totalExp > expFirstClaim {
		t.Errorf("NGUY HIỂM: Tài nguyên bị cấp 2 lần (Double Grant)!")
	}

	// Claim lần 3 KHÁC key (Hacker cố thử)
	_, err = svc.ClaimReward(ctx, "u1", session.ID, "idem_hacker_key")
	if err != combat.ErrRewardAlreadyClaimed {
		t.Errorf("Mong đợi lỗi ErrRewardAlreadyClaimed, nhưng nhận %v", err)
	}
}

func TestPvECombat_CannotClaimBeforeWin(t *testing.T) {
	svc, grantSvc, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")

	// Session đang active, chưa đánh thắng
	_, err := svc.ClaimReward(ctx, "u1", session.ID, "idem_abc")
	if err != combat.ErrRewardSessionNotWon {
		t.Errorf("Mong đợi lỗi ErrRewardSessionNotWon, nhận %v", err)
	}
	if grantSvc.totalExp > 0 {
		t.Errorf("Quái chưa chết mà đã rớt đồ!")
	}
}

func TestPvECombat_CannotClaimLostSession(t *testing.T) {
	svc, _, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateLost
	svc.repo.UpdateSession(ctx, session)

	_, err := svc.ClaimReward(ctx, "u1", session.ID, "idem_abc")
	if err != combat.ErrRewardSessionNotWon {
		t.Errorf("Thua trận không được nhận thưởng")
	}
}

func TestPvECombat_RewardGrantFailure_DoesNotMarkClaimed(t *testing.T) {
	svc, grantSvc, repo, pveProv := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	// Giả lập túi đồ đầy / Lỗi DB
	grantSvc.failNext = true

	_, err := svc.ClaimReward(ctx, "u1", session.ID, "idem_fail")
	if err == nil {
		t.Fatalf("Đáng lẽ phải lỗi grant")
	}

	// Kiểm tra xem session có bị đánh dấu là đã nhận không
	updatedSession, _ := repo.GetSession(ctx, session.ID)
	if updatedSession.RewardClaimed {
		t.Errorf("Lỗi lớn: Grant lỗi nhưng Session vẫn bị đánh dấu là đã nhận (Mất đồ của user)")
	}
	if pveProv.clearedStage > 0 {
		t.Errorf("Grant lỗi thì không được lưu Progress")
	}

	// Thử lại lần 2 (Sau khi user dọn túi)
	grantSvc.failNext = false
	_, err = svc.ClaimReward(ctx, "u1", session.ID, "idem_fail")
	if err != nil {
		t.Fatalf("Lần 2 phải thành công, nhận %v", err)
	}
}

func TestPvECombat_ProgressOnlyIncreases(t *testing.T) {
	svc, _, _, pveProv := setupIntegrationEnv()
	ctx := context.Background()

	pveProv.clearedStage = 5 // User đã đánh tới ải 5

	// User farm lại ải 3
	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.Stage = 3
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	_, _ = svc.ClaimReward(ctx, "u1", session.ID, "idem_abc")

	if pveProv.clearedStage != 5 {
		t.Errorf("Farm lại ải cũ làm tụt tiến trình! Mong đợi 5, nhận %d", pveProv.clearedStage)
	}
}
