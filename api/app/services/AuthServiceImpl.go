// app/Services/auth_service_impl.go

package services

import (
	"combined-crawler/api/app/models"
	"combined-crawler/api/app/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	AuthRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{AuthRepo: authRepo}
}

func (s *AuthServiceImpl) Authenticate(username, password string) (bool, error) {
	user, err := s.AuthRepo.FindUserByUsername(username)
	if err != nil {
		return false, err
	}

	// Check if the provided password matches the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, nil // Passwords do not match
	}

	return true, nil // Authentication successful
}

func (s *AuthServiceImpl) GetUserByUsername(username string) (*models.User, error) {
	return s.AuthRepo.FindUserByUsername(username)
}

func (s *AuthServiceImpl) Register(user *models.User) error {
	// Hash the user's password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the user's password with the hashed password
	user.Password = string(hashedPassword)

	// Save the user to the database
	return s.AuthRepo.CreateUser(user)
}

func (s *AuthServiceImpl) UpdateProfile(username string, updatedUser *models.User) error {
	// Implement logic to update user profile in the database
	return s.AuthRepo.UpdateUser(username, updatedUser)
}

func (s *AuthServiceImpl) ForgetPassword(username string) (string, error) {
	// Implement logic to handle forget password (e.g., generate and send reset token)
	// Return a reset token for further processing
	return "reset_token_here", nil
}
