package cultivation

import "github.com/whiskey/tu-tien-bot/internal/apperrors"

var RealmOrder = []string{
	"Phàm Nhân",
	"Luyện Khí",
	"Trúc Cơ",
	"Kim Đan",
	"Nguyên Anh",
	"Hóa Thần",
}

// GetRealmIndex trả về vị trí của cảnh giới trong mảng.
func GetRealmIndex(realm string) int {
	for i, r := range RealmOrder {
		if r == realm {
			return i
		}
	}
	return 0 // Fallback
}

// NextRealm tính toán cảnh giới tiếp theo dựa trên logic v0.2 (10 tầng).
func NextRealm(currentRealm string, currentLevel int) (nextRealm string, nextLevel int, advancedRealm bool, err error) {
	idx := GetRealmIndex(currentRealm)

	// Đang ở đỉnh phong Hóa Thần
	if idx == len(RealmOrder)-1 && currentLevel >= 10 {
		return currentRealm, currentLevel, false, apperrors.ErrMaxRealmReached
	}

	// Tăng tầng
	if currentLevel < 10 {
		return currentRealm, currentLevel + 1, false, nil
	}

	// Đột phá cảnh giới lớn
	return RealmOrder[idx+1], 1, true, nil
}

// CalculateNextExpRequired tính tu vi cần thiết theo công thức v0.2.
// Công thức: base(100) * (realmIndex + 1) * (realmLevel + 1)
func CalculateNextExpRequired(realm string, realmLevel int) int64 {
	base := int64(100)
	realmMultiplier := int64(GetRealmIndex(realm) + 1)
	return base * realmMultiplier * int64(realmLevel+1)
}

// CalculateBreakthroughCost tính linh thạch cần tiêu hao khi đột phá.
func CalculateBreakthroughCost(realmLevel int) int64 {
	return int64(100 * realmLevel)
}
