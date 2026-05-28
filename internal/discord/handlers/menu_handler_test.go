package handlers

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
	"github.com/whiskey/tu-tien-bot/internal/game/profile"
)

func init() {
	item.RegisterItems(map[string]item.ItemDefinition{
		"usable_pill": {ID: "usable_pill", Name: "Pill", Usable: true, Type: item.TypePill},
		"equip_sword": {ID: "equip_sword", Name: "Sword", Usable: false, Type: item.TypeEquipment},
	})
}

func TestToInventoryMenuVM_MapsItemsIntoVM(t *testing.T) {
	session := &menu.Session{CurrentCategory: "1"}
	player := &profile.Player{DaoName: "DaoHuu"}

	items := make([]*item.ItemInstance, 6)
	for i := 0; i < 6; i++ {
		items[i] = &item.ItemInstance{InstanceID: "id", DefinitionID: "usable_pill", Quantity: 1}
	}

	vm := toInventoryMenuVM(session, player, items)
	if len(vm.Items) != 6 {
		t.Errorf("Expected 6 items, got %d", len(vm.Items))
	}
}

func TestToInventoryMenuVM_DoesNotUseUsableItemsAsMainList(t *testing.T) {
	session := &menu.Session{CurrentCategory: "1"}
	player := &profile.Player{DaoName: "DaoHuu"}

	items := make([]*item.ItemInstance, 6)
	for i := 0; i < 6; i++ {
		items[i] = &item.ItemInstance{InstanceID: "id", DefinitionID: "equip_sword", Quantity: 1}
	}

	vm := toInventoryMenuVM(session, player, items)
	if len(vm.Items) != 6 {
		t.Errorf("Expected 6 items in main list, got %d", len(vm.Items))
	}
	if len(vm.UsableItems) != 0 {
		t.Errorf("Expected 0 usable items, got %d", len(vm.UsableItems))
	}
}
