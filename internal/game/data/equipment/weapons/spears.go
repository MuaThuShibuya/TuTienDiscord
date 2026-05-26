// File: internal/data/equipment/weapons/spears.go
package weapons

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_weapon_truc_con_d":        {ID: "eq_weapon_truc_con_d", Name: "Trúc Côn", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Côn tre chắc chắn.", Stats: map[string]int{"attack": 7, "defense": 2}, SellPrice: 20},
		"eq_weapon_thiet_thuong_d":    {ID: "eq_weapon_thiet_thuong_d", Name: "Thiết Thương", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Thương sắt bình thường.", Stats: map[string]int{"attack": 14}, SellPrice: 40},
		"eq_weapon_ngan_lan_thuong_c": {ID: "eq_weapon_ngan_lan_thuong_c", Name: "Ngân Lân Thương", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Mũi thương lấp lánh như vảy bạc.", Stats: map[string]int{"attack": 30}, SellPrice: 100},
		"eq_weapon_pha_san_con_c":     {ID: "eq_weapon_pha_san_con_c", Name: "Phá Sơn Côn", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Côn nặng ngàn cân.", Stats: map[string]int{"attack": 35, "defense": 5}, SellPrice: 150},
		"eq_weapon_du_long_thuong_b":  {ID: "eq_weapon_du_long_thuong_b", Name: "Du Long Thương", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Đâm ra tựa rồng bơi.", Stats: map[string]int{"attack": 68, "speed": 5}, SellPrice: 450},
		"eq_weapon_pha_san_b":         {ID: "eq_weapon_pha_san_b", Name: "Phá Sơn Thương", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Thương nặng nề.", Stats: map[string]int{"attack": 60}, SellPrice: 450},
		"eq_weapon_toai_tinh_kich_b":  {ID: "eq_weapon_toai_tinh_kich_b", Name: "Toái Tinh Kích", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Kích đập nát tinh thạch.", Stats: map[string]int{"attack": 82, "crit_rate": 2}, SellPrice: 600},
		"eq_weapon_loi_dinh_thuong_a": {ID: "eq_weapon_loi_dinh_thuong_a", Name: "Lôi Đình Thương", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Thương mang sấm sét.", Stats: map[string]int{"attack": 145, "crit_damage": 10}, SellPrice: 1800},
		"eq_weapon_cuu_loi_s":         {ID: "eq_weapon_cuu_loi_s", Name: "Cửu Lôi Thương", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Triệu hoán cửu lôi.", Stats: map[string]int{"attack": 210, "crit_damage": 15}, SellPrice: 7500},
		"eq_weapon_long_van_ss":       {ID: "eq_weapon_long_van_ss", Name: "Long Văn Chiến Kích", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Khắc hình chân long.", Stats: map[string]int{"attack": 320, "crit_rate": 7}, SellPrice: 20000},
		"eq_weapon_thien_kiep_sss":    {ID: "eq_weapon_thien_kiep_sss", Name: "Thiên Kiếp Lôi Thương", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Mang sức mạnh thiên kiếp.", Stats: map[string]int{"attack": 580, "breakthrough_chance": 3}, SellPrice: 90000},
	})
}
