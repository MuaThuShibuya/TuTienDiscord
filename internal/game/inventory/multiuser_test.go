package inventory_test

import (
	"context"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_test": {ID: "eq_test", Name: "Test Equip", Stackable: false},
	})
}

func TestInventory_ListInventory_ReturnsFreshData(t *testing.T) {
	itemRepo := newMockItemRepo()
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, nil)
	ctx := context.Background()

	// Giả lập nhận thưởng sau PvE
	_ = svc.AddItem(ctx, "user_fresh", "", "eq_test", 1)

	// Lấy inventory ngay sau đó
	_, items, err := svc.GetInventory(ctx, "user_fresh", "")
	if err != nil || len(items) != 1 {
		t.Errorf("GetInventory phải đọc được ngay item vừa ghi vào DB. Lỗi: %v, Số lượng: %d", err, len(items))
	}
}

func TestInventory_UserGuildScoped(t *testing.T) {
	itemRepo := newMockItemRepo()
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, nil)
	ctx := context.Background()

	// User A ở Guild 1 nhận đồ
	_ = svc.AddItem(ctx, "userA", "guild1", "eq_test", 1)

	// User B ở Guild 1 nhận đồ
	_ = svc.AddItem(ctx, "userB", "guild1", "eq_test", 1)

	// User A ở Guild 2 nhận đồ
	_ = svc.AddItem(ctx, "userA", "guild2", "eq_test", 1)

	// Kiểm tra sự độc lập
	_, itemsA_G1, _ := svc.GetInventory(ctx, "userA", "guild1")
	if len(itemsA_G1) != 1 {
		t.Errorf("UserA-Guild1 phải có 1 món, nhận %d", len(itemsA_G1))
	}

	_, itemsB_G1, _ := svc.GetInventory(ctx, "userB", "guild1")
	if len(itemsB_G1) != 1 {
		t.Errorf("UserB-Guild1 phải có 1 món, đồ không được gộp chung, nhận %d", len(itemsB_G1))
	}
	if itemsA_G1[0].InstanceID == itemsB_G1[0].InstanceID {
		t.Errorf("InstanceID của 2 người phải khác nhau hoàn toàn")
	}

	_, itemsA_G2, _ := svc.GetInventory(ctx, "userA", "guild2")
	if itemsA_G1[0].InstanceID == itemsA_G2[0].InstanceID {
		t.Errorf("Đồ của cùng 1 user nhưng ở 2 Guild khác nhau phải bị tách biệt biệt lập")
	}
}
