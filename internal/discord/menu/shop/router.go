package shopmenu

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	npcgame "github.com/whiskey/tu-tien-bot/internal/game/shop/npc"
	playergame "github.com/whiskey/tu-tien-bot/internal/game/shop/player"
)

type Router struct {
	npcSvc       npcgame.Service
	playerSvc    playergame.Service
	economySvc   economy.Service
	inventorySvc inventory.Service
	sessionSvc   menu.SessionService
	log          *zap.Logger
}

func NewRouter(npcSvc npcgame.Service, playerSvc playergame.Service, economySvc economy.Service, invSvc inventory.Service, sessionSvc menu.SessionService, log *zap.Logger) *Router {
	return &Router{npcSvc: npcSvc, playerSvc: playerSvc, economySvc: economySvc, inventorySvc: invSvc, sessionSvc: sessionSvc, log: log.Named("menu.shop")}
}

func getModalValue(i *discordgo.Interaction, customID string) string {
	for _, comp := range i.ModalSubmitData().Components {
		if row, ok := comp.(*discordgo.ActionsRow); ok {
			for _, inner := range row.Components {
				if textInput, ok := inner.(*discordgo.TextInput); ok && textInput.CustomID == customID {
					return textInput.Value
				}
			}
		}
	}
	return ""
}

// sendPopup tạo thông báo nổi phía trên Menu chính (không ghi đè giao diện cửa hàng)
func sendPopup(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	_, err := s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		zap.L().Error("Lỗi gửi Followup popup", zap.Error(err))
	}
}

func (r *Router) RenderMarketMain(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	embed := &discordgo.MessageEmbed{
		Title:       "⚖️ Phường Thị Giao Thương",
		Description: "Nơi giao lưu, trao đổi kỳ trân dị bảo của giới tu chân.",
		Color:       ui.ColorEconomy,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Vạn Bảo Các", Value: "Thương hội chuyên bán đan dược, nguyên liệu và thu mua phế phẩm.", Inline: true},
			{Name: "Thiên Bảo Đấu Giá Các", Value: "Sàn đấu giá tự do giữa các vị đạo hữu. Nơi pháp bảo được định giá bằng linh thạch.", Inline: true},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	comps := []discordgo.MessageComponent{
		ui.ActionRow(
			ui.Button("Vạn Bảo Các", menu.Build(menu.DomainShop, menu.ActionShopGoNPC, session.SessionID), ui.BtnPrimary, nil, false),
			ui.Button("Sàn Đấu Giá", menu.Build(menu.DomainShop, menu.ActionShopGoPlayer, session.SessionID), ui.BtnSuccess, nil, false),
		),
		ui.NavRow(session.SessionID, string(menu.PageMarket), string(menu.PageMain)),
	}
	return &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Components: comps}, nil
}

func (r *Router) RenderPlayerShop(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	wallet, err := r.economySvc.GetWallet(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, err
	}
	listings, err := r.playerSvc.GetActiveListings(ctx, session.GuildID, 25, 0)
	if err != nil {
		return nil, err
	}
	vm := AuctionViewModel{SpiritStones: wallet.SpiritStones}
	for _, l := range listings {
		name := l.ItemDefID
		if def, ok := item.GetDefinition(l.ItemDefID); ok {
			name = def.Name
		}
		vm.Listings = append(vm.Listings, AuctionListingVM{
			ListingID:  l.ID.Hex(),
			SellerName: fmt.Sprintf("<@%s>", l.SellerID), // Tự động ping để Discord hiển thị Tên thay vì ID
			ItemName:   name,
			Quantity:   int(l.Quantity),
			TotalPrice: l.TotalPrice,
			ExpiresAt:  l.ExpiresAt,
		})
	}
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{BuildAuctionHouseEmbed(vm)},
		Components: BuildAuctionHouseComponents(session.SessionID, vm),
	}, nil
}

func (r *Router) RenderNPCShop(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	shopDef, err := npcgame.GetShop("van_bao_cac")
	if err != nil {
		return nil, err
	}

	wallet, err := r.economySvc.GetWallet(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, err
	}

	// Phân tích trạng thái "mode|category|page" (VD: "sell|Đan Dược|1")
	state := session.CurrentCategory
	parts := strings.Split(state, "|")
	mode := "buy"
	selectedCat := ""
	page := 1

	if len(parts) > 0 && parts[0] != "" {
		mode = parts[0]
	}
	if len(parts) > 1 {
		selectedCat = parts[1]
	}
	if len(parts) > 2 {
		if p, err := strconv.Atoi(parts[2]); err == nil && p > 0 {
			page = p
		}
	}

	categoryMap := make(map[string]bool)
	for _, it := range shopDef.Items {
		cat := it.Category
		if cat == "" {
			cat = "Khác"
		}
		categoryMap[cat] = true
	}
	var categories []string
	for k := range categoryMap {
		categories = append(categories, k)
	}
	sort.Strings(categories)

	if selectedCat == "" && len(categories) > 0 {
		selectedCat = categories[0]
	}

	vm := NPCShopViewModel{
		ShopName:         shopDef.Name,
		NPCName:          "Mạc Chưởng Quỹ",
		SpiritStones:     wallet.SpiritStones,
		Mode:             mode,
		SelectedCategory: selectedCat,
		Categories:       categories,
		CurrentPage:      page,
	}

	var rawItems []NPCShopItemVM

	// Thu thập dữ liệu theo Mode
	if mode == "sell" {
		_, items, err := r.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
		if err == nil {
			for _, it := range items {
				if it.Quantity <= 0 {
					continue
				}
				shopItem, ok := shopDef.Items[it.DefinitionID]
				if !ok || !shopItem.Enabled || shopItem.SellPrice <= 0 {
					continue
				}
				cat := shopItem.Category
				if cat == "" {
					cat = "Khác"
				}
				if cat != selectedCat {
					continue
				}
				name := it.DefinitionID
				if def, ok := item.GetDefinition(it.DefinitionID); ok {
					name = def.Name
				}
				rawItems = append(rawItems, NPCShopItemVM{ID: it.DefinitionID, Name: name, Price: shopItem.SellPrice, Quantity: it.Quantity})
			}
		}
	} else {
		for _, it := range shopDef.Items {
			if !it.Enabled {
				continue
			}
			cat := it.Category
			if cat == "" {
				cat = "Khác"
			}
			if cat != selectedCat {
				continue
			}
			name := it.ItemID
			itemDesc := "Vật phẩm thần bí"
			if def, ok := item.GetDefinition(it.ItemID); ok {
				name = def.Name
				itemDesc = string(def.Rarity) + " - " + def.Description
			}
			rawItems = append(rawItems, NPCShopItemVM{ID: it.ItemID, Name: name, Price: it.BuyPrice, Stock: it.Stock, Desc: itemDesc})
		}
	}

	// Thuật toán Phân trang (20 mục/trang)
	itemsPerPage := 20
	totalItems := len(rawItems)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if end > totalItems {
		end = totalItems
	}

	vm.Items = rawItems[start:end]
	vm.CurrentPage = page
	vm.TotalPages = totalPages

	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildNPCShopEmbed(vm)},
		Components: buildNPCShopComponents(session.SessionID, vm),
	}, nil
}

func (r *Router) HandleShopInteraction(s *discordgo.Session, i *discordgo.Interaction, session *menu.Session, action, extra string) {
	ctx := context.Background()
	switch action {
	case menu.ActionShopGoNPC:
		session.CurrentCategory = ""
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, "")
		res, err := r.RenderNPCShop(ctx, session)
		if err == nil {
			_, err = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
			if err != nil {
				r.log.Error("Lỗi Edit NPCShop", zap.Error(err))
			}
		}
	case menu.ActionShopGoPlayer:
		session.CurrentCategory = ""
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, "")
		res, err := r.RenderPlayerShop(ctx, session)
		if err == nil {
			_, err = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
			if err != nil {
				r.log.Error("Lỗi Edit PlayerShop", zap.Error(err))
			}
		} else {
			sendPopup(s, i, ui.ErrorEmbed(err.Error()))
		}
	case menu.ActionShopNPCModeBuy:
		state := session.CurrentCategory
		cat := ""
		if strings.Contains(state, "|") {
			cat = strings.SplitN(state, "|", 2)[1]
		}
		session.CurrentCategory = "buy|" + cat
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, session.CurrentCategory)
		res, err := r.RenderNPCShop(ctx, session)
		if err == nil {
			_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
		}
	case menu.ActionShopNPCModeSell:
		state := session.CurrentCategory
		cat := ""
		if strings.Contains(state, "|") {
			cat = strings.SplitN(state, "|", 2)[1]
		}
		session.CurrentCategory = "sell|" + cat
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, session.CurrentCategory)
		res, err := r.RenderNPCShop(ctx, session)
		if err == nil {
			_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
		}
	case menu.ActionShopNPCCategory:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			state := session.CurrentCategory
			if !strings.Contains(state, "|") {
				state = "buy|"
			}
			mode := strings.SplitN(state, "|", 2)[0]
			session.CurrentCategory = mode + "|" + data.Values[0]
			_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, session.CurrentCategory)
			res, err := r.RenderNPCShop(ctx, session)
			if err == nil {
				_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
			}
		}
	case menu.ActionShopNPCBuy:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			itemID := data.Values[0]
			walletAfter, err := r.npcSvc.BuyItem(ctx, session.UserID, session.GuildID, "van_bao_cac", itemID, 1)
			if err != nil {
				sendPopup(s, i, ui.ErrorEmbed(err.Error()))
				return
			}
			itemName := itemID
			if def, ok := item.GetDefinition(itemID); ok {
				itemName = def.Name
			}
			sendPopup(s, i, ui.SuccessEmbed("Giao Dịch Thành Công", fmt.Sprintf("Đạo hữu đã mua **1x %s** cất vào túi.\nSố dư còn lại: **%d** linh thạch.", itemName, walletAfter.SpiritStones)))

			// Làm mới lại shop để hiển thị linh thạch mới
			res, err := r.RenderNPCShop(ctx, session)
			if err == nil {
				_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
			}
		}
	case menu.ActionShopNPCSell:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			itemID := data.Values[0]
			walletAfter, err := r.npcSvc.SellItem(ctx, session.UserID, session.GuildID, "van_bao_cac", itemID, 1)
			if err != nil {
				sendPopup(s, i, ui.ErrorEmbed(err.Error()))
				return
			}
			itemName := itemID
			if def, ok := item.GetDefinition(itemID); ok {
				itemName = def.Name
			}
			sendPopup(s, i, ui.SuccessEmbed("Bán Thành Công", fmt.Sprintf("Đạo hữu đã giao **1x %s** cho cửa hàng.\nSố dư hiện tại: **%d** linh thạch.", itemName, walletAfter.SpiritStones)))
		}
		res, err := r.RenderNPCShop(ctx, session)
		if err == nil {
			_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
		}
	case menu.ActionShopNPCRefresh:
		session.CurrentCategory = "" // Đặt lại về giao diện mua hàng
		res, err := r.RenderNPCShop(ctx, session)
		if err == nil {
			_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
		}
	case menu.ActionShopPlayerBuy:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			listingID := data.Values[0]
			err := r.playerSvc.PurchaseListing(ctx, session.UserID, session.GuildID, listingID)
			if err != nil {
				sendPopup(s, i, ui.ErrorEmbed(err.Error()))
				// Vẫn render lại sàn đấu giá để xóa đồ đã bị người khác mua mất
				res, err := r.RenderPlayerShop(ctx, session)
				if err == nil {
					_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
				}
				return
			}
			sendPopup(s, i, ui.SuccessEmbed("Mua Thành Công", "Đấu giá thành công, vật phẩm đã cất vào túi!"))
			res, err := r.RenderPlayerShop(ctx, session)
			if err == nil {
				_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
			}
		}
	case menu.ActionShopPlayerList:
		_, items, err := r.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
		if err != nil {
			sendPopup(s, i, ui.ErrorEmbed("Không thể lấy dữ liệu túi đồ."))
			return
		}
		var opts []discordgo.SelectMenuOption
		for _, it := range items {
			if it.Quantity <= 0 {
				continue
			}
			name := it.DefinitionID
			if def, ok := item.GetDefinition(it.DefinitionID); ok {
				name = def.Name
			}
			opts = append(opts, ui.SelectOption(name, it.DefinitionID, fmt.Sprintf("Số lượng có: %d", it.Quantity), nil, false))
			if len(opts) >= 25 {
				break
			}
		}
		if len(opts) == 0 {
			sendPopup(s, i, ui.ErrorEmbed("Đạo hữu không có vật phẩm nào để bán."))
			return
		}
		embed := &discordgo.MessageEmbed{Title: "Chọn Vật Phẩm Đăng Bán", Description: "Hãy chọn 1 vật phẩm từ túi đồ để đưa lên Sàn Đấu Giá.", Color: ui.ColorEconomy}
		comps := []discordgo.MessageComponent{
			ui.ActionRow(ui.SelectMenu(menu.Build(menu.DomainShop, menu.ActionShopPlayerListSelect_modal, session.SessionID), "Chọn vật phẩm...", opts)),
			ui.ActionRow(ui.Button("Hủy Bỏ", menu.Build(menu.DomainShop, menu.ActionShopGoPlayer, session.SessionID), ui.BtnSecondary, nil, false)),
		}
		_, err = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
		if err != nil {
			r.log.Error("Lỗi Edit p_list", zap.Error(err))
		}
	case menu.ActionShopPlayerListSelect_modal:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			return
		}
		itemID := data.Values[0]
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: menu.Build(menu.DomainShop, menu.ActionShopPlayerListApply, session.SessionID, itemID),
				Title:    "Đăng Bán",
				Components: []discordgo.MessageComponent{
					ui.ActionRow(discordgo.TextInput{CustomID: "quantity", Label: "Số lượng bán", Style: discordgo.TextInputShort, Placeholder: "Ví dụ: 1", Required: true}),
					ui.ActionRow(discordgo.TextInput{CustomID: "price", Label: "Tổng Giá (Linh Thạch)", Style: discordgo.TextInputShort, Placeholder: "Ví dụ: 500", Required: true}),
				},
			},
		})
	case menu.ActionShopPlayerListApply:
		itemID := extra
		qty, _ := strconv.ParseInt(getModalValue(i, "quantity"), 10, 64)
		price, _ := strconv.ParseInt(getModalValue(i, "price"), 10, 64)
		if qty <= 0 || price <= 0 {
			sendPopup(s, i, ui.WarningEmbed("Số lượng và giá bán phải lớn hơn 0."))
			return
		}
		if err := r.playerSvc.CreateListing(ctx, session.UserID, session.GuildID, itemID, qty, price); err != nil {
			sendPopup(s, i, ui.WarningEmbed(err.Error()))
			return
		}
		sendPopup(s, i, ui.SuccessEmbed("Thành Công", "Vật phẩm đã được đưa lên Thiên Bảo Các."))
		res, err := r.RenderPlayerShop(ctx, session)
		if err == nil {
			_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &res.Embeds, Components: &res.Components})
		}
	default:
		sendPopup(s, i, ui.WarningEmbed("Tính năng cửa hàng đang tải lên linh mạch, vui lòng thử lại sau!"))
	}
}
