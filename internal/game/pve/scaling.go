// File: internal/game/pve/scaling.go
// Chức năng: Tính toán độ khó tăng dần theo ải, loại hoạt động và vai trò quái vật.

package pve

type DifficultyScalingConfig struct {
	BaseMultiplier     float64
	PerStageMultiplier float64
	EliteMultiplier    float64
	BossMultiplier     float64
	BiCanhMultiplier   float64
}

var DefaultScalingConfig = DifficultyScalingConfig{
	BaseMultiplier:     1.0,
	PerStageMultiplier: 0.1, // Tăng 10% mỗi ải
	EliteMultiplier:    1.5,
	BossMultiplier:     3.0,
	BiCanhMultiplier:   1.3, // Bí cảnh khó hơn Du Ngoạn 30%
}

// ScaleMonsterStats biến đổi chỉ số gốc của quái vật theo tiến độ ải.
func ScaleMonsterStats(base MonsterStats, stage int, actType ActivityType, role MonsterRole, cfg DifficultyScalingConfig) MonsterStats {
	mult := cfg.BaseMultiplier + (float64(stage-1) * cfg.PerStageMultiplier)

	if actType == ActivityBiCanh {
		mult *= cfg.BiCanhMultiplier
	}
	if role == MonsterRoleElite {
		mult *= cfg.EliteMultiplier
	}
	if role == MonsterRoleBoss {
		mult *= cfg.BossMultiplier
	}

	// HP tăng mạnh hơn ATK
	scaledHP := float64(base.MaxHP) * (mult * 1.2)
	scaledATK := float64(base.ATK) * mult
	scaledDEF := float64(base.DEF) * (mult * 0.8) // DEF tăng chậm hơn

	// Tốc độ đánh giữ nguyên hoặc chỉ tăng rất nhẹ (chặn buff speed quá mức)
	scaledSpeed := float64(base.Speed) + float64(stage)
	if scaledSpeed > float64(base.Speed)+50 {
		scaledSpeed = float64(base.Speed) + 50
	}

	base.MaxHP = int64(scaledHP)
	base.ATK = int64(scaledATK)
	base.DEF = int64(scaledDEF)
	base.Speed = int64(scaledSpeed)
	return base
}
