package menu_test

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
)

func TestParse_AdminModalSubmit(t *testing.T) {
	customID := "admin:reset_all_apply:41e29830eeb3fedc5c31d00e16fb0088"
	parsed, err := menu.Parse(customID)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if parsed.Domain != "admin" {
		t.Errorf("Expected domain 'admin', got %q", parsed.Domain)
	}
	if parsed.Action != "reset_all_apply" {
		t.Errorf("Expected action 'reset_all_apply', got %q", parsed.Action)
	}
	if parsed.SessionID != "41e29830eeb3fedc5c31d00e16fb0088" {
		t.Errorf("Expected session ID '41e29830eeb3fedc5c31d00e16fb0088', got %q", parsed.SessionID)
	}
}

func TestParse_MalformedID(t *testing.T) {
	_, err := menu.Parse("admin:invalid")
	if err == nil {
		t.Errorf("Expected error for malformed customID, got nil")
	}
}
