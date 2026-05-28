// File: internal/discord/menu/pve/session_adapter.go
package pve

import (
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

// FormatNumber thêm dấu phẩy phân cách hàng nghìn.
func FormatNumber(n int64) string {
	in := fmt.Sprintf("%d", n)
	out := make([]byte, len(in)+(len(in)-1)/3)
	for i, j, k := len(in)-1, len(out)-1, 0; i >= 0; i, j = i-1, j-1 {
		out[j] = in[i]
		k++
		if k == 3 && i > 0 {
			j--
			out[j] = ','
			k = 0
		}
	}
	return string(out)
}

// BuildHPBar sinh ra thanh HP trực quan (10 blocks) kèm phần trăm và text.
func BuildHPBar(current, max int64) string {
	if max <= 0 {
		return fmt.Sprintf("`%s` **0%%** (0/0)", utils.ProgressBar(0, 1, 10))
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	percent := float64(current) / float64(max)
	percentInt := int(percent * 100)
	bar := utils.ProgressBar(int(current), int(max), 10)
	return fmt.Sprintf("`%s` **%d%%** (%s/%s)", bar, percentInt, FormatNumber(current), FormatNumber(max))
}

func CombatSessionToViewModel(session *combat.CombatSession, areaName string) CombatViewModel {
	vm := CombatViewModel{
		SessionID:    session.ID, // Trùng CombatSession.ID (khác menu SessionID)
		AreaName:     areaName,
		Stage:        session.Stage,
		State:        session.State,
		Turn:         session.Turn,
		PlayerName:   session.Player.Name,
		PlayerHPStr:  BuildHPBar(session.Player.CurrentHP, session.Player.Stats.MaxHP),
		PlayerRage:   session.Player.CurrentRage,
		PlayerStats:  session.Player.Stats,
		IsPlayerTurn: session.IsPlayerTurn(),
	}

	for _, e := range session.Enemies {
		isDead := e.CurrentHP <= 0
		vm.Enemies = append(vm.Enemies, EnemyViewModel{
			ID: e.ID, Name: e.Name, Level: e.Level,
			HPStr: BuildHPBar(e.CurrentHP, e.Stats.MaxHP), IsDead: isDead, Stats: e.Stats,
		})
		if !isDead && vm.TargetID == "" {
			vm.TargetID = e.ID // Auto-target con đầu tiên sống
		}
	}

	// Lấy 5 log gần nhất (nếu có)
	start := len(session.Logs) - 5
	if start < 0 {
		start = 0
	}
	for _, l := range session.Logs[start:] {
		vm.Logs = append(vm.Logs, fmt.Sprintf("**[Hiệp %d]** %s", l.Turn, l.Message))
	}

	return vm
}
