// File: internal/discord/handlers/menu_handler.go
// Phiên bản: v0.1.1
// Mục đích: Controller cho lệnh /menu — mở giao diện game tổng hợp.
// Bảo mật: Xác minh người chơi đã đăng ký trước khi tạo phiên menu.
//           Session dùng sessionId sinh ngẫu nhiên (crypto/rand). Mọi tương tác sau đó
//           đều được xác thực qua menu.Router.
// Ghi chú: /menu luôn mở trang Main. Điều hướng giữa các trang xảy ra qua menu.Router.
//           PageLoader được inject để handler không trực tiếp gọi DB.

package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/equipment"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/profile"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

// MenuHandler xử lý lệnh /menu và cung cấp PageLoader cho menu router.
type MenuHandler struct {
	cfg            *config.Config
	profileSvc     profile.Service
	cultivationSvc cultivation.Service
	economySvc     economy.Service
	inventorySvc   inventory.Service
	equipSvc       equipment.Service
	sessionSvc     menu.SessionService
	log            *zap.Logger
}

// NewMenuHandler tạo MenuHandler với các service đã inject.
func NewMenuHandler(
	cfg *config.Config,
	profileSvc profile.Service,
	cultivationSvc cultivation.Service,
	economySvc economy.Service,
	inventorySvc inventory.Service,
	equipSvc equipment.Service,
	sessionSvc menu.SessionService,
) *MenuHandler {
	return &MenuHandler{
		cfg:            cfg,
		profileSvc:     profileSvc,
		cultivationSvc: cultivationSvc,
		economySvc:     economySvc,
		inventorySvc:   inventorySvc,
		equipSvc:       equipSvc,
		sessionSvc:     sessionSvc,
		log:            logger.L().Named("handler.menu"),
	}
}

// Handle xử lý lệnh slash /menu.
func (h *MenuHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		ui.RespondEphemeralError(s, i.Interaction, "Lệnh này chỉ dùng được trong server Discord.")
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID
	channelID := i.ChannelID

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.log.Debug("/menu được gọi", zap.String("userId", userID), zap.String("guildId", guildID))

	// 1. Kiểm tra người chơi đã đăng ký chưa
	player, err := h.profileSvc.GetPlayer(ctx, userID, guildID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			ui.RespondEphemeralError(s, i.Interaction, ui.MsgNotRegistered)
			return
		}
		h.log.Error("/menu: GetPlayer thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 2. Tải hồ sơ tu luyện và ví
	cult, err := h.cultivationSvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate cultivation thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	wallet, err := h.economySvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate wallet thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 3. Tạo phiên menu mới
	session, err := h.sessionSvc.OpenMenu(ctx, userID, guildID, channelID, h.cfg.Menu.SessionTTL)
	if err != nil {
		h.log.Error("/menu: OpenMenu thất bại", zap.String("userId", userID), zap.Error(err))
		ui.RespondEphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 4. Map domain models → ViewModel → gọi UI Builder
	vm := toMainMenuVM(session, player, cult, wallet)
	responseData := menu.BuildMainMenuResponse(vm)

	// 5. Gửi response (tạo message mới)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	})
	if err != nil {
		h.log.Error("/menu: InteractionRespond thất bại", zap.String("userId", userID), zap.Error(err))
		return
	}

	// 6. Lưu message ID để có thể chỉnh sửa sau
	msg, err := s.InteractionResponse(i.Interaction)
	if err == nil && msg != nil {
		_ = h.sessionSvc.SetMessageID(ctx, session.SessionID, msg.ID)
	}

	// 7. Cập nhật lần cuối hoạt động
	h.profileSvc.TouchLastActive(ctx, userID, guildID)
}

// PageLoaders trả về map page → loader function để menu.Router dùng khi điều hướng.
func (h *MenuHandler) PageLoaders() map[menu.Page]menu.PageLoader {
	return map[menu.Page]menu.PageLoader{
		menu.PageMain:        h.loadMainPage,
		menu.PageProfile:     h.loadProfilePage,
		menu.PageCultivation: h.loadCultivationPage,
		menu.PageInventory:   h.loadInventoryPage,
		// TODO v0.3+: thêm PageInventory, PageSkills, PagePets, PageGacha, PageMarket, PageSect
	}
}

func (h *MenuHandler) loadInventoryPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, err
	}

	_, items, err := h.inventorySvc.GetInventory(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, err
	}

	vm := toInventoryMenuVM(session, player, items)
	return menu.BuildInventoryMenuResponse(vm), nil
}

// loadMainPage tải dữ liệu và render trang Main Menu.
func (h *MenuHandler) loadMainPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadMainPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadMainPage cultivation: %w", err)
	}
	wallet, err := h.economySvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadMainPage wallet: %w", err)
	}
	return menu.BuildMainMenuEdit(toMainMenuVM(session, player, cult, wallet)), nil
}

// loadProfilePage tải dữ liệu và render trang Hồ Sơ.
func (h *MenuHandler) loadProfilePage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadProfilePage profile: %w", err)
	}
	wallet, err := h.economySvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadProfilePage wallet: %w", err)
	}
	return menu.BuildProfileMenuResponse(toProfileMenuVM(session, player, wallet)), nil
}

// loadCultivationPage tải dữ liệu và render trang Tu Luyện.
func (h *MenuHandler) loadCultivationPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadCultivationPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadCultivationPage cultivation: %w", err)
	}
	return menu.BuildCultivationMenuResponse(toCultivationMenuVM(session, player, cult)), nil
}

// --- ViewModel mapping functions ---
// Tất cả logic format dữ liệu (số, progress bar, timestamp) nằm ở đây,
// UI Builder chỉ nhận chuỗi đã format sẵn.

func toMainMenuVM(session *menu.Session, player *profile.Player, cult *cultivation.CultivationProfile, wallet *economy.Wallet) *menu.MainMenuVM {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	tip := ui.DailyTips[r.Intn(len(ui.DailyTips))]

	staminaBar := fmt.Sprintf("`%s` %d/%d",
		utils.ProgressBar(cult.Stamina, cult.MaxStamina, 10),
		cult.Stamina, cult.MaxStamina)

	expBar := fmt.Sprintf("`%s` %s/%s",
		utils.ProgressBar(int(cult.CultivationExp), int(cult.CultivationExpRequired), 10),
		utils.FormatNumber(cult.CultivationExp),
		utils.FormatNumber(cult.CultivationExpRequired))

	return &menu.MainMenuVM{
		SessionID:    session.SessionID,
		DaoName:      player.DaoName,
		RealmDisplay: fmt.Sprintf("%s tầng %d", cult.Realm.DisplayName(), cult.RealmLevel),
		CombatPower:  utils.FormatNumber(cult.CombatPower),
		MindState:    fmt.Sprintf("%s (%d/100)", cult.MindStateDisplayName(), cult.MindState),
		PathDisplay:  cult.Path.DisplayName(),
		StaminaBar:   staminaBar,
		ExpBar:       expBar,
		SpiritStones: utils.FormatNumber(wallet.SpiritStones),
		SpiritJades:  utils.FormatNumber(wallet.SpiritJades),
		FateTickets:  fmt.Sprintf("%d vé", wallet.FateTickets),
		DailyTip:     tip,
	}
}

func toProfileMenuVM(session *menu.Session, player *profile.Player, wallet *economy.Wallet) *menu.ProfileMenuVM {
	return &menu.ProfileMenuVM{
		SessionID:    session.SessionID,
		DaoName:      player.DaoName,
		DisplayName:  player.DisplayName,
		JoinedAt:     utils.DiscordTimestamp(player.CreatedAt, "D"),
		LastActive:   utils.DiscordTimestamp(player.LastActiveAt, "R"),
		SpiritStones: utils.FormatNumber(wallet.SpiritStones),
		SpiritJades:  utils.FormatNumber(wallet.SpiritJades),
		FateTickets:  fmt.Sprintf("%d vé", wallet.FateTickets),
	}
}

func toCultivationMenuVM(session *menu.Session, player *profile.Player, cult *cultivation.CultivationProfile) *menu.CultivationMenuVM {
	staminaBar := fmt.Sprintf("`%s` %d/%d",
		utils.ProgressBar(cult.Stamina, cult.MaxStamina, 10),
		cult.Stamina, cult.MaxStamina)

	expBar := fmt.Sprintf("`%s`\n%s / %s tu vi",
		utils.ProgressBar(int(cult.CultivationExp), int(cult.CultivationExpRequired), 12),
		utils.FormatNumber(cult.CultivationExp),
		utils.FormatNumber(cult.CultivationExpRequired))

	return &menu.CultivationMenuVM{
		SessionID:       session.SessionID,
		DaoName:         player.DaoName,
		RealmDisplay:    fmt.Sprintf("%s tầng %d", cult.Realm.DisplayName(), cult.RealmLevel),
		MindState:       fmt.Sprintf("%s (%d/100)", cult.MindStateDisplayName(), cult.MindState),
		PathDisplay:     cult.Path.DisplayName(),
		HasPath:         cult.Path != cultivation.PathNone,
		StaminaBar:      staminaBar,
		ExpBar:          expBar,
		CombatPower:     utils.FormatNumber(cult.CombatPower),
		CanBreakthrough: cult.CanBreakthrough(),
	}
}

func toInventoryMenuVM(session *menu.Session, player *profile.Player, items []*item.ItemInstance) *menu.InventoryMenuVM {
	// Logic Phân trang (Pagination) bảo vệ giới hạn 25 của Discord
	const itemsPerPage = 20
	page := 1 // Mặc định trang 1 (Có thể mở rộng lấy từ session.Data hoặc URL router sau này)

	totalItems := len(items)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	start := (page - 1) * itemsPerPage
	end := start + itemsPerPage
	if start >= totalItems {
		start = 0
	}
	if end > totalItems {
		end = totalItems
	}

	var itemVMs []menu.InventoryItemVM
	for _, it := range items[start:end] {
		name := it.DefinitionID
		// Cố gắng lấy tên thật của vật phẩm từ từ điển (nếu có)
		if def, ok := item.GetDefinition(it.DefinitionID); ok {
			name = def.Name
		}

		itemVMs = append(itemVMs, menu.InventoryItemVM{
			InstanceID: it.InstanceID,
			Name:       name,
			Quantity:   it.Quantity,
		})
	}

	return &menu.InventoryMenuVM{
		SessionID: session.SessionID,
		DaoName:   player.DaoName,
		// Đã xóa SlotLimit vì view model hiện tại không định nghĩa trường này
		Items: itemVMs,
		// Cần thêm 2 trường này vào menu.InventoryMenuVM của bạn bên gói menu
		// CurrentPage: page,
		// TotalPages:  totalPages,
	}
}
