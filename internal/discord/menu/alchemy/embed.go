// File: internal/discord/menu/alchemy/embed.go
// Chức năng: Xây dựng giao diện tĩnh cho Lò Đan.

package alchemymenu

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func BuildMenuResponse(vm *menu.AlchemyMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildEmbed(vm)},
		Components: buildComponents(vm),
	}
}

func buildEmbed(vm *menu.AlchemyMenuVM) *discordgo.MessageEmbed {
	// Nếu đang ở trạng thái xem chi tiết 1 Đan Dược
	if vm.SelectedRecipe != nil {
		r := vm.SelectedRecipe
		return &discordgo.MessageEmbed{
			Title:       emoji.Alchemy.String() + " Chi tiết: " + r.Name,
			Description: fmt.Sprintf("**Yêu cầu Cấp:** %d\n**Tỷ lệ thành công:** %s\n\n**Nguyên liệu yêu cầu:**\n%s", r.LevelRequired, r.SuccessRate, r.Materials),
			Color:       ui.ColorCultivate,
		}
	}

	// Trạng thái tổng quan của Lò Đan
	desc := fmt.Sprintf("Danh hiệu: **%s**\nCấp độ Luyện Đan Sư: **Cấp %d**\nKinh nghiệm: %s\n\n", vm.Title, vm.Level, vm.ExpBar)
	desc += "Sử dụng đan dược giúp gia tăng tốc độ tu luyện hoặc đột phá rào cản cảnh giới."

	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "💡 Mẹo",
		Value:  "_" + vm.DailyTip + "_",
		Inline: false,
	})

	embed := &discordgo.MessageEmbed{
		Title:       emoji.Alchemy.String() + " Lò Luyện Đan",
		Description: desc,
		Color:       ui.ColorCultivate,
		Fields:      fields,
	}

	return embed
}
