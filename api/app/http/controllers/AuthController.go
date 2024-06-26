// Controllers/auth_controller.go

package controllers

import (
	authrequests "combined-crawler/api/app/http/requests/auth"
	"combined-crawler/api/app/http/responses"
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
	"combined-crawler/api/app/services"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AuthController struct {
	BaseController *BaseController
	AuthService    services.AuthService
	JWTService     services.JWTService
}

func NewAuthController() *AuthController {
	that := NewBaseController()
	authRepo := repositories.NewAuthRepository(that.DB)
	authService := services.NewAuthService(authRepo)
	jwtService := services.NewJWTService(that.Config.App.JwtSecret)

	return &AuthController{
		BaseController: that,
		AuthService:    authService,
		JWTService:     jwtService,
	}
}

// Login handles the login route
func (that *AuthController) Login(c *fiber.Ctx) error {
	var request authrequests.LoginRequest

	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return responses.Error(c, err.Error())
	}

	// Validate the request
	if err := request.Validate(&request); err != nil {
		return responses.Error(c, err.Error())
	}

	// Authenticate user
	authenticated, err := that.AuthService.Authenticate(request.Username, request.Password)
	if err != nil {
		return responses.Error(c, "Authentication failed")
	}

	if !authenticated {
		return responses.Error(c, "Invalid credentials")
	}

	// Generate JWT token
	user, _ := that.AuthService.GetUserByUsername(request.Username) // Assuming you have a GetUserByUsername method
	token, err := that.JWTService.GenerateToken(user)
	if err != nil {
		return responses.Error(c, "Failed to generate token")
	}

	// Send JWT token in the response
	return responses.Success(c, fiber.Map{
		"message": "Login successful",
		"token":   token,
	})
}

// Register handles the registration route
func (that *AuthController) Register(c *fiber.Ctx) error {
	var request authrequests.RegisterRequest
	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return responses.Error(c, err.Error())
	}
	if err := request.Validate(&request); err != nil {
		return responses.Error(c, err.Error())
	}

	user := &models.User{
		Name:     request.Name,
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
		// Add other user properties as needed
	}

	if err := that.AuthService.Register(user); err != nil {
		return responses.Error(c, "Registration failed")
	}

	return responses.Success(c, fiber.Map{"message": "Registration successful"})
}

// UpdateProfile handles the update profile route
func (that *AuthController) UpdateProfile(c *fiber.Ctx) error {
	var request authrequests.UpdateProfileRequest
	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return responses.Error(c, err.Error())
	}
	if err := request.Validate(&request); err != nil {
		return responses.Error(c, err.Error())
	}

	// Assuming you have a way to identify the current user (e.g., from the JWT token)
	username := "current_username"

	updatedUser := &models.User{
		// Update user properties as needed
	}

	if err := that.AuthService.UpdateProfile(username, updatedUser); err != nil {
		return responses.Error(c, "Profile update failed")
	}

	return responses.Success(c, fiber.Map{"message": "Profile updated successfully"})
}

// ForgetPasswordHandler handles the forget password route
func (that *AuthController) ForgetPassword(c *fiber.Ctx) error {
	var request authrequests.ForgetPasswordRequest
	// Parse the request body
	if err := c.BodyParser(&request); err != nil {
		return responses.Error(c, err.Error())
	}
	if err := request.Validate(&request); err != nil {
		return responses.Error(c, err.Error())
	}

	// Assuming you have a way to identify the current user (e.g., from the JWT token)
	username := "current_username"

	resetToken, err := that.AuthService.ForgetPassword(username)
	if err != nil {
		return responses.Error(c, "Failed to initiate password reset")
	}
	log.Println(resetToken)
	// Send resetToken to the user (e.g., via email)

	return responses.Success(c, fiber.Map{"message": "Password reset initiated"})
}
func (that *AuthController) Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	return responses.Success(c, fiber.Map{"message": "Profile updated successfully", "user": user})
}
