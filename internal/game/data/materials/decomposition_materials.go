// File: internal/data/materials/decomposition_materials.go
package materials

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"mat_scrap_linh_tran_d":      {ID: "mat_scrap_linh_tran_d", Name: "Linh Trần", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_scrap_thiet_vun_d":      {ID: "mat_scrap_thiet_vun_d", Name: "Thiết Vụn", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_scrap_tan_linh_d":       {ID: "mat_scrap_tan_linh_d", Name: "Tàn Linh Mảnh", Type: item.TypeMaterial, Rarity: item.RarityD, Stackable: true, MaxStack: 9999},
		"mat_scrap_thanh_dong_c":     {ID: "mat_scrap_thanh_dong_c", Name: "Thanh Đồng Mảnh", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_scrap_moc_linh_c":       {ID: "mat_scrap_moc_linh_c", Name: "Mộc Linh Mảnh", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_scrap_hoa_linh_c":       {ID: "mat_scrap_hoa_linh_c", Name: "Hỏa Linh Mảnh", Type: item.TypeMaterial, Rarity: item.RarityC, Stackable: true, MaxStack: 9999},
		"mat_scrap_huyen_thiet_b":    {ID: "mat_scrap_huyen_thiet_b", Name: "Huyền Thiết Mảnh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_scrap_han_bang_b":       {ID: "mat_scrap_han_bang_b", Name: "Hàn Băng Mảnh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_scrap_tu_van_b":         {ID: "mat_scrap_tu_van_b", Name: "Tử Vân Mảnh", Type: item.TypeMaterial, Rarity: item.RarityB, Stackable: true, MaxStack: 9999},
		"mat_scrap_bach_ngoc_a":      {ID: "mat_scrap_bach_ngoc_a", Name: "Bạch Ngọc Mảnh", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_scrap_linh_hoa_a":       {ID: "mat_scrap_linh_hoa_a", Name: "Linh Hỏa Mảnh", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_scrap_huyen_am_a":       {ID: "mat_scrap_huyen_am_a", Name: "Huyền Âm Mảnh", Type: item.TypeMaterial, Rarity: item.RarityA, Stackable: true, MaxStack: 9999},
		"mat_scrap_long_van_s":       {ID: "mat_scrap_long_van_s", Name: "Long Văn Mảnh", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_scrap_phuong_hoa_s":     {ID: "mat_scrap_phuong_hoa_s", Name: "Phượng Hỏa Mảnh", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_scrap_cuu_loi_s":        {ID: "mat_scrap_cuu_loi_s", Name: "Cửu Lôi Mảnh", Type: item.TypeMaterial, Rarity: item.RarityS, Stackable: true, MaxStack: 9999},
		"mat_scrap_tinh_ha_ss":       {ID: "mat_scrap_tinh_ha_ss", Name: "Tinh Hà Mảnh", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999},
		"mat_scrap_van_phap_ss":      {ID: "mat_scrap_van_phap_ss", Name: "Vạn Pháp Mảnh", Type: item.TypeMaterial, Rarity: item.RaritySS, Stackable: true, MaxStack: 9999},
		"mat_scrap_hu_khong_sss":     {ID: "mat_scrap_hu_khong_sss", Name: "Hư Không Mảnh", Type: item.TypeMaterial, Rarity: item.RaritySSS, Stackable: true, MaxStack: 9999},
		"mat_scrap_luan_hoi_sssp":    {ID: "mat_scrap_luan_hoi_sssp", Name: "Luân Hồi Mảnh", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999},
		"mat_scrap_nghich_menh_sssp": {ID: "mat_scrap_nghich_menh_sssp", Name: "Nghịch Mệnh Mảnh", Type: item.TypeMaterial, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 9999},
	})
}
