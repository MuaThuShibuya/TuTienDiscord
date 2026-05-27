// File: internal/discord/menu/admin/router_test.go
package admin

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
)

func TestValidateAdminAction(t *testing.T) {
	cfg := &config.Config{
		Discord: config.DiscordConfig{OwnerIDs: []string{"owner_123"}},
	}

	tests := []struct {
		name          string
		userID        string
		action        string
		confirmPhrase string
		wantErr       bool
	}{
		{
			name:          "Chặn người chơi thường",
			userID:        "normal_user",
			action:        menu.ActionAdminMain,
			confirmPhrase: "",
			wantErr:       true,
		},
		{
			name:          "Cho phép Owner vào Menu",
			userID:        "owner_123",
			action:        menu.ActionAdminMain,
			confirmPhrase: "",
			wantErr:       false,
		},
		{
			name:          "Chặn thao tác nguy hiểm nếu sai mật khẩu xác nhận",
			userID:        "owner_123",
			action:        menu.ActionAdminMigrateApply,
			confirmPhrase: "SAI MAT KHAU",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAdminAction(cfg, tt.userID, tt.action, tt.confirmPhrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAdminAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
