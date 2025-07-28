package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App         *App
		Redis       *Redis
		HTTP        *HTTP
		RateLimiter *RateLimiter
	}

	App struct {
		Name string
		Env  string
	}

	Redis struct {
		Addr     string
		Password string
		Db       int
	}

	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}

	RateLimiter struct {
		Limit       int
		Window      time.Duration
		BanDuration time.Duration
	}
)

func New() (*Config, error) {
	LoadConfig()

	app := &App{
		Name: getEnv("APP_NAME"),
		Env:  getEnv("APP_ENV"),
	}

	http := &HTTP{
		Env:            getEnv("APP_ENV"),
		URL:            getEnv("HTTP_URL"),
		Port:           getEnv("HTTP_PORT"),
		AllowedOrigins: getEnv("HTTP_ALLOWED_ORIGINS"),
	}

	window, _ := time.ParseDuration(getEnv("REDIS_WINDOW"))
	banDuration, _ := time.ParseDuration(getEnv("REDIS_BAN_DURATION"))

	redis := &Redis{
		Addr:     getEnv("REDIS_ADDR"),
		Password: getEnv("REDIS_PASSWORD"),
	}

	rateLimiter := &RateLimiter{
		Limit:       viper.GetInt("REDIS_LIMIT"),
		Window:      window,
		BanDuration: banDuration,
	}

	return &Config{
		App:         app,
		Redis:       redis,
		HTTP:        http,
		RateLimiter: rateLimiter,
	}, nil
}

func LoadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("no .env file found, using system environment variables")
	}
	return nil
}

func getEnv(key string) string {
	return viper.GetString(key)
}
