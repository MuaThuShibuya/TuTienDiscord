// File: internal/game/pve/monster.go
// Chức năng: Định nghĩa cấu trúc của quái vật và các hồ quái (Monster Pool).

package pve

import "github.com/whiskey/tu-tien-bot/internal/game/combat"

type MonsterRole string

const (
	MonsterRoleNormal MonsterRole = "normal"
	MonsterRoleElite  MonsterRole = "elite"
	MonsterRoleBoss   MonsterRole = "boss"
)

type MonsterDefinition struct {
	ID          string             `bson:"id" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Role        MonsterRole        `bson:"role" json:"role"`
	Element     string             `bson:"element" json:"element"` // kim, moc, thuy, hoa, tho
	BaseLevel   int                `bson:"baseLevel" json:"baseLevel"`
	BaseStats   combat.CombatStats `bson:"baseStats" json:"baseStats"` // Được import an toàn từ combat
	SkillIDs    []string           `bson:"skillIds" json:"skillIds"`
	Tags        []string           `bson:"tags" json:"tags"`
}

// Cấu trúc MonsterPool nằm trong registry.go (hoặc tách riêng nếu lớn).
