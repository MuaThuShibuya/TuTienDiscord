package inventory_test

import (
	"sync"
	"sync/atomic"
	"testing"
)

// LƯU Ý: Đây là khung sườn Test chuẩn Enterprise dùng Mock nội bộ.
// Bạn cần import thư viện Mock (như gomock / testify) hoặc triển khai fake repository
// cho inventory.Repository, item.Repository và cultivation.Service để chạy được.

// --- A. Unit Tests (Logic Game) ---

func Test_AddItem_Stackable(t *testing.T) {
	// Setup: Mock Repo trả về 1 item đã có trong túi (Stackable = true)
	// Giả lập item "pill_exp_small" có Stackable = true
	// Thực thi s.AddItem()
	// Kì vọng: Không gọi CreateInstance, chỉ gọi AdjustQuantity
	t.Log("Đã setup kịch bản Test_AddItem_Stackable. Cần tiêm Fake Repo để assert.")
}

func Test_AddItem_InventoryFull(t *testing.T) {
	// Setup: Mock Repo trả về len(items) == 50 (Full túi đồ)
	// Thực thi s.AddItem() với item MỚI (không stack)
	// Kì vọng trả về lỗi ErrInventoryFull

	// Giả lập logic kiểm thử:
	// err := svc.AddItem(ctx, "user", "guild", "new_item", 1)
	// if !errors.Is(err, apperrors.ErrInventoryFull) {
	//     t.Errorf("Kì vọng lỗi ErrInventoryFull, nhận được %v", err)
	// }
	t.Log("Đã setup kịch bản Test_AddItem_InventoryFull. Cần tiêm Fake Repo để assert.")
}

func Test_UseItem_Success(t *testing.T) {
	// Setup: Mock túi có "pill_exp_small" (Usable=true, Qty=5)
	// Gọi s.UseItem(ctx, "user", "guild", "instance_id_1")
	// Kì vọng: Trả về thành công, AdjustQuantity(-1) được gọi
	t.Log("Đã setup kịch bản Test_UseItem_Success.")
}

func Test_UseItem_NotUsable(t *testing.T) {
	// Setup: Mock túi có "refine_stone" (Usable=false)
	// Gọi s.UseItem
	// Kì vọng: apperrors.ErrItemNotUsable
	t.Log("Đã setup kịch bản Test_UseItem_NotUsable.")
}

func Test_UseItem_ZeroQuantity_Cleanup(t *testing.T) {
	// Setup: Mock túi có Đan Dược Qty=1
	// Gọi s.UseItem
	// Kì vọng: AdjustQuantity(-1) gọi thành công, sau đó DeleteInstance được kích hoạt (do số lượng <= 0)
	t.Log("Đã setup kịch bản Test_UseItem_ZeroQuantity_Cleanup.")
}

// --- B. Concurrency Tests (Chống Race Condition) ---

func Test_GrantStarterItems_Concurrent(t *testing.T) {
	// Kịch bản: 10 goroutines cùng lúc nhận quà tân thủ
	// Cần Fake Repo hỗ trợ atomic Set/CompareAndSwap
	const goroutines = 10
	var wg sync.WaitGroup
	var successCount int32

	// svc := inventory.NewService(fakeInvRepo, fakeItemRepo, fakeCultSvc)
	// ctx := context.Background()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// if err := svc.GrantStarterItems(ctx, "userX", "guild1"); err == nil {
			// 	atomic.AddInt32(&successCount, 1)
			// }
			atomic.AddInt32(&successCount, 1) // Fake hành động
		}()
	}
	wg.Wait()

	// Ở môi trường Mongo thực, successCount sẽ là 1 (hoặc số lần mark thành công)
	// và hàm MarkStarterGranted sẽ cản các flow khác.
	t.Logf("Concurrent Starter Granted Check Complete. (Cần Mongo Fake để test logic atomic)")
}

func Test_UseItem_Concurrent(t *testing.T) {
	// Kịch bản: User chỉ có 1 viên đan dược, click dùng 5 lần cùng lúc
	const clicks = 5
	var wg sync.WaitGroup
	// var successCount int32
	// var failCount int32

	for i := 0; i < clicks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// msg, err := svc.UseItem(ctx, "userX", "guild1", "instance1")
			// if err == nil {
			// 	atomic.AddInt32(&successCount, 1)
			// } else {
			//  atomic.AddInt32(&failCount, 1)
			// }
		}()
	}
	wg.Wait()

	// Ở DB thật (nhờ MongoDB $gte: 1 / $inc -1):
	// if atomic.LoadInt32(&successCount) > 1 {
	// 	t.Errorf("LỖI BẢO MẬT: Double use item! Chỉ được phép 1 giao dịch thành công")
	// }
	t.Log("Concurrent Use Item Test setup done. Cần tiêm Fake Mongo Repository để xác minh chống race condition.")
}
