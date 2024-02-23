package app

import (
	"github.com/EvgeniyBudaev/love-server/internal/handler/profile"
	"github.com/EvgeniyBudaev/love-server/internal/handler/user"
	"github.com/gofiber/fiber/v2"
)

func InitPublicRoutes(grp fiber.Router, imh *user.HandlerUser, ph *profile.HandlerProfile) {
	grp.Post("/user/register", imh.PostRegisterHandler())
	grp.Put("/user/update", imh.UpdateUserHandler())
	grp.Delete("/user/delete", imh.DeleteUserHandler())
	grp.Post("/profile/add", ph.AddProfileHandler())
	grp.Get("/profile/list", ph.GetProfileListHandler())
	grp.Get("/profile/telegram/:id", ph.GetProfileByTelegramIDHandler())
	grp.Get("/profile/keycloak/:id", ph.GetProfileByUserIDHandler())
	grp.Get("/profile/:id", ph.GetProfileByIDHandler())
	grp.Get("/profile/detail/:id", ph.GetProfileDetailHandler())

	grp.Post("/profile/edit", ph.UpdateProfileHandler())
	grp.Post("/profile/delete", ph.DeleteProfileHandler())
	grp.Post("/profile/image/delete", ph.DeleteProfileImageHandler())

	grp.Post("/review/add", ph.AddReviewHandler())
	grp.Post("/review/update", ph.UpdateReviewHandler())
	grp.Post("/review/delete", ph.DeleteReviewHandler())
	grp.Get("/review/list", ph.GetReviewListHandler())
	grp.Get("/review/detail/:id", ph.GetReviewByIDHandler())
}

func InitProtectedRoutes(grp fiber.Router, ph *profile.HandlerProfile) {
}
