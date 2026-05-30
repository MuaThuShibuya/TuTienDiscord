// File: internal/game/shop/npc/service_test.go
// Chức năng: Kiểm tra tính chính xác của Service Mua/Bán (Chống thất thoát dữ liệu).

package npc_test

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/shop/npc"
	"go.uber.org/zap"
)

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"test_item_valid":    {ID: "test_item_valid"},
		"test_item_disabled": {ID: "test_item_disabled"},
		"test_item_nosell":   {ID: "test_item_nosell"},
	})

	npc.Registry["test_shop"] = npc.ShopDef{
		ID:   "test_shop",
		Name: "Test Shop",
		Items: map[string]npc.ItemDef{
			"test_item_valid":    {ItemID: "test_item_valid", BuyPrice: 100, SellPrice: 20, Enabled: true},
			"test_item_disabled": {ItemID: "test_item_disabled", BuyPrice: 100, SellPrice: 20, Enabled: false},
			"test_item_nosell":   {ItemID: "test_item_nosell", BuyPrice: 100, SellPrice: 0, Enabled: true},
		},
	}
}

// --- Mock Dependencies ---

type mockEconomy struct {
	economy.Service
	mu           sync.Mutex
	spent        int64
	earn         int64
	failNext     error
	insufficient bool
}

func (m *mockEconomy) SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.insufficient {
		return nil, apperrors.ErrInsufficientFunds
	}
	if m.failNext != nil {
		err := m.failNext
		m.failNext = nil
		return nil, err
	}
	m.spent += amount
	return &economy.Wallet{SpiritStones: 10000 - m.spent + m.earn}, nil
}

func (m *mockEconomy) EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failNext != nil {
		err := m.failNext
		m.failNext = nil
		return nil, err
	}
	m.earn += amount
	return &economy.Wallet{SpiritStones: 10000 - m.spent + m.earn}, nil
}

type mockInventory struct {
	inventory.Service
	mu      sync.Mutex
	added   int64
	removed int64
	failAdd bool
	failRem bool
}

func (m *mockInventory) AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failAdd {
		return apperrors.ErrInventoryFull
	}
	m.added += quantity
	return nil
}

func (m *mockInventory) RemoveItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failRem {
		return apperrors.ErrInsufficientItemQuantity
	}
	m.removed += quantity
	return nil
}

// --- Tests ---

func TestNPCService_BuyItem_Table(t *testing.T) {
	cases := []struct {
		name          string
		shopID        string
		itemID        string
		qty           int
		setupEco      func(*mockEconomy)
		setupInv      func(*mockInventory)
		expectErr     bool
		expectErrCode string
		expectSpent   int64
		expectAdded   int64
	}{
		{
			name:        "1. Mua hợp lệ 1 vật phẩm",
			shopID:      "test_shop",
			itemID:      "test_item_valid",
			qty:         1,
			expectErr:   false,
			expectSpent: 100,
			expectAdded: 1,
		},
		{
			name:        "2. Mua hợp lệ nhiều vật phẩm",
			shopID:      "test_shop",
			itemID:      "test_item_valid",
			qty:         5,
			expectErr:   false,
			expectSpent: 500,
			expectAdded: 5,
		},
		{
			name:          "3. Số lượng mua 0",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           0,
			expectErr:     true,
			expectErrCode: "INVALID_QTY",
		},
		{
			name:          "4. Số lượng mua âm",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           -5,
			expectErr:     true,
			expectErrCode: "INVALID_QTY",
		},
		{
			name:          "5. Thương hội không tồn tại",
			shopID:        "shop_ma",
			itemID:        "test_item_valid",
			qty:           1,
			expectErr:     true,
			expectErrCode: "SHOP_NOT_FOUND",
		},
		{
			name:          "6. Vật phẩm không bán tại thương hội",
			shopID:        "test_shop",
			itemID:        "item_ma",
			qty:           1,
			expectErr:     true,
			expectErrCode: "ITEM_NOT_FOR_SALE",
		},
		{
			name:          "7. Vật phẩm bị vô hiệu hóa (Enabled=false)",
			shopID:        "test_shop",
			itemID:        "test_item_disabled",
			qty:           1,
			expectErr:     true,
			expectErrCode: "ITEM_NOT_FOR_SALE",
		},
		{
			name:          "8. Không đủ linh thạch",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           1,
			setupEco:      func(m *mockEconomy) { m.insufficient = true },
			expectErr:     true,
			expectErrCode: "NO_MONEY",
			expectSpent:   0,
		},
		{
			name:        "9. Lỗi DB Economy bất ngờ",
			shopID:      "test_shop",
			itemID:      "test_item_valid",
			qty:         1,
			setupEco:    func(m *mockEconomy) { m.failNext = errors.New("db timeout") },
			expectErr:   true,
			expectSpent: 0,
		},
		{
			name:          "10. Túi đồ đầy (SAGA Rollback)",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           1,
			setupInv:      func(m *mockInventory) { m.failAdd = true },
			expectErr:     true,
			expectErrCode: "INV_FULL",
			expectSpent:   0, // Đã trừ 100, nhưng SAGA Rollback hoàn 100 -> Net Spent = 0
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ecoMock := &mockEconomy{}
			invMock := &mockInventory{}
			if tc.setupEco != nil {
				tc.setupEco(ecoMock)
			}
			if tc.setupInv != nil {
				tc.setupInv(invMock)
			}

			svc := npc.NewService(ecoMock, invMock, zap.NewNop())
			_, err := svc.BuyItem(context.Background(), "u1", "g1", tc.shopID, tc.itemID, tc.qty)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Mong đợi lỗi nhưng nhận nil")
				}
				if tc.expectErrCode != "" {
					var appErr *apperrors.AppError
					if errors.As(err, &appErr) {
						if appErr.Code != tc.expectErrCode {
							t.Errorf("Sai mã lỗi: mong đợi %s, nhận %s", tc.expectErrCode, appErr.Code)
						}
					} else {
						t.Errorf("Lỗi không phải là AppError: %v", err)
					}
				}
			} else if err != nil {
				t.Fatalf("Không mong đợi lỗi: %v", err)
			}

			netSpent := ecoMock.spent - ecoMock.earn
			if netSpent != tc.expectSpent {
				t.Errorf("Net Spent: mong đợi %d, nhận %d (spent: %d, earn: %d)", tc.expectSpent, netSpent, ecoMock.spent, ecoMock.earn)
			}
			if invMock.added != tc.expectAdded {
				t.Errorf("Added Items: mong đợi %d, nhận %d", tc.expectAdded, invMock.added)
			}
		})
	}
}

func TestNPCService_SellItem_Table(t *testing.T) {
	cases := []struct {
		name          string
		shopID        string
		itemID        string
		qty           int
		setupEco      func(*mockEconomy)
		setupInv      func(*mockInventory)
		expectErr     bool
		expectErrCode string
		expectEarn    int64
		expectRemoved int64
	}{
		{
			name:          "11. Bán hợp lệ 1 vật phẩm",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           1,
			expectErr:     false,
			expectEarn:    20,
			expectRemoved: 1,
		},
		{
			name:          "12. Bán hợp lệ nhiều vật phẩm",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           10,
			expectErr:     false,
			expectEarn:    200,
			expectRemoved: 10,
		},
		{
			name:          "13. Bán số lượng 0",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           0,
			expectErr:     true,
			expectErrCode: "INVALID_QTY",
		},
		{
			name:          "14. Bán số lượng âm",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           -2,
			expectErr:     true,
			expectErrCode: "INVALID_QTY",
		},
		{
			name:          "15. Thương hội không tồn tại",
			shopID:        "shop_ma",
			itemID:        "test_item_valid",
			qty:           1,
			expectErr:     true,
			expectErrCode: "SHOP_NOT_FOUND",
		},
		{
			name:          "16. Vật phẩm không thu mua",
			shopID:        "test_shop",
			itemID:        "item_khong_thu_mua",
			qty:           1,
			expectErr:     true,
			expectErrCode: "ITEM_NOT_BOUGHT",
		},
		{
			name:          "17. Vật phẩm có giá thu mua = 0",
			shopID:        "test_shop",
			itemID:        "test_item_nosell",
			qty:           1,
			expectErr:     true,
			expectErrCode: "ITEM_NO_VALUE",
		},
		{
			name:          "18. Thiếu vật phẩm trong túi (Inventory báo lỗi)",
			shopID:        "test_shop",
			itemID:        "test_item_valid",
			qty:           1,
			setupInv:      func(m *mockInventory) { m.failRem = true },
			expectErr:     true,
			expectErrCode: "NO_ITEM",
			expectEarn:    0,
			expectRemoved: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ecoMock := &mockEconomy{}
			invMock := &mockInventory{}
			if tc.setupEco != nil {
				tc.setupEco(ecoMock)
			}
			if tc.setupInv != nil {
				tc.setupInv(invMock)
			}

			svc := npc.NewService(ecoMock, invMock, zap.NewNop())
			_, err := svc.SellItem(context.Background(), "u1", "g1", tc.shopID, tc.itemID, tc.qty)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Mong đợi lỗi nhưng nhận nil")
				}
				if tc.expectErrCode != "" {
					var appErr *apperrors.AppError
					if errors.As(err, &appErr) {
						if appErr.Code != tc.expectErrCode {
							t.Errorf("Sai mã lỗi: mong đợi %s, nhận %s", tc.expectErrCode, appErr.Code)
						}
					} else {
						t.Errorf("Lỗi không phải là AppError: %v", err)
					}
				}
			} else if err != nil {
				t.Fatalf("Không mong đợi lỗi: %v", err)
			}

			if ecoMock.earn != tc.expectEarn {
				t.Errorf("Earn: mong đợi %d, nhận %d", tc.expectEarn, ecoMock.earn)
			}
			if invMock.removed != tc.expectRemoved {
				t.Errorf("Removed Items: mong đợi %d, nhận %d", tc.expectRemoved, invMock.removed)
			}
		})
	}
}

// --- Concurrent Tests ---

func TestNPCService_ConcurrentBuy(t *testing.T) {
	eco := &mockEconomy{}
	inv := &mockInventory{}
	svc := npc.NewService(eco, inv, zap.NewNop())

	var wg sync.WaitGroup
	const routines = 50

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.BuyItem(context.Background(), "u1", "g1", "test_shop", "test_item_valid", 1)
		}()
	}
	wg.Wait()

	netSpent := eco.spent - eco.earn
	if netSpent != 100*routines {
		t.Errorf("Race condition on economy: net spent %d, want %d", netSpent, 100*routines)
	}
	if inv.added != routines {
		t.Errorf("Race condition on inventory: added %d, want %d", inv.added, routines)
	}
}

func TestNPCService_ConcurrentSell(t *testing.T) {
	eco := &mockEconomy{}
	inv := &mockInventory{}
	svc := npc.NewService(eco, inv, zap.NewNop())

	var wg sync.WaitGroup
	const routines = 50

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.SellItem(context.Background(), "u1", "g1", "test_shop", "test_item_valid", 1)
		}()
	}
	wg.Wait()

	if eco.earn != 20*routines {
		t.Errorf("Race condition on economy: earn %d, want %d", eco.earn, 20*routines)
	}
	if inv.removed != routines {
		t.Errorf("Race condition on inventory: removed %d, want %d", inv.removed, routines)
	}
}

func TestNPCService_ConcurrentBuy_WithRollbacks(t *testing.T) {
	eco := &mockEconomy{}
	inv := &mockInventory{failAdd: true} // Always fail -> should rollback everything
	svc := npc.NewService(eco, inv, zap.NewNop())

	var wg sync.WaitGroup
	const routines = 30

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.BuyItem(context.Background(), "u1", "g1", "test_shop", "test_item_valid", 1)
		}()
	}
	wg.Wait()

	if eco.spent != 100*routines {
		t.Errorf("Eco spent expected %d, got %d", 100*routines, eco.spent)
	}
	if eco.earn != 100*routines {
		t.Errorf("Eco earn (rollback) expected %d, got %d", 100*routines, eco.earn)
	}
	netSpent := eco.spent - eco.earn
	if netSpent != 0 {
		t.Errorf("SAGA Lỗi: Lệnh Rollback đã đánh mất linh thạch trên đường đi! Net Cost: %d", netSpent)
	}
	if inv.added != 0 {
		t.Errorf("Inventory added should be 0, got %d", inv.added)
	}
}

func TestNPCService_ChaosSpam_NoMoneyLost(t *testing.T) {
	// Giả lập 1 kho bạc dùng chung cho tất cả giao dịch để đo lường dòng tiền
	eco := &mockEconomy{}
	inv := &mockInventory{}
	svc := npc.NewService(eco, inv, zap.NewNop())

	var wg sync.WaitGroup
	const routines = 100

	// Mảng ghi chép nội bộ của Test để đối chiếu với Backend Service
	var expectedNetSpent int64
	var expectedNetItems int64
	var mu sync.Mutex

	for i := 0; i < routines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			// Randomize hành động: 50% mua, 50% bán
			isBuy := rand.Intn(2) == 0
			if isBuy {
				_, err := svc.BuyItem(context.Background(), "u1", "g1", "test_shop", "test_item_valid", 1)
				if err == nil {
					mu.Lock()
					expectedNetSpent += 100 // BuyPrice
					expectedNetItems += 1
					mu.Unlock()
				}
			} else {
				_, err := svc.SellItem(context.Background(), "u1", "g1", "test_shop", "test_item_valid", 1)
				if err == nil {
					mu.Lock()
					expectedNetSpent -= 20 // SellPrice (Backend trả tiền cho user -> NetSpent giảm)
					expectedNetItems -= 1
					mu.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()

	actualNetSpent := eco.spent - eco.earn
	if actualNetSpent != expectedNetSpent {
		t.Fatalf("THẤT THOÁT DÒNG TIỀN: Kỳ vọng chênh lệch %d, Thực tế %d", expectedNetSpent, actualNetSpent)
	}
	actualNetItems := inv.added - inv.removed
	if actualNetItems != expectedNetItems {
		t.Fatalf("THẤT THOÁT VẬT PHẨM: Kỳ vọng chênh lệch %d, Thực tế %d", expectedNetItems, actualNetItems)
	}
}
