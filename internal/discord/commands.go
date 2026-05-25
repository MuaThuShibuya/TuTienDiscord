// File: internal/discord/commands.go
// Version: v0.1
// Purpose: Define all Discord slash command schemas registered with the Discord API.
// Security: Commands are registered server-side; descriptions must not reveal internal details.
// Notes: Add new commands here as new features are built. Re-register on every bot start.

package discord

import "github.com/bwmarrin/discordgo"

// AllCommands returns all slash command definitions to register with Discord.
func AllCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "start",
			Description: "Bắt đầu hành trình tu tiên của đạo hữu.",
		},
		{
			Name:        "menu",
			Description: "Mở giao diện chính của Vạn Pháp Tiên Nghịch.",
		},
		// TODO v0.2: /cultivate — shortcut to tĩnh tu without opening full menu
		// TODO v0.5: /gacha  — shortcut to gacha
		// TODO v0.8: /market — shortcut to market
	}
}
