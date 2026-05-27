package cultivation

import (
	"strings"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

type RealmStep string

const (
	RealmStepTungHoanh RealmStep = "tung_hoanh"
	RealmStepPhiThien  RealmStep = "phi_thien"
	RealmStepVoBien    RealmStep = "vo_bien"
	RealmStepSieuThoat RealmStep = "sieu_thoat"
)

type RealmDefinition struct {
	ID        string
	Name      string
	Step      RealmStep
	Order     int
	MaxLevel  int
	Aliases   []string
	BaseStats combat.CombatStats
	PerLevel  combat.CombatStats
}

var RealmRegistry = map[string]RealmDefinition{}
var RealmOrderList []string

func init() {
	realms := []RealmDefinition{
		{ID: "pham_nhan", Name: "Phàm Nhân", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 100, ATK: 10, DEF: 5, Speed: 90}, PerLevel: combat.CombatStats{MaxHP: 10, ATK: 1, DEF: 1}},

		// Bước 1: Tung Hoành
		{ID: "ngung_khi", Name: "Ngưng Khí", Step: RealmStepTungHoanh, MaxLevel: 10, Aliases: []string{"Luyện Khí", "Linh Động Kỳ"}, BaseStats: combat.CombatStats{MaxHP: 500, ATK: 50, DEF: 20, Speed: 100}, PerLevel: combat.CombatStats{MaxHP: 50, ATK: 5, DEF: 2}},
		{ID: "truc_co", Name: "Trúc Cơ", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 2000, ATK: 200, DEF: 100, Speed: 105}, PerLevel: combat.CombatStats{MaxHP: 200, ATK: 20, DEF: 10}},
		{ID: "ket_dan", Name: "Kết Đan", Step: RealmStepTungHoanh, MaxLevel: 10, Aliases: []string{"Kim Đan"}, BaseStats: combat.CombatStats{MaxHP: 8000, ATK: 800, DEF: 400, Speed: 110}, PerLevel: combat.CombatStats{MaxHP: 800, ATK: 80, DEF: 40}},
		{ID: "nguyen_anh", Name: "Nguyên Anh", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 30000, ATK: 3000, DEF: 1500, Speed: 115}, PerLevel: combat.CombatStats{MaxHP: 3000, ATK: 300, DEF: 150}},
		{ID: "hoa_than", Name: "Hóa Thần", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 100000, ATK: 10000, DEF: 5000, Speed: 120}, PerLevel: combat.CombatStats{MaxHP: 10000, ATK: 1000, DEF: 500}},
		{ID: "anh_bien", Name: "Anh Biến", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 300000, ATK: 30000, DEF: 15000, Speed: 125}, PerLevel: combat.CombatStats{MaxHP: 30000, ATK: 3000, DEF: 1500}},
		{ID: "van_dinh", Name: "Vấn Đỉnh", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 800000, ATK: 80000, DEF: 40000, Speed: 130}, PerLevel: combat.CombatStats{MaxHP: 80000, ATK: 8000, DEF: 4000}},
		{ID: "am_hu", Name: "Âm Hư", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 2000000, ATK: 200000, DEF: 100000, Speed: 135}, PerLevel: combat.CombatStats{MaxHP: 200000, ATK: 20000, DEF: 10000}},
		{ID: "duong_thuc", Name: "Dương Thực", Step: RealmStepTungHoanh, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 5000000, ATK: 500000, DEF: 250000, Speed: 140}, PerLevel: combat.CombatStats{MaxHP: 500000, ATK: 50000, DEF: 25000}},

		// Bước 2: Phi Thiên
		{ID: "khuy_niet", Name: "Khuy Niết", Step: RealmStepPhiThien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 15000000, ATK: 1500000, DEF: 750000, Speed: 145}, PerLevel: combat.CombatStats{MaxHP: 1500000, ATK: 150000, DEF: 75000}},
		{ID: "tinh_niet", Name: "Tịnh Niết", Step: RealmStepPhiThien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 40000000, ATK: 4000000, DEF: 2000000, Speed: 150}, PerLevel: combat.CombatStats{MaxHP: 4000000, ATK: 400000, DEF: 200000}},
		{ID: "toai_niet", Name: "Toái Niết", Step: RealmStepPhiThien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 100000000, ATK: 10000000, DEF: 5000000, Speed: 155}, PerLevel: combat.CombatStats{MaxHP: 10000000, ATK: 1000000, DEF: 500000}},
		{ID: "thien_nhan_ngu_suy", Name: "Thiên Nhân Ngũ Suy", Step: RealmStepPhiThien, MaxLevel: 5, Aliases: []string{"Phá Không Ngũ Chỉ"}, BaseStats: combat.CombatStats{MaxHP: 300000000, ATK: 30000000, DEF: 15000000, Speed: 160}, PerLevel: combat.CombatStats{MaxHP: 60000000, ATK: 6000000, DEF: 3000000}},

		// Bước 3: Vô Biên
		{ID: "khong_niet", Name: "Không Niết", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 1000000000, ATK: 100000000, DEF: 50000000, Speed: 165}, PerLevel: combat.CombatStats{MaxHP: 100000000, ATK: 10000000, DEF: 5000000}},
		{ID: "khong_linh", Name: "Không Linh", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 3000000000, ATK: 300000000, DEF: 150000000, Speed: 170}, PerLevel: combat.CombatStats{MaxHP: 300000000, ATK: 30000000, DEF: 15000000}},
		{ID: "khong_huyen", Name: "Không Huyền", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 10000000000, ATK: 1000000000, DEF: 500000000, Speed: 175}, PerLevel: combat.CombatStats{MaxHP: 1000000000, ATK: 100000000, DEF: 50000000}},

		{ID: "ngoai_kiep", Name: "Ngoại Kiếp", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 30000000000, ATK: 3000000000, DEF: 1500000000, Speed: 180}, PerLevel: combat.CombatStats{MaxHP: 3000000000, ATK: 300000000, DEF: 150000000}},
		{ID: "noi_kiep", Name: "Nội Kiếp", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 90000000000, ATK: 9000000000, DEF: 4500000000, Speed: 185}, PerLevel: combat.CombatStats{MaxHP: 9000000000, ATK: 900000000, DEF: 450000000}},
		{ID: "hon_kiep", Name: "Hồn Kiếp", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 270000000000, ATK: 27000000000, DEF: 13500000000, Speed: 190}, PerLevel: combat.CombatStats{MaxHP: 27000000000, ATK: 2700000000, DEF: 1350000000}},

		{ID: "dai_ton", Name: "Đại Tôn", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 800000000000, ATK: 80000000000, DEF: 40000000000, Speed: 195}, PerLevel: combat.CombatStats{MaxHP: 80000000000, ATK: 8000000000, DEF: 4000000000}},
		{ID: "kim_ton", Name: "Kim Tôn", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 2400000000000, ATK: 240000000000, DEF: 120000000000, Speed: 200}, PerLevel: combat.CombatStats{MaxHP: 240000000000, ATK: 24000000000, DEF: 12000000000}},
		{ID: "thien_ton", Name: "Thiên Tôn", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 7200000000000, ATK: 720000000000, DEF: 360000000000, Speed: 205}, PerLevel: combat.CombatStats{MaxHP: 720000000000, ATK: 72000000000, DEF: 36000000000}},
		{ID: "duoc_thien_ton", Name: "Dược Thiên Tôn", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 21600000000000, ATK: 2160000000000, DEF: 1080000000000, Speed: 210}, PerLevel: combat.CombatStats{MaxHP: 2160000000000, ATK: 216000000000, DEF: 108000000000}},
		{ID: "dai_thien_ton", Name: "Đại Thiên Tôn", Step: RealmStepVoBien, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 64800000000000, ATK: 6480000000000, DEF: 3240000000000, Speed: 215}, PerLevel: combat.CombatStats{MaxHP: 6480000000000, ATK: 648000000000, DEF: 324000000000}},

		{ID: "dap_thien_nhat_kieu", Name: "Đạp Thiên Nhất Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 100000000000000, ATK: 10000000000000, DEF: 5000000000000, Speed: 220}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_nhi_kieu", Name: "Đạp Thiên Nhị Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 300000000000000, ATK: 30000000000000, DEF: 15000000000000, Speed: 225}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_tam_kieu", Name: "Đạp Thiên Tam Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 900000000000000, ATK: 90000000000000, DEF: 45000000000000, Speed: 230}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_tu_kieu", Name: "Đạp Thiên Tứ Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 2700000000000000, ATK: 270000000000000, DEF: 135000000000000, Speed: 235}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_ngu_kieu", Name: "Đạp Thiên Ngũ Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 8100000000000000, ATK: 810000000000000, DEF: 405000000000000, Speed: 240}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_luc_kieu", Name: "Đạp Thiên Lục Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 24300000000000000, ATK: 2430000000000000, DEF: 1215000000000000, Speed: 245}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_that_kieu", Name: "Đạp Thiên Thất Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 72900000000000000, ATK: 7290000000000000, DEF: 3645000000000000, Speed: 250}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_bat_kieu", Name: "Đạp Thiên Bát Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 218700000000000000, ATK: 21870000000000000, DEF: 10935000000000000, Speed: 255}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},
		{ID: "dap_thien_cuu_kieu", Name: "Đạp Thiên Cửu Kiều", Step: RealmStepVoBien, MaxLevel: 1, BaseStats: combat.CombatStats{MaxHP: 656100000000000000, ATK: 65610000000000000, DEF: 32805000000000000, Speed: 260}, PerLevel: combat.CombatStats{MaxHP: 0, ATK: 0, DEF: 0}},

		// Bước 4: Siêu Thoát
		{ID: "sieu_thoat", Name: "Siêu Thoát Cảnh", Step: RealmStepSieuThoat, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 2000000000000000000, ATK: 200000000000000000, DEF: 100000000000000000, Speed: 265}, PerLevel: combat.CombatStats{MaxHP: 200000000000000000, ATK: 20000000000000000, DEF: 10000000000000000}},
		{ID: "dap_thien_canh", Name: "Đạp Thiên Cảnh", Step: RealmStepSieuThoat, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 6000000000000000000, ATK: 600000000000000000, DEF: 300000000000000000, Speed: 270}, PerLevel: combat.CombatStats{MaxHP: 600000000000000000, ATK: 60000000000000000, DEF: 30000000000000000}},
		{ID: "khong_diet", Name: "Không Diệt", Step: RealmStepSieuThoat, MaxLevel: 10, BaseStats: combat.CombatStats{MaxHP: 9000000000000000000, ATK: 900000000000000000, DEF: 450000000000000000, Speed: 275}, PerLevel: combat.CombatStats{MaxHP: 900000000000000000, ATK: 90000000000000000, DEF: 45000000000000000}},
	}

	for i, r := range realms {
		r.Order = i
		RealmRegistry[r.ID] = r
		RealmOrderList = append(RealmOrderList, r.ID)
	}
}

// NormalizeRealmID chuẩn hóa tên hoặc alias thành ID chuẩn.
func NormalizeRealmID(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return "pham_nhan"
	}
	// Dò theo ID
	if _, ok := RealmRegistry[input]; ok {
		return input
	}
	// Dò theo Name hoặc Aliases
	lowerInput := strings.ToLower(input)
	for _, def := range RealmRegistry {
		if strings.ToLower(def.Name) == lowerInput {
			return def.ID
		}
		for _, alias := range def.Aliases {
			if strings.ToLower(alias) == lowerInput {
				return def.ID
			}
		}
	}
	return "pham_nhan"
}
