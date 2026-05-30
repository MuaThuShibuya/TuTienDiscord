package player

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListingStatus string

const (
	StatusActive    ListingStatus = "active"
	StatusSold      ListingStatus = "sold"
	StatusCancelled ListingStatus = "cancelled"
	StatusExpired   ListingStatus = "expired"
)

type Listing struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	SellerID       string             `bson:"sellerId"`
	BuyerID        string             `bson:"buyerId,omitempty"` // ID người mua thành công
	GuildID        string             `bson:"guildId"`
	ItemInstanceID string             `bson:"itemInstanceId"`
	ItemDefID      string             `bson:"itemDefId"`
	Quantity       int64              `bson:"quantity"`
	TotalPrice     int64              `bson:"totalPrice"`
	Status         ListingStatus      `bson:"status"`
	CreatedAt      time.Time          `bson:"createdAt"`
	PurchasedAt    time.Time          `bson:"purchasedAt,omitempty"`
	ExpiresAt      time.Time          `bson:"expiresAt"`
}
