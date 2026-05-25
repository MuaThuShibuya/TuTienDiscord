// File: internal/config/config.go
// Version: v0.1
// Purpose: Load and validate all application configuration from environment variables.
// Security: Never log or expose secret values. All secrets come from env vars only.
// Notes: Call config.Load() once at startup. Fail fast if required vars are missing.

package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	App      AppConfig
	Discord  DiscordConfig
	MongoDB  MongoDBConfig
	Menu     MenuConfig
	Log      LogConfig
	Server   ServerConfig
}

type AppConfig struct {
	Env     string
	Name    string
	Version string
}

type DiscordConfig struct {
	Token             string
	AppID             string
	GuildID           string
	OwnerID           string
	CommandRegisterMode string // "guild" or "global"
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type MenuConfig struct {
	SessionTTL time.Duration
}

type LogConfig struct {
	Level string // debug, info, warn, error
}

type ServerConfig struct {
	Port            string
	KeepaliveEnabled bool
	KeepaliveURL    string
}

// Load reads config from environment variables and returns a validated Config.
// It will return an error if any required variable is missing.
func Load() (*Config, error) {
	cfg := &Config{}

	cfg.App.Env = getEnv("APP_ENV", "development")
	cfg.App.Name = getEnv("APP_NAME", "tu-tien-discord-bot")
	cfg.App.Version = getEnv("APP_VERSION", "0.1.0")

	// Discord — required
	token, err := requireEnv("DISCORD_TOKEN")
	if err != nil {
		return nil, err
	}
	cfg.Discord.Token = token

	appID, err := requireEnv("DISCORD_APP_ID")
	if err != nil {
		return nil, err
	}
	cfg.Discord.AppID = appID

	cfg.Discord.GuildID = os.Getenv("DISCORD_GUILD_ID")
	cfg.Discord.OwnerID = os.Getenv("DISCORD_OWNER_ID")
	cfg.Discord.CommandRegisterMode = getEnv("COMMAND_REGISTER_MODE", "guild")

	// MongoDB — required
	mongoURI, err := requireEnv("MONGODB_URI")
	if err != nil {
		return nil, err
	}
	cfg.MongoDB.URI = mongoURI
	cfg.MongoDB.Database = getEnv("MONGODB_DATABASE", "tu_tien_bot")

	// Menu session TTL
	ttlMinutes := getEnvInt("MENU_SESSION_TTL_MINUTES", 15)
	cfg.Menu.SessionTTL = time.Duration(ttlMinutes) * time.Minute

	// Logging
	cfg.Log.Level = getEnv("LOG_LEVEL", "info")

	// Server
	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.KeepaliveEnabled = getEnvBool("KEEPALIVE_ENABLED", false)
	cfg.Server.KeepaliveURL = os.Getenv("KEEPALIVE_URL")

	return cfg, nil
}

func requireEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return val, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}
