// File: internal/game/item/definitions.go
package item

var Definitions = make(map[string]ItemDefinition)

// RegisterItems được gọi từ hàm init() của các file data để nạp vật phẩm vào hệ thống.
func RegisterItems(items map[string]ItemDefinition) {
	for k, v := range items {
		Definitions[k] = v
	}
}

func GetDefinition(id string) (ItemDefinition, bool) {
	def, ok := Definitions[id]
	return def, ok
}
