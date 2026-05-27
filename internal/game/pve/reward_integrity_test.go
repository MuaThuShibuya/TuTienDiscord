// File: internal/game/pve/reward_integrity_test.go
// Chức năng: Đảm bảo không có rác cấu hình. Mọi ID phần thưởng trong Pool đều phải tồn tại thật sự.

package pve

import (
	"testing"

	_ "github.com/whiskey/tu-tien-bot/internal/game/data/loader"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/skill"
)

func TestRewardPoolRegistry_AllRefIDsResolvable(t *testing.T) {
	for poolID, pool := range RewardPoolRegistry {
		for _, entry := range pool.Entries {
			switch entry.Type {
			case "exp", "stones":
				continue
			case "material", "item", "skill_scroll", "equipment", "artifact":
				_, ok := item.GetDefinition(entry.RefID)
				if !ok {
					t.Errorf("CRITICAL: Reward Pool [%s] chứa RefID [%s] (Type: %s) không tồn tại trong Item Registry!", poolID, entry.RefID, entry.Type)
				}
			case "skill":
				_, ok := skill.Registry[entry.RefID]
				if !ok {
					t.Errorf("CRITICAL: Reward Pool [%s] chứa Kỹ Năng [%s] không tồn tại trong Skill Registry!", poolID, entry.RefID)
				}
			default:
				t.Errorf("Unknown reward type: %s", entry.Type)
			}
		}
	}
}

func TestRewardPoolRegistry_NoInvalidQuantity(t *testing.T) {
	for poolID, pool := range RewardPoolRegistry {
		for _, entry := range pool.Entries {
			if entry.Type == "exp" || entry.Type == "stones" {
				continue
			}
			if entry.MinQuantity <= 0 {
				t.Errorf("Pool [%s]: Item %s có MinQuantity = %d (Yêu cầu phải > 0)", poolID, entry.RefID, entry.MinQuantity)
			}
			if entry.MaxQuantity < entry.MinQuantity {
				t.Errorf("Pool [%s]: MaxQuantity (%d) nhỏ hơn MinQuantity (%d)", poolID, entry.MaxQuantity, entry.MinQuantity)
			}
		}
	}
}

func TestEquipmentRewardDefinitionsUsable(t *testing.T) {
	for poolID, pool := range RewardPoolRegistry {
		for _, entry := range pool.Entries {
			if entry.Type == "equipment" || entry.Type == "artifact" {
				def, ok := item.GetDefinition(entry.RefID)
				if !ok {
					continue // Lỗi này được TestRewardPoolRegistry_AllRefIDsResolvable báo rồi
				}

				slot := equipment.GetSlotForDefinition(entry.RefID)
				if slot == "" {
					t.Errorf("Pool [%s]: Trang bị [%s] không có Slot hợp lệ (Phải bắt đầu bằng eq_weapon_, eq_armor_, eq_artifact_,...)", poolID, entry.RefID)
				}

				if len(def.Stats) == 0 {
					t.Errorf("Pool [%s]: Trang bị [%s] đang có Stats = 0, vô dụng trong combat", poolID, entry.RefID)
				}
			}
		}
	}
}

func TestSkillScrollReferencesValidSkill(t *testing.T) {
	// Hardcode map kiểm tra cuộn kỹ năng trỏ đúng kỹ năng (vì chưa có struct SkillScroll riêng)
	scrollMap := map[string]string{
		"scroll_skill_than_hanh_bo": "skill_than_hanh_bo",
		"scroll_skill_huyet_sat":    "skill_huyet_sat",
	}

	for scrollID, targetSkillID := range scrollMap {
		if _, ok := item.GetDefinition(scrollID); ok {
			if _, skillOk := skill.Registry[targetSkillID]; !skillOk {
				t.Errorf("Cuộn kỹ năng [%s] tồn tại nhưng Kỹ Năng Đích [%s] không tồn tại trong Skill Registry!", scrollID, targetSkillID)
			}
		}
	}
}

func TestMonsterPoolsExist(t *testing.T) {
	for areaID, area := range AreaRegistry {
		pool, ok := MonsterPoolRegistry[area.MonsterPoolID]
		if !ok {
			t.Errorf("Khu vực [%s] đang trỏ tới Monster Pool [%s] không tồn tại", areaID, area.MonsterPoolID)
		}
		if len(pool.MonsterIDs) == 0 {
			t.Errorf("Monster Pool [%s] không chứa bất kỳ quái vật nào", area.MonsterPoolID)
		}
	}
}
