// File: internal/game/combat/auto_battle.go
// Chức năng: Nền tảng cấu hình và giao diện ra quyết định cho hệ thống Đánh Tự Động (Auto Battle).

package combat

import (
	"context"
	"fmt"
)

type AutoBattlePolicy struct {
	Enabled                  bool    `bson:"enabled" json:"enabled"`
	PreferSkill              bool    `bson:"preferSkill" json:"preferSkill"`
	UseUltimateWhenAvailable bool    `bson:"useUltimateWhenAvailable" json:"useUltimateWhenAvailable"`
	StopWhenLowHPPercent     float64 `bson:"stopWhenLowHpPercent" json:"stopWhenLowHpPercent"`
}

// AutoBattleDecider định nghĩa hợp đồng ra quyết định hành động tự động.
// TODO: v0.4.x - Implement decider sử dụng Skill Service và Turn Order.
type AutoBattleDecider interface {
	// DecideAction nhận vào Session và Actor đang tới lượt, trả về ID Kỹ Năng và Mục Tiêu.
	DecideAction(session *CombatSession, actor *CombatActor) (skillID string, targetID string)
}

// DefaultAutoDecider là AI cốt lõi cho chế độ tự động chiến đấu (Auto).
// Được thiết kế dạng module (Engine-level) để dùng chung cho cả PvE, PvP và mọi trận chiến.
type DefaultAutoDecider struct {
	Policy AutoBattlePolicy
}

// NewDefaultAutoDecider khởi tạo AI cơ bản kèm theo chính sách chiến thuật (Policy).
func NewDefaultAutoDecider(policy AutoBattlePolicy) AutoBattleDecider {
	return &DefaultAutoDecider{Policy: policy}
}

// DecideAction quét toàn bộ chiến trường để chọn mục tiêu và chiêu thức tối ưu.
func (d *DefaultAutoDecider) DecideAction(session *CombatSession, actor *CombatActor) (skillID string, targetID string) {
	isPlayer := actor.ID == session.Player.ID

	var validTargets []*CombatActor
	// Phân loại phe: Nếu actor là phe ta -> địch là đối thủ (quái/người). Ngược lại.
	if isPlayer {
		for i := range session.Enemies {
			if session.Enemies[i].CurrentHP > 0 {
				validTargets = append(validTargets, &session.Enemies[i])
			}
		}
	} else {
		if session.Player.CurrentHP > 0 {
			validTargets = append(validTargets, &session.Player)
		}
	}

	if len(validTargets) == 0 {
		return "", "" // Trận chiến đã kết thúc hoặc không còn mục tiêu
	}

	// Chiến thuật mặc định: Tập trung tiêu diệt kẻ địch có sinh lực (HP) thấp nhất
	bestTarget := validTargets[0]
	for _, t := range validTargets[1:] {
		if t.CurrentHP < bestTarget.CurrentHP {
			bestTarget = t
		}
	}

	// Trả về rỗng cho skillID (sử dụng Basic Attack).
	// Sẽ tích hợp thêm Policy PreferSkill sau khi hệ thống Kỹ Năng hoàn thiện.
	return "", bestTarget.ID
}

type AutoBattleOptions struct {
	MaxActions     int
	IdempotencyKey string
	PreferSkill    bool
}

type AutoBattleResult struct {
	Session       *CombatSession
	ActionsTaken  int
	StoppedReason string
}

// ExecuteAutoBattle chạy vòng lặp chiến đấu tự động ở tầng Engine.
// Dùng chung được cho PvE, PvP, Tông chiến, Boss...
func (s *Service) ExecuteAutoBattle(ctx context.Context, userID, sessionID string, opts AutoBattleOptions) (*AutoBattleResult, error) {
	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.UserID != userID {
		return nil, ErrCombatSessionForbidden
	}
	if !session.IsActive() {
		return nil, ErrCombatSessionNotActive
	}

	decider := NewDefaultAutoDecider(AutoBattlePolicy{PreferSkill: opts.PreferSkill})
	actionsTaken := 0
	stoppedReason := "max_actions"

	for actionsTaken < opts.MaxActions {
		if !session.IsActive() {
			stoppedReason = string(session.State)
			break
		}
		if !session.IsPlayerTurn() {
			stoppedReason = "not_player_turn"
			break
		}

		skillID, targetID := decider.DecideAction(session, &session.Player)
		if targetID == "" {
			stoppedReason = "no_target"
			break
		}

		// Tạo IdempotencyKey riêng biệt cho mỗi nhịp đánh để tránh bị chặn
		loopKey := fmt.Sprintf("%s_auto_%d", opts.IdempotencyKey, session.Turn)

		if skillID == "" {
			session, err = s.PlayerBasicAttack(ctx, userID, session.ID, targetID, loopKey)
			if err != nil {
				return nil, fmt.Errorf("lỗi đánh tự động nhịp %d: %w", actionsTaken+1, err)
			}
		} else {
			// TODO (v0.5+): Gọi s.PlayerSkillAttack khi hệ thống chiêu thức ra mắt
			stoppedReason = "skill_not_implemented"
			break
		}
		actionsTaken++
	}

	return &AutoBattleResult{
		Session:       session,
		ActionsTaken:  actionsTaken,
		StoppedReason: stoppedReason,
	}, nil
}
