package bootstrap

import (
	"context"

	"github.com/aldotp/rate-limiter/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Option func(*Bootstrap)

type Bootstrap struct {
	ctx         context.Context
	RedisClient *redis.Client
	Log         *zap.Logger
	Config      *config.Config
}

func NewBootstrap(ctx context.Context) *Bootstrap {
	return &Bootstrap{
		ctx: ctx,
	}
}
