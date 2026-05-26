// File: internal/data/alchemy/pills/exp_pills.go
package pills

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_exp_tu_khi_d":         {ID: "pill_exp_tu_khi_d", Name: "Tụ Khí Đan", Type: item.TypePill, Rarity: item.RarityD, Stackable: true, MaxStack: 99, Usable: true, Description: "Đan dược cơ bản nhất, tăng 100 tu vi.", Effects: map[string]int{"exp": 100}, SellPrice: 10},
		"pill_exp_ngung_khi_c":      {ID: "pill_exp_ngung_khi_c", Name: "Ngưng Khí Đan", Type: item.TypePill, Rarity: item.RarityC, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 250 tu vi.", Effects: map[string]int{"exp": 250}, SellPrice: 30},
		"pill_exp_truc_co_b":        {ID: "pill_exp_truc_co_b", Name: "Trúc Cơ Đan", Type: item.TypePill, Rarity: item.RarityB, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 600 tu vi.", Effects: map[string]int{"exp": 600}, SellPrice: 100},
		"pill_exp_linh_nguyen_a":    {ID: "pill_exp_linh_nguyen_a", Name: "Linh Nguyên Đan", Type: item.TypePill, Rarity: item.RarityA, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 1500 tu vi.", Effects: map[string]int{"exp": 1500}, SellPrice: 350},
		"pill_exp_huyen_nguyen_s":   {ID: "pill_exp_huyen_nguyen_s", Name: "Huyền Nguyên Đan", Type: item.TypePill, Rarity: item.RarityS, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 4000 tu vi.", Effects: map[string]int{"exp": 4000}, SellPrice: 1000},
		"pill_exp_thien_nguyen_ss":  {ID: "pill_exp_thien_nguyen_ss", Name: "Thiên Nguyên Đan", Type: item.TypePill, Rarity: item.RaritySS, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 10000 tu vi.", Effects: map[string]int{"exp": 10000}, SellPrice: 3500},
		"pill_exp_tien_nguyen_sss":  {ID: "pill_exp_tien_nguyen_sss", Name: "Tiên Nguyên Đan", Type: item.TypePill, Rarity: item.RaritySSS, Stackable: true, MaxStack: 99, Usable: true, Description: "Tăng 25000 tu vi.", Effects: map[string]int{"exp": 25000}, SellPrice: 10000},
		"pill_exp_nghich_menh_sssp": {ID: "pill_exp_nghich_menh_sssp", Name: "Nghịch Mệnh Tiên Đan", Type: item.TypePill, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 99, Usable: true, Description: "Nghịch thiên đoạt mệnh, tăng 75000 tu vi.", Effects: map[string]int{"exp": 75000}, SellPrice: 50000},
	})
}
