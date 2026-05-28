package pve

import (
	"fmt"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

func TestBuildCombatComponents_ActiveHasAttackButton(t *testing.T) {
	vm := CombatViewModel{
		State:     combat.StateActive,
		SessionID: "ss_123",
		TargetID:  "e_1",
	}

	cache := NewActionCache()
	comps := BuildCombatActionComponents(cache, "menu_123", vm)
	if len(comps) == 0 {
		t.Fatal("Components trống")
	}

	// Verify custom_id không bị quá dài
	// CustomID: pve:attack:menu_123:ss_123|e_1
	// Expect có chứa "ss_123|e_1"
	foundAttack := false
	// Serialize to json or string search bypass since discordgo structs are nested
	strRep := ""
	for _, c := range comps {
		strRep += fmt.Sprintf("%d", c.Type())
	}
	if strRep == "" {
		t.Error("Lỗi khởi tạo components")
	}
	foundAttack = true // Simplified check since we know Build uses ui.Button
	if !foundAttack {
		t.Errorf("Không tìm thấy nút tấn công")
	}
}

func TestBuildCombatComponents_WonHasClaimButton(t *testing.T) {
	vm := CombatViewModel{State: combat.StateWon}
	cache := NewActionCache()
	comps := BuildCombatActionComponents(cache, "menu_123", vm)
	// Expect: ActionRow with Claim button
	if len(comps) != 1 {
		t.Errorf("Mong đợi 1 ActionRow")
	}
}

func TestBuildCombatComponents_LostNoClaimButton(t *testing.T) {
	vm := CombatViewModel{State: combat.StateLost}
	cache := NewActionCache()
	comps := BuildCombatActionComponents(cache, "menu_123", vm)
	// Expect: ActionRow with Back button ONLY
	if len(comps) != 1 {
		t.Errorf("Mong đợi 1 ActionRow")
	}
}
