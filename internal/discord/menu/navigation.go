// File: internal/discord/menu/navigation.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa hằng Page, map trang cha-con, và helpers điều hướng menu.
// Ghi chú: Thêm Page mới ở đây khi xây dựng tính năng theo version roadmap.

package menu

// Page định danh trang hiện tại trong hệ thống menu.
type Page string

const (
	PageMain        Page = "main"        // Trang chủ — Main Menu
	PageProfile     Page = "profile"     // Hồ sơ người chơi
	PageCultivation Page = "cultivation" // Tu luyện
	PageInventory   Page = "inventory"   // Túi đồ / Trang bị — TODO v0.3
	PageSkills      Page = "skills"      // Kỹ năng / Công pháp — TODO v0.4
	PagePets        Page = "pets"        // Linh thú / Con rối — TODO v0.6
	PageGacha       Page = "gacha"       // Cơ duyên / Gacha — TODO v0.5
	PageMarket      Page = "market"      // Chợ / Đấu giá — TODO v0.8
	PageSect        Page = "sect"        // Tông môn / NPC / Đạo lữ — TODO v1.0
)

// parentPage map trang hiện tại → trang cha (để nút Quay lại hoạt động đúng).
var parentPage = map[Page]Page{
	PageProfile:     PageMain,
	PageCultivation: PageMain,
	PageInventory:   PageMain,
	PageSkills:      PageMain,
	PagePets:        PageMain,
	PageGacha:       PageMain,
	PageMarket:      PageMain,
	PageSect:        PageMain,
}

// ParentOf trả về trang cha của trang đã cho.
// Nếu không có trang cha (ví dụ: PageMain), trả về chuỗi rỗng.
func ParentOf(page Page) string {
	if parent, ok := parentPage[page]; ok {
		return string(parent)
	}
	return "" // Không có trang cha — nút Quay lại bị disable
}

// IsValidPage kiểm tra xem page có phải là trang hợp lệ không.
func IsValidPage(page Page) bool {
	switch page {
	case PageMain, PageProfile, PageCultivation,
		PageInventory, PageSkills, PagePets,
		PageGacha, PageMarket, PageSect:
		return true
	}
	return false
}
