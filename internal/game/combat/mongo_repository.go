// File: internal/game/combat/mongo_repository.go
package combat

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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
	filter := bson.M{"_id": session.ID, "userId": session.UserID} // Thêm userId để đảm bảo chỉ chủ sở hữu mới update được
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

func (r *mongoCombatRepo) TryStartRewardClaim(ctx context.Context, sessionID string, claimID string, now time.Time) (*CombatSession, error) {
	filter := bson.M{
		"_id":           sessionID,
		"rewardClaimed": bson.M{"$ne": true},
		"$or": []bson.M{
			{"rewardClaimStatus": bson.M{"$exists": false}},
			{"rewardClaimStatus": ""},
			{"rewardClaimStatus": "pending"},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"rewardClaimStatus":    "claiming",
			"rewardClaimId":        claimID,
			"rewardClaimError":     "",
			"rewardClaimStartedAt": now,
			"updatedAt":            now,
		},
	}
	var session CombatSession
	err := r.col.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&session)
	if errors.Is(err, mongo.ErrNoDocuments) {
		var current CombatSession
		if findErr := r.col.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&current); findErr == nil {
			if current.RewardClaimed || current.RewardClaimStatus == "claimed" {
				return nil, ErrRewardAlreadyClaimed
			}
			if current.RewardClaimStatus == "claiming" {
				return nil, ErrRewardClaimInProgress
			}
			if current.RewardClaimStatus == "claim_failed" {
				return nil, ErrRewardClaimFailedNeedsAdmin
			}
		}
		return nil, apperrors.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *mongoCombatRepo) CompleteRewardClaim(ctx context.Context, sessionID string, claimID string, details []ClaimedReward, now time.Time) error {
	filter := bson.M{"_id": sessionID, "rewardClaimStatus": "claiming", "rewardClaimId": claimID}
	update := bson.M{"$set": bson.M{
		"rewardClaimed": true, "rewardClaimStatus": "claimed", "rewardClaimedAt": now, "claimedRewards": details, "updatedAt": now,
	}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("session not found or invalid state for completion")
	}
	return nil
}

func (r *mongoCombatRepo) FailRewardClaim(ctx context.Context, sessionID string, claimID string, reason string, now time.Time) error {
	filter := bson.M{"_id": sessionID, "rewardClaimStatus": "claiming", "rewardClaimId": claimID}
	update := bson.M{"$set": bson.M{"rewardClaimStatus": "claim_failed", "rewardClaimError": reason, "updatedAt": now}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("session not found or invalid state for failure")
	}
	return nil
}
