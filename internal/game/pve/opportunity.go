// File: internal/game/pve/opportunity.go
// Chức năng: Quản lý Cơ duyên / Bí cảnh đặc biệt xuất hiện theo thời gian thực (chuẩn bị cho v0.7 Boss/Event).

package pve

import (
	"context"
	"time"
)

type TimedOpportunityDefinition struct {
	ID                  string       `bson:"id" json:"id"`
	Name                string       `bson:"name" json:"name"`
	Description         string       `bson:"description" json:"description"`
	ActivityType        ActivityType `bson:"activityType" json:"activityType"`
	AreaID              string       `bson:"areaId" json:"areaId"` // Trỏ tới PvEAreaDefinition
	StartsAt            time.Time    `bson:"startsAt" json:"startsAt"`
	EndsAt              time.Time    `bson:"endsAt" json:"endsAt"`
	NotifyBeforeMinutes int          `bson:"notifyBeforeMinutes" json:"notifyBeforeMinutes"`
	RewardPoolID        string       `bson:"rewardPoolId" json:"rewardPoolId"`
	RequiredRealm       string       `bson:"requiredRealm" json:"requiredRealm"`
	RequiredCombatPower int64        `bson:"requiredCombatPower" json:"requiredCombatPower"`
	MaxAttemptsPerUser  int          `bson:"maxAttemptsPerUser" json:"maxAttemptsPerUser"`
	GlobalAttemptLimit  int          `bson:"globalAttemptLimit" json:"globalAttemptLimit"` // Số lần tối đa toàn server có thể tham gia (vd: Boss thế giới)
	Tags                []string     `bson:"tags" json:"tags"`
}

// OpportunityRepository định nghĩa interface quản lý cơ duyên.
// TODO: Triển khai MongoDB implementation ở block sau.
type OpportunityRepository interface {
	ListActiveOpportunities(ctx context.Context, now time.Time) ([]*TimedOpportunityDefinition, error)
}
