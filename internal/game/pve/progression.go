// File: internal/game/pve/progression.go

package pve

import "time"

type AreaProgress struct {
	AreaID              string       `bson:"areaId" json:"areaId"`
	ActivityType        ActivityType `bson:"activityType" json:"activityType"`
	HighestStageCleared int          `bson:"highestStageCleared" json:"highestStageCleared"`
	AttemptsToday       int          `bson:"attemptsToday" json:"attemptsToday"`
	LastAttemptAt       time.Time    `bson:"lastAttemptAt" json:"lastAttemptAt"`
	UpdatedAt           time.Time    `bson:"updatedAt" json:"updatedAt"`
}

type UserPvEProgress struct {
	UserID    string                  `bson:"userId" json:"userId"`
	Areas     map[string]AreaProgress `bson:"areas" json:"areas"`
	UpdatedAt time.Time               `bson:"updatedAt" json:"updatedAt"`
}
