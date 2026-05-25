// File: internal/game/profile/mongo_repository.go
// Phiên bản: v0.1.1
// Mục đích: Implementation MongoDB cho Repository interface của Player.
//           Chỉ thao tác với MongoDB — không chứa business logic.
// Bảo mật: Mọi query phải có filter guildId để tránh đọc dữ liệu chéo server.
// Ghi chú: Không gọi Discord API hay bất kỳ package discord nào.

package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
)

const collectionPlayers = "players"

// mongoRepository là implementation MongoDB của Repository.
type mongoRepository struct {
	col *mongo.Collection
}

// NewMongoRepository tạo repository MongoDB cho Player.
func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection(collectionPlayers)}
}

func (r *mongoRepository) FindByUserID(ctx context.Context, userID, guildID string) (*Player, error) {
	var player Player
	err := r.col.FindOne(ctx, bson.M{"userId": userID, "guildId": guildID}).Decode(&player)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	if err != nil {
		return nil, fmt.Errorf("profile.FindByUserID: %w", err)
	}
	return &player, nil
}

func (r *mongoRepository) Create(ctx context.Context, player *Player) error {
	now := time.Now().UTC()
	player.CreatedAt = now
	player.UpdatedAt = now
	player.LastActiveAt = now

	_, err := r.col.InsertOne(ctx, player)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("%w: userId=%s guildId=%s", apperrors.ErrAlreadyExists, player.UserID, player.GuildID)
	}
	if err != nil {
		return fmt.Errorf("profile.Create: %w", err)
	}
	return nil
}

func (r *mongoRepository) UpdateLastActive(ctx context.Context, userID, guildID string) error {
	now := time.Now().UTC()
	_, err := r.col.UpdateOne(ctx,
		bson.M{"userId": userID, "guildId": guildID},
		bson.M{"$set": bson.M{"lastActiveAt": now, "updatedAt": now}},
		options.Update().SetUpsert(false),
	)
	if err != nil {
		return fmt.Errorf("profile.UpdateLastActive: %w", err)
	}
	return nil
}

func (r *mongoRepository) UpdateDaoName(ctx context.Context, userID, guildID, daoName string) error {
	result, err := r.col.UpdateOne(ctx,
		bson.M{"userId": userID, "guildId": guildID},
		bson.M{"$set": bson.M{"daoName": daoName, "updatedAt": time.Now().UTC()}},
	)
	if err != nil {
		return fmt.Errorf("profile.UpdateDaoName: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	return nil
}
