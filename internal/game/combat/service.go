// File: internal/game/combat/service.go
// Chức năng: Khởi tạo và điều phối logic trận đấu, tuân thủ nguyên tắc không phụ thuộc ngược.

package combat

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"go.uber.org/zap"
)

type Service struct {
	repo      Repository
	turnOrder *TurnOrderService
	calc      Calculator
	rng       *rand.Rand
	now       func() time.Time
	log       *zap.Logger
}

func NewService(
	repo Repository,
	turnOrder *TurnOrderService,
	rng *rand.Rand,
	log *zap.Logger,
) (*Service, error) {
	if repo == nil {
		return nil, errors.New("combat repo cannot be nil")
	}
	if turnOrder == nil {
		return nil, errors.New("turnOrder cannot be nil")
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	if log == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Service{
		repo:      repo,
		turnOrder: turnOrder,
		calc:      NewDamageCalculator(),
		rng:       rng,
		now:       time.Now,
		log:       log.Named("game.combat"),
	}, nil
}

// PlayerBasicAttack xử lý lệnh đánh thường của người chơi.
func (s *Service) PlayerBasicAttack(ctx context.Context, userID, sessionID, targetID, idempotencyKey string) (*CombatSession, error) {
	if userID == "" || sessionID == "" || targetID == "" || idempotencyKey == "" {
		return nil, errors.New("thiếu tham số bắt buộc")
	}

	session, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, ErrCombatSessionNotFound
		}
		return nil, err
	}

	if session.UserID != userID {
		return nil, ErrCombatSessionForbidden
	}

	// Chống Double Click ưu tiên xử lý trước để có thể mượt mà trả về StateWon/Lost nếu click đúp lúc quái vừa chết
	if session.HasIdempotencyKey(idempotencyKey) {
		return session, nil
	}

	if !session.IsActive() {
		return nil, ErrCombatSessionNotActive
	}
	if session.IsExpired(s.now().UTC()) {
		return nil, ErrCombatSessionExpired
	}

	if !session.IsPlayerTurn() {
		return nil, ErrNotYourTurn
	}

	enemyIdx := session.FindEnemyIndex(targetID)
	if enemyIdx == -1 {
		return nil, ErrTargetNotFound
	}
	if session.Enemies[enemyIdx].CurrentHP <= 0 {
		return nil, ErrTargetAlreadyDead
	}

	// Calculate & Apply Damage
	dmgRes := s.calc.CalculateBasicAttack(&session.Player, &session.Enemies[enemyIdx], s.rng)
	session.Enemies[enemyIdx].CurrentHP -= dmgRes.Damage
	if session.Enemies[enemyIdx].CurrentHP < 0 {
		session.Enemies[enemyIdx].CurrentHP = 0
		// Loại bỏ ngay khỏi hàng đợi để không bị trôi Turn
		session.TurnOrder = s.turnOrder.RemoveActor(session.TurnOrder, targetID)
	}

	// Tích Nộ
	session.Player.CurrentRage += 10
	if session.Player.CurrentRage > 100 {
		session.Player.CurrentRage = 100
	}

	// Log
	msg := fmt.Sprintf("Đạo hữu tung quyền cước vào %s, gây %d sát thương.", session.Enemies[enemyIdx].Name, dmgRes.Damage)
	if dmgRes.IsCrit {
		msg = fmt.Sprintf("Chí mạng! Đạo hữu giáng đòn sấm sét vào %s, gây %d sát thương.", session.Enemies[enemyIdx].Name, dmgRes.Damage)
	}
	if session.Enemies[enemyIdx].CurrentHP == 0 {
		msg += " Mục tiêu đã bị hạ gục!"
	}
	session.AppendLog(CombatLogEntry{Turn: session.Turn, ActorID: session.Player.ID, Action: "attack", Message: msg, Damage: dmgRes.Damage, IsCrit: dmgRes.IsCrit, CreatedAt: s.now().UTC()})

	s.log.Info("Player Basic Attack",
		zap.String("userId", userID),
		zap.String("sessionId", sessionID),
		zap.String("targetId", targetID),
		zap.Int64("damage", dmgRes.Damage),
		zap.Bool("isCrit", dmgRes.IsCrit),
		zap.Int64("targetHpAfter", session.Enemies[enemyIdx].CurrentHP),
		zap.String("idempotencyKey", idempotencyKey),
	)

	// Check Win
	if session.AreAllEnemiesDead() {
		session.State = StateWon
		session.AppendLog(CombatLogEntry{Turn: session.Turn, ActorID: "system", Action: "win", Message: "Toàn bộ kẻ địch đã bị tiêu diệt. Đạo hữu giành chiến thắng!", CreatedAt: s.now().UTC()})
	} else {
		// Advance Turn & Auto-resolve Enemy Turns
		s.advanceTurn(session)
		s.processEnemyTurns(session)
	}

	session.AddIdempotencyKey(idempotencyKey)
	session.TrimLogs(20) // Giữ 20 dòng log gần nhất
	// Chống phình mảng IdempotencyKeys trong Document MongoDB
	if len(session.IdempotencyKeys) > 10 {
		session.IdempotencyKeys = session.IdempotencyKeys[len(session.IdempotencyKeys)-10:]
	}

	if err := s.repo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *Service) processEnemyTurns(session *CombatSession) {
	for session.IsActive() && !session.IsPlayerTurn() {
		enemyIdx := session.FindEnemyIndex(session.CurrentActorID)
		if enemyIdx == -1 || session.Enemies[enemyIdx].CurrentHP <= 0 {
			// Nếu chết do DoT, xóa khỏi thanh hành động thay vì push AV làm sai lệch lượt
			session.TurnOrder = s.turnOrder.RemoveActor(session.TurnOrder, session.CurrentActorID)
			if len(session.TurnOrder) > 0 {
				session.CurrentActorID = session.TurnOrder[0].ActorID
			}
			continue
		}

		enemy := &session.Enemies[enemyIdx]

		// Tạm thời quái luôn nhắm vào Player
		dmgRes := s.calc.CalculateBasicAttack(enemy, &session.Player, s.rng)
		session.Player.CurrentHP -= dmgRes.Damage
		if session.Player.CurrentHP < 0 {
			session.Player.CurrentHP = 0
		}

		msg := fmt.Sprintf("%s tấn công đạo hữu, gây %d sát thương.", enemy.Name, dmgRes.Damage)
		if dmgRes.IsCrit {
			msg = fmt.Sprintf("Nguy hiểm! %s tung đòn chí mạng gây %d sát thương.", enemy.Name, dmgRes.Damage)
		}
		session.AppendLog(CombatLogEntry{Turn: session.Turn, ActorID: enemy.ID, Action: "attack", Message: msg, Damage: dmgRes.Damage, IsCrit: dmgRes.IsCrit, CreatedAt: s.now().UTC()})

		// Check Lose
		if session.IsPlayerDead() {
			session.State = StateLost
			session.AppendLog(CombatLogEntry{Turn: session.Turn, ActorID: "system", Action: "lose", Message: "Đạo hữu đã trọng thương, không thể chiến đấu tiếp. Thất bại!", CreatedAt: s.now().UTC()})
			break
		}

		s.advanceTurn(session)
	}
}

func (s *Service) advanceTurn(session *CombatSession) {
	session.Turn++ // Hành động thực tế

	// Tính ActionValue tịnh tiến
	_, newOrder := s.turnOrder.NextActor(session.TurnOrder)
	session.TurnOrder = newOrder

	if len(newOrder) > 0 {
		session.CurrentActorID = newOrder[0].ActorID
	}
}
