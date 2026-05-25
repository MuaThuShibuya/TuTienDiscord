// File: internal/game/cultivation/service_test.go
// Phiên bản: v0.1.1
// Mục đích: Unit test cho cultivation.Service dùng in-memory repository.
// Ghi chú: Không kết nối MongoDB thật. Fake data chỉ dùng trong test, không dùng trong runtime.

package cultivation_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cultivation"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- In-memory repository ---

type memCultivationRepo struct {
	mu       sync.Mutex
	profiles map[string]*cultivation.CultivationProfile // key: userID+":"+guildID
}

func newMemCultivationRepo() *memCultivationRepo {
	return &memCultivationRepo{profiles: make(map[string]*cultivation.CultivationProfile)}
}

func key(userID, guildID string) string {
	return userID + ":" + guildID
}

func (r *memCultivationRepo) FindByUserID(_ context.Context, userID, guildID string) (*cultivation.CultivationProfile, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.profiles[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w: userId=%s", apperrors.ErrNotFound, userID)
	}
	copy := *p
	return &copy, nil
}

func (r *memCultivationRepo) Upsert(_ context.Context, profile *cultivation.CultivationProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	copy := *profile
	r.profiles[key(profile.UserID, profile.GuildID)] = &copy
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

func TestGetOrCreate_NewProfile(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo())
	ctx := context.Background()

	profile, err := svc.GetOrCreate(ctx, "user1", "guild1")
	if err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}
	if profile.UserID != "user1" {
		t.Errorf("UserID sai: muốn 'user1', có '%s'", profile.UserID)
	}
	// Cảnh giới khởi đầu phải là Luyện Khí
	if profile.Realm != cultivation.RealmQiRefining {
		t.Errorf("Realm khởi đầu sai: muốn %s, có %s", cultivation.RealmQiRefining, profile.Realm)
	}
	if profile.RealmLevel != 1 {
		t.Errorf("RealmLevel khởi đầu phải là 1, có %d", profile.RealmLevel)
	}
}

func TestGetOrCreate_ExistingProfile(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo())
	ctx := context.Background()

	first, _ := svc.GetOrCreate(ctx, "user2", "guild1")
	second, err := svc.GetOrCreate(ctx, "user2", "guild1")
	if err != nil {
		t.Fatalf("GetOrCreate lần 2 thất bại: %v", err)
	}
	// Phải là cùng cảnh giới
	if first.Realm != second.Realm {
		t.Error("Realm phải giống nhau cho cùng người chơi")
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo())
	ctx := context.Background()

	_, err := svc.GetProfile(ctx, "ghost", "guild1")
	if !apperrors.IsNotFound(err) {
		t.Errorf("Phải trả về ErrNotFound, có: %v", err)
	}
}

func TestGetProfile_Found(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user3", "guild1")
	p, err := svc.GetProfile(ctx, "user3", "guild1")
	if err != nil {
		t.Fatalf("GetProfile thất bại: %v", err)
	}
	if p.UserID != "user3" {
		t.Errorf("UserID sai: muốn 'user3', có '%s'", p.UserID)
	}
}

func TestDefaultValues(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo())
	ctx := context.Background()

	p, _ := svc.GetOrCreate(ctx, "user4", "guild1")

	// Tâm cảnh khởi đầu phải là Bình Tĩnh
	if p.MindState != cultivation.MindStateCalm {
		t.Errorf("MindState khởi đầu sai: muốn %s, có %s", cultivation.MindStateCalm, p.MindState)
	}
	// Thể lực phải > 0
	if p.MaxStamina <= 0 {
		t.Error("MaxStamina phải lớn hơn 0")
	}
	// Exp cần đột phá phải > 0
	if p.CultivationExpRequired <= 0 {
		t.Error("CultivationExpRequired phải lớn hơn 0")
	}
	// CanBreakthrough phải là false khi chưa có exp
	if p.CanBreakthrough() {
		t.Error("Người chơi mới không nên có thể đột phá ngay")
	}
}
