package http

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aldotp/rate-limiter/internal/bootstrap"
	"github.com/aldotp/rate-limiter/internal/handler"
	"github.com/aldotp/rate-limiter/internal/middleware"
	"github.com/aldotp/rate-limiter/internal/router"
	redis_rate_limiter "github.com/aldotp/rate-limiter/pkg/rate_limiter"
)

// NewHTTPServer builds the HTTP server components and returns the serve function
func NewHTTPServer(ctx context.Context) (*router.Router, string, error) {
	bootstrap, err := bootstrap.NewBootstrap(ctx).BuildDependencies()
	if err != nil {
		return nil, "", fmt.Errorf("error building bootstrap: %w", err)
	}

	ratelimit := redis_rate_limiter.NewRedisLimiter(
		bootstrap.RedisClient,
		bootstrap.Config.RateLimiter.Limit,
		bootstrap.Config.RateLimiter.Window,
		bootstrap.Config.RateLimiter.BanDuration,
	)

	mdl := middleware.NewMiddleware(bootstrap.RedisClient, bootstrap.Config, bootstrap.Log, ratelimit)
	hdl := handler.NewHandler()

	routes := router.NewRouter(bootstrap.Config, mdl, hdl)
	routes.SetupRouter()

	listenAddr := fmt.Sprintf("%s:%s", bootstrap.Config.HTTP.URL, bootstrap.Config.HTTP.Port)
	return routes, listenAddr, nil
}

func RunHTTPServer(ctx context.Context) {
	routes, listenAddr, err := NewHTTPServer(ctx)
	if err != nil {
		log.Println("Error building HTTP server:", err)
		os.Exit(1)
	}

	err = routes.Serve(listenAddr)
	if err != nil {
		log.Println("Error starting the HTTP server:", err)
		os.Exit(1)
	}
}
