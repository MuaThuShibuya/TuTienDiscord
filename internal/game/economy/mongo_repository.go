// File: internal/game/economy/mongo_repository.go
// Phiên bản: v0.1.1
// Mục đích: MongoDB implementation cho economy Repository.
//           Dùng $inc atomic cho mọi thay đổi tiền — tránh race condition.
// Bảo mật: Filter tiền âm: "field >= |amount|" trước khi $inc để không cho số dư âm.

package economy

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

const collectionWallets = "wallets"

type mongoEconomyRepo struct {
	col *mongo.Collection
}

// NewMongoRepository tạo economy repository MongoDB.
func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoEconomyRepo{col: db.Collection(collectionWallets)}
}

func (r *mongoEconomyRepo) FindByUserID(ctx context.Context, userID, guildID string) (*Wallet, error) {
	var wallet Wallet
	err := r.col.FindOne(ctx, bson.M{"userId": userID, "guildId": guildID}).Decode(&wallet)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("%w: wallet userId=%s", apperrors.ErrNotFound, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("economy.FindByUserID: %w", err)
	}
	return &wallet, nil
}

func (r *mongoEconomyRepo) Upsert(ctx context.Context, wallet *Wallet) error {
	now := time.Now().UTC()
	if wallet.CreatedAt.IsZero() {
		wallet.CreatedAt = now
	}
	wallet.UpdatedAt = now

	filter := bson.M{"userId": wallet.UserID, "guildId": wallet.GuildID}
	_, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": wallet}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("economy.Upsert: %w", err)
	}
	return nil
}

func (r *mongoEconomyRepo) AdjustSpiritStones(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error) {
	return r.atomicAdjust(ctx, userID, guildID, "spiritStones", amount)
}

func (r *mongoEconomyRepo) AdjustSpiritJades(ctx context.Context, userID, guildID string, amount int64) (*Wallet, error) {
	return r.atomicAdjust(ctx, userID, guildID, "spiritJades", amount)
}

func (r *mongoEconomyRepo) AdjustFateTickets(ctx context.Context, userID, guildID string, amount int) (*Wallet, error) {
	return r.atomicAdjust(ctx, userID, guildID, "fateTickets", int64(amount))
}

// atomicAdjust thực hiện $inc atomic trên một trường currency.
// Nếu amount âm, thêm filter "field >= |amount|" để chặn số dư âm.
func (r *mongoEconomyRepo) atomicAdjust(ctx context.Context, userID, guildID, field string, amount int64) (*Wallet, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	if amount < 0 {
		// Chặn số dư âm: chỉ update khi field >= |amount|
		filter[field] = bson.M{"$gte": -amount}
	}

	update := bson.M{
		"$inc": bson.M{field: amount},
		"$set": bson.M{"updatedAt": time.Now().UTC()},
	}

	var updated Wallet
	err := r.col.FindOneAndUpdate(ctx, filter, update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	if errors.Is(err, mongo.ErrNoDocuments) {
		if amount < 0 {
			return nil, fmt.Errorf("%w: field=%s userId=%s", apperrors.ErrInsufficientFunds, field, userID)
		}
		return nil, fmt.Errorf("%w: wallet userId=%s", apperrors.ErrNotFound, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("economy.atomicAdjust field=%s: %w", field, err)
	}
	return &updated, nil
}
