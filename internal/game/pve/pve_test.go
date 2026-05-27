// File: internal/game/pve/pve_test.go

package pve

import (
	"math/rand"
	"testing"
)

func TestPvERegistry(t *testing.T) {
	duNgoan := AreaRegistry["area_du_ngoan_rung_truc"]
	biCanh := AreaRegistry["area_bi_canh_thach_dong"]

	if duNgoan.ActivityType != ActivityDuNgoan {
		t.Errorf("Rừng trúc phải là Du Ngoạn")
	}
	if biCanh.ActivityType != ActivityBiCanh {
		t.Errorf("Thạch động phải là Bí Cảnh")
	}
	if duNgoan.RewardPoolID == biCanh.RewardPoolID {
		t.Errorf("Reward pool của Bí Cảnh phải khác Du Ngoạn")
	}
}

func TestRewardPoolRegistry(t *testing.T) {
	biCanhPool, ok := RewardPoolRegistry["reward_bi_canh_rare"]
	if !ok {
		t.Fatalf("Không tìm thấy reward_bi_canh_rare")
	}

	hasRareType := false
	for _, entry := range biCanhPool.Entries {
		if entry.Type == "artifact" || entry.Type == "equipment" || entry.Type == "skill_scroll" {
			hasRareType = true
			break
		}
	}
	if !hasRareType {
		t.Errorf("Bí cảnh phải có rớt đồ hiếm (artifact/equipment/skill_scroll)")
	}
}

func TestStageLimits(t *testing.T) {
	area := AreaRegistry["area_du_ngoan_rung_truc"]

	if _, err := GetStageDefinition(area, 1); err != nil {
		t.Errorf("Stage 1 phải hợp lệ: %v", err)
	}
	if _, err := GetStageDefinition(area, 30); err != nil {
		t.Errorf("Stage 30 phải hợp lệ: %v", err)
	}

	_, err := GetStageDefinition(area, 31)
	if err == nil {
		t.Errorf("Stage > 30 phải bị chặn (trả về lỗi)")
	}

	stage30, _ := GetStageDefinition(area, 30)
	if !stage30.IsBossStage {
		t.Errorf("Stage 30 phải được đánh dấu là Boss Stage")
	}
}

func TestScaling(t *testing.T) {
	base := MonsterStats{MaxHP: 100, ATK: 10, DEF: 5, Speed: 100}
	stage1 := ScaleMonsterStats(base, 1, ActivityDuNgoan, MonsterRoleNormal, DefaultScalingConfig)
	stage10 := ScaleMonsterStats(base, 10, ActivityDuNgoan, MonsterRoleNormal, DefaultScalingConfig)
	stage30 := ScaleMonsterStats(base, 30, ActivityDuNgoan, MonsterRoleNormal, DefaultScalingConfig)
	biCanh := ScaleMonsterStats(base, 30, ActivityBiCanh, MonsterRoleBoss, DefaultScalingConfig)

	if stage10.MaxHP <= stage1.MaxHP || stage10.ATK <= stage1.ATK {
		t.Errorf("Stage 10 phải mạnh hơn Stage 1")
	}
	if stage30.MaxHP <= stage10.MaxHP || stage30.ATK <= stage10.ATK {
		t.Errorf("Stage 30 phải mạnh hơn Stage 10")
	}
	if biCanh.MaxHP <= stage30.MaxHP {
		t.Errorf("Bí Cảnh Boss phải mạnh hơn Du Ngoạn Normal cùng Stage")
	}
	if biCanh.Speed > 150 {
		t.Errorf("Speed bị rò rỉ tăng vô hạn. Giá trị hiện tại: %v", biCanh.Speed)
	}
}

func TestEncounterLimits(t *testing.T) {
	area := AreaRegistry["area_du_ngoan_rung_truc"]
	rng := rand.New(rand.NewSource(1))

	// Test stage sâu nhất để xem số lượng quái có vượt 9 không
	enc, err := GenerateEncounter(area, 30, rng)
	if err != nil {
		t.Fatalf("Không mong đợi lỗi khi generate: %v", err)
	}

	if len(enc.Enemies) > 9 {
		t.Errorf("Số lượng quái vượt quá 9: nhận %d", len(enc.Enemies))
	}
}

func TestRewardResolver(t *testing.T) {
	rng := rand.New(rand.NewSource(1))

	// Roll Bí Cảnh
	res := ResolveStageRewards("reward_bi_canh_rare", "reward_bi_canh_bonus", rng)

	hasGuaranteed := false
	for _, g := range res.Guaranteed {
		if g.Type == "exp" || g.Type == "stones" {
			hasGuaranteed = true
		}
	}
	if !hasGuaranteed {
		t.Errorf("Phải có exp/stones trong guaranteed pool")
	}

	// Roll bonus nhiều lần để test gacha logic
	foundBonus := false
	for i := 0; i < 50; i++ {
		res2 := ResolveStageRewards("reward_bi_canh_rare", "reward_bi_canh_bonus", rand.New(rand.NewSource(int64(i))))
		if len(res2.Bonus) > 0 {
			foundBonus = true
			break
		}
	}
	if !foundBonus {
		t.Errorf("Bonus pool (Weight) đang không rớt món nào sau 50 lần thử")
	}
}
