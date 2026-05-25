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
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// PageLoader là hàm tải toàn bộ dữ liệu cho một trang và trả về response data.
type PageLoader func(ctx context.Context, session *Session) (*discordgo.InteractionResponseData, error)

// Router xử lý mọi tương tác component (button / select menu) bên trong hệ thống menu.
type Router struct {
	cfg            *config.Config
	sessionSvc     SessionService
	cultivationSvc cultivation.Service
	inventorySvc   inventory.Service
	equipSvc       equipment.Service
	pageLoaders    map[Page]PageLoader
	log            *zap.Logger
}

// NewRouter tạo menu interaction router.
// pageLoaders ánh xạ mỗi Page sang hàm tải dữ liệu và render trang đó.
func NewRouter(cfg *config.Config, sessionSvc SessionService, cultSvc cultivation.Service, invSvc inventory.Service, equipSvc equipment.Service, loaders map[Page]PageLoader) *Router {
	return &Router{
		cfg:            cfg,
		sessionSvc:     sessionSvc,
		cultivationSvc: cultSvc,
		inventorySvc:   invSvc,
		equipSvc:       equipSvc,
		pageLoaders:    loaders,
		log:            logger.L().Named("menu.router"),
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
	case DomainInventory:
		r.handleInventoryAction(s, i, session, parsed.Action)
	default:
		r.log.Warn("menu.Router: domain không xác định",
			zap.String("domain", parsed.Domain),
			zap.String("customID", customID))
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
	}
}

func (r *Router) handleInventoryAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	ctx := context.Background()
	var err error
	var msg string

	switch action {
	case ActionInventoryUse:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			msg, err = r.inventorySvc.UseItem(ctx, session.UserID, session.GuildID, data.Values[0])
		}
	default:
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
		return
	}

	if err != nil {
		ui.RespondEphemeralWarning(s, i, apperrors.UserFacing(err, "Không thể sử dụng vật phẩm này!"))
		return
	}

	if msg != "" {
		r.renderPage(s, i, session, PageInventory)
		_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{ui.SuccessEmbed("Túi Đồ", msg)},
			Flags:  discordgo.MessageFlagsEphemeral,
		})
	}
}

// handleNav xử lý các nút Làm mới / Quay lại / Đóng.
func (r *Router) handleNav(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionClose:
		_ = r.sessionSvc.Close(ctx, session.SessionID)
		// Gửi tín hiệu phản hồi ngay để tránh Discord báo lỗi timeout
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredMessageUpdate,
		})
		// Xóa hoàn toàn tin nhắn menu để giữ kênh gọn gàng
		_ = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

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
func (r *Router) handleProfileAction(s *discordgo.Session, i *discordgo.Interaction, _ *Session, action string) {
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
	ctx := context.Background()
	in := cultivation.CultivationActionInput{UserID: session.UserID, GuildID: session.GuildID, Now: time.Now().UTC()}
	var err error
	var msg string

	switch action {
	case ActionMeditate:
		res, e := r.cultivationSvc.Meditate(ctx, in)
		err = e
		if res != nil {
			msg = res.Message
		}
	case ActionClosedDoor:
		res, e := r.cultivationSvc.Seclusion(ctx, in)
		err = e
		if res != nil {
			msg = res.Message
		}
	case ActionBodyTraining:
		res, e := r.cultivationSvc.BodyTraining(ctx, in)
		err = e
		if res != nil {
			msg = res.Message
		}
	case ActionBreakthrough:
		bIn := cultivation.BreakthroughInput{UserID: session.UserID, GuildID: session.GuildID, Now: time.Now().UTC(), Rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
		res, e := r.cultivationSvc.Breakthrough(ctx, bIn)
		err = e
		if res != nil {
			msg = res.Message
		}
	case ActionChoosePath:
		data := i.MessageComponentData()
		if len(data.Values) > 0 {
			selectedPath := cultivation.CultivationPath(data.Values[0])
			err = r.cultivationSvc.ChoosePath(ctx, session.UserID, session.GuildID, selectedPath)
			if err == nil {
				msg = fmt.Sprintf("Cảm ngộ thiên địa thành công! Đạo hữu đã chính thức bước lên con đường **%s**.", selectedPath.DisplayName())
			}
		} else {
			return
		}
	default:
		ui.RespondEphemeralError(s, i, ui.MsgComingSoon)
		return
	}

	// Xử lý báo lỗi (nếu có)
	if err != nil {
		var cdErr *apperrors.CooldownError
		switch {
		case errors.As(err, &cdErr):
			ui.RespondEphemeralWarning(s, i, fmt.Sprintf("Đạo hữu cần nghỉ ngơi. Vui lòng chờ %s.", cdErr.Remaining))
		case apperrors.IsInsufficientStamina(err):
			ui.RespondEphemeralWarning(s, i, "Thể lực không đủ! Xin hãy nghỉ ngơi.")
		case apperrors.IsInsufficientFunds(err):
			ui.RespondEphemeralWarning(s, i, "Linh thạch không đủ để chuẩn bị trận pháp đột phá!")
		case errors.Is(err, apperrors.ErrInsufficientCultivationExp):
			ui.RespondEphemeralWarning(s, i, "Tu vi chưa đạt bình cảnh, không thể miễn cưỡng đột phá.")
		case errors.Is(err, apperrors.ErrInsufficientMindState):
			ui.RespondEphemeralWarning(s, i, "Tâm cảnh quá thấp, đột phá lúc này chắc chắn tẩu hỏa nhập ma!")
		case errors.Is(err, apperrors.ErrMaxRealmReached):
			ui.RespondEphemeralWarning(s, i, "Đạo hữu đã đứng trên đỉnh phong vạn giới, không thể tiến thêm.")
		case errors.Is(err, apperrors.ErrPathAlreadyChosen):
			ui.RespondEphemeralWarning(s, i, "Đạo hữu đã có đạo lộ của riêng mình, không thể cải tu!")
		default:
			r.log.Error("Cultivation action error", zap.Error(err))
			ui.RespondEphemeralError(s, i, ui.MsgGenericError)
		}
		return
	}

	// 1. Cập nhật message gốc ngay lập tức (Real-time UI Edit)
	r.renderPage(s, i, session, PageCultivation)

	// 2. Gửi thông báo kết quả dạng popup ẩn (Followup Message)
	_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{ui.SuccessEmbed("Tu Luyện", msg)},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
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
