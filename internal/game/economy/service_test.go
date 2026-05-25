// File: internal/game/economy/service_test.go
// Phiên bản: v0.1.1
// Mục đích: Unit test cho economy.Service dùng in-memory repository.
//           Test bao gồm cộng/trừ tài nguyên, guard số âm, và các edge case.
// Ghi chú: Không kết nối MongoDB thật. Fake data chỉ dùng trong test, không dùng trong runtime.

package economy_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- In-memory repository ---

type memEconomyRepo struct {
	mu      sync.Mutex
	wallets map[string]*economy.Wallet // key: userID+":"+guildID
}

func newMemEconomyRepo() *memEconomyRepo {
	return &memEconomyRepo{wallets: make(map[string]*economy.Wallet)}
}

func key(userID, guildID string) string {
	return userID + ":" + guildID
}

func (r *memEconomyRepo) FindByUserID(_ context.Context, userID, guildID string) (*economy.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.wallets[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w: userId=%s", apperrors.ErrNotFound, userID)
	}
	copy := *w
	return &copy, nil
}

func (r *memEconomyRepo) Upsert(_ context.Context, wallet *economy.Wallet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if wallet.CreatedAt.IsZero() {
		wallet.CreatedAt = time.Now().UTC()
	}
	wallet.UpdatedAt = time.Now().UTC()
	copy := *wallet
	r.wallets[key(wallet.UserID, wallet.GuildID)] = &copy
	return nil
}

func (r *memEconomyRepo) AdjustSpiritStones(_ context.Context, userID, guildID string, amount int64) (*economy.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.wallets[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	// Guard: không cho số dư âm khi trừ
	if amount < 0 && w.SpiritStones < -amount {
		return nil, fmt.Errorf("%w: spirit_stones", apperrors.ErrInsufficientFunds)
	}
	w.SpiritStones += amount
	w.UpdatedAt = time.Now().UTC()
	copy := *w
	return &copy, nil
}

func (r *memEconomyRepo) AdjustSpiritJades(_ context.Context, userID, guildID string, amount int64) (*economy.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.wallets[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	if amount < 0 && w.SpiritJades < -amount {
		return nil, fmt.Errorf("%w: spirit_jades", apperrors.ErrInsufficientFunds)
	}
	w.SpiritJades += amount
	w.UpdatedAt = time.Now().UTC()
	copy := *w
	return &copy, nil
}

func (r *memEconomyRepo) AdjustFateTickets(_ context.Context, userID, guildID string, amount int) (*economy.Wallet, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	w, ok := r.wallets[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	if amount < 0 && w.FateTickets < -amount {
		return nil, fmt.Errorf("%w: fate_tickets", apperrors.ErrInsufficientFunds)
	}
	w.FateTickets += amount
	w.UpdatedAt = time.Now().UTC()
	copy := *w
	return &copy, nil
}

// --- Test setup ---

func TestMain(m *testing.M) {
	// Format json + level error: im lặng khi test pass
	if err := logger.Init(logger.Options{Level: "error", Format: "json"}); err != nil {
		panic("logger init thất bại: " + err.Error())
	}
	m.Run()
}

// --- Tests ---

func TestGetOrCreate_NewWallet(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	wallet, err := svc.GetOrCreate(ctx, "user1", "guild1")
	if err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}
	// Linh thạch khởi đầu phải là 500
	if wallet.SpiritStones != 500 {
		t.Errorf("SpiritStones khởi đầu sai: muốn 500, có %d", wallet.SpiritStones)
	}
	// Vé cơ duyên khởi đầu phải là 3
	if wallet.FateTickets != 3 {
		t.Errorf("FateTickets khởi đầu sai: muốn 3, có %d", wallet.FateTickets)
	}
}

func TestEarnSpiritStones(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user2", "guild1")
	w, err := svc.EarnSpiritStones(ctx, "user2", "guild1", 100, "test")
	if err != nil {
		t.Fatalf("EarnSpiritStones thất bại: %v", err)
	}
	if w.SpiritStones != 600 { // 500 + 100
		t.Errorf("SpiritStones sai: muốn 600, có %d", w.SpiritStones)
	}
}

func TestSpendSpiritStones_Success(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user3", "guild1") // 500 linh thạch
	w, err := svc.SpendSpiritStones(ctx, "user3", "guild1", 200, "test")
	if err != nil {
		t.Fatalf("SpendSpiritStones thất bại: %v", err)
	}
	if w.SpiritStones != 300 {
		t.Errorf("SpiritStones sai: muốn 300, có %d", w.SpiritStones)
	}
}

func TestSpendSpiritStones_InsufficientFunds(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user4", "guild1") // 500 linh thạch
	_, err := svc.SpendSpiritStones(ctx, "user4", "guild1", 1000, "test")
	if !apperrors.IsInsufficientFunds(err) {
		t.Errorf("Phải trả về ErrInsufficientFunds, có: %v", err)
	}
}

func TestEarnSpend_InvalidAmount(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user5", "guild1")

	// amount = 0 phải trả lỗi
	if _, err := svc.EarnSpiritStones(ctx, "user5", "guild1", 0, "test"); err == nil {
		t.Error("amount=0 phải trả về lỗi")
	}
	// amount âm phải trả lỗi (dùng Earn với số âm)
	if _, err := svc.EarnSpiritStones(ctx, "user5", "guild1", -1, "test"); err == nil {
		t.Error("amount âm phải trả về lỗi")
	}
}

func TestSpendFateTickets_Insufficient(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user6", "guild1") // 3 vé
	_, err := svc.SpendFateTickets(ctx, "user6", "guild1", 5, "test")
	if !apperrors.IsInsufficientFunds(err) {
		t.Errorf("Phải trả về ErrInsufficientFunds khi hết vé, có: %v", err)
	}
}

func TestGetWallet_NotFound(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	_, err := svc.GetWallet(ctx, "ghost", "guild1")
	if !apperrors.IsNotFound(err) {
		t.Errorf("Phải trả về ErrNotFound, có: %v", err)
	}
}
