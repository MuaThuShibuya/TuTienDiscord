package player

import "context"

type Repository interface {
	CreateListing(ctx context.Context, listing *Listing) error
	GetListing(ctx context.Context, listingID string) (*Listing, error)
	GetActiveListings(ctx context.Context, guildID string, limit, offset int) ([]*Listing, error)
	GetUserActiveListings(ctx context.Context, userID, guildID string) ([]*Listing, error)

	// Lệnh DB chỉ cập nhật status -> "sold" NẾU status hiện tại là "active".
	AtomicPurchase(ctx context.Context, listingID string, buyerID string) error
	CancelListing(ctx context.Context, listingID, sellerID string) (*Listing, error)
}
