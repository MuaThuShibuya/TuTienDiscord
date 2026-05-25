// File: cmd/bot/main.go
// Version: v0.1
// Purpose: Application entrypoint. Wires all dependencies and manages the service lifecycle.
// Security: Loads all secrets from environment variables via config.Load(). Never hardcodes tokens.
//           Fails fast on missing required environment variables.
// Notes: Dependency injection order: config → logger → DB → repositories → services → handlers → bot.
//        Graceful shutdown on SIGINT/SIGTERM flushes logs and closes all connections cleanly.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/database"
	"github.com/yourname/tu-tien-bot/internal/discord"
	"github.com/yourname/tu-tien-bot/internal/discord/handlers"
	discordmenu "github.com/yourname/tu-tien-bot/internal/discord/menu"
	"github.com/yourname/tu-tien-bot/internal/logger"
	"github.com/yourname/tu-tien-bot/internal/scheduler"
	"github.com/yourname/tu-tien-bot/internal/server"

	cultivationrepo "github.com/yourname/tu-tien-bot/internal/game/cultivation"
	economyrepo "github.com/yourname/tu-tien-bot/internal/game/economy"
	profilerepo "github.com/yourname/tu-tien-bot/internal/game/profile"
	cooldownrepo "github.com/yourname/tu-tien-bot/internal/game/cooldown"
)

func main() {
	// Load .env for local development (ignored in production where env vars are set directly)
	_ = godotenv.Load()

	// --- 1. Config ---
	cfg, err := config.Load()
	if err != nil {
		// Logger not yet initialized; use fmt
		println("FATAL: config error:", err.Error())
		os.Exit(1)
	}

	// --- 2. Logger ---
	if err := logger.Init(cfg.Log.Level); err != nil {
		println("FATAL: logger init failed:", err.Error())
		os.Exit(1)
	}
	defer logger.Sync()
	log := logger.L()

	log.Info("Starting Tu Tien Bot",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
	)

	// --- 3. Database ---
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, err := database.Connect(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	cancel()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer func() {
		shutdownCtx, sc := context.WithTimeout(context.Background(), 10*time.Second)
		defer sc()
		_ = db.Disconnect(shutdownCtx)
	}()

	// Ensure indexes
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := database.EnsureIndexes(idxCtx, db.DB()); err != nil {
		log.Fatal("Failed to ensure MongoDB indexes", zap.Error(err))
	}
	idxCancel()

	// --- 4. Repositories ---
	profileRepo     := profilerepo.NewRepository(db.DB())
	cultivationRepo := cultivationrepo.NewRepository(db.DB())
	economyRepo     := economyrepo.NewRepository(db.DB())
	cooldownRepo    := cooldownrepo.NewRepository(db.DB())
	sessionRepo     := discordmenu.NewRepository(db.DB())

	// --- 5. Services ---
	profileSvc     := profilerepo.NewService(profileRepo)
	cultivationSvc := cultivationrepo.NewService(cultivationRepo)
	economySvc     := economyrepo.NewService(economyRepo)
	_              = cooldownrepo.NewService(cooldownRepo) // registered; used from v0.2
	sessionSvc     := discordmenu.NewService(sessionRepo)

	// --- 6. Handlers (Controllers) ---
	startHandler := handlers.NewStartHandler(profileSvc, cultivationSvc, economySvc)
	menuHandler  := handlers.NewMenuHandler(cfg, profileSvc, cultivationSvc, economySvc, sessionSvc)

	// --- 7. Menu router ---
	menuRouter := discordmenu.NewRouter(cfg, sessionSvc, menuHandler.PageLoaders())

	// --- 8. Discord top-level router ---
	discordRouter := discord.NewRouter(startHandler, menuHandler, menuRouter)

	// --- 9. Discord bot ---
	bot, err := discord.NewBot(cfg, discordRouter)
	if err != nil {
		log.Fatal("Failed to create Discord bot", zap.Error(err))
	}
	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := bot.Start(startCtx); err != nil {
		startCancel()
		log.Fatal("Failed to start Discord bot", zap.Error(err))
	}
	startCancel()

	// --- 10. HTTP server (health check) ---
	httpServer := server.NewHTTPServer(cfg, db)
	httpServer.Start()

	// --- 11. Schedulers ---
	keepalive := scheduler.NewKeepalive(cfg)
	keepalive.Start()

	sessionCleaner := scheduler.NewSessionCleaner(db.DB())
	sessionCleaner.Start()

	log.Info("Bot is running. Press Ctrl+C to stop.")

	// --- 12. Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received — cleaning up...")

	keepalive.Stop()
	sessionCleaner.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	httpServer.Stop(shutdownCtx)

	bot.Stop()

	log.Info("Shutdown complete.")
}
