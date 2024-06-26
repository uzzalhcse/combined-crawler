package routes

import (
	"combined-crawler/api/app/http/controllers"
	"combined-crawler/api/app/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetUpAuthRoutes(api fiber.Router) {
	AuthController := controllers.NewAuthController()
	api.Post("/login", AuthController.Login)
	api.Post("/register", AuthController.Register)
	api.Post("/forget-password", AuthController.ForgetPassword)

	auth := api.Group("", middleware.Auth())
	auth.Get("/update-profile", AuthController.UpdateProfile)
	auth.Get("/me", AuthController.Me)

}
