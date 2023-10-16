package utils

import (
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort                 string
	DBHost                     string
	DBName                     string
	DBUser                     string
	DBPass                     string
	DBPort                     string
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int
	SSL_Mode                   string
}

func NewConfig(logger hclog.Logger) *Config {
	// viper.SetConfigFile(".env")
	// viper.ReadInConfig()

	viper.AutomaticEnv()
	viper.SetDefault("ServerPort", "5000")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_Name", "postgresDB")
	viper.SetDefault("DB_User", "postgresUser")
	viper.SetDefault("DB_Pass", "postgresPW")
	viper.SetDefault("DB_Port", "5432")
	viper.SetDefault("SSL_Mode", "disable")
	viper.SetDefault("ACCESS_TOKEN_PRIVATE_KEY_PATH", "./access-private.pem")
	viper.SetDefault("ACCESS_TOKEN_PUBLIC_KEY_PATH", "./access-public.pem")
	viper.SetDefault("REFRESH_TOKEN_PRIVATE_KEY_PATH", "./refresh-private.pem")
	viper.SetDefault("REFRESH_TOKEN_PUBLIC_KEY_PATH", "./refresh-public.pem")
	viper.SetDefault("JWT_EXPIRATION", 10)
	config := &Config{
		ServerPort:                 viper.GetString("ServerPort"),
		DBHost:                     viper.GetString("DB_HOST"),
		DBName:                     viper.GetString("DB_NAME"),
		DBUser:                     viper.GetString("DB_User"),
		DBPass:                     viper.GetString("DB_Pass"),
		DBPort:                     viper.GetString("DB_Port"),
		JwtExpiration:              viper.GetInt("JWT_EXPIRATION"),
		AccessTokenPrivateKeyPath:  viper.GetString("ACCESS_TOKEN_PRIVATE_KEY_PATH"),
		AccessTokenPublicKeyPath:   viper.GetString("ACCESS_TOKEN_PUBLIC_KEY_PATH"),
		RefreshTokenPrivateKeyPath: viper.GetString("REFRESH_TOKEN_PRIVATE_KEY_PATH"),
		RefreshTokenPublicKeyPath:  viper.GetString("REFRESH_TOKEN_PUBLIC_KEY_PATH"),
		SSL_Mode:                   viper.GetString("SSL_Mode"),
	}

	// logger.Debug("db host", config.DBHost)
	// logger.Debug("db name", config.DBName)
	// logger.Debug("db port", config.DBPort)
	// logger.Debug("db user", config.DBUser)
	// logger.Debug("ssl_mode", config.SSL_Mode)
	// logger.Debug("jwt expiration", config.JwtExpiration)
	// logger.Debug("path to accessToken pem", config.AccessTokenPublicKeyPath)

	return config
}
