// File: internal/discord/bot.go
// Version: v0.1
// Purpose: Bootstrap and lifecycle management for the Discord bot session.
//          Connects to Discord, registers commands, wires handlers, and cleanly shuts down.
// Security: Token is read from config (env var). Never logged. Intent is minimal (no privileged intents needed for slash commands).
// Notes: RegisterCommands supports guild-mode (instant) or global-mode (up to 1h propagation).

package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// Bot wraps the discordgo session and manages its lifecycle.
type Bot struct {
	session          *discordgo.Session
	cfg              *config.Config
	router           *Router
	registeredCmdIDs []string
	log              *zap.Logger
}

// NewBot creates a new Discord bot session but does not yet connect.
func NewBot(cfg *config.Config, router *Router) (*Bot, error) {
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return nil, fmt.Errorf("discord.NewBot: failed to create session: %w", err)
	}

	// Only need guild messages intent for prefix commands; slash commands work with no intents.
	dg.Identify.Intents = discordgo.IntentsGuilds

	return &Bot{
		session: dg,
		cfg:     cfg,
		router:  router,
		log:     logger.L().Named("discord.bot"),
	}, nil
}

// Start opens the WebSocket connection to Discord and registers slash commands.
func (b *Bot) Start(ctx context.Context) error {
	// Register interaction handler before opening connection
	b.session.AddHandler(b.router.HandleInteraction)

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.log.Info("Discord bot is ready",
			zap.String("username", r.User.Username),
			zap.String("id", r.User.ID),
		)
	})

	if err := b.session.Open(); err != nil {
		return fmt.Errorf("discord.Bot.Start: failed to open connection: %w", err)
	}

	b.log.Info("Discord WebSocket connection opened")

	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("discord.Bot.Start: failed to register commands: %w", err)
	}

	return nil
}

// Stop cleanly closes the Discord WebSocket connection and deregisters commands if needed.
func (b *Bot) Stop() {
	b.log.Info("Shutting down Discord bot")
	_ = b.session.Close()
}

// registerCommands registers all slash commands with Discord.
// In "guild" mode they appear instantly. In "global" mode they can take up to 1 hour.
func (b *Bot) registerCommands() error {
	commands := AllCommands()
	mode := b.cfg.Discord.CommandRegisterMode
	guildID := ""

	if mode == "guild" {
		guildID = b.cfg.Discord.GuildID
		if guildID == "" {
			b.log.Warn("CommandRegisterMode=guild but DISCORD_GUILD_ID is empty; falling back to global")
		}
	}

	b.log.Info("Registering slash commands",
		zap.String("mode", mode),
		zap.String("guildId", guildID),
		zap.Int("count", len(commands)),
	)

	for _, cmd := range commands {
		registered, err := b.session.ApplicationCommandCreate(
			b.cfg.Discord.AppID,
			guildID,
			cmd,
		)
		if err != nil {
			return fmt.Errorf("registerCommands: failed to register %q: %w", cmd.Name, err)
		}
		b.registeredCmdIDs = append(b.registeredCmdIDs, registered.ID)
		b.log.Debug("Command registered", zap.String("name", cmd.Name), zap.String("id", registered.ID))
	}

	b.log.Info("All commands registered successfully", zap.Int("count", len(b.registeredCmdIDs)))
	return nil
}
