// File: internal/discord/menu/equipment/components.go
// Chức năng: Tạo Discord components (nút tháo, select menu mặc trang bị) cho Trang Bị.

package equipmenu

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
	"github.com/whiskey/tu-tien-bot/internal/discord/ui/emoji"
)

func buildComponents(vm *menu.EquipmentMenuVM) []discordgo.MessageComponent {
	var components []discordgo.MessageComponent

	// Hàng 1: Nút tháo trang bị (chỉ hiện nếu có slot đang mặc)
	unequipSlots := []struct {
		slot string
		name string
		item *menu.EquippedItemVM
	}{
		{"weapon", "Vũ Khí", vm.Weapon},
		{"armor", "Giáp", vm.Armor},
		{"artifact", "Pháp Bảo", vm.Artifact},
		{"treasure", "Bảo Vật", vm.Treasure},
	}

	var unequipBtns []discordgo.MessageComponent
	for _, s := range unequipSlots {
		if s.item != nil {
			unequipBtns = append(unequipBtns, ui.Button(
				"Tháo "+s.name,
				menu.Build(menu.DomainEquipment, menu.ActionEquipmentUnequip, vm.SessionID, s.slot),
				ui.BtnDanger, emoji.Equip, false,
			))
		}
	}
	if len(unequipBtns) > 0 {
		components = append(components, ui.ActionRow(unequipBtns...))
	}

	// Hàng 2: Select menu mặc trang bị từ túi đồ (chỉ hiện nếu có trang bị khả dụng)
	if len(vm.Equippable) > 0 {
		var equipOptions []discordgo.SelectMenuOption
		for _, it := range vm.Equippable {
			if len(equipOptions) >= 25 {
				break
			}
			// Giá trị: "instanceID:definitionID" — router parse để xác định slot
			value := it.InstanceID + ":" + it.DefinitionID
			equipOptions = append(equipOptions, ui.SelectOption(
				fmt.Sprintf("[%s] %s", it.Rarity, it.Name),
				value,
				"Vị trí: "+it.SlotName,
				emoji.Equip,
				false,
			))
		}
		equipSelect := ui.SelectMenu(
			menu.Build(menu.DomainEquipment, menu.ActionEquipmentEquip, vm.SessionID),
			"⚔️ Mặc trang bị từ túi đồ...",
			equipOptions,
		)
		components = append(components, ui.ActionRow(equipSelect))
	}

	// Hàng cuối: Điều hướng
	components = append(components, ui.NavRow(vm.SessionID, string(menu.PageEquipment), string(menu.PageMain)))
	return components
}
