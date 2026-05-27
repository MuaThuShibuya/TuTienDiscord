// File: internal/game/combat/session_helpers.go
package combat

import "time"

func (s *CombatSession) IsActive() bool {
	return s.State == StateActive
}

func (s *CombatSession) IsExpired(now time.Time) bool {
	return now.After(s.ExpiresAt)
}

func (s *CombatSession) IsPlayerTurn() bool {
	return s.CurrentActorID == s.Player.ID
}

func (s *CombatSession) IsPlayerDead() bool {
	return s.Player.CurrentHP <= 0
}

func (s *CombatSession) FindEnemyIndex(enemyID string) int {
	for i, e := range s.Enemies {
		if e.ID == enemyID {
			return i
		}
	}
	return -1
}

func (s *CombatSession) AreAllEnemiesDead() bool {
	for _, e := range s.Enemies {
		if e.CurrentHP > 0 {
			return false
		}
	}
	return true
}

func (s *CombatSession) HasIdempotencyKey(key string) bool {
	for _, k := range s.IdempotencyKeys {
		if k == key {
			return true
		}
	}
	return false
}

func (s *CombatSession) AddIdempotencyKey(key string) {
	if !s.HasIdempotencyKey(key) {
		s.IdempotencyKeys = append(s.IdempotencyKeys, key)
	}
}

func (s *CombatSession) AppendLog(entry CombatLogEntry) {
	s.Logs = append(s.Logs, entry)
}

// TrimLogs giữ lại tối đa số lượng log gần nhất để tránh phình to Document trên MongoDB
func (s *CombatSession) TrimLogs(max int) {
	if len(s.Logs) > max {
		s.Logs = s.Logs[len(s.Logs)-max:]
	}
}
