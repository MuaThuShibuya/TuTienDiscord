package player

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"go.uber.org/zap"
)

type Service interface {
	PurchaseListing(ctx context.Context, buyerID, guildID, listingID string) error
	GetActiveListings(ctx context.Context, guildID string, limit, offset int) ([]*Listing, error)
	CreateListing(ctx context.Context, userID, guildID, itemDefID string, quantity, price int64) error
	GetUserActiveListings(ctx context.Context, userID, guildID string) ([]*Listing, error)
	CancelListing(ctx context.Context, userID, guildID, listingID string) error
}

type service struct {
	repo   Repository
	ecoSvc economy.Service
	invSvc inventory.Service
	log    *zap.Logger
}

func NewService(repo Repository, ecoSvc economy.Service, invSvc inventory.Service, log *zap.Logger) Service {
	return &service{repo: repo, ecoSvc: ecoSvc, invSvc: invSvc, log: log.Named("game.shop.player")}
}

func (s *service) GetActiveListings(ctx context.Context, guildID string, limit, offset int) ([]*Listing, error) {
	return s.repo.GetActiveListings(ctx, guildID, limit, offset)
}

func (s *service) GetUserActiveListings(ctx context.Context, userID, guildID string) ([]*Listing, error) {
	return s.repo.GetUserActiveListings(ctx, userID, guildID)
}

func (s *service) CancelListing(ctx context.Context, userID, guildID, listingID string) error {
	listing, err := s.repo.CancelListing(ctx, listingID, userID)
	if err != nil {
		return apperrors.New("CANCEL_FAIL", "Không thể thu hồi. Có thể đã bị người khác mua mất hoặc không tồn tại.", err)
	}
	return s.invSvc.AddItem(ctx, userID, guildID, listing.ItemDefID, listing.Quantity)
}

func (s *service) CreateListing(ctx context.Context, userID, guildID, itemDefID string, quantity, price int64) error {
	if err := s.invSvc.RemoveItem(ctx, userID, guildID, itemDefID, quantity); err != nil {
		return apperrors.New("NO_ITEM", "Không đủ vật phẩm trong túi để đăng bán.", err)
	}

	listing := &Listing{
		SellerID:   userID,
		GuildID:    guildID,
		ItemDefID:  itemDefID,
		Quantity:   quantity,
		TotalPrice: price,
		Status:     StatusActive,
		CreatedAt:  time.Now().UTC(),
		ExpiresAt:  time.Now().UTC().Add(24 * time.Hour),
	}
	if err := s.repo.CreateListing(ctx, listing); err != nil {
		_ = s.invSvc.AddItem(ctx, userID, guildID, itemDefID, quantity)
		return apperrors.New("DB_ERROR", "Lỗi tạo phiếu đấu giá.", err)
	}
	return nil
}

func (s *service) PurchaseListing(ctx context.Context, buyerID, guildID, listingID string) error {
	listing, err := s.repo.GetListing(ctx, listingID)
	if err != nil {
		return apperrors.New("LISTING_404", "Phiếu đấu giá không tồn tại hoặc đã bị gỡ.", err)
	}

	if listing.Status != StatusActive {
		return apperrors.New("LISTING_NOT_ACTIVE", "Chậm chân một bước! Vật phẩm này đã có chủ hoặc hết hạn.", errors.New("not active"))
	}

	if listing.SellerID == buyerID {
		return apperrors.New("SELF_BUY", "Không thể tự mua vật phẩm của chính mình!", errors.New("self buy"))
	}

	// Trừ tiền người mua (Kiểm tra ví)
	reasonBuy := fmt.Sprintf("mua_dau_gia_%s", listingID)
	_, err = s.ecoSvc.SpendSpiritStones(ctx, buyerID, guildID, listing.TotalPrice, reasonBuy)
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return apperrors.New("NO_MONEY", "Đạo hữu không đủ linh thạch để chốt đơn này.", err)
		}
		return err
	}

	// Khóa Atomic cấp DB (Compare-And-Swap)
	err = s.repo.AtomicPurchase(ctx, listingID, buyerID)
	if err != nil {
		// ROLLBACK
		s.log.Warn("Race condition bị chặn: Rollback tiền mua đấu giá", zap.String("buyerId", buyerID), zap.String("listingId", listingID))
		_, _ = s.ecoSvc.EarnSpiritStones(ctx, buyerID, guildID, listing.TotalPrice, "hoan_tien_truot_dau_gia")
		return apperrors.New("LISTING_SOLD_OUT", "Vật phẩm đã bị người khác nhanh tay mua mất. Linh thạch đã được hoàn trả.", err)
	}

	// Chuyển tiền cho Seller
	reasonSell := fmt.Sprintf("ban_dau_gia_%s", listingID)
	_, _ = s.ecoSvc.EarnSpiritStones(ctx, listing.SellerID, guildID, listing.TotalPrice, reasonSell)

	// Chuyển đồ cho Buyer
	_ = s.invSvc.AddItem(ctx, buyerID, guildID, listing.ItemDefID, listing.Quantity)

	return nil
}
