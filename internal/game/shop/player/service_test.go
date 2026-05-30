// File: internal/game/shop/player/service_test.go
// Chức năng: Kiểm tra bạo lực tính toàn vẹn của Sàn Đấu Giá (Player Auction).

package player_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/shop/player"
	"go.uber.org/zap"
)

// --- Mock Dependencies ---

type mockRepo struct {
	mu       sync.Mutex
	listings map[string]*player.Listing
}

func (m *mockRepo) CreateListing(ctx context.Context, listing *player.Listing) error { return nil }
func (m *mockRepo) GetListing(ctx context.Context, listingID string) (*player.Listing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.listings[listingID]
	if !ok {
		return nil, errors.New("not found")
	}
	return l, nil
}

func (m *mockRepo) GetActiveListings(ctx context.Context, guildID string, limit, offset int) ([]*player.Listing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var active []*player.Listing
	for _, l := range m.listings {
		if l.Status == player.StatusActive {
			if guildID == "" || l.GuildID == guildID {
				active = append(active, l)
			}
		}
	}
	return active, nil
}

func (m *mockRepo) GetUserActiveListings(ctx context.Context, userID, guildID string) ([]*player.Listing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return nil, nil // Mock đơn giản để pass compiler
}

func (m *mockRepo) CancelListing(ctx context.Context, listingID, sellerID string) (*player.Listing, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.listings[listingID]
	if !ok || l.SellerID != sellerID {
		return nil, errors.New("not found or not owner")
	}
	return l, nil
}

// AtomicPurchase giả lập lệnh FindOneAndUpdate của MongoDB
func (m *mockRepo) AtomicPurchase(ctx context.Context, listingID string, buyerID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	l, ok := m.listings[listingID]
	if !ok {
		return errors.New("not found")
	}
	// Bắt buộc status phải đang "active" mới cho mua
	if l.Status != player.StatusActive {
		return errors.New("race condition blocked: already sold")
	}
	l.Status = player.StatusSold
	l.BuyerID = buyerID
	return nil
}

type mockEconomy struct {
	economy.Service
	mu           sync.Mutex
	balances     map[string]int64
	insufficient bool
}

func (m *mockEconomy) SpendSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.insufficient {
		return nil, apperrors.ErrInsufficientFunds
	}
	m.balances[userID] -= amount
	return &economy.Wallet{SpiritStones: m.balances[userID]}, nil
}

func (m *mockEconomy) EarnSpiritStones(ctx context.Context, userID, guildID string, amount int64, reason string) (*economy.Wallet, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.balances[userID] += amount
	return &economy.Wallet{SpiritStones: m.balances[userID]}, nil
}

type mockInventory struct {
	inventory.Service
}

// --- Tests ---

func TestAuctionService_Table(t *testing.T) {
	activeListing := &player.Listing{SellerID: "seller_1", TotalPrice: 1000, Status: player.StatusActive}
	soldListing := &player.Listing{SellerID: "seller_1", TotalPrice: 1000, Status: player.StatusSold}

	cases := []struct {
		name          string
		listingID     string
		buyerID       string
		listing       *player.Listing
		insufficient  bool
		expectErr     bool
		expectErrCode string
	}{
		{
			name:      "1. Mua hợp lệ",
			listingID: "list_1",
			buyerID:   "buyer_1",
			listing:   activeListing,
			expectErr: false,
		},
		{
			name:          "2. Không thể tự mua đồ của mình",
			listingID:     "list_2",
			buyerID:       "seller_1", // Trùng ID người bán
			listing:       activeListing,
			expectErr:     true,
			expectErrCode: "SELF_BUY",
		},
		{
			name:          "3. Đồ đã bị bán hoặc hết hạn",
			listingID:     "list_3",
			buyerID:       "buyer_1",
			listing:       soldListing,
			expectErr:     true,
			expectErrCode: "LISTING_NOT_ACTIVE",
		},
		{
			name:          "4. Người mua không đủ tiền",
			listingID:     "list_4",
			buyerID:       "poor_buyer",
			listing:       activeListing,
			insufficient:  true,
			expectErr:     true,
			expectErrCode: "NO_MONEY",
		},
		{
			name:          "5. Không tìm thấy phiếu đấu giá",
			listingID:     "list_ma",
			buyerID:       "buyer_1",
			listing:       nil, // Repo sẽ trả về not found
			expectErr:     true,
			expectErrCode: "LISTING_404",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mockRepo{listings: make(map[string]*player.Listing)}
			if tc.listing != nil {
				// Clone struct để tránh ô nhiễm vùng nhớ giữa các test
				clone := *tc.listing
				repo.listings[tc.listingID] = &clone
			}

			eco := &mockEconomy{balances: make(map[string]int64), insufficient: tc.insufficient}
			eco.balances[tc.buyerID] = 5000 // Khởi tạo ví

			svc := player.NewService(repo, eco, &mockInventory{}, zap.NewNop())
			err := svc.PurchaseListing(context.Background(), tc.buyerID, "g1", tc.listingID)

			if tc.expectErr {
				if err == nil {
					t.Fatalf("Mong đợi lỗi nhưng nhận nil")
				}
				var appErr *apperrors.AppError
				if errors.As(err, &appErr) && appErr.Code != tc.expectErrCode {
					t.Errorf("Sai mã lỗi: mong %s, nhận %s", tc.expectErrCode, appErr.Code)
				}
			} else if err != nil {
				t.Fatalf("Không mong đợi lỗi: %v", err)
			}
		})
	}
}

func TestAuctionService_ConcurrentDoubleSpend(t *testing.T) {
	// BÀI TEST QUAN TRỌNG NHẤT: Bắn 50 request mua CÙNG 1 MÓN ĐỒ cùng lúc.
	// CHỈ ĐƯỢC PHÉP 1 NGƯỜI MUA ĐƯỢC. Tiền của 49 người còn lại phải được hoàn trả 100%.
	repo := &mockRepo{listings: map[string]*player.Listing{
		"hot_item": {SellerID: "seller_1", TotalPrice: 1000, Status: player.StatusActive},
	}}
	eco := &mockEconomy{balances: make(map[string]int64)}
	svc := player.NewService(repo, eco, &mockInventory{}, zap.NewNop())

	var wg sync.WaitGroup
	var successCount int32

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := svc.PurchaseListing(context.Background(), "buyer_bot", "g1", "hot_item")
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("LỖ HỔNG NGHIÊM TRỌNG: %d người mua được cùng 1 món đồ!", successCount)
	}

	// Tiền của Seller phải nhận đúng 1000
	if eco.balances["seller_1"] != 1000 {
		t.Errorf("Người bán nhận sai tiền: %d", eco.balances["seller_1"])
	}

	// Người mua dù bắn 50 request nhưng chỉ mua được 1 -> Chỉ bị trừ 1000.
	// (Những luồng xịt đã trừ 1000 tạm thời nhưng được SAGA nhả về ngay lập tức).
	if eco.balances["buyer_bot"] != -1000 {
		t.Errorf("Hoàn tiền thất bại (Lỗi SAGA Rollback)! Số dư người mua: %d", eco.balances["buyer_bot"])
	}
}
