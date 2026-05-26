// File: internal/game/alchemy/mongo_repository.go
// Chức năng: Impl interface Repository cho MongoDB.

package alchemy

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionAlchemy = "alchemy_profiles"

type mongoAlchemyRepo struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoAlchemyRepo{col: db.Collection(collectionAlchemy)}
}

func (r *mongoAlchemyRepo) Get(ctx context.Context, userID, guildID string) (*AlchemyProfile, error) {
	var profile AlchemyProfile
	err := r.col.FindOne(ctx, bson.M{"userId": userID, "guildId": guildID}).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Nếu chưa có, trả về profile mặc định cấp 1
			return &AlchemyProfile{UserID: userID, GuildID: guildID, Level: 1, Exp: 0}, nil
		}
		return nil, err
	}
	return &profile, nil
}

func (r *mongoAlchemyRepo) Upsert(ctx context.Context, profile *AlchemyProfile) error {
	now := time.Now().UTC()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.UpdatedAt = now

	filter := bson.M{"userId": profile.UserID, "guildId": profile.GuildID}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": profile}, options.Update().SetUpsert(true))
	return err
}
