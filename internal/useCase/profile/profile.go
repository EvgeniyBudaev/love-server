package profile

import (
	"context"
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type profileRepo interface {
	Add(ctx context.Context, p *profile.Profile) (*profile.Profile, error)
	FindById(ctx context.Context, id uint64) (*profile.Profile, error)
	AddTelegram(ctx context.Context, p *profile.TelegramProfile) (*profile.TelegramProfile, error)
	FindTelegramById(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error)
	AddImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error)
	SelectListPublicImage(ctx context.Context, id uint64) ([]*profile.ImageProfile, error)
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

func (u *UseCaseProfile) Add(ctf *fiber.Ctx, r *profile.AddRequestProfile) (*profile.Profile, error) {
	filePath := "./static/uploads/profile/images/defaultImage.jpg"
	directoryPath := fmt.Sprintf("./static/uploads/profile/images")
	form, err := ctf.MultipartForm()
	if err != nil {
		u.logger.Debug("error func Add, method MultipartForm by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	imageFiles := form.File["image"]
	imagesFilePath := make([]string, 0, len(imageFiles))
	imagesProfile := make([]*profile.ImageProfile, 0, len(imagesFilePath))
	for _, file := range imageFiles {
		filePath = fmt.Sprintf("%s/%s", directoryPath, file.Filename)
		if err := ctf.SaveFile(file, filePath); err != nil {
			u.logger.Debug("error func Add, method SaveFile by path internal/useCase/profile/profile.go",
				zap.Error(err))
			return nil, err
		}
		image := profile.ImageProfile{
			Name:      file.Filename,
			Url:       filePath,
			Size:      file.Size,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			IsDeleted: false,
			IsBlocked: false,
			IsPrimary: false,
			IsPrivate: false,
		}
		imagesFilePath = append(imagesFilePath, filePath)
		imagesProfile = append(imagesProfile, &image)
	}
	profileDto := &profile.Profile{
		DisplayName: r.DisplayName,
		Age:         r.Age,
		Gender:      r.Gender,
		Location:    r.Location,
		Description: r.Description,
		IsDeleted:   false,
		IsBlocked:   false,
		IsPremium:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastOnline:  time.Now(),
		Images:      imagesProfile,
	}
	newProfile, err := u.profileRepo.Add(ctf.Context(), profileDto)
	if err != nil {
		u.logger.Debug("error func Add, method profileRepo.Add by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	for _, i := range profileDto.Images {
		image := &profile.ImageProfile{
			ProfileID: newProfile.ID,
			Name:      i.Name,
			Url:       i.Url,
			Size:      i.Size,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
			IsDeleted: i.IsDeleted,
			IsBlocked: i.IsBlocked,
			IsPrimary: i.IsPrimary,
			IsPrivate: i.IsPrivate,
		}
		_, err := u.profileRepo.AddImage(ctf.Context(), image)
		if err != nil {
			u.logger.Debug("error func Add, method AddImage by path internal/useCase/profile/profile.go",
				zap.Error(err))
			return nil, err
		}
	}
	telegramDto := &profile.TelegramProfile{
		ProfileID:       newProfile.ID,
		TelegramID:      r.TelegramID,
		UserName:        r.UserName,
		Firstname:       r.Firstname,
		Lastname:        r.Lastname,
		LanguageCode:    r.LanguageCode,
		AllowsWriteToPm: r.AllowsWriteToPm,
		QueryID:         r.QueryID,
	}
	t, err := u.profileRepo.AddTelegram(ctf.Context(), telegramDto)
	if err != nil {
		u.logger.Debug("error func Add, method AddTelegram by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	p, err := u.profileRepo.FindById(ctf.Context(), newProfile.ID)
	if err != nil {
		u.logger.Debug("error func Add, method FindById by path internal/useCase/profile/profile.go", zap.Error(err))
		return nil, err
	}
	i, err := u.profileRepo.SelectListPublicImage(ctf.Context(), p.ID)
	if err != nil {
		u.logger.Debug("error func FindById, method SelectListPublicImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	response := &profile.Profile{
		ID:          p.ID,
		DisplayName: p.DisplayName,
		Age:         p.Age,
		Gender:      p.Gender,
		Location:    p.Location,
		Description: p.Description,
		IsDeleted:   p.IsDeleted,
		IsBlocked:   p.IsBlocked,
		IsPremium:   p.IsPremium,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		LastOnline:  p.LastOnline,
		Telegram:    t,
		Images:      i,
	}
	return response, nil
}

func (u *UseCaseProfile) FindById(ctf *fiber.Ctx) (*profile.Profile, error) {
	idStr := ctf.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		u.logger.Debug("error func FindById, method ParseUint by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	p, err := u.profileRepo.FindById(ctf.Context(), id)
	if err != nil {
		u.logger.Debug("error func FindById, method FindById by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	t, err := u.profileRepo.FindTelegramById(ctf.Context(), id)
	if err != nil {
		u.logger.Debug("error func FindById, method FindTelegramById by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	i, err := u.profileRepo.SelectListPublicImage(ctf.Context(), id)
	if err != nil {
		u.logger.Debug("error func FindById, method SelectListPublicImage by path internal/useCase/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	response := &profile.Profile{
		ID:          p.ID,
		DisplayName: p.DisplayName,
		Age:         p.Age,
		Gender:      p.Gender,
		Location:    p.Location,
		Description: p.Description,
		IsDeleted:   p.IsDeleted,
		IsBlocked:   p.IsBlocked,
		IsPremium:   p.IsPremium,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		LastOnline:  p.LastOnline,
		Telegram:    t,
		Images:      i,
	}
	return response, nil
}
