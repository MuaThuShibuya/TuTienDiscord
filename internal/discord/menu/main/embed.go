// File: internal/discord/menu/main/embed.go
// Chức năng: Tạo Discord embed cho trang Main Menu.
// Ghi chú: Chỉ nhận ViewModel đã xử lý, không gọi DB hay service.

package mainmenu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

// BuildMenuResponse tạo response Discord đầy đủ cho trang Main Menu (lần đầu mở menu).
func BuildMenuResponse(vm *menu.MainMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildEmbed(vm)},
		Components: buildComponents(vm),
	}
}

// BuildMenuEdit tạo response cập nhật message khi điều hướng về Main Menu.
func BuildMenuEdit(vm *menu.MainMenuVM) *discordgo.InteractionResponseData {
	return BuildMenuResponse(vm)
}

func buildEmbed(vm *menu.MainMenuVM) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{Name: emoji.Realm.String() + " Cảnh Giới", Value: vm.RealmDisplay, Inline: true},
		{Name: emoji.CombatPower.String() + " Chiến Lực", Value: vm.CombatPower, Inline: true},
		{Name: emoji.MindState.String() + " Tâm Cảnh", Value: vm.MindState, Inline: true},
		{Name: emoji.Skill.String() + " Đạo Lộ", Value: vm.PathDisplay, Inline: true},
		{Name: emoji.Stamina.String() + " Thể Lực", Value: vm.StaminaBar, Inline: false},
		{Name: emoji.Cultivate.String() + " Tu Vi", Value: vm.ExpBar, Inline: false},
		{Name: emoji.SpiritStone.String() + " Linh Thạch", Value: vm.SpiritStones, Inline: true},
		{Name: emoji.SpiritJade.String() + " Linh Ngọc", Value: vm.SpiritJades, Inline: true},
		{Name: emoji.FateTicket.String() + " Vé Cơ Duyên", Value: vm.FateTickets, Inline: true},
		{Name: emoji.Info.String() + " Gợi ý hôm nay", Value: "_" + vm.DailyTip + "_"},
	}

	return &discordgo.MessageEmbed{
		Title:       emoji.Profile.String() + " Vạn Pháp Tiên Nghịch — " + vm.DaoName,
		Description: "Chào mừng trở lại, **" + vm.DaoName + "**!\nHãy chọn chức năng bên dưới.",
		Color:       ui.ColorDefault,
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · v0.1"},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
