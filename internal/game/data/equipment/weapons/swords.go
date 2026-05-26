// File: internal/data/equipment/weapons/swords.go
package weapons

import (
	"github.com/whiskey/tu-tien-bot/internal/game/data"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

func init() {
	data.RegisterItems(map[string]item.ItemDefinition{
		"eq_weapon_moc_kiem_d":       {ID: "eq_weapon_moc_kiem_d", Name: "Mộc Kiếm", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Vũ khí thô sơ.", Stats: map[string]int{"attack": 8}, SellPrice: 20},
		"eq_weapon_thanh_thiet_d":    {ID: "eq_weapon_thanh_thiet_d", Name: "Thanh Thiết Kiếm", Type: item.TypeEquipment, Rarity: item.RarityD, Stackable: false, MaxStack: 1, Description: "Kiếm sắt bén.", Stats: map[string]int{"attack": 15}, SellPrice: 40},
		"eq_weapon_xich_luyen_c":     {ID: "eq_weapon_xich_luyen_c", Name: "Xích Luyện Kiếm", Type: item.TypeEquipment, Rarity: item.RarityC, Stackable: false, MaxStack: 1, Description: "Kiếm nung từ xích sa.", Stats: map[string]int{"attack": 35, "crit_rate": 1}, SellPrice: 150},
		"eq_weapon_huyen_thiet_b":    {ID: "eq_weapon_huyen_thiet_b", Name: "Huyền Thiết Kiếm", Type: item.TypeEquipment, Rarity: item.RarityB, Stackable: false, MaxStack: 1, Description: "Kiếm huyền thiết cứng rắn.", Stats: map[string]int{"attack": 75, "crit_rate": 2}, SellPrice: 600},
		"eq_weapon_tu_dien_a":        {ID: "eq_weapon_tu_dien_a", Name: "Tử Điện Kiếm", Type: item.TypeEquipment, Rarity: item.RarityA, Stackable: false, MaxStack: 1, Description: "Lưỡi kiếm chứa lôi điện.", Stats: map[string]int{"attack": 110, "crit_rate": 3}, SellPrice: 1500},
		"eq_weapon_huyet_anh_s":      {ID: "eq_weapon_huyet_anh_s", Name: "Huyết Ảnh Ma Kiếm", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Ma kiếm hút máu.", Stats: map[string]int{"attack": 180, "crit_rate": 5}, SellPrice: 6000},
		"eq_weapon_thai_am_s":        {ID: "eq_weapon_thai_am_s", Name: "Thái Âm Kiếm", Type: item.TypeEquipment, Rarity: item.RarityS, Stackable: false, MaxStack: 1, Description: "Mang sức mạnh của Mặt Trăng.", Stats: map[string]int{"attack": 190, "stamina_cost_reduce": 3}, SellPrice: 8000},
		"eq_weapon_phuong_hoa_ss":    {ID: "eq_weapon_phuong_hoa_ss", Name: "Phượng Hỏa Kiếm", Type: item.TypeEquipment, Rarity: item.RaritySS, Stackable: false, MaxStack: 1, Description: "Cháy mãi không tắt.", Stats: map[string]int{"attack": 350, "crit_damage": 25}, SellPrice: 25000},
		"eq_weapon_hu_khong_sss":     {ID: "eq_weapon_hu_khong_sss", Name: "Hư Không Đạo Kiếm", Type: item.TypeEquipment, Rarity: item.RaritySSS, Stackable: false, MaxStack: 1, Description: "Chém rách hư không.", Stats: map[string]int{"attack": 520, "crit_rate": 10, "crit_damage": 35}, SellPrice: 80000},
		"eq_weapon_luan_hoi_sssp":    {ID: "eq_weapon_luan_hoi_sssp", Name: "Luân Hồi Cổ Kiếm", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Một kiếm đứt luân hồi.", Stats: map[string]int{"attack": 850, "crit_rate": 15, "crit_damage": 50}, SellPrice: 300000},
		"eq_weapon_nghich_menh_sssp": {ID: "eq_weapon_nghich_menh_sssp", Name: "Nghịch Mệnh Trảm Thiên Kiếm", Type: item.TypeEquipment, Rarity: item.RaritySSSP, Stackable: false, MaxStack: 1, Description: "Vũ khí chí tôn, nghịch lại thiên mệnh.", Stats: map[string]int{"attack": 1000, "cultivation_power": 1000, "breakthrough_chance": 5}, SellPrice: 500000},
	})
}
