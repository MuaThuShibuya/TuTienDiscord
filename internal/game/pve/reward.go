// File: internal/game/pve/reward.go
// Chức năng: Định nghĩa cấu trúc Reward Pool (Hồ phần thưởng) riêng biệt cho từng khu vực PvE.

package pve

type RewardPoolDefinition struct {
	ID           string        `bson:"id" json:"id"`
	Name         string        `bson:"name" json:"name"`
	ActivityType ActivityType  `bson:"activityType" json:"activityType"`
	Entries      []RewardEntry `bson:"entries" json:"entries"`
}

type RewardEntry struct {
	Type        string  `bson:"type" json:"type"`               // exp, stones, item, material, equipment, artifact, skill_scroll
	RefID       string  `bson:"refId" json:"refId"`             // ID của vật phẩm/trang bị/kỹ năng
	MinQuantity int64   `bson:"minQuantity" json:"minQuantity"` // Số lượng tối thiểu
	MaxQuantity int64   `bson:"maxQuantity" json:"maxQuantity"` // Số lượng tối đa
	Chance      float64 `bson:"chance" json:"chance"`           // Tỷ lệ rớt (0.0 -> 1.0)
	Weight      int64   `bson:"weight" json:"weight"`           // Trọng số (dùng cho gacha logic nếu cần)
	Rarity      string  `bson:"rarity" json:"rarity"`
	IsUnique    bool    `bson:"isUnique" json:"isUnique"` // Chỉ rớt 1 lần duy nhất trong pool này?
}

// TODO: v0.4.x - Tạo RewardResolver để tính toán drop dựa trên RewardPoolDefinition.
