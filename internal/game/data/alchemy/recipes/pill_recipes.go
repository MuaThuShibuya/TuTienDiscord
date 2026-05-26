// File: internal/data/alchemy/recipes/pill_recipes.go
package recipes

import (
	"github.com/whiskey/tu-tien-bot/internal/game/alchemy"
	"github.com/whiskey/tu-tien-bot/internal/game/data"
)

func init() {
	data.RegisterRecipes(map[string]alchemy.Recipe{
		// --- ĐAN TĂNG TU VI (EXP) ---
		"recipe_tu_khi_d": {
			ID: "recipe_tu_khi_d", Name: "Tụ Khí Đan", OutputItem: "pill_exp_tu_khi_d", OutputQuantity: 1,
			SuccessRate: 0.90, LevelRequired: 1, ExpReward: 5,
			RequiredItems: map[string]int64{"mat_herb_linh_thao_d": 3, "mat_herb_thanh_tam_d": 1},
		},
		"recipe_ngung_khi_c": {
			ID: "recipe_ngung_khi_c", Name: "Ngưng Khí Đan", OutputItem: "pill_exp_ngung_khi_c", OutputQuantity: 1,
			SuccessRate: 0.85, LevelRequired: 2, ExpReward: 12,
			RequiredItems: map[string]int64{"mat_herb_tu_khi_c": 3, "mat_herb_moc_linh_c": 1},
		},
		"recipe_truc_co_b": {
			ID: "recipe_truc_co_b", Name: "Trúc Cơ Đan", OutputItem: "pill_exp_truc_co_b", OutputQuantity: 1,
			SuccessRate: 0.75, LevelRequired: 4, ExpReward: 30,
			RequiredItems: map[string]int64{"mat_herb_huyen_linh_b": 3, "mat_herb_tu_van_b": 2},
		},
		"recipe_linh_nguyen_a": {
			ID: "recipe_linh_nguyen_a", Name: "Linh Nguyên Đan", OutputItem: "pill_exp_linh_nguyen_a", OutputQuantity: 1,
			SuccessRate: 0.65, LevelRequired: 7, ExpReward: 80,
			RequiredItems: map[string]int64{"mat_herb_bach_ngoc_a": 3, "mat_herb_xich_duong_a": 2},
		},
		"recipe_huyen_nguyen_s": {
			ID: "recipe_huyen_nguyen_s", Name: "Huyền Nguyên Đan", OutputItem: "pill_exp_huyen_nguyen_s", OutputQuantity: 1,
			SuccessRate: 0.50, LevelRequired: 10, ExpReward: 200,
			RequiredItems: map[string]int64{"mat_herb_cuu_duong_s": 4, "mat_herb_long_huyet_s": 2},
		},
		"recipe_thien_nguyen_ss": {
			ID: "recipe_thien_nguyen_ss", Name: "Thiên Nguyên Đan", OutputItem: "pill_exp_thien_nguyen_ss", OutputQuantity: 1,
			SuccessRate: 0.35, LevelRequired: 14, ExpReward: 500,
			RequiredItems: map[string]int64{"mat_herb_tinh_ha_ss": 4, "mat_herb_van_nien_ss": 3},
		},
		"recipe_tien_nguyen_sss": {
			ID: "recipe_tien_nguyen_sss", Name: "Tiên Nguyên Đan", OutputItem: "pill_exp_tien_nguyen_sss", OutputQuantity: 1,
			SuccessRate: 0.20, LevelRequired: 18, ExpReward: 1500,
			RequiredItems: map[string]int64{"mat_herb_hu_khong_sss": 5},
		},
		"recipe_nghich_menh_exp_sssp": {
			ID: "recipe_nghich_menh_exp_sssp", Name: "Nghịch Mệnh Tiên Đan", OutputItem: "pill_exp_nghich_menh_sssp", OutputQuantity: 1,
			SuccessRate: 0.05, LevelRequired: 20, ExpReward: 5000,
			RequiredItems: map[string]int64{"mat_herb_luan_hoi_sssp": 3, "mat_herb_nghich_menh_sssp": 3},
		},

		// --- ĐAN HỒI THỂ LỰC (STAMINA) ---
		"recipe_hoi_luc_d": {
			ID: "recipe_hoi_luc_d", Name: "Hồi Lực Đan", OutputItem: "pill_stm_hoi_luc_d", OutputQuantity: 1,
			SuccessRate: 0.95, LevelRequired: 1, ExpReward: 4,
			RequiredItems: map[string]int64{"mat_herb_hoi_luc_d": 3, "mat_herb_linh_thao_d": 1},
		},
		"recipe_thanh_moc_c": {
			ID: "recipe_thanh_moc_c", Name: "Thanh Mộc Đan", OutputItem: "pill_stm_thanh_moc_c", OutputQuantity: 1,
			SuccessRate: 0.88, LevelRequired: 2, ExpReward: 10,
			RequiredItems: map[string]int64{"mat_herb_moc_linh_c": 3, "mat_herb_hoi_luc_d": 2},
		},
		"recipe_sinh_cot_b": {
			ID: "recipe_sinh_cot_b", Name: "Sinh Cốt Đan", OutputItem: "pill_stm_sinh_cot_b", OutputQuantity: 1,
			SuccessRate: 0.80, LevelRequired: 4, ExpReward: 25,
			RequiredItems: map[string]int64{"mat_herb_han_nguyet_b": 3, "mat_herb_tu_van_b": 2},
		},
		"recipe_linh_tuc_a": {
			ID: "recipe_linh_tuc_a", Name: "Linh Tức Đan", OutputItem: "pill_stm_linh_tuc_a", OutputQuantity: 1,
			SuccessRate: 0.70, LevelRequired: 6, ExpReward: 60,
			RequiredItems: map[string]int64{"mat_herb_huyen_am_a": 3, "mat_herb_bach_ngoc_a": 2},
		},
		"recipe_cuu_chuyen_s": {
			ID: "recipe_cuu_chuyen_s", Name: "Cửu Chuyển Hồi Lực Đan", OutputItem: "pill_stm_cuu_chuyen_s", OutputQuantity: 1,
			SuccessRate: 0.55, LevelRequired: 10, ExpReward: 150,
			RequiredItems: map[string]int64{"mat_herb_phuong_tuc_s": 4, "mat_herb_long_huyet_s": 2},
		},

		// --- ĐAN HỖ TRỢ ĐỘT PHÁ (BREAKTHROUGH) ---
		"recipe_pha_canh_d": {
			ID: "recipe_pha_canh_d", Name: "Tiểu Phá Cảnh Đan", OutputItem: "pill_brk_pha_canh_d", OutputQuantity: 1,
			SuccessRate: 0.85, LevelRequired: 2, ExpReward: 8,
			RequiredItems: map[string]int64{"mat_herb_thanh_tam_d": 3, "mat_herb_tu_khi_c": 1},
		},
		"recipe_pha_canh_c": {
			ID: "recipe_pha_canh_c", Name: "Phá Cảnh Đan", OutputItem: "pill_brk_pha_canh_c", OutputQuantity: 1,
			SuccessRate: 0.75, LevelRequired: 3, ExpReward: 20,
			RequiredItems: map[string]int64{"mat_herb_hoa_tam_c": 3, "mat_herb_tu_khi_c": 2},
		},
		"recipe_ho_mach_b": {
			ID: "recipe_ho_mach_b", Name: "Hộ Mạch Đan", OutputItem: "pill_brk_ho_mach_b", OutputQuantity: 1,
			SuccessRate: 0.60, LevelRequired: 5, ExpReward: 45,
			RequiredItems: map[string]int64{"mat_herb_tu_van_b": 3, "mat_herb_han_nguyet_b": 2},
		},
		"recipe_huyen_mon_a": {
			ID: "recipe_huyen_mon_a", Name: "Huyền Môn Phá Cảnh Đan", OutputItem: "pill_brk_huyen_mon_a", OutputQuantity: 1,
			SuccessRate: 0.45, LevelRequired: 8, ExpReward: 100,
			RequiredItems: map[string]int64{"mat_herb_xich_duong_a": 4, "mat_herb_huyen_am_a": 2},
		},
		"recipe_thien_kiep_s": {
			ID: "recipe_thien_kiep_s", Name: "Thiên Kiếp Hộ Mạch Đan", OutputItem: "pill_brk_thien_kiep_s", OutputQuantity: 1,
			SuccessRate: 0.30, LevelRequired: 12, ExpReward: 300,
			RequiredItems: map[string]int64{"mat_herb_cuu_duong_s": 4, "mat_herb_phuong_tuc_s": 3},
		},
		"recipe_ngo_dao_ss": {
			ID: "recipe_ngo_dao_ss", Name: "Ngộ Đạo Đan", OutputItem: "pill_brk_ngo_dao_ss", OutputQuantity: 1,
			SuccessRate: 0.15, LevelRequired: 16, ExpReward: 800,
			RequiredItems: map[string]int64{"mat_herb_van_nien_ss": 5, "mat_herb_tinh_ha_ss": 3},
		},
		"recipe_luan_hoi_brk_sssp": {
			ID: "recipe_luan_hoi_brk_sssp", Name: "Luân Hồi Ngộ Đạo Đan", OutputItem: "pill_brk_luan_hoi_sssp", OutputQuantity: 1,
			SuccessRate: 0.03, LevelRequired: 20, ExpReward: 10000,
			RequiredItems: map[string]int64{"mat_herb_luan_hoi_sssp": 5, "mat_herb_nghich_menh_sssp": 2},
		},
	})
}
