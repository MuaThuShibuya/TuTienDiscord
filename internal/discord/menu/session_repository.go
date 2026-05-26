// File: internal/discord/menu/session_repository.go
// Chức năng: Interface SessionRepository và MongoDB implementation.
// Ghi chú: MongoDB TTL index trên expiresAt tự xóa session hết hạn.

package menu

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

const sessionCollection = "menu_sessions"

// SessionRepository định nghĩa các thao tác lưu trữ cho Session.
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindBySessionID(ctx context.Context, sessionID string) (*Session, error)
	UpdatePage(ctx context.Context, sessionID string, page Page, category string) error
	UpdateMessageID(ctx context.Context, sessionID, messageID string) error
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
	DeleteExpiredByUser(ctx context.Context, userID, guildID string) error
}

type mongoSessionRepo struct {
	col *mongo.Collection
}

// NewSessionRepository tạo session repository MongoDB.
func NewSessionRepository(db *mongo.Database) SessionRepository {
	return &mongoSessionRepo{col: db.Collection(sessionCollection)}
}

func (r *mongoSessionRepo) Create(ctx context.Context, session *Session) error {
	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}
	if session.SessionID == "" {
		session.SessionID = utils.NewSessionID()
	}
	now := time.Now().UTC()
	session.CreatedAt = now
	session.UpdatedAt = now

	_, err := r.col.InsertOne(ctx, session)
	if err != nil {
		return fmt.Errorf("session.Create: %w", err)
	}
	return nil
}

func (r *mongoSessionRepo) FindBySessionID(ctx context.Context, sessionID string) (*Session, error) {
	var s Session
	err := r.col.FindOne(ctx, bson.M{"sessionId": sessionID}).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: sessionId=%s", apperrors.ErrNotFound, sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("session.FindBySessionID: %w", err)
	}
	return &s, nil
}

func (r *mongoSessionRepo) UpdatePage(ctx context.Context, sessionID string, page Page, category string) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"sessionId": sessionID},
		bson.M{"$set": bson.M{
			"currentPage":     page,
			"currentCategory": category,
			"updatedAt":       time.Now().UTC(),
		}},
	)
	return err
}

func (r *mongoSessionRepo) UpdateMessageID(ctx context.Context, sessionID, messageID string) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"sessionId": sessionID},
		bson.M{"$set": bson.M{"messageId": messageID, "updatedAt": time.Now().UTC()}},
	)
	return err
}

func (r *mongoSessionRepo) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	_, err := r.col.UpdateOne(ctx,
		bson.M{"sessionId": sessionID},
		bson.M{"$set": bson.M{
			"expiresAt": time.Now().UTC().Add(ttl),
			"updatedAt": time.Now().UTC(),
		}},
	)
	return err
}

func (r *mongoSessionRepo) Delete(ctx context.Context, sessionID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"sessionId": sessionID})
	return err
}

func (r *mongoSessionRepo) DeleteExpiredByUser(ctx context.Context, userID, guildID string) error {
	_, err := r.col.DeleteMany(ctx, bson.M{
		"userId":    userID,
		"guildId":   guildID,
		"expiresAt": bson.M{"$lte": time.Now().UTC()},
	})
	return err
}
