// File: internal/discord/menu/admin/router.go
package admin

import (
	"context"
	"fmt"

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

func (r *Router) HandleAdminInteraction(s *discordgo.Session, i *discordgo.Interaction, menuSession *menu.Session, action, extra string) {
	// BẢO MẬT: Gate keeping tuyệt đối
	if !r.CheckOwner(menuSession.UserID) {
		ui.EditEphemeralEmbed(s, i, ui.ErrorEmbed("Thiên Cơ Các khép cửa. Đạo hữu không có thiên mệnh chấp chưởng pháp trận này."))
		return
	}

	ctx := context.Background()

	switch action {
	case menu.ActionAdminMain:
		r.renderMainPanel(s, i, menuSession)

	case menu.ActionAdminMigrateDryRun:
		report, err := r.adminSvc.PreviewMigration(ctx)
		if err != nil {
			ui.EditEphemeralEmbed(s, i, ui.ErrorEmbed(fmt.Sprintf("Linh mạch nghịch chuyển: %v", err)))
			return
		}
		// Gửi preview bằng Ephemeral, không thay đổi message gốc
		ui.EditEphemeralEmbed(s, i, ui.SuccessEmbed("Phân Tích Thiên Cơ (Dry-run)", report))

	case menu.ActionAdminMigrateModal:
		// Trả về Modal form nhập chữ XAC NHAN
		modal := &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: &discordgo.InteractionResponseData{
				CustomID: menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateApply, menuSession.SessionID),
				Title:    "Xác Nhận Chuẩn Hóa Dữ Liệu",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.TextInput{
								CustomID:    "confirm_phrase",
								Label:       "Nhập: CHUAN HOA DU LIEU",
								Style:       discordgo.TextInputShort,
								Placeholder: "CHUAN HOA DU LIEU",
								Required:    true,
							},
						},
					},
				},
			},
		}
		_ = s.InteractionRespond(i, modal)

	case menu.ActionAdminMigrateApply:
		// Xử lý khi user submit Modal
		data := i.ModalSubmitData()
		phrase := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		if phrase != "CHUAN HOA DU LIEU" {
			ui.EditEphemeralEmbed(s, i, ui.WarningEmbed("Pháp ấn chưa khớp. Hành động đã bị Thiên Đạo chặn lại."))
			return
		}

		// Đã xác nhận, tiến hành ghi DB
		report, err := r.adminSvc.ApplyMigration(ctx, menuSession.UserID)
		if err != nil {
			ui.EditEphemeralEmbed(s, i, ui.ErrorEmbed(fmt.Sprintf("Linh mạch nghịch chuyển: %v", err)))
			return
		}
		ui.EditEphemeralEmbed(s, i, ui.SuccessEmbed("Pháp Trận Hoàn Tất", report))

	default:
		ui.EditEphemeralEmbed(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}

func (r *Router) renderMainPanel(s *discordgo.Session, i *discordgo.Interaction, session *menu.Session) {
	embed := &discordgo.MessageEmbed{
		Title:       emoji.Admin.String() + " Thiên Cơ Các — Owner Panel",
		Description: "Cảnh báo: Đây là pháp trận tối cao. Mọi hành động thao túng nhân quả đều sẽ được ghi chép vào Sổ Thiên Cơ.",
		Color:       ui.ColorCombat,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Trạng thái Hệ thống", Value: fmt.Sprintf("Môi trường: `%s`\nCho phép Danger Zone: `%v`", r.cfg.App.Env, r.cfg.App.AllowDangerousAdmin), Inline: false},
		},
	}

	comps := []discordgo.MessageComponent{
		ui.ActionRow(
			ui.SelectMenu(menu.Build(menu.DomainAdmin, menu.ActionAdminSelectMenu, session.SessionID), "Chọn khu vực quản lý...", []discordgo.SelectMenuOption{
				ui.SelectOption("Quản Lý Người Chơi", "player_tools", "Xem, reset, cấp vật phẩm", emoji.Profile, false),
				ui.SelectOption("Chuẩn Hóa Dữ Liệu Cũ", "legacy_tools", "Dọn dẹp data test v0.3", emoji.Migrate, false),
				ui.SelectOption("Giám Sát Combat", "combat_tools", "Cleanup session treo", emoji.Sword, false),
				ui.SelectOption("Nguy Hiểm (Danger Zone)", "danger_zone", "Wipe data, reset toàn cục", emoji.Danger, false),
			}),
		),
		// Các nút truy cập nhanh cho Legacy Migration (để dọn dữ liệu theo yêu cầu)
		ui.ActionRow(
			ui.Button("Preview Chuẩn Hóa", menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateDryRun, session.SessionID), ui.BtnPrimary, emoji.Migrate, false),
			ui.Button("Áp Dụng Chuẩn Hóa", menu.Build(menu.DomainAdmin, menu.ActionAdminMigrateModal, session.SessionID), ui.BtnDanger, emoji.Database, false),
		),
		ui.ActionRow(
			ui.Button("Quay Lại", menu.Build(menu.DomainNav, menu.ActionRefresh, session.SessionID, string(menu.PageMain)), ui.BtnSecondary, emoji.Back, false),
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
