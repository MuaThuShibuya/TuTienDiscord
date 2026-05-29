// File: internal/game/cultivation/paths_test.go
// Phiên bản: v0.2.1
// Mục đích: Test trực tiếp các helper tính toán Bonus Đạo lộ mà không cần setup toàn bộ Service.

package cultivation_test

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
)

func TestApplyPathExpBonus(t *testing.T) {
	tests := []struct {
		path     cultivation.CultivationPath
		action   string
		base     int64
		expected int64
	}{
		{cultivation.PathNone, "meditate", 100, 100},
		{cultivation.PathSword, "meditate", 100, 105},   // +5%
		{cultivation.PathSpirit, "meditate", 100, 110},  // +10%
		{cultivation.PathSpirit, "seclusion", 100, 110}, // +10%
		{cultivation.PathPoison, "seclusion", 100, 115}, // +15%
		{cultivation.PathSword, "seclusion", 100, 100},  // Không bonus
	}

	for _, tt := range tests {
		res := cultivation.ApplyPathExpBonus(tt.path, tt.action, tt.base)
		if res != tt.expected {
			t.Errorf("Path %s, action %s, base %d: want %d, got %d", tt.path, tt.action, tt.base, tt.expected, res)
		}
	}
}

func TestApplyPathCombatPowerBonus(t *testing.T) {
	// Kiếm tu +10%, Thể tu +5% khi luyện thể
	if cultivation.ApplyPathCombatPowerBonus(cultivation.PathSword, "body_training", 100) != 110 {
		t.Error("Kiếm tu +10% CP")
	}
	if cultivation.ApplyPathCombatPowerBonus(cultivation.PathBody, "body_training", 100) != 105 {
		t.Error("Thể tu +5% CP")
	}
	if cultivation.ApplyPathCombatPowerBonus(cultivation.PathSpirit, "body_training", 100) != 100 {
		t.Error("Linh tu không có bonus CP")
	}
}

func TestApplyPathBreakthroughBonuses(t *testing.T) {
	rate := cultivation.ApplyPathBreakthroughRateBonus(cultivation.PathPoison, 0.50)
	if rate != 0.55 {
		t.Error("Độc tu phải được +5% rate")
	}

	penalty := cultivation.ApplyPathBreakthroughFailurePenalty(cultivation.PathPoison, 5)
	if penalty != 7 {
		t.Error("Độc tu phải bị +2 mind penalty khi xịt")
	}

	penaltySafe := cultivation.ApplyPathBreakthroughFailurePenalty(cultivation.PathSword, 5)
	if penaltySafe != 5 {
		t.Error("Đạo lộ khác không bị thêm penalty")
	}
}
