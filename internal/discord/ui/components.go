// File: internal/discord/ui/components.go
// Version: v0.1
// Purpose: Builder functions for Discord UI components — buttons, select menus, action rows.
// Security: custom_id values must encode sessionId to allow ownership validation in handlers.
// Notes: custom_id format: "<action>:<sessionId>[:<extra>]" — parsed in menu router.

package ui

import "github.com/bwmarrin/discordgo"

// ButtonStyle aliases for readability.
const (
	BtnPrimary   = discordgo.PrimaryButton
	BtnSecondary = discordgo.SecondaryButton
	BtnSuccess   = discordgo.SuccessButton
	BtnDanger    = discordgo.DangerButton
	BtnLink      = discordgo.LinkButton
)

// ActionRow wraps components into a single action row.
func ActionRow(components ...discordgo.MessageComponent) discordgo.ActionsRow {
	return discordgo.ActionsRow{Components: components}
}

// Button builds a standard button component.
// customID format: "action:sessionId" or "action:sessionId:extra"
func Button(label, customID string, style discordgo.ButtonStyle, emoji *Emoji, disabled bool) discordgo.Button {
	btn := discordgo.Button{
		Label:    label,
		Style:    style,
		CustomID: customID,
		Disabled: disabled,
	}
	if emoji != nil {
		btn.Emoji = &discordgo.ComponentEmoji{
			Name:     emoji.Name,
			ID:       emoji.ID,
			Animated: emoji.Animated,
		}
	}
	return btn
}

// NavRow builds the standard 3-button navigation row (Làm mới / Quay lại / Đóng).
func NavRow(sessionID, currentPage, parentPage string) discordgo.ActionsRow {
	refreshID := "nav:refresh:" + sessionID + ":" + currentPage
	backID := "nav:back:" + sessionID + ":" + parentPage
	closeID := "nav:close:" + sessionID

	return ActionRow(
		Button("Làm mới", refreshID, BtnSecondary, &EmojiRefresh, false),
		Button("Quay lại", backID, BtnSecondary, &EmojiBack, parentPage == ""),
		Button("Đóng", closeID, BtnDanger, &EmojiClose, false),
	)
}

// SelectMenu builds a standard string select menu component.
func SelectMenu(customID, placeholder string, options []discordgo.SelectMenuOption) discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
	}
}

// SelectOption builds a select menu option with an optional emoji.
func SelectOption(label, value, description string, emoji *Emoji, isDefault bool) discordgo.SelectMenuOption {
	opt := discordgo.SelectMenuOption{
		Label:       label,
		Value:       value,
		Description: description,
		Default:     isDefault,
	}
	if emoji != nil {
		opt.Emoji = &discordgo.ComponentEmoji{
			Name:     emoji.Name,
			ID:       emoji.ID,
			Animated: emoji.Animated,
		}
	}
	return opt
}
