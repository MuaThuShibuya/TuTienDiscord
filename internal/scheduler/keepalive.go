// File: internal/scheduler/keepalive.go
// Phiên bản: v0.1.1
// Mục đích: Self-ping keepalive cho Render free tier — ngăn service bị sleep.
// Bảo mật: Chỉ ping endpoint /health của chính bot. Không ping URL ngoài trừ khi được cấu hình.
//           URL keepalive chỉ lấy từ env var KEEPALIVE_URL.
// Ghi chú: Render free tier sleep sau ~15 phút không có traffic. Bot tự ping mỗi 10 phút.
//          Trên VPS/Docker có process manager, tắt keepalive bằng KEEPALIVE_ENABLED=false.

package scheduler

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

const (
	keepaliveDefaultInterval = 10 * time.Minute // dùng nếu KEEPALIVE_INTERVAL_SECONDS không được đặt
	keepaliveTimeout         = 15 * time.Second
)

// Keepalive ping endpoint health của bot theo chu kỳ định kỳ.
type Keepalive struct {
	cfg    *config.Config
	client *http.Client
	log    *zap.Logger
	stop   chan struct{}
}

// NewKeepalive tạo keepalive scheduler.
func NewKeepalive(cfg *config.Config) *Keepalive {
	return &Keepalive{
		cfg:    cfg,
		client: &http.Client{Timeout: keepaliveTimeout},
		log:    logger.L().Named("scheduler.keepalive"),
		stop:   make(chan struct{}),
	}
}

// Start bắt đầu vòng lặp keepalive trong goroutine. Không blocking.
func (k *Keepalive) Start() {
	if !k.cfg.Server.KeepaliveEnabled {
		k.log.Info("Keepalive đã tắt (KEEPALIVE_ENABLED=false)")
		return
	}
	url := k.cfg.Server.KeepaliveURL
	if url == "" {
		url = "http://localhost:" + k.cfg.Server.Port + "/health"
	}
	interval := keepaliveDefaultInterval
	if secs := k.cfg.Server.KeepaliveIntervalSecs; secs > 0 {
		interval = time.Duration(secs) * time.Second
	}
	k.log.Info("Keepalive đã bật", zap.String("url", url), zap.Duration("interval", interval))
	go k.run(url, interval)
}

// Stop phát tín hiệu dừng vòng lặp keepalive.
func (k *Keepalive) Stop() {
	close(k.stop)
}

func (k *Keepalive) run(url string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			k.ping(url)
		case <-k.stop:
			k.log.Info("Keepalive đã dừng")
			return
		}
	}
}

func (k *Keepalive) ping(url string) {
	ctx, cancel := context.WithTimeout(context.Background(), keepaliveTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		k.log.Error("Keepalive: không tạo được request", zap.Error(err))
		return
	}

	resp, err := k.client.Do(req)
	if err != nil {
		k.log.Warn("Keepalive ping thất bại", zap.String("url", url), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	k.log.Debug("Keepalive ping OK", zap.String("url", url), zap.Int("status", resp.StatusCode))
}
