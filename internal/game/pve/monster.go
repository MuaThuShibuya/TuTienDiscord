// File: internal/game/pve/monster.go
// Chức năng: Định nghĩa cấu trúc của quái vật và các hồ quái (Monster Pool).

package pve

type MonsterRole string

const (
	MonsterRoleNormal MonsterRole = "normal"
	MonsterRoleElite  MonsterRole = "elite"
	MonsterRoleBoss   MonsterRole = "boss"
)

type MonsterStats struct {
	MaxHP int64 `bson:"maxHp" json:"maxHp"`
	ATK   int64 `bson:"atk" json:"atk"`
	DEF   int64 `bson:"def" json:"def"`
	Speed int64 `bson:"speed" json:"speed"`
}

type MonsterDefinition struct {
	ID          string       `bson:"id" json:"id"`
	Name        string       `bson:"name" json:"name"`
	Description string       `bson:"description" json:"description"`
	Role        MonsterRole  `bson:"role" json:"role"`
	Element     string       `bson:"element" json:"element"` // kim, moc, thuy, hoa, tho
	BaseLevel   int          `bson:"baseLevel" json:"baseLevel"`
	BaseStats   MonsterStats `bson:"baseStats" json:"baseStats"`
	SkillIDs    []string     `bson:"skillIds" json:"skillIds"`
	Tags        []string     `bson:"tags" json:"tags"`
}

// Cấu trúc MonsterPool nằm trong registry.go (hoặc tách riêng nếu lớn).
