// File: internal/game/data/pve_rewards.go
package data

import (
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/skill"
)

func init() {
	RegisterItems(map[string]item.ItemDefinition{
		"mat_rare_tinh_thiet":        {ID: "mat_rare_tinh_thiet", Name: "Tinh Thiết Bí Cảnh", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, Description: "Khoáng tinh ngưng tụ trong bí cảnh, dùng cường hóa pháp khí cao phẩm."},
		"eq_weapon_huyet_kiem":       {ID: "eq_weapon_huyet_kiem", Name: "Huyết Kiếm", Type: item.TypeEquipment, Rarity: item.RarityA, MaxEnhanceLevel: 15, Stats: map[string]int{"atk": 80, "crit": 5}, Description: "Kiếm khí nhuốm huyết sát, thường xuất hiện trong bí cảnh âm hàn."},
		"eq_artifact_guong_bat_quai": {ID: "eq_artifact_guong_bat_quai", Name: "Gương Bát Quái", Type: item.TypeEquipment, Rarity: item.RarityS, MaxEnhanceLevel: 5, Stats: map[string]int{"def": 120, "hp": 1000}, Description: "Pháp bảo hộ thân, phản chiếu một tia thiên cơ."},
		"scroll_skill_than_hanh_bo":  {ID: "scroll_skill_than_hanh_bo", Name: "Tàn Quyển Thần Hành Bộ", Type: item.TypeSkillScroll, Rarity: item.RarityA, Usable: true, Description: "Cuộn trục lưu giữ tâm pháp Thần Hành Bộ.", Effects: map[string]int{"unlock_skill": 1}},
		"scroll_skill_huyet_sat":     {ID: "scroll_skill_huyet_sat", Name: "Tàn Quyển Huyết Sát", Type: item.TypeSkillScroll, Rarity: item.RarityA, Usable: true, Description: "Cuộn trục lưu giữ tà pháp Huyết Sát.", Effects: map[string]int{"unlock_skill": 1}},
	})

	skill.Registry["skill_huyet_sat"] = skill.SkillDefinition{
		ID:                "skill_huyet_sat",
		Name:              "Huyết Sát",
		Description:       "Hiến tế máu để sát thương địch nhân.",
		SkillType:         "active",
		AllowedActorTypes: []string{"player"},
		Effects:           []skill.SkillEffect{{Type: "damage", Value: 3.0, ScalingStat: "atk"}},
	}
}
