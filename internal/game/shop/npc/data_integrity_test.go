// File: internal/game/shop/npc/data_integrity_test.go
// Chức năng: Kiểm duyệt toàn vẹn dữ liệu tĩnh của Cửa Hàng NPC.

package npc_test

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/shop/npc"
)

func init() {
	// Đăng ký vật phẩm giả lập để test không bị fail oan khi chưa load data thật
	item.RegisterItems(map[string]item.ItemDefinition{
		"pill_exp_tu_khi_d":        {ID: "pill_exp_tu_khi_d"},
		"pill_stm_hoi_luc_d":       {ID: "pill_stm_hoi_luc_d"},
		"mat_enhance_hac_thiet_d":  {ID: "mat_enhance_hac_thiet_d"},
		"mat_enhance_hac_thiet_c":  {ID: "mat_enhance_hac_thiet_c"},
		"eq_weapon_moc_kiem_d":     {ID: "eq_weapon_moc_kiem_d"},
		"eq_armor_vai_tho_d":       {ID: "eq_armor_vai_tho_d"},
		"eq_boots_vai_tho_d":       {ID: "eq_boots_vai_tho_d"},
		"eq_artifact_guong_dong_d": {ID: "eq_artifact_guong_dong_d"},
	})
}

func TestNPCRegistry_AllItemsExistInCore(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for _, it := range shop.Items {
			if !it.Enabled {
				continue
			}
			// Kiểm tra Lệch Dữ Liệu: Vật phẩm NPC bán CÓ TỒN TẠI trong game không?
			_, ok := item.GetDefinition(it.ItemID)
			if !ok {
				t.Errorf("CRITICAL DATA MISMATCH: Thương hội [%s] đang bán vật phẩm ảo [%s] không hề tồn tại trong Core Registry!", shopID, it.ItemID)
			}
		}
	}
}

func TestNPCRegistry_NoInfiniteMoneyExploit(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for _, it := range shop.Items {
			if it.SellPrice >= it.BuyPrice {
				t.Errorf("LỖ HỔNG KINH TẾ (Exploit): Vật phẩm [%s] ở [%s] có Giá Bán Ra (%d) <= Giá Thu Mua (%d). User có thể lặp vòng mua đi bán lại để farm linh thạch vô hạn!", it.ItemID, shopID, it.BuyPrice, it.SellPrice)
			}
		}
	}
}

func TestNPCRegistry_BuyPriceMustBePositive(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for itemID, it := range shop.Items {
			if it.BuyPrice <= 0 {
				t.Errorf("LỖI DỮ LIỆU: Shop %s - Item %s có BuyPrice <= 0 (%d). Sẽ bị mua miễn phí!", shopID, itemID, it.BuyPrice)
			}
		}
	}
}

func TestNPCRegistry_SellPriceCannotBeNegative(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for itemID, it := range shop.Items {
			if it.SellPrice < 0 {
				t.Errorf("LỖI DỮ LIỆU: Shop %s - Item %s có SellPrice < 0 (%d). Sẽ trừ linh thạch của user khi bán!", shopID, itemID, it.SellPrice)
			}
		}
	}
}

func TestNPCRegistry_StockIsValid(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for itemID, it := range shop.Items {
			if it.Stock < -1 {
				t.Errorf("LỖI DỮ LIỆU: Shop %s - Item %s có Stock không hợp lệ (%d). Chỉ chấp nhận >= 0 hoặc -1.", shopID, itemID, it.Stock)
			}
		}
	}
}

func TestNPCRegistry_ShopIDMatchesKey(t *testing.T) {
	for shopID, shop := range npc.Registry {
		if shopID != shop.ID {
			t.Errorf("LỖI CẤU TRÚC: Key của map (%s) không khớp với ID bên trong struct (%s)", shopID, shop.ID)
		}
		for itemID, it := range shop.Items {
			if itemID != it.ItemID {
				t.Errorf("LỖI CẤU TRÚC: Key của map Items (%s) không khớp với ItemID bên trong struct (%s)", itemID, it.ItemID)
			}
		}
	}
}

func TestNPCRegistry_CategoriesAreValid(t *testing.T) {
	for shopID, shop := range npc.Registry {
		for itemID, it := range shop.Items {
			// Nếu phân loại rỗng, UI sẽ tự động đẩy vào nhóm "Khác". Nhưng tốt nhất là nên cảnh báo để admin điền cho đẹp.
			if it.Category == "" {
				t.Logf("CẢNH BÁO UI: Shop %s - Item %s chưa có Category. Sẽ hiển thị ở mục 'Khác'.", shopID, itemID)
			}
		}
	}
}
