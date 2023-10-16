package authService

import (
	"time"

	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
)

type AuthService interface {
	Authenticate(user *userModel.User, reqUser *userModel.User) bool
	GenerateRefreshToken(user *userModel.User) (string, error)
	GenerateAccessToken(user *userModel.User) (string, time.Time, error)
	GenerateCustomKey(userID string, tokenHash string) string
	// return userId custom key from token payload
	ValidateRefreshToken(tokenString string) (string, string, error)

	// return userId from token payload
	ValidateAccessToken(tokenString string) (string, error)
}
