package shopmenu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

type AuctionViewModel struct {
	SpiritStones int64
	Listings     []AuctionListingVM
}

type AuctionListingVM struct {
	ListingID  string
	SellerName string
	ItemName   string
	Quantity   int
	TotalPrice int64
	ExpiresAt  time.Time
}

func BuildAuctionHouseEmbed(vm AuctionViewModel) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{Title: "Thiên Bảo Đấu Giá Các", Description: fmt.Sprintf("💰 Linh thạch hiện có: **%s**\n\n", utils.FormatNumber(vm.SpiritStones)), Color: ui.ColorGacha, Timestamp: time.Now().UTC().Format(time.RFC3339)}
	for _, l := range vm.Listings {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: fmt.Sprintf("%s x%d", l.ItemName, l.Quantity), Value: fmt.Sprintf("Giá: %s LT\nNgười bán: %s\nHết hạn: <t:%d:R>", utils.FormatNumber(l.TotalPrice), l.SellerName, l.ExpiresAt.Unix()), Inline: false})
	}
	if len(vm.Listings) == 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Trống Trải", Value: "*Hiện tại không có đạo hữu nào đăng bán.*", Inline: false})
	}
	return embed
}
