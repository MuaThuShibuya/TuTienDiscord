// File: internal/database/indexes.go
// Version: v0.1
// Purpose: Create and ensure MongoDB indexes for all collections on startup.
// Security: No secrets involved. Index creation is idempotent and safe to run multiple times.
// Notes: Add new indexes here as new collections are introduced in later versions.

package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/yourname/tu-tien-bot/internal/logger"
)

const indexTimeout = 30 * time.Second

// EnsureIndexes creates all required indexes for the application.
// This is idempotent — safe to call on every startup.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	log := logger.L()

	indexCtx, cancel := context.WithTimeout(ctx, indexTimeout)
	defer cancel()

	var errs []error

	if err := ensurePlayerIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Failed to create player indexes", zap.Error(err))
	}

	if err := ensureCultivationIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Failed to create cultivation indexes", zap.Error(err))
	}

	if err := ensureWalletIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Failed to create wallet indexes", zap.Error(err))
	}

	if err := ensureCooldownIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Failed to create cooldown indexes", zap.Error(err))
	}

	if err := ensureMenuSessionIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Failed to create menu_sessions indexes", zap.Error(err))
	}

	if len(errs) > 0 {
		return errs[0]
	}

	log.Info("MongoDB indexes ensured successfully")
	return nil
}

func ensurePlayerIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("players")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_player_user_guild"),
		},
		{
			Keys:    bson.D{{Key: "lastActiveAt", Value: -1}},
			Options: options.Index().SetName("idx_player_last_active"),
		},
	})
	return err
}

func ensureCultivationIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("cultivation_profiles")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_cultivation_user_guild"),
		},
	})
	return err
}

func ensureWalletIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("wallets")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_wallet_user_guild"),
		},
	})
	return err
}

func ensureCooldownIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("cooldowns")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}, {Key: "action", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_cooldown_user_guild_action"),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("idx_cooldown_ttl"),
		},
	})
	return err
}

func ensureMenuSessionIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("menu_sessions")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sessionId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_session_id"),
		},
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetName("idx_session_user_guild"),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("idx_session_ttl"),
		},
	})
	return err
}
