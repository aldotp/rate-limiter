package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aldotp/rate-limiter/config"
	"github.com/aldotp/rate-limiter/internal/handler"
	"github.com/aldotp/rate-limiter/internal/middleware"
)

func TestNewRouter(t *testing.T) {
	cfg := &config.Config{
		App: &config.App{
			Name: "test-app",
			Env:  "test",
		},
		Redis:       &config.Redis{},
		HTTP:        &config.HTTP{},
		RateLimiter: &config.RateLimiter{},
	}
	mdl := new(middleware.MockMiddleware)
	hdl := new(handler.MockHandler)

	router := NewRouter(cfg, mdl, hdl)

	assert.NotNil(t, router, "Router should not be nil")
	assert.Equal(t, cfg, router.cfg, "Config should be set correctly")
	assert.Equal(t, mdl, router.mdl, "Middleware should be set correctly")
	assert.Equal(t, hdl, router.hdl, "Handler should be set correctly")
}

func TestSetupRouter(t *testing.T) {
	tests := []struct {
		name     string
		env      string
		setMocks func(*middleware.MockMiddleware, *handler.MockHandler)
	}{
		{
			name: "development environment",
			env:  "development",
			setMocks: func(mdl *middleware.MockMiddleware, hdl *handler.MockHandler) {
				mdl.On("Log").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() }))
				mdl.On("RateLimiter").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() })).Once()
				hdl.On("Ping", mock.MatchedBy(func(c *gin.Context) bool { return true })).Run(func(args mock.Arguments) {
					c := args.Get(0).(*gin.Context)
					c.JSON(http.StatusOK, gin.H{"message": "pong"})
				}).Once()
			},
		},
		{
			name: "production environment",
			env:  "production",
			setMocks: func(mdl *middleware.MockMiddleware, hdl *handler.MockHandler) {
				mdl.On("Log").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() }))
				mdl.On("RateLimiter").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() })).Once()
				hdl.On("Ping", mock.MatchedBy(func(c *gin.Context) bool { return true })).Run(func(args mock.Arguments) {
					c := args.Get(0).(*gin.Context)
					c.JSON(http.StatusOK, gin.H{"message": "pong"})
				}).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req := httptest.NewRequest("GET", "/ping", nil)
			c.Request = req

			cfg := &config.Config{
				App: &config.App{
					Name: "test-app",
					Env:  tt.env,
				},
				Redis:       &config.Redis{},
				HTTP:        &config.HTTP{},
				RateLimiter: &config.RateLimiter{},
			}

			mdl := new(middleware.MockMiddleware)
			hdl := new(handler.MockHandler)

			// Setup mocks
			tt.setMocks(mdl, hdl)

			// Create and setup router
			r := NewRouter(cfg, mdl, hdl)
			r.SetupRouter()

			// Verify Gin mode
			if tt.env == "development" {
				assert.Equal(t, gin.DebugMode, gin.Mode())
			} else {
				assert.Equal(t, gin.ReleaseMode, gin.Mode())
			}

			r.Engine.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			mdl.AssertExpectations(t)
			hdl.AssertExpectations(t)
		})
	}
}

func TestServe(t *testing.T) {
	// Setup test config
	cfg := &config.Config{
		App: &config.App{
			Name: "test-app",
			Env:  "test",
		},
		Redis:       &config.Redis{},
		HTTP:        &config.HTTP{},
		RateLimiter: &config.RateLimiter{},
	}

	mdl := new(middleware.MockMiddleware)
	hdl := new(handler.MockHandler)

	mdl.On("RateLimiter").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() }))
	mdl.On("Log").Return(gin.HandlerFunc(func(c *gin.Context) { c.Next() }))

	r := NewRouter(cfg, mdl, hdl)
	r.SetupRouter()

	err := r.Serve("invalid-address:invalid-port")
	assert.Error(t, err, "Expected error for invalid address")
}
