// File: internal/discord/ui/emojis.go
// Version: v0.1
// Purpose: Central registry for all Discord custom emoji IDs used in the bot UI.
// Security: Emoji IDs are server-side constants. Never accept emoji IDs from user input.
// Notes: Replace the placeholder IDs with your actual Discord emoji IDs from your emoji server.
//        Format for inline use: <:name:id> for static, <a:name:id> for animated.
//        Use E() helper to build the formatted string.

package ui

import "fmt"

// Emoji wraps a Discord custom emoji definition.
type Emoji struct {
	Name     string
	ID       string
	Animated bool
}

// String returns the Discord embed-compatible emoji string.
func (e Emoji) String() string {
	if e.Animated {
		return fmt.Sprintf("<a:%s:%s>", e.Name, e.ID)
	}
	return fmt.Sprintf("<:%s:%s>", e.Name, e.ID)
}

// All custom emojis used in the Tu Tien bot.
// TODO: Replace placeholder IDs (000000000000000001, etc.) with your real Discord emoji IDs.
var (
	// --- Currency ---
	EmojiSpiritStone = Emoji{Name: "spirit_stone", ID: "000000000000000001"}   // Linh Thạch
	EmojiSpiritJade  = Emoji{Name: "spirit_jade",  ID: "000000000000000002"}   // Linh Ngọc
	EmojiFateTicket  = Emoji{Name: "fate_ticket",  ID: "000000000000000003"}   // Vé Cơ Duyên

	// --- Cultivation ---
	EmojiCultivate   = Emoji{Name: "cultivate",    ID: "000000000000000010"}   // Tu luyện
	EmojiRealm       = Emoji{Name: "realm",        ID: "000000000000000011"}   // Cảnh giới
	EmojiBreakthrough= Emoji{Name: "breakthrough", ID: "000000000000000012"}   // Đột phá
	EmojiStamina     = Emoji{Name: "stamina",      ID: "000000000000000013"}   // Thể lực
	EmojiMindState   = Emoji{Name: "mind_state",   ID: "000000000000000014"}   // Tâm cảnh

	// --- Combat ---
	EmojiSword       = Emoji{Name: "sword",        ID: "000000000000000020"}   // Chiến đấu
	EmojiCombatPower = Emoji{Name: "combat_power", ID: "000000000000000021"}   // Chiến lực
	EmojiPvP         = Emoji{Name: "pvp",          ID: "000000000000000022"}   // PvP
	EmojiBoss        = Emoji{Name: "boss",         ID: "000000000000000023"}   // Boss

	// --- Profile / Social ---
	EmojiProfile     = Emoji{Name: "profile",      ID: "000000000000000030"}   // Hồ sơ
	EmojiSect        = Emoji{Name: "sect",         ID: "000000000000000031"}   // Tông môn
	EmojiPartner     = Emoji{Name: "partner",      ID: "000000000000000032"}   // Đạo lữ
	EmojiNPC         = Emoji{Name: "npc",          ID: "000000000000000033"}   // NPC

	// --- Inventory / Equipment ---
	EmojiBag         = Emoji{Name: "bag",          ID: "000000000000000040"}   // Túi đồ
	EmojiEquip       = Emoji{Name: "equip",        ID: "000000000000000041"}   // Trang bị
	EmojiSkill       = Emoji{Name: "skill",        ID: "000000000000000042"}   // Kỹ năng
	EmojiPet         = Emoji{Name: "pet",          ID: "000000000000000043"}   // Linh thú

	// --- Market ---
	EmojiGacha       = Emoji{Name: "gacha",        ID: "000000000000000050"}   // Cơ duyên / Gacha
	EmojiMarket      = Emoji{Name: "market",       ID: "000000000000000051"}   // Chợ
	EmojiAuction     = Emoji{Name: "auction",      ID: "000000000000000052"}   // Đấu giá

	// --- System ---
	EmojiRefresh     = Emoji{Name: "refresh",      ID: "000000000000000060"}   // Làm mới
	EmojiBack        = Emoji{Name: "back",         ID: "000000000000000061"}   // Quay lại
	EmojiClose       = Emoji{Name: "close",        ID: "000000000000000062"}   // Đóng
	EmojiInfo        = Emoji{Name: "info",         ID: "000000000000000063"}   // Thông tin
	EmojiSuccess     = Emoji{Name: "success",      ID: "000000000000000064"}   // Thành công
	EmojiWarning     = Emoji{Name: "warning",      ID: "000000000000000065"}   // Cảnh báo
	EmojiError       = Emoji{Name: "error_icon",   ID: "000000000000000066"}   // Lỗi
	EmojiLock        = Emoji{Name: "lock",         ID: "000000000000000067"}   // Khoá (coming soon)
)
