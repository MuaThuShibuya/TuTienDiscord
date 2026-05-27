// File: internal/game/combat/mongo_repository.go
package combat

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
)

type mongoCombatRepo struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	// TODO: Cần chạy script để tạo các index:
	// 1. _id: unique
	// 2. (userId, state): partial unique index cho state="active" để ngăn multiple active sessions.
	// 3. expiresAt: TTL index để dọn dẹp các session bị treo (expireAfterSeconds: 0)
	return &mongoCombatRepo{col: db.Collection("combat_sessions")}
}

func (r *mongoCombatRepo) CreateSession(ctx context.Context, session *CombatSession) error {
	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = session.CreatedAt

	_, err := r.col.InsertOne(ctx, session)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrCombatSessionAlreadyActive
		}
		return err
	}
	return nil
}

func (r *mongoCombatRepo) GetSession(ctx context.Context, sessionID string) (*CombatSession, error) {
	var session CombatSession
	err := r.col.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.ErrNotFound
	}
	return &session, err
}

func (r *mongoCombatRepo) GetActiveSessionByUser(ctx context.Context, userID string) (*CombatSession, error) {
	var session CombatSession
	filter := bson.M{
		"userId":    userID,
		"state":     StateActive,
		"expiresAt": bson.M{"$gt": time.Now().UTC()}, // Chỉ lấy nếu chưa hết hạn
	}
	err := r.col.FindOne(ctx, filter).Decode(&session)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.ErrNotFound
	}
	return &session, err
}

func (r *mongoCombatRepo) UpdateSession(ctx context.Context, session *CombatSession) error {
	session.UpdatedAt = time.Now().UTC()
	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": session}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err == nil && res.MatchedCount == 0 {
		return apperrors.ErrNotFound
	}
	return err
}

func (r *mongoCombatRepo) MarkSessionState(ctx context.Context, sessionID string, state SessionState) error {
	filter := bson.M{"_id": sessionID}
	update := bson.M{"$set": bson.M{"state": state, "updatedAt": time.Now().UTC()}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
