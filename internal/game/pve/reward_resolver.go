// File: internal/game/pve/reward_resolver.go
// Chức năng: Đổ xúc xắc (Roll) để lấy phần thưởng từ Registry. Không thao tác với Database.

package pve

import (
	"fmt"
	"math/rand"
)

type ResolvedReward struct {
	Type     string
	RefID    string
	Quantity int64
	Rarity   string
	IsBonus  bool
}

type RewardRollResult struct {
	Guaranteed []ResolvedReward
	Bonus      []ResolvedReward
	Logs       []string
}

// ResolveStageRewards roll cả 2 pool để trả về phần thưởng cuối cùng.
func ResolveStageRewards(guaranteedPoolID, bonusPoolID string, rng *rand.Rand) RewardRollResult {
	res := RewardRollResult{}

	// 1. Guaranteed Roll (Quét toàn bộ, Chance quyết định có rớt hay không)
	if gPool, ok := RewardPoolRegistry[guaranteedPoolID]; ok {
		for _, entry := range gPool.Entries {
			if entry.Chance >= 1.0 || rng.Float64() <= entry.Chance {
				qty := entry.MinQuantity
				if entry.MaxQuantity > entry.MinQuantity {
					qty += rng.Int63n(entry.MaxQuantity - entry.MinQuantity + 1)
				}
				res.Guaranteed = append(res.Guaranteed, ResolvedReward{Type: entry.Type, RefID: entry.RefID, Quantity: qty, Rarity: entry.Rarity, IsBonus: false})
			}
		}
	}

	// 2. Bonus Gacha Roll (Trọng số - Weight). Chỉ roll trúng 1 món duy nhất.
	if bPool, ok := RewardPoolRegistry[bonusPoolID]; ok && len(bPool.Entries) > 0 {
		var totalWeight int64
		for _, entry := range bPool.Entries {
			totalWeight += entry.Weight
		}

		if totalWeight > 0 {
			roll := rng.Int63n(totalWeight)
			var current int64
			for _, entry := range bPool.Entries {
				current += entry.Weight
				if roll < current {
					// Trúng món này, kiểm tra Chance phụ (nếu có)
					if entry.Chance >= 1.0 || rng.Float64() <= entry.Chance {
						qty := entry.MinQuantity
						if entry.MaxQuantity > entry.MinQuantity {
							qty += rng.Int63n(entry.MaxQuantity - entry.MinQuantity + 1)
						}
						res.Bonus = append(res.Bonus, ResolvedReward{Type: entry.Type, RefID: entry.RefID, Quantity: qty, Rarity: entry.Rarity, IsBonus: true})

						if entry.Rarity == "S" || entry.Rarity == "SS" {
							res.Logs = append(res.Logs, fmt.Sprintf("Cơ duyên giáng lâm! Nhận được kỳ trân dị bảo cấp %s.", entry.Rarity))
						}
					}
					break
				}
			}
		}
	}

	if len(res.Bonus) > 0 {
		res.Logs = append(res.Logs, "Phát hiện thêm phần thưởng ẩn sau trận chiến!")
	}
	return res
}
