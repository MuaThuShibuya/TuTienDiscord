// File: internal/game/aptitude/service.go
package aptitude

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/whiskey/tu-tien-bot/internal/apperrors"
)

type Service interface {
	RollForNewCharacter(ctx context.Context, userID string) (*AptitudeProfile, *AptitudeDefinition, error)
	GetProfile(ctx context.Context, userID string) (*AptitudeProfile, *AptitudeDefinition, error)
}

type serviceImpl struct{ repo Repository }

func NewService(repo Repository) Service { return &serviceImpl{repo: repo} }

func (s *serviceImpl) RollForNewCharacter(ctx context.Context, userID string) (*AptitudeProfile, *AptitudeDefinition, error) {
	existing, err := s.repo.GetByUserID(ctx, userID)
	if err == nil && existing != nil {
		def := Registry[existing.AptitudeID]
		return existing, &def, nil
	}
	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return nil, nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	def := GetRandomAptitude(rng.Intn)
	profile := &AptitudeProfile{UserID: userID, AptitudeID: def.ID, Rarity: def.Rarity, RolledAt: time.Now().UTC(), Locked: true, RerollCount: 0}
	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, nil, err
	}
	return profile, &def, nil
}

func (s *serviceImpl) GetProfile(ctx context.Context, userID string) (*AptitudeProfile, *AptitudeDefinition, error) {
	profile, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	def, ok := Registry[profile.AptitudeID]
	if !ok {
		def = Registry["apt_pham_tu"]
	}
	return profile, &def, nil
}
