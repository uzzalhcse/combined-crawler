package services

import (
	"combined-crawler/api/app/models"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type JWTServiceImpl struct {
	SecretKey string
}

func NewJWTService(secretKey string) *JWTServiceImpl {
	return &JWTServiceImpl{SecretKey: secretKey}
}

func (s *JWTServiceImpl) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   strconv.Itoa(int(user.ID)),
		"iss":   "my-app",                              // Replace with your desired issuer
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expiration time (adjust as needed)
		"iat":   time.Now().Unix(),
		"name":  user.Name,
		"email": user.Email,
		// Add other custom claims as needed
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
