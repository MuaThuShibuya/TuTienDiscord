// File: internal/discord/menu/pve/components.go
package pve

import (
	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

func BuildPvEMainComponents(menuSessionID string) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		ui.ActionRow(
			ui.Button("Du Ngoạn", menu.Build(menu.DomainPvE, menu.ActionPvEDuNgoan, menuSessionID), ui.BtnPrimary, emoji.Map, false),
			ui.Button("Bí Cảnh", menu.Build(menu.DomainPvE, menu.ActionPvEBiCanh, menuSessionID), ui.BtnDanger, emoji.Dungeon, false),
		),
		ui.NavRow(menuSessionID, string(menu.PagePvE), string(menu.PageMain)),
	}
}

func BuildAreaSelectComponents(menuSessionID string, vm PvEMenuViewModel) []discordgo.MessageComponent {
	var opts []discordgo.SelectMenuOption
	for _, area := range vm.Areas {
		opts = append(opts, ui.SelectOption(area.Name, area.ID, "Bắt đầu khám phá khu vực này", emoji.Sword, false))
	}

	return []discordgo.MessageComponent{
		ui.ActionRow(ui.SelectMenu(menu.Build(menu.DomainPvE, menu.ActionPvEStart, menuSessionID), "Chọn khu vực muốn khiêu chiến...", opts)),
		ui.NavRow(menuSessionID, string(menu.PagePvE), string(menu.PagePvE)),
	}
}

func BuildCombatActionComponents(cache *ActionCache, menuSessionID string, vm CombatViewModel) []discordgo.MessageComponent {
	if vm.State == combat.StateWon {
		return []discordgo.MessageComponent{
			ui.ActionRow(
				ui.Button("Nhận Thưởng", menu.Build(menu.DomainPvE, menu.ActionPvEClaim, menuSessionID, vm.SessionID), ui.BtnSuccess, emoji.Reward, false),
				ui.Button("Rời Khỏi", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSessionID, string(menu.PagePvE)), ui.BtnSecondary, emoji.Back, false),
			),
		}
	}

	// Nút quay lại dành cho màn thua hoặc cần lối thoát an toàn
	backBtn := ui.Button("Rời Khỏi", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSessionID, string(menu.PagePvE)), ui.BtnSecondary, emoji.Escape, false)

	if vm.State == combat.StateLost {
		return []discordgo.MessageComponent{ui.ActionRow(backBtn)}
	}

	// Đang đánh (StateOngoing)
	// Gói payload vào RAM Cache, trả về token 8 ký tự
	payload := PvEActionPayload{
		OwnerID:         vm.PlayerID,
		MenuSessionID:   menuSessionID,
		CombatSessionID: vm.SessionID,
		TargetID:        vm.TargetID,
	}
	atkToken := cache.Save(payload)
	autoToken := cache.Save(payload)

	attackID := menu.Build(menu.DomainPvE, menu.ActionPvEAttack, menuSessionID, atkToken)
	skillID := menu.Build(menu.DomainPvE, menu.ActionPvESkill, menuSessionID, "dummy")
	autoID := menu.Build(menu.DomainPvE, menu.ActionPvEAuto, menuSessionID, autoToken)
	escapeID := menu.Build(menu.DomainPvE, menu.ActionPvEEscape, menuSessionID, "dummy")

	return []discordgo.MessageComponent{
		ui.ActionRow(
			ui.Button("Tấn Công", attackID, ui.BtnPrimary, emoji.Sword, !vm.IsPlayerTurn),
			ui.Button("Kỹ Năng", skillID, ui.BtnPrimary, emoji.Skill, !vm.IsPlayerTurn),
			ui.Button("Auto x100", autoID, ui.BtnSuccess, emoji.Auto, !vm.IsPlayerTurn),
			ui.Button("Bỏ Chạy", escapeID, ui.BtnDanger, emoji.Escape, false),
		),
	}
}
