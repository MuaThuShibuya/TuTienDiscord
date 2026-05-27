// File: internal/game/skill/skill_test.go

package skill

import (
	"testing"
)

func TestSkillDefinitionParse(t *testing.T) {
	skill := Registry["skill_linh_khi_tram"]
	if skill.Cost.Rage != 50 {
		t.Errorf("Linh khí trảm phải tốn 50 Nộ, nhưng nhận %d", skill.Cost.Rage)
	}

	support := Registry["skill_than_hanh_bo"]
	if len(support.Effects) > 0 && support.Effects[0].Type != "turn_advance" {
		t.Errorf("Thần hành bộ phải có effect turn_advance")
	}
}

func TestEquippedSkillSet_SetSlot_Success(t *testing.T) {
	set := NewEquippedSkillSet("player", "u1")
	if err := set.SetSlot(0, "skill_basic_slash"); err != nil {
		t.Errorf("Không mong đợi lỗi khi gài slot 0: %v", err)
	}
	if err := set.SetSlot(2, "skill_linh_khi_tram"); err != nil {
		t.Errorf("Không mong đợi lỗi khi gài slot 2: %v", err)
	}
}

func TestEquippedSkillSet_SetSlot_OutOfRange(t *testing.T) {
	set := NewEquippedSkillSet("player", "u1")
	if err := set.SetSlot(3, "skill_invalid"); err == nil {
		t.Errorf("Mong đợi lỗi khi gài slot ngoài mảng")
	}
}

func TestEquippedSkillSet_SetSlot_LockedSlot(t *testing.T) {
	set := NewEquippedSkillSet("player", "u1")
	set.Slots[1].Unlocked = false
	if err := set.SetSlot(1, "skill_basic_slash"); err == nil {
		t.Errorf("Mong đợi lỗi khi gài vào slot chưa mở khóa")
	}
}
