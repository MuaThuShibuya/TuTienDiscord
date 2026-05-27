package cultivation

import "github.com/whiskey/tu-tien-bot/internal/apperrors"

// GetRealmIndex trả về vị trí của cảnh giới trong mảng.
func GetRealmIndex(realm string) int {
	id := NormalizeRealmID(realm)
	if def, ok := RealmRegistry[id]; ok {
		return def.Order
	}
	return 0 // Fallback
}

// NextRealm tính toán cảnh giới tiếp theo, tuân thủ MaxLevel của từng cảnh giới.
func NextRealm(currentRealm string, currentLevel int) (nextRealm string, nextLevel int, advancedRealm bool, err error) {
	id := NormalizeRealmID(currentRealm)
	def, ok := RealmRegistry[id]
	if !ok {
		return "pham_nhan", 1, false, nil
	}

	// Tăng tầng
	if currentLevel < def.MaxLevel {
		return id, currentLevel + 1, false, nil
	}

	// Đột phá cảnh giới lớn
	if def.Order >= len(RealmOrderList)-1 {
		return currentRealm, currentLevel, false, apperrors.ErrMaxRealmReached
	}

	nextID := RealmOrderList[def.Order+1]
	return nextID, 1, true, nil
}

// CalculateNextExpRequired tính tu vi cần thiết.
func CalculateNextExpRequired(realm string, realmLevel int) int64 {
	base := int64(100)
	idx := GetRealmIndex(realm)
	if idx == 0 {
		idx = 1
	} // Tránh x0

	realmMultiplier := int64(idx * 5)
	return base * realmMultiplier * int64(realmLevel+1)
}

func CalculateBreakthroughCost(realm string, realmLevel int) int64 {
	idx := GetRealmIndex(realm)
	return int64(100 * (idx + 1) * realmLevel)
}
