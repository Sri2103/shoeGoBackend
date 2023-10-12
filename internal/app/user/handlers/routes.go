package userHandler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	authService "github.com/sri2103/shoeMart/internal/app/auth/service"
	userRepository "github.com/sri2103/shoeMart/internal/app/user/repository"
	userService "github.com/sri2103/shoeMart/internal/app/user/service"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

func SetupUserRoutes(userRepo userRepository.UserRepo, router *mux.Router, config *utils.Config, logger hclog.Logger, validation *utils.Validation) {
	auth := authService.NewAuthService(logger, config)
	userService := userService.NewUserService(logger, userRepo, auth)
	userHandler := NewUser(userService, logger, config,validation)

	authRouter := router.Methods(http.MethodPost).Subrouter()
	authRouter.Use(userHandler.MiddlewareValidateUser)
	authRouter.HandleFunc("/signin", userHandler.SignIn)
	authRouter.HandleFunc("/login", userHandler.LogIn)

}
