// File: internal/discord/menu/pve/session_adapter.go
package pve

import (
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

// BuildHPBar sinh ra thanh HP trực quan (10 blocks).
func BuildHPBar(current, max int64) string {
	if max <= 0 {
		return "🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩 0/0"
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	percent := float64(current) / float64(max)
	filled := int(percent * 10)
	if current > 0 && filled == 0 {
		filled = 1 // Vẫn còn sống thì ít nhất 1 vạch
	}

	bar := ""
	for i := 0; i < 10; i++ {
		if i < filled {
			bar += "🟩"
		} else {
			bar += "🟥"
		}
	}
	return fmt.Sprintf("`%s` %d/%d", bar, current, max)
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
		IsPlayerTurn: session.IsPlayerTurn(),
	}

	for _, e := range session.Enemies {
		isDead := e.CurrentHP <= 0
		vm.Enemies = append(vm.Enemies, EnemyViewModel{
			ID: e.ID, Name: e.Name, Level: e.Level,
			HPStr: BuildHPBar(e.CurrentHP, e.Stats.MaxHP), IsDead: isDead,
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
