package shopmenu

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildNPCShopComponents(sessionID string, vm NPCShopViewModel) []discordgo.MessageComponent {
	var comps []discordgo.MessageComponent
	var catOptions []discordgo.SelectMenuOption
	for _, cat := range vm.Categories {
		catOptions = append(catOptions, ui.SelectOption(cat, cat, "Xem danh mục "+cat, nil, cat == vm.SelectedCategory))
	}
	if len(catOptions) > 0 {
		comps = append(comps, ui.ActionRow(ui.SelectMenu(menu.Build(menu.DomainShop, menu.ActionShopNPCCategory, sessionID), "Chuyển danh mục khác...", catOptions)))
	}

	var itemOptions []discordgo.SelectMenuOption
	for _, item := range vm.Items {
		label := fmt.Sprintf("Mua với giá %d LT", item.Price)
		if vm.Mode == "sell" {
			label = fmt.Sprintf("Bán 1x lấy %d LT", item.Price)
		}
		itemOptions = append(itemOptions, ui.SelectOption(item.Name, item.ID, label, emoji.Reward, false))
		if len(itemOptions) >= 25 {
			break
		}
	}

	if len(itemOptions) > 0 {
		actionID, placeholder := menu.ActionShopNPCBuy, "Chọn vật phẩm để mua (1 cái)..."
		if vm.Mode == "sell" {
			actionID, placeholder = menu.ActionShopNPCSell, "Chọn vật phẩm để bán (1 cái)..."
		}
		comps = append(comps, ui.ActionRow(ui.SelectMenu(menu.Build(menu.DomainShop, actionID, sessionID), placeholder, itemOptions)))
	}

	// Hàng Phân Trang
	if vm.TotalPages > 1 {
		prevPage := vm.CurrentPage - 1
		if prevPage < 1 {
			prevPage = vm.TotalPages
		}
		nextPage := vm.CurrentPage + 1
		if nextPage > vm.TotalPages {
			nextPage = 1
		}

		comps = append(comps, ui.ActionRow(
			ui.Button(fmt.Sprintf("Trang %d/%d", vm.CurrentPage, vm.TotalPages), "dummy", ui.BtnSecondary, nil, true),
			ui.Button("◀ Trước", menu.Build(menu.DomainShop, menu.ActionShopNPCPage, sessionID, fmt.Sprintf("%d", prevPage)), ui.BtnPrimary, nil, false),
			ui.Button("Sau ▶", menu.Build(menu.DomainShop, menu.ActionShopNPCPage, sessionID, fmt.Sprintf("%d", nextPage)), ui.BtnPrimary, nil, false),
		))
	}

	modeBtn := ui.Button("Chuyển Sang Bán Hàng", menu.Build(menu.DomainShop, menu.ActionShopNPCModeSell, sessionID), ui.BtnSecondary, emoji.Reward, false)
	if vm.Mode == "sell" {
		modeBtn = ui.Button("Chuyển Sang Mua Hàng", menu.Build(menu.DomainShop, menu.ActionShopNPCModeBuy, sessionID), ui.BtnPrimary, emoji.Reward, false)
	}

	comps = append(comps, ui.ActionRow(modeBtn, ui.Button("Làm Mới", menu.Build(menu.DomainShop, menu.ActionShopNPCRefresh, sessionID), ui.BtnSecondary, emoji.Refresh, false)))
	comps = append(comps, ui.NavRow(sessionID, string(menu.PageMarket), string(menu.PageMain)))

	return comps
}
