package pvecombat

import (
	"context"
	"errors"

	"github.com/whiskey/tu-tien-bot/internal/game/combat"
)

type AutoBattleOptions struct {
	MaxActions     int
	IdempotencyKey string
	PreferSkill    bool
}

type AutoBattleResult struct {
	Session       *combat.CombatSession
	ActionsTaken  int
	StoppedReason string
}

// AutoBattle chạy vòng lặp tấn công tự động tối đa MaxActions lần.
func (s *Service) AutoBattle(ctx context.Context, userID, sessionID string, opts AutoBattleOptions) (*AutoBattleResult, error) {
	// TODO: Tính năng AutoBattle đang trong quá trình chuyển giao sang Combat Engine v0.5.
	// Vòng lặp ủy thác tự động gọi combatSvc sẽ được viết lại tại đây để tránh kẹt vòng lặp.
	return nil, errors.New("hàm AutoBattle đang trong quá trình chuyển giao sang Combat Engine v0.5")
}
