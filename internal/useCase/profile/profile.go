package profile

import (
	"context"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
)

type profileRepo interface {
	Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
}

type UseCaseProfile struct {
	profileRepo profileRepo
}

func NewUseCaseProfile(profileRepo profileRepo) *UseCaseProfile {
	return &UseCaseProfile{
		profileRepo: profileRepo,
	}
}

func (u *UseCaseProfile) Create(ctx context.Context, r *profile.CreateRequestProfile) (*profile.Profile, error) {
	profileDto := &profile.Profile{
		DisplayName: r.DisplayName,
	}
	return u.profileRepo.Create(ctx, profileDto)
}
