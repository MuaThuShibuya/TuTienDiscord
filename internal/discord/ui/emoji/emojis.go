// File: internal/discord/ui/emoji/emojis.go
// Chức năng: Khai báo toàn bộ emoji dùng trong UI bot — với icon Unicode mặc định.
// Ghi chú: Để thêm emoji mới: khai báo var mới + thêm vào All().
//          Để dùng custom Discord emoji: gọi LoadCustomEmojis() với ID tương ứng.
//          Key phải dùng snake_case, không dấu — dùng để tra trong LoadCustomEmojis().

package emoji

// --- Tài nguyên (Economy) ---
var (
	SpiritStone = &Emoji{Key: "spirit_stone", Fallback: "💎"} // Linh thạch
	SpiritJade  = &Emoji{Key: "spirit_jade", Fallback: "🔮"}  // Linh ngọc
	FateTicket  = &Emoji{Key: "fate_ticket", Fallback: "🎫"}  // Vé cơ duyên
)

// --- Tu luyện (Cultivation) ---
var (
	Cultivate    = &Emoji{Key: "cultivate", Fallback: "🧘"}    // Tĩnh tu
	Realm        = &Emoji{Key: "realm", Fallback: "🌟"}        // Cảnh giới
	Breakthrough = &Emoji{Key: "breakthrough", Fallback: "⚡"}  // Đột phá
	Stamina      = &Emoji{Key: "stamina", Fallback: "💪"}      // Thể lực
	MindState    = &Emoji{Key: "mind_state", Fallback: "🌸"}   // Tâm cảnh
)

// --- Chiến đấu (Combat) ---
var (
	Sword       = &Emoji{Key: "sword", Fallback: "⚔️"}       // Chiến đấu / Kiếm
	CombatPower = &Emoji{Key: "combat_power", Fallback: "🔥"} // Chiến lực
	PvP         = &Emoji{Key: "pvp", Fallback: "🥊"}         // PvP
	Boss        = &Emoji{Key: "boss", Fallback: "🐉"}        // Boss server
)

// --- Hồ sơ / Xã hội (Social) ---
var (
	Profile = &Emoji{Key: "profile", Fallback: "👤"} // Hồ sơ người chơi
	Sect    = &Emoji{Key: "sect", Fallback: "🏯"}    // Tông môn
	Partner = &Emoji{Key: "partner", Fallback: "💕"} // Đạo lữ
	NPC     = &Emoji{Key: "npc", Fallback: "🤖"}    // NPC
)

// --- Túi đồ / Trang bị (Inventory & Equipment) ---
var (
	Bag   = &Emoji{Key: "bag", Fallback: "🎒"}   // Túi đồ
	Equip = &Emoji{Key: "equip", Fallback: "🛡️"} // Trang bị
	Skill = &Emoji{Key: "skill", Fallback: "📖"} // Kỹ năng / Công pháp
	Pet   = &Emoji{Key: "pet", Fallback: "🐾"}   // Linh thú
)

// --- Luyện đan / Đặc biệt ---
var (
	Alchemy = &Emoji{Key: "alchemy", Fallback: "⚗️"} // Luyện đan
	Quest   = &Emoji{Key: "quest", Fallback: "📜"}   // Nhiệm vụ
)

// --- Thị trường (Market) ---
var (
	Gacha   = &Emoji{Key: "gacha", Fallback: "🎰"}  // Cơ duyên / Gacha
	Market  = &Emoji{Key: "market", Fallback: "🏪"} // Chợ trao đổi
	Auction = &Emoji{Key: "auction", Fallback: "🔨"} // Đấu giá
)

// --- Hệ thống (System) ---
var (
	Refresh = &Emoji{Key: "refresh", Fallback: "🔄"}  // Làm mới
	Back    = &Emoji{Key: "back", Fallback: "◀️"}     // Quay lại
	Close   = &Emoji{Key: "close", Fallback: "✖️"}   // Đóng
	Info    = &Emoji{Key: "info", Fallback: "ℹ️"}    // Thông tin
	Success = &Emoji{Key: "success", Fallback: "✅"}  // Thành công
	Warning = &Emoji{Key: "warning", Fallback: "⚠️"} // Cảnh báo
	Error   = &Emoji{Key: "error", Fallback: "❌"}   // Lỗi
	Lock    = &Emoji{Key: "lock", Fallback: "🔒"}    // Khoá (coming soon)
)

// All trả về toàn bộ emoji trong registry.
// Dùng bởi LoadCustomEmojis để duyệt và áp dụng custom ID.
func All() []*Emoji {
	return []*Emoji{
		// Tài nguyên
		SpiritStone, SpiritJade, FateTicket,
		// Tu luyện
		Cultivate, Realm, Breakthrough, Stamina, MindState,
		// Chiến đấu
		Sword, CombatPower, PvP, Boss,
		// Xã hội
		Profile, Sect, Partner, NPC,
		// Túi đồ / Trang bị
		Bag, Equip, Skill, Pet,
		// Đặc biệt
		Alchemy, Quest,
		// Thị trường
		Gacha, Market, Auction,
		// Hệ thống
		Refresh, Back, Close, Info, Success, Warning, Error, Lock,
	}
}
