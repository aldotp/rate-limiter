package middleware

import (
	"github.com/aldotp/rate-limiter/config"
	redis_rate_limiter "github.com/aldotp/rate-limiter/pkg/rate_limiter"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Middleware struct {
	client    *redis.Client
	cfg       *config.Config
	log       *zap.Logger
	rateLimit *redis_rate_limiter.RedisLimiter
}

var _ MiddlewareInterface = (*Middleware)(nil)

func NewMiddleware(client *redis.Client, cfg *config.Config, log *zap.Logger, rateLimit *redis_rate_limiter.RedisLimiter) *Middleware {
	return &Middleware{
		client:    client,
		cfg:       cfg,
		log:       log,
		rateLimit: rateLimit,
	}
}
