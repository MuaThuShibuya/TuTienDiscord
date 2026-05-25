// File: internal/discord/router.go
// Phiên bản: v0.1.1
// Mục đích: Router tương tác Discord cấp cao nhất — phân luồng slash command và component interaction
//           đến đúng handler. Đây là điểm vào duy nhất cho mọi Discord event.
// Bảo mật: Chỉ xử lý interaction từ guild (bỏ qua DM). Lệnh không xác định được bỏ qua graceful.
// Ghi chú: Thêm route lệnh mới ở đây khi có tính năng mới. File này chỉ phân luồng, không có logic.

package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/discord/handlers"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Router chứa tất cả handler đã đăng ký và phân luồng interaction đến chúng.
type Router struct {
	startHandler *handlers.StartHandler
	menuHandler  *handlers.MenuHandler
	menuRouter   *menu.Router
	log          *zap.Logger
}

// NewRouter tạo top-level interaction router.
func NewRouter(
	startHandler *handlers.StartHandler,
	menuHandler *handlers.MenuHandler,
	menuRouter *menu.Router,
) *Router {
	return &Router{
		startHandler: startHandler,
		menuHandler:  menuHandler,
		menuRouter:   menuRouter,
		log:          logger.L().Named("discord.router"),
	}
}

// HandleInteraction là discordgo event handler được đăng ký trên session.
// Phân luồng mỗi interaction đến handler phù hợp dựa trên type và tên.
func (r *Router) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Bỏ qua tất cả interaction không từ guild member (DM)
	if i.Member == nil {
		r.log.Debug("Bỏ qua DM interaction")
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		r.routeCommand(s, i)

	case discordgo.InteractionMessageComponent:
		// Tất cả component interaction đi qua menu router
		r.menuRouter.Handle(s, i.Interaction)

	case discordgo.InteractionModalSubmit:
		// TODO v0.1+: phân luồng modal submit (ví dụ: modal đổi đạo hiệu)
		r.log.Debug("Nhận modal submit (chưa xử lý)",
			zap.String("customID", i.ModalSubmitData().CustomID))

	default:
		r.log.Debug("Loại interaction không xác định", zap.Int("type", int(i.Type)))
	}
}

// routeCommand phân luồng slash command đến handler tương ứng.
func (r *Router) routeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	r.log.Debug("Nhận slash command",
		zap.String("command", name),
		zap.String("userId", i.Member.User.ID),
		zap.String("guildId", i.GuildID),
	)

	switch name {
	case "start":
		r.startHandler.Handle(s, i)
	case "menu":
		r.menuHandler.Handle(s, i)
	default:
		r.log.Warn("Lệnh không xác định", zap.String("command", name))
	}
}
