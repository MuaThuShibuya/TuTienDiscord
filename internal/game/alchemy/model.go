// File: internal/game/alchemy/model.go
// Chức năng: Định nghĩa cấu trúc dữ liệu và công thức luyện đan (Recipes).

package alchemy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AlchemyProfile lưu trữ cấp độ và kinh nghiệm luyện đan của người chơi.
type AlchemyProfile struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"userId"        json:"userId"`
	GuildID   string             `bson:"guildId"       json:"guildId"`
	Level     int                `bson:"level"         json:"level"`
	Exp       int64              `bson:"exp"           json:"exp"`
	CreatedAt time.Time          `bson:"createdAt"     json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt"     json:"updatedAt"`
}

// Recipe định nghĩa một công thức luyện chế đan dược.
type Recipe struct {
	ID             string
	Name           string
	RequiredItems  map[string]int64 // DefinitionID -> Số lượng yêu cầu
	OutputItem     string           // DefinitionID của đan dược tạo ra
	OutputQuantity int64
	SuccessRate    float64 // Tỉ lệ thành công (0.0 -> 1.0)
	LevelRequired  int     // Cấp luyện đan tối thiểu để luyện
	ExpReward      int64   // Điểm kinh nghiệm luyện đan nhận được nếu thành công
}

// Recipes danh sách công thức luyện đan của server. (Có thể mở rộng thêm)
var Recipes = map[string]Recipe{
	"recipe_qi_pill": {
		ID:   "recipe_qi_pill",
		Name: "Tụ Khí Đan",
		RequiredItems: map[string]int64{
			"mat_spirit_grass": 2, // 2 Nhất giai linh thảo
		},
		OutputItem:     "item_qi_pill",
		OutputQuantity: 1,
		SuccessRate:    0.80, // 80% thành công
		LevelRequired:  1,
		ExpReward:      10,
	},
	// TODO: Thêm Tẩy Tủy Đan, Trúc Cơ Đan...
}
