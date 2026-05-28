// File: internal/discord/menu/inventory/embed.go
// Chức năng: Tạo Discord embed cho trang Túi Đồ.
// Ghi chú: Items là danh sách đã phân trang. UsableItems là tất cả item có thể dùng (cho select menu).

package invmenu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"go.uber.org/zap"
)

// BuildMenuResponse tạo response Discord đầy đủ cho trang Túi Đồ.
func BuildMenuResponse(vm *menu.InventoryMenuVM) *discordgo.InteractionResponseData {
	embed := buildEmbed(vm)
	comps := buildComponents(vm)

	logger.L().Debug("BuildMenuResponse trace",
		zap.Int("len(vm.Items)", len(vm.Items)),
		zap.Int("len(vm.UsableItems)", len(vm.UsableItems)),
		zap.Int("len(components)", len(comps)),
		zap.Int("currentPage", vm.CurrentPage),
		zap.Int("totalPages", vm.TotalPages),
	)

	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: comps,
	}
}

func buildEmbed(vm *menu.InventoryMenuVM) *discordgo.MessageEmbed {
	desc := fmt.Sprintf("📦 **Sức chứa:** `%s`\n\n", vm.SlotUsage)

	if len(vm.Items) == 0 {
		if vm.TotalPages <= 1 {
			desc += "_Túi càn khôn của đạo hữu hiện chưa có vật phẩm nào._"
		} else {
			desc += "_Trang này không có vật phẩm._"
		}
	} else {
		for idx, it := range vm.Items {
			rarityTag := it.Rarity
			if rarityTag == "" {
				rarityTag = "?"
			}
			typeTag := "📦"
			if it.IsEquip {
				typeTag = "🛡️"
			} else if it.IsUsable {
				typeTag = "💊"
			}
			desc += fmt.Sprintf("**%d.** %s `[%s]` %s × %d\n",
				idx+1, typeTag, rarityTag, it.Name, it.Quantity)
		}
	}

	return &discordgo.MessageEmbed{
		Title:       emoji.Bag.String() + " Túi Đồ — " + vm.DaoName,
		Description: desc,
		Color:       ui.ColorEconomy,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Vạn Pháp Tiên Nghịch · Trang %d/%d", vm.CurrentPage, vm.TotalPages)},
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}
}
