// File: internal/game/aptitude/registry.go
package aptitude

var Registry = map[string]AptitudeDefinition{
	// --- Mortal / Common ---
	"apt_pham_tu": {
		ID: "apt_pham_tu", Name: "Phàm Tư", Rarity: AptitudeMortal, Weight: 5000,
		Description:          "Tư chất của phàm nhân, con đường tu tiên mịt mờ, gian nan muôn vàn.",
		BaseStats:            BaseStatBonus{MaxHP: 100, ATK: 10, DEF: 5, Speed: 90},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.0, ATKMultiplier: 1.0, DEFMultiplier: 1.0, SpeedMultiplier: 1.0},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 0.8},
	},
	"apt_tap_linh_can": {
		ID: "apt_tap_linh_can", Name: "Tạp Linh Căn", Rarity: AptitudeCommon, Weight: 3000,
		Description:          "Sở hữu linh căn nhưng hỗn tạp, khó đột phá đại đạo.",
		BaseStats:            BaseStatBonus{MaxHP: 150, ATK: 15, DEF: 8, Speed: 100},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.05, ATKMultiplier: 1.05, DEFMultiplier: 1.05, SpeedMultiplier: 1.0},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.0},
	},

	// --- Uncommon / Rare ---
	"apt_ngu_hanh_linh_can": {
		ID: "apt_ngu_hanh_linh_can", Name: "Ngũ Hành Linh Căn", Rarity: AptitudeUncommon, Weight: 1200,
		Description:          "Ngũ hành tương sinh tương khắc, bình ổn cân bằng.",
		BaseStats:            BaseStatBonus{MaxHP: 200, ATK: 20, DEF: 12, Speed: 105},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.1, ATKMultiplier: 1.1, DEFMultiplier: 1.1, SpeedMultiplier: 1.05},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.1},
	},
	"apt_kiem_tam_so_khai": {
		ID: "apt_kiem_tam_so_khai", Name: "Kiếm Tâm Sơ Khải", Rarity: AptitudeRare, Weight: 500,
		Description:          "Trời sinh thân cận với kiếm, một kiếm phá vạn pháp.",
		BaseStats:            BaseStatBonus{MaxHP: 180, ATK: 30, DEF: 10, Speed: 115, CritRate: 0.05},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.05, ATKMultiplier: 1.25, DEFMultiplier: 1.05, SpeedMultiplier: 1.1},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.15},
		Affinities:           []AptitudeAffinity{{Type: "weapon_type", RefID: "sword", Multiplier: 1.2}},
	},
	"apt_luyen_the_can_cot": {
		ID: "apt_luyen_the_can_cot", Name: "Thiên Sinh Thần Lực", Rarity: AptitudeRare, Weight: 500,
		Description:          "Khí huyết cuồn cuộn, căn cốt như rồng.",
		BaseStats:            BaseStatBonus{MaxHP: 300, ATK: 20, DEF: 20, Speed: 95},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.3, ATKMultiplier: 1.1, DEFMultiplier: 1.2, SpeedMultiplier: 1.0},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.15},
	},

	// --- Epic ---
	"apt_thuan_duong_hoa_linh_can": {
		ID: "apt_thuan_duong_hoa_linh_can", Name: "Thuần Dương Hỏa Linh Căn", Rarity: AptitudeEpic, Weight: 200,
		Description:          "Lửa thuần dương thiêu đốt vạn vật, tinh lực vô biên.",
		BaseStats:            BaseStatBonus{MaxHP: 250, ATK: 40, DEF: 15, Speed: 110},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.15, ATKMultiplier: 1.35, DEFMultiplier: 1.1, SpeedMultiplier: 1.1},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.3, BreakthroughChanceBonus: 0.05},
	},
	"apt_huyen_am_thuy_linh_can": {
		ID: "apt_huyen_am_thuy_linh_can", Name: "Huyền Âm Thủy Linh Căn", Rarity: AptitudeEpic, Weight: 200,
		Description:          "Âm nhu tột cùng, linh lực mềm mại nhưng thâm hậu.",
		BaseStats:            BaseStatBonus{MaxHP: 350, ATK: 25, DEF: 30, Speed: 105},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.25, ATKMultiplier: 1.15, DEFMultiplier: 1.35, SpeedMultiplier: 1.1},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.3, MeditationCooldownMultiplier: 0.9},
	},
	"apt_doc_dao_linh_the": {
		ID: "apt_doc_dao_linh_the", Name: "Độc Đạo Linh Thể", Rarity: AptitudeEpic, Weight: 200,
		Description:          "Bách độc bất xâm, giơ tay nhấc chân đều là kịch độc.",
		BaseStats:            BaseStatBonus{MaxHP: 220, ATK: 35, DEF: 20, Speed: 120},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.15, ATKMultiplier: 1.3, DEFMultiplier: 1.2, SpeedMultiplier: 1.2},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.25},
	},

	// --- Legendary ---
	"apt_tien_thien_kiem_cot": {
		ID: "apt_tien_thien_kiem_cot", Name: "Tiên Thiên Kiếm Cốt", Rarity: AptitudeLegendary, Weight: 80,
		Description:          "Sinh ra để dùng kiếm, sát phạt quả đoán, kinh tài tuyệt diễm.",
		BaseStats:            BaseStatBonus{MaxHP: 300, ATK: 60, DEF: 25, Speed: 130, CritRate: 0.1, CritDamage: 1.8},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.2, ATKMultiplier: 1.5, DEFMultiplier: 1.15, SpeedMultiplier: 1.25},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.5, BreakthroughChanceBonus: 0.1},
	},
	"apt_van_thu_linh_tam": {
		ID: "apt_van_thu_linh_tam", Name: "Vạn Thú Linh Tâm", Rarity: AptitudeLegendary, Weight: 80,
		Description:          "Tâm hồn thông linh với vạn yêu, ngự thú vô song.",
		BaseStats:            BaseStatBonus{MaxHP: 400, ATK: 30, DEF: 40, Speed: 110},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.4, ATKMultiplier: 1.2, DEFMultiplier: 1.4, SpeedMultiplier: 1.1},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 1.4},
	},

	// --- Mythic / Heavenly ---
	"apt_hon_don_dao_the": {
		ID: "apt_hon_don_dao_the", Name: "Hỗn Độn Đạo Thể", Rarity: AptitudeMythic, Weight: 15,
		Description:          "Thể chất trong truyền thuyết, bao hàm vạn tượng, sinh ra từ thời hồng hoang.",
		BaseStats:            BaseStatBonus{MaxHP: 500, ATK: 70, DEF: 50, Speed: 140, CritRate: 0.15},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 1.6, ATKMultiplier: 1.6, DEFMultiplier: 1.6, SpeedMultiplier: 1.4},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 2.0, BreakthroughChanceBonus: 0.2, StaminaCostMultiplier: 0.8},
	},
	"apt_nghich_thien_dao_thai": {
		ID: "apt_nghich_thien_dao_thai", Name: "Nghịch Thiên Đạo Thai", Rarity: AptitudeHeavenly, Weight: 5,
		Description:          "Mệnh cách nghịch thiên, thiên địa bất dung. Khởi điểm cực mạnh nhưng độ kiếp sinh tử nan lường.",
		BaseStats:            BaseStatBonus{MaxHP: 800, ATK: 100, DEF: 80, Speed: 160, CritRate: 0.2, CritDamage: 2.0},
		GrowthStats:          GrowthStatMultiplier{MaxHPMultiplier: 2.0, ATKMultiplier: 2.0, DEFMultiplier: 2.0, SpeedMultiplier: 1.5},
		CultivationModifiers: CultivationModifier{ExpGainMultiplier: 3.0, BreakthroughChanceBonus: -0.1}, // Buff bự nhưng penalty nhẹ lúc đột phá
	},
}

// GetRandomAptitude roll tư chất dựa trên tỷ lệ (Weight).
func GetRandomAptitude(rngFunc func(int) int) AptitudeDefinition {
	totalWeight := 0
	for _, def := range Registry {
		totalWeight += int(def.Weight)
	}
	roll := rngFunc(totalWeight)
	for _, def := range Registry {
		roll -= int(def.Weight)
		if roll < 0 {
			return def
		}
	}
	return Registry["apt_pham_tu"]
}
