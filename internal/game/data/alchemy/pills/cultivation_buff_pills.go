// File: internal/data/alchemy/pills/cultivation_buff_pills.go
package pills

import "github.com/whiskey/tu-tien-bot/internal/game/item"

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_buff_thanh_tam_d":   {ID: "pill_buff_thanh_tam_d", Name: "Thanh Tâm Tán", Type: item.TypePill, Rarity: item.RarityD, Stackable: true, MaxStack: 99, Usable: true, Description: "Làm mát nguyên thần, hồi 5 tâm cảnh.", Effects: map[string]int{"mind": 5}, SellPrice: 20},
		"pill_buff_tinh_than_d":   {ID: "pill_buff_tinh_than_d", Name: "Tĩnh Thần Đan", Type: item.TypePill, Rarity: item.RarityD, Stackable: true, MaxStack: 99, Usable: true, Description: "Giữ tinh thần ổn định, hồi 8 tâm cảnh.", Effects: map[string]int{"mind": 8}, SellPrice: 40},
		"pill_buff_thanh_tam_c":   {ID: "pill_buff_thanh_tam_c", Name: "Thanh Tâm Đan", Type: item.TypePill, Rarity: item.RarityC, Stackable: true, MaxStack: 99, Usable: true, Description: "Bảo hộ tâm cảnh, hồi 15 tâm cảnh.", Effects: map[string]int{"mind": 15}, SellPrice: 100},
		"pill_buff_dinh_hon_c":    {ID: "pill_buff_dinh_hon_c", Name: "Định Hồn Đan", Type: item.TypePill, Rarity: item.RarityC, Stackable: true, MaxStack: 99, Usable: true, Description: "Cố định linh hồn, hồi 20 tâm cảnh.", Effects: map[string]int{"mind": 20}, SellPrice: 150},
		"pill_buff_chan_hon_b":    {ID: "pill_buff_chan_hon_b", Name: "Chấn Hồn Đan", Type: item.TypePill, Rarity: item.RarityB, Stackable: true, MaxStack: 99, Usable: true, Description: "Hồi 35 điểm tâm cảnh.", Effects: map[string]int{"mind": 35}, SellPrice: 350},
		"pill_buff_tay_tuy_b":     {ID: "pill_buff_tay_tuy_b", Name: "Tẩy Tủy Đan", Type: item.TypePill, Rarity: item.RarityB, Stackable: true, MaxStack: 99, Usable: true, Description: "Thanh tẩy kinh mạch, hồi 40 tâm cảnh.", Effects: map[string]int{"mind": 40}, SellPrice: 450},
		"pill_buff_ngo_dao_a":     {ID: "pill_buff_ngo_dao_a", Name: "Tiểu Ngộ Đạo Đan", Type: item.TypePill, Rarity: item.RarityA, Stackable: true, MaxStack: 99, Usable: true, Description: "Giúp ngộ đạo, hồi 50 tâm cảnh.", Effects: map[string]int{"mind": 50}, SellPrice: 1000},
		"pill_buff_bang_tam_a":    {ID: "pill_buff_bang_tam_a", Name: "Băng Tâm Quyết Đan", Type: item.TypePill, Rarity: item.RarityA, Stackable: true, MaxStack: 99, Usable: true, Description: "Tâm như băng, hồi 65 tâm cảnh.", Effects: map[string]int{"mind": 65}, SellPrice: 1400},
		"pill_buff_vo_niem_s":     {ID: "pill_buff_vo_niem_s", Name: "Vô Niệm Tâm Đan", Type: item.TypePill, Rarity: item.RarityS, Stackable: true, MaxStack: 99, Usable: true, Description: "Đoạn tuyệt tạp niệm, hồi 100 tâm cảnh.", Effects: map[string]int{"mind": 100}, SellPrice: 3000},
		"pill_buff_bo_de_s":       {ID: "pill_buff_bo_de_s", Name: "Bồ Đề Tử", Type: item.TypePill, Rarity: item.RarityS, Stackable: true, MaxStack: 99, Usable: true, Description: "Quả bồ đề ngàn năm, hồi 100 tâm cảnh.", Effects: map[string]int{"mind": 100}, SellPrice: 4000},
		"pill_buff_bich_hai_ss":   {ID: "pill_buff_bich_hai_ss", Name: "Bích Hải Triều Sinh Đan", Type: item.TypePill, Rarity: item.RaritySS, Stackable: true, MaxStack: 99, Usable: true, Description: "Tâm bình như biển, hồi 100 tâm cảnh.", Effects: map[string]int{"mind": 100}, SellPrice: 10000},
		"pill_buff_thai_hu_ss":    {ID: "pill_buff_thai_hu_ss", Name: "Thái Hư Tĩnh Tâm Đan", Type: item.TypePill, Rarity: item.RaritySS, Stackable: true, MaxStack: 99, Usable: true, Description: "Vào cảnh thái hư, hồi 100 tâm cảnh.", Effects: map[string]int{"mind": 100}, SellPrice: 12000},
		"pill_buff_thien_dao_sss": {ID: "pill_buff_thien_dao_sss", Name: "Thiên Đạo Vô Tâm Đan", Type: item.TypePill, Rarity: item.RaritySSS, Stackable: true, MaxStack: 99, Usable: true, Description: "Trở thành vô tâm như thiên đạo, hồi 100 tâm cảnh.", Effects: map[string]int{"mind": 100}, SellPrice: 50000},
		"pill_buff_dai_dao_sssp":  {ID: "pill_buff_dai_dao_sssp", Name: "Đại Đạo Bản Nguyên Đan", Type: item.TypePill, Rarity: item.RaritySSSP, Stackable: true, MaxStack: 99, Usable: true, Description: "Hòa mình vào đại đạo, bất sinh bất diệt.", Effects: map[string]int{"mind": 100}, SellPrice: 200000},
	})
}
