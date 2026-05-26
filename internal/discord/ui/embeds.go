// File: internal/discord/ui/embeds.go
// Chức năng: Các hàm tạo Discord embed chuẩn dùng xuyên suốt bot.
// Ghi chú: Hàm gửi/chỉnh response lỗi nằm trong ui/errors.go.

package ui

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

// ErrorEmbed tạo embed thông báo lỗi màu đỏ.
func ErrorEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: emoji.Error.String() + " " + message,
		Color:       ColorError,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessEmbed tạo embed thông báo thành công màu xanh lá.
func SuccessEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       emoji.Success.String() + " " + title,
		Description: message,
		Color:       ColorSuccess,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// InfoEmbed tạo embed thông tin màu xanh dương.
func InfoEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       emoji.Info.String() + " " + title,
		Description: message,
		Color:       ColorInfo,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// WarningEmbed tạo embed cảnh báo màu vàng.
func WarningEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: emoji.Warning.String() + " " + message,
		Color:       ColorWarning,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
