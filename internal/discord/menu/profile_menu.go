// File: internal/discord/menu/profile_menu.go
// Version: v0.1
// Purpose: UI builder for the Profile page — renders player info embed and buttons.
// Security: Pure rendering only. No DB calls. Receives pre-fetched data from handler.
// Notes: Handler fetches data → passes to BuildProfileMenuResponse → sends to Discord.

package menu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	profile "github.com/yourname/tu-tien-bot/internal/game/profile"
	economy "github.com/yourname/tu-tien-bot/internal/game/economy"
	"github.com/yourname/tu-tien-bot/pkg/utils"
)

// ProfileMenuData bundles all data needed to render the profile page.
type ProfileMenuData struct {
	Session *Session
	Player  *profile.Player
	Wallet  *economy.Wallet
}

// BuildProfileMenuResponse constructs the interaction response data for the profile page.
func BuildProfileMenuResponse(data *ProfileMenuData) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildProfileEmbed(data)},
		Components: buildProfileComponents(data.Session.SessionID),
	}
}

func buildProfileEmbed(data *ProfileMenuData) *discordgo.MessageEmbed {
	p := data.Player
	w := data.Wallet

	joinedAt := utils.DiscordTimestamp(p.CreatedAt, "D")
	lastActive := utils.DiscordTimestamp(p.LastActiveAt, "R")

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiProfile.String() + " Hồ Sơ Tu Sĩ",
		Description: fmt.Sprintf("**%s** — Đạo hiệu đã được khắc vào thiên địa.", p.DaoName),
		Color:       ui.ColorDefault,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Đạo Hiệu", Value: p.DaoName, Inline: true},
			{Name: "Tên Discord", Value: p.DisplayName, Inline: true},
			{Name: "Tham Gia", Value: joinedAt, Inline: true},
			{Name: "Lần Cuối Online", Value: lastActive, Inline: true},
			{
				Name: ui.EmojiSpiritStone.String() + " Linh Thạch",
				Value: utils.FormatNumber(w.SpiritStones),
				Inline: true,
			},
			{
				Name: ui.EmojiSpiritJade.String() + " Linh Ngọc",
				Value: utils.FormatNumber(w.SpiritJades),
				Inline: true,
			},
			{
				Name: ui.EmojiFateTicket.String() + " Vé Cơ Duyên",
				Value: fmt.Sprintf("%d vé", w.FateTickets),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Hồ Sơ"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildProfileComponents(sessionID string) []discordgo.MessageComponent {
	actionRow := ui.ActionRow(
		ui.Button("Đổi Đạo Hiệu", "profile:rename:"+sessionID, ui.BtnPrimary, &ui.EmojiProfile, false),
		// TODO v0.2+: Button("Xem Thành Tích", ...)
		// TODO v0.9+: Button("Chia Sẻ Hồ Sơ", ...)
	)
	navRow := ui.NavRow(sessionID, string(PageProfile), string(PageMain))

	return []discordgo.MessageComponent{actionRow, navRow}
}
