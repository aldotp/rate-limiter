package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	os.Setenv("APP_NAME", "TestApp")
	os.Setenv("APP_ENV", "development")
	os.Setenv("REDIS_ADDR", "localhost:6379")
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("HTTP_URL", "http://localhost")
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("HTTP_ALLOWED_ORIGINS", "*")
	os.Setenv("REDIS_LIMIT", "3")
	os.Setenv("REDIS_WINDOW", "5s")
	os.Setenv("REDIS_BAN_DURATION", "1m")

	cfg, err := New()
	assert.NoError(t, err)
	assert.Equal(t, "TestApp", cfg.App.Name)
	assert.Equal(t, "development", cfg.App.Env)
	assert.Equal(t, "localhost:6379", cfg.Redis.Addr)
	assert.Equal(t, "", cfg.Redis.Password)
	assert.Equal(t, "http://localhost", cfg.HTTP.URL)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "*", cfg.HTTP.AllowedOrigins)
	assert.Equal(t, 3, cfg.RateLimiter.Limit)
	assert.Equal(t, 5*time.Second, cfg.RateLimiter.Window)
	assert.Equal(t, 1*time.Minute, cfg.RateLimiter.BanDuration)
}
