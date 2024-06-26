// middleware/auth.go

package middleware

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/bootstrap"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// Your secret key for signing and validating JWT tokens
var secretKey = []byte(bootstrap.App().Config.App.JwtSecret)

// Auth is a middleware for handling JWT token authentication
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract the token from the request header or query parameter
		token := c.Get("Authorization")

		// Verify the token
		user, err := verifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"data":  err.Error(),
			})
		}

		// Set the authenticated user into the request context
		c.Locals("user", user)

		// Continue to the next handler
		return c.Next()
	}
}

// verifyToken verifies the JWT token and returns the user claims
func verifyToken(tokenString string) (*models.User, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid token claims")
	}
	// Get the ID claim as a string
	userIDStr := getStringClaim(claims, "sub")
	if userIDStr == "" {
		return nil, fmt.Errorf("Empty user ID claim")
	}

	// Convert the ID claim to uint
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("Error converting user ID to int: %v", err)
	}
	return &models.User{
		ID:    uint(userID), // Convert to uint
		Name:  getStringClaim(claims, "name"),
		Email: getStringClaim(claims, "email"),
	}, nil
}

// Helper function to safely retrieve a string claim from the JWT claims
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
