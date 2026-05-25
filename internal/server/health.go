// File: internal/server/health.go
// Phiên bản: v0.1.1
// Mục đích: HTTP health check handler cho Render/uptime monitor và keepalive ping.
// Bảo mật: Không trả về chi tiết nội bộ — chỉ trả về status, tên app, version, và trạng thái DB.
// Ghi chú: GET /health là endpoint duy nhất. Render dùng endpoint này để xác minh service đang chạy.

package server

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/database"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// HealthResponse là JSON body trả về từ GET /health.
type HealthResponse struct {
	Status   string `json:"status"`
	App      string `json:"app"`
	Version  string `json:"version"`
	Database string `json:"database"`
}

// HealthHandler trả về http.HandlerFunc cho endpoint /health.
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
			log.Warn("Health check: MongoDB không kết nối được")
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
			log.Error("Health handler: không ghi được response", zap.Error(err))
		}
	}
}
