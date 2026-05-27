// File: internal/game/skill/registry.go
// Chức năng: Lưu trữ tĩnh các định nghĩa kỹ năng cơ bản.

package skill

var Registry = map[string]SkillDefinition{
	"skill_basic_slash": {
		ID:                "skill_basic_slash",
		Name:              "Phách Trảm",
		Description:       "Tấn công cơ bản, gây 100% sát thương.",
		SkillType:         "basic",
		AllowedActorTypes: []string{"player", "monster"},
		Effects: []SkillEffect{
			{Type: "damage", Value: 1.0, ScalingStat: "atk"},
			{Type: "rage_gain", Value: 20},
		},
	},
	"skill_linh_khi_tram": {
		ID:                "skill_linh_khi_tram",
		Name:              "Linh Khí Trảm",
		Description:       "Tụ linh khí xuất kiếm, tiêu hao 50 Nộ.",
		SkillType:         "active",
		Cost:              SkillCost{Rage: 50},
		CooldownTurns:     2,
		AllowedActorTypes: []string{"player"},
		Effects: []SkillEffect{
			{Type: "damage", Value: 2.5, ScalingStat: "atk"},
		},
	},
	"skill_than_hanh_bo": {
		ID:                "skill_than_hanh_bo",
		Name:              "Thần Hành Bộ",
		Description:       "Thân pháp quỷ mị, tăng tốc độ xuất thủ kế tiếp.",
		SkillType:         "support",
		CooldownTurns:     3,
		AllowedActorTypes: []string{"player"},
		Effects: []SkillEffect{
			{Type: "turn_advance", TurnAdvanceValue: 0.3}, // Kéo 30% thanh hành động
		},
	},
}
