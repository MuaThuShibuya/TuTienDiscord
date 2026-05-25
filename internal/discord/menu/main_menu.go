// File: internal/discord/menu/main_menu.go
// Version: v0.1
// Purpose: Render the Main Menu embed and component layout for the Tu Tien RPG bot.
// Security: Only the session owner can interact. sessionId is embedded in all custom_ids.
// Notes: Main menu shows player snapshot + category navigation. Edit existing message, don't spam new ones.

package menu

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	cultivation "github.com/yourname/tu-tien-bot/internal/game/cultivation"
	economy "github.com/yourname/tu-tien-bot/internal/game/economy"
	profile "github.com/yourname/tu-tien-bot/internal/game/profile"
	"github.com/yourname/tu-tien-bot/pkg/utils"
)

// MainMenuData bundles all data needed to render the main menu.
type MainMenuData struct {
	Session     *Session
	Player      *profile.Player
	Cultivation *cultivation.CultivationProfile
	Wallet      *economy.Wallet
}

// BuildMainMenuResponse constructs the full Discord interaction response for the main menu.
func BuildMainMenuResponse(data *MainMenuData) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildMainEmbed(data)},
		Components: buildMainComponents(data.Session.SessionID),
	}
}

// BuildMainMenuEdit constructs an edit-message response for navigating back to main.
func BuildMainMenuEdit(data *MainMenuData) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildMainEmbed(data)},
		Components: buildMainComponents(data.Session.SessionID),
	}
}

func buildMainEmbed(data *MainMenuData) *discordgo.MessageEmbed {
	p := data.Player
	c := data.Cultivation
	w := data.Wallet

	realmDisplay := fmt.Sprintf("%s tầng %d", c.Realm.DisplayName(), c.RealmLevel)
	staminaBar := utils.ProgressBar(c.Stamina, c.MaxStamina, 10)
	expBar := utils.ProgressBar(int(c.CultivationExp), int(c.CultivationExpRequired), 10)

	tip := ui.DailyTips[rand.New(rand.NewSource(time.Now().Unix())).Intn(len(ui.DailyTips))]

	fields := []*discordgo.MessageEmbedField{
		{
			Name: ui.EmojiRealm.String() + " Cảnh Giới",
			Value: realmDisplay,
			Inline: true,
		},
		{
			Name: ui.EmojiCombatPower.String() + " Chiến Lực",
			Value: utils.FormatNumber(c.CombatPower),
			Inline: true,
		},
		{
			Name: ui.EmojiMindState.String() + " Tâm Cảnh",
			Value: c.MindState.DisplayName(),
			Inline: true,
		},
		{
			Name: ui.EmojiStamina.String() + " Thể Lực",
			Value: fmt.Sprintf("`%s` %d/%d", staminaBar, c.Stamina, c.MaxStamina),
			Inline: false,
		},
		{
			Name: ui.EmojiCultivate.String() + " Tu Vi",
			Value: fmt.Sprintf("`%s` %s/%s",
				expBar,
				utils.FormatNumber(c.CultivationExp),
				utils.FormatNumber(c.CultivationExpRequired)),
			Inline: false,
		},
		{
			Name: ui.EmojiSpiritStone.String() + " Linh Thạch",
			Value: utils.FormatNumber(w.SpiritStones),
			Inline: true,
		},
		{
			Name: ui.EmojiSpiritJade.String() + " Linh Ngọc",
			Value: utils.FormatNumber(w.SpiritJades),
			Inline: true,
		},
		{
			Name: ui.EmojiFateTicket.String() + " Vé Cơ Duyên",
			Value: fmt.Sprintf("%d vé", w.FateTickets),
			Inline: true,
		},
		{
			Name:  ui.EmojiInfo.String() + " Gợi ý hôm nay",
			Value: "_" + tip + "_",
		},
	}

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiProfile.String() + " Vạn Pháp Tiên Nghịch — " + p.DaoName,
		Description: fmt.Sprintf("Chào mừng trở lại, **%s**!\nHãy chọn chức năng bên dưới.", p.DaoName),
		Color:       ui.ColorDefault,
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Vạn Pháp Tiên Nghịch · v0.1",
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildMainComponents(sessionID string) []discordgo.MessageComponent {
	// Row 1: Category select menu
	categorySelect := ui.SelectMenu(
		"menu:nav:"+sessionID,
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

	// Row 2: Second set of categories
	categorySelect2 := ui.SelectMenu(
		"menu:nav2:"+sessionID,
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

	// Row 3: Nav buttons
	navRow := ui.NavRow(sessionID, string(PageMain), "")

	return []discordgo.MessageComponent{
		ui.ActionRow(categorySelect),
		ui.ActionRow(categorySelect2),
		navRow,
	}
}
