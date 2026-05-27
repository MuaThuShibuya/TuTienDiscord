package handlers

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

type DevHandler struct {
	cfg        *config.Config
	sessionSvc menu.SessionService
	log        *zap.Logger
}

func NewDevHandler(cfg *config.Config, sessionSvc menu.SessionService) *DevHandler {
	return &DevHandler{cfg: cfg, sessionSvc: sessionSvc, log: zap.L().Named("handler.dev")}
}

func (h *DevHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !h.cfg.IsOwner(i.Member.User.ID) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{ui.WarningEmbed("Thiên Cơ Các khép cửa. Đạo hữu không có thiên mệnh chấp chưởng pháp trận này.")},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	ctx := context.Background()
	session, err := h.sessionSvc.OpenMenu(ctx, i.Member.User.ID, i.GuildID, i.ChannelID, h.cfg.Menu.SessionTTL)
	if err != nil {
		h.log.Error("Không thể tạo session admin", zap.Error(err))
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{ui.ErrorEmbed("Linh mạch nghịch chuyển, không thể tạo pháp trận lúc này.")},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	embed := ui.SuccessEmbed("Thiên Cơ Các", "Pháp trận đã sẵn sàng. Nhấn nút bên dưới để tiến vào không gian quản trị.")
	comps := []discordgo.MessageComponent{ui.ActionRow(ui.Button("Tiến Vào Thiên Cơ Các", menu.Build(menu.DomainAdmin, menu.ActionAdminMain, session.SessionID), ui.BtnPrimary, emoji.Admin, false))}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}, Components: comps, Flags: discordgo.MessageFlagsEphemeral},
	})
}
