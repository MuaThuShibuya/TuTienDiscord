package pve

import (
	"strings"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

func TestBuildHPBar_ClampAndFormat(t *testing.T) {
	// Test 1: Bình thường
	bar1 := BuildHPBar(50, 100)
	if !strings.Contains(bar1, "**50%** (50/100)") {
		t.Errorf("Lỗi hiển thị bar 50%%: %s", bar1)
	}

	// Test 2: Vượt quá MaxHP (Buff/Lỗi data)
	bar2 := BuildHPBar(120, 100)
	if !strings.Contains(bar2, "**100%** (100/100)") {
		t.Errorf("Lỗi clamp max HP: %s", bar2)
	}

	// Test 3: Âm HP (Chết)
	bar3 := BuildHPBar(-10, 100)
	if !strings.Contains(bar3, "**0%** (0/100)") {
		t.Errorf("Lỗi clamp âm HP: %s", bar3)
	}

	// Test 4: Lỗi MaxHP 0
	bar4 := BuildHPBar(10, 0)
	if !strings.Contains(bar4, "**0%** (0/0)") {
		t.Errorf("Lỗi chia cho 0: %s", bar4)
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		in  int64
		out string
	}{
		{0, "0"},
		{999, "999"},
		{1000, "1,000"},
		{1234567, "1,234,567"},
		{-1000, "-1,000"}, // Format cả số âm nếu bị debuff
	}
	for _, tc := range tests {
		if got := FormatNumber(tc.in); got != tc.out {
			t.Errorf("FormatNumber(%d) = %s; mong đợi %s", tc.in, got, tc.out)
		}
	}
}

func TestCombatSessionToViewModel_TrimLogs_UX(t *testing.T) {
	cSession := &combat.CombatSession{
		ID:    "ss_logs",
		State: combat.StateActive,
		Logs: []combat.CombatLogEntry{
			{Turn: 1, Message: "Log 1"},
			{Turn: 2, Message: "Log 2"},
			{Turn: 3, Message: "Log 3"},
			{Turn: 4, Message: "Log 4"},
			{Turn: 5, Message: "Log 5"},
			{Turn: 6, Message: "Log 6"}, // Phải cắt bỏ Log 1, chỉ show Log 2->6 lên UI
		},
	}

	vm := CombatSessionToViewModel(cSession, "Bí Cảnh")

	if len(vm.Logs) != 5 {
		t.Errorf("UI Adapter chỉ nên lấy 5 logs gần nhất để tránh tràn Embed, nhận: %d", len(vm.Logs))
	}
	if !strings.Contains(vm.Logs[0], "Log 2") {
		t.Errorf("Log đầu tiên render ra phải là Log 2")
	}
}

func TestCombatSessionToViewModel_Active(t *testing.T) {
	cSession := &combat.CombatSession{
		ID:    "ss_123",
		State: combat.StateActive,
		Player: combat.CombatActor{
			ID: "u1", Name: "Đạo Hữu Test", CurrentHP: 100, Stats: combat.CombatStats{MaxHP: 100},
		},
		Enemies: []combat.CombatActor{
			{ID: "e1", Name: "Thỏ Yêu", CurrentHP: 50, Stats: combat.CombatStats{MaxHP: 100}},
		},
		CurrentActorID: "u1",
	}

	vm := CombatSessionToViewModel(cSession, "Rừng Trúc")

	if vm.State != combat.StateActive {
		t.Errorf("State map sai")
	}
	if vm.TargetID != "e1" {
		t.Errorf("Chưa tự động target quái đầu tiên")
	}
}
