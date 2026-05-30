package npc

type ShopDef struct {
	ID    string
	Name  string
	Items map[string]ItemDef
}

type ItemDef struct {
	ItemID    string
	BuyPrice  int64  // Giá người chơi mua từ NPC
	SellPrice int64  // Giá NPC thu mua từ người chơi
	Stock     int    // -1 = vô hạn
	Category  string // Phân loại: Đan Dược, Trang Bị, Nguyên Liệu...
	Enabled   bool
}
