// File: internal/discord/menu/cultivation_menu.go
// Version: v0.1
// Purpose: UI builder for the Cultivation page — shows realm, exp, and tu luyện actions.
// Security: Pure rendering only. No DB calls or business logic.
// Notes: Cultivation actions (tĩnh tu, đột phá) will be wired to cooldown checks in v0.2.

package menu

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/yourname/tu-tien-bot/internal/discord/ui"
	cultivation "github.com/yourname/tu-tien-bot/internal/game/cultivation"
	profile "github.com/yourname/tu-tien-bot/internal/game/profile"
	"github.com/yourname/tu-tien-bot/pkg/utils"
)

// CultivationMenuData bundles all data needed to render the cultivation page.
type CultivationMenuData struct {
	Session     *Session
	Player      *profile.Player
	Cultivation *cultivation.CultivationProfile
}

// BuildCultivationMenuResponse constructs the interaction response data for cultivation page.
func BuildCultivationMenuResponse(data *CultivationMenuData) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Embeds:     []*discordgo.MessageEmbed{buildCultivationEmbed(data)},
		Components: buildCultivationComponents(data.Session.SessionID, data.Cultivation),
	}
}

func buildCultivationEmbed(data *CultivationMenuData) *discordgo.MessageEmbed {
	c := data.Cultivation
	p := data.Player

	realmFull := fmt.Sprintf("%s tầng %d", c.Realm.DisplayName(), c.RealmLevel)
	expBar := utils.ProgressBar(int(c.CultivationExp), int(c.CultivationExpRequired), 12)
	expText := fmt.Sprintf("`%s`\n%s / %s tu vi",
		expBar,
		utils.FormatNumber(c.CultivationExp),
		utils.FormatNumber(c.CultivationExpRequired),
	)

	pathDisplay := "Chưa chọn đạo lộ"
	if c.Path != cultivation.PathNone {
		pathDisplay = string(c.Path)
	}

	return &discordgo.MessageEmbed{
		Title:       ui.EmojiCultivate.String() + " Tu Luyện — " + p.DaoName,
		Description: "Con đường tu tiên vạn dặm bắt đầu từ một bước nhỏ.",
		Color:       ui.ColorCultivate,
		Fields: []*discordgo.MessageEmbedField{
			{Name: ui.EmojiRealm.String() + " Cảnh Giới", Value: realmFull, Inline: true},
			{Name: ui.EmojiMindState.String() + " Tâm Cảnh", Value: c.MindState.DisplayName(), Inline: true},
			{Name: ui.EmojiSkill.String() + " Đạo Lộ", Value: pathDisplay, Inline: true},
			{Name: ui.EmojiStamina.String() + " Thể Lực",
				Value: fmt.Sprintf("`%s` %d / %d",
					utils.ProgressBar(c.Stamina, c.MaxStamina, 10),
					c.Stamina, c.MaxStamina,
				),
				Inline: false,
			},
			{Name: ui.EmojiCultivate.String() + " Tiến Độ Tu Vi", Value: expText, Inline: false},
			{Name: ui.EmojiCombatPower.String() + " Chiến Lực", Value: utils.FormatNumber(c.CombatPower), Inline: true},
		},
		Footer:    &discordgo.MessageEmbedFooter{Text: "Vạn Pháp Tiên Nghịch · Tu Luyện"},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func buildCultivationComponents(sessionID string, c *cultivation.CultivationProfile) []discordgo.MessageComponent {
	canBreakthrough := c.CultivationExp >= c.CultivationExpRequired

	row1 := ui.ActionRow(
		// TODO v0.2: wire these to real cooldown-checked actions
		ui.Button("Tĩnh Tu", "cultivation:meditate:"+sessionID, ui.BtnPrimary, &ui.EmojiCultivate, false),
		ui.Button("Bế Quan", "cultivation:closeddoor:"+sessionID, ui.BtnPrimary, &ui.EmojiLock, false),
		ui.Button("Đột Phá", "cultivation:breakthrough:"+sessionID, ui.BtnSuccess, &ui.EmojiBreakthrough, !canBreakthrough),
	)

	navRow := ui.NavRow(sessionID, string(PageCultivation), string(PageMain))

	return []discordgo.MessageComponent{row1, navRow}
}
