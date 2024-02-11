package profile

import (
	"context"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"go.uber.org/zap"
)

type profileRepo interface {
	Add(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	UpdateLastOnline(ctx context.Context, profileID uint64) error
	Delete(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	SelectList(ctx context.Context, qp *profile.QueryParamsProfileList) (*profile.ResponseListProfile, error)
	FindById(ctx context.Context, id uint64) (*profile.Profile, error)
	FindByTelegramId(ctx context.Context, telegramID uint64) (*profile.Profile, error)
	AddTelegram(ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error)
	UpdateTelegram(ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error)
	DeleteTelegram(ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error)
	FindTelegramByProfileID(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error)
	AddNavigator(ctx context.Context, p *profile.NavigatorProfile) (*profile.NavigatorProfile, error)
	UpdateNavigator(ctx context.Context, p *profile.NavigatorProfile) (*profile.NavigatorProfile, error)
	FindNavigatorByProfileID(ctx context.Context, profileID uint64) (*profile.NavigatorProfile, error)
	DeleteImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	AddImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	UpdateImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	FindImageById(ctx context.Context, imageID uint64) (*profile.ImageProfile, error)
	SelectListPublicImage(ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error)
	SelectListImage(ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error)
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
		u.logger.Debug("error func Update, method Update by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) UpdateLastOnline(ctx context.Context, profileID uint64) error {
	err := u.profileRepo.UpdateLastOnline(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func UpdateLastOnline, method UpdateLastOnline by path"+
			" internal/useCase/profile/profile.go", zap.Error(err))
		return err
	}
	return nil
}

func (u *UseCaseProfile) Delete(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	response, err := u.profileRepo.Delete(ctx, p)
	if err != nil {
		u.logger.Debug("error func Delete, method Delete by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) SelectList(
	ctx context.Context, qp *profile.QueryParamsProfileList) (*profile.ResponseListProfile, error) {
	response, err := u.profileRepo.SelectList(ctx, qp)
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

func (u *UseCaseProfile) FindByTelegramId(ctx context.Context, telegramID uint64) (*profile.Profile, error) {
	response, err := u.profileRepo.FindByTelegramId(ctx, telegramID)
	if err != nil {
		u.logger.Debug("error func FindByTelegramId, methodFindByTelegramId by path"+
			" internal/useCase/profile/profile.go", zap.Error(err))
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
		u.logger.Debug("error func UpdateImage, method UpdateImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) DeleteImage(ctx context.Context, i *profile.ImageProfile) (*profile.ImageProfile, error) {
	response, err := u.profileRepo.DeleteImage(ctx, i)
	if err != nil {
		u.logger.Debug("error func DeleteImage, method DeleteImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) FindImageById(ctx context.Context, imageID uint64) (*profile.ImageProfile, error) {
	response, err := u.profileRepo.FindImageById(ctx, imageID)
	if err != nil {
		u.logger.Debug("error func FindImageById, method FindImageById by path internal/useCase/profile/profile.go",
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

func (u *UseCaseProfile) SelectListImage(ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error) {
	response, err := u.profileRepo.SelectListImage(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func SelectListImage, method SelectListImage by path"+
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
		u.logger.Debug("error func UpdateTelegram, method UpdateTelegram by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) DeleteTelegram(
	ctx context.Context, t *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	response, err := u.profileRepo.DeleteTelegram(ctx, t)
	if err != nil {
		u.logger.Debug("error func DeleteTelegram, method DeleteTelegram by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) FindTelegramByProfileID(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error) {
	response, err := u.profileRepo.FindTelegramByProfileID(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func FindTelegramByProfileID, method FindTelegramByProfileID by path "+
			"internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) AddNavigator(
	ctx context.Context, n *profile.NavigatorProfile) (*profile.NavigatorProfile, error) {
	response, err := u.profileRepo.AddNavigator(ctx, n)
	if err != nil {
		u.logger.Debug("error func AddNavigator, method AddNavigator by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) UpdateNavigator(
	ctx context.Context, n *profile.NavigatorProfile) (*profile.NavigatorProfile, error) {
	response, err := u.profileRepo.UpdateNavigator(ctx, n)
	if err != nil {
		u.logger.Debug("error func UpdateTNavigator, method UpdateNavigator by path"+
			" internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}

func (u *UseCaseProfile) FindNavigatorByProfileID(
	ctx context.Context, profileID uint64) (*profile.NavigatorProfile, error) {
	response, err := u.profileRepo.FindNavigatorByProfileID(ctx, profileID)
	if err != nil {
		u.logger.Debug("error func FindNavigatorByProfileID, method FindNavigatorByProfileID by path "+
			"internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	return response, nil
}
