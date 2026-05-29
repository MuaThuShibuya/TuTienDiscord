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
		val := fmt.Sprintf("Tiến trình: **Ải %d**\nĐề xuất: %s **%s**",
			area.NextStage, emoji.CombatPower.String(), area.RecommendCP)
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

	// Current Turn Indicator
	turnDesc := ""
	if vm.State == combat.StateActive {
		if vm.IsPlayerTurn {
			turnDesc = "\n> " + emoji.Sword.String() + " **Lượt của đạo hữu! Hãy xuất chiêu.**"
		} else {
			turnDesc = "\n> ⏳ *Yêu vật đang vận khởi yêu khí...*"
		}
	}
	if turnDesc != "" {
		embed.Description = turnDesc
	}

	// Status Player
	playerStatus := fmt.Sprintf("%s **%s**\n%s\n%s Nộ khí: %d/100\n*ATK: %s | DEF: %s | Tốc: %s*",
		emoji.Profile.String(), vm.PlayerName, vm.PlayerHPStr, emoji.CombatPower.String(), vm.PlayerRage, FormatNumber(vm.PlayerStats.ATK), FormatNumber(vm.PlayerStats.DEF), FormatNumber(vm.PlayerStats.Speed))
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Đạo Hữu", Value: playerStatus, Inline: false})

	// Status Enemies
	var enemyStrs []string
	for _, e := range vm.Enemies {
		icon := emoji.Monster.String()
		if e.IsDead {
			icon = emoji.Cross.String()
		}
		targetMark := ""
		if vm.TargetID == e.ID && !e.IsDead {
			targetMark = " " + emoji.Sword.String() + " *(Mục tiêu)*"
		}
		enemyStrs = append(enemyStrs, fmt.Sprintf("%s **%s** (Lv.%d)%s\n%s\n*ATK: %s | DEF: %s | Tốc: %s*", icon, e.Name, e.Level, targetMark, e.HPStr, FormatNumber(e.Stats.ATK), FormatNumber(e.Stats.DEF), FormatNumber(e.Stats.Speed)))
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Kẻ Địch Yêu Khí", Value: strings.Join(enemyStrs, "\n\n"), Inline: false})

	// Battle Log
	logStr := "*Chưa có động tĩnh gì.*"
	if len(vm.Logs) > 0 {
		logStr = strings.Join(vm.Logs, "\n")
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Diễn Biến Trận Đấu", Value: logStr, Inline: false})

	// Footer Status
	switch vm.State {
	case combat.StateWon:
		embed.Color = ui.ColorSuccess
		embed.Description = "🎉 **CHIẾN THẮNG!** Đạo hữu đã chém giết yêu ma, thiên địa trở lại thanh minh."
	case combat.StateLost:
		embed.Color = ui.ColorError
		embed.Description = "💀 **THẤT BẠI!** Đạo thể trọng thương, linh khí cạn kiệt. Hãy tĩnh tu thêm rồi quay lại."
	}
	return embed
}
