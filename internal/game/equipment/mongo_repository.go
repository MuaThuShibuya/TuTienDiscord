// File: internal/game/equipment/mongo_repository.go
package equipment

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoEquipRepo struct{ col *mongo.Collection }

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoEquipRepo{col: db.Collection("equipment_sets")}
}

func (r *mongoEquipRepo) Get(ctx context.Context, userID, guildID string) (*EquipmentSet, error) {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{
		"$setOnInsert": bson.M{
			"userId":    userID,
			"guildId":   guildID,
			"slots":     bson.M{},
			"createdAt": time.Now().UTC(),
			"updatedAt": time.Now().UTC(),
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	var eq EquipmentSet
	err := r.col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&eq)
	return &eq, err
}

func (r *mongoEquipRepo) Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{"$set": bson.M{"slots." + string(slot): instanceID, "updatedAt": time.Now().UTC()}}
	_, err := r.col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *mongoEquipRepo) Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error {
	filter := bson.M{"userId": userID, "guildId": guildID}
	update := bson.M{"$unset": bson.M{"slots." + string(slot): ""}, "$set": bson.M{"updatedAt": time.Now().UTC()}}
	_, err := r.col.UpdateOne(ctx, filter, update)
	return err
}
