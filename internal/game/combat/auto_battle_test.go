package combat

import (
	"testing"
)

func TestDefaultAutoDecider_DecideAction(t *testing.T) {
	decider := NewDefaultAutoDecider(AutoBattlePolicy{})

	session := &CombatSession{
		Player: CombatActor{
			ID:        "player_1",
			CurrentHP: 100,
		},
		Enemies: []CombatActor{
			{ID: "enemy_full", CurrentHP: 100},
			{ID: "enemy_low", CurrentHP: 20}, // Quái thấp máu nhất -> Phải bị target
			{ID: "enemy_dead", CurrentHP: 0}, // Đã chết -> Phải bị bỏ qua
		},
	}

	t.Run("Người chơi tự động đánh mục tiêu thấp máu nhất", func(t *testing.T) {
		skillID, targetID := decider.DecideAction(session, &session.Player)

		if skillID != "" {
			t.Errorf("Mặc định phải dùng đánh thường (skillID rỗng), nhận: %s", skillID)
		}
		if targetID != "enemy_low" {
			t.Errorf("Phải target mục tiêu thấp máu nhất (enemy_low), nhận: %s", targetID)
		}
	})

	t.Run("Quái vật tự động target người chơi", func(t *testing.T) {
		skillID, targetID := decider.DecideAction(session, &session.Enemies[0])

		if skillID != "" {
			t.Errorf("Quái mặc định đánh thường, nhận: %s", skillID)
		}
		if targetID != "player_1" {
			t.Errorf("Quái phải target người chơi, nhận: %s", targetID)
		}
	})

	t.Run("Trả về rỗng khi không còn mục tiêu hợp lệ", func(t *testing.T) {
		// Mô phỏng tất cả quái đều đã chết
		session.Enemies[0].CurrentHP = 0
		session.Enemies[1].CurrentHP = 0

		skillID, targetID := decider.DecideAction(session, &session.Player)
		if targetID != "" {
			t.Errorf("Phải trả về rỗng khi hết quái, nhận targetID: %s", targetID)
		}
		if skillID != "" {
			t.Errorf("Phải trả về rỗng khi hết quái, nhận skillID: %s", skillID)
		}
	})
}
