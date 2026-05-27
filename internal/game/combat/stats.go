// File: internal/game/combat/stats.go
// Chức năng: Định nghĩa các chỉ số chiến đấu cơ bản (CombatStats), cơ chế biến đổi chỉ số (StatModifier) và hiệu ứng trạng thái (StatusEffect) cho v0.4+.

package combat

// CombatStats đại diện cho toàn bộ chỉ số chiến đấu của một thực thể (Player/Monster).
type CombatStats struct {
	MaxHP       int64   `bson:"maxHp" json:"maxHp"`
	ATK         int64   `bson:"atk" json:"atk"`
	DEF         int64   `bson:"def" json:"def"`
	CritRate    float64 `bson:"critRate" json:"critRate"`       // 0.0 -> 1.0 (ví dụ: 0.15 = 15%)
	CritDamage  float64 `bson:"critDamage" json:"critDamage"`   // Tỷ lệ nhân sát thương bạo kích (mặc định 1.5)
	Speed       int64   `bson:"speed" json:"speed"`             // Quyết định thứ tự đi trước/sau
	HitRate     float64 `bson:"hitRate" json:"hitRate"`         // Tỷ lệ chính xác
	DodgeRate   float64 `bson:"dodgeRate" json:"dodgeRate"`     // Tỷ lệ né tránh
	CombatPower int64   `bson:"combatPower" json:"combatPower"` // Lực chiến tổng (chỉ dùng hiển thị)
}

type ModifierMode string

const (
	ModeFlat       ModifierMode = "flat"       // Cộng/trừ thẳng (ví dụ: +50 ATK)
	ModePercent    ModifierMode = "percent"    // Cộng/trừ theo phần trăm Base (ví dụ: +10% ATK)
	ModeMultiplier ModifierMode = "multiplier" // Nhân tổng lực sau cùng (ví dụ: x1.2 Damage)
)

// StatModifier lưu trữ thông tin về một phép biến đổi chỉ số.
type StatModifier struct {
	SourceType    string       `bson:"sourceType" json:"sourceType"` // eq: weapon, buff: skill, aura: sect
	SourceID      string       `bson:"sourceId" json:"sourceId"`
	Stat          string       `bson:"stat" json:"stat"` // atk, def, maxHp...
	Mode          ModifierMode `bson:"mode" json:"mode"`
	Value         float64      `bson:"value" json:"value"`
	DurationTurns int          `bson:"durationTurns" json:"durationTurns"` // -1 = vĩnh viễn (passive)
}

// StatusEffect đại diện cho Buff/Debuff/DoT/HoT áp dụng lên Actor.
// TODO: v0.4 chưa dùng nhiều, nhưng đặt nền móng cho v0.7 (Buff/Debuff/DoT).
type StatusEffect struct {
	ID             string         `bson:"id" json:"id"`
	Name           string         `bson:"name" json:"name"`
	Type           string         `bson:"type" json:"type"` // buff, debuff, dot, hot, control
	Modifiers      []StatModifier `bson:"modifiers" json:"modifiers"`
	RemainingTurns int            `bson:"remainingTurns" json:"remainingTurns"`
	Stack          int            `bson:"stack" json:"stack"`
	MaxStack       int            `bson:"maxStack" json:"maxStack"`
	SourceActorID  string         `bson:"sourceActorId" json:"sourceActorId"` // Ai buff cái này?
}

// GetCritDamage an toàn, nếu chưa config sẽ mặc định là 1.5 (150%)
func (s *CombatStats) GetCritDamage() float64 {
	if s.CritDamage <= 0 {
		return 1.5
	}
	return s.CritDamage
}
