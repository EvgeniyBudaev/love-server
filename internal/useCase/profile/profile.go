package profile

import (
	"context"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type profileRepo interface {
	Add(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	SelectList(ctx context.Context, qp *profile.QueryParamsProfileList) (*profile.ResponseListProfile, error)
	FindById(ctx context.Context, id uint64) (*profile.Profile, error)
	AddTelegram(ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error)
	UpdateTelegram(ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error)
	FindTelegramById(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error)
	AddImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	UpdateImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	SelectListPublicImage(ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error)
	CheckIfCommonImageExists(ctx context.Context, profileID uint64, fileName string) (bool, uint64, error)
}

type UseCaseProfile struct {
	logger      logger.Logger
	profileRepo profileRepo
}

func NewUseCaseProfile(l logger.Logger, pr profileRepo) *UseCaseProfile {
	return &UseCaseProfile{
		logger:      l,
		profileRepo: pr,
	}
}

func (u *UseCaseProfile) Add(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	response, err := u.profileRepo.Add(ctx, p)
	if err != nil {
		u.logger.Debug("error func Add, method Add by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	response, err := u.profileRepo.Update(ctx, p)
	if err != nil {
		u.logger.Debug("error func response, method Update by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) SelectList(ctf *fiber.Ctx) (*profile.ResponseListProfile, error) {
	var params profile.QueryParamsProfileList
	if err := ctf.QueryParser(&params); err != nil {
		u.logger.Debug("error func SelectList, method QueryParser by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	response, err := u.profileRepo.SelectList(ctf.Context(), &params)
	if err != nil {
		u.logger.Debug("error func SelectList, method SelectList by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) FindById(ctx context.Context, id uint64) (*profile.Profile, error) {
	response, err := u.profileRepo.FindById(ctx, id)
	if err != nil {
		u.logger.Debug("error func FindById, method FindById by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) AddImage(ctx context.Context, i *profile.ImageProfile) (*profile.ImageProfile, error) {
	response, err := u.profileRepo.AddImage(ctx, i)
	if err != nil {
		u.logger.Debug("error func AddImage, method AddImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) UpdateImage(ctx context.Context, i *profile.ImageProfile) (*profile.ImageProfile, error) {
	response, err := u.profileRepo.UpdateImage(ctx, i)
	if err != nil {
		u.logger.Debug("error func AddImage, method UpdateImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) SelectListPublicImage(ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error) {
	response, err := u.profileRepo.SelectListPublicImage(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func SelectListPublicImage, method SelectListPublicImage by path"+
			" internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) CheckIfCommonImageExists(
	ctx context.Context, profileID uint64, fileName string) (bool, uint64, error) {
	return u.profileRepo.CheckIfCommonImageExists(ctx, profileID, fileName)
}

func (u *UseCaseProfile) AddTelegram(
	ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	response, err := u.profileRepo.AddTelegram(ctx, t)
	if err != nil {
		u.logger.Debug("error func AddTelegram, method AddTelegram by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) UpdateTelegram(
	ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	response, err := u.profileRepo.UpdateTelegram(ctx, t)
	if err != nil {
		u.logger.Debug("error func UpdateTelegram, method AddTelegram by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) FindTelegramById(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error) {
	response, err := u.profileRepo.FindTelegramById(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func FindTelegramById, method FindTelegramById by path "+
			"internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}
