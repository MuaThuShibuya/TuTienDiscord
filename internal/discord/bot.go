// File: internal/discord/bot.go
// Phiên bản: v0.1.1
// Mục đích: Khởi tạo và quản lý vòng đời Discord bot session.
//           Kết nối Discord, đăng ký lệnh cho nhiều guild, wire handler, tắt graceful.
// Bảo mật: Token đọc từ env var, không bao giờ log. Application ID tự phát hiện từ sự kiện Ready.
// Ghi chú: Chế độ guild (instant) hoặc global (tối đa 1 tiếng truyền bá).
//          Hỗ trợ nhiều guild ID — đăng ký lệnh cho từng guild riêng biệt.

package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// Bot quản lý Discord session và vòng đời của nó.
type Bot struct {
	session          *discordgo.Session
	cfg              *config.Config
	router           *Router
	appID            string // Tự phát hiện từ sự kiện Ready — không cần env var
	registeredCmdIDs []string
	log              *zap.Logger
}

// NewBot tạo Discord bot session mới nhưng chưa kết nối.
func NewBot(cfg *config.Config, router *Router) (*Bot, error) {
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		return nil, fmt.Errorf("discord.NewBot: không tạo được session: %w", err)
	}

	// Slash command không cần privileged intents. IntentsGuilds đủ để nhận guild info.
	dg.Identify.Intents = discordgo.IntentsGuilds

	return &Bot{
		session: dg,
		cfg:     cfg,
		router:  router,
		log:     logger.L().Named("discord.bot"),
	}, nil
}

// Start mở kết nối WebSocket đến Discord và đăng ký slash command.
// Chờ sự kiện Ready để lấy Application ID trước khi đăng ký lệnh.
func (b *Bot) Start(ctx context.Context) error {
	// Channel nhận Application ID từ sự kiện Ready
	readyCh := make(chan string, 1)

	b.session.AddHandler(b.router.HandleInteraction)

	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		b.log.Info("Discord bot đã sẵn sàng",
			zap.String("username", r.User.Username),
			zap.String("id", r.User.ID),
			zap.Int("guilds", len(r.Guilds)),
		)
		readyCh <- r.Application.ID
	})

	if err := b.session.Open(); err != nil {
		return fmt.Errorf("discord.Bot.Start: không mở được kết nối: %w", err)
	}

	b.log.Info("Discord WebSocket đã kết nối")

	// Chờ sự kiện Ready để lấy Application ID
	select {
	case appID := <-readyCh:
		b.appID = appID
	case <-ctx.Done():
		return fmt.Errorf("discord.Bot.Start: context bị huỷ khi chờ sự kiện Ready")
	}

	if err := b.registerCommands(); err != nil {
		return fmt.Errorf("discord.Bot.Start: không đăng ký được lệnh: %w", err)
	}

	return nil
}

// Stop đóng kết nối Discord WebSocket một cách graceful.
func (b *Bot) Stop() {
	b.log.Info("Đang tắt Discord bot")
	_ = b.session.Close()
}

// registerCommands đăng ký tất cả slash command với Discord.
// Chế độ guild: đăng ký cho từng guild trong DISCORD_GUILD_IDS (hiệu lực tức thì).
// Chế độ global: đăng ký toàn cầu (tối đa 1 tiếng để truyền bá).
func (b *Bot) registerCommands() error {
	commands := AllCommands()
	mode := b.cfg.Discord.CommandRegisterMode

	if mode == "guild" {
		if len(b.cfg.Discord.GuildIDs) == 0 {
			b.log.Warn("CommandRegisterMode=guild nhưng DISCORD_GUILD_IDS rỗng, chuyển sang chế độ global")
			mode = "global"
		}
	}

	if mode == "guild" {
		for _, guildID := range b.cfg.Discord.GuildIDs {
			if err := b.registerForGuild(guildID, commands); err != nil {
				return err
			}
		}
	} else {
		// Global — guildID rỗng = đăng ký toàn cầu
		if err := b.registerForGuild("", commands); err != nil {
			return err
		}
	}

	b.log.Info("Tất cả lệnh đã đăng ký", zap.Int("tổng", len(b.registeredCmdIDs)))
	return nil
}

// registerForGuild đăng ký lệnh cho một guild cụ thể (hoặc global nếu guildID rỗng).
func (b *Bot) registerForGuild(guildID string, commands []*discordgo.ApplicationCommand) error {
	scope := guildID
	if scope == "" {
		scope = "global"
	}

	b.log.Info("Đăng ký lệnh",
		zap.String("scope", scope),
		zap.Int("số lệnh", len(commands)),
	)

	for _, cmd := range commands {
		registered, err := b.session.ApplicationCommandCreate(b.appID, guildID, cmd)
		if err != nil {
			return fmt.Errorf("registerForGuild %q: không đăng ký được lệnh %q: %w", scope, cmd.Name, err)
		}
		b.registeredCmdIDs = append(b.registeredCmdIDs, registered.ID)
		b.log.Debug("Lệnh đã đăng ký", zap.String("name", cmd.Name), zap.String("id", registered.ID))
	}

	return nil
}
