package npc

import "errors"

var (
	ErrShopNotFound   = errors.New("không tìm thấy thương hội này")
	ErrItemNotForSale = errors.New("vật phẩm không được bán tại đây")
)

// Registry là cơ sở dữ liệu tĩnh của các cửa hàng NPC
var Registry = map[string]ShopDef{
	"van_bao_cac": {
		ID:   "van_bao_cac",
		Name: "Vạn Bảo Các (Thương Hội)",
		Items: map[string]ItemDef{
			// --- ĐAN DƯỢC ---
			"pill_exp_tu_khi_d":  {ItemID: "pill_exp_tu_khi_d", BuyPrice: 50, SellPrice: 10, Stock: -1, Category: "Đan Dược", Enabled: true},
			"pill_stm_hoi_luc_d": {ItemID: "pill_stm_hoi_luc_d", BuyPrice: 30, SellPrice: 5, Stock: -1, Category: "Đan Dược", Enabled: true},

			// --- NGUYÊN LIỆU ---
			"mat_enhance_hac_thiet_d": {ItemID: "mat_enhance_hac_thiet_d", BuyPrice: 100, SellPrice: 20, Stock: -1, Category: "Nguyên Liệu", Enabled: true},
			"mat_enhance_hac_thiet_c": {ItemID: "mat_enhance_hac_thiet_c", BuyPrice: 500, SellPrice: 100, Stock: 5, Category: "Nguyên Liệu", Enabled: true},

			// --- TRANG BỊ TÂN THỦ ---
			"eq_weapon_moc_kiem_d":     {ItemID: "eq_weapon_moc_kiem_d", BuyPrice: 200, SellPrice: 40, Stock: -1, Category: "Trang Bị", Enabled: true},
			"eq_armor_vai_tho_d":       {ItemID: "eq_armor_vai_tho_d", BuyPrice: 200, SellPrice: 40, Stock: -1, Category: "Trang Bị", Enabled: true},
			"eq_boots_vai_tho_d":       {ItemID: "eq_boots_vai_tho_d", BuyPrice: 150, SellPrice: 30, Stock: -1, Category: "Trang Bị", Enabled: true},
			"eq_artifact_guong_dong_d": {ItemID: "eq_artifact_guong_dong_d", BuyPrice: 800, SellPrice: 150, Stock: 1, Category: "Pháp Bảo", Enabled: true},
		},
	},
}

func GetShop(shopID string) (ShopDef, error) {
	shop, ok := Registry[shopID]
	if !ok {
		return ShopDef{}, ErrShopNotFound
	}
	return shop, nil
}

func (s *ShopDef) GetItem(itemID string) (ItemDef, error) {
	item, ok := s.Items[itemID]
	if !ok || !item.Enabled {
		return ItemDef{}, ErrItemNotForSale
	}
	return item, nil
}
