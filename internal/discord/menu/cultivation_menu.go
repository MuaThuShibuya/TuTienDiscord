// File: internal/discord/menu/cultivation_menu.go
// Phiên bản: v0.1.1
// Mục đích: UI Builder cho trang Tu Luyện — render embed cảnh giới, tiến độ tu vi và các nút hành động.
// Bảo mật: Chỉ nhận ViewModel đã xử lý sẵn, không gọi DB hay service, không có side effect.
// Ghi chú: Các action tĩnh tu/đột phá sẽ được nối với cooldown check trong v0.2.

package menu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
)

// BuildCultivationMenuResponse tạo response Discord cho trang Tu Luyện.
func BuildCultivationMenuResponse(vm *CultivationMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildCultivationEmbed(vm)},
		Components: buildCultivationComponents(vm.SessionID, vm.CanBreakthrough),
	}
}

func buildCultivationEmbed(vm *CultivationMenuVM) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       ui.EmojiCultivate.String() + " Tu Luyện — " + vm.DaoName,
		Description: "Con đường tu tiên vạn dặm bắt đầu từ một bước nhỏ.",
		Color:       ui.ColorCultivate,
		Fields: []*discordgo.MessageEmbedField{
			{Name: ui.EmojiRealm.String() + " Cảnh Giới", Value: vm.RealmDisplay, Inline: true},
			{Name: ui.EmojiMindState.String() + " Tâm Cảnh", Value: vm.MindState, Inline: true},
			{Name: ui.EmojiSkill.String() + " Đạo Lộ", Value: vm.PathDisplay, Inline: true},
			{Name: ui.EmojiStamina.String() + " Thể Lực", Value: vm.StaminaBar, Inline: false},
			{Name: ui.EmojiCultivate.String() + " Tiến Độ Tu Vi", Value: vm.ExpBar, Inline: false},
			{Name: ui.EmojiCombatPower.String() + " Chiến Lực", Value: vm.CombatPower, Inline: true},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Tu Luyện"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildCultivationComponents(sessionID string, canBreakthrough bool) []discordgo.MessageComponent {
	row1 := ui.ActionRow(
		ui.Button("Tĩnh Tu", Build(DomainCultivation, ActionMeditate, sessionID), ui.BtnPrimary, &ui.EmojiCultivate, false),
		ui.Button("Bế Quan", Build(DomainCultivation, ActionClosedDoor, sessionID), ui.BtnPrimary, &ui.EmojiLock, false),
		ui.Button("Luyện Thể", Build(DomainCultivation, ActionBodyTraining, sessionID), ui.BtnSecondary, &ui.EmojiStamina, false),
		ui.Button("Đột Phá", Build(DomainCultivation, ActionBreakthrough, sessionID), ui.BtnDanger, &ui.EmojiBreakthrough, false),
	)
	navRow := ui.NavRow(sessionID, string(PageCultivation), string(PageMain))

	return []discordgo.MessageComponent{row1, navRow}
}
