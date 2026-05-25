// File: internal/scheduler/keepalive.go
// Version: v0.1
// Purpose: Self-ping keepalive for Render free tier — prevents the service from sleeping.
// Security: Only pings the bot's own /health endpoint. No external URLs unless configured.
//           Keepalive URL comes from env var KEEPALIVE_URL only.
// Notes: Render free tier sleeps after ~15 min of no inbound traffic. This pings every 10 min.
//        On VPS/Docker with a process manager, disable keepalive via KEEPALIVE_ENABLED=false.
//        Keepalive is intentionally not a goroutine flood — it pings once every interval with a timeout.

package scheduler

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/config"
	"github.com/yourname/tu-tien-bot/internal/logger"
)

const (
	keepaliveInterval = 10 * time.Minute
	keepaliveTimeout  = 15 * time.Second
)

// Keepalive pings the bot's own health endpoint at regular intervals.
type Keepalive struct {
	cfg    *config.Config
	client *http.Client
	log    *zap.Logger
	stop   chan struct{}
}

// NewKeepalive creates a keepalive scheduler.
func NewKeepalive(cfg *config.Config) *Keepalive {
	return &Keepalive{
		cfg:    cfg,
		client: &http.Client{Timeout: keepaliveTimeout},
		log:    logger.L().Named("scheduler.keepalive"),
		stop:   make(chan struct{}),
	}
}

// Start begins the keepalive loop in a goroutine. Non-blocking.
func (k *Keepalive) Start() {
	if !k.cfg.Server.KeepaliveEnabled {
		k.log.Info("Keepalive disabled (KEEPALIVE_ENABLED=false)")
		return
	}
	url := k.cfg.Server.KeepaliveURL
	if url == "" {
		url = "http://localhost:" + k.cfg.Server.Port + "/health"
	}
	k.log.Info("Keepalive started", zap.String("url", url), zap.Duration("interval", keepaliveInterval))
	go k.run(url)
}

// Stop signals the keepalive loop to exit.
func (k *Keepalive) Stop() {
	close(k.stop)
}

func (k *Keepalive) run(url string) {
	ticker := time.NewTicker(keepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			k.ping(url)
		case <-k.stop:
			k.log.Info("Keepalive stopped")
			return
		}
	}
}

func (k *Keepalive) ping(url string) {
	ctx, cancel := context.WithTimeout(context.Background(), keepaliveTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		k.log.Error("Keepalive: failed to build request", zap.Error(err))
		return
	}

	resp, err := k.client.Do(req)
	if err != nil {
		k.log.Warn("Keepalive ping failed", zap.String("url", url), zap.Error(err))
		return
	}
	defer resp.Body.Close()

	k.log.Debug("Keepalive ping OK", zap.String("url", url), zap.Int("status", resp.StatusCode))
}
