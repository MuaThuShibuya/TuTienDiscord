// File: internal/game/cooldown/repository.go
// Version: v0.1
// Purpose: Data access for cooldowns — check, set, and delete per-action cooldowns.
// Security: action is always a server-side constant, never user-provided text.
// Notes: MongoDB TTL index handles expiry cleanup. No manual deletion needed in most cases.

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

	apperrors "github.com/yourname/tu-tien-bot/internal/errors"
)

const collectionName = "cooldowns"

// Repository defines data access for cooldowns.
type Repository interface {
	Get(ctx context.Context, userID, guildID string, action Action) (*Cooldown, error)
	Set(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error
	Delete(ctx context.Context, userID, guildID string, action Action) error
}

type mongoRepository struct {
	col *mongo.Collection
}

// NewRepository creates a new MongoDB-backed cooldown repository.
func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection(collectionName)}
}

func (r *mongoRepository) Get(ctx context.Context, userID, guildID string, action Action) (*Cooldown, error) {
	filter := bson.M{
		"userId":  userID,
		"guildId": guildID,
		"action":  action,
		"expiresAt": bson.M{"$gt": time.Now().UTC()}, // Only return active cooldowns
	}
	var cd Cooldown
	err := r.col.FindOne(ctx, filter).Decode(&cd)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: no active cooldown for action=%s userId=%s", apperrors.ErrNotFound, action, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("cooldown.Get: %w", err)
	}
	return &cd, nil
}

func (r *mongoRepository) Set(ctx context.Context, userID, guildID string, action Action, duration time.Duration) error {
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
	update := bson.M{"$set": cd}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("cooldown.Set action=%s: %w", action, err)
	}
	return nil
}

func (r *mongoRepository) Delete(ctx context.Context, userID, guildID string, action Action) error {
	filter := bson.M{"userId": userID, "guildId": guildID, "action": action}
	_, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("cooldown.Delete action=%s: %w", action, err)
	}
	return nil
}
