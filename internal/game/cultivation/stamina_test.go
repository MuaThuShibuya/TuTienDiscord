package cultivation

import (
	"testing"
)

// Giả lập cấu trúc StaminaRestoreResult theo thiết kế
type StaminaRestoreResult struct {
	Before   int64
	After    int64
	Max      int64
	Restored int64
	Capped   bool
}

func restoreLogic(current, max, amount int64) StaminaRestoreResult {
	res := StaminaRestoreResult{Before: current, Max: max, Restored: amount, Capped: false}
	if amount <= 0 {
		return res // Handle error out of band
	}
	res.After = current + amount
	if res.After >= max {
		res.Restored = max - current
		res.After = max
		res.Capped = true
	}
	return res
}

func TestRestoreStamina_IncreasesStamina(t *testing.T) {
	res := restoreLogic(20, 100, 30)
	if res.After != 50 || res.Restored != 30 || res.Capped != false {
		t.Errorf("Kỳ vọng 50, restored 30, capped false. Nhận: %+v", res)
	}
}

func TestRestoreStamina_ClampToMax(t *testing.T) {
	res := restoreLogic(90, 100, 30)
	if res.After != 100 || res.Restored != 10 || res.Capped != true {
		t.Errorf("Kỳ vọng clamp tại 100, restored 10, capped true. Nhận: %+v", res)
	}
}
