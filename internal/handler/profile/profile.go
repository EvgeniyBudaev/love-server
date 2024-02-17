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
	"math"
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
		height := 0
		if req.Height != "" {
			heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			height = int(heightUint64)
		}
		weight := 0
		if req.Weight != "" {
			weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			weight = int(weightUint64)
		}
		profileDto := &profile.Profile{
			UserID:         req.UserID,
			DisplayName:    req.DisplayName,
			Birthday:       req.Birthday,
			Gender:         req.Gender,
			Location:       req.Location,
			Description:    req.Description,
			Height:         uint8(height),
			Weight:         uint8(weight),
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
				"error func AddProfileHandler, method ParseUint by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		allowsWriteToPm, err := strconv.ParseBool(req.AllowsWriteToPm)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseBool by path internal/handler/profile/profile.go",
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
		ageFrom := 0
		if req.AgeFrom != "" {
			ageFromUint32, err := strconv.ParseUint(req.AgeFrom, 10, 32)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageFrom = int(ageFromUint32)
		}
		ageTo := 0
		if req.AgeTo != "" {
			ageToUint32, err := strconv.ParseUint(req.AgeTo, 10, 32)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageTo = int(ageToUint32)
		}
		distance := 0
		if req.Distance != "" {
			distance32, err := strconv.ParseUint(req.Distance, 10, 64)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			distance = int(distance32)
		}
		page := 0
		if req.Page != "" {
			page32, err := strconv.ParseUint(req.Page, 10, 64)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			page = int(page32)
		}
		size := 0
		if req.Size != "" {
			size32, err := strconv.ParseUint(req.Size, 10, 64)
			if err != nil {
				h.logger.Debug(
					"error func AddProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			size = int(size32)
		}
		filterDto := &profile.FilterProfile{
			ProfileID:    newProfile.ID,
			SearchGender: req.SearchGender,
			LookingFor:   req.LookingFor,
			AgeFrom:      uint32(ageFrom),
			AgeTo:        uint32(ageTo),
			Distance:     uint64(distance),
			Page:         uint64(page),
			Size:         uint64(size),
		}
		_, err = h.uc.AddFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method AddFilter by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitude, err := strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseFloat height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		longitude, err := strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method ParseFloat height by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		point := &profile.Point{
			Latitude:  latitude,
			Longitude: longitude,
		}
		navigatorDto := &profile.NavigatorProfile{
			ProfileID: newProfile.ID,
			Location:  point,
		}
		_, err = h.uc.AddNavigator(ctf.Context(), navigatorDto)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method AddNavigator by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), newProfile.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindById by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindFilterByProfileID by path"+
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
			UserID:         p.UserID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
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
			Filter:         f,
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileListHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/list")
		params := profile.QueryParamsProfileList{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug(
				"error func GetProfileListHandler, method QueryParser by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(params.ProfileID, 10, 64)
		if err != nil {
			h.logger.Debug(
				"error func GetProfileListHandler, method ParseUint by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug(
				"error func GetProfileListHandler, method FindFilterByProfileID by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		ageFrom := 0
		if params.AgeFrom != "" {
			ageFromUint32, err := strconv.ParseUint(params.AgeFrom, 10, 32)
			if err != nil {
				h.logger.Debug(
					"error func GetProfileListHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageFrom = int(ageFromUint32)
		}
		ageTo := 0
		if params.AgeTo != "" {
			ageToUint32, err := strconv.ParseUint(params.AgeTo, 10, 32)
			if err != nil {
				h.logger.Debug(
					"error func GetProfileListHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageTo = int(ageToUint32)
		}
		distance := 0
		if params.Distance != "" {
			distance32, err := strconv.ParseUint(params.Distance, 10, 64)
			if err != nil {
				h.logger.Debug(
					"error func GetProfileListHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			distance = int(distance32)
		}
		filterDto := &profile.FilterProfile{
			ID:           f.ID,
			ProfileID:    profileID,
			SearchGender: params.SearchGender,
			AgeFrom:      uint32(ageFrom),
			AgeTo:        uint32(ageTo),
			Distance:     uint64(distance),
			Page:         params.Page,
			Size:         params.Size,
		}
		_, err = h.uc.UpdateFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method UpdateFilter by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileListHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.SelectList(ctf.Context(), &params)
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

func (h *HandlerProfile) GetProfileByTelegramIDHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/telegram/:id")
		idStr := ctf.Params("id")
		telegramID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method ParseUint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		params := profile.QueryParamsGetProfileByTelegramID{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug(
				"error func GetProfileByTelegramIDHandler, method QueryParser by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindByTelegramId(ctf.Context(), telegramID)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByTelegramIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByTelegramIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			point := &profile.Point{
				Latitude:  latitude,
				Longitude: longitude,
			}
			navigatorDto := &profile.NavigatorProfile{
				ProfileID: p.ID,
				Location:  point,
			}
			_, err = h.uc.UpdateNavigator(ctf.Context(), navigatorDto)
			if err != nil {
				h.logger.Debug("error func GetProfileByTelegramIDHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByTelegramIDHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.ResponseProfile{
			ID:       p.ID,
			UserID:   p.UserID,
			Image:    nil,
			Telegram: &profile.ResponseTelegramProfile{TelegramID: t.TelegramID},
			Filter: &profile.ResponseFilterProfile{
				ID:           f.ID,
				SearchGender: f.SearchGender,
				LookingFor:   f.LookingFor,
				AgeFrom:      f.AgeFrom,
				AgeTo:        f.AgeTo,
				Distance:     f.Distance,
				Page:         f.Page,
				Size:         f.Size,
			},
		}
		if len(i) > 0 {
			i := profile.ResponseImageProfile{
				Url: i[0].Url,
			}
			response.Image = &i
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileByUserIDHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/keycloak/:id")
		userID := ctf.Params("id")
		params := profile.QueryParamsGetProfileByUserID{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug(
				"error func GetProfileByUserIDHandler, method QueryParser by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindByKeycloakID(ctf.Context(), userID)
		if err != nil {
			h.logger.Debug("error func GetProfileByUserIDHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByUserIDHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByUserIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByUserIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			point := &profile.Point{
				Latitude:  latitude,
				Longitude: longitude,
			}
			navigatorDto := &profile.NavigatorProfile{
				ProfileID: p.ID,
				Location:  point,
			}
			_, err = h.uc.UpdateNavigator(ctf.Context(), navigatorDto)
			if err != nil {
				h.logger.Debug("error func GetProfileByUserIDHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByUserIDHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByUserIDHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByUserIDHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.ResponseProfile{
			ID:       p.ID,
			UserID:   p.UserID,
			Image:    nil,
			Telegram: &profile.ResponseTelegramProfile{TelegramID: t.TelegramID},
			Filter: &profile.ResponseFilterProfile{
				ID:           f.ID,
				SearchGender: f.SearchGender,
				LookingFor:   f.LookingFor,
				AgeFrom:      f.AgeFrom,
				AgeTo:        f.AgeTo,
				Distance:     f.Distance,
				Page:         f.Page,
				Size:         f.Size,
			},
		}
		if len(i) > 0 {
			i := profile.ResponseImageProfile{
				Url: i[0].Url,
			}
			response.Image = &i
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileByIDHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/:id")
		idStr := ctf.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method ParseUint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		params := profile.QueryParamsGetProfileByID{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug(
				"error func GetProfileByIDHandler, method QueryParser by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileByIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			point := &profile.Point{
				Latitude:  latitude,
				Longitude: longitude,
			}
			navigatorDto := &profile.NavigatorProfile{
				ProfileID: p.ID,
				Location:  point,
			}
			_, err = h.uc.UpdateNavigator(ctf.Context(), navigatorDto)
			if err != nil {
				h.logger.Debug("error func GetProfileByIDHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.ResponseProfile{
			ID:       p.ID,
			UserID:   p.UserID,
			Image:    nil,
			Telegram: &profile.ResponseTelegramProfile{TelegramID: t.TelegramID},
			Filter: &profile.ResponseFilterProfile{
				ID:           f.ID,
				SearchGender: f.SearchGender,
				LookingFor:   f.LookingFor,
				AgeFrom:      f.AgeFrom,
				AgeTo:        f.AgeTo,
				Distance:     f.Distance,
				Page:         f.Page,
				Size:         f.Size,
			},
		}
		if len(i) > 0 {
			i := profile.ResponseImageProfile{
				Url: i[0].Url,
			}
			response.Image = &i
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileDetailHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/detail/:id")
		idStr := ctf.Params("id")
		profileID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method ParseUint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		params := profile.QueryParamsGetProfileDetail{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug(
				"error func GetProfileDetailHandler, method QueryParser by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		v, err := h.uc.FindByKeycloakID(ctf.Context(), params.ViewerID)
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindByKeycloakID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileDetailHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileDetailHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			point := &profile.Point{
				Latitude:  latitude,
				Longitude: longitude,
			}
			navigatorDto := &profile.NavigatorProfile{
				ProfileID: p.ID,
				Location:  point,
			}
			_, err = h.uc.UpdateNavigator(ctf.Context(), navigatorDto)
			if err != nil {
				h.logger.Debug("error func GetProfileDetailHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		n, err := h.uc.FindNavigatorByProfileIDAndViewerID(ctf.Context(), p.ID, v.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindNavigatorByProfileIDAndViewerID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.Profile{
			ID:             p.ID,
			UserID:         p.UserID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
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
			Navigator:      n,
			Filter:         f,
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
		height := 0
		if req.Height != "" {
			heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
			if err != nil {
				h.logger.Debug(
					"error func UpdateProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			height = int(heightUint64)
		}
		weight := 0
		if req.Weight != "" {
			weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
			if err != nil {
				h.logger.Debug(
					"error func UpdateProfileHandler, method ParseUint height by path internal/handler/profile/profile.go",
					zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			weight = int(weightUint64)
		}
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
				UserID:         profileInDB.UserID,
				DisplayName:    req.DisplayName,
				Birthday:       req.Birthday,
				Gender:         req.Gender,
				Location:       req.Location,
				Description:    req.Description,
				Height:         uint8(height),
				Weight:         uint8(weight),
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
				UserID:         profileInDB.UserID,
				DisplayName:    req.DisplayName,
				Birthday:       req.Birthday,
				Gender:         req.Gender,
				Location:       req.Location,
				Description:    req.Description,
				Height:         uint8(height),
				Weight:         uint8(weight),
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
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), profileUpdated.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		telegramDto := &profile.TelegramProfile{
			ID:              t.ID,
			ProfileID:       profileUpdated.ID,
			TelegramID:      telegramID,
			UserName:        req.UserName,
			Firstname:       req.Firstname,
			Lastname:        req.Lastname,
			LanguageCode:    req.LanguageCode,
			AllowsWriteToPm: allowsWriteToPm,
			QueryID:         req.QueryID,
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), profileUpdated.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		_, err = h.uc.UpdateTelegram(ctf.Context(), telegramDto)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method UpdateTelegram by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filterDto := &profile.FilterProfile{
			ID:           f.ID,
			ProfileID:    profileUpdated.ID,
			SearchGender: req.SearchGender,
			LookingFor:   req.LookingFor,
		}
		_, err = h.uc.UpdateFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method UpdateFilter by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := req.Latitude
		longitudeStr := req.Longitude
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			point := &profile.Point{
				Latitude:  latitude,
				Longitude: longitude,
			}
			navigatorDto := &profile.NavigatorProfile{
				ProfileID: profileID,
				Location:  point,
			}
			_, err = h.uc.UpdateNavigator(ctf.Context(), navigatorDto)
			if err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		if err != nil {
			h.logger.Debug(
				"error func UpdateProfileHandler, method UpdateNavigator by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileUpdated.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		t, err = h.uc.FindTelegramByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler method FindTelegramByProfileID by path"+
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
			UserID:         p.UserID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
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
			Filter:         f,
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerProfile) DeleteProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/delete")
		req := profile.RequestDeleteProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler, method ParseUint roomIdStr by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileInDB, err := h.uc.FindById(ctf.Context(), profileID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler,"+
					" method FindById by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if profileInDB.IsDeleted == true {
			msg := errors.Wrap(err, "user has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		imageList, err := h.uc.SelectListImage(ctf.Context(), profileID)
		if len(imageList) > 0 {
			for _, i := range imageList {
				filePath := i.Url
				if err := os.Remove(filePath); err != nil {
					h.logger.Debug("error func DeleteProfileHandler, method Remove by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				imageDTO := &profile.ImageProfile{
					ID:        i.ID,
					ProfileID: i.ProfileID,
					Name:      "",
					Url:       "",
					Size:      0,
					CreatedAt: i.CreatedAt,
					UpdatedAt: time.Now(),
					IsDeleted: true,
					IsBlocked: i.IsBlocked,
					IsPrimary: i.IsPrimary,
					IsPrivate: i.IsPrivate,
				}
				_, err := h.uc.DeleteImage(ctf.Context(), imageDTO)
				if err != nil {
					h.logger.Debug("error func DeleteProfileHandler, method DeleteImage by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), profileInDB.ID)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		telegramDto := &profile.TelegramProfile{
			ID:              t.ID,
			ProfileID:       profileInDB.ID,
			TelegramID:      0,
			UserName:        "",
			Firstname:       "",
			Lastname:        "",
			LanguageCode:    "",
			AllowsWriteToPm: false,
			QueryID:         "",
		}
		_, err = h.uc.DeleteTelegram(ctf.Context(), telegramDto)
		if err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler, method DeleteTelegram by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		n, err := h.uc.FindNavigatorByProfileID(ctf.Context(), profileInDB.ID)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method FindNavigatorByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		point := &profile.Point{
			Latitude:  0.0,
			Longitude: 0.0,
		}
		navigatorDto := &profile.NavigatorProfile{
			ID:        n.ID,
			ProfileID: profileInDB.ID,
			Location:  point,
		}
		_, err = h.uc.DeleteNavigator(ctf.Context(), navigatorDto)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method DeleteNavigator by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), profileInDB.ID)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filterDto := &profile.FilterProfile{
			ID:           f.ID,
			ProfileID:    profileInDB.ID,
			SearchGender: "",
			LookingFor:   "",
			AgeFrom:      0,
			AgeTo:        0,
			Distance:     0,
			Page:         0,
			Size:         0,
		}
		_, err = h.uc.DeleteFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method DeleteFilter by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileDto := &profile.Profile{
			ID:             profileID,
			UserID:         "",
			DisplayName:    "",
			Birthday:       profileInDB.Birthday,
			Gender:         "",
			Location:       "",
			Description:    "",
			Height:         0,
			Weight:         0,
			IsDeleted:      true,
			IsBlocked:      false,
			IsPremium:      false,
			IsShowDistance: false,
			IsInvisible:    false,
			CreatedAt:      profileInDB.CreatedAt,
			UpdatedAt:      time.Now(),
			LastOnline:     time.Now(),
		}
		_, err = h.uc.Delete(ctf.Context(), profileDto)
		if err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler, method Delete by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func DeleteProfileHandler,"+
					" method FindById by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		response := &profile.Profile{
			ID:             p.ID,
			UserID:         p.UserID,
			DisplayName:    p.DisplayName,
			Birthday:       p.Birthday,
			Gender:         p.Gender,
			Location:       p.Location,
			Description:    p.Description,
			Height:         p.Height,
			Weight:         p.Weight,
			IsDeleted:      p.IsDeleted,
			IsBlocked:      p.IsBlocked,
			IsPremium:      p.IsPremium,
			IsShowDistance: p.IsShowDistance,
			IsInvisible:    p.IsInvisible,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			LastOnline:     p.LastOnline,
			Images:         nil,
			Telegram:       nil,
			Navigator:      nil,
			Filter:         nil,
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

func (h *HandlerProfile) hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func (h *HandlerProfile) Distance(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, rad float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180
	rad = 6378100
	hs := h.hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*h.hsin(lo2-lo1)
	return 2 * rad * math.Asin(math.Sqrt(hs))
}
