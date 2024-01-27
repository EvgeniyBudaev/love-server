package app

import (
	profileRepo "github.com/EvgeniyBudaev/love-server/internal/adapter/psqlRepo/profile"
	"github.com/EvgeniyBudaev/love-server/internal/handler/profile"
	profileUseCase "github.com/EvgeniyBudaev/love-server/internal/useCase/profile"
	"go.uber.org/zap"
)

var prefix = "/api/v1"

func (app *App) StartHTTPServer() error {
	app.fiber.Static("/static", "./static")
	pr := profileRepo.NewRepositoryProfile(app.Logger, app.db.psql)
	puc := profileUseCase.NewUseCaseProfile(pr)
	ph := profile.NewHandlerProfile(app.Logger, puc)
	grp := app.fiber.Group(prefix)
	grp.Post("/profile/create", ph.CreateProfileHandler())
	if err := app.fiber.Listen(app.config.Port); err != nil {
		app.Logger.Fatal("error func StartHTTPServer, method Listen by path internal/app/http.go", zap.Error(err))
	}
	return nil
}
