// File: internal/server/http_server.go
// Phiên bản: v0.1.1
// Mục đích: HTTP server nhẹ cho health check và các endpoint admin/REST trong tương lai.
// Bảo mật: Hiện chỉ expose /health. Endpoint này không yêu cầu auth (không chứa thông tin nhạy cảm).
//           Endpoint admin tương lai phải có auth middleware trước khi thêm vào.
// Ghi chú: Trên VPS/Docker, có thể bỏ server này và dùng Docker health check hoặc systemd thay thế.

package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/database"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// HTTPServer bọc HTTP server của standard library.
type HTTPServer struct {
	server *http.Server
	log    *zap.Logger
}

// NewHTTPServer tạo và cấu hình HTTP server với tất cả route.
func NewHTTPServer(cfg *config.Config, db *database.Client) *HTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg, db))

	// TODO: Thêm route admin/dashboard ở đây trong phiên bản tương lai (cần auth middleware).

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &HTTPServer{
		server: srv,
		log:    logger.L().Named("server.http"),
	}
}

// Start bắt đầu lắng nghe trong goroutine. Không blocking.
func (s *HTTPServer) Start() {
	go func() {
		s.log.Info("HTTP server đang lắng nghe", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("HTTP server lỗi", zap.Error(err))
		}
	}()
}

// Stop tắt HTTP server graceful.
func (s *HTTPServer) Stop(ctx context.Context) {
	s.log.Info("Đang tắt HTTP server")
	_ = s.server.Shutdown(ctx)
}
