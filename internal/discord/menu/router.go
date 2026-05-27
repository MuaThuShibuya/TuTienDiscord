// File: internal/discord/menu/router.go
// Chức năng: Phân luồng tương tác component (button/select-menu) đến handler đúng.
//            Xác thực session, điều hướng trang, xử lý hành động tu luyện/túi đồ/trang bị.
// Bảo mật: Mọi tương tác đều được ACK ngay (defer) để tránh timeout 3 giây Discord.
//          Sau khi defer, dùng InteractionResponseEdit và FollowupMessageCreate cho response.
//          ValidateOwner xác thực chủ sở hữu phiên trước khi xử lý bất kỳ hành động nào.
// Ghi chú: Format custom_id: "<domain>:<action>:<sessionId>[:<extra>]"

package menu

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/alchemy"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// PageLoader là hàm tải toàn bộ dữ liệu cho một trang và trả về response data.
type PageLoader func(ctx context.Context, session *Session) (*discordgo.InteractionResponseData, error)

// Router xử lý mọi tương tác component bên trong hệ thống menu.
type Router struct {
	cfg            *config.Config
	sessionSvc     SessionService
	cultivationSvc cultivation.Service
	inventorySvc   inventory.Service
	equipSvc       equipment.Service
	alchemySvc     alchemy.Service
	pageLoaders    map[Page]PageLoader
	log            *zap.Logger
}

// NewRouter tạo menu interaction router.
func NewRouter(
	cfg *config.Config,
	sessionSvc SessionService,
	cultSvc cultivation.Service,
	invSvc inventory.Service,
	equipSvc equipment.Service,
	alchemySvc alchemy.Service,
	loaders map[Page]PageLoader,
) *Router {
	return &Router{
		cfg:            cfg,
		sessionSvc:     sessionSvc,
		cultivationSvc: cultSvc,
		inventorySvc:   invSvc,
		equipSvc:       equipSvc,
		alchemySvc:     alchemySvc,
		pageLoaders:    loaders,
		log:            logger.L().Named("menu.router"),
	}
}

// Handle phân luồng một Discord component interaction đến handler phù hợp.
// Luôn defer trước để tránh timeout 3 giây, sau đó xử lý và edit response.
func (r *Router) Handle(s *discordgo.Session, i *discordgo.Interaction) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID

	// Phân tích custom_id — nếu lỗi format, trả về lỗi ngay (chưa defer, dùng Respond thường)
	parsed, err := Parse(customID)
	if err != nil {
		r.log.Warn("custom_id không hợp lệ",
			zap.String("customID", customID),
			zap.String("userId", i.Member.User.ID),
		)
		_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{ui.ErrorEmbed(ui.MsgGenericError)},
				Flags:  discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// ACK ngay để Discord không báo "Tương tác này không thành công".
	// DeferredMessageUpdate = ACK nhưng không thay đổi message ngay, chờ edit sau.
	if ackErr := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	}); ackErr != nil {
		r.log.Warn("Không thể ACK interaction", zap.Error(ackErr))
		return
	}

	// Từ đây dùng InteractionResponseEdit (cập nhật menu) và FollowupMessageCreate (ephemeral)
	ctx := context.Background()
	userID := i.Member.User.ID

	// Xác thực chủ sở hữu phiên
	session, err := r.sessionSvc.ValidateOwner(ctx, parsed.SessionID, userID)
	if err != nil {
		switch {
		case apperrors.IsSessionExpired(err):
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgSessionExpired))
		default:
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgNotYourMenu))
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
		r.handleInventoryAction(s, i, session, parsed.Action, parsed.Extra)
	case DomainEquipment:
		r.handleEquipmentAction(s, i, session, parsed.Action, parsed.Extra)
	case DomainAlchemy:
		r.handleAlchemyAction(s, i, session, parsed.Action, parsed.Extra)
	default:
		r.log.Warn("domain không xác định",
			zap.String("domain", parsed.Domain),
			zap.String("customID", customID),
			zap.String("userId", userID),
		)
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}

// handleAlchemyAction xử lý các action thuộc trang Lò Đan.
func (r *Router) handleAlchemyAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionAlchemyView:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		recipeID := data.Values[0]
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, PageAlchemy, recipeID)
		session.CurrentCategory = recipeID
		r.renderPage(s, i, session, PageAlchemy)

	case ActionAlchemyCancel:
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, PageAlchemy, "")
		session.CurrentCategory = ""
		r.renderPage(s, i, session, PageAlchemy)

	case ActionAlchemyCraft:
		recipeID := extra
		if recipeID == "" {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}

		// Rand source an toàn cục bộ
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

		res, err := r.alchemySvc.Craft(ctx, session.UserID, session.GuildID, recipeID, rnd)
		if err != nil {
			r.log.Error("Craft thất bại",
				zap.String("userId", session.UserID),
				zap.String("guildId", session.GuildID),
				zap.String("recipeId", recipeID),
				zap.String("step", "alchemySvc.Craft"),
				zap.Error(err),
			)
			r.sendEphemeral(s, i, ui.ErrorEmbed(apperrors.UserFacing(err, "Không đủ nguyên liệu hoặc không thể luyện đan lúc này!")))
			return
		}

		r.renderPage(s, i, session, PageAlchemy)
		if res.Success {
			r.sendEphemeral(s, i, ui.SuccessEmbed("Luyện Đan", res.Message))
		} else {
			r.sendEphemeral(s, i, ui.WarningEmbed(res.Message))
		}
	default:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}

// --- Helpers response (phải dùng sau khi đã defer) ---

// sendEphemeral gửi followup message ẩn chỉ user thấy.
func (r *Router) sendEphemeral(s *discordgo.Session, i *discordgo.Interaction, embed *discordgo.MessageEmbed) {
	_, _ = s.FollowupMessageCreate(i, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}

// renderPage gọi PageLoader tương ứng rồi cập nhật message menu.
// Phải gọi sau khi đã defer interaction.
func (r *Router) renderPage(s *discordgo.Session, i *discordgo.Interaction, session *Session, page Page) {
	loader, ok := r.pageLoaders[page]
	if !ok {
		// Trang chưa có loader → hiển thị "đang phát triển" dạng popup, menu giữ nguyên
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseData, err := loader(ctx, session)
	if err != nil {
		r.log.Error("page loader thất bại",
			zap.String("page", string(page)),
			zap.String("userId", session.UserID),
			zap.String("guildId", session.GuildID),
			zap.Error(err),
		)
		// Cập nhật menu thành thông báo lỗi, xóa components để tránh spam
		emptyComps := []discordgo.MessageComponent{}
		_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
			Embeds:     &[]*discordgo.MessageEmbed{ui.ErrorEmbed(ui.MsgGenericError)},
			Components: &emptyComps,
		})
		return
	}

	comps := responseData.Components
	_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Embeds:     &responseData.Embeds,
		Components: &comps,
	})
}

// --- Domain Handlers ---

// handleNav xử lý các nút Làm mới / Quay lại / Đóng.
func (r *Router) handleNav(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionClose:
		_ = r.sessionSvc.Close(ctx, session.SessionID)
		// Xóa message menu, loading state tự biến mất khi message không còn tồn tại
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
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, targetPage, "")
		session.CurrentPage = targetPage
		r.renderPage(s, i, session, targetPage)

	default:
		r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
	}
}

// handleMenuSelect xử lý select menu điều hướng trên trang chính.
func (r *Router) handleMenuSelect(s *discordgo.Session, i *discordgo.Interaction, session *Session, data discordgo.MessageComponentInteractionData) {
	if len(data.Values) == 0 {
		r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
		return
	}

	selected := Page(data.Values[0])

	// Check xem màn hình đích đã được code chưa, nếu chưa thì báo "Đang phát triển" và ngắt luôn
	if _, ok := r.pageLoaders[selected]; !ok {
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
		return
	}

	ctx := context.Background()
	_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, selected, "")
	session.CurrentPage = selected
	r.renderPage(s, i, session, selected)
}

// handleProfileAction xử lý các action thuộc trang Hồ Sơ.
func (r *Router) handleProfileAction(s *discordgo.Session, i *discordgo.Interaction, _ *Session, action string) {
	switch action {
	case ActionRename:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	default:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}

// handleCultivationAction xử lý các action thuộc trang Tu Luyện.
func (r *Router) handleCultivationAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	ctx := context.Background()
	in := cultivation.CultivationActionInput{
		UserID:  session.UserID,
		GuildID: session.GuildID,
		Now:     time.Now().UTC(),
	}

	var (
		err error
		msg string
	)

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
		bIn := cultivation.BreakthroughInput{
			UserID:  session.UserID,
			GuildID: session.GuildID,
			Now:     time.Now().UTC(),
			Rand:    rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		}
		res, e := r.cultivationSvc.Breakthrough(ctx, bIn)
		err = e
		if res != nil {
			msg = res.Message
		}

	case ActionChoosePath:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			return
		}
		selectedPath := cultivation.CultivationPath(data.Values[0])
		err = r.cultivationSvc.ChoosePath(ctx, session.UserID, session.GuildID, selectedPath)
		if err == nil {
			msg = fmt.Sprintf("Cảm ngộ thiên địa thành công! Đạo hữu đã bước lên con đường **%s**.", selectedPath.DisplayName())
		}

	default:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
		return
	}

	if err != nil {
		var cdErr *apperrors.CooldownError
		switch {
		case errors.As(err, &cdErr):
			r.sendEphemeral(s, i, ui.WarningEmbed(
				fmt.Sprintf("Đạo hữu cần nghỉ ngơi. Vui lòng chờ **%s**.", cdErr.Remaining)))
		case apperrors.IsInsufficientStamina(err):
			r.sendEphemeral(s, i, ui.WarningEmbed("Thể lực không đủ! Hãy nghỉ ngơi để hồi phục."))
		case apperrors.IsInsufficientFunds(err):
			r.sendEphemeral(s, i, ui.WarningEmbed("Linh thạch không đủ để chuẩn bị trận pháp đột phá!"))
		case errors.Is(err, apperrors.ErrInsufficientCultivationExp):
			r.sendEphemeral(s, i, ui.WarningEmbed("Tu vi chưa đạt bình cảnh, không thể miễn cưỡng đột phá."))
		case errors.Is(err, apperrors.ErrInsufficientMindState):
			r.sendEphemeral(s, i, ui.WarningEmbed("Tâm cảnh quá thấp, đột phá lúc này chắc chắn tẩu hỏa nhập ma!"))
		case errors.Is(err, apperrors.ErrMaxRealmReached):
			r.sendEphemeral(s, i, ui.WarningEmbed("Đạo hữu đã đứng trên đỉnh phong vạn giới, không thể tiến thêm."))
		case errors.Is(err, apperrors.ErrPathAlreadyChosen):
			r.sendEphemeral(s, i, ui.WarningEmbed("Đạo hữu đã có đạo lộ của riêng mình, không thể cải tu!"))
		default:
			r.log.Error("cultivation action lỗi",
				zap.String("action", action),
				zap.String("userId", session.UserID),
				zap.Error(err),
			)
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
		}
		return
	}

	// Cập nhật menu tu luyện (hiển thị số liệu mới)
	r.renderPage(s, i, session, PageCultivation)

	// Gửi popup kết quả ẩn (chỉ user thấy)
	if msg != "" {
		r.sendEphemeral(s, i, ui.SuccessEmbed("Tu Luyện", msg))
	}
}

// handleInventoryAction xử lý các action thuộc trang Túi Đồ.
func (r *Router) handleInventoryAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionInventoryPage:
		// extra là số trang mới (string)
		if extra == "" {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, session.CurrentPage, extra)
		session.CurrentCategory = extra
		r.renderPage(s, i, session, PageInventory)

	case ActionInventoryUse:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		instanceID := data.Values[0]
		msg, err := r.inventorySvc.UseItem(ctx, session.UserID, session.GuildID, instanceID)
		if err != nil {
			r.log.Warn("UseItem thất bại",
				zap.String("instanceId", instanceID),
				zap.String("userId", session.UserID),
				zap.Error(err),
			)
			r.sendEphemeral(s, i, ui.WarningEmbed(apperrors.UserFacing(err, "Không thể sử dụng vật phẩm này!")))
			return
		}
		// Làm mới túi đồ sau khi dùng item
		r.renderPage(s, i, session, PageInventory)
		if msg != "" {
			r.sendEphemeral(s, i, ui.SuccessEmbed("Túi Đồ", msg))
		}

	case ActionInventoryDismantle:
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		instanceID := data.Values[0]

		// 1. Chặn phân giải nếu đang được mặc trên người
		eqSet, err := r.equipSvc.GetEquipment(ctx, session.UserID, session.GuildID)
		if err == nil && eqSet != nil {
			for _, eqInstID := range eqSet.Slots {
				if eqInstID == instanceID {
					r.sendEphemeral(s, i, ui.WarningEmbed("Trang bị đang được mặc, không thể phân giải! Hãy tháo ra trước."))
					return
				}
			}
		}

		// 2. Thay thế menu hiện tại bằng cảnh báo Confirm Nội tuyến (Tuyệt đỉnh UI/UX)
		embed := ui.WarningEmbed("Bạn có chắc chắn muốn phân giải vật phẩm này? Quá trình không thể hoàn tác.")
		comps := []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: Build(DomainInventory, ActionInventoryDismantleConfirm, session.SessionID, instanceID),
						Label:    "Xác nhận Phân giải",
						Style:    discordgo.DangerButton,
					},
					discordgo.Button{
						CustomID: Build(DomainNav, ActionRefresh, session.SessionID, string(PageInventory)),
						Label:    "Hủy",
						Style:    discordgo.SecondaryButton,
					},
				},
			},
		}
		_, _ = s.InteractionResponseEdit(i, &discordgo.WebhookEdit{Embeds: &[]*discordgo.MessageEmbed{embed}, Components: &comps})
		return

	case ActionInventoryDismantleConfirm:
		instanceID := extra
		msg, err := r.inventorySvc.DismantleItem(ctx, session.UserID, session.GuildID, instanceID)
		if err != nil {
			r.sendEphemeral(s, i, ui.WarningEmbed(apperrors.UserFacing(err, "Không thể phân giải vật phẩm lúc này!")))
			return
		}
		r.renderPage(s, i, session, PageInventory)
		r.sendEphemeral(s, i, ui.SuccessEmbed("Phân Giải", msg))

	default:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}

// handleEquipmentAction xử lý các action thuộc trang Trang Bị.
func (r *Router) handleEquipmentAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case ActionEquipmentEquip:
		// Giá trị select menu: "instanceID:definitionID"
		data := i.MessageComponentData()
		if len(data.Values) == 0 {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		parts := strings.SplitN(data.Values[0], ":", 2)
		if len(parts) != 2 {
			r.log.Warn("equip value không hợp lệ", zap.String("value", data.Values[0]))
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}
		instanceID, definitionID := parts[0], parts[1]

		slot := equipment.GetSlotForDefinition(definitionID)
		if slot == "" {
			r.sendEphemeral(s, i, ui.WarningEmbed("Vật phẩm này không thể mặc được!"))
			return
		}

		if err := r.equipSvc.Equip(ctx, session.UserID, session.GuildID, slot, instanceID); err != nil {
			r.log.Error("Equip thất bại",
				zap.String("userId", session.UserID),
				zap.String("instanceId", instanceID),
				zap.String("slot", string(slot)),
				zap.Error(err),
			)
			r.sendEphemeral(s, i, ui.ErrorEmbed(apperrors.UserFacing(err, "Không thể mặc trang bị lúc này!")))
			return
		}

		r.renderPage(s, i, session, PageEquipment)
		r.sendEphemeral(s, i, ui.SuccessEmbed("Trang Bị", "Mặc trang bị thành công!"))

	case ActionEquipmentUnequip:
		// extra chứa tên slot cần tháo ("weapon", "armor", ...)
		slot := equipment.EquipmentSlot(extra)
		if slot == "" {
			r.sendEphemeral(s, i, ui.ErrorEmbed(ui.MsgGenericError))
			return
		}

		if err := r.equipSvc.Unequip(ctx, session.UserID, session.GuildID, slot); err != nil {
			r.log.Error("Unequip thất bại",
				zap.String("userId", session.UserID),
				zap.String("slot", string(slot)),
				zap.Error(err),
			)
			r.sendEphemeral(s, i, ui.ErrorEmbed(apperrors.UserFacing(err, "Không thể tháo trang bị lúc này!")))
			return
		}

		r.renderPage(s, i, session, PageEquipment)
		r.sendEphemeral(s, i, ui.SuccessEmbed("Trang Bị", "Tháo trang bị thành công!"))

	case ActionEquipmentEnhance:
		slot := equipment.EquipmentSlot(extra)
		if err := r.equipSvc.Enhance(ctx, session.UserID, session.GuildID, slot); err != nil {
			r.sendEphemeral(s, i, ui.WarningEmbed(apperrors.UserFacing(err, err.Error())))
			return
		}
		r.renderPage(s, i, session, PageEquipment)
		r.sendEphemeral(s, i, ui.SuccessEmbed("Cường Hóa", "Chúc mừng! Trang bị đã được thăng cấp."))

	default:
		r.sendEphemeral(s, i, ui.WarningEmbed(ui.MsgComingSoon))
	}
}
