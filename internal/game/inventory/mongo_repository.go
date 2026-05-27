// File: internal/game/inventory/mongo_repository.go
package inventory

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoInventoryRepo struct{ col *mongo.Collection }

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoInventoryRepo{col: db.Collection("inventories")}
}

func (r *mongoInventoryRepo) GetOrCreate(ctx context.Context, userID, guildID string) (*Inventory, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{
		"$setOnInsert": bson.M{
			"userId":         userID,
			"guildId":        guildID,
			"slotLimit":      50,
			"slotUsage":      0,
			"starterGranted": false,
			"createdAt":      time.Now().UTC(),
			"updatedAt":      time.Now().UTC(),
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var inv Inventory
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&inv)
	return &inv, err
}

func (r *mongoInventoryRepo) MarkStarterGranted(ctx context.Context, userID, guildID string) error {
	// THÊM: Filter kiểm tra starterGranted == false để đảm bảo Atomic
	filter := bson.M{"userId": userID, "guildId": guildID, "starterGranted": false}
	update := bson.M{"$set": bson.M{"starterGranted": true, "updatedAt": time.Now().UTC()}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("already granted or not found")
	}
	return nil
}

func (r *mongoInventoryRepo) AcquireSlot(ctx context.Context, userID, guildID string) error {
	// So sánh trực tiếp slotUsage < slotLimit trên DB (Atomic)
	filter := bson.M{"userId": userID, "guildId": guildID, "$expr": bson.M{"$lt": []interface{}{"$slotUsage", "$slotLimit"}}}
	update := bson.M{"$inc": bson.M{"slotUsage": 1}, "$set": bson.M{"updatedAt": time.Now().UTC()}}
	res, err := r.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("không còn ô trống trong túi đồ")
	}
	return nil
}

func (r *mongoInventoryRepo) ReleaseSlot(ctx context.Context, userID, guildID string) error {
	filter := bson.M{"userId": userID, "guildId": guildID, "slotUsage": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"slotUsage": -1}, "$set": bson.M{"updatedAt": time.Now().UTC()}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
