// File: internal/game/cooldown/service_test.go
// Phiên bản: v0.1.1
// Mục đích: Unit test cho cooldown.Service dùng in-memory repository.
//           Test bao gồm thiết lập, kiểm tra, và xóa cooldown.
// Ghi chú: Không kết nối MongoDB thật. Fake data chỉ dùng trong test, không dùng trong runtime.

package cooldown_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cooldown"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- In-memory repository ---

type memCooldownRepo struct {
	mu        sync.Mutex
	cooldowns map[string]*cooldown.Cooldown // key: userID+":"+guildID+":"+action
}

func newMemCooldownRepo() *memCooldownRepo {
	return &memCooldownRepo{cooldowns: make(map[string]*cooldown.Cooldown)}
}

func key(userID, guildID string, action cooldown.Action) string {
	return fmt.Sprintf("%s:%s:%s", userID, guildID, string(action))
}

func (r *memCooldownRepo) Get(_ context.Context, userID, guildID string, action cooldown.Action) (*cooldown.Cooldown, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	cd, ok := r.cooldowns[key(userID, guildID, action)]
	if !ok {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	// Nếu đã hết hạn, coi như không có
	if cd.IsExpired() {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	copy := *cd
	return &copy, nil
}

func (r *memCooldownRepo) Set(_ context.Context, userID, guildID string, action cooldown.Action, duration time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cooldowns[key(userID, guildID, action)] = &cooldown.Cooldown{
		UserID:    userID,
		GuildID:   guildID,
		Action:    action,
		ExpiresAt: time.Now().UTC().Add(duration),
		CreatedAt: time.Now().UTC(),
	}
	return nil
}

func (r *memCooldownRepo) Delete(_ context.Context, userID, guildID string, action cooldown.Action) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cooldowns, key(userID, guildID, action))
	return nil
}

// --- Test setup ---

func TestMain(m *testing.M) {
	// Format json + level error: im lặng khi test pass
	if err := logger.Init(logger.Options{Level: "error", Format: "json"}); err != nil {
		panic("logger init thất bại: " + err.Error())
	}
	m.Run()
}

// --- Tests ---

func TestIsOnCooldown_None(t *testing.T) {
	svc := cooldown.NewService(newMemCooldownRepo())
	ctx := context.Background()

	onCooldown, remaining := svc.IsOnCooldown(ctx, "user1", "guild1", cooldown.ActionCultivate)
	if onCooldown {
		t.Error("Người dùng mới không nên có cooldown")
	}
	if remaining != 0 {
		t.Errorf("Remaining phải là 0, có %v", remaining)
	}
}

func TestSetAndCheckCooldown(t *testing.T) {
	svc := cooldown.NewService(newMemCooldownRepo())
	ctx := context.Background()

	// Thiết lập cooldown 5 phút
	err := svc.SetCooldown(ctx, "user2", "guild1", cooldown.ActionCultivate, 5*time.Minute)
	if err != nil {
		t.Fatalf("SetCooldown thất bại: %v", err)
	}

	// Phải đang cooldown
	onCooldown, remaining := svc.IsOnCooldown(ctx, "user2", "guild1", cooldown.ActionCultivate)
	if !onCooldown {
		t.Error("Phải đang trong cooldown sau khi SetCooldown")
	}
	if remaining <= 0 || remaining > 5*time.Minute {
		t.Errorf("Remaining không hợp lệ: %v", remaining)
	}
}

func TestClearCooldown(t *testing.T) {
	svc := cooldown.NewService(newMemCooldownRepo())
	ctx := context.Background()

	svc.SetCooldown(ctx, "user3", "guild1", cooldown.ActionDaily, 24*time.Hour)
	err := svc.ClearCooldown(ctx, "user3", "guild1", cooldown.ActionDaily)
	if err != nil {
		t.Fatalf("ClearCooldown thất bại: %v", err)
	}

	// Sau khi clear, không còn cooldown
	onCooldown, _ := svc.IsOnCooldown(ctx, "user3", "guild1", cooldown.ActionDaily)
	if onCooldown {
		t.Error("Cooldown phải được xóa sau ClearCooldown")
	}
}

func TestCooldownExpired(t *testing.T) {
	ctx := context.Background()

	// Dùng duration âm để giả lập cooldown đã hết hạn
	repo := newMemCooldownRepo()
	// Thêm trực tiếp vào repo một cooldown đã hết hạn
	_ = repo.Set(ctx, "user4", "guild1", cooldown.ActionGacha, -1*time.Second)

	svc2 := cooldown.NewService(repo)
	onCooldown, _ := svc2.IsOnCooldown(ctx, "user4", "guild1", cooldown.ActionGacha)
	if onCooldown {
		t.Error("Cooldown đã hết hạn không nên được tính là đang cooldown")
	}
}

func TestCooldownIsolatedByAction(t *testing.T) {
	svc := cooldown.NewService(newMemCooldownRepo())
	ctx := context.Background()

	// Chỉ đặt cooldown cho ActionCultivate
	svc.SetCooldown(ctx, "user5", "guild1", cooldown.ActionCultivate, 10*time.Minute)

	// ActionDaily phải không bị ảnh hưởng
	onCooldown, _ := svc.IsOnCooldown(ctx, "user5", "guild1", cooldown.ActionDaily)
	if onCooldown {
		t.Error("ActionDaily không nên bị cooldown khi chỉ đặt ActionCultivate")
	}
}

func TestCooldownIsolatedByGuild(t *testing.T) {
	svc := cooldown.NewService(newMemCooldownRepo())
	ctx := context.Background()

	// Cooldown trong guild1
	svc.SetCooldown(ctx, "user6", "guild1", cooldown.ActionPvP, 10*time.Minute)

	// Guild2 không bị ảnh hưởng
	onCooldown, _ := svc.IsOnCooldown(ctx, "user6", "guild2", cooldown.ActionPvP)
	if onCooldown {
		t.Error("Cooldown trong guild1 không nên ảnh hưởng đến guild2")
	}
}
