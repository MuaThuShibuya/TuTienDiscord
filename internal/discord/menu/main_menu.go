// File: internal/discord/menu/main_menu.go
// Phiên bản: v0.1.1
// Mục đích: UI Builder cho trang Main Menu — render embed tổng quan và select menu điều hướng.
// Bảo mật: Chỉ nhận ViewModel đã xử lý sẵn, không gọi DB hay service, không có side effect.
// Ghi chú: Handler map domain models → MainMenuVM trước khi gọi hàm này.

package menu

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/whiskey/tu-tien-bot/internal/discord/ui"
)

// BuildMainMenuResponse tạo response Discord đầy đủ cho trang Main Menu.
func BuildMainMenuResponse(vm *MainMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildMainEmbed(vm)},
		Components: buildMainComponents(vm.SessionID),
	}
}

// BuildMainMenuEdit tạo response chỉnh sửa message khi điều hướng về Main Menu.
func BuildMainMenuEdit(vm *MainMenuVM) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildMainEmbed(vm)},
		Components: buildMainComponents(vm.SessionID),
	}
}

func buildMainEmbed(vm *MainMenuVM) *discordgo.MessageEmbed {
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   ui.EmojiRealm.String() + " Cảnh Giới",
			Value:  vm.RealmDisplay,
			Inline: true,
		},
		{
			Name:   ui.EmojiCombatPower.String() + " Chiến Lực",
			Value:  vm.CombatPower,
			Inline: true,
		},
		{
			Name:   ui.EmojiMindState.String() + " Tâm Cảnh",
			Value:  vm.MindState,
			Inline: true,
		},
		{
			Name:   ui.EmojiSkill.String() + " Đạo Lộ",
			Value:  vm.PathDisplay,
			Inline: true,
		},
		{
			Name:   ui.EmojiStamina.String() + " Thể Lực",
			Value:  vm.StaminaBar,
			Inline: false,
		},
		{
			Name:   ui.EmojiCultivate.String() + " Tu Vi",
			Value:  vm.ExpBar,
			Inline: false,
		},
		{
			Name:   ui.EmojiSpiritStone.String() + " Linh Thạch",
			Value:  vm.SpiritStones,
			Inline: true,
		},
		{
			Name:   ui.EmojiSpiritJade.String() + " Linh Ngọc",
			Value:  vm.SpiritJades,
			Inline: true,
		},
		{
			Name:   ui.EmojiFateTicket.String() + " Vé Cơ Duyên",
			Value:  vm.FateTickets,
			Inline: true,
		},
		{
			Name:  ui.EmojiInfo.String() + " Gợi ý hôm nay",
			Value: "_" + vm.DailyTip + "_",
		},
	}

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiProfile.String() + " Vạn Pháp Tiên Nghịch — " + vm.DaoName,
		Description: "Chào mừng trở lại, **" + vm.DaoName + "**!\nHãy chọn chức năng bên dưới.",
		Color:       ui.ColorDefault,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Vạn Pháp Tiên Nghịch · v0.1",
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildMainComponents(sessionID string) []discordgo.MessageComponent {
	// Hàng 1: Select menu danh mục chính
	categorySelect := ui.SelectMenu(
		Build(DomainMenuSelect, ActionNavSelect, sessionID),
		"✦ Chọn chức năng...",
		[]discordgo.SelectMenuOption{
			ui.SelectOption("Hồ Sơ", string(PageProfile),
				"Xem và chỉnh sửa thông tin đạo hữu", &ui.EmojiProfile, false),
			ui.SelectOption("Tu Luyện", string(PageCultivation),
				"Tĩnh tu, bế quan, đột phá cảnh giới", &ui.EmojiCultivate, false),
			ui.SelectOption("Chiến Đấu", "combat",
				"PvE, PvP, Boss server", &ui.EmojiSword, false),
			ui.SelectOption("Túi Đồ / Trang Bị", "inventory",
				"Quản lý vật phẩm và trang bị", &ui.EmojiBag, false),
			ui.SelectOption("Kỹ Năng / Công Pháp", "skills",
				"Tu luyện kỹ năng và công pháp", &ui.EmojiSkill, false),
		},
	)

	// Hàng 2: Select menu danh mục phụ
	categorySelect2 := ui.SelectMenu(
		Build(DomainMenuSelect, ActionNav2Select, sessionID),
		"✦ Chọn thêm...",
		[]discordgo.SelectMenuOption{
			ui.SelectOption("Linh Thú / Con Rối", "pets",
				"Nuôi và ra trận cùng linh thú", &ui.EmojiPet, false),
			ui.SelectOption("Cơ Duyên / Gacha", "gacha",
				"Quay cơ duyên bằng vé (không dùng tiền thật)", &ui.EmojiGacha, false),
			ui.SelectOption("Chợ / Đấu Giá", "market",
				"Mua bán và đấu giá vật phẩm", &ui.EmojiMarket, false),
			ui.SelectOption("Tông Môn / NPC / Đạo Lữ", "sect",
				"Tông môn, NPC, và hợp tác đồng hành", &ui.EmojiSect, false),
		},
	)

	// Hàng 3: Nút điều hướng (làm mới / quay lại / đóng) — Main Menu không có trang cha
	navRow := ui.NavRow(sessionID, string(PageMain), "")

	return []discordgo.MessageComponent{
		ui.ActionRow(categorySelect),
		ui.ActionRow(categorySelect2),
		navRow,
	}
}
