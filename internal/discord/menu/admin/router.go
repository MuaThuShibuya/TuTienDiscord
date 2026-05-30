// File: internal/discord/menu/admin/router.go
package admin

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
	gameadmin "github.com/whiskey/tu-tien-bot/internal/game/admin"
)

type Router struct {
	cfg      *config.Config
	adminSvc gameadmin.Service
	log      *zap.Logger
}

func NewRouter(cfg *config.Config, adminSvc gameadmin.Service, log *zap.Logger) *Router {
	return &Router{cfg: cfg, adminSvc: adminSvc, log: log.Named("menu.admin")}
}

// CheckOwner là Middleware chặn cửa Thiên Cơ Các
func (r *Router) CheckOwner(userID string) bool {
	return r.cfg.IsOwner(userID)
}

// ValidateAdminAction xác minh tính hợp lệ và bảo mật của hành động trước khi gọi service
func ValidateAdminAction(cfg *config.Config, userID, action, confirmPhrase string) error {
	if !cfg.IsOwner(userID) {
		return fmt.Errorf("Thiên Cơ Các khép cửa. Đạo hữu không có thiên mệnh chấp chưởng pháp trận này.")
	}

	// Các hành động nguy hiểm cần xác thực thêm
	isDangerousAction := strings.Contains(action, "reset") || strings.Contains(action, "wipe") || strings.Contains(action, "apply")
	if isDangerousAction && !cfg.CanExecuteDangerZone() {
		return fmt.Errorf("Thiên Đạo đã phong ấn khu vực này trong môi trường `production`. Hãy cấu hình `ALLOW_DANGEROUS_ADMIN=true` để mở khóa.")
	}

	switch action {
	case menu.ActionAdminMigrateApply:
		if confirmPhrase != "XACNHAN" {
			return fmt.Errorf("Pháp ấn chưa khớp. Thiên Đạo đã chặn hành động này.")
		}
	case menu.ActionAdminResetUserApply:
		if !strings.HasPrefix(confirmPhrase, "XACNHAN") {
			return fmt.Errorf("Pháp ấn chưa khớp. Thiên Đạo đã chặn hành động này.")
		}
	case menu.ActionAdminResetAllApply:
		if confirmPhrase != "XACNHAN" {
			return fmt.Errorf("Pháp ấn chưa khớp. Thiên Đạo đã chặn hành động này.")
		}
	case menu.ActionAdminCombatCleanApply:
		if confirmPhrase != "XACNHAN" {
			return fmt.Errorf("Pháp ấn chưa khớp. Thiên Đạo đã chặn hành động này.")
		}
	}
	return nil
}

func (r *Router) sendError(s *discordgo.Session, i *discordgo.Interaction, err error) {
	_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{ui.ErrorEmbed(fmt.Sprintf("Linh mạch nghịch chuyển: %v", err))},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

func (r *Router) sendSuccess(s *discordgo.Session, i *discordgo.Interaction, title, description string) {
	_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{ui.SuccessEmbed(title, description)},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

func (r *Router) sendWarning(s *discordgo.Session, i *discordgo.Interaction, message string) {
	_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{ui.WarningEmbed(message)},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

func getModalValue(i *discordgo.Interaction, customID string) string {
	for _, comp := range i.ModalSubmitData().Components {
		if row, ok := comp.(*discordgo.ActionsRow); ok {
			for _, inner := range row.Components {
				if textInput, ok := inner.(*discordgo.TextInput); ok && textInput.CustomID == customID {
					return textInput.Value
				}
			}
		}
	}
	return ""
}

func (r *Router) HandleAdminInteraction(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, action, extra string) {
	isModalSubmit := i.Type == discordgo.InteractionModalSubmit
	var phrase string
	if isModalSubmit {
		// Lấy confirm phrase từ modal nếu có
		phrase = getModalValue(i, "confirm_phrase")
		r.log.Debug("Xử lý Admin Modal Submit",
			zap.String("action", action),
			zap.String("userId", menuSession.UserID),
			zap.String("sessionId", menuSession.SessionID),
		)
	}

	// BẢO MẬT: Gate keeping tuyệt đối
	if err := ValidateAdminAction(r.cfg, menuSession.UserID, action, phrase); err != nil {
		r.log.Warn("Chặn truy cập Admin trái phép", zap.String("userId", menuSession.UserID), zap.Error(err))
		r.sendWarning(s, i, err.Error())
		return
	}

	if isModalSubmit {
		r.log.Debug("Gate keeping Admin Modal Submit thành công", zap.String("action", action))
	}

	ctx := context.Background()

	switch action {
	case menu.ActionAdminMain:
		r.renderMainPanel(s, i, menuSession)

	case menu.ActionAdminMigrateDryRun:
		report, err := r.adminSvc.PreviewMigration(ctx)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.sendSuccess(s, i, "Thiên Cơ Soi Chiếu (Huyễn Ảnh)", report)

	case menu.ActionAdminMigrateModal:
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateApply, menuSession.SessionID),
				Title:    "Xác Nhận Chuẩn Hóa Dữ Liệu",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "confirm_phrase",
								Label:       "Nhập: XACNHAN",
								Style:       discordgo.TextInputShort,
								Placeholder: "XACNHAN",
								Required:    true,
							},
						},
					},
				},
			},
		})

	case menu.ActionAdminMigrateApply:
		report, err := r.adminSvc.ApplyMigration(ctx, menuSession.UserID)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.sendSuccess(s, i, "Pháp Trận Vận Chuyển Hoàn Tất", report)

	case menu.ActionAdminPlayerLookupModal:
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID:   menu.Build(menu.DomainAdmin, menu.ActionAdminPlayerLookupApply, menuSession.SessionID),
				Title:      "Tra Cứu Đạo Hữu",
				Components: []discordgo.MessageComponent{ui.ActionRow(discordgo.TextInput{CustomID: "target_id", Label: "Discord ID", Style: discordgo.TextInputShort, Required: true})},
			},
		})

	case menu.ActionAdminPlayerLookupApply:
		targetID := getModalValue(i, "target_id")
		info, err := r.adminSvc.GetPlayerInfo(ctx, targetID)
		if err != nil {
			r.sendWarning(s, i, err.Error())
			return
		}
		r.sendSuccess(s, i, fmt.Sprintf("Nhân Quả: %s", targetID), info)

	case menu.ActionAdminResetUserModal:
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID:   menu.Build(menu.DomainAdmin, menu.ActionAdminResetUserPreview, menuSession.SessionID),
				Title:      "Thiên Phạt - Reset Đạo Hữu",
				Components: []discordgo.MessageComponent{ui.ActionRow(discordgo.TextInput{CustomID: "target_id", Label: "Discord ID của Đạo Hữu", Style: discordgo.TextInputShort, Required: true})},
			},
		})

	case menu.ActionAdminResetUserPreview:
		targetID := getModalValue(i, "target_id")
		opts := gameadmin.ResetOptions{Scope: gameadmin.ResetScopeUser, TargetUserID: targetID, DryRun: true}
		preview, err := r.adminSvc.PreviewReset(ctx, opts)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.renderResetPreview(s, i, menuSession, preview)

	case menu.ActionAdminResetUserApply:
		targetID := extra // Lấy từ custom_id của nút confirm
		fullPhrase := getModalValue(i, "confirm_phrase")
		expectedPhrase := fmt.Sprintf("XACNHAN %s", targetID)

		if fullPhrase != expectedPhrase {
			r.sendWarning(s, i, "Pháp ấn chưa khớp. Thiên Đạo đã chặn hành động này.")
			return
		}

		opts := gameadmin.ResetOptions{Scope: gameadmin.ResetScopeUser, TargetUserID: targetID, RequestedBy: menuSession.UserID}
		result, err := r.adminSvc.ApplyReset(ctx, opts)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.sendSuccess(s, i, "Thiên Phạt Hoàn Tất", result.Summary())

	case menu.ActionAdminResetAllPreview:
		opts := gameadmin.ResetOptions{Scope: gameadmin.ResetScopeAll, DryRun: true}
		preview, err := r.adminSvc.PreviewReset(ctx, opts)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.renderResetPreview(s, i, menuSession, preview)

	case menu.ActionAdminResetAllApply:
		opts := gameadmin.ResetOptions{Scope: gameadmin.ResetScopeAll, RequestedBy: menuSession.UserID}
		result, err := r.adminSvc.ApplyReset(ctx, opts)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.sendSuccess(s, i, "Càn Khôn Tái Lập Hoàn Tất", result.Summary())

	case menu.ActionAdminConfirmResetModal:
		target := strings.TrimSpace(extra)
		var actionConfirm, phrase string
		if target == "all" {
			actionConfirm = menu.ActionAdminResetAllApply
			phrase = "XACNHAN"
		} else {
			actionConfirm = menu.ActionAdminResetUserApply
			phrase = fmt.Sprintf("XACNHAN %s", target)
		}

		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: menu.Build(menu.DomainAdmin, actionConfirm, menuSession.SessionID, target),
				Title:    "Xác Nhận Pháp Ấn Tối Cao",
				Components: []discordgo.MessageComponent{
					ui.ActionRow(discordgo.TextInput{
						CustomID: "confirm_phrase", Label: fmt.Sprintf("Nhập chính xác: %s", phrase),
						Style: discordgo.TextInputShort, Placeholder: phrase, Required: true,
					}),
				},
			},
		})

	case menu.ActionAdminAuditLogs:
		report, err := r.adminSvc.GetRecentAudits(ctx)
		if err != nil {
			r.sendError(s, i, err)
			return
		}
		r.sendSuccess(s, i, "Sổ Thiên Cơ (Audit)", report)

	default:
		r.log.Debug("Tính năng Admin chưa mở", zap.String("action", action))
		r.sendWarning(s, i, ui.MsgComingSoon)
	}
}

func (r *Router) renderMainPanel(s *discordgo.Session, i *discordgo.Interaction, session *menu.Session) {
	embed := &discordgo.MessageEmbed{
		Title:       emoji.Admin.String() + " Thiên Cơ Các — Chủ Quyền Pháp Trận",
		Description: "Không gian quản trị tối cao. Mọi thao tác xoay chuyển càn khôn đều được Thiên Đạo ghi chép lại.",
		Color:       ui.ColorCombat,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Trạng thái Hệ thống", Value: fmt.Sprintf("Môi trường: `%s`\nDanger Zone: `%v`", r.cfg.App.Env, r.cfg.CanExecuteDangerZone()), Inline: false},
		},
	}

	comps := []discordgo.MessageComponent{
		ui.ActionRow(
			ui.Button("Tra Cứu Đạo Hữu", menu.Build(menu.DomainAdmin, menu.ActionAdminPlayerLookupModal, session.SessionID), ui.BtnPrimary, emoji.Profile, false),
			ui.Button("Sổ Thiên Cơ", menu.Build(menu.DomainAdmin, menu.ActionAdminAuditLogs, session.SessionID), ui.BtnSecondary, emoji.Info, false),
		),
		ui.ActionRow(
			ui.Button("Preview Chuẩn Hóa", menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateDryRun, session.SessionID), ui.BtnPrimary, emoji.Migrate, false),
			ui.Button("Áp Dụng Chuẩn Hóa", menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateModal, session.SessionID), ui.BtnDanger, emoji.Database, false),
		),
		ui.ActionRow(
			ui.Button("Reset Đạo Hữu", menu.Build(menu.DomainAdmin, menu.ActionAdminResetUserModal, session.SessionID), ui.BtnDanger, emoji.Profile, false),
			ui.Button("Reset Toàn Bộ", menu.Build(menu.DomainAdmin, menu.ActionAdminResetAllPreview, session.SessionID), ui.BtnDanger, emoji.Danger, false),
		),
		ui.ActionRow(
			ui.Button("Đóng Menu Admin", menu.Build(menu.DomainNav, menu.ActionClose, session.SessionID), ui.BtnSecondary, emoji.Close, false),
		),
	}

	if i.Type == discordgo.InteractionApplicationCommand {
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Components: comps, Flags: discordgo.MessageFlagsEphemeral},
		})
	} else {
		_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
	}
}

func (r *Router) renderResetPreview(s *discordgo.Session, i *discordgo.Interaction, session *menu.Session, preview *gameadmin.ResetPreview) {
	var title, confirmLabel, extra string
	if preview.Scope == gameadmin.ResetScopeUser {
		title = fmt.Sprintf("Thiên Phạt - Huyễn Ảnh Reset Đạo Hữu %s", preview.TargetUserID)
		confirmLabel = "Xác Nhận Thiên Phạt"
		extra = preview.TargetUserID
	} else {
		title = "Càn Khôn Tái Lập - Huyễn Ảnh Reset Toàn Bộ"
		confirmLabel = "Xác Nhận Tái Lập"
		extra = "all"
	}

	embed := &discordgo.MessageEmbed{
		Title:       emoji.Danger.String() + " " + title,
		Description: preview.Summary() + "\n\n**CẢNH BÁO: HÀNH ĐỘNG NÀY KHÔNG THỂ HOÀN TÁC!**",
		Color:       0xED4245,
	}

	confirmButton := ui.Button(confirmLabel, menu.Build(menu.DomainAdmin, menu.ActionAdminConfirmResetModal, session.SessionID, extra), ui.BtnDanger, emoji.Check, false)
	cancelButton := ui.Button("Hủy", menu.Build(menu.DomainAdmin, menu.ActionAdminMain, session.SessionID), ui.BtnSecondary, emoji.Cross, false)

	// Hiển thị preview cùng với nút mở Modal
	comps := []discordgo.MessageComponent{ui.ActionRow(confirmButton, cancelButton)}
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
}
