// File: internal/discord/menu/profile/components.go
// Chức năng: Tạo Discord components (button, select menu) cho trang Hồ Sơ.

package profilemenu

import (
	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildComponents(sessionID string) []discordgo.MessageComponent {
	actionRow := ui.ActionRow(
		ui.Button("Đổi Đạo Hiệu",
			menu.Build(menu.DomainProfile, menu.ActionRename, sessionID),
			ui.BtnPrimary, emoji.Profile, false),
	)
	navRow := ui.NavRow(sessionID, string(menu.PageProfile), string(menu.PageMain))
	return []discordgo.MessageComponent{actionRow, navRow}
}
