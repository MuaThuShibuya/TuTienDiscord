// File: cmd/bot/main.go
// Phiên bản: v0.1.1
// Mục đích: Điểm khởi động ứng dụng. Wire toàn bộ dependency và quản lý vòng đời service.
// Bảo mật: Tất cả secret (token, URI, ID) đọc từ env var qua config.Load(), không hardcode.
//           Thoát ngay nếu thiếu env var bắt buộc.
// Ghi chú: Thứ tự injection: config → logger → DB → repository → service → handler → bot.
//          Graceful shutdown khi nhận SIGINT/SIGTERM: xả log, đóng DB, ngắt Discord.

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/database"
	"github.com/whiskey/tu-tien-bot/internal/discord"
	"github.com/whiskey/tu-tien-bot/internal/discord/handlers"
	discordmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/logger"
	"github.com/whiskey/tu-tien-bot/internal/scheduler"
	"github.com/whiskey/tu-tien-bot/internal/server"

	alchemypkg "github.com/whiskey/tu-tien-bot/internal/game/alchemy"
	cooldownpkg "github.com/whiskey/tu-tien-bot/internal/game/cooldown"
	cultivationpkg "github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	economypkg "github.com/whiskey/tu-tien-bot/internal/game/economy"
	equipmentpkg "github.com/whiskey/tu-tien-bot/internal/game/equipment"
	inventorypkg "github.com/whiskey/tu-tien-bot/internal/game/inventory"
	itempkg "github.com/whiskey/tu-tien-bot/internal/game/item"
	profilepkg "github.com/whiskey/tu-tien-bot/internal/game/profile"

	_ "github.com/whiskey/tu-tien-bot/internal/game/data/loader"

	"github.com/bwmarrin/discordgo"
	adminmenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/admin"
	pvemenu "github.com/whiskey/tu-tien-bot/internal/discord/menu/pve"
	gameadminpkg "github.com/whiskey/tu-tien-bot/internal/game/admin"
	aptitudepkg "github.com/whiskey/tu-tien-bot/internal/game/aptitude"
	characterstatspkg "github.com/whiskey/tu-tien-bot/internal/game/characterstats"
	combatpkg "github.com/whiskey/tu-tien-bot/internal/game/combat"
	pvepkg "github.com/whiskey/tu-tien-bot/internal/game/pve"
	pvecombatpkg "github.com/whiskey/tu-tien-bot/internal/game/pvecombat"
)

func main() {
	// Tải .env cho môi trường local (bỏ qua nếu không có file — production dùng env trực tiếp)
	_ = godotenv.Load()

	// --- 1. Cấu hình ---
	cfg, err := config.Load()
	if err != nil {
		// Logger chưa khởi tạo — dùng println
		println("FATAL: lỗi cấu hình:", err.Error())
		os.Exit(1)
	}

	// --- 2. Logger ---
	if err := logger.Init(logger.Options{
		Level:         cfg.Log.Level,
		Format:        cfg.Log.Format,
		Color:         cfg.Log.Color,
		CallerEnabled: cfg.Log.CallerEnabled,
	}); err != nil {
		println("FATAL: không khởi tạo được logger:", err.Error())
		os.Exit(1)
	}
	defer logger.Sync()
	log := logger.L()

	log.Info("Khởi động Vạn Pháp Tiên Nghịch Bot",
		zap.String("app", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
	)

	// --- 3. Kết nối Database ---
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	db, err := database.Connect(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	cancel()
	if err != nil {
		log.Fatal("Không kết nối được MongoDB", zap.Error(err))
	}
	defer func() {
		shutdownCtx, sc := context.WithTimeout(context.Background(), 10*time.Second)
		defer sc()
		_ = db.Disconnect(shutdownCtx)
	}()

	// Tạo indexes MongoDB (TTL, unique, sparse)
	idxCtx, idxCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := database.EnsureIndexes(idxCtx, db.DB()); err != nil {
		log.Fatal("Không tạo được MongoDB indexes", zap.Error(err))
	}
	idxCancel()

	// Chạy Auto Migration (Tự động quét và vá lỗi dữ liệu trên DB)
	migCtx, migCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := database.AutoMigrate(migCtx, db.DB()); err != nil {
		log.Fatal("Lỗi auto migration schema", zap.Error(err))
	}
	migCancel()

	// --- 4. Repositories (chỉ MongoDB) ---
	profileRepo := profilepkg.NewMongoRepository(db.DB())
	cultivationRepo := cultivationpkg.NewMongoRepository(db.DB())
	economyRepo := economypkg.NewMongoRepository(db.DB())
	cooldownRepo := cooldownpkg.NewMongoRepository(db.DB())
	sessionRepo := discordmenu.NewSessionRepository(db.DB())
	itemRepo := itempkg.NewMongoRepository(db.DB())
	invRepo := inventorypkg.NewMongoRepository(db.DB())
	equipRepo := equipmentpkg.NewMongoRepository(db.DB())
	alchemyRepo := alchemypkg.NewMongoRepository(db.DB())
	combatRepo := combatpkg.NewMongoRepository(db.DB())
	pveProgRepo := pvepkg.NewMongoProgressRepository(db.DB())
	aptitudeRepo := aptitudepkg.NewMongoRepository(db.DB())

	// --- 5. Services (business logic) ---
	profileSvc := profilepkg.NewService(profileRepo)
	economySvc := economypkg.NewService(economyRepo)
	cooldownSvc := cooldownpkg.NewService(cooldownRepo)
	cultivationSvc := cultivationpkg.NewService(cultivationRepo, cooldownSvc, economySvc)
	inventorySvc := inventorypkg.NewService(invRepo, itemRepo, cultivationSvc)
	equipSvc := equipmentpkg.NewService(equipRepo, itemRepo, inventorySvc)
	sessionSvc := discordmenu.NewSessionService(sessionRepo)
	alchemySvc := alchemypkg.NewService(alchemyRepo, inventorySvc)
	aptitudeSvc := aptitudepkg.NewService(aptitudeRepo)
	charStatsSvc := characterstatspkg.NewPipelineService(aptitudeSvc, cultivationSvc, equipSvc)

	// Khởi tạo các services và adapter cho hệ thống Combat PvE
	turnOrderSvc := combatpkg.NewTurnOrderService()
	combatSvc, err := combatpkg.NewService(combatRepo, turnOrderSvc, nil)
	if err != nil {
		log.Fatal("Không khởi tạo được combat service", zap.Error(err))
	}

	pveProgSvc := pvepkg.NewProgressService(pveProgRepo)
	statsProv := pvecombatpkg.NewStatsAdapter(charStatsSvc)
	pveProv := pvecombatpkg.NewPvEAdapter(pveProgSvc)
	grantSvc := pvecombatpkg.NewGrantAdapter(inventorySvc)

	pveCombatSvc, err := pvecombatpkg.NewService(combatRepo, statsProv, pveProv, grantSvc, turnOrderSvc, nil)
	if err != nil {
		log.Fatal("Không khởi tạo được pvecombat service", zap.Error(err))
	}

	// --- Admin Service ---
	adminSvc := gameadminpkg.NewService(db.DB(), log)

	// --- 6. Handlers (Controllers) ---
	startHandler := handlers.NewStartHandler(profileSvc, cultivationSvc, economySvc, inventorySvc, aptitudeSvc)
	menuHandler := handlers.NewMenuHandler(cfg, profileSvc, cultivationSvc, economySvc, inventorySvc, equipSvc, alchemySvc, sessionSvc)

	// --- PvE Menu Router ---
	pveRouter := pvemenu.NewRouter(pveCombatSvc, combatSvc, pveProgRepo, log)
	pveActionHandler := func(s *discordgo.Session, i *discordgo.Interaction, session *discordmenu.Session, action string, extra string) {
		pveRouter.HandlePvEInteraction(s, i, session, action, extra)
	}

	// --- Admin Menu Router ---
	adminRouter := adminmenu.NewRouter(cfg, adminSvc, log)
	adminActionHandler := func(s *discordgo.Session, i *discordgo.Interaction, session *discordmenu.Session, action string, extra string) {
		adminRouter.HandleAdminInteraction(s, i, session, action, extra)
	}

	// --- 7. Menu router ---
	menuRouter := discordmenu.NewRouter(cfg, sessionSvc, cultivationSvc, inventorySvc, equipSvc, alchemySvc, pveActionHandler, adminActionHandler, menuHandler.PageLoaders())

	// --- 8. Discord top-level router ---
	discordRouter := discord.NewRouter(startHandler, menuHandler, menuRouter)

	// --- 9. Discord bot ---
	bot, err := discord.NewBot(cfg, discordRouter)
	if err != nil {
		log.Fatal("Không tạo được Discord bot", zap.Error(err))
	}
	// Timeout đủ rộng để chờ Ready event trước khi đăng ký lệnh
	startCtx, startCancel := context.WithTimeout(context.Background(), 30*time.Second)
	if err := bot.Start(startCtx); err != nil {
		startCancel()
		log.Fatal("Không khởi động được Discord bot", zap.Error(err))
	}
	startCancel()

	// --- 10. HTTP server (health check cho Render / uptime monitor) ---
	httpServer := server.NewHTTPServer(cfg, db)
	httpServer.Start()

	// --- 11. Schedulers ---
	keepalive := scheduler.NewKeepalive(cfg)
	keepalive.Start()

	sessionCleaner := scheduler.NewSessionCleaner(db.DB())
	sessionCleaner.Start()

	log.Info("Bot đang chạy. Nhấn Ctrl+C để tắt.")

	// --- 12. Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Nhận tín hiệu tắt — đang dọn dẹp...")

	keepalive.Stop()
	sessionCleaner.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	httpServer.Stop(shutdownCtx)

	bot.Stop()

	log.Info("Tắt hoàn tất.")
}
