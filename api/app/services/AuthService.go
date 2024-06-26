// Services/auth_service.go

package services

import "combined-crawler/api/app/models"

type AuthService interface {
	Authenticate(username, password string) (bool, error)
	GetUserByUsername(username string) (*models.User, error)
	Register(user *models.User) error
	UpdateProfile(username string, updatedUser *models.User) error
	ForgetPassword(username string) (string, error)
}
