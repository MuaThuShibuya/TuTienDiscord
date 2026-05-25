// File: internal/game/cultivation/repository.go
// Version: v0.1
// Purpose: Data access layer for CultivationProfile — find and upsert cultivation data.
// Security: All queries scoped to (userId, guildId). No cross-user or cross-guild access.
// Notes: Uses atomic $set updates to avoid overwriting fields not in scope.

package cultivation

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

const collectionName = "cultivation_profiles"

// Repository defines data access for cultivation profiles.
type Repository interface {
	FindByUserID(ctx context.Context, userID, guildID string) (*CultivationProfile, error)
	Upsert(ctx context.Context, profile *CultivationProfile) error
	// TODO v0.2: AddCultivationExp, SetRealm, SetMindState
}

type mongoRepository struct {
	col *mongo.Collection
}

// NewRepository creates a new MongoDB-backed cultivation repository.
func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection(collectionName)}
}

func (r *mongoRepository) FindByUserID(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	var profile CultivationProfile
	err := r.col.FindOne(ctx, filter).Decode(&profile)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: cultivation userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	if err != nil {
		return nil, fmt.Errorf("cultivation.FindByUserID: %w", err)
	}
	return &profile, nil
}

func (r *mongoRepository) Upsert(ctx context.Context, profile *CultivationProfile) error {
	if profile.ID.IsZero() {
		profile.ID = primitive.NewObjectID()
	}
	now := time.Now().UTC()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.UpdatedAt = now

	filter := bson.M{"userId": profile.UserID, "guildId": profile.GuildID}
	update := bson.M{"$set": profile}
	opts := options.Update().SetUpsert(true)

	_, err := r.col.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("cultivation.Upsert: %w", err)
	}
	return nil
}
