// File: internal/discord/ui/embeds.go
// Phiên bản: v0.1.1
// Mục đích: Các hàm tạo Discord embed chuẩn dùng xuyên suốt bot.
// Bảo mật: Không nhúng trực tiếp input của user vào embed mà không sanitize.
// Ghi chú: Hàm gửi/chỉnh response lỗi nằm trong ui/errors.go.

package ui

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// ErrorEmbed tạo embed thông báo lỗi màu đỏ.
func ErrorEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: EmojiError.String() + " " + message,
		Color:       ColorError,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// SuccessEmbed tạo embed thông báo thành công màu xanh lá.
func SuccessEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       EmojiSuccess.String() + " " + title,
		Description: message,
		Color:       ColorSuccess,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// InfoEmbed tạo embed thông tin màu xanh dương.
func InfoEmbed(title, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       EmojiInfo.String() + " " + title,
		Description: message,
		Color:       ColorInfo,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

// WarningEmbed tạo embed cảnh báo màu vàng.
func WarningEmbed(message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: EmojiWarning.String() + " " + message,
		Color:       ColorWarning,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
