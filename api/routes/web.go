package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetUpWebRoutes(api fiber.Router) {
	api.Get("/", func(c *fiber.Ctx) error { // Middleware for /api/v1
		return c.JSON(fiber.Map{
			"message": "Backend REST API Starter Kit : V1",
			"status":  "Ok",
		})
	})
	api.Static("/categories", "./api/static/categories.html")
	api.Static("/product-lists", "./api/static/product-lists.html")
	api.Static("/product-details", "./api/static/product-details.html")
}
