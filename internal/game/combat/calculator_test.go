// File: internal/game/combat/calculator_test.go
// Chức năng: Unit Test cho bộ máy tính toán sát thương (Damage Calculator), đảm bảo công thức bạo kích và chỉ số cơ bản luôn chính xác.

package combat

import (
	"math/rand"
	"testing"
)

func TestCalculator_CalculateBasicAttack(t *testing.T) {
	calc := NewDamageCalculator()

	tests := []struct {
		name       string
		attacker   CombatActor
		defender   CombatActor
		seed       int64 // Seed cố định cho rng
		wantDamage int64
		wantCrit   bool
	}{
		{
			name: "Sát thương bình thường (ATK > DEF)",
			attacker: CombatActor{
				Stats: CombatStats{ATK: 100, CritRate: 0.0}, // Không crit
			},
			defender: CombatActor{
				Stats: CombatStats{DEF: 50},
			},
			seed:       1,
			wantDamage: 75, // 100 - (50 * 0.5) = 75
			wantCrit:   false,
		},
		{
			name: "Sát thương tối thiểu bằng 1 (DEF cực lớn)",
			attacker: CombatActor{
				Stats: CombatStats{ATK: 10, CritRate: 0.0},
			},
			defender: CombatActor{
				Stats: CombatStats{DEF: 500},
			},
			seed:       1,
			wantDamage: 1, // 10 - 250 = -240 -> bị clamp về 1
			wantCrit:   false,
		},
		{
			name: "Đòn đánh bạo kích (CritRate = 1.0)",
			attacker: CombatActor{
				Stats: CombatStats{
					ATK:        100,
					CritRate:   1.0, // Chắc chắn crit
					CritDamage: 1.5,
				},
			},
			defender: CombatActor{
				Stats: CombatStats{DEF: 50},
			},
			seed:       1,
			wantDamage: 112, // (100 - 25) * 1.5 = 112.5 -> ép kiểu int64 thành 112
			wantCrit:   true,
		},
		{
			name: "CritDamage fallback mặc định 1.5",
			attacker: CombatActor{
				Stats: CombatStats{
					ATK:        100,
					CritRate:   1.0,
					CritDamage: 0.0, // Bỏ trống, fallback tự động về 1.5
				},
			},
			defender: CombatActor{
				Stats: CombatStats{DEF: 50},
			},
			seed:       1,
			wantDamage: 112,
			wantCrit:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sử dụng local rng với seed cố định cho từng test case
			rng := rand.New(rand.NewSource(tt.seed))

			res := calc.CalculateBasicAttack(&tt.attacker, &tt.defender, rng)

			if res.Damage != tt.wantDamage {
				t.Errorf("Damage mong đợi %d, nhận được %d", tt.wantDamage, res.Damage)
			}
			if res.IsCrit != tt.wantCrit {
				t.Errorf("Crit mong đợi %v, nhận được %v", tt.wantCrit, res.IsCrit)
			}
		})
	}
}
