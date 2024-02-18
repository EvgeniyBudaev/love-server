package user

import (
	r "github.com/EvgeniyBudaev/love-server/internal/handler/http/api/v1/response"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"github.com/EvgeniyBudaev/love-server/internal/useCase/user"
	userUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/user"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
)

type HandlerUser struct {
	logger logger.Logger
	uc     *userUseCase.UseCaseUser
}

func NewHandlerUser(l logger.Logger, uc *userUseCase.UseCaseUser) *HandlerUser {
	return &HandlerUser{logger: l, uc: uc}
}

func (h *HandlerUser) PostRegisterHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		var ctx = ctf.UserContext()
		h.logger.Info("POST /api/v1/user/register")
		var request = user.RegisterRequest{}
		err := ctf.BodyParser(&request)
		if err != nil {
			h.logger.Debug("error func PostRegisterHandler, method BodyParser by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.Register(ctx, request)
		if err != nil {
			h.logger.Debug("error func PostRegisterHandler, method Register by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerUser) UpdateUserHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		var ctx = ctf.UserContext()
		h.logger.Info("POST /api/v1/user/update")
		var request = user.RequestUpdateUser{}
		err := ctf.BodyParser(&request)
		if err != nil {
			h.logger.Debug("error func UpdateUserHandle, method BodyParser by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		response, err := h.uc.UpdateUser(ctx, request)
		if err != nil {
			h.logger.Debug("error func UpdateUserHandle, method UpdateUser by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, response)
	}
}

func (h *HandlerUser) DeleteUserHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		var ctx = ctf.UserContext()
		h.logger.Info("POST /api/v1/user/delete")
		var request = user.RequestDeleteUser{}
		err := ctf.BodyParser(&request)
		if err != nil {
			h.logger.Debug("error func DeleteUserHandler, method BodyParser by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		err = h.uc.DeleteUser(ctx, request)
		if err != nil {
			h.logger.Debug("error func DeleteUserHandler, method DeleteUser by path internal/handler/user/user.go",
				zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapCreated(ctf, nil)
	}
}

func (h *HandlerUser) GetUserListHandler() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		var ctx = ctf.UserContext()
		h.logger.Info("GET /api/v1/user/list")
		query := user.QueryParamsUserList{}
		if err := ctf.QueryParser(&query); err != nil {
			h.logger.Debug("error func GetUserListHandler, method QueryParser by path"+
				" internal/handler/user/user.go", zap.Error(err))
			return err
		}
		response, err := h.uc.GetUserList(ctx, query)
		if err != nil {
			h.logger.Debug("error func GetUserListHandler, method GetUserList by path"+
				" internal/handler/user/user.go", zap.Error(err))
			return r.WrapError(ctf, err, http.StatusBadRequest)
		}
		return r.WrapOk(ctf, response)
	}
}
