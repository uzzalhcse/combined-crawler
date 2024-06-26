package responses

import "github.com/gofiber/fiber/v2"

// Success sends a successful JSON response
func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Error sends a JSON error response
func Error(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": message,
		//"errors":  nil,
	})
}
