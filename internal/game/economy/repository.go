// File: internal/game/economy/repository.go
// Version: v0.1
// Purpose: Data access layer for Wallet — find, upsert, and atomic currency operations.
// Security: Currency changes use atomic $inc to prevent race conditions in v0.2+.
//           Never allow negative balances via direct $set.
// Notes: AdjustCurrency uses atomic findOneAndUpdate to safely change balances.

package economy

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

const collectionName = "wallets"

// Repository defines data access for player wallets.
type Repository interface {
	FindByUserID(ctx context.Context, userID, guildID string) (*Wallet, error)
	Upsert(ctx context.Context, wallet *Wallet) error
	// AdjustCurrency atomically adds (or subtracts) a currency amount.
	// amount can be negative for deductions. Will return ErrInsufficientFunds if balance < 0 after deduction.
	AdjustSpiritStones(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error)
	AdjustSpiritJades(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error)
	AdjustFateTickets(ctx context.Context, userID, guildID string, amount int) (*Wallet, error)
}

type mongoRepository struct {
	col *mongo.Collection
}

// NewRepository creates a new MongoDB-backed economy repository.
func NewRepository(db *mongo.Database) Repository {
	return &mongoRepository{col: db.Collection(collectionName)}
}

func (r *mongoRepository) FindByUserID(ctx context.Context, userID, guildID string) (*Wallet, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	var wallet Wallet
	err := r.col.FindOne(ctx, filter).Decode(&wallet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: wallet userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	if err != nil {
		return nil, fmt.Errorf("economy.FindByUserID: %w", err)
	}
	return &wallet, nil
}

func (r *mongoRepository) Upsert(ctx context.Context, wallet *Wallet) error {
	if wallet.ID.IsZero() {
		wallet.ID = primitive.NewObjectID()
	}
	now := time.Now().UTC()
	if wallet.CreatedAt.IsZero() {
		wallet.CreatedAt = now
	}
	wallet.UpdatedAt = now

	filter := bson.M{"userId": wallet.UserID, "guildId": wallet.GuildID}
	update := bson.M{"$set": wallet}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("economy.Upsert: %w", err)
	}
	return nil
}

func (r *mongoRepository) AdjustSpiritStones(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error) {
	return r.adjustCurrency(ctx, userID, guildID, "spiritStones", amount)
}

func (r *mongoRepository) AdjustSpiritJades(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error) {
	return r.adjustCurrency(ctx, userID, guildID, "spiritJades", amount)
}

func (r *mongoRepository) AdjustFateTickets(ctx context.Context, userID, guildID string, amount int) (*Wallet, error) {
	wallet, err := r.adjustCurrency(ctx, userID, guildID, "fateTickets", int64(amount))
	return wallet, err
}

// adjustCurrency uses atomic $inc; for deductions it also checks the resulting balance >= 0.
func (r *mongoRepository) adjustCurrency(ctx context.Context, userID, guildID, field string, amount int64) (*Wallet, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}

	// For deductions, add a minimum value guard so balance never goes negative.
	if amount < 0 {
		filter[field] = bson.M{"$gte": -amount}
	}

	update := bson.M{
		"$inc": bson.M{field: amount},
		"$set": bson.M{"updatedAt": time.Now().UTC()},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated Wallet
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if errors.Is(err, mongo.ErrNoDocuments) {
		if amount < 0 {
			return nil, fmt.Errorf("%w: field=%s userId=%s", apperrors.ErrInsufficientFunds, field, userID)
		}
		return nil, fmt.Errorf("%w: wallet userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	if err != nil {
		return nil, fmt.Errorf("economy.adjustCurrency field=%s: %w", field, err)
	}
	return &updated, nil
}
