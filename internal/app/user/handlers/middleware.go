package userHandler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	"github.com/sri2103/shoeMart/internal/app/utils"
)

func (uh *User) MiddlewareValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		user := &userModel.User{}

		err := utils.FromJson(user, r.Body)

		if err != nil {
			uh.logger.Error("error in parsing the request body ", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{
				Success: false,
				Message: err.Error(),
			}, w)
			return
		}

		errs := uh.validator.Validate(user)

		if len(errs) != 0 {
			uh.logger.Error("valid of user json failed", errs)
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{
				Success: false,
				Message: strings.Join(errs.Errors(), ","),
			}, w)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey{}, *user)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (uh *User) MiddlewareValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		uh.logger.Debug("Validating accessToken")

		token, err := extractToken(r)

		if err != nil {
			uh.logger.Error("Token not provided or malformed")
			w.WriteHeader(http.StatusBadRequest)
			// data.ToJSON(&GenericError{Error: err.Error()}, w)
			utils.ToJson(&utils.GenericResponse{Success: false, Message: "Authentication failed. Token not provided or malformed"}, w)
			return
		}

		uh.logger.Debug("token present in header", token)

		userID, err := uh.userService.ValidateAccessToken(token)

		if err != nil {
			uh.logger.Error("validation failed", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{
				Success: false,
				Message: "Authentication failed. Invalid token",
			}, w)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey{}, userID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (uh *User) MiddlewareValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type","application/json")

		uh.logger.Debug("Validating refresh token")

		token, err := extractToken(r)

		uh.logger.Debug("Token present in header",token)
		if err != nil {
			uh.logger.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{Success: false, Message: err.Error()}, w)
			return
		}

		uh.logger.Debug("Token present in request", token)

		userID, customKey,err := uh.userService.ValidateRefreshToken(token)

		if err != nil {
			uh.logger.Error("Token validation failed",err)
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{Success: false, Message: "Invalid Refresh Token"}, w)
			return
		}

		uh.logger.Debug("Refresh token validated")

		user,err := uh.userService.GetUserById(userID)

		if err != nil {
			uh.logger.Error("Failed to get user by id ",err )
			w.WriteHeader(http.StatusInternalServerError)
			utils.ToJson(&utils.GenericResponse{
				Success: false,
				Message: "Unable to fetch the user",
			},w)

			return
		}

		actualCustomKey := uh.userService.GenerateCustomKey(user.ID,user.TokenHash)
		if customKey != actualCustomKey {
			uh.logger.Debug("wrong token: authentication failed")
			w.WriteHeader(http.StatusBadRequest)
			utils.ToJson(&utils.GenericResponse{
				Success: false,
				Message: "Authentication failed",
			}, w)
			return
		}

		ctx := context.WithValue(r.Context(),UserKey{},*user)

		r = r.WithContext(ctx)
		next.ServeHTTP(w,r)
	})
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.SplitN(authHeader, " ", 2)
	if len(authHeaderContent) != 2 {
		return "", errors.New("Token not provided or malformed")
	}
	return authHeaderContent[1], nil
}
