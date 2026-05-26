// File: internal/data/equipment/armor/robes.go
package armor

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_armor_vai_tho_d":      {ID: "eq_armor_vai_tho_d", Name: "Vải Thô Đạo Bào", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Áo vải cơ bản.", Stats: map[string]int{"defense": 5, "hp": 20}, SellPrice: 15},
		"eq_armor_thanh_sam_d":    {ID: "eq_armor_thanh_sam_d", Name: "Thanh Sam Tu Sĩ", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Áo tu sĩ phổ thông.", Stats: map[string]int{"defense": 8, "hp": 35}, SellPrice: 30},
		"eq_armor_linh_van_c":     {ID: "eq_armor_linh_van_c", Name: "Linh Văn Đạo Bào", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Có khắc linh văn.", Stats: map[string]int{"defense": 15, "stamina_max": 5}, SellPrice: 150},
		"eq_armor_thanh_moc_b":    {ID: "eq_armor_thanh_moc_b", Name: "Thanh Mộc Linh Bào", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Hấp thụ linh khí mộc.", Stats: map[string]int{"defense": 30, "stamina_max": 10}, SellPrice: 450},
		"eq_armor_tu_van_a":       {ID: "eq_armor_tu_van_a", Name: "Tử Vân Pháp Bào", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Tỏa ra khói tím.", Stats: map[string]int{"defense": 75, "hp": 400}, SellPrice: 1500},
		"eq_armor_huyen_am_s":     {ID: "eq_armor_huyen_am_s", Name: "Huyền Âm Đạo Bào", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Tuyệt đỉnh phòng thủ âm nhu.", Stats: map[string]int{"defense": 140, "stamina_max": 25}, SellPrice: 6000},
		"eq_armor_phuong_vu_ss":   {ID: "eq_armor_phuong_vu_ss", Name: "Phượng Vũ Tiên Bào", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Lông phượng hoàng dệt thành.", Stats: map[string]int{"defense": 280, "stamina_max": 50}, SellPrice: 22000},
		"eq_armor_van_phap_ss":    {ID: "eq_armor_van_phap_ss", Name: "Vạn Pháp Đạo Bào", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Miễn nhiễm vạn pháp.", Stats: map[string]int{"defense": 300, "breakthrough_protection": 5}, SellPrice: 28000},
		"eq_armor_thien_kiep_sss": {ID: "eq_armor_thien_kiep_sss", Name: "Thiên Kiếp Hộ Mệnh Bào", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Chống lại lôi kiếp.", Stats: map[string]int{"defense": 520, "breakthrough_protection": 8}, SellPrice: 95000},
	})
}
