// File: internal/discord/menu/admin/router_test.go
package admin

import (
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/config"
	"github.com/whiskey/tu-tien-bot/internal/discord/menu"
)

func TestValidateAdminAction(t *testing.T) {
	// Setup base configs
	cfgDev := &config.Config{
		App:     config.AppConfig{Env: "development", AllowDangerousAdmin: false},
		Discord: config.DiscordConfig{OwnerIDs: []string{"owner_123"}},
	}
	cfgProdBlocked := &config.Config{
		App:     config.AppConfig{Env: "production", AllowDangerousAdmin: false},
		Discord: config.DiscordConfig{OwnerIDs: []string{"owner_123"}},
	}
	cfgProdAllowed := &config.Config{
		App:     config.AppConfig{Env: "production", AllowDangerousAdmin: true},
		Discord: config.DiscordConfig{OwnerIDs: []string{"owner_123"}},
	}
	cfgNoOwner := &config.Config{
		App:     config.AppConfig{Env: "development", AllowDangerousAdmin: true},
		Discord: config.DiscordConfig{OwnerIDs: []string{}},
	}

	tests := []struct {
		name          string
		cfg           *config.Config
		userID        string
		action        string
		confirmPhrase string
		wantErr       bool
	}{
		{
			name:          "1. Non-owner bị chặn",
			cfg:           cfgDev,
			userID:        "normal_user",
			action:        menu.ActionAdminMain,
			confirmPhrase: "",
			wantErr:       true,
		},
		{
			name:          "2. Owner được mở panel thường",
			cfg:           cfgDev,
			userID:        "owner_123",
			action:        menu.ActionAdminMain,
			confirmPhrase: "",
			wantErr:       false,
		},
		{
			name:          "3. Owner sai confirm phrase bị chặn",
			cfg:           cfgDev,
			userID:        "owner_123",
			action:        menu.ActionAdminMigrateApply,
			confirmPhrase: "SAI MAT KHAU",
			wantErr:       true,
		},
		{
			name:          "4. Owner đúng confirm phrase + production + Allow=false bị chặn",
			cfg:           cfgProdBlocked,
			userID:        "owner_123",
			action:        menu.ActionAdminResetAllApply,
			confirmPhrase: "XACNHAN",
			wantErr:       true,
		},
		{
			name:          "5. Owner đúng confirm phrase + development được phép",
			cfg:           cfgDev,
			userID:        "owner_123",
			action:        menu.ActionAdminResetAllApply,
			confirmPhrase: "XACNHAN",
			wantErr:       false,
		},
		{
			name:          "6. Owner đúng confirm phrase + production + Allow=true được phép",
			cfg:           cfgProdAllowed,
			userID:        "owner_123",
			action:        menu.ActionAdminResetUserApply,
			confirmPhrase: "XACNHAN target123",
			wantErr:       false,
		},
		{
			name:          "7. OWNER_ID rỗng bị chặn",
			cfg:           cfgNoOwner,
			userID:        "owner_123",
			action:        menu.ActionAdminCombatCleanApply,
			confirmPhrase: "XACNHAN",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAdminAction(tt.cfg, tt.userID, tt.action, tt.confirmPhrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAdminAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
