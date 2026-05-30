package equipment

import (
	"context"
	"errors"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/game/inventory"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

// --- Mocks ---
type mockEquipRepo struct {
	equipped map[EquipmentSlot]string
}

func (m *mockEquipRepo) Get(ctx context.Context, userID, guildID string) (*EquipmentSet, error) {
	slots := make(map[string]string)
	if m.equipped != nil {
		for k, v := range m.equipped {
			slots[string(k)] = v
		}
	}
	return &EquipmentSet{
		UserID:  userID,
		GuildID: guildID,
		Slots:   slots,
	}, nil
}
func (m *mockEquipRepo) Equip(ctx context.Context, userID, guildID string, slot EquipmentSlot, instanceID string) error {
	if m.equipped == nil {
		m.equipped = make(map[EquipmentSlot]string)
	}
	m.equipped[slot] = instanceID
	return nil
}
func (m *mockEquipRepo) Unequip(ctx context.Context, userID, guildID string, slot EquipmentSlot) error {
	if m.equipped != nil {
		delete(m.equipped, slot)
	}
	return nil
}

type mockItemRepo struct {
	validInstance string
}

func (m *mockItemRepo) CreateInstance(ctx context.Context, inst *item.ItemInstance) error {
	return nil
}

func (m *mockItemRepo) GetInstancesByUser(ctx context.Context, userID, guildID string) ([]*item.ItemInstance, error) {
	return nil, nil
}

func (m *mockItemRepo) GetInstanceByID(ctx context.Context, instanceID, userID, guildID string) (*item.ItemInstance, error) {
	if instanceID == "inst_wpn_max" {
		return &item.ItemInstance{InstanceID: instanceID, DefinitionID: "eq_weapon_moc_kiem_d", UserID: userID, Metadata: map[string]interface{}{"level": int32(10)}}, nil
	}
	if instanceID == m.validInstance || instanceID == "inst_wpn_1" {
		return &item.ItemInstance{InstanceID: instanceID, DefinitionID: "eq_weapon_moc_kiem_d", UserID: userID}, nil
	}
	switch instanceID {
	case "inst_armor_1":
		return &item.ItemInstance{InstanceID: instanceID, DefinitionID: "eq_armor_vai_tho_d", UserID: userID}, nil
	case "inst_pill_1":
		return &item.ItemInstance{InstanceID: instanceID, DefinitionID: "pill_exp_tu_khi_d", UserID: userID}, nil
	case "inst_missing_def":
		return &item.ItemInstance{InstanceID: instanceID, DefinitionID: "", UserID: userID}, nil
	}
	return nil, errors.New("ErrNotFound")
}

func (m *mockItemRepo) AdjustQuantity(ctx context.Context, instanceID, userID, guildID string, amount int64) error {
	return nil
}

func (m *mockItemRepo) DeleteInstance(ctx context.Context, instanceID, userID, guildID string) error {
	return nil
}
func (m *mockItemRepo) UpdateMetadata(ctx context.Context, instanceID, userID, guildID string, metadata map[string]interface{}) error {
	return nil
}

type mockInvSvc struct {
	failConsume   bool
	consumeCalled bool
}

func (m *mockInvSvc) ConsumeItems(ctx context.Context, userID, guildID string, itemsToConsume map[string]int64) error {
	m.consumeCalled = true
	if m.failConsume {
		return errors.New("mock error: thiếu nguyên liệu")
	}
	return nil
}
func (m *mockInvSvc) AddItem(ctx context.Context, userID, guildID, defID string, qty int64) error {
	return nil
}
func (m *mockInvSvc) GetInventory(ctx context.Context, userID, guildID string) (*inventory.Inventory, []*item.ItemInstance, error) {
	return nil, nil, nil
}
func (m *mockInvSvc) GrantStarterItems(ctx context.Context, userID, guildID string) error { return nil }
func (m *mockInvSvc) UseItem(ctx context.Context, userID, guildID, instanceID string) (string, error) {
	return "", nil
}
func (m *mockInvSvc) DismantleItem(ctx context.Context, userID, guildID, instanceID string) (string, error) {
	return "", nil
}
func (m *mockInvSvc) RemoveItem(ctx context.Context, userID, guildID, definitionID string, quantity int64) error {
	return nil
}

// --- Tests ---

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"eq_weapon_moc_kiem_d": {ID: "eq_weapon_moc_kiem_d", Name: "Mộc Kiếm", Type: item.TypeEquipment, MaxEnhanceLevel: 10},
		"eq_armor_vai_tho_d":   {ID: "eq_armor_vai_tho_d", Name: "Vải Thô", Type: item.TypeEquipment},
		"pill_exp_tu_khi_d":    {ID: "pill_exp_tu_khi_d", Name: "Tụ Khí Đan", Type: item.TypePill},
	})
}

func TestEquipment_EquipSuccess(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_1"}
	mockEq := &mockEquipRepo{}
	svc := NewService(mockEq, mockItem, &mockInvSvc{})

	err := svc.Equip(context.Background(), "user1", "guild1", "weapon", "inst_wpn_1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if mockEq.equipped["weapon"] != "inst_wpn_1" {
		t.Errorf("Trang bị không được lưu đúng slot")
	}
}

func TestEquipment_EquipNotOwned(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_1"}
	mockEq := &mockEquipRepo{}
	svc := NewService(mockEq, mockItem, &mockInvSvc{})

	// Thử mặc item của người khác (instance ID không tồn tại cho user này)
	err := svc.Equip(context.Background(), "user1", "guild1", "weapon", "inst_hacker_1")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi mặc trang bị không sở hữu")
	}
}

func TestEquipment_EquipRejectWrongSlot(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_1"}
	mockEq := &mockEquipRepo{}
	svc := NewService(mockEq, mockItem, &mockInvSvc{})

	err := svc.Equip(context.Background(), "user1", "guild1", "armor", "inst_wpn_1")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi mặc sai vị trí (Vũ khí vào ô Giáp)")
	}
}

func TestEquipment_EquipRejectMissingDefinitionID(t *testing.T) {
	mockItem := &mockItemRepo{}
	mockEq := &mockEquipRepo{}
	svc := NewService(mockEq, mockItem, &mockInvSvc{})

	err := svc.Equip(context.Background(), "user1", "guild1", "weapon", "inst_missing_def")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi thiếu DefinitionID")
	}
}

func TestEquipment_EquipRejectNonEquipment(t *testing.T) {
	mockItem := &mockItemRepo{}
	mockEq := &mockEquipRepo{}
	svc := NewService(mockEq, mockItem, &mockInvSvc{})

	err := svc.Equip(context.Background(), "user1", "guild1", "weapon", "inst_pill_1")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi equip vật phẩm không phải trang bị (Đan dược)")
	}
}

func TestEquipment_Enhance_Success(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_1"}
	mockEq := &mockEquipRepo{equipped: map[EquipmentSlot]string{"weapon": "inst_wpn_1"}}
	mockInv := &mockInvSvc{}
	svc := NewService(mockEq, mockItem, mockInv)

	err := svc.Enhance(context.Background(), "user1", "guild1", "weapon")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi khi cường hóa hợp lệ: %v", err)
	}
}

func TestEquipment_Enhance_MaxLevel(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_max"}
	mockEq := &mockEquipRepo{equipped: map[EquipmentSlot]string{"weapon": "inst_wpn_max"}}
	mockInv := &mockInvSvc{}
	svc := NewService(mockEq, mockItem, mockInv)

	err := svc.Enhance(context.Background(), "user1", "guild1", "weapon")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi cường hóa trang bị max level")
	}
	if mockInv.consumeCalled {
		t.Error("Không được phép tiêu thụ nguyên liệu khi đạt cấp cường hóa tối đa")
	}
}

func TestEquipment_Enhance_MissingMaterials(t *testing.T) {
	mockItem := &mockItemRepo{validInstance: "inst_wpn_1"}
	mockEq := &mockEquipRepo{equipped: map[EquipmentSlot]string{"weapon": "inst_wpn_1"}}
	mockInv := &mockInvSvc{failConsume: true}
	svc := NewService(mockEq, mockItem, mockInv)

	err := svc.Enhance(context.Background(), "user1", "guild1", "weapon")
	if err == nil {
		t.Fatal("Mong đợi lỗi khi thiếu nguyên liệu")
	}
}

func TestEquipment_Enhance_EmptySlot(t *testing.T) {
	svc := NewService(&mockEquipRepo{}, &mockItemRepo{}, &mockInvSvc{})
	if err := svc.Enhance(context.Background(), "u1", "g1", "weapon"); err == nil {
		t.Fatal("Mong đợi lỗi khi cường hóa slot trống")
	}
}

func TestEquipment_Enhance_NotOwned(t *testing.T) {
	mockEq := &mockEquipRepo{equipped: map[EquipmentSlot]string{"weapon": "inst_hacker_1"}}
	svc := NewService(mockEq, &mockItemRepo{}, &mockInvSvc{})
	if err := svc.Enhance(context.Background(), "u1", "g1", "weapon"); err == nil {
		t.Fatal("Mong đợi lỗi khi cường hóa trang bị không thuộc sở hữu/không tồn tại")
	}
}
