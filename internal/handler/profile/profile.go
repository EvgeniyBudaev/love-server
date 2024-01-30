package profile

import (
	profileUseCase "github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	r "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/response"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/EvgeniyBudaev/love-server/internal/useCase/profile"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
)

type HandlerProfile struct {
	logger logger.Logger
	uc     *profile.UseCaseProfile
}

func NewHandlerProfile(l logger.Logger, uc *profile.UseCaseProfile) *HandlerProfile {
	return &HandlerProfile{logger: l, uc: uc}
}

func (h *HandlerProfile) AddProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/add")
		req := profileUseCase.AddRequestProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func AddProfileHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.Add(ctf, &req)
		if err != nil {
			h.logger.Debug(
				"error func AddProfileHandler, method Add by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerProfile) GetProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("GET /api/v1/profile/:id")
		response, err := h.uc.FindById(ctf)
		if err != nil {
			h.logger.Debug(
				"error func GetProfileHandler, method FindById by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}
