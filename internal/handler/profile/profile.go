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

func (h *HandlerProfile) CreateProfileHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		h.logger.Info("POST /api/v1/profile/create")
		req := profileUseCase.CreateRequestProfile{}
		if err := ctf.BodyParser(&req); err != nil {
			h.logger.Debug(
				"error func CreateProfileHandler,"+
					" method BodyParser by path internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.Create(ctf.Context(), &req)
		if err != nil {
			h.logger.Debug(
				"error func CreateProfileHandler, method Create by path"+
					" internal/handler/profile/profile.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, response)
	}
}
