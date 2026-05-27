// File: internal/game/skill/types.go
// Chức năng: Định nghĩa cấu trúc kỹ năng dùng chung cho Người chơi, Linh thú, Boss.

package skill

type SkillCost struct {
	Rage      int64   `bson:"rage" json:"rage"`           // Nộ khí
	Energy    int64   `bson:"energy" json:"energy"`       // Năng lượng (Mana/Linh lực)
	Stamina   int64   `bson:"stamina" json:"stamina"`     // Thể lực
	HPPercent float64 `bson:"hpPercent" json:"hpPercent"` // Tự hao tổn % máu
}

type SkillEffect struct {
	Type             string   `bson:"type" json:"type"` // damage, heal, shield, apply_status, turn_advance...
	Target           string   `bson:"target" json:"target"`
	Value            float64  `bson:"value" json:"value"`
	ScalingStat      string   `bson:"scalingStat" json:"scalingStat"` // atk, def, maxHp
	Chance           float64  `bson:"chance" json:"chance"`
	DurationTurns    int      `bson:"durationTurns" json:"durationTurns"`
	StatusEffectID   string   `bson:"statusEffectId" json:"statusEffectId"`
	TurnAdvanceValue float64  `bson:"turnAdvanceValue" json:"turnAdvanceValue"`
	Tags             []string `bson:"tags" json:"tags"`
}

type SkillDefinition struct {
	ID                string        `bson:"id" json:"id"`
	Name              string        `bson:"name" json:"name"`
	Description       string        `bson:"description" json:"description"`
	SkillType         string        `bson:"skillType" json:"skillType"` // basic, active, ultimate, passive
	Element           string        `bson:"element" json:"element"`     // kim, moc, thuy, hoa, tho
	TargetType        string        `bson:"targetType" json:"targetType"`
	Cost              SkillCost     `bson:"cost" json:"cost"`
	CooldownTurns     int           `bson:"cooldownTurns" json:"cooldownTurns"`
	Effects           []SkillEffect `bson:"effects" json:"effects"`
	Tags              []string      `bson:"tags" json:"tags"`
	AllowedActorTypes []string      `bson:"allowedActorTypes" json:"allowedActorTypes"` // player, pet, monster
}
