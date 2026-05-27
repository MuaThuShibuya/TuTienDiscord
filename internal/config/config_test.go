// File: internal/config/config_test.go
package config

import (
	"testing"
)

func TestConfig_IsOwner(t *testing.T) {
	cfg := &Config{
		Discord: DiscordConfig{
			OwnerIDs: []string{"123", "456"},
		},
	}

	if !cfg.IsOwner("123") {
		t.Errorf("Mong đợi 123 là Owner")
	}
	if cfg.IsOwner("789") {
		t.Errorf("Không mong đợi 789 là Owner")
	}

	var nilCfg *Config
	if nilCfg.IsOwner("123") {
		t.Errorf("Nil config phải trả về false")
	}
}

func TestConfig_CanExecuteDangerZone(t *testing.T) {
	// Case 1: Không có owner -> Cấm tuyệt đối
	cfgNoOwner := &Config{Discord: DiscordConfig{OwnerIDs: []string{}}}
	if cfgNoOwner.CanExecuteDangerZone() {
		t.Errorf("Rỗng OwnerID phải chặn Danger Zone")
	}

	// Case 2: Môi trường Production nhưng Allow=false -> Cấm
	cfgProdBlock := &Config{
		Discord: DiscordConfig{OwnerIDs: []string{"123"}},
		App:     AppConfig{Env: "production", AllowDangerousAdmin: false},
	}
	if cfgProdBlock.CanExecuteDangerZone() {
		t.Errorf("Production và Allow=false phải chặn Danger Zone")
	}

	// Case 3: Môi trường Production và Allow=true -> Cho phép
	cfgProdAllow := &Config{
		Discord: DiscordConfig{OwnerIDs: []string{"123"}},
		App:     AppConfig{Env: "production", AllowDangerousAdmin: true},
	}
	if !cfgProdAllow.CanExecuteDangerZone() {
		t.Errorf("Production và Allow=true phải cho phép Danger Zone")
	}
}
