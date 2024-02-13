package app

import (
	"github.com/EvgeniyBudaev/love-server/internal/handler/profile"
	"github.com/EvgeniyBudaev/love-server/internal/handler/user"
	"github.com/gofiber/fiber/v2"
)

func InitPublicRoutes(grp fiber.Router, imh *user.HandlerUser, ph *profile.HandlerProfile) {
	grp.Post("/user/register", imh.PostRegisterHandler())
	grp.Post("/profile/add", ph.AddProfileHandler())
	grp.Get("/profile/list", ph.GetProfileListHandler())
	grp.Get("/profile/telegram/:id", ph.GetProfileByTelegramIDHandler())
	grp.Get("/profile/:id", ph.GetProfileByIDHandler())
	grp.Get("/profile/detail/:id", ph.GetProfileDetailHandler())
}

func InitProtectedRoutes(grp fiber.Router, ph *profile.HandlerProfile) {
	grp.Post("/profile/edit", ph.UpdateProfileHandler())
	grp.Post("/profile/delete", ph.DeleteProfileHandler())
	grp.Post("/profile/image/delete", ph.DeleteProfileImageHandler())
}
