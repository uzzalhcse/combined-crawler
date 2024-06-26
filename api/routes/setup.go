package routes

import "github.com/gofiber/fiber/v2"

// RegisterRoutes registers all routes
func RegisterRoutes(router fiber.Router) {
	web := router.Group("")
	SetUpWebRoutes(web)

	api := router.Group("/api")
	SetUpApiRoutes(api)

	auth := router.Group("/api/auth")
	SetUpAuthRoutes(auth)
}
