// File: internal/discord/ui/emojis.go
// Phiên bản: v0.1.2
// Mục đích: Registry emoji dùng trong toàn bộ UI bot.
// Bảo mật: Emoji là server-side constant — không nhận emoji ID từ user input.
// Ghi chú: Hiện dùng standard Unicode emoji (không cần upload lên server Discord).
//          Khi muốn dùng custom emoji riêng: đặt Name=tên_emoji, ID=snowflake_id.
//          Discord phân biệt 2 loại ComponentEmoji:
//            - Unicode  : ID="" → Discord dùng Name là ký tự unicode ("💎")
//            - Custom   : ID≠"" → Discord validate emoji tồn tại trong server bot đã join

package ui

import "fmt"

// Emoji bọc định nghĩa emoji dùng trong UI Discord.
// ID để trống ("") = unicode emoji, ID có giá trị = custom server emoji.
type Emoji struct {
	Name     string // Ký tự unicode ("💎") HOẶC tên custom emoji ("spirit_stone")
	ID       string // Snowflake ID của custom emoji — để "" nếu dùng unicode
	Animated bool   // true = custom animated emoji (chỉ dùng khi ID ≠ "")
}

// String trả về chuỗi embed-compatible cho emoji.
// - Unicode emoji (ID=""): trả về ký tự trực tiếp, ví dụ "💎"
// - Custom emoji (ID≠""): trả về format "<:name:id>" hoặc "<a:name:id>"
func (e Emoji) String() string {
	if e.ID == "" {
		return e.Name // unicode — dùng trực tiếp trong embed field name, description
	}
	if e.Animated {
		return fmt.Sprintf("<a:%s:%s>", e.Name, e.ID)
	}
	return fmt.Sprintf("<:%s:%s>", e.Name, e.ID)
}

// --- Tài nguyên ---
var (
	EmojiSpiritStone = Emoji{Name: "💎"} // Linh Thạch
	EmojiSpiritJade  = Emoji{Name: "🔮"} // Linh Ngọc
	EmojiFateTicket  = Emoji{Name: "🎫"} // Vé Cơ Duyên
)

// --- Tu luyện ---
var (
	EmojiCultivate    = Emoji{Name: "🧘"} // Tu luyện / Tĩnh tu
	EmojiRealm        = Emoji{Name: "🌟"} // Cảnh giới
	EmojiBreakthrough = Emoji{Name: "⚡"} // Đột phá
	EmojiStamina      = Emoji{Name: "💪"} // Thể lực
	EmojiMindState    = Emoji{Name: "🌸"} // Tâm cảnh
)

// --- Chiến đấu ---
var (
	EmojiSword       = Emoji{Name: "⚔️"} // Chiến đấu
	EmojiCombatPower = Emoji{Name: "⚡"}  // Chiến lực
	EmojiPvP         = Emoji{Name: "🥊"}  // PvP
	EmojiBoss        = Emoji{Name: "🐉"}  // Boss
)

// --- Hồ sơ / Xã hội ---
var (
	EmojiProfile = Emoji{Name: "👤"} // Hồ sơ
	EmojiSect    = Emoji{Name: "🏯"} // Tông môn
	EmojiPartner = Emoji{Name: "💕"} // Đạo lữ
	EmojiNPC     = Emoji{Name: "🤖"} // NPC
)

// --- Túi đồ / Trang bị ---
var (
	EmojiBag   = Emoji{Name: "🎒"}  // Túi đồ
	EmojiEquip = Emoji{Name: "🛡️"} // Trang bị
	EmojiSkill = Emoji{Name: "📖"}  // Kỹ năng / Công pháp
	EmojiPet   = Emoji{Name: "🐾"}  // Linh thú
)

// --- Thị trường ---
var (
	EmojiGacha   = Emoji{Name: "🎰"} // Cơ duyên / Gacha
	EmojiMarket  = Emoji{Name: "🏪"} // Chợ
	EmojiAuction = Emoji{Name: "🔨"} // Đấu giá
)

// --- Hệ thống ---
var (
	EmojiRefresh = Emoji{Name: "🔄"}  // Làm mới
	EmojiBack    = Emoji{Name: "◀️"} // Quay lại
	EmojiClose   = Emoji{Name: "✖️"} // Đóng
	EmojiInfo    = Emoji{Name: "ℹ️"} // Thông tin
	EmojiSuccess = Emoji{Name: "✅"}  // Thành công
	EmojiWarning = Emoji{Name: "⚠️"} // Cảnh báo
	EmojiError   = Emoji{Name: "❌"}  // Lỗi
	EmojiLock    = Emoji{Name: "🔒"}  // Khoá (coming soon)
)
