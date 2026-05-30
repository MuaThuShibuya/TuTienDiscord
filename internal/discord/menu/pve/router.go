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
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/pve"
	"github.com/whiskey/tu-tien-bot/internal/game/pvecombat"
)

type Router struct {
	pvecombatSvc *pvecombat.Service
	combatSvc    *combat.Service
	pveRepo      pve.ProgressRepository
	actionCache  *ActionCache
	log          *zap.Logger
}

func NewRouter(pvecombatSvc *pvecombat.Service, combatSvc *combat.Service, pveRepo pve.ProgressRepository, log *zap.Logger) *Router {
	return &Router{pvecombatSvc: pvecombatSvc, combatSvc: combatSvc, pveRepo: pveRepo, actionCache: NewActionCache(), log: log.Named("menu.pve")}
}

// formatCombatError dịch lỗi backend sang UI tu tiên
func formatCombatError(err error) string {
	if errors.Is(err, combat.ErrNotYourTurn) {
		return "Thiên đạo chưa xoay tới lượt đạo hữu. Hãy chờ khí cơ ổn định."
	}
	if errors.Is(err, combat.ErrTargetAlreadyDead) {
		return "Yêu vật này đã tan biến, không thể truy kích thêm."
	}
	if errors.Is(err, combat.ErrInvalidCombatStats) {
		return fmt.Sprintf("Đạo thể bất ổn, chỉ số chiến đấu không hợp lệ.\n*Debug: %v*", err)
	}
	if errors.Is(err, combat.ErrRewardAlreadyClaimed) {
		return "Đạo hữu đã gom sạch thiên tài địa bảo nơi này rồi, không còn gì sót lại."
	}
	if errors.Is(err, combat.ErrRewardClaimInProgress) {
		return "Phần thưởng đang được xử lý, vui lòng chờ."
	}
	if errors.Is(err, combat.ErrRewardClaimFailedNeedsAdmin) {
		return "Quá trình nhận thưởng gặp lỗi hệ thống và đã được khóa an toàn. Vui lòng báo Admin để kiểm tra, không bấm nhận lại nhiều lần."
	}
	if errors.Is(err, combat.ErrRewardSessionNotWon) {
		return "Yêu ma chưa dẹp yên, làm sao có thể tranh đoạt cơ duyên?"
	}
	if errors.Is(err, combat.ErrCombatSessionNotFound) || errors.Is(err, combat.ErrCombatSessionExpired) {
		return "Cơ duyên đã tàn, linh khí nơi này đã tản đi. Hãy mở lại một chuyến du ngoạn mới."
	}
	if errors.Is(err, combat.ErrRewardGrantFailed) || strings.Contains(err.Error(), "vật phẩm không tồn tại") || strings.Contains(err.Error(), "tồn tại") {
		return fmt.Sprintf("Thiên tài địa bảo nơi này chưa được Thiên Cơ Các ghi vào bảo lục. Hãy báo lại cho chưởng quản.\n*Debug: %v*", err)
	}
	if strings.Contains(err.Error(), "inventory full") || strings.Contains(err.Error(), "hết chỗ") {
		return "Túi đồ của đạo hữu đã đầy. Hãy dọn túi trước khi nhận thưởng. Phần thưởng chưa được nhận và sẽ không bị mất."
	}

	// Lỗi Generic
	return fmt.Sprintf("Linh mạch dao động, pháp trận tạm thời bất ổn. Hãy thử lại sau.\n*Debug: %v*", err)
}

func (r *Router) HandlePvEInteraction(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, action, extra string) {
	ctx := context.Background()
	userID := menuSession.UserID

	switch action {
	case menu.ActionPvEDuNgoan, menu.ActionPvEBiCanh:
		isBiCanh := action == menu.ActionPvEBiCanh
		r.renderAreaSelect(s, i, menuSession, isBiCanh)

	case menu.ActionPvEStart:
		r.log.Info("Nhận action bắt đầu PvE",
			zap.String("userId", userID),
			zap.String("customID", i.MessageComponentData().CustomID),
			zap.String("sessionID_menu", menuSession.SessionID),
		)
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			ui.EditEphemeralError(s, i, ui.MsgGenericError)
			return
		}
		areaID := data.Values[0]

		cSession, err := r.pvecombatSvc.StartPvECombat(ctx, userID, areaID)
		if err != nil {
			r.log.Error("StartPvECombat lỗi", zap.String("userId", userID), zap.Error(err))
			ui.EditEphemeralError(s, i, fmt.Sprintf("Không thể khởi tạo ải: %v", err))
			return
		}

		r.log.Info("StartPvECombat thành công, chuẩn bị render UI",
			zap.String("combatSessionID", cSession.ID),
			zap.Int("stage", cSession.Stage),
			zap.Int("enemyCount", len(cSession.Enemies)),
			zap.String("state", string(cSession.State)),
		)
		r.renderCombatScreen(s, i, menuSession, cSession)

	case menu.ActionPvEAttack:
		token := extra
		payload, ok := r.actionCache.Get(token)
		if !ok {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Cơ duyên đã tàn, pháp ấn đã mất hiệu lực. Hãy làm mới trận chiến."))
			return
		}
		if payload.OwnerID != userID {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Thiên cơ đã định, đây không phải trận chiến của đạo hữu."))
			return
		}

		r.log.Info("Xử lý lệnh Attack",
			zap.String("userId", userID),
			zap.String("sessionID", payload.CombatSessionID),
			zap.String("targetID", payload.TargetID),
		)

		cSession, err := r.combatSvc.PlayerBasicAttack(ctx, userID, payload.CombatSessionID, payload.TargetID, i.ID)
		if err != nil {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed(formatCombatError(err)))
			return
		}
		r.log.Info("Attack thành công", zap.String("stateAfter", string(cSession.State)))
		r.renderCombatScreen(s, i, menuSession, cSession)

	case menu.ActionPvESkill:
		ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Công pháp chưa khai mở, đạo hữu hãy tạm dùng kiếm pháp căn cơ."))
		return

	case menu.ActionPvEAuto:
		token := extra
		payload, ok := r.actionCache.Get(token)
		if !ok {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Cơ duyên đã tàn, pháp ấn mất hiệu lực."))
			return
		}

		opts := combat.AutoBattleOptions{
			MaxActions: 100, // Đánh liên tục cho đến khi kết thúc trận (hoặc chạm mốc an toàn 100 hiệp)
			// FIX: Không dùng interaction ID. Dùng SessionID kết hợp suffix để đảm bảo
			// dù user spam nút Auto, backend chỉ nhận 1 luồng xử lý duy nhất cho trận này.
			IdempotencyKey: fmt.Sprintf("auto_%s", payload.CombatSessionID),
			PreferSkill:    false,
		}

		res, err := r.combatSvc.ExecuteAutoBattle(ctx, userID, payload.CombatSessionID, opts)
		if err != nil {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed(formatCombatError(err)))
			return
		}

		r.log.Info("AutoBattle hoàn tất", zap.Int("actionsTaken", res.ActionsTaken), zap.String("stoppedReason", res.StoppedReason))
		r.renderCombatScreen(s, i, menuSession, res.Session)

	case menu.ActionPvEEscape:
		ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Độn thuật đang được nghiên cứu, hiện tại chỉ có tử chiến tới cùng!"))
		return

	case menu.ActionPvEClaim:
		combatSessionID := extra

		claimed, err := r.pvecombatSvc.ClaimReward(ctx, userID, combatSessionID)
		if err != nil {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed(formatCombatError(err)))
			return
		}

		r.log.Info("Nhận thưởng thành công",
			zap.String("userId", userID),
			zap.String("sessionID", combatSessionID),
			zap.Int("itemCount", len(claimed)),
			zap.Bool("rewardClaimed", true),
		)

		// Render màn hình thưởng
		desc := "Đạo hữu nhận được:\n"
		for _, c := range claimed {
			prefix := ""
			if c.IsBonus {
				prefix = "[Hiếm] "
			}

			// Phân giải hiển thị
			displayName := c.RefID
			switch c.Type {
			case "exp":
				displayName = "Tu vi"
			case "stones":
				displayName = "Linh thạch"
			default:
				if def, ok := item.GetDefinition(c.RefID); ok {
					rarityMark := ""
					if def.Rarity != "" && def.Rarity != "D" {
						rarityMark = fmt.Sprintf(" [%s]", def.Rarity)
					}
					displayName = fmt.Sprintf("%s%s", def.Name, rarityMark)
				}
			}
			desc += fmt.Sprintf("- %s %s x%d\n", prefix, displayName, c.Quantity)
		}
		embed := ui.SuccessEmbed("Phần Thưởng", desc)
		comps := []discordgo.MessageComponent{
			ui.ActionRow(ui.Button("Quay Lại", menu.Build(menu.DomainNav, menu.ActionRefresh, menuSession.SessionID, string(menu.PageMain)), ui.BtnPrimary, nil, false)),
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
	comps := BuildCombatActionComponents(r.actionCache, menuSession.SessionID, vm)

	r.log.Debug("Edit UI Combat",
		zap.String("embedType", "combat_active"),
		zap.Int("componentCount", len(comps)),
		zap.String("currentActorID", cSession.CurrentActorID),
	)

	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
	if err != nil {
		r.log.Error("Edit UI Combat thất bại (có thể custom_id quá dài)", zap.Error(err), zap.String("userId", cSession.UserID), zap.String("combatSessionID", cSession.ID))
	}
}

// PageLoader tĩnh cho màn Main PvE
func PvEMainLoader(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	embed := &discordgo.MessageEmbed{
		Title:       emoji.Map.String() + " Tu Tiên Giới",
		Description: "Lựa chọn con đường chinh phạt tiếp theo của đạo hữu.",
		Color:       ui.ColorCombat,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Du Ngoạn", Value: "Rớt kinh nghiệm, linh thạch, tài nguyên cường hóa cơ bản.", Inline: true},
			{Name: "Bí Cảnh", Value: "Rớt trang bị, pháp bảo, kỳ trân dị thảo.", Inline: true},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	comps := BuildPvEMainComponents(session.SessionID)

	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: comps,
	}, nil
}
