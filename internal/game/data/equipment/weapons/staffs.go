// File: internal/data/equipment/weapons/staffs.go
package weapons

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_weapon_thiet_truong_d":      {ID: "eq_weapon_thiet_truong_d", Name: "Thiết Trượng", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Gậy sắt trơn tuột.", Stats: map[string]int{"attack": 12}, SellPrice: 30},
		"eq_weapon_thanh_dang_truong_d": {ID: "eq_weapon_thanh_dang_truong_d", Name: "Thanh Đằng Trượng", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Làm từ dây leo dai chắc.", Stats: map[string]int{"attack": 10, "cultivation_power": 5}, SellPrice: 35},
		"eq_weapon_kho_moc_truong_c":    {ID: "eq_weapon_kho_moc_truong_c", Name: "Khô Mộc Trượng", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Gỗ khô chứa tà khí.", Stats: map[string]int{"attack": 25, "cultivation_power": 15}, SellPrice: 120},
		"eq_weapon_xa_cot_truong_c":     {ID: "eq_weapon_xa_cot_truong_c", Name: "Xà Cốt Trượng", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Làm từ xương rắn khổng lồ.", Stats: map[string]int{"attack": 32}, SellPrice: 160},
		"eq_weapon_linh_moc_b":          {ID: "eq_weapon_linh_moc_b", Name: "Linh Mộc Trượng", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Trượng gỗ tụ linh.", Stats: map[string]int{"attack": 48, "cultivation_power": 30}, SellPrice: 400},
		"eq_weapon_phan_thien_truong_b": {ID: "eq_weapon_phan_thien_truong_b", Name: "Phần Thiên Trượng", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Đầu trượng có đốm lửa vĩnh cửu.", Stats: map[string]int{"attack": 75, "crit_rate": 3}, SellPrice: 650},
		"eq_weapon_bach_ngoc_a":         {ID: "eq_weapon_bach_ngoc_a", Name: "Bạch Ngọc Linh Trượng", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Gia tăng linh lực mạnh.", Stats: map[string]int{"attack": 90, "cultivation_power": 120}, SellPrice: 2000},
		"eq_weapon_thien_ma_truong_s":   {ID: "eq_weapon_thien_ma_truong_s", Name: "Thiên Ma Tế Cốt Trượng", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Triệu hoán linh hồn thiên ma.", Stats: map[string]int{"attack": 230, "cultivation_power": 200, "crit_damage": 10}, SellPrice: 8500},
		"eq_weapon_tinh_ha_ss":          {ID: "eq_weapon_tinh_ha_ss", Name: "Tinh Hà Pháp Trượng", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Chứa cả dải ngân hà.", Stats: map[string]int{"attack": 270, "cultivation_power": 400}, SellPrice: 22000},
	})
}
