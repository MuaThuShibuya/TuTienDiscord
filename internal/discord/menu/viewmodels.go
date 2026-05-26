// File: internal/discord/menu/viewmodels.go
// Chức năng: Định nghĩa các ViewModel — dữ liệu đã xử lý sẵn để UI Builder render.
// Ghi chú: Handler chuyển domain models → ViewModel trước khi gọi UI Builder.
//          UI Builder chỉ nhận ViewModel, không chứa business logic.

package menu

// MainMenuVM dữ liệu đã xử lý để render trang Main Menu.
type MainMenuVM struct {
	SessionID string

	DaoName      string
	RealmDisplay string // "Luyện Khí tầng 1"
	CombatPower  string // "1.234"
	MindState    string // "Bình Tĩnh (50/100)"
	PathDisplay  string // "Kiếm Tu"
	StaminaBar   string // "`████░░` 80/100"
	ExpBar       string // "`███░░░` 500/1.000"

	SpiritStones string // "1.234"
	SpiritJades  string // "0"
	FateTickets  string // "3 vé"

	DailyTip string
}

// ProfileMenuVM dữ liệu đã xử lý để render trang Hồ Sơ.
type ProfileMenuVM struct {
	SessionID   string
	DaoName     string
	DisplayName string
	JoinedAt    string // Discord timestamp "D"
	LastActive  string // Discord timestamp "R"

	SpiritStones string
	SpiritJades  string
	FateTickets  string
}

// CultivationMenuVM dữ liệu đã xử lý để render trang Tu Luyện.
type CultivationMenuVM struct {
	SessionID       string
	DaoName         string
	RealmDisplay    string
	MindState       string
	PathDisplay     string
	HasPath         bool
	StaminaBar      string
	ExpBar          string
	CombatPower     string
	CanBreakthrough bool
}

// InventoryMenuVM dữ liệu đã xử lý để render trang Túi Đồ.
type InventoryMenuVM struct {
	SessionID   string
	DaoName     string
	SlotUsage   string            // "5/50"
	Items       []InventoryItemVM // Items hiển thị trên trang hiện tại (đã phân trang)
	UsableItems []InventoryItemVM // Tất cả items có thể dùng (dùng cho select menu, tối đa 25)
	CurrentPage int
	TotalPages  int
}

// InventoryItemVM thông tin một vật phẩm trong túi đồ.
type InventoryItemVM struct {
	InstanceID string
	Name       string
	Quantity   int64
	Rarity     string // "D", "C", "B", "A", "S"
	IsUsable   bool   // Có thể dùng trực tiếp không (đan dược)
	IsEquip    bool   // Có phải trang bị không
}

// EquipmentMenuVM dữ liệu đã xử lý để render trang Trang Bị.
type EquipmentMenuVM struct {
	SessionID   string
	DaoName     string
	CombatPower string

	// Trang bị đang mặc theo từng vị trí (nil = trống)
	Weapon   *EquippedItemVM
	Armor    *EquippedItemVM
	Artifact *EquippedItemVM
	Treasure *EquippedItemVM
	Boots    *EquippedItemVM

	// Danh sách trang bị trong túi có thể mặc (tối đa 25 mục)
	Equippable []EquippableItemVM
}

// EquippedItemVM thông tin một trang bị đang mặc.
type EquippedItemVM struct {
	Slot     string // "weapon", "armor", ...
	SlotName string // "Vũ Khí", "Giáp", ...
	Name     string
	Rarity   string
}

// EquippableItemVM thông tin một trang bị trong túi có thể mặc.
type EquippableItemVM struct {
	InstanceID   string
	DefinitionID string
	Name         string
	Rarity       string
	SlotName     string // Vị trí trang bị phù hợp
}

// AlchemyMenuVM dữ liệu đã xử lý để render trang Lò Đan.
type AlchemyMenuVM struct {
	SessionID      string
	Level          int
	ExpBar         string
	Title          string     // Danh hiệu, VD: "Dược Đồng"
	DailyTip       string     // Mẹo luyện đan
	Recipes        []RecipeVM // Tất cả công thức
	SelectedRecipe *RecipeVM  // Nếu != nil, UI sẽ render trang chi tiết đan dược
}

// RecipeVM dữ liệu hiển thị một công thức luyện đan.
type RecipeVM struct {
	ID            string
	Name          string
	SuccessRate   string
	LevelRequired int
	Materials     string // Hiển thị nguyên liệu yêu cầu
	CanCraft      bool   // Trạng thái đủ nguyên liệu hay không
}
