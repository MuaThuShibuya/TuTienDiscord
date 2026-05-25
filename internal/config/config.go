// File: internal/config/config.go
// Phiên bản: v0.1.1
// Mục đích: Đọc và validate toàn bộ cấu hình từ biến môi trường.
//           Application ID của Discord được tự động lấy từ Ready event — không cần env var riêng.
// Bảo mật: Không log secret. Tất cả secret đọc từ env var, không hardcode.
// Ghi chú: Gọi config.Load() một lần duy nhất lúc khởi động. Fail fast nếu thiếu biến bắt buộc.

package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config chứa toàn bộ cấu hình ứng dụng.
type Config struct {
	App     AppConfig
	Discord DiscordConfig
	MongoDB MongoDBConfig
	Menu    MenuConfig
	Log     LogConfig
	Server  ServerConfig
}

// AppConfig — thông tin chung của ứng dụng.
type AppConfig struct {
	Env     string
	Name    string
	Version string
}

// DiscordConfig — cấu hình Discord bot.
// AppID KHÔNG cần đặt trong env — được tự động lấy từ Discord Ready event.
// GuildIDs và OwnerIDs hỗ trợ nhiều giá trị phân cách bằng dấu phẩy.
type DiscordConfig struct {
	Token               string
	GuildIDs            []string // Nhiều guild ID, phân cách bằng dấu phẩy
	OwnerIDs            []string // Nhiều owner ID, phân cách bằng dấu phẩy
	CommandRegisterMode string   // "guild" hoặc "global"
}

// IsOwner kiểm tra xem userID có phải là owner không.
func (d *DiscordConfig) IsOwner(userID string) bool {
	for _, id := range d.OwnerIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// MongoDBConfig — kết nối MongoDB Atlas.
type MongoDBConfig struct {
	URI      string
	Database string
}

// MenuConfig — cấu hình hệ thống menu.
type MenuConfig struct {
	SessionTTL time.Duration
}

// LogConfig — cấu hình logging.
// Format và Color chỉ ảnh hưởng đến output format, không ảnh hưởng logic nghiệp vụ.
type LogConfig struct {
	Level         string // "debug" | "info" | "warn" | "error"
	Format        string // "console" (local đẹp) | "json" (production/Render)
	Color         bool   // true = màu ANSI, chỉ có tác dụng khi Format="console"
	CallerEnabled bool   // true = in file:dòng vào mỗi log entry
}

// ServerConfig — cấu hình HTTP server và keepalive.
type ServerConfig struct {
	Port                  string
	KeepaliveEnabled      bool
	KeepaliveURL          string
	KeepaliveIntervalSecs int // giây giữa các lần ping, mặc định 600 (10 phút)
}

// Load đọc và validate cấu hình từ biến môi trường.
// Trả về lỗi nếu thiếu biến bắt buộc.
func Load() (*Config, error) {
	cfg := &Config{}

	// Thông tin ứng dụng
	cfg.App.Env = getEnv("APP_ENV", "development")
	cfg.App.Name = getEnv("APP_NAME", "tu-tien-discord-bot")
	cfg.App.Version = getEnv("APP_VERSION", "0.1.1")

	// Discord — bắt buộc
	token, err := requireEnv("DISCORD_TOKEN")
	if err != nil {
		return nil, err
	}
	cfg.Discord.Token = token
	cfg.Discord.GuildIDs = parseStringSlice(os.Getenv("DISCORD_GUILD_IDS"))
	cfg.Discord.OwnerIDs = parseStringSlice(os.Getenv("DISCORD_OWNER_IDS"))
	cfg.Discord.CommandRegisterMode = getEnv("COMMAND_REGISTER_MODE", "guild")

	// Kiểm tra: guild mode cần ít nhất 1 guild ID
	if cfg.Discord.CommandRegisterMode == "guild" && len(cfg.Discord.GuildIDs) == 0 {
		return nil, fmt.Errorf("COMMAND_REGISTER_MODE=guild nhưng DISCORD_GUILD_IDS trống")
	}

	// MongoDB — bắt buộc
	mongoURI, err := requireEnv("MONGODB_URI")
	if err != nil {
		return nil, err
	}
	cfg.MongoDB.URI = mongoURI
	cfg.MongoDB.Database = getEnv("MONGODB_DATABASE", "tu_tien_bot")

	// Menu session TTL
	ttlMin := getEnvInt("MENU_SESSION_TTL_MINUTES", 15)
	cfg.Menu.SessionTTL = time.Duration(ttlMin) * time.Minute

	// Logging
	cfg.Log.Level = getEnv("LOG_LEVEL", "info")
	cfg.Log.Format = getEnv("LOG_FORMAT", "json")  // mặc định json — an toàn cho CI/production
	cfg.Log.Color = getEnvBool("LOG_COLOR", false) // tắt mặc định — tránh ANSI code lọt vào file log
	cfg.Log.CallerEnabled = getEnvBool("LOG_CALLER", false)

	// HTTP Server
	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.KeepaliveEnabled = getEnvBool("KEEPALIVE_ENABLED", false)
	cfg.Server.KeepaliveURL = os.Getenv("KEEPALIVE_URL")
	cfg.Server.KeepaliveIntervalSecs = getEnvInt("KEEPALIVE_INTERVAL_SECONDS", 600)

	return cfg, nil
}

// --- Helpers ---

// requireEnv lấy biến env bắt buộc. Trả về lỗi nếu không có.
func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("biến môi trường bắt buộc %q chưa được thiết lập", key)
	}
	return val, nil
}

// getEnv lấy biến env với giá trị mặc định.
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// getEnvInt lấy biến env số nguyên với giá trị mặc định.
func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}

// getEnvBool lấy biến env boolean với giá trị mặc định.
func getEnvBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

// parseStringSlice tách chuỗi phân cách bằng dấu phẩy thành slice, bỏ khoảng trắng.
// Ví dụ: "id1, id2, id3" → ["id1", "id2", "id3"]
func parseStringSlice(val string) []string {
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
