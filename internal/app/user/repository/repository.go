package userRepository

import userModel "github.com/sri2103/shoeMart/internal/app/user/models"

type UserRepo interface {
	AddUser(user *userModel.User)(*userModel.User, error)
	// FindById(id string)(*userModel.User, error)
	FindUser(email string)(*userModel.User,error)

	GetUserById(userId string) (*userModel.User, error)
}