// File: internal/discord/ui/components.go
// Chức năng: Builder functions cho Discord UI components — button, select menu, action row.
// Ghi chú: Mọi emoji truyền vào phải là *emoji.Emoji từ package ui/emoji — không định nghĩa emoji ở đây.
//          customID phải encode sessionID để handler xác thực quyền sở hữu.

package ui

import (
	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

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
// em có thể nil nếu không cần icon.
func Button(label, customID string, style discordgo.ButtonStyle, em *emoji.Emoji, disabled bool) discordgo.Button {
	btn := discordgo.Button{
		Label:    label,
		Style:    style,
		CustomID: customID,
		Disabled: disabled,
	}
	if em != nil {
		btn.Emoji = em.Component()
	}
	return btn
}

// NavRow builds the standard navigation row (Quay lại / Đóng).
// parentPage trống → nút Quay lại bị disabled.
func NavRow(sessionID, currentPage, parentPage string) discordgo.ActionsRow {
	backID := "nav:back:" + sessionID + ":" + parentPage
	closeID := "nav:close:" + sessionID

	return ActionRow(
		Button("Quay lại", backID, BtnSecondary, emoji.Back, parentPage == ""),
		Button("Đóng", closeID, BtnDanger, emoji.Close, false),
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
// em có thể nil nếu không cần icon.
func SelectOption(label, value, description string, em *emoji.Emoji, isDefault bool) discordgo.SelectMenuOption {
	opt := discordgo.SelectMenuOption{
		Label:       label,
		Value:       value,
		Description: description,
		Default:     isDefault,
	}
	if em != nil {
		opt.Emoji = em.Component()
	}
	return opt
}
