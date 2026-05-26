// File: internal/data/equipment/armor/armor.go
package armor

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_armor_hac_thiet_c":      {ID: "eq_armor_hac_thiet_c", Name: "Hắc Thiết Giáp", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Giáp sắt đen.", Stats: map[string]int{"defense": 18, "hp": 80}, SellPrice: 100},
		"eq_armor_bach_ngoc_b":      {ID: "eq_armor_bach_ngoc_b", Name: "Bạch Ngọc Khinh Giáp", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Nhẹ nhưng cực cứng.", Stats: map[string]int{"defense": 35, "hp": 180}, SellPrice: 400},
		"eq_armor_huyen_thiet_b":    {ID: "eq_armor_huyen_thiet_b", Name: "Huyền Thiết Trọng Giáp", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Nặng nề, thủ cao.", Stats: map[string]int{"defense": 50, "hp": 250}, SellPrice: 600},
		"eq_armor_han_nguyet_a":     {ID: "eq_armor_han_nguyet_a", Name: "Hàn Nguyệt Linh Giáp", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Bảo vệ nguyên thần.", Stats: map[string]int{"defense": 85, "breakthrough_protection": 2}, SellPrice: 2000},
		"eq_armor_xich_hoa_a":       {ID: "eq_armor_xich_hoa_a", Name: "Xích Hỏa Chiến Giáp", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Rực lửa chiến đấu.", Stats: map[string]int{"defense": 95, "hp": 500}, SellPrice: 2200},
		"eq_armor_cuu_duong_s":      {ID: "eq_armor_cuu_duong_s", Name: "Cửu Dương Linh Giáp", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Nóng rực như 9 mặt trời.", Stats: map[string]int{"defense": 160, "hp": 900}, SellPrice: 7500},
		"eq_armor_long_lan_s":       {ID: "eq_armor_long_lan_s", Name: "Long Lân Hộ Giáp", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Làm từ vảy rồng.", Stats: map[string]int{"defense": 190, "damage_reduce": 3}, SellPrice: 9000},
		"eq_armor_tinh_ha_ss":       {ID: "eq_armor_tinh_ha_ss", Name: "Tinh Hà Hộ Thân Giáp", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Lấp lánh ánh sao.", Stats: map[string]int{"defense": 320, "hp": 1800}, SellPrice: 25000},
		"eq_armor_hu_khong_sss":     {ID: "eq_armor_hu_khong_sss", Name: "Hư Không Tiên Giáp", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Ẩn hiện trong hư không.", Stats: map[string]int{"defense": 480, "damage_reduce": 8}, SellPrice: 85000},
		"eq_armor_luan_hoi_sssp":    {ID: "eq_armor_luan_hoi_sssp", Name: "Luân Hồi Bất Diệt Giáp", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Bất tử bất diệt.", Stats: map[string]int{"defense": 800, "hp": 5000}, SellPrice: 350000},
		"eq_armor_nghich_menh_sssp": {ID: "eq_armor_nghich_menh_sssp", Name: "Nghịch Mệnh Vạn Đạo Bào", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Mặc vào không sợ thiên đạo.", Stats: map[string]int{"defense": 950, "damage_reduce": 12, "breakthrough_protection": 12}, SellPrice: 500000},
	})
}
