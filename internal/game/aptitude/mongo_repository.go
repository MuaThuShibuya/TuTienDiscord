// File: internal/game/aptitude/mongo_repository.go
package aptitude

import (
	"context"
	"errors"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepo struct{ col *mongo.Collection }

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepo{col: db.Collection("aptitudes")}
}

func (r *mongoRepo) GetByUserID(ctx context.Context, userID string) (*AptitudeProfile, error) {
	var p AptitudeProfile
	err := r.col.FindOne(ctx, bson.M{"userId": userID}).Decode(&p)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *mongoRepo) Create(ctx context.Context, profile *AptitudeProfile) error {
	_, err := r.col.InsertOne(ctx, profile)
	if mongo.IsDuplicateKeyError(err) {
		return nil
	} // Idempotent
	return err
}
