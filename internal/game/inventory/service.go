// File: internal/game/inventory/service.go
package inventory

import (
	"context"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/pkg/utils"
)

type Service interface {
	GetInventory(ctx context.Context, userID, guildID string) (*Inventory, []*item.ItemInstance, error)
	AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error
	GrantStarterItems(ctx context.Context, userID, guildID string) error
	UseItem(ctx context.Context, userID, guildID, instanceID string) (string, error)
}

type inventoryService struct {
	invRepo  Repository
	itemRepo item.Repository
	cultSvc  cultivation.Service
}

func NewService(invRepo Repository, itemRepo item.Repository, cultSvc cultivation.Service) Service {
	return &inventoryService{invRepo: invRepo, itemRepo: itemRepo, cultSvc: cultSvc}
}

func (s *inventoryService) GetInventory(ctx context.Context, userID, guildID string) (*Inventory, []*item.ItemInstance, error) {
	inv, err := s.invRepo.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		return nil, nil, err
	}
	items, err := s.itemRepo.GetInstancesByUser(ctx, userID, guildID)
	return inv, items, err
}

func (s *inventoryService) AddItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	def, ok := item.GetDefinition(definitionID)
	if !ok {
		return apperrors.ErrInvalidInput
	}

	inv, items, err := s.GetInventory(ctx, userID, guildID)
	if err != nil {
		return err
	}

	if def.Stackable {
		for _, it := range items {
			if it.DefinitionID == definitionID {
				return s.itemRepo.AdjustQuantity(ctx, it.InstanceID, userID, guildID, quantity)
			}
		}
	}

	// Cần tạo ô mới
	if len(items) >= inv.SlotLimit {
		return apperrors.ErrInventoryFull
	}

	inst := &item.ItemInstance{
		InstanceID:   utils.NewSessionID(),
		DefinitionID: definitionID,
		UserID:       userID,
		GuildID:      guildID,
		Quantity:     quantity,
	}
	return s.itemRepo.CreateInstance(ctx, inst)
}

func (s *inventoryService) GrantStarterItems(ctx context.Context, userID, guildID string) error {
	inv, err := s.invRepo.GetOrCreate(ctx, userID, guildID)
	if err != nil {
		return err
	}
	if inv.StarterGranted {
		return nil
	}

	_ = s.AddItem(ctx, userID, guildID, "pill_exp_small", 3)
	_ = s.AddItem(ctx, userID, guildID, "pill_stamina_small", 2)
	_ = s.AddItem(ctx, userID, guildID, "eq_wood_sword", 1)
	_ = s.AddItem(ctx, userID, guildID, "eq_cloth_robe", 1)
	_ = s.AddItem(ctx, userID, guildID, "refine_stone", 3)

	return s.invRepo.MarkStarterGranted(ctx, userID, guildID)
}

func (s *inventoryService) UseItem(ctx context.Context, userID, guildID, instanceID string) (string, error) {
	inst, err := s.itemRepo.GetInstanceByID(ctx, instanceID, userID, guildID)
	if err != nil {
		return "", err
	}

	def, ok := item.GetDefinition(inst.DefinitionID)
	if !ok || !def.Usable {
		return "", apperrors.ErrItemNotUsable
	}

	// 1. Trừ số lượng (Atomic)
	if err := s.itemRepo.AdjustQuantity(ctx, instanceID, userID, guildID, -1); err != nil {
		return "", err
	}

	// 2. Gọi hàm xóa (Hàm này đã được bảo vệ Atomic $lte: 0 ở Tầng Repository)
	_ = s.itemRepo.DeleteInstance(ctx, instanceID, userID, guildID)

	return "Sử dụng vật phẩm thành công!", nil
}