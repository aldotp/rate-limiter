package bootstrap

import (
	"github.com/aldotp/rate-limiter/config"
	"github.com/aldotp/rate-limiter/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func (b *Bootstrap) setConfig() error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	b.Config = cfg
	return nil
}

func (b *Bootstrap) setLogger() error {
	log, err := logger.InitLogger(b.Config.App.Env)
	if err != nil {
		return err
	}
	b.Log = log
	return nil
}

func (b *Bootstrap) setRedisClient() error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     b.Config.Redis.Addr,
		Password: b.Config.Redis.Password,
		DB:       b.Config.Redis.Db,
	})

	if err := redisClient.Ping(b.ctx).Err(); err != nil {
		return err
	}

	b.RedisClient = redisClient

	return nil
}
