package userService

import userModel "github.com/sri2103/shoeMart/internal/app/user/models"

type UserService interface {
	AddUser(user *userModel.User) (*userModel.User, error)
	// FindById(id string)(*userModel.User, error)
	ValidateLogIn(reqUser *userModel.User) (*userModel.User,  string,  string,error)

	ValidateAccessToken(token string) (string, error)
	ValidateRefreshToken(token string) (string, string, error)
	GetUserById(userId string)(*userModel.User,error)
	GenerateCustomKey(userId string, tokenHash string) string
}
