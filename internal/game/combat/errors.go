// File: internal/game/combat/errors.go
package combat

import "errors"

var (
	// ErrCombatSessionAlreadyActive trả về khi người chơi đang trong một trận chiến khác chưa kết thúc.
	ErrCombatSessionAlreadyActive = errors.New("đạo hữu đang trong một trận chiến khác chưa kết thúc")

	ErrAreaNotFound       = errors.New("không tìm thấy khu vực")
	ErrInvalidStage       = errors.New("ải không hợp lệ")
	ErrEncounterEmpty     = errors.New("không có kẻ địch nào xuất hiện, khu vực này đang yên bình")
	ErrInvalidCombatStats = errors.New("chỉ số chiến đấu không hợp lệ (máu hoặc tốc độ bất thường)")
	ErrEnemyLimitExceeded = errors.New("số lượng kẻ địch vượt quá giới hạn hệ thống cho phép")

	// Combat Loop Errors
	ErrCombatSessionNotFound  = errors.New("không tìm thấy trận đấu")
	ErrCombatSessionExpired   = errors.New("trận đấu đã hết hạn")
	ErrCombatSessionNotActive = errors.New("trận đấu không còn hoạt động")
	ErrCombatSessionForbidden = errors.New("đạo hữu không có quyền truy cập trận đấu này")
	ErrNotYourTurn            = errors.New("chưa tới lượt của đạo hữu")
	ErrTargetNotFound         = errors.New("không tìm thấy mục tiêu")
	ErrTargetAlreadyDead      = errors.New("mục tiêu đã bị tiêu diệt")

	// Reward Errors
	ErrRewardSessionNotWon         = errors.New("trận chiến chưa giành chiến thắng, không thể nhận thưởng")
	ErrRewardAlreadyClaimed        = errors.New("phần thưởng của trận này đã được nhận")
	ErrRewardInvalidIdempotencyKey = errors.New("idempotency key không hợp lệ khi nhận thưởng")
	ErrRewardGrantFailed           = errors.New("lỗi trao phần thưởng")
)
