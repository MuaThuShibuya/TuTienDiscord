// File: internal/game/combat/turn_order.go
// Chức năng: Quản lý thứ tự lượt đi dựa trên Tốc Độ (Speed) và Thanh Hành Động (Action Value).

package combat

import "sort"

type TurnOrderEntry struct {
	ActorID     string  `json:"actorId"`
	ActionValue float64 `json:"actionValue"` // Càng thấp càng tới lượt nhanh
	Speed       int64   `json:"speed"`
}

type TurnOrderService struct{}

func NewTurnOrderService() *TurnOrderService {
	return &TurnOrderService{}
}

// BuildInitialOrder khởi tạo thứ tự đánh ban đầu.
// Công thức cơ bản: ActionValue = 10000 / Speed
func (s *TurnOrderService) BuildInitialOrder(actors []CombatActor) []TurnOrderEntry {
	var order []TurnOrderEntry
	for _, a := range actors {
		sp := a.Stats.Speed
		if sp <= 0 {
			sp = 100 // Tốc độ mặc định nếu chưa cấu hình
		}
		order = append(order, TurnOrderEntry{
			ActorID:     a.ID,
			Speed:       sp,
			ActionValue: 10000.0 / float64(sp),
		})
	}
	s.sortOrder(order)
	return order
}

// NextActor lấy actor có ActionValue thấp nhất (đầu danh sách),
// sau đó tịnh tiến thời gian của tất cả mọi người và đẩy actor vừa đánh về sau.
func (s *TurnOrderService) NextActor(order []TurnOrderEntry) (TurnOrderEntry, []TurnOrderEntry) {
	if len(order) == 0 {
		return TurnOrderEntry{}, order
	}
	acting := order[0]
	passedAV := acting.ActionValue

	for i := range order {
		order[i].ActionValue -= passedAV
		if order[i].ActionValue < 0 {
			order[i].ActionValue = 0
		}
	}

	baseAV := 10000.0 / float64(acting.Speed)
	order[0].ActionValue += baseAV
	s.sortOrder(order)
	return acting, order
}

// AdvanceActor kéo thanh hành động (Turn Advance). amount = 0.3 nghĩa là giảm 30% AV gốc.
func (s *TurnOrderService) AdvanceActor(order []TurnOrderEntry, actorID string, amount float64) []TurnOrderEntry {
	for i, entry := range order {
		if entry.ActorID == actorID {
			baseAV := 10000.0 / float64(entry.Speed)
			order[i].ActionValue -= baseAV * amount
			if order[i].ActionValue < 0 {
				order[i].ActionValue = 0
			}
		}
	}
	s.sortOrder(order)
	return order
}

// DelayActor làm chậm thanh hành động (Turn Delay). amount = 0.3 nghĩa là tăng 30% AV gốc.
func (s *TurnOrderService) DelayActor(order []TurnOrderEntry, actorID string, amount float64) []TurnOrderEntry {
	for i, entry := range order {
		if entry.ActorID == actorID {
			baseAV := 10000.0 / float64(entry.Speed)
			order[i].ActionValue += baseAV * amount
		}
	}
	s.sortOrder(order)
	return order
}

// RemoveActor loại bỏ một thực thể khỏi TurnOrder (ví dụ khi đã chết).
func (s *TurnOrderService) RemoveActor(order []TurnOrderEntry, actorID string) []TurnOrderEntry {
	var newOrder []TurnOrderEntry
	for _, entry := range order {
		if entry.ActorID != actorID {
			newOrder = append(newOrder, entry)
		}
	}
	return newOrder
}

func (s *TurnOrderService) sortOrder(order []TurnOrderEntry) {
	sort.Slice(order, func(i, j int) bool {
		return order[i].ActionValue < order[j].ActionValue
	})
}
