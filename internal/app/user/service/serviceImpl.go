package userService

import (
	"errors"

	"github.com/hashicorp/go-hclog"
	authService "github.com/sri2103/shoeMart/internal/app/auth/service"
	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	userRepository "github.com/sri2103/shoeMart/internal/app/user/repository"
	"github.com/sri2103/shoeMart/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	logger   hclog.Logger
	userRepo userRepository.UserRepo
	auth     authService.AuthService
}

func NewUserService(l hclog.Logger, userRepo userRepository.UserRepo, auth authService.AuthService) *UserServiceImpl {
	return &UserServiceImpl{
		logger:   l,
		userRepo: userRepo,
		auth:     auth,
	}
}

func (u *UserServiceImpl) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("unable to hash password", "error", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func (u *UserServiceImpl) AddUser(user *userModel.User) (*userModel.User, error) {

	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		u.logger.Error("Unable to hash password")
		return nil, err
	}
	user.Password = hashedPassword
	user.TokenHash = utils.GenerateRandomString(15)

	userCreated, err := u.userRepo.AddUser(user)

	if err != nil {
		u.logger.Error("Unable to add new user in DB.")
		return nil, err
	}

	return userCreated, nil
}

func (u *UserServiceImpl) ValidateLogIn(reqUser *userModel.User) (*userModel.User, string, string, error) {
	user, err := u.userRepo.FindUser(reqUser.Email)

	if err != nil {
		u.logger.Error("Unable to find the user with email ", "Email:", reqUser.Email, ", error:", err)
		return nil, "", "", err
	}

	if valid := u.auth.Authenticate(user, reqUser); !valid {
		u.logger.Error("Passwords does not match")
		return nil, "", "", errors.New("passwords does not match")

	}

	accessToken, err := u.auth.GenerateAccessToken(user)

	if err != nil {
		u.logger.Error("Unable to generate access token")
		return nil, "", "", err

	}

	refreshToken, err := u.auth.GenerateRefreshToken(user)

	if err != nil {
		u.logger.Error("Unable to generate refresh token")
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (u *UserServiceImpl) ValidateAccessToken(token string) (string, error) {
	return u.auth.ValidateAccessToken(token)
}

func (u *UserServiceImpl) ValidateRefreshToken(token string) (string, string, error) {
	return u.auth.ValidateRefreshToken(token)
}

func (u *UserServiceImpl) GetUserById(userId string) (*userModel.User, error) {
	return u.userRepo.GetUserById(userId)
}

func (u *UserServiceImpl) GenerateCustomKey(userId string, tokenHash string) string {
	return u.auth.GenerateCustomKey(userId, tokenHash)
}
