// File: internal/game/economy/concurrent_test.go
// Phiên bản: v0.1.2
// Mục đích: Kiểm tra bảo mật và tính toàn vẹn dữ liệu khi nhiều goroutine đồng thời
//           thực hiện giao dịch trên cùng một ví — phòng chống double-spend và số dư âm.
// Bối cảnh bảo mật:
//   - Double-spend: user gửi nhiều request đồng thời để tiêu cùng một khoản tiền hai lần.
//   - Race condition: không có mutex → số dư có thể âm khi nhiều goroutine cùng đọc-ghi.
//   - In-memory repo dùng mutex nên test này xác nhận behavior đúng.
//   - MongoDB repo PHẢI dùng atomic conditional $inc (xem economy/repository_mongo.go).
// Ghi chú: Chạy với -race flag để phát hiện data race: go test -race ./internal/game/economy/...

package economy_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/economy"
)

// TestConcurrentSpend_NoDoublespend kiểm tra kịch bản double-spend:
// 3 goroutine đồng thời chi tiêu 300 linh thạch từ ví có 500.
// Chỉ 1 giao dịch được phép thành công (300 <= 500), 2 còn lại phải thất bại với ErrInsufficientFunds.
func TestConcurrentSpend_NoDoublespend(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	// Khởi tạo ví với 500 linh thạch (giá trị mặc định từ GetOrCreate)
	if _, err := svc.GetOrCreate(ctx, "userX", "guild1"); err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}

	const (
		goroutines  = 3
		spendAmount = 300 // 300 > 500/2 → chỉ có thể 1 lần thành công
	)

	var (
		wg           sync.WaitGroup
		successCount int32
		failCount    int32
	)

	// Phóng goroutines đồng thời
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.SpendSpiritStones(ctx, "userX", "guild1", spendAmount, "concurrent-test")
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				atomic.AddInt32(&failCount, 1)
			}
		}()
	}
	wg.Wait()

	// Tối đa 1 giao dịch thành công (500 / 300 = 1 lần)
	if got := atomic.LoadInt32(&successCount); got > 1 {
		t.Errorf("Double-spend! Có %d giao dịch thành công nhưng chỉ được phép 1 (ví có 500, mỗi lần tiêu 300)", got)
	}

	// Số dư KHÔNG được âm — đây là điều kiện an toàn tuyệt đối
	wallet, err := svc.GetWallet(ctx, "userX", "guild1")
	if err != nil {
		t.Fatalf("GetWallet thất bại: %v", err)
	}
	if wallet.SpiritStones < 0 {
		t.Errorf("LỖ HỔNG NGHIÊM TRỌNG: số dư âm %d (double-spend thành công!)", wallet.SpiritStones)
	}

	t.Logf("Kết quả: %d thành công, %d thất bại, số dư cuối: %d linh thạch",
		atomic.LoadInt32(&successCount), atomic.LoadInt32(&failCount), wallet.SpiritStones)
}

// TestConcurrentSpend_StressTest kiểm tra với 20 goroutine đồng thời chi tiêu 60 linh thạch
// từ ví có 500. Tối đa 8 giao dịch được phép thành công (500 / 60 = 8, dư 20).
func TestConcurrentSpend_StressTest(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	if _, err := svc.GetOrCreate(ctx, "userStress", "guild1"); err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}

	const (
		goroutines  = 20
		spendAmount = 60
		initialBal  = 500
		maxSuccess  = initialBal / spendAmount // = 8
	)

	var (
		wg           sync.WaitGroup
		successCount int32
	)

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := svc.SpendSpiritStones(ctx, "userStress", "guild1", spendAmount, "stress-test"); err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	// Không vượt quá số lần tối đa có thể thành công
	if got := atomic.LoadInt32(&successCount); got > maxSuccess {
		t.Errorf("Quá nhiều giao dịch thành công: %d (tối đa %d với số dư %d, mỗi lần tiêu %d)",
			got, maxSuccess, initialBal, spendAmount)
	}

	// Số dư không âm
	wallet, err := svc.GetWallet(ctx, "userStress", "guild1")
	if err != nil {
		t.Fatalf("GetWallet thất bại: %v", err)
	}
	if wallet.SpiritStones < 0 {
		t.Errorf("LỖ HỔNG NGHIÊM TRỌNG: số dư âm %d sau stress test", wallet.SpiritStones)
	}

	t.Logf("Stress test: %d/%d goroutine thành công, số dư cuối: %d",
		atomic.LoadInt32(&successCount), goroutines, wallet.SpiritStones)
}

// TestConcurrentEarnAndSpend kiểm tra cùng lúc cộng và trừ linh thạch.
// 10 goroutine earn 50, 10 goroutine spend 50 — số dư cuối phải nhất quán.
func TestConcurrentEarnAndSpend(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	// Ví bắt đầu với 500 linh thạch
	if _, err := svc.GetOrCreate(ctx, "userMixed", "guild1"); err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}

	const (
		earnGoroutines  = 10
		spendGoroutines = 10
		amount          = 50
	)

	var (
		wg           sync.WaitGroup
		earnSuccess  int32
		spendSuccess int32
	)

	// Tung đồng thời earn và spend
	for i := 0; i < earnGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := svc.EarnSpiritStones(ctx, "userMixed", "guild1", amount, "earn-concurrent"); err == nil {
				atomic.AddInt32(&earnSuccess, 1)
			}
		}()
	}
	for i := 0; i < spendGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := svc.SpendSpiritStones(ctx, "userMixed", "guild1", amount, "spend-concurrent"); err == nil {
				atomic.AddInt32(&spendSuccess, 1)
			}
		}()
	}
	wg.Wait()

	// Số dư cuối = 500 + earnSuccess*50 - spendSuccess*50
	wallet, err := svc.GetWallet(ctx, "userMixed", "guild1")
	if err != nil {
		t.Fatalf("GetWallet thất bại: %v", err)
	}
	if wallet.SpiritStones < 0 {
		t.Errorf("LỖ HỔNG: số dư âm %d sau concurrent earn+spend", wallet.SpiritStones)
	}

	// Kiểm tra tính nhất quán: số dư phải khớp với số giao dịch thành công
	expected := int64(500) + int64(atomic.LoadInt32(&earnSuccess))*amount - int64(atomic.LoadInt32(&spendSuccess))*amount
	if wallet.SpiritStones != expected {
		t.Errorf("Số dư không nhất quán: expect %d, got %d (earn=%d spend=%d)",
			expected, wallet.SpiritStones,
			atomic.LoadInt32(&earnSuccess), atomic.LoadInt32(&spendSuccess))
	}

	t.Logf("Earn: %d/%d OK, Spend: %d/%d OK, số dư cuối: %d (expect %d)",
		atomic.LoadInt32(&earnSuccess), earnGoroutines,
		atomic.LoadInt32(&spendSuccess), spendGoroutines,
		wallet.SpiritStones, expected)
}

// TestConcurrentMultiUser kiểm tra isolation giữa các user:
// Giao dịch của user A không ảnh hưởng đến số dư của user B.
func TestConcurrentMultiUser(t *testing.T) {
	svc := economy.NewService(newMemEconomyRepo())
	ctx := context.Background()

	users := []string{"userA", "userB", "userC"}
	for _, u := range users {
		if _, err := svc.GetOrCreate(ctx, u, "guild1"); err != nil {
			t.Fatalf("GetOrCreate thất bại cho %s: %v", u, err)
		}
	}

	var wg sync.WaitGroup

	// Mỗi user chi tiêu 200 linh thạch trong goroutine riêng
	for _, u := range users {
		wg.Add(1)
		u := u // capture
		go func() {
			defer wg.Done()
			svc.SpendSpiritStones(ctx, u, "guild1", 200, "multi-user-test")
		}()
	}
	wg.Wait()

	// Kiểm tra từng user độc lập — không bị ảnh hưởng lẫn nhau
	for _, u := range users {
		w, err := svc.GetWallet(ctx, u, "guild1")
		if err != nil {
			t.Errorf("GetWallet thất bại cho %s: %v", u, err)
			continue
		}
		if w.SpiritStones < 0 {
			t.Errorf("User %s bị số dư âm: %d", u, w.SpiritStones)
		}
		// Sau khi spend 200 từ 500: phải còn đúng 300
		if w.SpiritStones != 300 {
			t.Errorf("User %s: expect 300 linh thạch, got %d", u, w.SpiritStones)
		}
	}
}
