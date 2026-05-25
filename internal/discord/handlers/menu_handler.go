// File: internal/discord/handlers/menu_handler.go
// Version: v0.1
// Purpose: Controller for the /menu slash command — opens or re-opens the all-in-one game menu.
// Security: Verifies player is registered before opening menu. Menu session uses cryptographic sessionId.
//           All subsequent button/select interactions are validated through menu.Router.
// Notes: /menu always opens the main page. Navigation between pages happens via menu.Router.
//        PageLoaders are injected so this handler stays decoupled from DB access.

package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/discord/menu"
	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/game/cultivation"
	"github.com/yourname/tu-tien-bot/internal/game/economy"
	"github.com/yourname/tu-tien-bot/internal/game/profile"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// MenuHandler handles the /menu command and provides PageLoaders for the menu router.
type MenuHandler struct {
	cfg            *config.Config
	profileSvc     profile.Service
	cultivationSvc cultivation.Service
	economySvc     economy.Service
	sessionSvc     menu.Service
	log            *zap.Logger
}

// NewMenuHandler creates a new MenuHandler.
func NewMenuHandler(
	cfg *config.Config,
	profileSvc profile.Service,
	cultivationSvc cultivation.Service,
	economySvc economy.Service,
	sessionSvc menu.Service,
) *MenuHandler {
	return &MenuHandler{
		cfg:            cfg,
		profileSvc:     profileSvc,
		cultivationSvc: cultivationSvc,
		economySvc:     economySvc,
		sessionSvc:     sessionSvc,
		log:            logger.L().Named("handler.menu"),
	}
}

// Handle processes the /menu slash command.
func (h *MenuHandler) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil || i.GuildID == "" {
		ui.EphemeralError(s, i.Interaction, "Lệnh này chỉ dùng được trong server Discord.")
		return
	}

	userID := i.Member.User.ID
	guildID := i.GuildID
	channelID := i.ChannelID

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	h.log.Debug("/menu invoked", zap.String("userId", userID), zap.String("guildId", guildID))

	// 1. Ensure player exists
	player, err := h.profileSvc.GetPlayer(ctx, userID, guildID)
	if err != nil {
		if apperrors.IsNotFound(err) {
			ui.EphemeralError(s, i.Interaction, ui.MsgNotRegistered)
			return
		}
		h.log.Error("/menu: GetPlayer failed",
			zap.String("userId", userID), zap.Error(err))
		ui.EphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 2. Load cultivation + wallet
	cult, err := h.cultivationSvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate cultivation failed",
			zap.String("userId", userID), zap.Error(err))
		ui.EphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	wallet, err := h.economySvc.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		h.log.Error("/menu: GetOrCreate wallet failed",
			zap.String("userId", userID), zap.Error(err))
		ui.EphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 3. Open a new menu session
	session, err := h.sessionSvc.OpenMenu(ctx, userID, guildID, channelID, h.cfg.Menu.SessionTTL)
	if err != nil {
		h.log.Error("/menu: OpenMenu failed",
			zap.String("userId", userID), zap.Error(err))
		ui.EphemeralError(s, i.Interaction, ui.MsgGenericError)
		return
	}

	// 4. Build main menu response
	responseData := menu.BuildMainMenuResponse(&menu.MainMenuData{
		Session:     session,
		Player:      player,
		Cultivation: cult,
		Wallet:      wallet,
	})

	// 5. Respond to the interaction (creates a new public or ephemeral message)
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: responseData,
	})
	if err != nil {
		h.log.Error("/menu: InteractionRespond failed",
			zap.String("userId", userID), zap.Error(err))
		return
	}

	// 6. Save the Discord message ID so future edits work
	msg, err := s.InteractionResponse(i.Interaction)
	if err == nil && msg != nil {
		_ = h.sessionSvc.SetMessageID(ctx, session.SessionID, msg.ID)
	}

	// 7. Update last active
	h.profileSvc.TouchLastActive(ctx, userID, guildID)
}

// PageLoaders returns a map of page → loader functions for use by menu.Router.
// Each loader fetches all required data and calls the corresponding UI builder.
func (h *MenuHandler) PageLoaders() map[menu.Page]menu.PageLoader {
	return map[menu.Page]menu.PageLoader{
		menu.PageMain:        h.loadMainPage,
		menu.PageProfile:     h.loadProfilePage,
		menu.PageCultivation: h.loadCultivationPage,
		// TODO v0.3+: add PageInventory, PageSkills, PagePets, PageGacha, PageMarket, PageSect
	}
}

// loadMainPage fetches all data needed for the main page and returns the rendered response.
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
	return menu.BuildMainMenuEdit(&menu.MainMenuData{
		Session:     session,
		Player:      player,
		Cultivation: cult,
		Wallet:      wallet,
	}), nil
}

// loadProfilePage fetches profile data and returns the profile page response.
func (h *MenuHandler) loadProfilePage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadProfilePage profile: %w", err)
	}
	wallet, err := h.economySvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadProfilePage wallet: %w", err)
	}
	return menu.BuildProfileMenuResponse(&menu.ProfileMenuData{
		Session: session,
		Player:  player,
		Wallet:  wallet,
	}), nil
}

// loadCultivationPage fetches cultivation data and returns the cultivation page response.
func (h *MenuHandler) loadCultivationPage(ctx context.Context, session *menu.Session) (*discordgo.InteractionResponseData, error) {
	player, err := h.profileSvc.GetPlayer(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadCultivationPage profile: %w", err)
	}
	cult, err := h.cultivationSvc.GetOrCreate(ctx, session.UserID, session.GuildID)
	if err != nil {
		return nil, fmt.Errorf("loadCultivationPage cultivation: %w", err)
	}
	return menu.BuildCultivationMenuResponse(&menu.CultivationMenuData{
		Session:     session,
		Player:      player,
		Cultivation: cult,
	}), nil
}
