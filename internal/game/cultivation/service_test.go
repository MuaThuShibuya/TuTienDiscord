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
	"time"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/cooldown"
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

func (r *memCultivationRepo) UpdateStats(_ context.Context, profile *cultivation.CultivationProfile) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.profiles[key(profile.UserID, profile.GuildID)]
	if !ok {
		return fmt.Errorf("%w: cultivation update failed", apperrors.ErrNotFound)
	}

	profile.UpdatedAt = time.Now().UTC()
	existing.Realm = profile.Realm
	existing.RealmLevel = profile.RealmLevel
	existing.CultivationExp = profile.CultivationExp
	existing.CultivationExpRequired = profile.CultivationExpRequired
	existing.CombatPower = profile.CombatPower
	existing.MindState = profile.MindState
	existing.Path = profile.Path
	existing.UpdatedAt = profile.UpdatedAt

	return nil
}

// --- Mocks ---

type mockCooldownSvc struct{}

func (m *mockCooldownSvc) IsOnCooldown(ctx context.Context, userID, guildID string, act cooldown.Action) (bool, time.Duration) {
	return false, 0
}
func (m *mockCooldownSvc) SetCooldown(ctx context.Context, userID, guildID string, act cooldown.Action, duration time.Duration) error {
	return nil
}
func (m *mockCooldownSvc) ClearCooldown(ctx context.Context, userID, guildID string, act cooldown.Action) error {
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
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
	ctx := context.Background()

	profile, err := svc.GetOrCreate(ctx, "user1", "guild1")
	if err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}
	if profile.UserID != "user1" {
		t.Errorf("UserID sai: muốn 'user1', có '%s'", profile.UserID)
	}
	// Cảnh giới khởi đầu phải là Phàm Nhân
	if profile.Realm != cultivation.DefaultRealm {
		t.Errorf("Realm khởi đầu sai: muốn %s, có %s", cultivation.DefaultRealm, profile.Realm)
	}
	if profile.RealmLevel != 1 {
		t.Errorf("RealmLevel khởi đầu phải là 1, có %d", profile.RealmLevel)
	}
}

func TestGetOrCreate_ExistingProfile(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
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
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
	ctx := context.Background()

	_, err := svc.GetProfile(ctx, "ghost", "guild1")
	if !apperrors.IsNotFound(err) {
		t.Errorf("Phải trả về ErrNotFound, có: %v", err)
	}
}

func TestGetProfile_Found(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
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
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
	ctx := context.Background()

	p, _ := svc.GetOrCreate(ctx, "user4", "guild1")

	// Tâm cảnh khởi đầu phải là 50 (Bình Tĩnh)
	if p.MindState != 50 {
		t.Errorf("MindState khởi đầu sai: muốn 50, có %d", p.MindState)
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

func TestChoosePath(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
	ctx := context.Background()

	// 1. Setup profile
	_, _ = svc.GetOrCreate(ctx, "user_path", "guild1")

	// 2. Chọn sai path -> ErrInvalidInput
	err := svc.ChoosePath(ctx, "user_path", "guild1", cultivation.CultivationPath("InvalidPath"))
	if !apperrors.IsInvalidInput(err) {
		t.Errorf("Phải trả về ErrInvalidInput, có: %v", err)
	}

	// 3. Chọn đúng path -> Thành công
	err = svc.ChoosePath(ctx, "user_path", "guild1", cultivation.PathSword)
	if err != nil {
		t.Fatalf("ChoosePath thất bại: %v", err)
	}

	// 4. Chọn lại lần 2 -> ErrPathAlreadyChosen
	err = svc.ChoosePath(ctx, "user_path", "guild1", cultivation.PathBody)
	if err == nil || err.Error() != apperrors.ErrPathAlreadyChosen.Error() {
		t.Errorf("Phải trả về ErrPathAlreadyChosen, có: %v", err)
	}
}

func TestCultivationActions_NoStaminaNeeded(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo(), &mockCooldownSvc{}, nil)
	ctx := context.Background()
	in := cultivation.CultivationActionInput{
		UserID:  "stamina_user",
		GuildID: "guild1",
		Now:     time.Now().UTC(),
	}

	// GetOrCreate sẽ tạo profile có Stamina=0 vì đã bỏ DefaultStamina
	_, _ = svc.GetOrCreate(ctx, in.UserID, in.GuildID)

	// 1. Tĩnh tu
	if _, err := svc.Meditate(ctx, in); err != nil {
		t.Errorf("Meditate lỗi khi stamina=0: %v", err)
	}

	// 2. Bế quan
	if _, err := svc.Seclusion(ctx, in); err != nil {
		t.Errorf("Seclusion lỗi khi stamina=0: %v", err)
	}

	// 3. Luyện thể
	if _, err := svc.BodyTraining(ctx, in); err != nil {
		t.Errorf("BodyTraining lỗi khi stamina=0: %v", err)
	}
}

func TestGetProfile_ReadOnly(t *testing.T) {
	svc := cultivation.NewService(newMemCultivationRepo(), nil, nil)
	ctx := context.Background()
	_, _ = svc.GetOrCreate(ctx, "user5", "guild1")
	_, err := svc.GetProfile(ctx, "user5", "guild1")
	if err != nil {
		t.Errorf("GetProfile lỗi: %v", err)
	}
	// Vì hàm GetProfile ở file gốc không gọi UpdateStats (đã xóa), test ngầm kiểm chứng không có runtime error.
}
