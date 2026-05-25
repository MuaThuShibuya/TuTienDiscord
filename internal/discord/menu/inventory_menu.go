// File: internal/discord/menu/inventory_menu.go
package menu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
)

func BuildInventoryMenuResponse(vm *InventoryMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildInventoryEmbed(vm)},
		Components: buildInventoryComponents(vm),
	}
}

func buildInventoryEmbed(vm *InventoryMenuVM) *discordgo.MessageEmbed {
	desc := fmt.Sprintf("📦 **Sức chứa:** `%s`\n\n", vm.SlotUsage)
	
	if len(vm.Items) == 0 {
		desc += "_Túi đồ của đạo hữu trống rỗng. Hãy đi rèn luyện để tìm kiếm kỳ ngộ!_"
	} else {
		for i, item := range vm.Items {
			desc += fmt.Sprintf("**%d.** [%s] %s x%d\n", i+1, item.Rarity, item.Name, item.Quantity)
		}
	}

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiBag.String() + " Túi Đồ — " + vm.DaoName,
		Description: desc,
		Color:       ui.ColorEconomy,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Túi Đồ"},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}

func buildInventoryComponents(vm *InventoryMenuVM) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	// Nút sử dụng Đan Dược
	var usableOptions []discordgo.SelectMenuOption
	for _, item := range vm.Items {
		if item.IsUsable {
			usableOptions = append(usableOptions, ui.SelectOption(
				item.Name, item.InstanceID, fmt.Sprintf("Số lượng: %d", item.Quantity), &ui.EmojiBag, false,
			))
		}
	}

	if len(usableOptions) > 0 {
		useSelect := ui.SelectMenu(
			Build(DomainInventory, ActionInventoryUse, vm.SessionID),
			"💊 Sử dụng đan dược...",
			usableOptions,
		)
		components = append(components, ui.ActionRow(useSelect))
	}

	components = append(components, ui.NavRow(vm.SessionID, string(PageInventory), string(PageMain)))
	return components
}
