// File: internal/discord/menu/equipment/embed.go
// Chức năng: Tạo Discord embed cho trang Trang Bị.
// Ghi chú: Chỉ nhận ViewModel đã xử lý, không gọi DB hay service.

package equipmenu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

// BuildMenuResponse tạo response Discord đầy đủ cho trang Trang Bị.
func BuildMenuResponse(vm *menu.EquipmentMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildEmbed(vm)},
		Components: buildComponents(vm),
	}
}

func buildEmbed(vm *menu.EquipmentMenuVM) *discordgo.MessageEmbed {
	slots := []struct {
		icon string
		name string
		item *menu.EquippedItemVM
	}{
		{"⚔️", "Vũ Khí", vm.Weapon},
		{"🛡️", "Giáp", vm.Armor},
		{"🔮", "Pháp Bảo", vm.Artifact},
		{"💎", "Bảo Vật", vm.Treasure},
	}

	var fields []*discordgo.MessageEmbedField
	for _, s := range slots {
		value := "_Chưa mặc_"
		if s.item != nil {
			value = fmt.Sprintf("`[%s]` %s", s.item.Rarity, s.item.Name)
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   s.icon + " " + s.name,
			Value:  value,
			Inline: true,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   emoji.CombatPower.String() + " Chiến Lực",
		Value:  vm.CombatPower,
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title:       emoji.Equip.String() + " Trang Bị — " + vm.DaoName,
		Description: "Trang bị hợp lý giúp chiến lực tăng vọt, bước đường tu tiên thêm vững chắc.",
		Color:       ui.ColorCombat,
		Fields:      fields,
		Footer:      &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Trang Bị"},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
