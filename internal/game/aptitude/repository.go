// File: internal/game/aptitude/repository.go
package aptitude

import (
	"context"
	"time"
)

type AptitudeProfile struct {
	UserID      string         `bson:"userId"`
	AptitudeID  string         `bson:"aptitudeId"`
	Rarity      AptitudeRarity `bson:"rarity"`
	RolledAt    time.Time      `bson:"rolledAt"`
	RerollCount int            `bson:"rerollCount"`
	Locked      bool           `bson:"locked"`
	Seed        string         `bson:"seed"`
}

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*AptitudeProfile, error)
	Create(ctx context.Context, profile *AptitudeProfile) error
}
