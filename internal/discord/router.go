// File: internal/discord/router.go
// Version: v0.1
// Purpose: Top-level Discord interaction router — dispatches slash commands and component interactions
//          to the correct handler. This is the single point of entry for all Discord events.
// Security: Only processes interactions from guilds (no DMs). Unknown commands are ignored gracefully.
// Notes: Add new command routes here as features grow. Keep this file thin — no business logic.

package discord

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/discord/handlers"
	"github.com/yourname/tu-tien-bot/internal/discord/menu"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// Router holds all registered handlers and routes interactions to them.
type Router struct {
	startHandler *handlers.StartHandler
	menuHandler  *handlers.MenuHandler
	menuRouter   *menu.Router
	log          *zap.Logger
}

// NewRouter creates the top-level interaction router.
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

// HandleInteraction is the discordgo event handler registered on the session.
// It routes each interaction to the correct handler based on type and name.
func (r *Router) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// All interactions must come from a guild member
	if i.Member == nil {
		r.log.Debug("Ignoring DM interaction")
		return
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		r.routeCommand(s, i)

	case discordgo.InteractionMessageComponent:
		// All component interactions go through the menu router
		r.menuRouter.Handle(s, i.Interaction)

	case discordgo.InteractionModalSubmit:
		// TODO v0.1+: route modal submissions (e.g., dao name rename modal)
		r.log.Debug("Modal submit received (not yet handled)",
			zap.String("customID", i.ModalSubmitData().CustomID))

	default:
		r.log.Debug("Unknown interaction type", zap.Int("type", int(i.Type)))
	}
}

// routeCommand dispatches slash commands to the appropriate handler.
func (r *Router) routeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Name

	r.log.Debug("Slash command received",
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
		r.log.Warn("Unknown command", zap.String("command", name))
	}
}
