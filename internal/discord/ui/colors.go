// File: internal/discord/ui/colors.go
// Version: v0.1
// Purpose: Discord embed color palette for the Tu Tien bot.
// Notes: Colors are standard decimal integers (Discord API format).

package ui

// Embed color constants. Use these in discordgo.MessageEmbed.Color.
const (
	ColorDefault   = 0x8B5CF6 // Purple — main theme
	ColorSuccess   = 0x10B981 // Green — success, reward
	ColorWarning   = 0xF59E0B // Amber — warning, cooldown
	ColorError     = 0xEF4444 // Red — error, failure
	ColorInfo      = 0x3B82F6 // Blue — info, tips
	ColorCultivate = 0x7C3AED // Deep purple — cultivation
	ColorCombat    = 0xDC2626 // Dark red — combat
	ColorGacha     = 0xD97706 // Gold — gacha
	ColorEconomy   = 0x059669 // Teal — economy
	ColorNeutral   = 0x6B7280 // Gray — neutral/closed
)
