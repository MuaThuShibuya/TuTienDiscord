package npc

import (
	"context"
	"errors"
	"fmt"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/economy"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"go.uber.org/zap"
)

type Service interface {
	BuyItem(ctx context.Context, userID, guildID, shopID, itemID string, quantity int) (*economy.Wallet, error)
	SellItem(ctx context.Context, userID, guildID, shopID, itemID string, quantity int) (*economy.Wallet, error)
}

type service struct {
	ecoSvc economy.Service
	invSvc inventory.Service
	log    *zap.Logger
}

func NewService(ecoSvc economy.Service, invSvc inventory.Service, log *zap.Logger) Service {
	return &service{ecoSvc: ecoSvc, invSvc: invSvc, log: log.Named("game.shop.npc")}
}

func (s *service) BuyItem(ctx context.Context, userID, guildID, shopID, itemID string, quantity int) (*economy.Wallet, error) {
	if quantity <= 0 {
		return nil, apperrors.New("INVALID_QTY", "Số lượng mua không hợp lệ.", errors.New("qty <= 0"))
	}

	shopDef, err := GetShop(shopID)
	if err != nil {
		return nil, apperrors.New("SHOP_NOT_FOUND", "Thương hội này không tồn tại.", err)
	}
	itemDef, err := shopDef.GetItem(itemID)
	if err != nil {
		return nil, apperrors.New("ITEM_NOT_FOR_SALE", "Vật phẩm này không có bán tại đây.", err)
	}

	totalCost := itemDef.BuyPrice * int64(quantity)

	wallet, err := s.ecoSvc.SpendSpiritStones(ctx, userID, guildID, totalCost, fmt.Sprintf("npc_buy_%s_%s", shopID, itemID))
	if err != nil {
		if apperrors.IsInsufficientFunds(err) {
			return nil, apperrors.New("NO_MONEY", "Đạo hữu không đủ linh thạch.", err)
		}
		return nil, err
	}

	if err := s.invSvc.AddItem(ctx, userID, guildID, itemID, int64(quantity)); err != nil {
		s.log.Warn("Túi đầy, rollback giao dịch mua", zap.String("userId", userID))
		_, _ = s.ecoSvc.EarnSpiritStones(ctx, userID, guildID, totalCost, "rollback_npc_buy")
		return nil, apperrors.New("INV_FULL", "Túi đồ đã đầy, linh thạch đã được hoàn trả.", err)
	}

	return wallet, nil
}

func (s *service) SellItem(ctx context.Context, userID, guildID, shopID, itemID string, quantity int) (*economy.Wallet, error) {
	if quantity <= 0 {
		return nil, apperrors.New("INVALID_QTY", "Số lượng bán không hợp lệ.", errors.New("qty <= 0"))
	}

	shopDef, err := GetShop(shopID)
	if err != nil {
		return nil, apperrors.New("SHOP_NOT_FOUND", "Thương hội này không tồn tại.", err)
	}
	itemDef, err := shopDef.GetItem(itemID)
	if err != nil {
		return nil, apperrors.New("ITEM_NOT_BOUGHT", "Thương hội không thu mua vật phẩm này.", err)
	}
	if itemDef.SellPrice <= 0 {
		return nil, apperrors.New("ITEM_NO_VALUE", "Vật phẩm này không có giá trị thu mua.", errors.New("sell price 0"))
	}

	// Yêu cầu hàm RemoveItem từ Inventory Service
	if err := s.invSvc.RemoveItem(ctx, userID, guildID, itemID, int64(quantity)); err != nil {
		return nil, apperrors.New("NO_ITEM", "Đạo hữu không có vật phẩm này hoặc số lượng không đủ.", err)
	}

	totalEarn := itemDef.SellPrice * int64(quantity)
	wallet, err := s.ecoSvc.EarnSpiritStones(ctx, userID, guildID, totalEarn, fmt.Sprintf("npc_sell_%s_%s", shopID, itemID))
	return wallet, err
}
