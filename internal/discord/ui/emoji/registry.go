// File: internal/discord/ui/emoji/registry.go
// Chức năng: Định nghĩa kiểu Emoji — hỗ trợ cả Unicode fallback lẫn Discord custom emoji ID.
// Ghi chú: Mọi file UI trong dự án phải import package này thay vì định nghĩa emoji riêng.
//          Để dùng custom emoji Discord: gọi emoji.LoadCustomEmojis() khi khởi động bot.
//          Khi chưa cấu hình custom ID → tự động dùng icon Unicode mặc định (fallback).

package emoji

import "github.com/bwmarrin/discordgo"

// Emoji đại diện cho một icon có thể dùng trong Discord embed và component.
// Hỗ trợ hai chế độ:
//   - Unicode fallback: dùng khi CustomID chưa được cấu hình (mặc định).
//   - Discord custom emoji: dùng khi CustomID được cấu hình qua LoadCustomEmojis().
type Emoji struct {
	Key      string // Định danh nội bộ, VD: "spirit_stone" — dùng để cấu hình custom ID
	Fallback string // Unicode emoji mặc định, VD: "💎" — luôn có sẵn khi không có custom ID
	customID string // Discord emoji snowflake ID (rỗng = dùng Fallback)
	animated bool   // true nếu là animated custom emoji
}

// String trả về chuỗi để nhúng vào embed description, field, hoặc title.
//   - Nếu có custom ID: trả về "<:name:id>" hoặc "<a:name:id>"
//   - Nếu không có custom ID: trả về Unicode fallback
func (e *Emoji) String() string {
	if e.customID != "" {
		if e.animated {
			return "<a:" + e.Key + ":" + e.customID + ">"
		}
		return "<:" + e.Key + ":" + e.customID + ">"
	}
	return e.Fallback
}

// Component trả về ComponentEmoji để dùng trong button emoji hoặc select option emoji.
//   - Nếu có custom ID: trả về cấu trúc custom emoji với Name + ID
//   - Nếu không có custom ID: trả về cấu trúc unicode emoji với Name = Fallback
func (e *Emoji) Component() *discordgo.ComponentEmoji {
	if e.customID != "" {
		return &discordgo.ComponentEmoji{
			Name:     e.Key,
			ID:       e.customID,
			Animated: e.animated,
		}
	}
	return &discordgo.ComponentEmoji{
		Name: e.Fallback,
	}
}

// CustomEmojiConfig cấu hình để áp dụng custom Discord emoji cho một icon.
type CustomEmojiConfig struct {
	ID       string // Discord emoji snowflake ID
	Animated bool   // true nếu là animated emoji
}

// LoadCustomEmojis áp dụng custom Discord emoji ID từ config map.
// Key của map phải khớp với Emoji.Key (VD: "spirit_stone", "cultivate", ...).
// Gọi hàm này một lần khi khởi động bot, sau khi có custom emoji trên server Discord.
//
// Ví dụ:
//
//	emoji.LoadCustomEmojis(map[string]emoji.CustomEmojiConfig{
//	    "spirit_stone": {ID: "1234567890123456789", Animated: false},
//	    "breakthrough":  {ID: "9876543210987654321", Animated: true},
//	})
func LoadCustomEmojis(configs map[string]CustomEmojiConfig) {
	for _, e := range All() {
		if cfg, ok := configs[e.Key]; ok && cfg.ID != "" {
			e.customID = cfg.ID
			e.animated = cfg.Animated
		}
	}
}
