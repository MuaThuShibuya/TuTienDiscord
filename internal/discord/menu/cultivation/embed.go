// File: internal/discord/menu/cultivation/embed.go
// Chức năng: Tạo Discord embed cho trang Tu Luyện.
// Ghi chú: Chỉ nhận ViewModel đã xử lý, không gọi DB hay service.

package cultivmenu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

// BuildMenuResponse tạo response Discord đầy đủ cho trang Tu Luyện.
func BuildMenuResponse(vm *menu.CultivationMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildEmbed(vm)},
		Components: buildComponents(vm),
	}
}

func buildEmbed(vm *menu.CultivationMenuVM) *discordgo.MessageEmbed {
	desc := "Con đường tu tiên vạn dặm bắt đầu từ một bước nhỏ."
	if !vm.HasPath {
		desc += "\n\n⚠️ **Đạo hữu chưa chọn đạo lộ.** Hãy chọn một hướng tu luyện bên dưới để mở khóa thiên phú tu luyện."
	}

	return &discordgo.MessageEmbed{
		Title:       emoji.Cultivate.String() + " Tu Luyện — " + vm.DaoName,
		Description: desc,
		Color:       ui.ColorCultivate,
		Fields: []*discordgo.MessageEmbedField{
			{Name: emoji.Realm.String() + " Cảnh Giới", Value: vm.RealmDisplay, Inline: true},
			{Name: emoji.MindState.String() + " Tâm Cảnh", Value: vm.MindState, Inline: true},
			{Name: emoji.Skill.String() + " Đạo Lộ", Value: vm.PathDisplay, Inline: true},
			{Name: emoji.Stamina.String() + " Thể Lực", Value: vm.StaminaBar, Inline: false},
			{Name: emoji.Cultivate.String() + " Tiến Độ Tu Vi", Value: vm.ExpBar, Inline: false},
			{Name: emoji.CombatPower.String() + " Chiến Lực", Value: vm.CombatPower, Inline: true},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Tu Luyện"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
