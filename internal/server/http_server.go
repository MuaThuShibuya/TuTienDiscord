// File: internal/server/http_server.go
// Version: v0.1
// Purpose: Lightweight HTTP server for health checks and future REST/admin endpoints.
// Security: Only exposes /health. No auth on health endpoint (non-sensitive). Future admin
//           endpoints must require auth before being added.
// Notes: On VPS/Docker, this can be removed in favor of Docker health checks or systemd.

package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/database"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// HTTPServer wraps the standard library HTTP server.
type HTTPServer struct {
	server *http.Server
	log    *zap.Logger
}

// NewHTTPServer creates and configures the HTTP server with all routes.
func NewHTTPServer(cfg *config.Config, db *database.Client) *HTTPServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthHandler(cfg, db))

	// TODO: Add admin/dashboard routes here in a future version (with auth middleware).

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

// Start begins listening in a goroutine. Non-blocking.
func (s *HTTPServer) Start() {
	go func() {
		s.log.Info("HTTP server listening", zap.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("HTTP server error", zap.Error(err))
		}
	}()
}

// Stop gracefully shuts down the HTTP server.
func (s *HTTPServer) Stop(ctx context.Context) {
	s.log.Info("Shutting down HTTP server")
	_ = s.server.Shutdown(ctx)
}
