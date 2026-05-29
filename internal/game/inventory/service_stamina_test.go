package inventory_test

import (
	"context"
	"strings"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

type mockCultSvc struct {
	cultivation.Service // embed interface để khỏi phải mock mọi hàm
	expAdded            int64
}

func (m *mockCultSvc) AddExperience(ctx context.Context, userID, guildID string, amount int64) error {
	m.expAdded += amount
	return nil
}

func TestUseItem_StaminaRemovedBehavior(t *testing.T) {
	// Setup Item definitions
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_stamina_only": {
			ID: "pill_stamina_only", Name: "Hồi Lực Đan", Type: item.TypePill, Usable: true,
			Effects: map[string]int{"stamina": 50},
		},
		"pill_exp_and_stamina": {
			ID: "pill_exp_and_stamina", Name: "Linh Đan Hỗn Hợp", Type: item.TypePill, Usable: true,
			Effects: map[string]int{"stamina": 20, "exp": 100},
		},
	})

	itemRepo := newMockItemRepo()
	cultSvc := &mockCultSvc{}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)
	ctx := context.Background()
	userID := "user_test"
	guildID := "guild1"

	// 1. Test Stamina Only
	_ = svc.AddItem(ctx, userID, guildID, "pill_stamina_only", 1)
	items, _ := itemRepo.GetInstancesByUser(ctx, userID, guildID)
	inst1 := items[0].InstanceID

	msg1, err1 := svc.UseItem(ctx, userID, guildID, inst1)
	if err1 != nil {
		t.Errorf("Dùng item stamina-only không nên lỗi, mà trả về message. Lỗi: %v", err1)
	}
	if !strings.Contains(msg1, "Cơ chế Thể Lực đã được gỡ bỏ") {
		t.Errorf("Thiếu message báo gỡ stamina: %s", msg1)
	}

	// Kiểm tra item không bị trừ
	afterItems, _ := itemRepo.GetInstancesByUser(ctx, userID, guildID)
	if len(afterItems) == 0 || afterItems[0].Quantity != 1 {
		t.Errorf("Item stamina-only đã bị trừ dù hệ thống đã bị gỡ bỏ!")
	}

	// 2. Test Exp + Stamina
	_ = svc.AddItem(ctx, userID, guildID, "pill_exp_and_stamina", 1)
	items2, _ := itemRepo.GetInstancesByUser(ctx, userID, guildID)
	var inst2 string
	for _, it := range items2 {
		if it.DefinitionID == "pill_exp_and_stamina" {
			inst2 = it.InstanceID
		}
	}

	msg2, err2 := svc.UseItem(ctx, userID, guildID, inst2)
	if err2 != nil {
		t.Errorf("Dùng item exp+stamina lỗi: %v", err2)
	}
	if cultSvc.expAdded != 100 {
		t.Errorf("Kỳ vọng expAdded = 100, có: %d", cultSvc.expAdded)
	}
	if strings.Contains(msg2, "thể lực") {
		t.Errorf("Message không nên hiển thị số thể lực hồi khi đã bỏ stamina: %s", msg2)
	}
}
