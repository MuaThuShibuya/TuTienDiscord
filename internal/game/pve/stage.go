// File: internal/game/pve/stage.go
// Chức năng: Quản lý thiết lập ải (Stage). Chặn giới hạn ải, tính toán lượng quái.

package pve

import "fmt"

const DefaultMaxStage = 30
const MaxEnemiesPerEncounter = 9 // Phù hợp với UI Discord (không làm rối màn hình)

type StageDefinition struct {
	AreaID                 string
	Stage                  int
	Name                   string
	RecommendedCombatPower int64
	MonsterPoolID          string
	EnemyCountMin          int
	EnemyCountMax          int
	GuaranteedRewardPoolID string
	BonusRewardPoolID      string
	IsBossStage            bool
	Tags                   []string
}

// GetStageDefinition sinh cấu hình ải động dựa trên khu vực và cấp độ ải.
// Không lưu database 30 object, tiết kiệm tài nguyên.
func GetStageDefinition(area PvEAreaDefinition, stage int) (StageDefinition, error) {
	maxStage := area.MaxStage
	if maxStage <= 0 {
		maxStage = DefaultMaxStage
	}
	if stage < 1 || stage > maxStage {
		return StageDefinition{}, fmt.Errorf("ải %d không hợp lệ (tối đa %d)", stage, maxStage)
	}

	return StageDefinition{
		AreaID:                 area.ID,
		Stage:                  stage,
		Name:                   fmt.Sprintf("%s - Ải %d", area.Name, stage),
		RecommendedCombatPower: int64(100 * stage), // Nháp: 100 * stage
		MonsterPoolID:          area.MonsterPoolID,
		EnemyCountMin:          1 + (stage / 10), // Càng sâu quái min càng tăng
		EnemyCountMax:          3 + (stage / 5),  // Càng sâu quái max càng tăng
		GuaranteedRewardPoolID: area.RewardPoolID,
		BonusRewardPoolID:      area.BonusRewardPoolID,
		IsBossStage:            stage%5 == 0 || stage == maxStage,
	}, nil
}
