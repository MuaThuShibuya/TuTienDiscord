// File: internal/data/materials/alchemy_materials.go
package materials

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"mat_herb_linh_thao_d":      {ID: "mat_herb_linh_thao_d", Name: "Linh Thảo", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_herb_thanh_tam_d":      {ID: "mat_herb_thanh_tam_d", Name: "Thanh Tâm Thảo", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_herb_hoi_luc_d":        {ID: "mat_herb_hoi_luc_d", Name: "Hồi Lực Hoa", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_herb_tu_khi_c":         {ID: "mat_herb_tu_khi_c", Name: "Tụ Khí Quả", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_herb_moc_linh_c":       {ID: "mat_herb_moc_linh_c", Name: "Mộc Linh Chi", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_herb_hoa_tam_c":        {ID: "mat_herb_hoa_tam_c", Name: "Hỏa Tâm Thảo", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_herb_han_nguyet_b":     {ID: "mat_herb_han_nguyet_b", Name: "Hàn Nguyệt Hoa", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_herb_tu_van_b":         {ID: "mat_herb_tu_van_b", Name: "Tử Vân Chi", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_herb_huyen_linh_b":     {ID: "mat_herb_huyen_linh_b", Name: "Huyền Linh Quả", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_herb_bach_ngoc_a":      {ID: "mat_herb_bach_ngoc_a", Name: "Bạch Ngọc Liên", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_herb_xich_duong_a":     {ID: "mat_herb_xich_duong_a", Name: "Xích Dương Quả", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_herb_huyen_am_a":       {ID: "mat_herb_huyen_am_a", Name: "Huyền Âm Chi", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_herb_cuu_duong_s":      {ID: "mat_herb_cuu_duong_s", Name: "Cửu Dương Hoa", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_herb_long_huyet_s":     {ID: "mat_herb_long_huyet_s", Name: "Long Huyết Thảo", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_herb_phuong_tuc_s":     {ID: "mat_herb_phuong_tuc_s", Name: "Phượng Tức Liên", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_herb_tinh_ha_ss":       {ID: "mat_herb_tinh_ha_ss", Name: "Tinh Hà Quả", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999},
		"mat_herb_van_nien_ss":      {ID: "mat_herb_van_nien_ss", Name: "Vạn Niên Linh Chi", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999},
		"mat_herb_hu_khong_sss":     {ID: "mat_herb_hu_khong_sss", Name: "Hư Không Đạo Hoa", Type: item.TypeMaterial, Rarity: item.RaritySSS, Stackable: true, MaxStack: 9999},
		"mat_herb_luan_hoi_sssp":    {ID: "mat_herb_luan_hoi_sssp", Name: "Luân Hồi Tiên Quả", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999},
		"mat_herb_nghich_menh_sssp": {ID: "mat_herb_nghich_menh_sssp", Name: "Nghịch Mệnh Đạo Liên", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999},
	})
}
