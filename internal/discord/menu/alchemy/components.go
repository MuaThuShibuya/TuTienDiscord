// File: internal/discord/menu/alchemy/components.go
// Chức năng: Tạo nút điều hướng cho trang Lò Đan.

package alchemymenu

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildComponents(vm *menu.AlchemyMenuVM) []discordgo.MessageComponent {
	// Nếu đang xem chi tiết, đổi Select Menu thành các Nút tương tác
	if vm.SelectedRecipe != nil {
		craftBtn := ui.Button("Luyện Đan", menu.Build(menu.DomainAlchemy, menu.ActionAlchemyCraft, vm.SessionID, vm.SelectedRecipe.ID), ui.BtnPrimary, emoji.Alchemy, !vm.SelectedRecipe.CanCraft)
		cancelBtn := ui.Button("Quay lại", menu.Build(menu.DomainAlchemy, menu.ActionAlchemyCancel, vm.SessionID), ui.BtnDanger, emoji.Close, false)
		return []discordgo.MessageComponent{
			ui.ActionRow(craftBtn, cancelBtn),
		}
	}

	var options []discordgo.SelectMenuOption

	for _, r := range vm.Recipes {
		options = append(options, ui.SelectOption(
			r.Name,
			r.ID,
			fmt.Sprintf("Yêu cầu Cấp %d - Tỷ lệ: %s", r.LevelRequired, r.SuccessRate),
			emoji.Alchemy,
			false,
		))
	}

	craftSelect := ui.SelectMenu(
		menu.Build(menu.DomainAlchemy, menu.ActionAlchemyView, vm.SessionID),
		emoji.Alchemy.Fallback+" Chọn đan dược để xem chi tiết...",
		options,
	)

	return []discordgo.MessageComponent{
		ui.ActionRow(craftSelect),
		ui.NavRow(vm.SessionID, string(menu.PageAlchemy), string(menu.PageMain)),
	}
}
