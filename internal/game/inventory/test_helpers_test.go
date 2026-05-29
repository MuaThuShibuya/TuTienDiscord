package inventory_test

import (
	"context"
	"sync"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

func init() {
	_ = logger.Init(logger.Options{Level: "error", Format: "json"})
}

type mockInvRepo struct{}

func (m *mockInvRepo) GetOrCreate(ctx context.Context, userID, guildID string) (*inventory.Inventory, error) {
	return &inventory.Inventory{SlotLimit: 50, UserID: userID, GuildID: guildID}, nil
}
func (m *mockInvRepo) MarkStarterGranted(ctx context.Context, userID, guildID string) error {
	return nil
}
func (m *mockInvRepo) AcquireSlot(ctx context.Context, userID, guildID string) error { return nil }
func (m *mockInvRepo) ReleaseSlot(ctx context.Context, userID, guildID string) error { return nil }

type mockItemRepo struct {
	sync.Mutex
	items map[string]*item.ItemInstance
}

func newMockItemRepo() *mockItemRepo {
	return &mockItemRepo{items: make(map[string]*item.ItemInstance)}
}
func (m *mockItemRepo) CreateInstance(ctx context.Context, inst *item.ItemInstance) error {
	m.Lock()
	defer m.Unlock()
	m.items[inst.InstanceID] = inst
	return nil
}
func (m *mockItemRepo) GetInstancesByUser(ctx context.Context, userID, guildID string) ([]*item.ItemInstance, error) {
	m.Lock()
	defer m.Unlock()
	var res []*item.ItemInstance
	for _, it := range m.items {
		if it.UserID == userID && it.GuildID == guildID {
			res = append(res, it)
		}
	}
	return res, nil
}
func (m *mockItemRepo) GetInstanceByID(ctx context.Context, instanceID, userID, guildID string) (*item.ItemInstance, error) {
	m.Lock()
	defer m.Unlock()
	it, ok := m.items[instanceID]
	if !ok {
		return nil, apperrors.ErrItemNotFound
	}
	return it, nil
}
func (m *mockItemRepo) AdjustQuantity(ctx context.Context, instanceID, userID, guildID string, amount int64) error {
	m.Lock()
	defer m.Unlock()
	it, ok := m.items[instanceID]
	if !ok {
		return apperrors.ErrItemNotFound
	}
	if it.Quantity+amount < 0 {
		return apperrors.ErrInsufficientItemQuantity
	}
	it.Quantity += amount
	return nil
}
func (m *mockItemRepo) DeleteInstance(ctx context.Context, instanceID, userID, guildID string) error {
	m.Lock()
	defer m.Unlock()
	it, ok := m.items[instanceID]
	if ok && it.Quantity <= 0 {
		delete(m.items, instanceID)
	}
	return nil
}
func (m *mockItemRepo) UpdateMetadata(ctx context.Context, instanceID, userID, guildID string, metadata map[string]interface{}) error {
	return nil
}
