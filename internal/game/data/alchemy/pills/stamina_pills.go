// File: internal/data/alchemy/pills/stamina_pills.go
package pills

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_stm_hoi_luc_d":    {ID: "pill_stm_hoi_luc_d", Name: "Hồi Lực Đan", Type: item.TypePill, Rarity: item.RarityD, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 10 thể lực.", Effects: map[string]int{"stamina": 10}, SellPrice: 15},
		"pill_stm_thanh_moc_c":  {ID: "pill_stm_thanh_moc_c", Name: "Thanh Mộc Đan", Type: item.TypePill, Rarity: item.RarityC, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 25 thể lực.", Effects: map[string]int{"stamina": 25}, SellPrice: 40},
		"pill_stm_sinh_cot_b":   {ID: "pill_stm_sinh_cot_b", Name: "Sinh Cốt Đan", Type: item.TypePill, Rarity: item.RarityB, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 50 thể lực.", Effects: map[string]int{"stamina": 50}, SellPrice: 120},
		"pill_stm_linh_tuc_a":   {ID: "pill_stm_linh_tuc_a", Name: "Linh Tức Đan", Type: item.TypePill, Rarity: item.RarityA, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 100 thể lực.", Effects: map[string]int{"stamina": 100}, SellPrice: 300},
		"pill_stm_cuu_chuyen_s": {ID: "pill_stm_cuu_chuyen_s", Name: "Cửu Chuyển Hồi Lực Đan", Type: item.TypePill, Rarity: item.RarityS, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 200 thể lực.", Effects: map[string]int{"stamina": 200}, SellPrice: 800},
	})
}
