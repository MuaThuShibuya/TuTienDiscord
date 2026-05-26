// File: internal/data/equipment/weapons/sabers.go
package weapons

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_weapon_o_moc_dao_d":       {ID: "eq_weapon_o_moc_dao_d", Name: "Ô Mộc Đao", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Đao bằng gỗ cứng đen.", Stats: map[string]int{"attack": 9}, SellPrice: 20},
		"eq_weapon_thanh_cuong_dao_d": {ID: "eq_weapon_thanh_cuong_dao_d", Name: "Thanh Cương Đao", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Đao phàm trần nhưng rất bén.", Stats: map[string]int{"attack": 18}, SellPrice: 40},
		"eq_weapon_han_thiet_c":       {ID: "eq_weapon_han_thiet_c", Name: "Hàn Thiết Đao", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Đao tỏa khí lạnh.", Stats: map[string]int{"attack": 28}, SellPrice: 100},
		"eq_weapon_bang_tinh_dao_c":   {ID: "eq_weapon_bang_tinh_dao_c", Name: "Băng Tinh Phá Cốt Đao", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Khảm băng tinh thạch.", Stats: map[string]int{"attack": 38, "crit_rate": 1}, SellPrice: 150},
		"eq_weapon_cuong_sa_dao_b":    {ID: "eq_weapon_cuong_sa_dao_b", Name: "Cuồng Sa Đao", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Mang theo bão cát.", Stats: map[string]int{"attack": 70}, SellPrice: 450},
		"eq_weapon_cuu_hoan_dao_b":    {ID: "eq_weapon_cuu_hoan_dao_b", Name: "Cửu Hoàn Huyết Đao", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Đao chín vòng đẫm máu.", Stats: map[string]int{"attack": 85, "crit_damage": 5}, SellPrice: 600},
		"eq_weapon_xich_huyet_dao_a":  {ID: "eq_weapon_xich_huyet_dao_a", Name: "Xích Huyết Lôi Giao Đao", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Làm từ xương lôi giao.", Stats: map[string]int{"attack": 130, "crit_rate": 2}, SellPrice: 1500},
		"eq_weapon_thanh_minh_a":      {ID: "eq_weapon_thanh_minh_a", Name: "Thanh Minh Đao", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Đao phán xét âm dương.", Stats: map[string]int{"attack": 125, "crit_damage": 8}, SellPrice: 1800},
		"eq_weapon_minh_vuong_dao_s":  {ID: "eq_weapon_minh_vuong_dao_s", Name: "Cửu U Minh Vương Đao", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Tụ tập cửu u oán khí.", Stats: map[string]int{"attack": 240, "crit_rate": 5}, SellPrice: 6500},
		"eq_weapon_hong_hoang_dao_ss": {ID: "eq_weapon_hong_hoang_dao_ss", Name: "Hồng Hoang Tế Nhật Đao", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Đao chém rớt mặt trời.", Stats: map[string]int{"attack": 400, "crit_damage": 25}, SellPrice: 28000},
	})
}
