// File: internal/data/alchemy/pills/breakthrough_pills.go
package pills

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_brk_pha_canh_d":    {ID: "pill_brk_pha_canh_d", Name: "Tiểu Phá Cảnh Đan", Type: item.TypePill, Rarity: item.RarityD, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 2% tỷ lệ đột phá.", Effects: map[string]int{"breakthrough_chance": 2}, SellPrice: 50},
		"pill_brk_pha_canh_c":    {ID: "pill_brk_pha_canh_c", Name: "Phá Cảnh Đan", Type: item.TypePill, Rarity: item.RarityC, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 4% tỷ lệ đột phá.", Effects: map[string]int{"breakthrough_chance": 4}, SellPrice: 150},
		"pill_brk_ho_mach_b":     {ID: "pill_brk_ho_mach_b", Name: "Hộ Mạch Đan", Type: item.TypePill, Rarity: item.RarityB, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 7% đột phá, giảm penalty thất bại.", Effects: map[string]int{"breakthrough_chance": 7, "fail_reduce": 2}, SellPrice: 400},
		"pill_brk_huyen_mon_a":   {ID: "pill_brk_huyen_mon_a", Name: "Huyền Môn Phá Cảnh Đan", Type: item.TypePill, Rarity: item.RarityA, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 10% tỷ lệ đột phá.", Effects: map[string]int{"breakthrough_chance": 10}, SellPrice: 1200},
		"pill_brk_thien_kiep_s":  {ID: "pill_brk_thien_kiep_s", Name: "Thiên Kiếp Hộ Mạch Đan", Type: item.TypePill, Rarity: item.RarityS, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 15% đột phá, giảm 5% thất bại.", Effects: map[string]int{"breakthrough_chance": 15, "fail_reduce": 5}, SellPrice: 4000},
		"pill_brk_ngo_dao_ss":    {ID: "pill_brk_ngo_dao_ss", Name: "Ngộ Đạo Đan", Type: item.TypePill, Rarity: item.RaritySS, Stackable: true, MaxStack: 99, Usable: false, Description: "Tăng 20% đột phá.", Effects: map[string]int{"breakthrough_chance": 20}, SellPrice: 15000},
		"pill_brk_luan_hoi_sssp": {ID: "pill_brk_luan_hoi_sssp", Name: "Luân Hồi Ngộ Đạo Đan", Type: item.TypePill, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 99, Usable: false, Description: "Cam kết nghịch thiên, tăng 35% đột phá.", Effects: map[string]int{"breakthrough_chance": 35}, SellPrice: 100000},
	})
}
