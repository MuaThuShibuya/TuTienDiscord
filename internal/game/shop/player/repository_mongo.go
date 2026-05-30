package player

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrListingSoldOrNotFound = errors.New("phiếu đấu giá không tồn tại hoặc đã bị người khác mua")

type mongoRepo struct {
	col *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) Repository {
	return &mongoRepo{col: db.Collection("auction_listings")}
}

func (r *mongoRepo) CreateListing(ctx context.Context, listing *Listing) error {
	if listing.ID.IsZero() {
		listing.ID = primitive.NewObjectID()
	}
	_, err := r.col.InsertOne(ctx, listing)
	return err
}

func (r *mongoRepo) GetListing(ctx context.Context, listingID string) (*Listing, error) {
	objID, err := primitive.ObjectIDFromHex(listingID)
	if err != nil {
		return nil, ErrListingSoldOrNotFound
	}
	var listing Listing
	err = r.col.FindOne(ctx, bson.M{"_id": objID}).Decode(&listing)
	return &listing, err
}

func (r *mongoRepo) GetActiveListings(ctx context.Context, guildID string, limit, offset int) ([]*Listing, error) {
	filter := bson.M{"status": StatusActive}
	if guildID != "" {
		filter["guildId"] = guildID
	}
	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var listings []*Listing
	if err := cursor.All(ctx, &listings); err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *mongoRepo) GetUserActiveListings(ctx context.Context, userID, guildID string) ([]*Listing, error) {
	filter := bson.M{"status": StatusActive, "sellerId": userID}
	if guildID != "" {
		filter["guildId"] = guildID
	}
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var listings []*Listing
	if err := cursor.All(ctx, &listings); err != nil {
		return nil, err
	}
	return listings, nil
}

func (r *mongoRepo) CancelListing(ctx context.Context, listingID string, sellerID string) (*Listing, error) {
	objID, err := primitive.ObjectIDFromHex(listingID)
	if err != nil {
		return nil, ErrListingSoldOrNotFound
	}
	filter := bson.M{"_id": objID, "sellerId": sellerID, "status": StatusActive}
	update := bson.M{"$set": bson.M{"status": StatusCancelled}}
	var listing Listing
	err = r.col.FindOneAndUpdate(ctx, filter, update).Decode(&listing)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrListingSoldOrNotFound
		}
		return nil, err
	}
	return &listing, nil
}

func (r *mongoRepo) AtomicPurchase(ctx context.Context, listingID string, buyerID string) error {
	objID, err := primitive.ObjectIDFromHex(listingID)
	if err != nil {
		return ErrListingSoldOrNotFound
	}

	// FILTER KHẮT KHE: Chỉ cập nhật nếu status đang là active
	filter := bson.M{"_id": objID, "status": StatusActive}
	update := bson.M{
		"$set": bson.M{"status": StatusSold, "buyerId": buyerID, "purchasedAt": time.Now().UTC()},
	}

	res := r.col.FindOneAndUpdate(ctx, filter, update)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return ErrListingSoldOrNotFound
		}
		return res.Err()
	}
	return nil
}
