// File: internal/discord/menu/profile_menu.go
// Phiên bản: v0.1.1
// Mục đích: UI Builder cho trang Hồ Sơ — render embed thông tin tu sĩ và các nút thao tác.
// Bảo mật: Chỉ nhận ViewModel đã xử lý sẵn, không gọi DB hay service, không có side effect.
// Ghi chú: Handler fetch dữ liệu → map sang ProfileMenuVM → gọi hàm này → gửi về Discord.

package menu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
)

// BuildProfileMenuResponse tạo response Discord cho trang Hồ Sơ.
func BuildProfileMenuResponse(vm *ProfileMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildProfileEmbed(vm)},
		Components: buildProfileComponents(vm.SessionID),
	}
}

func buildProfileEmbed(vm *ProfileMenuVM) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       ui.EmojiProfile.String() + " Hồ Sơ Tu Sĩ",
		Description: "**" + vm.DaoName + "** — Đạo hiệu đã được khắc vào thiên địa.",
		Color:       ui.ColorDefault,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Đạo Hiệu", Value: vm.DaoName, Inline: true},
			{Name: "Tên Discord", Value: vm.DisplayName, Inline: true},
			{Name: "Tham Gia", Value: vm.JoinedAt, Inline: true},
			{Name: "Lần Cuối Online", Value: vm.LastActive, Inline: true},
			{
				Name:   ui.EmojiSpiritStone.String() + " Linh Thạch",
				Value:  vm.SpiritStones,
				Inline: true,
			},
			{
				Name:   ui.EmojiSpiritJade.String() + " Linh Ngọc",
				Value:  vm.SpiritJades,
				Inline: true,
			},
			{
				Name:   ui.EmojiFateTicket.String() + " Vé Cơ Duyên",
				Value:  vm.FateTickets,
				Inline: true,
			},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Hồ Sơ"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildProfileComponents(sessionID string) []discordgo.MessageComponent {
	actionRow := ui.ActionRow(
		ui.Button("Đổi Đạo Hiệu", Build(DomainProfile, ActionRename, sessionID), ui.BtnPrimary, &ui.EmojiProfile, false),
	)
	navRow := ui.NavRow(sessionID, string(PageProfile), string(PageMain))

	return []discordgo.MessageComponent{actionRow, navRow}
}
