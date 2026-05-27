// File: internal/game/pve/types.go

package pve

type ActivityType string

const (
	ActivityDuNgoan ActivityType = "du_ngoan"
	ActivityBiCanh  ActivityType = "bi_canh"
)

type EntryCost struct {
	Stamina           int64 `bson:"stamina" json:"stamina"`
	Stones            int64 `bson:"stones" json:"stones"`
	OpportunityTicket int64 `bson:"opportunityTicket" json:"opportunityTicket"`
}

type UnlockRule struct {
	Type        string `bson:"type" json:"type"`
	Value       string `bson:"value" json:"value"`
	NumberValue int64  `bson:"numberValue" json:"numberValue"`
}

type PvEAreaDefinition struct {
	ID                  string       `bson:"id" json:"id"`
	Name                string       `bson:"name" json:"name"`
	ActivityType        ActivityType `bson:"activityType" json:"activityType"`
	Description         string       `bson:"description" json:"description"`
	RequiredRealm       string       `bson:"requiredRealm" json:"requiredRealm"`
	RequiredCombatPower int64        `bson:"requiredCombatPower" json:"requiredCombatPower"`
	MinStage            int          `bson:"minStage" json:"minStage"`
	MaxStage            int          `bson:"maxStage" json:"maxStage"`
	MonsterPoolID       string       `bson:"monsterPoolId" json:"monsterPoolId"`
	RewardPoolID        string       `bson:"rewardPoolId" json:"rewardPoolId"`
	BonusRewardPoolID   string       `bson:"bonusRewardPoolId" json:"bonusRewardPoolId"` // Gacha bonus reward
	EntryCost           EntryCost    `bson:"entryCost" json:"entryCost"`
	UnlockRules         []UnlockRule `bson:"unlockRules" json:"unlockRules"`
	IsTimed             bool         `bson:"isTimed" json:"isTimed"`
	ScheduleID          string       `bson:"scheduleId" json:"scheduleId"`
	Tags                []string     `bson:"tags" json:"tags"`
}
