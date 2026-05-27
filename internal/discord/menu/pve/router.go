// File: internal/discord/menu/pve/router.go
package pve

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	"github.com/whiskey/tu-tien-bot/internal/game/combat"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"github.com/whiskey/tu-tien-bot/internal/game/pvecombat"
)

type Router struct {
	pvecombatSvc *pvecombat.Service
	combatSvc    *combat.Service
	pveRepo      pve.ProgressRepository
	log          *zap.Logger
}

func NewRouter(pvecombatSvc *pvecombat.Service, combatSvc *combat.Service, pveRepo pve.ProgressRepository, log *zap.Logger) *Router {
	return &Router{pvecombatSvc: pvecombatSvc, combatSvc: combatSvc, pveRepo: pveRepo, log: log.Named("menu.pve")}
}

func (r *Router) HandlePvEInteraction(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, action, extra string) {
	ctx := context.Background()
	userID := menuSession.UserID

	switch action {
	case menu.ActionPvEDuNgoan, menu.ActionPvEBiCanh:
		isBiCanh := action == menu.ActionPvEBiCanh
		r.renderAreaSelect(s, i, menuSession, isBiCanh)

	case menu.ActionPvEStart:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			ui.EditEphemeralError(s, i, ui.MsgGenericError)
			return
		}
		areaID := data.Values[0]

		cSession, err := r.pvecombatSvc.StartPvECombat(ctx, userID, areaID)
		if err != nil {
			ui.EditEphemeralError(s, i, fmt.Sprintf("Không thể khởi tạo ải: %v", err))
			return
		}

		r.renderCombatScreen(s, i, menuSession, cSession)

	case menu.ActionPvEAttack:
		parts := strings.SplitN(extra, "|", 2)
		if len(parts) != 2 {
			ui.EditEphemeralError(s, i, ui.MsgGenericError)
			return
		}
		combatSessionID, targetID := parts[0], parts[1]
		idempotencyKey := i.ID // Dùng Discord Interaction ID làm Nonce chống double click hoàn hảo

		cSession, err := r.combatSvc.PlayerBasicAttack(ctx, userID, combatSessionID, targetID, idempotencyKey)
		if err != nil {
			if errors.Is(err, combat.ErrNotYourTurn) || errors.Is(err, combat.ErrTargetAlreadyDead) {
				ui.RespondEphemeralWarning(s, i, err.Error())
			} else {
				ui.EditEphemeralError(s, i, err.Error())
			}
			return
		}
		r.renderCombatScreen(s, i, menuSession, cSession)

	case menu.ActionPvEClaim:
		combatSessionID := extra
		idempotencyKey := i.ID

		claimed, err := r.pvecombatSvc.ClaimReward(ctx, userID, combatSessionID, idempotencyKey)
		if err != nil {
			ui.EditEphemeralError(s, i, fmt.Sprintf("Nhận thưởng thất bại: %v", err))
			return
		}

		// Render màn hình thưởng
		desc := "Đạo hữu nhận được:\n"
		for _, c := range claimed {
			prefix := ""
			if c.IsBonus {
				prefix = "[Hiếm] "
			}
			desc += fmt.Sprintf("- %s %s x%d\n", prefix, c.Type, c.Quantity)
		}
		embed := ui.SuccessEmbed("Phần Thưởng", desc)
		comps := []discordgo.MessageComponent{
			ui.ActionRow(ui.Button("Quay Lại", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSession.SessionID, string(menu.PagePvE)), ui.BtnPrimary, nil, false)),
		}
		_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})

	default:
		ui.EditEphemeralError(s, i, ui.MsgGenericError)
	}
}

func (r *Router) renderAreaSelect(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, isBiCanh bool) {
	vm := PvEMenuViewModel{SessionID: menuSession.SessionID}
	actType := pve.ActivityDuNgoan
	if isBiCanh {
		actType = pve.ActivityBiCanh
	}

	for _, def := range pve.AreaRegistry {
		if def.ActivityType == actType {
			nextStage := 1
			prog, err := r.pveRepo.GetAreaProgress(context.Background(), menuSession.UserID, def.ID)
			if err == nil && prog != nil {
				nextStage = prog.HighestStageCleared + 1
			}
			if nextStage > def.MaxStage {
				nextStage = def.MaxStage
			}

			vm.Areas = append(vm.Areas, AreaViewModel{
				ID: def.ID, Name: def.Name, NextStage: nextStage,
				RecommendCP: fmt.Sprintf("%d", 100*nextStage), // Nháp config CP
				StaminaCost: def.EntryCost.Stamina,
			})
		}
	}

	embed := BuildAreaSelectEmbed(vm, isBiCanh)
	comps := BuildAreaSelectComponents(menuSession.SessionID, vm)
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
}

func (r *Router) renderCombatScreen(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, cSession *combat.CombatSession) {
	areaName := "Không Rõ"
	if area, ok := pve.AreaRegistry[cSession.AreaID]; ok {
		areaName = area.Name
	}
	vm := CombatSessionToViewModel(cSession, areaName)
	embed := BuildCombatEmbed(vm)
	comps := BuildCombatActionComponents(menuSession.SessionID, vm)
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
}

// PageLoader tĩnh cho màn Main PvE
func PvEMainLoader(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	embed := &discordgo.MessageEmbed{
		Title:       emoji.Map.String() + " Tu Tiên Giới",
		Description: "Lựa chọn con đường chinh phạt tiếp theo của đạo hữu.",
		Color:       ui.ColorCombat,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Du Ngoạn", Value: "Tiêu hao Thể lực. Rớt kinh nghiệm, linh thạch, tài nguyên cường hóa cơ bản.", Inline: true},
			{Name: "Bí Cảnh", Value: "Tiêu hao Thể lực lớn. Rớt trang bị, pháp bảo, kỳ trân dị thảo.", Inline: true},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	comps := BuildPvEMainComponents(session.SessionID)

	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: comps,
	}, nil
}
