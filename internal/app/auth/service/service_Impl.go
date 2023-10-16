package authService

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/go-hclog"
	userModel "github.com/sri2103/shoeMart/internal/app/user/models"
	"github.com/sri2103/shoeMart/internal/app/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	logger hclog.Logger
	config *utils.Config
}

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID  string
	KeyType string
	jwt.RegisteredClaims
}
type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.RegisteredClaims
}

func NewAuthService(logger hclog.Logger, config *utils.Config) *AuthServiceImpl {
	return &AuthServiceImpl{
		logger: logger,
		config: config,
	}
}

func (auth *AuthServiceImpl) Authenticate(user *userModel.User, reqUser *userModel.User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		auth.logger.Debug("Password hashes are not same")
		return false
	}
	return true
}

// GenerateCustomKey creates a new key for our jwt payload
// the key is a hashed combination of the userID and user tokenhash
func (auth *AuthServiceImpl) GenerateCustomKey(userId string, tokenHash string) string {

	// data := userID + tokenHash

	h := hmac.New(sha256.New, []byte(tokenHash))
	h.Write([]byte(userId))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func (auth *AuthServiceImpl) GenerateAccessToken(user *userModel.User) (string, time.Time,error) {
	userId := user.ID
	tokenType := "access"

	//  Read private key from key path
	signBytes, err := os.ReadFile(auth.config.AccessTokenPrivateKeyPath)

	if err != nil {
		auth.logger.Error("unable to read privatekey", err)
		return "", time.Time{}, errors.New("could not generate access token. please try again later")
	}

	// Parse key to RSA format
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		auth.logger.Error("unable to parse private key", "error", err)
		return "", time.Time{}, errors.New("could not generate access token. please try again later")
	}

	expirationTime := time.Now().Add(time.Minute * time.Duration(auth.config.JwtExpiration))
	// Set claims to add to JWT
	claims := AccessTokenCustomClaims{
		userId,
		tokenType,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "shoeGo.auth.service",
		},
	}

	// generate token from attached token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	key,err := token.SignedString(signKey)

	return  key,expirationTime, err

}


func (auth *AuthServiceImpl) GenerateRefreshToken(user *userModel.User) (string, error) {

	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		user.ID,
		auth.GenerateCustomKey(user.ID, user.TokenHash),
		tokenType,
		jwt.RegisteredClaims{
			Issuer: "shoeGo.auth.service",
		},
	}

	signBytes, err := os.ReadFile(auth.config.RefreshTokenPrivateKeyPath)

	if err != nil {
		auth.logger.Error("unable to read privateKey", err)
		return "", errors.New("could not generate refresh token, please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		auth.logger.Error("unable to parse private key to RSA", err)
		return "", errors.New("could not generate refresh token please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func (auth *AuthServiceImpl) ValidateAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				auth.logger.Error("Unexpected signing method in auth Token")
				return nil, errors.New("unexpected signing method in auth token")
			}

			verifyBytes, err := os.ReadFile(auth.config.AccessTokenPublicKeyPath)

			if err != nil {
				auth.logger.Error("Unable to read public key from path")
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

			if err != nil {
				auth.logger.Error("Could Not Parse Public Key From PEM")
				return nil, errors.New("could not parse public key from PEM")
			}

			return verifyKey, nil

		})

	if err != nil {
		auth.logger.Error("Unable to parse claims")
		return "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return "", errors.New("invalid token: auth failed")
	}

	return claims.UserID, nil

}

func (auth *AuthServiceImpl) ValidateRefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)
			if !ok {
				auth.logger.Error("Invalid signing method")
				return nil, errors.New("unexpected signing method in authToken")
			}
			auth.logger.Info(auth.config.RefreshTokenPublicKeyPath,"key path")
			verifyBytes, err := os.ReadFile(auth.config.RefreshTokenPublicKeyPath)
			

			if err != nil {
				auth.logger.Error("unable to parse public key", err)
				return nil, err
			}

			verifyToken, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

			if err != nil {
				auth.logger.Error("Could Not Parse Public Key From PEM")
				return nil, errors.New("could not parse public key from PEM")
			}

			return verifyToken, nil

		})

	if err != nil {
		auth.logger.Error("Unable to get token")
		return "", "", err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" {
		auth.logger.Error("could not extract claims from token")
		return "", "", errors.New("invalid token: authentication failed")
	}

	return claims.UserID, claims.CustomKey, nil
}
