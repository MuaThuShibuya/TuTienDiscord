package invmenu

import (
	"strings"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
)

func TestInventoryBuildMenuResponse_WithItemsShowsItemNames(t *testing.T) {
	vm := &menu.InventoryMenuVM{
		Items: []menu.InventoryItemVM{
			{Name: "Kiếm Gỗ Test", Quantity: 1},
		},
		TotalPages:  1,
		CurrentPage: 1,
	}
	res := BuildMenuResponse(vm)
	desc := res.Embeds[0].Description
	if !strings.Contains(desc, "Kiếm Gỗ Test") {
		t.Errorf("Expected embed to contain item name, got: %s", desc)
	}
}

func TestInventoryBuildMenuResponse_EmptyShowsEmptyState(t *testing.T) {
	vm := &menu.InventoryMenuVM{
		Items:       []menu.InventoryItemVM{},
		TotalPages:  1,
		CurrentPage: 1,
	}
	res := BuildMenuResponse(vm)
	desc := res.Embeds[0].Description
	if !strings.Contains(desc, "chưa có vật phẩm nào") {
		t.Errorf("Expected empty state message, got: %s", desc)
	}
}
