// File: internal/discord/menu/profile/embed.go
// Chức năng: Tạo Discord embed cho trang Hồ Sơ tu sĩ.
// Ghi chú: Chỉ nhận ViewModel đã xử lý, không gọi DB hay service.

package profilemenu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

// BuildMenuResponse tạo response Discord đầy đủ cho trang Hồ Sơ.
func BuildMenuResponse(vm *menu.ProfileMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildEmbed(vm)},
		Components: buildComponents(vm.SessionID),
	}
}

func buildEmbed(vm *menu.ProfileMenuVM) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       emoji.Profile.String() + " Hồ Sơ Tu Sĩ",
		Description: "**" + vm.DaoName + "** — Đạo hiệu đã được khắc vào thiên địa.",
		Color:       ui.ColorDefault,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Đạo Hiệu", Value: vm.DaoName, Inline: true},
			{Name: "Tên Discord", Value: vm.DisplayName, Inline: true},
			{Name: "Tham Gia", Value: vm.JoinedAt, Inline: true},
			{Name: "Lần Cuối Online", Value: vm.LastActive, Inline: true},
			{Name: emoji.SpiritStone.String() + " Linh Thạch", Value: vm.SpiritStones, Inline: true},
			{Name: emoji.SpiritJade.String() + " Linh Ngọc", Value: vm.SpiritJades, Inline: true},
			{Name: emoji.FateTicket.String() + " Vé Cơ Duyên", Value: vm.FateTickets, Inline: true},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Hồ Sơ"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
