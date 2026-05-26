// File: internal/discord/menu/navigation.go
// Chức năng: Định nghĩa hằng Page, map trang cha-con, và helpers điều hướng menu.
// Ghi chú: Thêm Page mới ở đây khi xây dựng tính năng. parentPage dùng bởi nút Quay Lại.

package menu

// Page định danh trang hiện tại trong hệ thống menu.
type Page string

const (
	PageMain        Page = "main"        // Trang chủ
	PageProfile     Page = "profile"     // Hồ sơ người chơi
	PageCultivation Page = "cultivation" // Tu luyện
	PageInventory   Page = "inventory"   // Túi đồ
	PageEquipment   Page = "equipment"   // Trang bị (v0.3)
	PageAlchemy     Page = "alchemy"     // Luyện đan — TODO v0.4
	PageCombat      Page = "combat"      // Chiến đấu / PvP — TODO v0.5
	PageSkills      Page = "skills"      // Kỹ năng / Công pháp — TODO v0.4
	PagePets        Page = "pets"        // Linh thú — TODO v0.6
	PageGacha       Page = "gacha"       // Cơ duyên / Gacha — TODO v0.5
	PageMarket      Page = "market"      // Chợ / Đấu giá — TODO v0.8
	PageSect        Page = "sect"        // Tông môn — TODO v1.0
)

// parentPage ánh xạ trang hiện tại → trang cha (để nút Quay Lại hoạt động đúng).
var parentPage = map[Page]Page{
	PageProfile:     PageMain,
	PageCultivation: PageMain,
	PageInventory:   PageMain,
	PageEquipment:   PageMain,
	PageAlchemy:     PageMain,
	PageCombat:      PageMain,
	PageSkills:      PageMain,
	PagePets:        PageMain,
	PageGacha:       PageMain,
	PageMarket:      PageMain,
	PageSect:        PageMain,
}

// ParentOf trả về trang cha của trang đã cho.
func ParentOf(page Page) string {
	if parent, ok := parentPage[page]; ok {
		return string(parent)
	}
	return ""
}

// IsValidPage kiểm tra xem page có phải là trang hợp lệ không.
func IsValidPage(page Page) bool {
	if page == PageMain {
		return true
	}
	_, ok := parentPage[page]
	return ok
}
