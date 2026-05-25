// File: internal/discord/menu/viewmodels.go
// Phiên bản: v0.1.1
// Mục đích: Định nghĩa các ViewModel — dữ liệu đã xử lý sẵn để UI Builder render.
//           Handler/Controller chuyển đổi domain models → ViewModel trước khi gọi UI Builder.
// Ghi chú: ViewModel không chứa business logic. Chỉ là data container đã format.
//          UI Builder (main_menu.go, profile_menu.go, ...) chỉ nhận ViewModel, không nhận domain model.

package menu

// MainMenuVM dữ liệu đã xử lý để render trang Main Menu.
type MainMenuVM struct {
	SessionID string // Nhúng vào custom_id để xác thực session
	DaoName   string // Đạo hiệu người chơi

	// Cảnh giới & tu luyện
	RealmDisplay string // Ví dụ: "Luyện Khí tầng 1"
	CombatPower  string // Ví dụ: "1.234" (đã format)
	MindState    string // Ví dụ: "Bình Tĩnh"
	StaminaBar   string // Ví dụ: "████░░░░░░ 80/100"
	ExpBar       string // Ví dụ: "███░░░░░░░ 500/1.000"

	// Tài nguyên
	SpiritStones string // Ví dụ: "1.234"
	SpiritJades  string // Ví dụ: "0"
	FateTickets  string // Ví dụ: "3 vé"

	// Gợi ý hàng ngày
	DailyTip string
}

// ProfileMenuVM dữ liệu đã xử lý để render trang Hồ Sơ.
type ProfileMenuVM struct {
	SessionID    string
	DaoName      string
	DisplayName  string // Tên Discord
	JoinedAt     string // Discord timestamp "D" — ngày tham gia
	LastActive   string // Discord timestamp "R" — lần cuối online
	SpiritStones string
	SpiritJades  string
	FateTickets  string
}

// CultivationMenuVM dữ liệu đã xử lý để render trang Tu Luyện.
type CultivationMenuVM struct {
	SessionID       string
	DaoName         string
	RealmDisplay    string // "Luyện Khí tầng 1"
	MindState       string // "Bình Tĩnh"
	PathDisplay     string // "Kiếm Tu" hoặc "Chưa chọn đạo lộ"
	StaminaBar      string // "████░░ 80/100"
	ExpBar          string // "███░░░ 500/1.000 tu vi"
	CombatPower     string // "1.234"
	CanBreakthrough bool   // Có thể đột phá không — dùng để enable/disable nút
}
