// File: internal/discord/menu/session.go
// Version: v0.1
// Purpose: MenuSession model and repository — track which user owns which active menu message.
// Security: sessionId is cryptographically random. Every button interaction validates sessionId
//           and userId to prevent users from operating each other's menus.
// Notes: MongoDB TTL index on expiresAt auto-expires sessions. See database/indexes.go.

package menu

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
	"github.com/yourname/tu-tien-bot/pkg/utils"
)

const sessionCollection = "menu_sessions"

// Page constants identify which page of the menu is currently shown.
type Page string

const (
	PageMain        Page = "main"
	PageProfile     Page = "profile"
	PageCultivation Page = "cultivation"
	PageInventory   Page = "inventory"    // TODO v0.3
	PageSkills      Page = "skills"       // TODO v0.4
	PagePets        Page = "pets"         // TODO v0.6
	PageGacha       Page = "gacha"        // TODO v0.5
	PageMarket      Page = "market"       // TODO v0.8
	PageSect        Page = "sect"         // TODO v1.0
)

// ParentOf returns the parent page for back-navigation.
var ParentOf = map[Page]Page{
	PageProfile:     PageMain,
	PageCultivation: PageMain,
	PageInventory:   PageMain,
	PageSkills:      PageMain,
	PagePets:        PageMain,
	PageGacha:       PageMain,
	PageMarket:      PageMain,
	PageSect:        PageMain,
}

// Session is one active menu instance per user.
type Session struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"    json:"id"`
	SessionID       string             `bson:"sessionId"        json:"sessionId"`       // cryptographic random ID
	UserID          string             `bson:"userId"           json:"userId"`
	GuildID         string             `bson:"guildId"          json:"guildId"`
	ChannelID       string             `bson:"channelId"        json:"channelId"`
	MessageID       string             `bson:"messageId"        json:"messageId"`         // Discord message ID of the menu
	CurrentPage     Page               `bson:"currentPage"      json:"currentPage"`
	CurrentCategory string             `bson:"currentCategory"  json:"currentCategory"`  // sub-navigation within a page
	ExpiresAt       time.Time          `bson:"expiresAt"        json:"expiresAt"`
	CreatedAt       time.Time          `bson:"createdAt"        json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"        json:"updatedAt"`
}

// IsExpired returns true if the session has exceeded its TTL.
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}

// OwnedBy returns true if this session belongs to the given user.
func (s *Session) OwnedBy(userID string) bool {
	return s.UserID == userID
}

// Repository defines data access for menu sessions.
type Repository interface {
	Create(ctx context.Context, session *Session) error
	FindBySessionID(ctx context.Context, sessionID string) (*Session, error)
	FindActiveByUser(ctx context.Context, userID, guildID string) (*Session, error)
	UpdatePage(ctx context.Context, sessionID string, page Page, category string) error
	UpdateMessageID(ctx context.Context, sessionID, messageID string) error
	Refresh(ctx context.Context, sessionID string, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
	DeleteExpiredByUser(ctx context.Context, userID, guildID string) error
}

type mongoSessionRepo struct {
	col *mongo.Collection
}

// NewRepository creates a MongoDB-backed session repository.
func NewRepository(db *mongo.Database) Repository {
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
		return nil, fmt.Errorf("%w: session sessionId=%s", apperrors.ErrNotFound, sessionID)
	}
	if err != nil {
		return nil, fmt.Errorf("session.FindBySessionID: %w", err)
	}
	return &s, nil
}

func (r *mongoSessionRepo) FindActiveByUser(ctx context.Context, userID, guildID string) (*Session, error) {
	filter := bson.M{
		"userId":    userID,
		"guildId":   guildID,
		"expiresAt": bson.M{"$gt": time.Now().UTC()},
	}
	var s Session
	err := r.col.FindOne(ctx, filter, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&s)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: active session userId=%s", apperrors.ErrNotFound, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("session.FindActiveByUser: %w", err)
	}
	return &s, nil
}

func (r *mongoSessionRepo) UpdatePage(ctx context.Context, sessionID string, page Page, category string) error {
	update := bson.M{"$set": bson.M{
		"currentPage":     page,
		"currentCategory": category,
		"updatedAt":       time.Now().UTC(),
	}}
	_, err := r.col.UpdateOne(ctx, bson.M{"sessionId": sessionID}, update)
	return err
}

func (r *mongoSessionRepo) UpdateMessageID(ctx context.Context, sessionID, messageID string) error {
	update := bson.M{"$set": bson.M{
		"messageId": messageID,
		"updatedAt": time.Now().UTC(),
	}}
	_, err := r.col.UpdateOne(ctx, bson.M{"sessionId": sessionID}, update)
	return err
}

func (r *mongoSessionRepo) Refresh(ctx context.Context, sessionID string, ttl time.Duration) error {
	update := bson.M{"$set": bson.M{
		"expiresAt": time.Now().UTC().Add(ttl),
		"updatedAt": time.Now().UTC(),
	}}
	_, err := r.col.UpdateOne(ctx, bson.M{"sessionId": sessionID}, update)
	return err
}

func (r *mongoSessionRepo) Delete(ctx context.Context, sessionID string) error {
	_, err := r.col.DeleteOne(ctx, bson.M{"sessionId": sessionID})
	return err
}

func (r *mongoSessionRepo) DeleteExpiredByUser(ctx context.Context, userID, guildID string) error {
	filter := bson.M{
		"userId":    userID,
		"guildId":   guildID,
		"expiresAt": bson.M{"$lte": time.Now().UTC()},
	}
	_, err := r.col.DeleteMany(ctx, filter)
	return err
}
