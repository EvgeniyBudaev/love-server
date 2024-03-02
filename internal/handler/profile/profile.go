package profile

import (
	"fmt"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	errorDomain "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/error"
	r "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/response"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	profileUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"image/jpeg"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
			h.logger.Debug("error func AddProfileHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filePath := fmt.Sprintf("static/uploads/profile/%s/images/defaultImage.jpg", req.UserName)
		directoryPath := fmt.Sprintf("static/uploads/profile/%s/images", req.UserName)
		if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
			if err := os.MkdirAll(directoryPath, 0755); err != nil {
				h.logger.Debug("error func AddProfileHandler, method MkdirAll by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		form, err := ctf.MultipartForm()
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method MultipartForm by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageFiles := form.File["image"]
		imagesFilePath := make([]string, 0, len(imageFiles))
		imagesProfile := make([]*profile.ImageProfile, 0, len(imagesFilePath))
		for _, file := range imageFiles {
			filePath = fmt.Sprintf("%s/%s", directoryPath, file.Filename)
			if err := ctf.SaveFile(file, filePath); err != nil {
				h.logger.Debug("error func AddProfileHandler, method SaveFile by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			fileImage, err := os.Open(filePath)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method os.Open by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			// The Decode function is used to read images from a file or other source and convert them into an image.
			// Image structure
			img, err := jpeg.Decode(fileImage)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method jpeg.Decode by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			newFileName := replaceExtension(file.Filename)
			newFilePath := fmt.Sprintf("%s/%s", directoryPath, newFileName)
			output, err := os.Create(directoryPath + "/" + newFileName)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method os.Create by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			defer output.Close()
			options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method NewLossyEncoderOptions by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			if err := webp.Encode(output, img, options); err != nil {
				h.logger.Debug("error func AddProfileHandler, method webp.Encode by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			if err := os.Remove(filePath); err != nil {
				h.logger.Debug("error func AddProfileHandler, method os.Remove by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			image := profile.ImageProfile{
				Name:      file.Filename,
				Url:       newFilePath,
				Size:      file.Size,
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				IsDeleted: false,
				IsBlocked: false,
				IsPrimary: false,
				IsPrivate: false,
			}
			imagesFilePath = append(imagesFilePath, newFilePath)
			imagesProfile = append(imagesProfile, &image)
		}
		height := 0
		if req.Height != "" {
			heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			height = int(heightUint64)
		}
		weight := 0
		if req.Weight != "" {
			weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			weight = int(weightUint64)
		}
		profileDto := &profile.Profile{
			SessionID:      req.SessionID,
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
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
			LastOnline:     time.Now().UTC(),
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
				h.logger.Debug("error func AddProfileHandler, method AddImage by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		telegramID, err := strconv.ParseUint(req.TelegramID, 10, 64)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method ParseUint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		allowsWriteToPm, err := strconv.ParseBool(req.AllowsWriteToPm)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method ParseBool by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		telegramDto := &profile.TelegramProfile{
			ProfileID:       newProfile.ID,
			TelegramID:      telegramID,
			UserName:        req.TelegramUserName,
			Firstname:       req.Firstname,
			Lastname:        req.Lastname,
			LanguageCode:    req.LanguageCode,
			AllowsWriteToPm: allowsWriteToPm,
			QueryID:         req.QueryID,
		}
		_, err = h.uc.AddTelegram(ctf.Context(), telegramDto)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method AddTelegram by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		ageFrom := 0
		if req.AgeFrom != "" {
			ageFromUint8, err := strconv.ParseUint(req.AgeFrom, 10, 8)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageFrom = int(ageFromUint8)
		}
		ageTo := 0
		if req.AgeTo != "" {
			ageToUint8, err := strconv.ParseUint(req.AgeTo, 10, 8)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageTo = int(ageToUint8)
		}
		distance := 0
		if req.Distance != "" {
			distance32, err := strconv.ParseUint(req.Distance, 10, 64)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			distance = int(distance32)
		}
		page := 0
		if req.Page != "" {
			page32, err := strconv.ParseUint(req.Page, 10, 64)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			page = int(page32)
		}
		size := 0
		if req.Size != "" {
			size32, err := strconv.ParseUint(req.Size, 10, 64)
			if err != nil {
				h.logger.Debug("error func AddProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			size = int(size32)
		}
		filterDto := &profile.FilterProfile{
			ProfileID:    newProfile.ID,
			SearchGender: req.SearchGender,
			LookingFor:   req.LookingFor,
			AgeFrom:      uint8(ageFrom),
			AgeTo:        uint8(ageTo),
			Distance:     uint64(distance),
			Page:         uint64(page),
			Size:         uint64(size),
		}
		_, err = h.uc.AddFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method AddFilter by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitude, err := strconv.ParseFloat(req.Latitude, 64)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method ParseFloat height by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		longitude, err := strconv.ParseFloat(req.Longitude, 64)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method ParseFloat height by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			h.logger.Debug("error func AddProfileHandler, method AddNavigator by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), newProfile.ID)
		if err != nil {
			h.logger.Debug("error func AddProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			SessionID:      p.SessionID,
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
			h.logger.Debug("error func GetProfileListHandler, method QueryParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindBySessionID(ctf.Context(), params.SessionID)
		if err != nil {
			h.logger.Debug("error func GetProfileListHandler, method FindBySessionID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileListHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		fmt.Println("List latitudeStr:", latitudeStr)
		fmt.Println("List longitudeStr:", longitudeStr)
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileBySessionIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileBySessionIDHandler, method ParseFloat height by path"+
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
				h.logger.Debug("error func GetProfileBySessionIDHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileListHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		ageFrom := 0
		if params.AgeFrom != "" {
			ageFromUint8, err := strconv.ParseUint(params.AgeFrom, 10, 8)
			if err != nil {
				h.logger.Debug("error func GetProfileListHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageFrom = int(ageFromUint8)
		}
		ageTo := 0
		if params.AgeTo != "" {
			ageToUint8, err := strconv.ParseUint(params.AgeTo, 10, 8)
			if err != nil {
				h.logger.Debug("error func GetProfileListHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			ageTo = int(ageToUint8)
		}
		distance := 0
		if params.Distance != "" {
			distance32, err := strconv.ParseUint(params.Distance, 10, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileListHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			distance = int(distance32)
		}
		filterDto := &profile.FilterProfile{
			ID:           f.ID,
			ProfileID:    p.ID,
			SearchGender: params.SearchGender,
			LookingFor:   params.LookingFor,
			AgeFrom:      uint8(ageFrom),
			AgeTo:        uint8(ageTo),
			Distance:     uint64(distance),
			Page:         params.Page,
			Size:         params.Size,
		}
		_, err = h.uc.UpdateFilter(ctf.Context(), filterDto)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method UpdateFilter by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.SelectList(ctf.Context(), &params)
		if err != nil {
			h.logger.Debug("error func GetProfileListHandler, method SelectList by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileBySessionIDHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/session/:id")
		sessionID := ctf.Params("id")
		params := profile.QueryParamsGetProfileByUserID{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method QueryParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindBySessionID(ctf.Context(), sessionID)
		if err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method FindBySessionID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		fmt.Println("Get session latitudeStr:", latitudeStr)
		fmt.Println("Get session longitudeStr:", longitudeStr)
		if latitudeStr != "" && longitudeStr != "" {
			latitude, err := strconv.ParseFloat(latitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileBySessionIDHandler, method ParseFloat height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			longitude, err := strconv.ParseFloat(longitudeStr, 64)
			if err != nil {
				h.logger.Debug("error func GetProfileBySessionIDHandler, method ParseFloat height by path"+
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
				h.logger.Debug("error func GetProfileBySessionIDHandler, method UpdateNavigator by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		t, err := h.uc.FindTelegramByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method FindTelegramByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		f, err := h.uc.FindFilterByProfileID(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method FindFilterByProfileID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		i, err := h.uc.SelectListPublicImage(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileBySessionIDHandler, method SelectListPublicImage by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response := &profile.ResponseProfile{
			ID:        p.ID,
			SessionID: p.SessionID,
			IsDeleted: p.IsDeleted,
			IsBlocked: p.IsBlocked,
			Image:     nil,
			Telegram:  &profile.ResponseTelegramProfile{TelegramID: t.TelegramID},
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
			h.logger.Debug("error func GetProfileDetailHandler, method QueryParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		v, err := h.uc.FindBySessionID(ctf.Context(), params.ViewerID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method FindBySessionID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), v.ID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		latitudeStr := params.Latitude
		longitudeStr := params.Longitude
		fmt.Println("Detail latitudeStr:", latitudeStr)
		fmt.Println("Detail longitudeStr:", longitudeStr)
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
				ProfileID: v.ID,
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
		l, isExistLike, err := h.uc.FindLikeByHumanID(ctf.Context(), v.ID, profileID)
		if err != nil {
			h.logger.Debug("error func GetProfileDetailHandler, FindLikeByHumanID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		var lDao *profile.ResponseLikeProfile
		if isExistLike {
			lDao = &profile.ResponseLikeProfile{
				ID: func() *uint64 {
					if isExistLike {
						return &l.ID
					}
					return nil
				}(),
				IsLiked: isExistLike && l.IsLiked,
				UpdatedAt: func() *time.Time {
					if isExistLike {
						return &l.UpdatedAt
					}
					return nil
				}(),
			}
		}
		response := &profile.ResponseProfileDetail{
			ID:             p.ID,
			SessionID:      p.SessionID,
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
			IsOnline:       false,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
			LastOnline:     p.LastOnline,
			Images:         i,
			Telegram:       t,
			Navigator:      n,
			Filter:         f,
			Like:           lDao,
		}
		elapsed := time.Since(p.LastOnline)
		if elapsed.Minutes() < 5 {
			response.IsOnline = true
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) UpdateProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/edit")
		req := profile.RequestUpdateProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			h.logger.Debug("error func UpdateProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
		err = h.uc.UpdateLastOnline(ctf.Context(), profileInDB.ID)
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		filePath := fmt.Sprintf("static/uploads/profile/%s/images/defaultImage.jpg", req.UserName)
		directoryPath := fmt.Sprintf("static/uploads/profile/%s/images", req.UserName)
		if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
			if err := os.MkdirAll(directoryPath, 0755); err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method MkdirAll by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		form, err := ctf.MultipartForm()
		if err != nil {
			h.logger.Debug("error func UpdateProfileHandler, method MultipartForm by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		height := 0
		if req.Height != "" {
			heightUint64, err := strconv.ParseUint(req.Height, 10, 8)
			if err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			height = int(heightUint64)
		}
		weight := 0
		if req.Weight != "" {
			weightUint64, err := strconv.ParseUint(req.Weight, 10, 8)
			if err != nil {
				h.logger.Debug("error func UpdateProfileHandler, method ParseUint height by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
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
					h.logger.Debug("error func UpdateProfileHandler, method SaveFile by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				fileImage, err := os.Open(filePath)
				if err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method os.Open by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				// The Decode function is used to read images from a file or other source and convert them into an image.
				// Image structure
				img, err := jpeg.Decode(fileImage)
				if err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method jpeg.Decode by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				newFileName := replaceExtension(file.Filename)
				newFilePath := fmt.Sprintf("%s/%s", directoryPath, newFileName)
				output, err := os.Create(directoryPath + "/" + newFileName)
				if err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method os.Create by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				defer output.Close()
				options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
				if err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method NewLossyEncoderOptions by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				if err := webp.Encode(output, img, options); err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method webp.Encode by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				if err := os.Remove(filePath); err != nil {
					h.logger.Debug("error func UpdateProfileHandler, method os.Remove by path"+
						" internal/handler/profile/profile.go", zap.Error(err))
					return r.WrapError(ctf, err, http.StatusBadRequest)
				}
				image := profile.ImageProfile{
					Name:      file.Filename,
					Url:       newFilePath,
					Size:      file.Size,
					CreatedAt: time.Now().UTC(),
					UpdatedAt: time.Now().UTC(),
					IsDeleted: false,
					IsBlocked: false,
					IsPrimary: false,
					IsPrivate: false,
				}
				imagesFilePath = append(imagesFilePath, newFilePath)
				imagesProfile = append(imagesProfile, &image)
			}
			profileDto = &profile.Profile{
				ID:             profileID,
				SessionID:      profileInDB.SessionID,
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
				UpdatedAt:      time.Now().UTC(),
				LastOnline:     time.Now().UTC(),
				Images:         imagesProfile,
			}
		} else {
			profileDto = &profile.Profile{
				ID:             profileID,
				SessionID:      profileInDB.SessionID,
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
				UpdatedAt:      time.Now().UTC(),
				LastOnline:     time.Now().UTC(),
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
			UserName:        req.TelegramUserName,
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
			h.logger.Debug("error func UpdateProfileHandler, method UpdateTelegram by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			h.logger.Debug("error func UpdateProfileHandler, method UpdateFilter by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			h.logger.Debug("error func UpdateProfileHandler, method UpdateNavigator by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			SessionID:      p.SessionID,
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
			h.logger.Debug("error func DeleteProfileHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method ParseUint roomIdStr by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileInDB, err := h.uc.FindById(ctf.Context(), profileID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if profileInDB.IsDeleted == true {
			msg := errors.Wrap(err, "user has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), profileInDB.ID)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
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
					UpdatedAt: time.Now().UTC(),
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
			h.logger.Debug("error func DeleteProfileHandler, method DeleteTelegram by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
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
			SessionID:      "",
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
			UpdatedAt:      time.Now().UTC(),
			LastOnline:     time.Now().UTC(),
		}
		_, err = h.uc.Delete(ctf.Context(), profileDto)
		if err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method Delete by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindById(ctf.Context(), profileID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func DeleteProfileHandler, method FindById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		response := &profile.Profile{
			ID:             p.ID,
			SessionID:      p.SessionID,
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
			h.logger.Debug("error func DeleteProfileImageHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		imageID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func DeleteProfileImageHandler, method ParseUint roomIdStr by path"+
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
			UpdatedAt: time.Now().UTC(),
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

func (h *HandlerProfile) AddReviewHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/review/add")
		req := profile.RequestAddReview{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func AddReviewHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(req.ProfileID, 10, 64)
		if err != nil {
			h.logger.Debug("error func AddReviewHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func AddReviewHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		rating, err := strconv.ParseFloat(req.Rating, 32)
		if err != nil {
			h.logger.Debug("error func AddReviewHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewDto := &profile.ReviewProfile{
			ProfileID:  profileID,
			Message:    req.Message,
			Rating:     float32(rating),
			HasDeleted: false,
			HasEdited:  false,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		review, err := h.uc.AddReview(ctf.Context(), reviewDto)
		if err != nil {
			h.logger.Debug("error func AddReviewHandler, method AddReview by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, review)
	}
}

func (h *HandlerProfile) UpdateReviewHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/review/update")
		req := profile.RequestUpdateReview{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(req.ProfileID, 10, 64)
		if err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewInDB, err := h.uc.FindReviewById(ctf.Context(), reviewID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method FindReviewById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if reviewInDB.HasDeleted == true {
			msg := errors.Wrap(err, "review has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		rating, err := strconv.ParseFloat(req.Rating, 32)
		if err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewDto := &profile.ReviewProfile{
			ID:         reviewID,
			ProfileID:  profileID,
			Message:    req.Message,
			Rating:     float32(rating),
			HasDeleted: reviewInDB.HasDeleted,
			HasEdited:  true,
			CreatedAt:  reviewInDB.CreatedAt,
			UpdatedAt:  time.Now().UTC(),
		}
		review, err := h.uc.UpdateReview(ctf.Context(), reviewDto)
		if err != nil {
			h.logger.Debug("error func UpdateReviewHandler, method UpdateReview by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, review)
	}
}

func (h *HandlerProfile) DeleteReviewHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/review/delete")
		req := profile.RequestDeleteReview{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func DeleteReviewHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func DeleteReviewHandler, method ParseUint roomIdStr by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		reviewInDB, err := h.uc.FindReviewById(ctf.Context(), reviewID)
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func DeleteReviewHandler, method FindReviewById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		if reviewInDB.HasDeleted == true {
			msg := errors.Wrap(err, "review has already been deleted")
			err = errorDomain.NewCustomError(msg, http.StatusNotFound)
			return r.WrapError(ctf, err, http.StatusNotFound)
		}
		reviewDto := &profile.ReviewProfile{
			ID:         reviewID,
			ProfileID:  reviewInDB.ProfileID,
			Message:    reviewInDB.Message,
			Rating:     reviewInDB.Rating,
			HasDeleted: true,
			HasEdited:  reviewInDB.HasEdited,
			CreatedAt:  reviewInDB.CreatedAt,
			UpdatedAt:  time.Now().UTC(),
		}
		review, err := h.uc.DeleteReview(ctf.Context(), reviewDto)
		if err != nil {
			h.logger.Debug("error func DeleteReviewHandler, method UpdateReview by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, review)
	}
}

func (h *HandlerProfile) GetReviewByIDHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/review/detail/:id")
		idStr := ctf.Params("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetReviewByIDHandler, method ParseUint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.FindReviewById(ctf.Context(), id)
		if err != nil {
			h.logger.Debug("error func GetProfileByIDHandler, method FindReviewById by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) GetReviewListHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/review/list")
		params := profile.QueryParamsReviewList{}
		if err := ctf.QueryParser(&params); err != nil {
			h.logger.Debug("error func GetReviewListHandler, method QueryParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		profileID, err := strconv.ParseUint(params.ProfileID, 10, 64)
		if err != nil {
			h.logger.Debug("error func GetReviewListHandler, method ParseUint roomIdStr by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), profileID)
		if err != nil {
			h.logger.Debug("error func GetReviewListHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.SelectReviewList(ctf.Context(), &params)
		if err != nil {
			h.logger.Debug("error func GetReviewListHandler, method SelectList by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}

func (h *HandlerProfile) AddLikeHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/like/add")
		req := profile.RequestAddLike{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func AddLikeHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		humanID, err := strconv.ParseUint(req.HumanID, 10, 64)
		if err != nil {
			h.logger.Debug("error func AddLikeHandler, method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindBySessionID(ctf.Context(), req.SessionID)
		if err != nil {
			h.logger.Debug("error func AddLikeHandler, method FindByKeycloakID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddLikeHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		likeDto := &profile.LikeProfile{
			ProfileID: p.ID,
			HumanID:   humanID,
			IsLiked:   true,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		like, err := h.uc.AddLike(ctf.Context(), likeDto)
		if err != nil {
			h.logger.Debug("error func AddLikeHandler, method AddLike by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, like)
	}
}

func (h *HandlerProfile) DeleteLikeHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/like/delete")
		req := profile.RequestDeleteLike{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func DeleteLikeHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		likeID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func DeleteLikeHandler, method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		l, isExistLike, err := h.uc.FindLikeByID(ctf.Context(), likeID)
		if err != nil {
			h.logger.Debug("error func DeleteLikeHandler, method FindByKeycloakID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		if !isExistLike {
			h.logger.Debug("error func DeleteLikeHandler, method !isExistLike by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			msg := errorDomain.ResponseError{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "not found",
			}
			return ctf.Status(http.StatusNotFound).JSON(msg)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), l.ProfileID)
		if err != nil {
			h.logger.Debug("error func DeleteLikeHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		likeDto := &profile.LikeProfile{
			ID:        likeID,
			ProfileID: l.ProfileID,
			HumanID:   l.HumanID,
			IsLiked:   false,
			CreatedAt: l.CreatedAt,
			UpdatedAt: time.Now().UTC(),
		}
		like, err := h.uc.DeleteLike(ctf.Context(), likeDto)
		if err != nil {
			h.logger.Debug("error func DeleteLikeHandler, method DeleteLike by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, like)
	}
}

func (h *HandlerProfile) UpdateLikeHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/like/update")
		req := profile.RequestUpdateLike{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func UpdateLikeHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		likeID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func UpdateLikeHandler, method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		l, isExist, err := h.uc.FindLikeByID(ctf.Context(), likeID)
		if err != nil {
			h.logger.Debug("error func UpdateLikeHandler, method FindLikeByID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		if !isExist {
			h.logger.Debug("error func UpdateLikeHandler, method !isExist by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			msg := errorDomain.ResponseError{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "not found",
			}
			return ctf.Status(http.StatusNotFound).JSON(msg)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), l.ProfileID)
		if err != nil {
			h.logger.Debug("error func UpdateLikeHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		likeDto := &profile.LikeProfile{
			ID:        likeID,
			ProfileID: l.ProfileID,
			HumanID:   l.HumanID,
			IsLiked:   true,
			CreatedAt: l.CreatedAt,
			UpdatedAt: time.Now().UTC(),
		}
		like, err := h.uc.UpdateLike(ctf.Context(), likeDto)
		if err != nil {
			h.logger.Debug("error func UpdateLikeHandler, method UpdateLike by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, like)
	}
}

func (h *HandlerProfile) AddBlockHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/block/add")
		req := profile.RequestAddBlock{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func AddBlockHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockedUserID, err := strconv.ParseUint(req.BlockedUserID, 10, 64)
		if err != nil {
			h.logger.Debug("error func AddBlockHandler, method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindBySessionID(ctf.Context(), req.SessionID)
		if err != nil {
			h.logger.Debug("error func AddBlockHandler, method FindBySessionID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddBlockHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockDto := &profile.BlockedProfile{
			ProfileID:     p.ID,
			BlockedUserID: blockedUserID,
			IsBlocked:     true,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		block, err := h.uc.AddBlock(ctf.Context(), blockDto)
		if err != nil {
			h.logger.Debug("error func AddBlockHandler, method AddBlock by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockForBlockedUserDto := &profile.BlockedProfile{
			ProfileID:     blockedUserID,
			BlockedUserID: p.ID,
			IsBlocked:     true,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		_, err = h.uc.AddBlock(ctf.Context(), blockForBlockedUserDto)
		if err != nil {
			h.logger.Debug("error func AddBlockHandler, method AddBlock by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, block)
	}
}

func (h *HandlerProfile) UpdateBlockHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/block/update")
		req := profile.RequestUpdateBlock{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func UpdateBlockHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockID, err := strconv.ParseUint(req.ID, 10, 64)
		if err != nil {
			h.logger.Debug("error func UpdateBlockHandler method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		b, isExist, err := h.uc.FindBlockByID(ctf.Context(), blockID)
		if err != nil {
			h.logger.Debug("error func UpdateBlockHandler, method FindBlockByID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		if !isExist {
			h.logger.Debug("error func UpdateBlockHandler, method !isExist by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			msg := errorDomain.ResponseError{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "not found",
			}
			return ctf.Status(http.StatusNotFound).JSON(msg)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), b.ProfileID)
		if err != nil {
			h.logger.Debug("error func UpdateBlockHandler, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockDto := &profile.BlockedProfile{
			ID:            blockID,
			ProfileID:     b.ProfileID,
			BlockedUserID: blockID,
			IsBlocked:     false,
			CreatedAt:     b.CreatedAt,
			UpdatedAt:     time.Now().UTC(),
		}
		like, err := h.uc.UpdateBlock(ctf.Context(), blockDto)
		if err != nil {
			h.logger.Debug("error func UpdateBlockHandler, method UpdateBlock by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, like)
	}
}

func (h *HandlerProfile) AddComplaintHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/complaint/add")
		req := profile.RequestAddComplaint{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug("error func AddComplaintHandler, method BodyParser by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		complaintUserId, err := strconv.ParseUint(req.ComplaintUserID, 10, 64)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandler, method ParseUint by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		p, err := h.uc.FindBySessionID(ctf.Context(), req.SessionID)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandler, method FindBySessionID by path "+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.UpdateLastOnline(ctf.Context(), p.ID)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandle, method UpdateLastOnline by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		complaintDto := &profile.ComplaintProfile{
			ProfileID:       p.ID,
			ComplaintUserID: complaintUserId,
			Reason:          req.Reason,
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		}
		complaint, err := h.uc.AddComplaint(ctf.Context(), complaintDto)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandler, method AddComplaint by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		blockDto := &profile.BlockedProfile{
			ProfileID:     p.ID,
			BlockedUserID: complaintUserId,
			IsBlocked:     true,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		_, err = h.uc.AddBlock(ctf.Context(), blockDto)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandler, method AddBlock by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		listComplaint, err := h.uc.SelectListComplaintByID(ctf.Context(), complaintUserId)
		if err != nil {
			h.logger.Debug("error func AddComplaintHandler, method SelectListComplaintByID by path"+
				" internal/handler/profile/profile.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		if len(filterComplaintsByCurrentMonth(listComplaint)) > 1 {
			p, err := h.uc.FindById(ctf.Context(), complaintUserId)
			if err != nil {
				h.logger.Debug("error func AddComplaintHandler, method FindById by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
			profileDto := &profile.Profile{
				ID:             p.ID,
				SessionID:      p.SessionID,
				DisplayName:    p.DisplayName,
				Birthday:       p.Birthday,
				Gender:         p.Gender,
				Location:       p.Location,
				Description:    p.Description,
				Height:         p.Height,
				Weight:         p.Weight,
				IsDeleted:      p.IsDeleted,
				IsBlocked:      true,
				IsPremium:      p.IsPremium,
				IsShowDistance: p.IsShowDistance,
				IsInvisible:    p.IsInvisible,
				CreatedAt:      p.CreatedAt,
				UpdatedAt:      p.UpdatedAt,
				LastOnline:     p.LastOnline,
			}
			_, err = h.uc.Update(ctf.Context(), profileDto)
			if err != nil {
				h.logger.Debug("error func AddComplaintHandler, method Update by path"+
					" internal/handler/profile/profile.go", zap.Error(err))
				return r.WrapError(ctf, err, http.StatusBadRequest)
			}
		}
		return r.WrapCreated(ctf, complaint)
	}
}

// filterComplaintsByCurrentMonth возвращает кол-во жалоб за текущий месяц
func filterComplaintsByCurrentMonth(complaints []*profile.ComplaintProfile) []*profile.ComplaintProfile {
	currentMonth := time.Now().UTC().Month()
	currentYear := time.Now().UTC().Year()
	filteredComplaints := make([]*profile.ComplaintProfile, 0)
	for _, complaint := range complaints {
		if complaint.CreatedAt.Month() == currentMonth && complaint.CreatedAt.Year() == currentYear {
			filteredComplaints = append(filteredComplaints, complaint)
		}
	}
	return filteredComplaints
}

func replaceExtension(filename string) string {
	// Удаляем текущее расширение
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	// Добавляем новое расширение .webp
	return filename + ".webp"
}
