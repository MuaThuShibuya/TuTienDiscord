// File: internal/game/combat/turn_order_test.go

package combat

import (
	"testing"
)

func TestTurnOrderService_BuildInitialOrder(t *testing.T) {
	svc := NewTurnOrderService()
	actors := []CombatActor{
		{ID: "slow_turtle", Stats: CombatStats{Speed: 50}},
		{ID: "fast_rabbit", Stats: CombatStats{Speed: 150}},
		{ID: "normal_human", Stats: CombatStats{Speed: 100}},
	}

	order := svc.BuildInitialOrder(actors)

	// Thỏ tốc 150 phải đi trước, rùa tốc 50 đi cuối
	if order[0].ActorID != "fast_rabbit" {
		t.Errorf("Kẻ nhanh nhất phải đi đầu tiên")
	}
	if order[2].ActorID != "slow_turtle" {
		t.Errorf("Kẻ chậm nhất phải đi cuối cùng")
	}
}

func TestTurnOrderService_AdvanceActor(t *testing.T) {
	svc := NewTurnOrderService()
	actors := []CombatActor{
		{ID: "p1", Stats: CombatStats{Speed: 100}}, // AV: 100
		{ID: "p2", Stats: CombatStats{Speed: 120}}, // AV: 83.33 (Đi trước)
	}

	order := svc.BuildInitialOrder(actors)
	// P1 dùng skill kéo 30% lượt (Turn Advance) -> Giảm 30 AV -> Còn 70 AV
	order = svc.AdvanceActor(order, "p1", 0.3)

	if order[0].ActorID != "p1" {
		t.Errorf("Sau khi được kéo lượt, P1 phải vượt lên trước P2")
	}
}

func TestTurnOrderService_DelayActor(t *testing.T) {
	svc := NewTurnOrderService()
	actors := []CombatActor{
		{ID: "p1", Stats: CombatStats{Speed: 120}}, // AV ban đầu thấp hơn
		{ID: "p2", Stats: CombatStats{Speed: 100}},
	}

	order := svc.BuildInitialOrder(actors)
	// P1 bị delay 50% lượt
	order = svc.DelayActor(order, "p1", 0.5)

	if order[0].ActorID != "p2" {
		t.Errorf("Sau khi P1 bị delay, P2 phải được đi trước")
	}
}

func TestTurnOrderService_NextActor(t *testing.T) {
	svc := NewTurnOrderService()
	actors := []CombatActor{
		{ID: "fast", Stats: CombatStats{Speed: 200}}, // AV = 50
		{ID: "slow", Stats: CombatStats{Speed: 100}}, // AV = 100
	}

	order := svc.BuildInitialOrder(actors)

	// fast đánh lần 1
	acting, order := svc.NextActor(order)
	if acting.ActorID != "fast" {
		t.Errorf("Kẻ nhanh phải đánh trước")
	}

	acting2, _ := svc.NextActor(order)
	if acting2.ActorID == "" {
		t.Errorf("NextActor không được trả về rỗng")
	}
}
