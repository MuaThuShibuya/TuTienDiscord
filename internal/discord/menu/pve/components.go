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

func BuildCombatActionComponents(menuSessionID string, vm CombatViewModel) []discordgo.MessageComponent {
	if vm.State == combat.StateWon {
		return []discordgo.MessageComponent{
			ui.ActionRow(
				ui.Button("Nhận Thưởng", menu.Build(menu.DomainPvE, menu.ActionPvEClaim, menuSessionID, vm.SessionID), ui.BtnSuccess, emoji.Reward, false),
				ui.Button("Rời Khỏi", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSessionID, string(menu.PagePvE)), ui.BtnSecondary, emoji.Back, false),
			),
		}
	}

	if vm.State == combat.StateLost {
		return []discordgo.MessageComponent{
			ui.ActionRow(ui.Button("Rời Khỏi (Thất Bại)", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSessionID, string(menu.PagePvE)), ui.BtnSecondary, emoji.Back, false)),
		}
	}

	// Đang đánh
	// Truyền TargetID và Interaction ID (sẽ tự gen nonce) vào Extra
	attackID := menu.Build(menu.DomainPvE, menu.ActionPvEAttack, menuSessionID, vm.SessionID+"|"+vm.TargetID)
	return []discordgo.MessageComponent{
		ui.ActionRow(
			ui.Button("Tấn Công", attackID, ui.BtnPrimary, emoji.Sword, !vm.IsPlayerTurn),
			ui.Button("Kỹ Năng", "none_skill", ui.BtnSecondary, emoji.Skill, true), // TODO: Block sau
			ui.Button("Bỏ Chạy", "none_escape", ui.BtnDanger, emoji.Escape, true),  // TODO: pvecombat.Escape
		),
	}
}
