// File: internal/game/item/mongo_repository.go
package item

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
)

type mongoItemRepo struct{ col *mongo.Collection }

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoItemRepo{col: db.Collection("item_instances")}
}

func (r *mongoItemRepo) CreateInstance(ctx context.Context, inst *ItemInstance) error {
	inst.CreatedAt = time.Now().UTC()
	inst.UpdatedAt = inst.CreatedAt
	_, err := r.col.InsertOne(ctx, inst)
	return err
}

func (r *mongoItemRepo) GetInstancesByUser(ctx context.Context, userID, guildID string) ([]*ItemInstance, error) {
	cursor, err := r.col.Find(ctx, bson.M{"userId": userID, "guildId": guildID})
	if err != nil {
		return nil, err
	}
	var items []*ItemInstance
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *mongoItemRepo) GetInstanceByID(ctx context.Context, instanceID, userID, guildID string) (*ItemInstance, error) {
	var item ItemInstance
	err := r.col.FindOne(ctx, bson.M{"instanceId": instanceID, "userId": userID, "guildId": guildID}).Decode(&item)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, apperrors.ErrItemNotFound
	}
	return &item, err
}

func (r *mongoItemRepo) AdjustQuantity(ctx context.Context, instanceID, userID, guildID string, amount int64) error {
	filter := bson.M{"instanceId": instanceID, "userId": userID, "guildId": guildID}
	if amount < 0 {
		filter["quantity"] = bson.M{"$gte": -amount}
	}
	update := bson.M{
		"$inc": bson.M{"quantity": amount},
		"$set": bson.M{"updatedAt": time.Now().UTC()},
	}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		if amount < 0 {
			return apperrors.ErrInsufficientItemQuantity
		}
		return apperrors.ErrItemNotFound
	}
	return nil
}

func (r *mongoItemRepo) DeleteInstance(ctx context.Context, instanceID, userID, guildID string) error {
	// Chỉ xóa khi quantity <= 0 (Atomic prevent accidental deletion of replenished stacks)
	_, err := r.col.DeleteOne(ctx, bson.M{"instanceId": instanceID, "userId": userID, "guildId": guildID, "quantity": bson.M{"$lte": 0}})
	return err
}

func (r *mongoItemRepo) UpdateMetadata(ctx context.Context, instanceID, userID, guildID string, metadata map[string]interface{}) error {
	filter := bson.M{"instanceId": instanceID, "userId": userID, "guildId": guildID}
	update := bson.M{"$set": bson.M{"metadata": metadata, "updatedAt": time.Now().UTC()}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
