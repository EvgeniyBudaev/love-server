package profile

import (
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	errorDomain "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/error"
	r "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/response"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	profileUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strconv"
	"time"
)

type HandlerProfile struct {
	logger logger.Logger
	uc     *profileUseCase.UseCaseProfile
}

func NewHandlerProfile(l logger.Logger, uc *profileUseCase.UseCaseProfile) *HandlerProfile {
	return &HandlerProfile{logger: l, uc: uc}
}

func (h *HandlerProfile) AddProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/add")
		req := profile.RequestAddProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func AddProfileHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filePath := fmt.Sprintf("static/uploads/profile/%s/images/defaultImage.jpg", req.UserName)
		directoryPath := fmt.Sprintf("static/uploads/profile/%s/images", req.UserName)
		if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
			if err := os.MkdirAll(directoryPath, 0755); err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method MkdirAll by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		form, err := ctf.MultipartForm()
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method MultipartForm by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageFiles := form.File["image"]
		imagesFilePath := make([]string, 0, len(imageFiles))
		imagesProfile := make([]*profile.ImageProfile, 0, len(imagesFilePath))
		for _, file := range imageFiles {
			filePath = fmt.Sprintf("%s/%s", directoryPath, file.Filename)
			if err := ctf.SaveFile(file, filePath); err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method SaveFile by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
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
		heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		height := uint8(heightUint64)
		weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		weight := uint8(weightUint64)
		profileDto := &profile.Profile{
			DisplayName:    req.DisplayName,
			Birthday:       req.Birthday,
			Gender:         req.Gender,
			SearchGender:   req.SearchGender,
			Location:       req.Location,
			Description:    req.Description,
			Height:         height,
			Weight:         weight,
			LookingFor:     req.LookingFor,
			IsDeleted:      false,
			IsBlocked:      false,
			IsPremium:      false,
			IsShowDistance: true,
			IsInvisible:    false,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			LastOnline:     time.Now(),
			Images:         imagesProfile,
		}
		newProfile, err := h.uc.Add(ctf.Context(), profileDto)
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
			_, err := h.uc.AddImage(ctf.Context(), image)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method AddImage by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		telegramID, err := strconv.ParseUint(req.TelegramID, 10, 64)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseUint roomIdStr by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		allowsWriteToPm, err := strconv.ParseBool(req.AllowsWriteToPm)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseBool roomIdStr by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		telegramDto := &profile.TelegramProfile{
			ProfileID:       newProfile.ID,
			TelegramID:      telegramID,
			UserName:        req.UserName,
			Firstname:       req.Firstname,
			Lastname:        req.Lastname,
			LanguageCode:    req.LanguageCode,
			AllowsWriteToPm: allowsWriteToPm,
			QueryID:         req.QueryID,
		}
		_, err = h.uc.AddTelegram(ctf.Context(), telegramDto)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method AddTelegram by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), newProfile.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindById by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		t, err := h.uc.FindTelegramById(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindTelegramById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.Profile{
			ID:             p.ID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			SearchGender:   p.SearchGender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
			LookingFor:     p.LookingFor,
			IsDeleted:      p.IsDeleted,
			IsBlocked:      p.IsBlocked,
			IsPremium:      p.IsPremium,
			IsShowDistance: p.IsShowDistance,
			IsInvisible:    p.IsInvisible,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			LastOnline:     p.LastOnline,
			Images:         i,
			Telegram:       t,
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileListHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/list")
		response, err := h.uc.SelectList(ctf)
		if err != nil {
			h.logger.Debug(
				"error func GetProfileListHandler, method SelectList by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/:id")
		idStr := ctf.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetProfileHandler, method ParseUint by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), id)
		if err != nil {
			h.logger.Debug(
				"error func GetProfileHandler, method FindById by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		t, err := h.uc.FindTelegramById(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileHandler, method FindTelegramById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.Profile{
			ID:             p.ID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			SearchGender:   p.SearchGender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
			LookingFor:     p.LookingFor,
			IsDeleted:      p.IsDeleted,
			IsBlocked:      p.IsBlocked,
			IsPremium:      p.IsPremium,
			IsShowDistance: p.IsShowDistance,
			IsInvisible:    p.IsInvisible,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			LastOnline:     p.LastOnline,
			Images:         i,
			Telegram:       t,
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) UpdateProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/edit")
		req := profile.RequestUpdateProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method ParseUint roomIdStr by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileInDB, err := h.uc.FindById(ctf.Context(), profileID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler,"+
					" method FindById by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if profileInDB.IsDeleted == true {
			msg := errors.Wrap(err, "user has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if profileInDB.IsBlocked == true {
			msg := errors.Wrap(err, "user has already been blocked")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		filePath := fmt.Sprintf("static/uploads/profile/%s/images/defaultImage.jpg", req.UserName)
		directoryPath := fmt.Sprintf("static/uploads/profile/%s/images", req.UserName)
		if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
			if err := os.MkdirAll(directoryPath, 0755); err != nil {
				h.logger.Debug(
					"error func UpdateProfileHandler, method MkdirAll by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		form, err := ctf.MultipartForm()
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method MultipartForm by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		height := uint8(heightUint64)
		weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		weight := uint8(weightUint64)
		imageFiles := form.File["image"]
		profileDto := &profile.Profile{}
		if len(imageFiles) > 0 {
			imagesFilePath := make([]string, 0, len(imageFiles))
			imagesProfile := make([]*profile.ImageProfile, 0, len(imagesFilePath))
			for _, file := range imageFiles {
				filePath = fmt.Sprintf("%s/%s", directoryPath, file.Filename)
				if err := ctf.SaveFile(file, filePath); err != nil {
					h.logger.Debug(
						"error func UpdateProfileHandler, method SaveFile by path internal/handler/profile/profile.go",
						zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
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
			profileDto = &profile.Profile{
				ID:             profileID,
				DisplayName:    req.DisplayName,
				Birthday:       req.Birthday,
				Gender:         req.Gender,
				SearchGender:   req.SearchGender,
				Location:       req.Location,
				Description:    req.Description,
				Height:         height,
				Weight:         weight,
				LookingFor:     req.LookingFor,
				IsDeleted:      profileInDB.IsDeleted,
				IsBlocked:      profileInDB.IsBlocked,
				IsPremium:      profileInDB.IsPremium,
				IsShowDistance: profileInDB.IsShowDistance,
				IsInvisible:    profileInDB.IsInvisible,
				CreatedAt:      profileInDB.CreatedAt,
				UpdatedAt:      time.Now(),
				LastOnline:     time.Now(),
				Images:         imagesProfile,
			}
		} else {
			profileDto = &profile.Profile{
				ID:             profileID,
				DisplayName:    req.DisplayName,
				Birthday:       req.Birthday,
				Gender:         req.Gender,
				SearchGender:   req.SearchGender,
				Location:       req.Location,
				Description:    req.Description,
				Height:         height,
				Weight:         weight,
				LookingFor:     req.LookingFor,
				IsDeleted:      profileInDB.IsDeleted,
				IsBlocked:      profileInDB.IsBlocked,
				IsPremium:      profileInDB.IsPremium,
				IsShowDistance: profileInDB.IsShowDistance,
				IsInvisible:    profileInDB.IsInvisible,
				CreatedAt:      profileInDB.CreatedAt,
				UpdatedAt:      time.Now(),
				LastOnline:     time.Now(),
			}
		}
		profileUpdated, err := h.uc.Update(ctf.Context(), profileDto)
		if len(imageFiles) > 0 {
			for _, i := range profileDto.Images {
				exists, imageID, err := h.uc.CheckIfCommonImageExists(ctf.Context(), profileUpdated.ID, i.Name)
				if err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method CheckIfCommonImageExists by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				if !exists {
					image := &profile.ImageProfile{
						ProfileID: profileUpdated.ID,
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
					_, err := h.uc.AddImage(ctf.Context(), image)
					if err != nil {
						h.logger.Debug("error func UpdateProfileHandler, method AddImage by path"+
							" internal/handler/profile/profile.go", zap.Error(err))
						return r.WrapError(ctf, err, http.StatusBadRequest)
					}
				} else {
					image := &profile.ImageProfile{
						ID:        imageID,
						ProfileID: profileUpdated.ID,
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
					_, err := h.uc.UpdateImage(ctf.Context(), image)
					if err != nil {
						h.logger.Debug("error func UpdateProfileHandler, method UpdateImage by path"+
							" internal/handler/profile/profile.go", zap.Error(err))
						return r.WrapError(ctf, err, http.StatusBadRequest)
					}
				}
			}
		}
		telegramID, err := strconv.ParseUint(req.TelegramID, 10, 64)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method ParseUint roomIdStr by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		allowsWriteToPm, err := strconv.ParseBool(req.AllowsWriteToPm)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method ParseBool roomIdStr by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		telegramDto := &profile.TelegramProfile{
			ID:              telegramID,
			ProfileID:       profileUpdated.ID,
			TelegramID:      telegramID,
			UserName:        req.UserName,
			Firstname:       req.Firstname,
			Lastname:        req.Lastname,
			LanguageCode:    req.LanguageCode,
			AllowsWriteToPm: allowsWriteToPm,
			QueryID:         req.QueryID,
		}
		_, err = h.uc.UpdateTelegram(ctf.Context(), telegramDto)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method AddTelegram by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileUpdated.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		t, err := h.uc.FindTelegramById(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler method FindTelegramById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.Profile{
			ID:             p.ID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			SearchGender:   p.SearchGender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
			LookingFor:     p.LookingFor,
			IsDeleted:      p.IsDeleted,
			IsBlocked:      p.IsBlocked,
			IsPremium:      p.IsPremium,
			IsShowDistance: p.IsShowDistance,
			IsInvisible:    p.IsInvisible,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			LastOnline:     p.LastOnline,
			Images:         i,
			Telegram:       t,
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerProfile) DeleteProfileImageHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/image/delete")
		req := profile.RequestDeleteProfileImage{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func DeleteProfileImageHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug(
				"error func DeleteProfileImageHandler, method ParseUint roomIdStr by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageInDB, err := h.uc.FindImageById(ctf.Context(), imageID)
		if err != nil {
			h.logger.Debug("error func DeleteProfileImageHandler, method FindImageById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		if imageInDB.IsDeleted == true {
			msg := errors.Wrap(err, "image has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filePath := imageInDB.Url
		if err := os.Remove(filePath); err != nil {
			h.logger.Debug("error func DeleteProfileImageHandler, method Remove by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageDTO := &profile.ImageProfile{
			ID:        imageInDB.ID,
			ProfileID: imageInDB.ProfileID,
			Name:      "",
			Url:       "",
			Size:      0,
			CreatedAt: imageInDB.CreatedAt,
			UpdatedAt: time.Now(),
			IsDeleted: true,
			IsBlocked: imageInDB.IsBlocked,
			IsPrimary: imageInDB.IsPrimary,
			IsPrivate: imageInDB.IsPrivate,
		}
		response, err := h.uc.DeleteImage(ctf.Context(), imageDTO)
		if err != nil {
			h.logger.Debug("error func DeleteProfileImageHandler, method DeleteImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, response)
	}
}
