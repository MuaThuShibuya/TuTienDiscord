// File: internal/database/indexes.go
// Phiên bản: v0.1.2
// Mục đích: Tạo và đảm bảo các MongoDB index cho toàn bộ collection khi khởi động.
// Bảo mật: Không liên quan đến secret. Tạo index là thao tác idempotent, an toàn khi gọi lại nhiều lần.
// Ghi chú: Thêm index mới vào đây khi bổ sung collection trong các phiên bản sau.

package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/whiskey/tu-tien-bot/internal/logger"
)

const indexTimeout = 30 * time.Second

// EnsureIndexes tạo toàn bộ index cần thiết cho ứng dụng.
// Idempotent — gọi lại nhiều lần không gây lỗi nếu index đã tồn tại.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	log := logger.L()

	indexCtx, cancel := context.WithTimeout(ctx, indexTimeout)
	defer cancel()

	var errs []error

	if err := ensurePlayerIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection players", zap.Error(err))
	}

	if err := ensureCultivationIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection cultivation_profiles", zap.Error(err))
	}

	if err := ensureWalletIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection wallets", zap.Error(err))
	}

	if err := ensureCooldownIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection cooldowns", zap.Error(err))
	}

	if err := ensureMenuSessionIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection menu_sessions", zap.Error(err))
	}

	if err := ensureInventoryIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection inventories", zap.Error(err))
	}

	if err := ensureItemIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection items", zap.Error(err))
	}

	if err := ensureEquipmentIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection equipment_sets", zap.Error(err))
	}

	if err := ensureAlchemyIndexes(indexCtx, db); err != nil {
		errs = append(errs, err)
		log.Error("Không tạo được index cho collection alchemy_profiles", zap.Error(err))
	}

	if len(errs) > 0 {
		// Trả về lỗi đầu tiên — caller (main.go) sẽ Fatal và không khởi động
		return errs[0]
	}

	log.Info("Đã đảm bảo MongoDB indexes")
	return nil
}

// ensurePlayerIndexes tạo index cho collection players.
// - Unique compound (userId + guildId): ngăn tạo trùng người chơi trong cùng server.
// - Index lastActiveAt: truy vấn nhanh người chơi hoạt động gần đây.
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

// ensureAlchemyIndexes tạo index cho collection alchemy_profiles.
// - Unique compound (userId + guildId).
func ensureAlchemyIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("alchemy_profiles")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_alchemy_user_guild"),
		},
	})
	return err
}

// ensureEquipmentIndexes tạo index cho collection equipment_sets.
func ensureEquipmentIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("equipment_sets")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_equipment_user_guild"),
		},
	})
	return err
}

// ensureInventoryIndexes tạo index cho collection inventories.
// - Unique compound (userId + guildId): mỗi người chơi chỉ có 1 túi đồ mỗi server.
func ensureInventoryIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("inventories")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_inventory_user_guild"),
		},
	})
	return err
}

// ensureItemIndexes tạo index cho collection items (chứa item instances).
// - Compound (userId + guildId): Truy xuất nhanh toàn bộ vật phẩm trong túi của 1 người.
// - Unique (instanceId): Định danh duy nhất cho từng món đồ/stack đồ.
func ensureItemIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("item_instances") // Khớp với tên dùng trong item/mongo_repository.go
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}},
			Options: options.Index().SetName("idx_item_user_guild"),
		},
		{
			Keys:    bson.D{{Key: "instanceId", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_item_instance_id"),
		},
	})
	return err
}

// ensureCultivationIndexes tạo index cho collection cultivation_profiles.
// - Unique compound (userId + guildId): mỗi người chơi chỉ có 1 hồ sơ tu luyện mỗi server.
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

// ensureWalletIndexes tạo index cho collection wallets.
// - Unique compound (userId + guildId): mỗi người chơi chỉ có 1 ví mỗi server.
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

// ensureCooldownIndexes tạo index cho collection cooldowns.
// - Unique compound (userId + guildId + action): mỗi hành động chỉ có 1 cooldown mỗi người chơi.
// - TTL index trên expiresAt: MongoDB tự động xóa document đã hết hạn (tránh rác dữ liệu).
func ensureCooldownIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection("cooldowns")
	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "guildId", Value: 1}, {Key: "action", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_cooldown_user_guild_action"),
		},
		{
			// TTL index: MongoDB daemon xóa document khi expiresAt < now
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("idx_cooldown_ttl"),
		},
	})
	return err
}

// ensureMenuSessionIndexes tạo index cho collection menu_sessions.
// - Unique sessionId: mỗi phiên menu có ID duy nhất (dùng để validate ownership).
// - Compound (userId + guildId): tìm nhanh phiên theo người chơi.
// - TTL index trên expiresAt: tự động xóa phiên hết hạn sau MENU_SESSION_TTL_MINUTES.
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
			// TTL index: phiên menu tự xóa sau khi hết hạn
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("idx_session_ttl"),
		},
	})
	return err
}
