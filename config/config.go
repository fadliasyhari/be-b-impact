package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type ApiConfig struct {
	ApiPort string
	ApiHost string
}

type RedisConfig struct {
	Url string
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Env      string
}

type TokenConfig struct {
	ApplicationName     string
	JwtSigantureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenlifeTime time.Duration
}

type FileConfig struct {
	UploadPath string
	LogPath    string
	Env        string
}

type Config struct {
	DbConfig
	ApiConfig
	FileConfig
	TokenConfig
	RedisConfig
}

func (c *Config) ReadConfigFile() error {

	c.RedisConfig = RedisConfig{
		Url: os.Getenv("REDIS_URL"),
	}

	c.DbConfig = DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort: os.Getenv("PORT"),
		ApiHost: os.Getenv("API_HOST"),
	}

	c.FileConfig = FileConfig{
		LogPath:    os.Getenv("FILE_PATH"),
		Env:        os.Getenv("ENV"),
		UploadPath: os.Getenv("UPLOAD_PATH"),
	}

	tokenDuration, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))

	c.TokenConfig = TokenConfig{
		ApplicationName:     os.Getenv("APP_NAME"),
		JwtSigantureKey:     os.Getenv("SECRET_KEY"),
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenlifeTime: time.Minute * time.Duration(tokenDuration),
	}

	if c.DbConfig.Host == "" || c.DbConfig.Port == "" || c.DbConfig.Name == "" ||
		c.DbConfig.User == "" || c.DbConfig.Password == "" ||
		c.ApiConfig.ApiHost == "" || c.ApiConfig.ApiPort == "" ||
		c.FileConfig.LogPath == "" || c.FileConfig.Env == "" || c.UploadPath == "" ||
		c.RedisConfig.Url == "" {
		return errors.New("missing required environment variables")
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := cfg.ReadConfigFile()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
