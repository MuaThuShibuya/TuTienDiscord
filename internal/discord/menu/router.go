// File: internal/discord/menu/router.go
// Phiên bản: v0.1.1
// Mục đích: Phân luồng tương tác button/select-menu bên trong hệ thống menu game.
//           Phân tích custom_id, xác thực chủ phiên, điều hướng đến page renderer tương ứng.
// Bảo mật: Mọi tương tác đều xác thực sessionId + userId trước khi thay đổi trạng thái hoặc ghi DB.
//           Người bấm menu của người khác nhận thông báo ephemeral riêng.
// Ghi chú: Format custom_id: "<domain>:<action>:<sessionId>[:<extra>]"
//          Ví dụ: "menu:nav:abc123" hoặc "nav:back:abc123:main"

package menu

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// PageLoader là hàm tải toàn bộ dữ liệu cho một trang và trả về response data.
type PageLoader func(ctx context.Context, session *Session) (*discordgo.InteractionResponseData, error)

// Router xử lý mọi tương tác component (button / select menu) bên trong hệ thống menu.
type Router struct {
	cfg         *config.Config
	sessionSvc  SessionService
	pageLoaders map[Page]PageLoader
	log         *zap.Logger
}

// NewRouter tạo menu interaction router.
// pageLoaders ánh xạ mỗi Page sang hàm tải dữ liệu và render trang đó.
func NewRouter(cfg *config.Config, sessionSvc SessionService, loaders map[Page]PageLoader) *Router {
	return &Router{
		cfg:         cfg,
		sessionSvc:  sessionSvc,
		pageLoaders: loaders,
		log:         logger.L().Named("menu.router"),
	}
}

// Handle phân luồng một Discord component interaction đến handler phù hợp.
// Được gọi bởi discord router cấp cao nhất cho mọi MessageComponent interaction.
func (r *Router) Handle(s *discordgo.Session, i *discordgo.Interaction) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID

	// Phân tích custom_id thành domain, action, sessionId, extra
	parsed, err := Parse(customID)
	if err != nil {
		r.log.Warn("menu.Router: custom_id không hợp lệ", zap.String("customID", customID))
		ui.RespondEphemeralError(s, i, ui.MsgGenericError)
		return
	}

	ctx := context.Background()
	userID := i.Member.User.ID

	// --- Bảo mật: xác thực chủ sở hữu phiên ---
	session, err := r.sessionSvc.ValidateOwner(ctx, parsed.SessionID, userID)
	if err != nil {
		switch {
		case apperrors.IsSessionExpired(err):
			ui.RespondEphemeralError(s, i, ui.MsgSessionExpired)
		default:
			ui.RespondEphemeralError(s, i, ui.MsgNotYourMenu)
		}
		return
	}

	// Gia hạn TTL sau mỗi tương tác hợp lệ
	_ = r.sessionSvc.Refresh(ctx, parsed.SessionID, r.cfg.Menu.SessionTTL)

	switch parsed.Domain {
	case DomainNav:
		r.handleNav(s, i, session, parsed.Action, parsed.Extra)
	case DomainMenuSelect:
		r.handleMenuSelect(s, i, session, i.MessageComponentData())
	case DomainProfile:
		r.handleProfileAction(s, i, session, parsed.Action)
	case DomainCultivation:
		r.handleCultivationAction(s, i, session, parsed.Action)
	default:
		r.log.Warn("menu.Router: domain không xác định",
			zap.String("domain", parsed.Domain),
			zap.String("customID", customID))
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	}
}

// handleNav xử lý các nút Làm mới / Quay lại / Đóng.
func (r *Router) handleNav(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionClose:
		_ = r.sessionSvc.Close(ctx, session.SessionID)
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					ui.InfoEmbed("Đã Đóng", "Giao diện đã được đóng. Dùng `/menu` để mở lại."),
				},
				Components: []discordgo.MessageComponent{},
			},
		})

	case ActionRefresh:
		page := Page(extra)
		if page == "" {
			page = session.CurrentPage
		}
		r.renderPage(s, i, session, page)

	case ActionBack:
		targetPage := Page(extra)
		if targetPage == "" {
			targetPage = PageMain
		}
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, targetPage)
		session.CurrentPage = targetPage
		r.renderPage(s, i, session, targetPage)

	default:
		ui.RespondEphemeralError(s, i, ui.MsgGenericError)
	}
}

// handleMenuSelect xử lý select menu danh mục trên trang chính.
func (r *Router) handleMenuSelect(s *discordgo.Session, i *discordgo.Interaction, session *Session, data discordgo.MessageComponentInteractionData) {
	if len(data.Values) == 0 {
		return
	}
	selected := Page(data.Values[0])
	ctx := context.Background()
	_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, selected)
	session.CurrentPage = selected
	r.renderPage(s, i, session, selected)
}

// handleProfileAction phân luồng các action button thuộc trang Hồ Sơ.
func (r *Router) handleProfileAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	switch action {
	case ActionRename:
		// TODO v0.1: mở modal đổi đạo hiệu
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	default:
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	}
}

// handleCultivationAction phân luồng các action button thuộc trang Tu Luyện.
func (r *Router) handleCultivationAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	switch action {
	case ActionMeditate:
		// TODO v0.2: tĩnh tu — kiểm tra cooldown → cộng exp → cập nhật DB
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	case ActionClosedDoor:
		// TODO v0.2: bế quan
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	case ActionBreakthrough:
		// TODO v0.2: đột phá cảnh giới
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	default:
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	}
}

// renderPage gọi PageLoader tương ứng và chỉnh sửa message menu hiện tại.
func (r *Router) renderPage(s *discordgo.Session, i *discordgo.Interaction, session *Session, page Page) {
	loader, ok := r.pageLoaders[page]
	if !ok {
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseData, err := loader(ctx, session)
	if err != nil {
		r.log.Error("menu.Router: page loader thất bại",
			zap.String("page", string(page)),
			zap.String("userId", session.UserID),
			zap.Error(err),
		)
		ui.UpdateWithError(s, i, ui.MsgGenericError)
		return
	}

	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: responseData,
	})
}
