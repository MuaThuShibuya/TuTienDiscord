package shopmenu

import (
	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func BuildAuctionHouseComponents(sessionID string, vm AuctionViewModel) []discordgo.MessageComponent {
	var comps []discordgo.MessageComponent
	var opts []discordgo.SelectMenuOption
	for _, l := range vm.Listings {
		opts = append(opts, ui.SelectOption(l.ItemName, l.ListingID, "Mua vật phẩm này", emoji.Reward, false))
	}
	if len(opts) > 0 {
		comps = append(comps, ui.ActionRow(ui.SelectMenu(menu.Build(menu.DomainShop, menu.ActionShopPlayerBuy, sessionID), "Chọn vật phẩm đấu giá...", opts)))
	}

	comps = append(comps, ui.ActionRow(ui.Button("Đăng Bán", menu.Build(menu.DomainShop, menu.ActionShopPlayerList, sessionID), ui.BtnSuccess, nil, false), ui.Button("Quản Lý", menu.Build(menu.DomainShop, menu.ActionShopPlayerManage, sessionID), ui.BtnSecondary, emoji.Profile, false)))
	comps = append(comps, ui.NavRow(sessionID, string(menu.PageMarket), string(menu.PageMarket)))
	return comps
}
