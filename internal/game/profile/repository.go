// File: internal/game/profile/repository.go
// Version: v0.1
// Purpose: Data access layer for the Player model. Only talks to MongoDB, no business logic.
// Security: All queries use (userId, guildId) filter to prevent cross-guild data leaks.
// Notes: Use database.NewContext() for every operation to enforce timeouts.

package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
)

const collectionName = "players"

// Repository defines the data access interface for Player records.
type Repository interface {
	FindByUserID(ctx context.Context, userID, guildID string) (*Player, error)
	Create(ctx context.Context, player *Player) error
	UpdateLastActive(ctx context.Context, userID, guildID string) error
	UpdateDaoName(ctx context.Context, userID, guildID, daoName string) error
}

type mongoRepository struct {
	col *mongo.Collection
}

// NewRepository creates a new MongoDB-backed player repository.
func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection(collectionName)}
}

func (r *mongoRepository) FindByUserID(ctx context.Context, userID, guildID string) (*Player, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	var player Player
	err := r.col.FindOne(ctx, filter).Decode(&player)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: player userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	if err != nil {
		return nil, fmt.Errorf("profile.FindByUserID: %w", err)
	}
	return &player, nil
}

func (r *mongoRepository) Create(ctx context.Context, player *Player) error {
	if player.ID.IsZero() {
		player.ID = primitive.NewObjectID()
	}
	now := time.Now().UTC()
	player.CreatedAt = now
	player.UpdatedAt = now
	player.LastActiveAt = now

	_, err := r.col.InsertOne(ctx, player)
	if mongo.IsDuplicateKeyError(err) {
		return fmt.Errorf("%w: player userId=%s guildId=%s", apperrors.ErrAlreadyExists, player.UserID, player.GuildID)
	}
	if err != nil {
		return fmt.Errorf("profile.Create: %w", err)
	}
	return nil
}

func (r *mongoRepository) UpdateLastActive(ctx context.Context, userID, guildID string) error {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{"$set": bson.M{
		"lastActiveAt": time.Now().UTC(),
		"updatedAt":    time.Now().UTC(),
	}}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return fmt.Errorf("profile.UpdateLastActive: %w", err)
	}
	return nil
}

func (r *mongoRepository) UpdateDaoName(ctx context.Context, userID, guildID, daoName string) error {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{"$set": bson.M{
		"daoName":   daoName,
		"updatedAt": time.Now().UTC(),
	}}
	result, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("profile.UpdateDaoName: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("%w: player userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	return nil
}
