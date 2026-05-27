// File: internal/game/combat/calculator.go
// Chức năng: Bộ máy tính toán sát thương chiến đấu độc lập, không phụ thuộc Database hay Discord. Mở rộng cho các hiệu ứng Buff/Debuff.

package combat

import (
	"math/rand"
)

// DamageResult chứa kết quả của một đòn tấn công.
type DamageResult struct {
	Damage  int64
	IsCrit  bool
	IsMiss  bool
	Message string // Log text tuỳ chỉnh nếu cần
}

// Calculator định nghĩa hợp đồng cho bộ máy tính toán sát thương.
type Calculator interface {
	CalculateBasicAttack(attacker *CombatActor, defender *CombatActor, rng *rand.Rand) DamageResult
}

type defaultCalculator struct{}

func NewDamageCalculator() Calculator {
	return &defaultCalculator{}
}

// CalculateBasicAttack thực thi công thức sát thương v0.4.
// Nhận rng (Random Number Generator) từ ngoài vào để đảm bảo unit test luôn chạy đúng (Deterministic).
func (c *defaultCalculator) CalculateBasicAttack(attacker *CombatActor, defender *CombatActor, rng *rand.Rand) DamageResult {
	// TODO (v0.7): ApplyPreDamageModifiers (kiểm tra buff sát thương)
	// TODO (v0.7): ApplyDefenseModifiers (kiểm tra khiên/giáp)

	// 1. Tính sát thương thô (ATK - DEF * 0.5)
	baseDamage := float64(attacker.Stats.ATK) - (float64(defender.Stats.DEF) * 0.5)
	if baseDamage < 1 {
		baseDamage = 1 // Sát thương tối thiểu luôn là 1
	}

	// 2. Tính tỷ lệ bạo kích (Crit)
	isCrit := false
	if rng.Float64() < attacker.Stats.CritRate {
		isCrit = true
		baseDamage *= attacker.Stats.GetCritDamage()
	}

	// TODO (v0.7): Tính né tránh (HitRate vs DodgeRate)
	// TODO (v0.7): ApplyPostDamageModifiers

	return DamageResult{
		Damage: int64(baseDamage),
		IsCrit: isCrit,
		IsMiss: false,
	}
}
