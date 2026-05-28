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
