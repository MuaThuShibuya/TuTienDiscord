// File: internal/data/equipment/boots/boots.go
package boots

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_boots_vai_bo_d":         {ID: "eq_boots_vai_bo_d", Name: "Vải Bố Ngoa", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Giày vải rẻ tiền.", Stats: map[string]int{"speed": 3}, SellPrice: 15},
		"eq_boots_thanh_phong_d":    {ID: "eq_boots_thanh_phong_d", Name: "Thanh Phong Ngoa", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Nhẹ như gió.", Stats: map[string]int{"speed": 5}, SellPrice: 30},
		"eq_boots_hac_thiet_c":      {ID: "eq_boots_hac_thiet_c", Name: "Hắc Thiết Chiến Ngoa", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Giày sắt bọc chân.", Stats: map[string]int{"speed": 8, "defense": 5}, SellPrice: 100},
		"eq_boots_linh_moc_c":       {ID: "eq_boots_linh_moc_c", Name: "Linh Mộc Khinh Ngoa", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Khảm mộc thạch.", Stats: map[string]int{"speed": 10, "stamina_max": 3}, SellPrice: 150},
		"eq_boots_tat_phong_b":      {ID: "eq_boots_tat_phong_b", Name: "Tật Phong Ngoa", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Đi mây về gió.", Stats: map[string]int{"speed": 18}, SellPrice: 400},
		"eq_boots_bach_van_b":       {ID: "eq_boots_bach_van_b", Name: "Bạch Vân Linh Ngoa", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Bồng bềnh như mây.", Stats: map[string]int{"speed": 20, "dodge_rate": 1}, SellPrice: 450},
		"eq_boots_huyen_anh_b":      {ID: "eq_boots_huyen_anh_b", Name: "Huyền Ảnh Ngoa", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Lưu lại tàn ảnh.", Stats: map[string]int{"speed": 25, "stamina_max": 8}, SellPrice: 600},
		"eq_boots_tu_dien_a":        {ID: "eq_boots_tu_dien_a", Name: "Tử Điện Bộ Ngoa", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Sải bước ra sấm sét.", Stats: map[string]int{"speed": 40, "dodge_rate": 2}, SellPrice: 1500},
		"eq_boots_han_nguyet_a":     {ID: "eq_boots_han_nguyet_a", Name: "Hàn Nguyệt Khinh Ngoa", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Lạnh lẽo tĩnh mịch.", Stats: map[string]int{"speed": 45, "stamina_regen_bonus": 1}, SellPrice: 2000},
		"eq_boots_xich_hoa_a":       {ID: "eq_boots_xich_hoa_a", Name: "Xích Hỏa Chiến Ngoa", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Cháy rực mỗi bước đi.", Stats: map[string]int{"speed": 50, "hp": 150}, SellPrice: 2200},
		"eq_boots_vo_anh_s":         {ID: "eq_boots_vo_anh_s", Name: "Vô Ảnh Ngoa", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Nhanh không thấy bóng.", Stats: map[string]int{"speed": 80, "dodge_rate": 4}, SellPrice: 6000},
		"eq_boots_cuu_bo_s":         {ID: "eq_boots_cuu_bo_s", Name: "Cửu Bộ Đạp Vân Ngoa", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Chín bước lên trời.", Stats: map[string]int{"speed": 90, "stamina_max": 20}, SellPrice: 7500},
		"eq_boots_loi_anh_s":        {ID: "eq_boots_loi_anh_s", Name: "Lôi Ảnh Chiến Ngoa", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Lôi điện bám gót.", Stats: map[string]int{"speed": 100, "crit_rate": 2}, SellPrice: 8500},
		"eq_boots_phuong_vu_ss":     {ID: "eq_boots_phuong_vu_ss", Name: "Phượng Vũ Linh Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Gót chân khinh khí.", Stats: map[string]int{"speed": 150, "dodge_rate": 6}, SellPrice: 22000},
		"eq_boots_tinh_ha_ss":       {ID: "eq_boots_tinh_ha_ss", Name: "Tinh Hà Bộ Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Đạp nát tinh hà.", Stats: map[string]int{"speed": 170, "stamina_regen_bonus": 2}, SellPrice: 25000},
		"eq_boots_hu_khong_ss":      {ID: "eq_boots_hu_khong_ss", Name: "Hư Không Tật Ảnh Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Dịch chuyển không gian.", Stats: map[string]int{"speed": 190, "stamina_max": 40}, SellPrice: 28000},
		"eq_boots_thien_hanh_sss":   {ID: "eq_boots_thien_hanh_sss", Name: "Thiên Hành Đạo Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Thay trời hành đạo.", Stats: map[string]int{"speed": 280, "dodge_rate": 9}, SellPrice: 85000},
		"eq_boots_cuu_thien_sss":    {ID: "eq_boots_cuu_thien_sss", Name: "Cửu Thiên Lôi Bộ", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Thần lôi nhập thể.", Stats: map[string]int{"speed": 320, "crit_rate": 5}, SellPrice: 95000},
		"eq_boots_luan_hoi_sssp":    {ID: "eq_boots_luan_hoi_sssp", Name: "Luân Hồi Vô Ảnh Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Chạy thoát khỏi luân hồi.", Stats: map[string]int{"speed": 500, "dodge_rate": 15}, SellPrice: 350000},
		"eq_boots_nghich_menh_sssp": {ID: "eq_boots_nghich_menh_sssp", Name: "Nghịch Mệnh Đạp Thiên Ngoa", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Giẫm nát thiên mệnh.", Stats: map[string]int{"speed": 620, "stamina_regen_bonus": 5, "dodge_rate": 18}, SellPrice: 500000},
	})
}
