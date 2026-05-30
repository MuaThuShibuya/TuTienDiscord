package shopmenu

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

type NPCShopViewModel struct {
	ShopName         string
	NPCName          string
	SpiritStones     int64
	Mode             string // "buy" hoặc "sell"
	SelectedCategory string
	Categories       []string
	Items            []NPCShopItemVM
	CurrentPage      int
	TotalPages       int
}

type NPCShopItemVM struct {
	ID       string
	Name     string
	Category string
	Price    int64
	Stock    int
	Desc     string
	Quantity int64 // Dùng cho màn hình Bán
}

func buildNPCShopEmbed(vm NPCShopViewModel) *discordgo.MessageEmbed {
	titleStr := "Mua Hàng"
	if vm.Mode == "sell" {
		titleStr = "Thu Mua"
	}
	desc := fmt.Sprintf("Mạc Chưởng Quỹ: \"Hoan nghênh đến **%s**.\"\n\n💰 Linh thạch hiện có: **%s**\n", vm.ShopName, utils.FormatNumber(vm.SpiritStones))

	embed := &discordgo.MessageEmbed{Title: fmt.Sprintf("%s - %s", vm.ShopName, titleStr), Description: desc, Color: ui.ColorEconomy, Timestamp: time.Now().UTC().Format(time.RFC3339)}

	var fieldsStr string
	for _, item := range vm.Items {
		if vm.Mode == "sell" {
			fieldsStr += fmt.Sprintf("▸ **%s** — Thu mua: %s LT *(Đang có: %d)*\n", item.Name, utils.FormatNumber(item.Price), item.Quantity)
		} else {
			stockStr := "Vô hạn"
			if item.Stock != -1 {
				stockStr = fmt.Sprintf("%d", item.Stock)
			}
			fieldsStr += fmt.Sprintf("▸ **%s** — %s LT *(Tồn: %s)*\n   └ *%s*\n", item.Name, utils.FormatNumber(item.Price), stockStr, item.Desc)
		}
	}
	if fieldsStr == "" {
		if vm.Mode == "sell" {
			fieldsStr = "*Không có vật phẩm nào trong mục này để bán.*"
		} else {
			fieldsStr = "*Mục này hiện đã bán hết.*"
		}
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: fmt.Sprintf("=== %s ===", strings.ToUpper(vm.SelectedCategory)), Value: fieldsStr})
	return embed
}
