package services

import "combined-crawler/api/app/models"

type JWTService interface {
	GenerateToken(user *models.User) (string, error)
}
