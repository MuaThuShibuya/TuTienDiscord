// File: internal/game/cultivation/paths.go
// Phiên bản: v0.2.1
// Mục đích: Helper cô lập các công thức áp dụng hiệu ứng (Bonus) của từng Đạo lộ.
// Ghi chú: Không nhúng trực tiếp logic vào service để giữ code sạch và dễ viết Unit Test.

package cultivation

import "math"

// ApplyPathExpBonus tăng tu vi theo Đạo lộ.
func ApplyPathExpBonus(path CultivationPath, action string, baseExp int64) int64 {
	bonus := float64(0)
	switch path {
	case PathSword:
		if action == "meditate" {
			bonus = 0.05
		}
	case PathSpirit:
		if action == "meditate" || action == "seclusion" {
			bonus = 0.10
		}
	case PathPoison:
		if action == "seclusion" {
			bonus = 0.15
		}
	}

	extra := int64(math.Floor(float64(baseExp) * bonus))
	if bonus > 0 && extra == 0 {
		extra = 1
	} // Đảm bảo luôn được ít nhất +1 nếu có bonus

	return baseExp + extra
}

// ApplyPathStaminaCostBonus giảm tiêu hao thể lực theo Đạo lộ (tối thiểu 1).
func ApplyPathStaminaCostBonus(path CultivationPath, action string, baseCost int) int {
	cost := baseCost
	if path == PathBody {
		if action == "meditate" {
			cost -= 1
		}
		if action == "body_training" {
			cost -= 2
		}
	}
	if cost < 1 {
		return 1
	}
	return cost
}

// ApplyPathCombatPowerBonus tăng chiến lực khi Luyện thể.
func ApplyPathCombatPowerBonus(path CultivationPath, action string, baseCP int64) int64 {
	bonus := float64(0)
	if action == "body_training" {
		if path == PathSword {
			bonus = 0.10
		}
		if path == PathBody {
			bonus = 0.05
		}
	}

	extra := int64(math.Floor(float64(baseCP) * bonus))
	if bonus > 0 && extra == 0 {
		extra = 1
	}

	return baseCP + extra
}

// ApplyPathBreakthroughRateBonus tăng tỉ lệ thành công khi Đột phá.
func ApplyPathBreakthroughRateBonus(path CultivationPath, baseRate float64) float64 {
	if path == PathSpirit {
		return baseRate + 0.03
	}
	if path == PathPoison {
		return baseRate + 0.05
	}
	return baseRate
}

// ApplyPathBreakthroughFailurePenalty tính toán lượng tâm cảnh mất thêm khi Đột phá thất bại.
func ApplyPathBreakthroughFailurePenalty(path CultivationPath, baseLoss int) int {
	if path == PathPoison {
		return baseLoss + 2
	}
	return baseLoss
}
