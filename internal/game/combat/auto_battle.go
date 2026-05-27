// File: internal/game/combat/auto_battle.go
// Chức năng: Nền tảng cấu hình và giao diện ra quyết định cho hệ thống Đánh Tự Động (Auto Battle).

package combat

type AutoBattlePolicy struct {
	Enabled                  bool    `bson:"enabled" json:"enabled"`
	PreferSkill              bool    `bson:"preferSkill" json:"preferSkill"`
	UseUltimateWhenAvailable bool    `bson:"useUltimateWhenAvailable" json:"useUltimateWhenAvailable"`
	StopWhenLowHPPercent     float64 `bson:"stopWhenLowHpPercent" json:"stopWhenLowHpPercent"`
}

// AutoBattleDecider định nghĩa hợp đồng ra quyết định hành động tự động.
// TODO: v0.4.x - Implement decider sử dụng Skill Service và Turn Order.
type AutoBattleDecider interface {
	// DecideAction nhận vào Session và Actor đang tới lượt, trả về ID Kỹ Năng và Mục Tiêu.
	DecideAction(session *CombatSession, actor *CombatActor) (skillID string, targetID string)
}
