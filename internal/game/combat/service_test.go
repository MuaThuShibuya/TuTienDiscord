package combat

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
	"go.uber.org/zap"
)

// --- Fake Repo cho Combat Engine ---
type fakeCombatRepo struct {
	sessions map[string]*CombatSession
}

func newFakeCombatRepo() *fakeCombatRepo {
	return &fakeCombatRepo{sessions: make(map[string]*CombatSession)}
}
func (r *fakeCombatRepo) CreateSession(ctx context.Context, session *CombatSession) error {
	r.sessions[session.ID] = session
	return nil
}
func (r *fakeCombatRepo) GetSession(ctx context.Context, sessionID string) (*CombatSession, error) {
	s, ok := r.sessions[sessionID]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return s, nil
}
func (r *fakeCombatRepo) GetActiveSessionByUser(ctx context.Context, userID string) (*CombatSession, error) {
	return nil, apperrors.ErrNotFound
}
func (r *fakeCombatRepo) UpdateSession(ctx context.Context, session *CombatSession) error {
	r.sessions[session.ID] = session
	return nil
}
func (r *fakeCombatRepo) MarkSessionState(ctx context.Context, sessionID string, state SessionState) error {
	return nil
}
func (r *fakeCombatRepo) TryStartRewardClaim(ctx context.Context, sessionID, claimID string, now time.Time) (*CombatSession, error) {
	return nil, nil
}
func (r *fakeCombatRepo) CompleteRewardClaim(ctx context.Context, sessionID, claimID string, details []ClaimedReward, now time.Time) error {
	return nil
}
func (r *fakeCombatRepo) FailRewardClaim(ctx context.Context, sessionID, claimID, reason string, now time.Time) error {
	return nil
}

// --- Tests ---

func TestService_PlayerBasicAttack(t *testing.T) {
	repo := newFakeCombatRepo()
	svc, _ := NewService(repo, NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())

	// Hàm tiện ích để reset một trận đấu mới tinh trước mỗi bài test nhỏ
	resetSession := func() *CombatSession {
		session := &CombatSession{
			ID:             "ss_1",
			UserID:         "u1",
			State:          StateActive,
			Turn:           1,
			Player:         CombatActor{ID: "u1", Name: "Player", Stats: CombatStats{ATK: 100, Speed: 100}},
			Enemies:        []CombatActor{{ID: "e1", Name: "Slime", CurrentHP: 50, Stats: CombatStats{DEF: 10, Speed: 90}}},
			CurrentActorID: "u1", // Tới lượt player
			ExpiresAt:      time.Now().Add(time.Hour),
		}
		repo.sessions[session.ID] = session
		return session
	}

	t.Run("Đánh thường thành công & Lưu Idempotency Key", func(t *testing.T) {
		resetSession()
		res, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "btn_click_1")
		if err != nil {
			t.Fatalf("Không mong đợi lỗi: %v", err)
		}
		if res.Enemies[0].CurrentHP >= 50 {
			t.Errorf("Quái chưa bị trừ máu")
		}
		if !res.HasIdempotencyKey("btn_click_1") {
			t.Errorf("Phải lưu IdempotencyKey để chống double-click")
		}
	})

	t.Run("Chống Double Click trả về ngay lập tức", func(t *testing.T) {
		session := resetSession()
		session.AddIdempotencyKey("btn_click_1") // Gài sẵn key

		res, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "btn_click_1")
		if err != nil {
			t.Fatalf("Không mong đợi lỗi khi click đúp: %v", err)
		}
		if res.Enemies[0].CurrentHP != 50 {
			t.Errorf("Máu quái bị trừ khi click đúp, mong đợi giữ nguyên")
		}
	})

	t.Run("Báo lỗi khi chưa tới lượt", func(t *testing.T) {
		session := resetSession()
		session.CurrentActorID = "e1" // Đổi turn cho quái

		_, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "btn_click_2")
		if err != ErrNotYourTurn {
			t.Errorf("Mong đợi lỗi chưa tới lượt, nhận: %v", err)
		}
	})

	t.Run("Đánh quái đã chết (Báo lỗi)", func(t *testing.T) {
		session := resetSession()
		session.Enemies[0].CurrentHP = 0 // Giả lập quái đã bị hạ gục trước đó

		_, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "btn_click_3")
		if err != ErrTargetAlreadyDead {
			t.Errorf("Mong đợi ErrTargetAlreadyDead, nhận: %v", err)
		}
	})

	t.Run("Không tìm thấy mục tiêu (ID sai)", func(t *testing.T) {
		resetSession()
		_, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "ghost_target", "btn_click_4")
		if err != ErrTargetNotFound {
			t.Errorf("Mong đợi ErrTargetNotFound, nhận: %v", err)
		}
	})

	t.Run("Trận đấu hết hạn (Timeout)", func(t *testing.T) {
		session := resetSession()
		session.ExpiresAt = time.Now().Add(-1 * time.Hour) // Lùi thời gian về quá khứ

		_, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "btn_click_5")
		if err != ErrCombatSessionExpired {
			t.Errorf("Mong đợi ErrCombatSessionExpired, nhận: %v", err)
		}
	})

	t.Run("Giới hạn mảng IdempotencyKeys (Chống phình data MongoDB)", func(t *testing.T) {
		session := resetSession()
		// Cố tình nhồi nhét 15 keys cũ vào session
		for i := 0; i < 15; i++ {
			session.AddIdempotencyKey("old_key")
		}

		res, _ := svc.PlayerBasicAttack(context.Background(), "u1", "ss_1", "e1", "new_key")
		if len(res.IdempotencyKeys) > 10 {
			t.Errorf("Mong đợi IdempotencyKeys bị trim xuống 10 phần tử, nhận: %d", len(res.IdempotencyKeys))
		}
	})
}

func TestService_EnemyTurn_KillPlayer(t *testing.T) {
	repo := newFakeCombatRepo()
	svc, _ := NewService(repo, NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())

	session := &CombatSession{
		ID:             "ss_lose",
		UserID:         "u1",
		State:          StateActive,
		Turn:           1,
		Player:         CombatActor{ID: "u1", CurrentHP: 1, Stats: CombatStats{ATK: 10, DEF: 0, Speed: 100}}, // Máu rất mỏng
		Enemies:        []CombatActor{{ID: "e1", Name: "Boss", CurrentHP: 1000, Stats: CombatStats{ATK: 100, Speed: 150}}},
		CurrentActorID: "u1",
		ExpiresAt:      time.Now().Add(time.Hour),
		TurnOrder: []TurnOrderEntry{
			{ActorID: "u1", ActionValue: 100, Speed: 100},
			{ActorID: "e1", ActionValue: 150, Speed: 150},
		},
	}
	repo.sessions[session.ID] = session

	// Người chơi tấn công, không hạ được Boss -> Boss phản công và giết người chơi
	res, err := svc.PlayerBasicAttack(context.Background(), "u1", "ss_lose", "e1", "btn_1")
	if err != nil {
		t.Fatalf("Không mong đợi lỗi: %v", err)
	}
	if res.State != StateLost {
		t.Errorf("Mong đợi StateLost vì HP người chơi bằng 0, nhận: %v", res.State)
	}
	if res.Player.CurrentHP != 0 {
		t.Errorf("Mong đợi HP người chơi = 0, nhận: %d", res.Player.CurrentHP)
	}
}

func TestService_ExecuteAutoBattle(t *testing.T) {
	repo := newFakeCombatRepo()
	svc, _ := NewService(repo, NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())

	// Khởi tạo một phiên với quái cực kỳ trâu để test lặp nhịp đánh
	session := &CombatSession{
		ID:             "ss_auto",
		UserID:         "u1",
		State:          StateActive,
		Turn:           1,
		Player:         CombatActor{ID: "u1", CurrentHP: 100, Stats: CombatStats{ATK: 50, Speed: 120}},
		Enemies:        []CombatActor{{ID: "e1", CurrentHP: 2000, Stats: CombatStats{DEF: 0, Speed: 90}}},
		CurrentActorID: "u1",
		ExpiresAt:      time.Now().Add(time.Hour),
		TurnOrder: []TurnOrderEntry{
			{ActorID: "u1", ActionValue: 50, Speed: 120},
			{ActorID: "e1", ActionValue: 100, Speed: 90},
		},
	}
	repo.sessions[session.ID] = session

	t.Run("Chạy đủ số lượt MaxActions", func(t *testing.T) {
		opts := AutoBattleOptions{
			MaxActions:     3,
			IdempotencyKey: "auto_btn",
		}

		res, err := svc.ExecuteAutoBattle(context.Background(), "u1", "ss_auto", opts)
		if err != nil {
			t.Fatalf("Không mong đợi lỗi: %v", err)
		}
		if res.ActionsTaken != 3 {
			t.Errorf("Mong đợi thực hiện đúng 3 hành động, nhận %d", res.ActionsTaken)
		}
		if res.StoppedReason != "max_actions" {
			t.Errorf("Mong đợi lý do dừng là max_actions, nhận %s", res.StoppedReason)
		}
		if res.Session.Enemies[0].CurrentHP >= 2000 {
			t.Errorf("Quái chưa bị mất máu từ chuỗi tự đánh")
		}
	})

	t.Run("Auto đánh đến khi quái chết (Chiến thắng)", func(t *testing.T) {
		repo := newFakeCombatRepo()
		svc, _ := NewService(repo, NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())

		session := &CombatSession{
			ID:             "ss_auto_win",
			UserID:         "u1",
			State:          StateActive,
			Turn:           1,
			Player:         CombatActor{ID: "u1", CurrentHP: 100, Stats: CombatStats{ATK: 500, Speed: 100}}, // ATK cao đánh 1 phát chết quái
			Enemies:        []CombatActor{{ID: "e1", CurrentHP: 100, Stats: CombatStats{DEF: 0, Speed: 90}}},
			CurrentActorID: "u1",
			ExpiresAt:      time.Now().Add(time.Hour),
			TurnOrder: []TurnOrderEntry{
				{ActorID: "u1", ActionValue: 100, Speed: 100},
				{ActorID: "e1", ActionValue: 110, Speed: 90},
			},
		}
		repo.sessions[session.ID] = session

		opts := AutoBattleOptions{
			MaxActions:     100,
			IdempotencyKey: "auto_btn_win",
		}

		res, err := svc.ExecuteAutoBattle(context.Background(), "u1", "ss_auto_win", opts)
		if err != nil {
			t.Fatalf("Không mong đợi lỗi: %v", err)
		}
		if res.Session.State != StateWon {
			t.Errorf("Mong đợi trạng thái won, nhận %v", res.Session.State)
		}
		if res.StoppedReason != string(StateWon) {
			t.Errorf("Mong đợi lý do dừng là won, nhận %s", res.StoppedReason)
		}
		if res.ActionsTaken != 1 {
			t.Errorf("Mong đợi đánh đúng 1 nhịp là win, nhận %d", res.ActionsTaken)
		}
	})

	t.Run("Auto đánh đến khi người chơi chết (Thất bại)", func(t *testing.T) {
		repo := newFakeCombatRepo()
		svc, _ := NewService(repo, NewTurnOrderService(), rand.New(rand.NewSource(1)), zap.NewNop())

		session := &CombatSession{
			ID:             "ss_auto_lose",
			UserID:         "u1",
			State:          StateActive,
			Turn:           1,
			Player:         CombatActor{ID: "u1", CurrentHP: 1, Stats: CombatStats{ATK: 10, Speed: 100}}, // HP người chơi cực thấp
			Enemies:        []CombatActor{{ID: "e1", CurrentHP: 1000, Stats: CombatStats{ATK: 100, DEF: 0, Speed: 90}}}, // ATK quái cao
			CurrentActorID: "u1",
			ExpiresAt:      time.Now().Add(time.Hour),
			TurnOrder: []TurnOrderEntry{
				{ActorID: "u1", ActionValue: 100, Speed: 100},
				{ActorID: "e1", ActionValue: 110, Speed: 90},
			},
		}
		repo.sessions[session.ID] = session

		opts := AutoBattleOptions{
			MaxActions:     100,
			IdempotencyKey: "auto_btn_lose",
		}

		res, err := svc.ExecuteAutoBattle(context.Background(), "u1", "ss_auto_lose", opts)
		if err != nil {
			t.Fatalf("Không mong đợi lỗi: %v", err)
		}
		if res.Session.State != StateLost {
			t.Errorf("Mong đợi trạng thái lost, nhận %v", res.Session.State)
		}
		if res.StoppedReason != string(StateLost) {
			t.Errorf("Mong đợi lý do dừng là lost, nhận %s", res.StoppedReason)
		}
		if res.ActionsTaken != 1 {
			t.Errorf("Mong đợi đánh 1 nhịp rồi bị quái vả chết, nhận %d", res.ActionsTaken)
		}
	})
}
