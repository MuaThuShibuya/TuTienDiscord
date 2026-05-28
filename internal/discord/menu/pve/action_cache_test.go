package pve

import "testing"

func TestActionCache_LengthUnder100(t *testing.T) {
	cache := NewActionCache()
	payload := PvEActionPayload{
		OwnerID:         "user123",
		MenuSessionID:   "menu_session_12345678901234567890",
		CombatSessionID: "ss_1234567890123456",
		TargetID:        "e_1",
	}
	token := cache.Save(payload)

	// Giả lập custom_id attack
	// pve:attack:<menu_session>:<token>
	customID := "pve:attack:menu_session_12345678901234567890:" + token
	if len(customID) > 100 {
		t.Errorf("Custom ID quá dài: %d ký tự", len(customID))
	}
	if len(customID) > 60 {
		t.Logf("Chiều dài an toàn lý tưởng: %d", len(customID))
	}
}
