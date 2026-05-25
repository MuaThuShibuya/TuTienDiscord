// File: internal/game/item/definitions.go
package item

var Definitions = map[string]ItemDefinition{
	// Đan Dược
	"pill_exp_small":               {ID: "pill_exp_small", Name: "Tiểu Tu Vi Đan", Type: TypePill, Rarity: RarityD, Stackable: true, MaxStack: 999, Usable: true},
	"pill_exp_medium":              {ID: "pill_exp_medium", Name: "Trung Tu Vi Đan", Type: TypePill, Rarity: RarityC, Stackable: true, MaxStack: 999, Usable: true},
	"pill_stamina_small":           {ID: "pill_stamina_small", Name: "Tiểu Hồi Thể Đan", Type: TypePill, Rarity: RarityD, Stackable: true, MaxStack: 999, Usable: true},
	"pill_breakthrough_rate_small": {ID: "pill_breakthrough_rate_small", Name: "Trúc Cơ Hộ Mệnh Đan", Type: TypePill, Rarity: RarityB, Stackable: true, MaxStack: 999, Usable: false},

	// Trang bị
	"eq_wood_sword":    {ID: "eq_wood_sword", Name: "Mộc Kiếm", Type: TypeEquipment, Rarity: RarityD, Stackable: false, MaxStack: 1, Usable: false},
	"eq_iron_sword":    {ID: "eq_iron_sword", Name: "Thiết Kiếm", Type: TypeEquipment, Rarity: RarityC, Stackable: false, MaxStack: 1, Usable: false},
	"eq_cloth_robe":    {ID: "eq_cloth_robe", Name: "Bố Y", Type: TypeEquipment, Rarity: RarityD, Stackable: false, MaxStack: 1, Usable: false},
	"eq_iron_armor":    {ID: "eq_iron_armor", Name: "Thiết Giáp", Type: TypeEquipment, Rarity: RarityC, Stackable: false, MaxStack: 1, Usable: false},
	"eq_spirit_bell":   {ID: "eq_spirit_bell", Name: "Hám Hồn Chung", Type: TypeEquipment, Rarity: RarityB, Stackable: false, MaxStack: 1, Usable: false},
	"eq_guardian_jade": {ID: "eq_guardian_jade", Name: "Hộ Mệnh Ngọc", Type: TypeEquipment, Rarity: RarityB, Stackable: false, MaxStack: 1, Usable: false},

	// Nguyên liệu
	"refine_stone": {ID: "refine_stone", Name: "Đá Cường Hóa", Type: TypeMaterial, Rarity: RarityD, Stackable: true, MaxStack: 999, Usable: false},
}

func GetDefinition(id string) (ItemDefinition, bool) {
	def, ok := Definitions[id]
	return def, ok
}
