// File: internal/game/aptitude/service_test.go
package aptitude_test

import (
	"context"
	"testing"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"github.com/whiskey/tu-tien-bot/internal/game/aptitude"
)

// --- Fake Repo ---
type fakeRepo struct {
	profiles    map[string]*aptitude.AptitudeProfile
	createCount int
}

func (r *fakeRepo) GetByUserID(ctx context.Context, userID string) (*aptitude.AptitudeProfile, error) {
	if p, ok := r.profiles[userID]; ok {
		return p, nil
	}
	return nil, apperrors.ErrNotFound
}
func (r *fakeRepo) Create(ctx context.Context, profile *aptitude.AptitudeProfile) error {
	r.createCount++
	r.profiles[profile.UserID] = profile
	return nil
}

// --- Tests ---
func TestRollForNewCharacter_CreateWhenMissing(t *testing.T) {
	repo := &fakeRepo{profiles: make(map[string]*aptitude.AptitudeProfile)}
	svc := aptitude.NewService(repo)

	profile, def, err := svc.RollForNewCharacter(context.Background(), "u1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if profile == nil || def == nil {
		t.Fatal("Profile hoặc Definition bị nil")
	}
	if repo.createCount != 1 {
		t.Errorf("Create phải được gọi đúng 1 lần, thực tế: %d", repo.createCount)
	}
}

func TestRollForNewCharacter_IdempotentWhenExists(t *testing.T) {
	repo := &fakeRepo{profiles: make(map[string]*aptitude.AptitudeProfile)}
	svc := aptitude.NewService(repo)

	p1, _, _ := svc.RollForNewCharacter(context.Background(), "u2")
	p2, _, _ := svc.RollForNewCharacter(context.Background(), "u2")

	if repo.createCount != 1 {
		t.Errorf("Chỉ được lưu DB 1 lần, nhưng Create được gọi %d lần", repo.createCount)
	}
	if p1.AptitudeID != p2.AptitudeID {
		t.Errorf("Idempotent lỗi: Lần 1 ra %s, lần 2 ra %s", p1.AptitudeID, p2.AptitudeID)
	}
}

func TestGetProfile_ReturnsDefinition(t *testing.T) {
	repo := &fakeRepo{profiles: map[string]*aptitude.AptitudeProfile{"u3": {UserID: "u3", AptitudeID: "apt_kiem_tam_so_khai"}}}
	svc := aptitude.NewService(repo)

	_, def, err := svc.GetProfile(context.Background(), "u3")
	if err != nil || def.ID != "apt_kiem_tam_so_khai" {
		t.Fatalf("Lỗi trả về definition: %v", err)
	}
}

func TestGetProfile_MissingDefinition_Fallback(t *testing.T) {
	repo := &fakeRepo{profiles: map[string]*aptitude.AptitudeProfile{"u4": {UserID: "u4", AptitudeID: "tu_chat_khong_ton_tai"}}}
	svc := aptitude.NewService(repo)

	_, def, _ := svc.GetProfile(context.Background(), "u4")
	if def == nil || def.ID != "apt_pham_tu" {
		t.Errorf("Phải fallback về phàm tư nếu ID không tồn tại, nhận: %v", def)
	}
}
