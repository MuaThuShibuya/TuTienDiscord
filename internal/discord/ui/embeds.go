// File: internal/discord/ui/embeds.go
// Version: v0.1
// Purpose: Builder functions for standard Discord message embeds used across the bot.
// Security: Never put raw user input directly into embed fields without sanitization.
// Notes: All embeds use the color palette from colors.go and emojis from emojis.go.

package ui

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// ErrorEmbed builds a standard error embed (ephemeral, shown only to the user).
func ErrorEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: EmojiError.String() + " " + message,
		Color:       ColorError,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessEmbed builds a standard success embed.
func SuccessEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       EmojiSuccess.String() + " " + title,
		Description: message,
		Color:       ColorSuccess,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// InfoEmbed builds a standard informational embed.
func InfoEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       EmojiInfo.String() + " " + title,
		Description: message,
		Color:       ColorInfo,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// WarningEmbed builds a warning embed.
func WarningEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: EmojiWarning.String() + " " + message,
		Color:       ColorWarning,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// EphemeralError returns an interaction response with an error embed, visible only to the user.
func EphemeralError(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{ErrorEmbed(message)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// EphemeralUpdate edits an existing interaction's deferred response with an error embed.
func EphemeralUpdate(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{ErrorEmbed(message)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
