package invmenu

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
)

func TestInventoryComponents_CustomIDsAreUniqueOnSinglePage(t *testing.T) {
	vm := &menu.InventoryMenuVM{
		SessionID:   "sess_test",
		CurrentPage: 1,
		TotalPages:  1,
		Items: []menu.InventoryItemVM{
			{Name: "Kiếm Gỗ", Quantity: 1},
		},
	}

	components := buildComponents(vm)
	assertNoDuplicateCustomIDs(t, components)
}

func TestInventoryComponents_CustomIDsAreUniqueOnMultiPage(t *testing.T) {
	vm := &menu.InventoryMenuVM{
		SessionID:   "sess_test",
		CurrentPage: 2,
		TotalPages:  3,
		Items: []menu.InventoryItemVM{
			{Name: "Đan Dược", Quantity: 1},
		},
	}

	components := buildComponents(vm)
	assertNoDuplicateCustomIDs(t, components)
}

func assertNoDuplicateCustomIDs(t *testing.T, rows []discordgo.MessageComponent) {
	t.Helper()
	seen := map[string]bool{}

	for _, row := range rows {
		actionRow, ok := row.(discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, c := range actionRow.Components {
			switch comp := c.(type) {
			case discordgo.Button:
				if comp.CustomID != "" {
					if seen[comp.CustomID] {
						t.Fatalf("duplicate custom_id found: %s", comp.CustomID)
					}
					seen[comp.CustomID] = true
				}
			case discordgo.SelectMenu:
				if comp.CustomID != "" {
					if seen[comp.CustomID] {
						t.Fatalf("duplicate custom_id found: %s", comp.CustomID)
					}
					seen[comp.CustomID] = true
				}
			}
		}
	}
}
