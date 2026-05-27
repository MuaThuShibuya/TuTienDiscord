// File: internal/discord/menu/pve/embed.go
package pve

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

func BuildAreaSelectEmbed(vm PvEMenuViewModel, isBiCanh bool) *discordgo.MessageEmbed {
	title := emoji.Map.String() + " Du Ngoạn Thiên Hạ"
	desc := "Đạo hữu có thể dạo bước tứ phương, diệt yêu trừ ma, thu thập kỳ trân dị thảo."
	if isBiCanh {
		title = emoji.Dungeon.String() + " Bí Cảnh Hung Hiểm"
		desc = "Nơi giấu vô vàn đại cơ duyên, nhưng cũng đầy rẫy hiểm nguy."
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       ui.ColorCombat,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	for _, area := range vm.Areas {
		val := fmt.Sprintf("Tiến trình: **Ải %d**\nĐề xuất: %s **%s**\nTiêu hao: %s **%d**",
			area.NextStage, emoji.CombatPower.String(), area.RecommendCP, emoji.Stamina.String(), area.StaminaCost)
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   area.Name,
			Value:  val,
			Inline: true,
		})
	}
	return embed
}

func BuildCombatEmbed(vm CombatViewModel) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s %s (Ải %d) — Hiệp %d", emoji.Sword.String(), vm.AreaName, vm.Stage, vm.Turn),
		Color: ui.ColorCombat,
	}

	// Status Player
	playerStatus := fmt.Sprintf("%s **%s**\n%s\n%s Nộ khí: %d/100",
		emoji.Profile.String(), vm.PlayerName, vm.PlayerHPStr, emoji.CombatPower.String(), vm.PlayerRage)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Đạo Hữu", Value: playerStatus, Inline: false})

	// Status Enemies
	var enemyStrs []string
	for _, e := range vm.Enemies {
		icon := emoji.Monster.String()
		if e.IsDead {
			icon = emoji.Cross.String()
		}
		enemyStrs = append(enemyStrs, fmt.Sprintf("%s **%s** (Lv.%d)\n%s", icon, e.Name, e.Level, e.HPStr))
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Kẻ Địch", Value: strings.Join(enemyStrs, "\n\n"), Inline: false})

	// Battle Log
	logStr := "*Chưa có động tĩnh gì.*"
	if len(vm.Logs) > 0 {
		logStr = strings.Join(vm.Logs, "\n")
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Diễn Biến Trận Đấu", Value: logStr, Inline: false})

	// Footer Status
	if vm.State == combat.StateWon {
		embed.Color = ui.ColorSuccess
	} else if vm.State == combat.StateLost {
		embed.Color = ui.ColorError
	}
	return embed
}
