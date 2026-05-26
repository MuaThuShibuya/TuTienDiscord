// File: internal/data/materials/enhancement_materials.go
package materials

import (
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"mat_enhance_hac_thiet_d":      {ID: "mat_enhance_hac_thiet_d", Name: "Hắc Thiết Tinh", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 10}},
		"mat_enhance_thanh_dong_d":     {ID: "mat_enhance_thanh_dong_d", Name: "Thanh Đồng Tinh", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 15}},
		"mat_enhance_linh_moc_c":       {ID: "mat_enhance_linh_moc_c", Name: "Linh Mộc Tinh", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 25}},
		"mat_enhance_xich_hoa_c":       {ID: "mat_enhance_xich_hoa_c", Name: "Xích Hỏa Tinh", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 35}},
		"mat_enhance_han_bang_b":       {ID: "mat_enhance_han_bang_b", Name: "Hàn Băng Tinh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 60}},
		"mat_enhance_tu_van_b":         {ID: "mat_enhance_tu_van_b", Name: "Tử Vân Tinh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 80}},
		"mat_enhance_huyen_thiet_b":    {ID: "mat_enhance_huyen_thiet_b", Name: "Huyền Thiết Tinh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 100}},
		"mat_enhance_bach_ngoc_a":      {ID: "mat_enhance_bach_ngoc_a", Name: "Bạch Ngọc Tủy", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 160}},
		"mat_enhance_linh_hoa_a":       {ID: "mat_enhance_linh_hoa_a", Name: "Linh Hỏa Tủy", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 200}},
		"mat_enhance_huyen_am_a":       {ID: "mat_enhance_huyen_am_a", Name: "Huyền Âm Tủy", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 240}},
		"mat_enhance_cuu_duong_s":      {ID: "mat_enhance_cuu_duong_s", Name: "Cửu Dương Tinh", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 400}},
		"mat_enhance_long_huyet_s":     {ID: "mat_enhance_long_huyet_s", Name: "Long Huyết Tinh", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 500}},
		"mat_enhance_phuong_hoa_s":     {ID: "mat_enhance_phuong_hoa_s", Name: "Phượng Hỏa Tủy", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 600}},
		"mat_enhance_tinh_ha_ss":       {ID: "mat_enhance_tinh_ha_ss", Name: "Tinh Hà Thạch", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 900}},
		"mat_enhance_van_linh_ss":      {ID: "mat_enhance_van_linh_ss", Name: "Vạn Linh Tinh", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 1100}},
		"mat_enhance_thien_hoa_ss":     {ID: "mat_enhance_thien_hoa_ss", Name: "Thiên Hỏa Tinh", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 1300}},
		"mat_enhance_hu_khong_sss":     {ID: "mat_enhance_hu_khong_sss", Name: "Hư Không Tinh", Type: item.TypeMaterial, Rarity: item.RaritySSS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 2000}},
		"mat_enhance_thien_kiep_sss":   {ID: "mat_enhance_thien_kiep_sss", Name: "Thiên Kiếp Tủy", Type: item.TypeMaterial, Rarity: item.RaritySSS, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 2600}},
		"mat_enhance_luan_hoi_sssp":    {ID: "mat_enhance_luan_hoi_sssp", Name: "Luân Hồi Tinh", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 5000}},
		"mat_enhance_nghich_menh_sssp": {ID: "mat_enhance_nghich_menh_sssp", Name: "Nghịch Mệnh Đạo Tinh", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999, Stats: map[string]int{"enhance_exp": 8000}},
	})
}
