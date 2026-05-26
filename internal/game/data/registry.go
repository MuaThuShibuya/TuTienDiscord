// File: internal/data/registry.go
package data

import (
	"github.com/whiskey/tu-tien-bot/internal/game/alchemy"
	"github.com/whiskey/tu-tien-bot/internal/game/item"
)

var (
	AllItems   = make(map[string]item.ItemDefinition)
	AllRecipes = make(map[string]alchemy.Recipe)
)

// RegisterItems đăng ký vật phẩm vào kho dữ liệu chung.
func RegisterItems(items map[string]item.ItemDefinition) {
	for k, v := range items {
		AllItems[k] = v
		item.Definitions[k] = v
	}
}

// RegisterRecipes đăng ký công thức luyện đan vào kho dữ liệu chung.
func RegisterRecipes(recipes map[string]alchemy.Recipe) {
	for k, v := range recipes {
		AllRecipes[k] = v
		alchemy.Recipes[k] = v
	}
}
