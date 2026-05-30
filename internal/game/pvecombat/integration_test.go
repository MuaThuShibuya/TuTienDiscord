// File: internal/game/pvecombat/integration_test.go
// Chức năng: Kiểm tra toàn bộ vòng đời Combat PvE từ lúc Start -> Won -> Claim Reward -> Tiến trình.

package pvecombat

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"go.uber.org/zap"
)

// --- Fake Services Dành Cho Integration ---

// fakeEconomyService giả lập service kinh tế để test.
type fakeEconomyService struct {
	economy.Service // Bọc Interface để thỏa mãn tự động các method không dùng tới
	wallets         map[string]*economy.Wallet
	failNext        bool
	mu              sync.Mutex
}

func (f *fakeEconomyService) EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.failNext {
		return nil, errors.New("mock db error")
	}
	if f.wallets == nil {
		f.wallets = make(map[string]*economy.Wallet)
	}
	if _, ok := f.wallets[userID]; !ok {
		f.wallets[userID] = &economy.Wallet{UserID: userID}
	}
	f.wallets[userID].SpiritStones += amount
	return f.wallets[userID], nil
}

// Các phương thức khác để thỏa mãn interface (không cần thiết cho test này)
func (f *fakeEconomyService) GetWallet(ctx context.Context, userID, guildID string) (*economy.Wallet, error) {
	return f.wallets[userID], nil
}
func (f *fakeEconomyService) GetOrCreate(ctx context.Context, userID, guildID string) (*economy.Wallet, error) {
	return f.EarnSpiritStones(ctx, userID, guildID, 0, "create")
}
func (f *fakeEconomyService) SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	return nil, nil
}

// fakeCultivationService giả lập service tu luyện.
type fakeCultivationService struct {
	cultivation.Service // Bọc Interface để thỏa mãn tự động các method không dùng tới
	exp                 map[string]int64
	failNext            bool
}

func (f *fakeCultivationService) AddExperience(ctx context.Context, userID, guildID string, amount int64) error {
	if f.failNext {
		return errors.New("mock db error")
	}
	if f.exp == nil {
		f.exp = make(map[string]int64)
	}
	f.exp[userID] += amount
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

func setupIntegrationEnv() (*Service, *fakeEconomyService, *fakeCultivationService, *fakeRepo, *integrationPvEProvider) {
	repo := newFakeRepo()

	// Khởi tạo các service giả lập
	ecoSvc := &fakeEconomyService{}
	cultSvc := &fakeCultivationService{}

	// Khởi tạo GrantAdapter với các service giả lập
	// inventory.Service có thể là nil vì không test grant item trong test này.
	grantSvc := NewGrantAdapter(nil, ecoSvc, cultSvc)

	pveProv := &integrationPvEProvider{}
	statsProv := &fakeStatsProvider{stats: combat.CombatStats{MaxHP: 1000, ATK: 50, Speed: 100}}

	svc, _ := NewService(repo, statsProv, pveProv, grantSvc, combat.NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())
	return svc, ecoSvc, cultSvc, repo, pveProv
}

// --- T E S T S ---

func TestPvECombat_EndToEnd_WinClaimRewardProgress(t *testing.T) {
	svc, ecoSvc, cultSvc, _, pveProv := setupIntegrationEnv()
	ctx := context.Background()
	_ = cultSvc // use cultSvc

	// 1. Start Combat
	session, err := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	if err != nil {
		t.Fatalf("StartPvECombat failed: %v", err)
	}

	// 2. Mock thắng trận (Mô phỏng hàm PlayerBasicAttack của combat engine)
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	// 3. Claim Reward
	rewards, err := svc.ClaimReward(ctx, "u1", session.ID)
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
	wallet, _ := ecoSvc.GetWallet(ctx, "u1", "")
	if wallet == nil || wallet.SpiritStones == 0 {
		t.Errorf("Economy service không ghi nhận linh thạch")
	}
	if cultSvc.exp["u1"] == 0 {
		t.Errorf("Cultivation service không ghi nhận tu vi (exp)")
	}
	if pveProv.clearedStage != 1 {
		t.Errorf("Tiến trình (Progress) chưa được cập nhật. Mong đợi 1, nhận %d", pveProv.clearedStage)
	}
}

func TestPvECombat_ConcurrentDoubleClaim_NoDoubleGrant(t *testing.T) {
	svc, _, _, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	var wg sync.WaitGroup
	var successCount int32

	// Bắn 100 goroutine đồng thời
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.ClaimReward(ctx, "u1", session.ID)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Errorf("Lỗ hổng Race Condition: Có %d request ăn được giải, mong đợi đúng 1", successCount)
	}
}

func TestPvECombat_CannotClaimBeforeWin(t *testing.T) {
	svc, ecoSvc, _, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")

	// Session đang active, chưa đánh thắng
	_, err := svc.ClaimReward(ctx, "u1", session.ID)
	if err != combat.ErrRewardSessionNotWon {
		t.Errorf("Mong đợi lỗi ErrRewardSessionNotWon, nhận %v", err)
	}
	if ecoSvc.wallets["u1"] != nil && ecoSvc.wallets["u1"].SpiritStones > 0 {
		t.Errorf("Quái chưa chết mà đã rớt đồ!")
	}
}

func TestPvECombat_CannotClaimLostSession(t *testing.T) {
	svc, ecoSvc, _, _, _ := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateLost
	svc.repo.UpdateSession(ctx, session)

	_, err := svc.ClaimReward(ctx, "u1", session.ID)
	if err != combat.ErrRewardSessionNotWon {
		t.Errorf("Thua trận không được nhận thưởng")
	}
	if ecoSvc.wallets["u1"] != nil && ecoSvc.wallets["u1"].SpiritStones > 0 {
		t.Errorf("Thua trận không được phép nhận thưởng exp")
	}
}

func TestPvECombat_RewardGrantFailure_LocksSession(t *testing.T) {
	svc, ecoSvc, _, repo, pveProv := setupIntegrationEnv()
	ctx := context.Background()

	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	// Giả lập Lỗi DB trong lúc grant (sau khi lock)
	ecoSvc.failNext = true

	_, err := svc.ClaimReward(ctx, "u1", session.ID)
	if err == nil {
		t.Fatalf("Đáng lẽ phải lỗi grant")
	}

	// Kiểm tra xem session có bị đánh dấu là đã nhận không
	updatedSession, _ := repo.GetSession(ctx, session.ID)
	if updatedSession.RewardClaimed {
		t.Errorf("Lỗi lớn: Grant lỗi nhưng Session vẫn bị đánh dấu là đã nhận (Mất đồ của user)")
	}
	if updatedSession.RewardClaimStatus != "claim_failed" {
		t.Errorf("Mong đợi trạng thái claim_failed, nhận %s", updatedSession.RewardClaimStatus)
	}
	if pveProv.clearedStage > 0 {
		t.Errorf("Grant lỗi thì không được lưu Progress")
	}

	// Thử lại lần 2 phải bị khóa (Hard Fail)
	ecoSvc.failNext = false
	_, err = svc.ClaimReward(ctx, "u1", session.ID)
	if err != combat.ErrRewardClaimFailedNeedsAdmin {
		t.Fatalf("Lần 2 phải báo lỗi khóa an toàn (cần admin), nhận %v", err)
	}
}

func TestPvECombat_ProgressOnlyIncreases(t *testing.T) {
	svc, _, _, _, pveProv := setupIntegrationEnv()
	ctx := context.Background()

	pveProv.clearedStage = 5 // User đã đánh tới ải 5

	// User farm lại ải 3
	session, _ := svc.StartPvECombat(ctx, "u1", "area_du_ngoan_rung_truc")
	session.Stage = 3
	session.State = combat.StateWon
	svc.repo.UpdateSession(ctx, session)

	_, _ = svc.ClaimReward(ctx, "u1", session.ID)

	if pveProv.clearedStage != 5 {
		t.Errorf("Farm lại ải cũ làm tụt tiến trình! Mong đợi 5, nhận %d", pveProv.clearedStage)
	}
}
