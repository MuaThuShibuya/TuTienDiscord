package inventory_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- Mocks ---
type mockInvRepo struct{}

func (m *mockInvRepo) GetOrCreate(ctx context.Context, userID, guildID string) (*inventory.Inventory, error) {
	return &inventory.Inventory{SlotLimit: 50}, nil
}
func (m *mockInvRepo) MarkStarterGranted(ctx context.Context, userID, guildID string) error {
	return nil
}
func (m *mockInvRepo) AcquireSlot(ctx context.Context, userID, guildID string) error {
	return nil // Chấp nhận luôn trong bài test hiện tại
}
func (m *mockInvRepo) ReleaseSlot(ctx context.Context, userID, guildID string) error {
	return nil
}

type mockCultSvc struct {
	cultivation.Service
	failNext bool
	exp      int64
	stamina  int
}

func (m *mockCultSvc) AddExperience(ctx context.Context, userID, guildID string, amount int64) error {
	if m.failNext {
		return errors.New("mock error")
	}
	m.exp += amount
	return nil
}
func (m *mockCultSvc) AddStamina(ctx context.Context, userID, guildID string, amount int) error {
	if m.failNext {
		return errors.New("mock error")
	}
	m.stamina += amount
	return nil
}

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
	return nil, nil
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

type mockCultRepo struct {
	prof *cultivation.CultivationProfile
}

func (m *mockCultRepo) FindByUserID(ctx context.Context, userID, guildID string) (*cultivation.CultivationProfile, error) {
	return m.prof, nil
}
func (m *mockCultRepo) Upsert(ctx context.Context, profile *cultivation.CultivationProfile) error {
	return nil
}
func (m *mockCultRepo) UpdateStats(ctx context.Context, profile *cultivation.CultivationProfile) error {
	m.prof = profile
	return nil
}

func init() {
	_ = logger.Init(logger.Options{Level: "error", Format: "json"})
	item.RegisterItems(map[string]item.ItemDefinition{
		"test_pill_exp": {ID: "test_pill_exp", Name: "EXP Pill", Type: item.TypePill, Usable: true, Effects: map[string]int{"exp": 100}},
		"test_pill_stm": {ID: "test_pill_stm", Name: "STM Pill", Type: item.TypePill, Usable: true, Effects: map[string]int{"stamina": 20}},
		"test_pill_brk": {ID: "test_pill_brk", Name: "BRK Pill", Type: item.TypePill, Usable: true, Effects: map[string]int{"breakthrough_chance": 5}},
		"test_mat":      {ID: "test_mat", Name: "Mat", Type: item.TypeMaterial, Usable: false},
	})
}

// --- Tests ---
func TestInventory_UseExpPill_AddsCultivationExp(t *testing.T) {
	cultSvc := &mockCultSvc{}
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_exp", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)

	msg, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if cultSvc.exp != 100 {
		t.Errorf("Exp mong đợi 100, nhận %d", cultSvc.exp)
	}
	if !strings.Contains(msg, "100") {
		t.Errorf("Message phải chứa số exp: %s", msg)
	}
}

func TestInventory_UseStaminaPill_AddsStamina(t *testing.T) {
	cultSvc := &mockCultSvc{}
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_stm", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)

	_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if cultSvc.stamina != 20 {
		t.Errorf("Stamina mong đợi 20, nhận %d", cultSvc.stamina)
	}
}

func TestInventory_UseStaminaPill_DoesNotExceedMax(t *testing.T) {
	cultRepo := &mockCultRepo{prof: &cultivation.CultivationProfile{UserID: "u1", GuildID: "g1", Stamina: 90, MaxStamina: 100}}
	realCultSvc := cultivation.NewService(cultRepo, nil, nil) // Inject mem repo

	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_stm", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, realCultSvc)

	_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if cultRepo.prof.Stamina != 100 {
		t.Errorf("Stamina mong đợi 100 (bị clamp max), nhận %d", cultRepo.prof.Stamina)
	}
}

func TestInventory_UsePill_ConsumesOneQuantity(t *testing.T) {
	cultSvc := &mockCultSvc{}
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_exp", Quantity: 5}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)

	_, _ = svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if itemRepo.items["inst1"].Quantity != 4 {
		t.Errorf("Quantity mong đợi 4, nhận %d", itemRepo.items["inst1"].Quantity)
	}
}

func TestInventory_UsePill_RollbackWhenEffectFails(t *testing.T) {
	cultSvc := &mockCultSvc{failNext: true}
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_exp", Quantity: 5}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)

	_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi effect fail")
	}
	if itemRepo.items["inst1"].Quantity != 5 {
		t.Errorf("Quantity mong đợi 5 (đã rollback), nhận %d", itemRepo.items["inst1"].Quantity)
	}
}

func TestInventory_UseItem_RejectUnknownDefinition(t *testing.T) {
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "unknown_pill", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, nil)

	_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err == nil || !strings.Contains(err.Error(), "tồn tại") {
		t.Errorf("Mong đợi lỗi unknown definition, nhận: %v", err)
	}
}

func TestInventory_UseItem_RejectNonUsableItem(t *testing.T) {
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_mat", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, nil)

	_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
	if err != apperrors.ErrItemNotUsable {
		t.Errorf("Mong đợi lỗi ErrItemNotUsable, nhận: %v", err)
	}
}

func Test_UseItem_Concurrent_DoubleSpendCheck(t *testing.T) {
	cultSvc := &mockCultSvc{}
	itemRepo := newMockItemRepo()
	itemRepo.items["inst1"] = &item.ItemInstance{InstanceID: "inst1", DefinitionID: "test_pill_exp", Quantity: 1}
	svc := inventory.NewService(&mockInvRepo{}, itemRepo, cultSvc)

	var wg sync.WaitGroup
	var successCount int32
	const concurrentRequests = 10

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := svc.UseItem(context.Background(), "u1", "g1", "inst1")
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	if successCount != 1 {
		t.Errorf("Chỉ mong đợi 1 lần UseItem thành công, nhận %d", successCount)
	}

	itemRepo.Lock()
	it, ok := itemRepo.items["inst1"]
	itemRepo.Unlock()

	if ok && it.Quantity < 0 {
		t.Errorf("Quantity không được âm, nhận %d", it.Quantity)
	}
	if cultSvc.exp != 100 {
		t.Errorf("Exp chỉ được cộng 1 lần (100), nhận %d", cultSvc.exp)
	}

	t.Log("Lưu ý: Test này chứng minh logic an toàn qua in-memory sync.Mutex. " +
		"Khi chạy trên MongoDB, tính toàn vẹn phụ thuộc vào atomic update $inc với điều kiện $gte: 0.")
}
