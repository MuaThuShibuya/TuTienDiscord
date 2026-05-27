// File: internal/discord/commands.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa schema tất cả slash command đăng ký với Discord API.
// Bảo mật: Mô tả lệnh không được tiết lộ chi tiết kỹ thuật nội bộ.
// Ghi chú: Thêm lệnh mới ở đây khi xây dựng tính năng. Bot tự đăng ký lại mỗi lần khởi động.

package discord

import "github.com/bwmarrin/discordgo"

// AllCommands trả về danh sách tất cả slash command cần đăng ký với Discord.
func AllCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "dev",
			Description: "Thiên Mệnh Cai Trị.",
		},
		{
			Name:        "start",
			Description: "Bắt đầu hành trình tu tiên của đạo hữu.",
		},
		{
			Name:        "menu",
			Description: "Mở giao diện chính của Vạn Pháp Tiên Nghịch.",
		},
		// TODO v0.2: /cultivate — tĩnh tu nhanh không cần mở menu
		// TODO v0.5: /gacha    — quay cơ duyên nhanh
		// TODO v0.8: /market   — truy cập chợ nhanh
	}
}
