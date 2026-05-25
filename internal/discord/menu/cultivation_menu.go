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
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
)

// BuildCultivationMenuResponse tạo response Discord cho trang Tu Luyện.
func BuildCultivationMenuResponse(vm *CultivationMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildCultivationEmbed(vm)},
		Components: buildCultivationComponents(vm),
	}
}

func buildCultivationEmbed(vm *CultivationMenuVM) *discordgo.MessageEmbed {
	desc := "Con đường tu tiên vạn dặm bắt đầu từ một bước nhỏ."
	if !vm.HasPath {
		desc += "\n\n⚠️ **Đạo hữu chưa chọn đạo lộ.** Hãy chọn một hướng tu luyện bên dưới để mở khóa thiên phú tu luyện."
	}

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiCultivate.String() + " Tu Luyện — " + vm.DaoName,
		Description: desc,
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

func buildCultivationComponents(vm *CultivationMenuVM) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	row1 := ui.ActionRow(
		ui.Button("Tĩnh Tu", Build(DomainCultivation, ActionMeditate, vm.SessionID), ui.BtnPrimary, &ui.EmojiCultivate, false),
		ui.Button("Bế Quan", Build(DomainCultivation, ActionClosedDoor, vm.SessionID), ui.BtnPrimary, &ui.EmojiLock, false),
		ui.Button("Luyện Thể", Build(DomainCultivation, ActionBodyTraining, vm.SessionID), ui.BtnSecondary, &ui.EmojiStamina, false),
		ui.Button("Đột Phá", Build(DomainCultivation, ActionBreakthrough, vm.SessionID), ui.BtnDanger, &ui.EmojiBreakthrough, !vm.CanBreakthrough),
	)
	components = append(components, row1)

	if !vm.HasPath {
		pathSelect := ui.SelectMenu(
			Build(DomainCultivation, ActionChoosePath, vm.SessionID),
			"✦ Lựa Chọn Đạo Lộ Tương Lai...",
			[]discordgo.SelectMenuOption{
				ui.SelectOption("Kiếm Tu", string(cultivation.PathSword), "Đạo của kiếm, chú trọng tấn công chí mạng", &ui.EmojiSword, false),
				ui.SelectOption("Thể Tu", string(cultivation.PathBody), "Lấy thân làm gốc, phòng ngự tuyệt đối", &ui.EmojiStamina, false),
				ui.SelectOption("Linh Tu", string(cultivation.PathSpirit), "Câu thông thiên địa, pháp thuật biến ảo", &ui.EmojiMindState, false),
				ui.SelectOption("Độc Tu", string(cultivation.PathPoison), "Sử dụng kịch độc, rút mòn sinh lực", &ui.EmojiSkill, false),
			},
		)
		components = append(components, ui.ActionRow(pathSelect))
	}

	navRow := ui.NavRow(vm.SessionID, string(PageCultivation), string(PageMain))
	components = append(components, navRow)

	return components
}
