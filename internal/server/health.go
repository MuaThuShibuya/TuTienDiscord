// File: internal/server/health.go
// Version: v0.1
// Purpose: HTTP health check handler for Render/uptime monitoring and keepalive pings.
// Security: Returns no internal details — only status, app name, version, and DB connectivity.
// Notes: GET /health is the only endpoint. Render uses this to verify the service is alive.

package server

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/database"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

// HealthResponse is the JSON body returned by GET /health.
type HealthResponse struct {
	Status   string `json:"status"`
	App      string `json:"app"`
	Version  string `json:"version"`
	Database string `json:"database"`
}

// HealthHandler returns an http.HandlerFunc for the /health endpoint.
func HealthHandler(cfg *config.Config, db *database.Client) http.HandlerFunc {
	log := logger.L().Named("server.health")

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		dbStatus := "connected"
		if !db.IsConnected() {
			dbStatus = "disconnected"
			log.Warn("Health check: MongoDB is not reachable")
		}

		resp := HealthResponse{
			Status:   "ok",
			App:      cfg.App.Name,
			Version:  cfg.App.Version,
			Database: dbStatus,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("Health handler: failed to write response", zap.Error(err))
		}
	}
}
