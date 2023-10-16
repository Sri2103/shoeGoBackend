package userHandler

import (
	"net/http"

	"github.com/hashicorp/go-hclog"
	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	userService "github.com/sri2103/shoeMart/internal/app/user/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

type User struct {
	userService userService.UserService
	logger      hclog.Logger
	config      *utils.Config
	validator   *utils.Validation
}

type UserKey struct{}

type UserIDKey struct{}

type VerificationDataKey struct{}

func NewUser(userService userService.UserService, l hclog.Logger, c *utils.Config, validator *utils.Validation) *User {
	return &User{
		userService: userService,
		logger:      l,
		config:      c,
		validator:   validator,
	}
}

// SignUp handle
func (u *User) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(UserKey{}).(userModel.User)

	userCreated, err := u.userService.AddUser(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)
		return
	}

	u.logger.Debug("User created successfully")
	w.WriteHeader(http.StatusCreated)
	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "User created",
		Data: &userModel.User{
			ID:       userCreated.ID,
			Username: userCreated.Username,
			Email:    userCreated.Email,
		},
	}, w)
}

func (u *User) LogIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqUser := r.Context().Value(UserKey{}).(userModel.User)

	user, accessToken, refreshToken, expirationTime, err := u.userService.ValidateLogIn(&reqUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: err.Error(),
		}, w)
		return
	}
	time := expirationTime.Unix()
	u.logger.Debug("User logged in successfully")
	w.WriteHeader(http.StatusAccepted)
	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "Logged in",
		Data: struct {
			ID           string `json:"id"`
			AccessToken  string `json:"accessToken"`
			RefreshToken string `json:"refreshToken"`
			Username     string `json:"username"`
			Email        string `json:"email"`
			Time         int64  `json:"accessExpiry"`
		}{user.ID, accessToken, refreshToken, user.Username, user.Email, time},
	}, w)

}

func (u *User) RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := r.Context().Value(UserKey{}).(userModel.User)
	accessToken, expirationTime, err := u.userService.GenerateAccessToken(&user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.logger.Error("unable to generate accessToken", err)
		utils.ToJson(&utils.GenericResponse{
			Success: false,
			Message: "Unable to generate token",
		}, w)
		return
	}

	w.WriteHeader(http.StatusOK)

	utils.ToJson(&utils.GenericResponse{
		Success: true,
		Message: "token generated",
		Data: struct {
			AccessToken    string `json:"accessToken"`
			ExpirationTime int64  `json:"expirationTime"`
		}{
			AccessToken: accessToken, 
			ExpirationTime: expirationTime.Unix()},
	}, w)
}
