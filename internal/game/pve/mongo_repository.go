// File: internal/game/pve/mongo_repository.go

package pve

import (
	"context"
	"errors"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProgressRepo struct {
	col *mongo.Collection
}

func NewMongoProgressRepository(db *mongo.Database) ProgressRepository {
	return &mongoProgressRepo{col: db.Collection("pve_progress")}
}

func (r *mongoProgressRepo) GetProgress(ctx context.Context, userID string) (*UserPvEProgress, error) {
	var prog UserPvEProgress
	err := r.col.FindOne(ctx, bson.M{"userId": userID}).Decode(&prog)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return &UserPvEProgress{UserID: userID, Areas: make(map[string]AreaProgress)}, nil
	}
	return &prog, err
}

func (r *mongoProgressRepo) GetAreaProgress(ctx context.Context, userID, areaID string) (*AreaProgress, error) {
	prog, err := r.GetProgress(ctx, userID)
	if err != nil {
		return nil, err
	}
	if prog.Areas == nil {
		return nil, apperrors.ErrNotFound
	}
	areaProg, ok := prog.Areas[areaID]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return &areaProg, nil
}

func (r *mongoProgressRepo) UpsertAreaProgress(ctx context.Context, userID string, progress AreaProgress) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"areas." + progress.AreaID: progress,
			"updatedAt":                time.Now().UTC(),
		},
		"$setOnInsert": bson.M{"userId": userID},
	}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *mongoProgressRepo) MarkStageCleared(ctx context.Context, userID, areaID string, stage int) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$max": bson.M{"areas." + areaID + ".highestStageCleared": stage},
		"$set": bson.M{"areas." + areaID + ".updatedAt": time.Now().UTC(), "updatedAt": time.Now().UTC()},
	}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoProgressRepo) IncrementAttempt(ctx context.Context, userID, areaID string) error {
	filter := bson.M{"userId": userID}
	update := bson.M{
		"$inc": bson.M{"areas." + areaID + ".attemptsToday": 1},
		"$set": bson.M{
			"areas." + areaID + ".lastAttemptAt": time.Now().UTC(),
			"areas." + areaID + ".updatedAt":     time.Now().UTC(),
			"updatedAt":                          time.Now().UTC(),
		},
	}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
