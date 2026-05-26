// File: internal/discord/menu/inventory/components.go
// Chức năng: Tạo Discord components (select menu sử dụng đan dược, phân trang) cho Túi Đồ.

package invmenu

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildComponents(vm *menu.InventoryMenuVM) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	// Hàng 1: Select menu sử dụng đan dược (UsableItems = toàn bộ, không giới hạn theo trang)
	if len(vm.UsableItems) > 0 {
		var usableOptions []discordgo.SelectMenuOption
		for _, it := range vm.UsableItems {
			if len(usableOptions) >= 25 {
				break
			}
			usableOptions = append(usableOptions, ui.SelectOption(
				it.Name,
				it.InstanceID,
				fmt.Sprintf("Số lượng: %d | Cấp: %s", it.Quantity, it.Rarity),
				emoji.Bag,
				false,
			))
		}
		if len(usableOptions) > 0 {
			useSelect := ui.SelectMenu(
				menu.Build(menu.DomainInventory, menu.ActionInventoryUse, vm.SessionID),
				"💊 Sử dụng đan dược...",
				usableOptions,
			)
			components = append(components, ui.ActionRow(useSelect))
		}
	}

	// Hàng 2: Nút phân trang
	prevPage := vm.CurrentPage - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := vm.CurrentPage + 1
	if nextPage > vm.TotalPages {
		nextPage = vm.TotalPages
	}

	pageInfoID := menu.Build(menu.DomainNav, "pageinfo", vm.SessionID)
	pageRow := ui.ActionRow(
		ui.Button("◀ Trước",
			menu.Build(menu.DomainInventory, menu.ActionInventoryPage, vm.SessionID, fmt.Sprintf("%d", prevPage)),
			ui.BtnSecondary, nil, vm.CurrentPage <= 1),
		ui.Button(fmt.Sprintf("Trang %d/%d", vm.CurrentPage, vm.TotalPages),
			pageInfoID, ui.BtnSecondary, nil, true),
		ui.Button("Sau ▶",
			menu.Build(menu.DomainInventory, menu.ActionInventoryPage, vm.SessionID, fmt.Sprintf("%d", nextPage)),
			ui.BtnSecondary, nil, vm.CurrentPage >= vm.TotalPages),
	)
	components = append(components, pageRow)

	// Hàng cuối: Điều hướng
	components = append(components, ui.NavRow(vm.SessionID, string(menu.PageInventory), string(menu.PageMain)))
	return components
}
