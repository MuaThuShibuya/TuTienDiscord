// File: internal/game/cooldown/mongo_repository.go
// Phiên bản: v0.1.1
// Mục đích: MongoDB implementation cho cooldown Repository.
// Ghi chú: TTL index trong indexes.go tự động dọn expired cooldown — không cần xóa thủ công.

package cooldown

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
)

const collectionCooldowns = "cooldowns"

type mongoCooldownRepo struct {
	col *mongo.Collection
}

// NewMongoRepository tạo cooldown repository MongoDB.
func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoCooldownRepo{col: db.Collection(collectionCooldowns)}
}

func (r *mongoCooldownRepo) Get(ctx context.Context, userID, guildID string, action Action) (*Cooldown, error) {
	// Chỉ lấy cooldown còn hiệu lực (chưa hết hạn)
	filter := bson.M{
		"userId":    userID,
		"guildId":   guildID,
		"action":    action,
		"expiresAt": bson.M{"$gt": time.Now().UTC()},
	}
	var cd Cooldown
	err := r.col.FindOne(ctx, filter).Decode(&cd)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: không có cooldown action=%s userId=%s", apperrors.ErrNotFound, action, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("cooldown.Get: %w", err)
	}
	return &cd, nil
}

func (r *mongoCooldownRepo) Set(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error {
	now := time.Now().UTC()
	cd := &Cooldown{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		GuildID:   guildID,
		Action:    action,
		ExpiresAt: now.Add(duration),
		CreatedAt: now,
	}
	filter := bson.M{"userId": userID, "guildId": guildID, "action": action}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": cd}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("cooldown.Set action=%s: %w", action, err)
	}
	return nil
}

func (r *mongoCooldownRepo) Delete(ctx context.Context, userID, guildID string, action Action) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"userId": userID, "guildId": guildID, "action": action})
	if err != nil {
		return fmt.Errorf("cooldown.Delete action=%s: %w", action, err)
	}
	return nil
}
