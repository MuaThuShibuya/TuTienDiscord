// File: internal/discord/ui/errors.go
// Phiên bản: v0.1.1
// Mục đích: Hàm tiện ích để trả về response lỗi cho Discord interaction.
//           Tách ra khỏi embeds.go để dễ tìm khi cần xử lý lỗi.
// Bảo mật: Không bao giờ trả về chi tiết lỗi nội bộ (stack trace, URI DB) cho user.
// Ghi chú: Dùng RespondEphemeralError cho lỗi đầu tiên (chưa defer).
//          Dùng EditEphemeralError sau khi đã defer response.

package ui

import "github.com/bwmarrin/discordgo"

// RespondEphemeralError gửi response lỗi ephemeral (chỉ user thấy) cho interaction chưa được respond.
func RespondEphemeralError(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{ErrorEmbed(message)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// RespondEphemeralWarning gửi response cảnh báo ephemeral.
func RespondEphemeralWarning(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{WarningEmbed(message)},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}

// UpdateWithError chỉnh sửa message hiện tại thành thông báo lỗi (dùng trong menu router).
func UpdateWithError(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{ErrorEmbed(message)},
			Components: []discordgo.MessageComponent{},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

// EditEphemeralError chỉnh sửa deferred response thành thông báo lỗi.
// Dùng sau khi đã gọi DeferredChannelMessageWithSource.
func EditEphemeralError(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{ErrorEmbed(message)},
	})
}

// EditEphemeralEmbed chỉnh sửa deferred response thành bất kỳ embed nào.
func EditEphemeralEmbed(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}
