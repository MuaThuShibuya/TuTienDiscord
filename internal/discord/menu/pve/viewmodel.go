// File: internal/discord/menu/pve/viewmodel.go
package pve

import "github.com/whiskey/tu-tien-bot/internal/game/combat"

type PvEMenuViewModel struct {
	SessionID   string
	DaoName     string
	CombatPower string
	Areas       []AreaViewModel
}

type AreaViewModel struct {
	ID          string
	Name        string
	Description string
	NextStage   int
	RecommendCP string
	StaminaCost int64
}

type CombatViewModel struct {
	SessionID string
	AreaName  string
	Stage     int
	State     combat.SessionState
	Turn      int

	PlayerName  string
	PlayerHPStr string // Bar + Text
	PlayerRage  int64
	PlayerStats combat.CombatStats

	Enemies      []EnemyViewModel
	Logs         []string
	IsPlayerTurn bool
	TargetID     string // ID kẻ địch đầu tiên còn sống (auto-target)
}

type EnemyViewModel struct {
	ID     string
	Name   string
	Level  int
	HPStr  string // Bar + Text
	IsDead bool
	Stats  combat.CombatStats
}

type CombatRewardViewModel struct {
	SessionID string
	Stage     int
	IsClaimed bool
	Rewards   []string
}
