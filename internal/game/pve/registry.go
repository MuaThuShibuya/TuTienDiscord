// File: internal/game/pve/registry.go
// Chức năng: Kho lưu trữ mẫu (Registry) cho các cấu hình PvE Area và Reward Pool ban đầu.

package pve

var MonsterRegistry = map[string]MonsterDefinition{
	// Du Ngoạn
	"monster_linh_tho": {ID: "monster_linh_tho", Name: "Linh Thỏ", BaseLevel: 1, BaseStats: MonsterStats{MaxHP: 50, ATK: 5, DEF: 2, Speed: 100}},
	"monster_da_lang":  {ID: "monster_da_lang", Name: "Dã Lang", BaseLevel: 2, BaseStats: MonsterStats{MaxHP: 80, ATK: 12, DEF: 5, Speed: 110}},
	"monster_truc_yeu": {ID: "monster_truc_yeu", Name: "Trúc Yêu", BaseLevel: 3, BaseStats: MonsterStats{MaxHP: 120, ATK: 8, DEF: 15, Speed: 80}},
	"monster_son_tac":  {ID: "monster_son_tac", Name: "Sơn Tặc", BaseLevel: 4, BaseStats: MonsterStats{MaxHP: 150, ATK: 20, DEF: 10, Speed: 105}},
	"monster_doc_xa":   {ID: "monster_doc_xa", Name: "Độc Xà", BaseLevel: 5, BaseStats: MonsterStats{MaxHP: 70, ATK: 25, DEF: 3, Speed: 130}},
	// Bí Cảnh
	"monster_huyet_bien_bat": {ID: "monster_huyet_bien_bat", Name: "Huyết Biển Bức", BaseLevel: 5, BaseStats: MonsterStats{MaxHP: 200, ATK: 30, DEF: 10, Speed: 140}},
	"monster_thach_ma":       {ID: "monster_thach_ma", Name: "Thạch Ma", BaseLevel: 6, BaseStats: MonsterStats{MaxHP: 500, ATK: 15, DEF: 40, Speed: 70}},
	"monster_am_linh":        {ID: "monster_am_linh", Name: "Âm Linh", BaseLevel: 6, BaseStats: MonsterStats{MaxHP: 180, ATK: 35, DEF: 5, Speed: 120}},
	// Boss
	"boss_huyet_giap_yeu": {ID: "boss_huyet_giap_yeu", Name: "Huyết Giáp Yêu", Role: MonsterRoleBoss, BaseLevel: 10, BaseStats: MonsterStats{MaxHP: 2000, ATK: 80, DEF: 50, Speed: 90}},
}

type MonsterPoolDefinition struct {
	ID              string
	ActivityType    ActivityType
	MonsterIDs      []string
	EliteMonsterIDs []string
	BossMonsterIDs  []string
}

var MonsterPoolRegistry = map[string]MonsterPoolDefinition{
	"pool_du_ngoan_1": {
		ID:              "pool_du_ngoan_1",
		ActivityType:    ActivityDuNgoan,
		MonsterIDs:      []string{"monster_linh_tho", "monster_da_lang", "monster_truc_yeu", "monster_son_tac", "monster_doc_xa"},
		EliteMonsterIDs: []string{"monster_son_tac"}, // Sơn tặc thỉnh thoảng xuất hiện dạng elite
	},
	"pool_bi_canh_1": {
		ID:              "pool_bi_canh_1",
		ActivityType:    ActivityBiCanh,
		MonsterIDs:      []string{"monster_huyet_bien_bat", "monster_thach_ma", "monster_am_linh"},
		EliteMonsterIDs: []string{"monster_thach_ma", "monster_am_linh"},
		BossMonsterIDs:  []string{"boss_huyet_giap_yeu"},
	},
}

var AreaRegistry = map[string]PvEAreaDefinition{
	"area_du_ngoan_rung_truc": {
		ID:                "area_du_ngoan_rung_truc",
		Name:              "Rừng Trúc Ngoại Thành",
		ActivityType:      ActivityDuNgoan,
		Description:       "Khu rừng trúc yên bình, thích hợp cho Luyện Khí kỳ rèn luyện.",
		MinStage:          1,
		MaxStage:          30,
		MonsterPoolID:     "pool_du_ngoan_1",
		RewardPoolID:      "reward_du_ngoan_basic",
		BonusRewardPoolID: "reward_du_ngoan_bonus",
		EntryCost:         EntryCost{Stamina: 5},
	},
	"area_bi_canh_thach_dong": {
		ID:                  "area_bi_canh_thach_dong",
		Name:                "Huyết Thạch Động",
		ActivityType:        ActivityBiCanh,
		Description:         "Bí cảnh nguy hiểm ngập tràn mùi máu, cơ duyên đi liền với hung hiểm.",
		MinStage:            1,
		MaxStage:            30,
		MonsterPoolID:       "pool_bi_canh_1",
		RewardPoolID:        "reward_bi_canh_rare",
		BonusRewardPoolID:   "reward_bi_canh_bonus",
		EntryCost:           EntryCost{Stamina: 20, OpportunityTicket: 1},
		RequiredCombatPower: 1000,
	},
}

var RewardPoolRegistry = map[string]RewardPoolDefinition{
	"reward_du_ngoan_basic": {
		ID:           "reward_du_ngoan_basic",
		ActivityType: ActivityDuNgoan,
		Entries: []RewardEntry{
			{Type: "exp", MinQuantity: 10, MaxQuantity: 50, Chance: 1.0},
			{Type: "stones", MinQuantity: 5, MaxQuantity: 20, Chance: 0.8},
			{Type: "material", RefID: "mat_enhance_hac_thiet_d", MinQuantity: 1, MaxQuantity: 2, Chance: 0.3},
		},
	},
	"reward_du_ngoan_bonus": {
		ID:           "reward_du_ngoan_bonus",
		ActivityType: ActivityDuNgoan,
		Entries: []RewardEntry{
			{Type: "material", RefID: "mat_enhance_hac_thiet_d", MinQuantity: 1, MaxQuantity: 5, Chance: 1.0, Weight: 80},
			{Type: "equipment", RefID: "eq_weapon_moc_kiem_d", MinQuantity: 1, MaxQuantity: 1, Chance: 1.0, Weight: 15, Rarity: "D"},
			{Type: "skill_scroll", RefID: "skill_than_hanh_bo", MinQuantity: 1, MaxQuantity: 1, Chance: 1.0, Weight: 5, Rarity: "A", IsUnique: true},
		},
	},
	"reward_bi_canh_rare": {
		ID:           "reward_bi_canh_rare",
		ActivityType: ActivityBiCanh,
		Entries: []RewardEntry{
			{Type: "exp", MinQuantity: 50, MaxQuantity: 200, Chance: 1.0},
			{Type: "stones", MinQuantity: 20, MaxQuantity: 100, Chance: 1.0},
			{Type: "material", RefID: "mat_rare_tinh_thiet", MinQuantity: 1, MaxQuantity: 3, Chance: 0.5},
			{Type: "equipment", RefID: "eq_weapon_huyet_kiem", Chance: 0.05, Rarity: "A", IsUnique: true},
			{Type: "artifact", RefID: "art_guong_bat_quai", Chance: 0.01, Rarity: "S", IsUnique: true},
			{Type: "skill_scroll", RefID: "skill_huyet_sat", Chance: 0.02, Rarity: "A", IsUnique: true},
		},
	},
	"reward_bi_canh_bonus": {
		ID:           "reward_bi_canh_bonus",
		ActivityType: ActivityBiCanh,
		Entries: []RewardEntry{
			{Type: "material", RefID: "mat_rare_tinh_thiet", MinQuantity: 3, MaxQuantity: 5, Chance: 1.0, Weight: 50},
			{Type: "equipment", RefID: "eq_weapon_huyet_kiem", MinQuantity: 1, MaxQuantity: 1, Chance: 1.0, Weight: 30, Rarity: "A"},
			{Type: "artifact", RefID: "art_guong_bat_quai", MinQuantity: 1, MaxQuantity: 1, Chance: 1.0, Weight: 10, Rarity: "S", IsUnique: true},
			{Type: "skill_scroll", RefID: "skill_huyet_sat", MinQuantity: 1, MaxQuantity: 1, Chance: 1.0, Weight: 10, Rarity: "A", IsUnique: true},
		},
	},
}
