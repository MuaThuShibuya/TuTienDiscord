// File: internal/discord/menu/router.go
// Version: v0.1
// Purpose: Route button/select-menu interactions inside the in-app menu system.
//          Parses custom_id, validates session ownership, dispatches to the correct page renderer.
// Security: Every interaction validates sessionId + userId before any state change or DB write.
//           Users operating someone else's menu receive an ephemeral error message.
// Notes: custom_id format: "<domain>:<action>:<sessionId>[:<extra>]"
//        Example: "menu:nav:abc123" or "nav:back:abc123:main"

package menu

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// PageLoader is a function that fetches all data needed for a page and returns the response data.
type PageLoader func(ctx context.Context, session *Session) (*discordgo.InteractionResponseData, error)

// Router handles all component interactions (button / select menu) inside the menu system.
type Router struct {
	cfg         *config.Config
	sessionSvc  Service
	pageLoaders map[Page]PageLoader
	log         *zap.Logger
}

// NewRouter creates a menu interaction router.
// pageLoaders maps each Page to a function that loads that page's data and returns its rendered response.
func NewRouter(cfg *config.Config, sessionSvc Service, loaders map[Page]PageLoader) *Router {
	return &Router{
		cfg:         cfg,
		sessionSvc:  sessionSvc,
		pageLoaders: loaders,
		log:         logger.L().Named("menu.router"),
	}
}

// Handle dispatches a Discord component interaction to the correct handler.
// Called by the top-level discord router for all MessageComponent interactions.
func (r *Router) Handle(s *discordgo.Session, i *discordgo.Interaction) {
	var customID string
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		customID = i.MessageComponentData().CustomID
	default:
		return
	}

	// Parse custom_id: "<domain>:<action>:<sessionId>[:<extra>]"
	parts := strings.SplitN(customID, ":", 4)
	if len(parts) < 3 {
		r.log.Warn("menu.Router: malformed custom_id", zap.String("customID", customID))
		ui.EphemeralError(s, i, ui.MsgGenericError)
		return
	}

	domain := parts[0]
	action := parts[1]
	sessionID := parts[2]
	extra := ""
	if len(parts) == 4 {
		extra = parts[3]
	}

	ctx := context.Background()
	userID := i.Member.User.ID

	// --- Security: validate session ownership ---
	session, err := r.sessionSvc.ValidateOwner(ctx, sessionID, userID)
	if err != nil {
		switch {
		case apperrors.IsSessionExpired(err):
			ui.EphemeralError(s, i, ui.MsgSessionExpired)
		default:
			ui.EphemeralError(s, i, ui.MsgNotYourMenu)
		}
		return
	}

	// Refresh TTL on every interaction.
	_ = r.sessionSvc.Refresh(ctx, sessionID, r.cfg.Menu.SessionTTL)

	switch domain {
	case "nav":
		r.handleNav(s, i, session, action, extra)
	case "menu":
		r.handleMenuSelect(s, i, session, action, i.MessageComponentData())
	case "profile":
		r.handleProfileAction(s, i, session, action)
	case "cultivation":
		r.handleCultivationAction(s, i, session, action)
	default:
		r.log.Warn("menu.Router: unknown domain", zap.String("domain", domain), zap.String("customID", customID))
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	}
}

// handleNav processes Làm mới / Quay lại / Đóng buttons.
func (r *Router) handleNav(s *discordgo.Session, i *discordgo.Interaction, session *Session, action, extra string) {
	ctx := context.Background()

	switch action {
	case "close":
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

	case "refresh":
		page := Page(extra)
		if page == "" {
			page = session.CurrentPage
		}
		r.renderPage(s, i, session, page)

	case "back":
		targetPage := Page(extra)
		if targetPage == "" {
			targetPage = PageMain
		}
		_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, targetPage)
		session.CurrentPage = targetPage
		r.renderPage(s, i, session, targetPage)

	default:
		ui.EphemeralError(s, i, ui.MsgGenericError)
	}
}

// handleMenuSelect processes the category select menus on the main page.
func (r *Router) handleMenuSelect(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string, data discordgo.MessageComponentInteractionData) {
	if len(data.Values) == 0 {
		return
	}
	selected := Page(data.Values[0])
	ctx := context.Background()
	_ = r.sessionSvc.NavigateTo(ctx, session.SessionID, selected)
	session.CurrentPage = selected
	r.renderPage(s, i, session, selected)
}

// handleProfileAction dispatches profile-specific button actions.
func (r *Router) handleProfileAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	switch action {
	case "rename":
		// TODO v0.1: open modal for dao name change
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	default:
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	}
}

// handleCultivationAction dispatches cultivation-specific button actions.
func (r *Router) handleCultivationAction(s *discordgo.Session, i *discordgo.Interaction, session *Session, action string) {
	switch action {
	case "meditate":
		// TODO v0.2: tĩnh tu — cooldown check → exp gain → update DB
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	case "closeddoor":
		// TODO v0.2: bế quan
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	case "breakthrough":
		// TODO v0.2: đột phá
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	default:
		ui.EphemeralError(s, i, ui.MsgComingSoon)
	}
}

// renderPage calls the appropriate PageLoader and edits the existing menu message.
func (r *Router) renderPage(s *discordgo.Session, i *discordgo.Interaction, session *Session, page Page) {
	loader, ok := r.pageLoaders[page]
	if !ok {
		ui.EphemeralError(s, i, ui.MsgComingSoon)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	responseData, err := loader(ctx, session)
	if err != nil {
		r.log.Error("menu.Router: page loader error",
			zap.String("page", string(page)),
			zap.String("userId", session.UserID),
			zap.Error(err),
		)
		ui.EphemeralError(s, i, ui.MsgGenericError)
		return
	}

	_ = s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: responseData,
	})
}
