// File: internal/game/profile/service_test.go
// Phiên bản: v0.1.1
// Mục đích: Unit test cho profile.Service dùng in-memory repository.
// Ghi chú: Không kết nối MongoDB thật. Fake data chỉ dùng trong test, không dùng trong runtime.

package profile_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	apperrors "github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/profile"
	"github.com/whiskey/tu-tien-bot/internal/logger"
)

// --- In-memory repository ---

type memProfileRepo struct {
	mu      sync.Mutex
	players map[string]*profile.Player // key: userID+":"+guildID
}

func newMemProfileRepo() *memProfileRepo {
	return &memProfileRepo{players: make(map[string]*profile.Player)}
}

func key(userID, guildID string) string {
	return userID + ":" + guildID
}

func (r *memProfileRepo) FindByUserID(_ context.Context, userID, guildID string) (*profile.Player, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.players[key(userID, guildID)]
	if !ok {
		return nil, fmt.Errorf("%w: userId=%s guildId=%s", apperrors.ErrNotFound, userID, guildID)
	}
	// Trả về bản copy để tránh data race trong test
	copy := *p
	return &copy, nil
}

func (r *memProfileRepo) Create(_ context.Context, player *profile.Player) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := key(player.UserID, player.GuildID)
	if _, ok := r.players[k]; ok {
		return fmt.Errorf("%w", apperrors.ErrAlreadyExists)
	}
	copy := *player
	copy.CreatedAt = time.Now().UTC()
	copy.LastActiveAt = time.Now().UTC()
	r.players[k] = &copy
	return nil
}

func (r *memProfileRepo) UpdateLastActive(_ context.Context, userID, guildID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.players[key(userID, guildID)]
	if !ok {
		return fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	p.LastActiveAt = time.Now().UTC()
	return nil
}

func (r *memProfileRepo) UpdateDaoName(_ context.Context, userID, guildID, daoName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.players[key(userID, guildID)]
	if !ok {
		return fmt.Errorf("%w", apperrors.ErrNotFound)
	}
	p.DaoName = daoName
	return nil
}

// --- Test setup ---

func TestMain(m *testing.M) {
	// Khởi tạo logger trước khi chạy test — services dùng logger.L()
	// Format json + level error: im lặng hoàn toàn khi test pass, chỉ in khi panic
	if err := logger.Init(logger.Options{Level: "error", Format: "json"}); err != nil {
		panic("logger init thất bại: " + err.Error())
	}
	m.Run()
}

// --- Tests ---

func TestGetOrCreate_NewPlayer(t *testing.T) {
	svc := profile.NewService(newMemProfileRepo())
	ctx := context.Background()

	player, err := svc.GetOrCreate(ctx, "user1", "guild1", "whisky", "Whisky")
	if err != nil {
		t.Fatalf("GetOrCreate thất bại: %v", err)
	}
	if player.UserID != "user1" {
		t.Errorf("UserID sai: muốn 'user1', có '%s'", player.UserID)
	}
	if player.GuildID != "guild1" {
		t.Errorf("GuildID sai")
	}
	// DaoName phải không rỗng
	if player.DaoName == "" {
		t.Error("DaoName không được rỗng")
	}
}

func TestGetOrCreate_ExistingPlayer(t *testing.T) {
	repo := newMemProfileRepo()
	svc := profile.NewService(repo)
	ctx := context.Background()

	// Tạo lần đầu
	first, _ := svc.GetOrCreate(ctx, "user2", "guild1", "user", "Player One")
	// Gọi lần hai phải trả về player cũ
	second, err := svc.GetOrCreate(ctx, "user2", "guild1", "user", "Player One")
	if err != nil {
		t.Fatalf("GetOrCreate lần 2 thất bại: %v", err)
	}
	if first.DaoName != second.DaoName {
		t.Errorf("DaoName phải giống nhau: '%s' vs '%s'", first.DaoName, second.DaoName)
	}
}

func TestGetPlayer_NotFound(t *testing.T) {
	svc := profile.NewService(newMemProfileRepo())
	ctx := context.Background()

	_, err := svc.GetPlayer(ctx, "ghost", "guild1")
	if !apperrors.IsNotFound(err) {
		t.Errorf("Phải trả về ErrNotFound, có: %v", err)
	}
}

func TestSetDaoName_Valid(t *testing.T) {
	repo := newMemProfileRepo()
	svc := profile.NewService(repo)
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user3", "guild1", "user", "Player")
	err := svc.SetDaoName(ctx, "user3", "guild1", "Kiếm Thần")
	if err != nil {
		t.Fatalf("SetDaoName thất bại: %v", err)
	}

	player, _ := svc.GetPlayer(ctx, "user3", "guild1")
	if player.DaoName != "Kiếm Thần" {
		t.Errorf("DaoName sai: muốn 'Kiếm Thần', có '%s'", player.DaoName)
	}
}

func TestSetDaoName_TooShort(t *testing.T) {
	svc := profile.NewService(newMemProfileRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user4", "guild1", "user", "Player")
	err := svc.SetDaoName(ctx, "user4", "guild1", "A") // 1 ký tự — quá ngắn
	if err == nil {
		t.Error("Phải trả về lỗi khi đạo hiệu quá ngắn")
	}
}

func TestSetDaoName_TooLong(t *testing.T) {
	svc := profile.NewService(newMemProfileRepo())
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user5", "guild1", "user", "Player")
	// 25 ký tự — quá dài (max 24)
	err := svc.SetDaoName(ctx, "user5", "guild1", "ABCDEFGHIJKLMNOPQRSTUVWXY")
	if err == nil {
		t.Error("Phải trả về lỗi khi đạo hiệu quá dài")
	}
}

func TestTouchLastActive(t *testing.T) {
	repo := newMemProfileRepo()
	svc := profile.NewService(repo)
	ctx := context.Background()

	svc.GetOrCreate(ctx, "user6", "guild1", "user", "Player")
	// TouchLastActive không trả về lỗi — chỉ log nếu thất bại
	svc.TouchLastActive(ctx, "user6", "guild1")
	// Không panic là pass
}
