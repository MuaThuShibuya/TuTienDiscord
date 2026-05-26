// File: internal/discord/menu/cultivation/components.go
// Chức năng: Tạo Discord components (button, select menu) cho trang Tu Luyện.

package cultivmenu

import (
	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
)

func buildComponents(vm *menu.CultivationMenuVM) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	// Hàng 1: Các nút hành động tu luyện
	row1 := ui.ActionRow(
		ui.Button("Tĩnh Tu",
			menu.Build(menu.DomainCultivation, menu.ActionMeditate, vm.SessionID),
			ui.BtnPrimary, emoji.Cultivate, false),
		ui.Button("Bế Quan",
			menu.Build(menu.DomainCultivation, menu.ActionClosedDoor, vm.SessionID),
			ui.BtnPrimary, emoji.Lock, false),
		ui.Button("Luyện Thể",
			menu.Build(menu.DomainCultivation, menu.ActionBodyTraining, vm.SessionID),
			ui.BtnSecondary, emoji.Stamina, false),
		ui.Button("Đột Phá",
			menu.Build(menu.DomainCultivation, menu.ActionBreakthrough, vm.SessionID),
			ui.BtnDanger, emoji.Breakthrough, !vm.CanBreakthrough),
	)
	components = append(components, row1)

	// Hàng 2: Select menu chọn đạo lộ (chỉ hiện khi chưa có đạo lộ)
	if !vm.HasPath {
		pathSelect := ui.SelectMenu(
			menu.Build(menu.DomainCultivation, menu.ActionChoosePath, vm.SessionID),
			"✦ Lựa Chọn Đạo Lộ Tương Lai...",
			[]discordgo.SelectMenuOption{
				ui.SelectOption("Kiếm Tu", string(cultivation.PathSword),
					"Đạo của kiếm, chú trọng tấn công chí mạng", emoji.Sword, false),
				ui.SelectOption("Thể Tu", string(cultivation.PathBody),
					"Lấy thân làm gốc, phòng ngự tuyệt đối", emoji.Stamina, false),
				ui.SelectOption("Linh Tu", string(cultivation.PathSpirit),
					"Câu thông thiên địa, pháp thuật biến ảo", emoji.MindState, false),
				ui.SelectOption("Độc Tu", string(cultivation.PathPoison),
					"Sử dụng kịch độc, rút mòn sinh lực", emoji.Skill, false),
			},
		)
		components = append(components, ui.ActionRow(pathSelect))
	}

	components = append(components, ui.NavRow(vm.SessionID, string(menu.PageCultivation), string(menu.PageMain)))
	return components
}
