// File: internal/discord/menu/main/components.go
// Chức năng: Xây dựng UI Component cho trang chính (Main Menu).
// Ghi chú: Gộp toàn bộ điều hướng vào 1 Select Menu duy nhất.
package mainmenu

import (
	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildComponents(vm *menu.MainMenuVM, isAdmin bool) []discordgo.MessageComponent {
	// Danh mục 1: Các tính năng cơ bản
	nav1 := ui.SelectMenu(
		menu.Build(menu.DomainMenuSelect, menu.ActionNavSelect, vm.SessionID),
		"Bế Quan Tu Luyện...",
		[]discordgo.SelectMenuOption{
			ui.SelectOption("Hồ Sơ", string(menu.PageProfile), "Xem thông tin cá nhân", emoji.Profile, false),
			ui.SelectOption("Tu Luyện", string(menu.PageCultivation), "Tĩnh tu, bế quan, đột phá", emoji.Cultivate, false),
			ui.SelectOption("Túi Đồ", string(menu.PageInventory), "Quản lý vật phẩm, đan dược", emoji.Bag, false),
			ui.SelectOption("Lò Đan", string(menu.PageAlchemy), "Luyện chế đan dược", emoji.Alchemy, false),
			ui.SelectOption("Trang Bị", string(menu.PageEquipment), "Mặc và tháo trang bị", emoji.Equip, false),
		},
	)

	// Danh mục 2: Các tính năng mở rộng
	nav2 := ui.SelectMenu(
		menu.Build(menu.DomainMenuSelect, menu.ActionNav2Select, vm.SessionID),
		"Hành Tẩu Giang Hồ...",
		[]discordgo.SelectMenuOption{
			ui.SelectOption("Du Ngoạn / Bí Cảnh", string(menu.PagePvE), "Khiêu chiến yêu ma, tìm kiếm cơ duyên.", emoji.Map, false),
			ui.SelectOption("Tông Môn", "sect", "Quản lý tông môn", emoji.Realm, false),
			ui.SelectOption("PvP", "pvp", "Chiến đấu với người chơi khác", emoji.CombatPower, false),
			ui.SelectOption("Nhiệm Vụ", "quest", "Nhiệm vụ hàng ngày", emoji.Skill, false),
			ui.SelectOption("Cửa Hàng", "shop", "Mua bán vật phẩm", emoji.SpiritStone, false),
		},
	)

	// Nút đóng ở hàng cuối
	navRow := ui.ActionRow(
		ui.Button("Đóng Menu", menu.Build(menu.DomainNav, menu.ActionClose, vm.SessionID), ui.BtnDanger, emoji.Close, false),
	)

	adminButton := ui.Button("Thiên Cơ Các", menu.Build(menu.DomainAdmin, menu.ActionAdminMain, vm.SessionID), ui.BtnSecondary, emoji.Admin, false)

	components := []discordgo.MessageComponent{
		ui.ActionRow(nav1),
		ui.ActionRow(nav2),
	}

	if isAdmin {
		components = append(components, ui.ActionRow(adminButton))
	}

	components = append(components, navRow)
	return components
}
