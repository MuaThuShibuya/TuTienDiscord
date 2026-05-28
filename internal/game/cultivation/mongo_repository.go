// File: internal/game/cultivation/mongo_repository.go
// Phiên bản: v0.1.1
// Mục đích: Implementation MongoDB cho cultivation Repository.
//           Dùng upsert để tránh race condition khi tạo profile lần đầu.
// Bảo mật: Query phải có guildId filter. Không gọi Discord API.

package cultivation

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

const collectionCultivation = "cultivation_profiles"

// mongoCultivationRepo là MongoDB implementation của Repository.
type mongoCultivationRepo struct {
	col *mongo.Collection
}

// NewMongoRepository tạo cultivation repository MongoDB.
func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoCultivationRepo{col: db.Collection(collectionCultivation)}
}

func (r *mongoCultivationRepo) FindByUserID(ctx context.Context, userID, guildID string) (*CultivationProfile, error) {
	var profile CultivationProfile
	filter := bson.M{"userId": userID}
	if guildID != "" {
		filter["guildId"] = guildID
	}
	err := r.col.FindOne(ctx, filter).Decode(&profile)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: cultivation userId=%s", apperrors.ErrNotFound, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("cultivation.FindByUserID: %w", err)
	}
	return &profile, nil
}

func (r *mongoCultivationRepo) Upsert(ctx context.Context, profile *CultivationProfile) error {
	now := time.Now().UTC()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.UpdatedAt = now

	filter := bson.M{"userId": profile.UserID, "guildId": profile.GuildID}
	update := bson.M{"$set": profile}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("cultivation.Upsert: %w", err)
	}
	return nil
}

func (r *mongoCultivationRepo) UpdateStats(ctx context.Context, p *CultivationProfile) error {
	p.UpdatedAt = time.Now().UTC()

	filter := bson.M{"userId": p.UserID, "guildId": p.GuildID}
	// Dùng $set để update một phần thay vì toàn bộ document, an toàn không đè CreatedAt
	update := bson.M{"$set": bson.M{
		"realm":                  p.Realm,
		"realmLevel":             p.RealmLevel,
		"cultivationExp":         p.CultivationExp,
		"cultivationExpRequired": p.CultivationExpRequired,
		"combatPower":            p.CombatPower,
		"stamina":                p.Stamina,
		"mindState":              p.MindState,
		"path":                   p.Path,
		"updatedAt":              p.UpdatedAt,
	}}

	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("cultivation.UpdateStats: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("%w: cultivation update failed", apperrors.ErrNotFound)
	}
	return nil
}
