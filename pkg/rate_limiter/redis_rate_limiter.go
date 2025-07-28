package redis_rate_limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const rateLimitScript = `
	local counter_key = KEYS[1]
	local ban_key = KEYS[2]

	-- define rate, window, and ban duration from ARGV
	local limit = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])
	local ban_duration = tonumber(ARGV[3])

	-- check if the client is banned
	local is_banned = redis.call("GET", ban_key)
	if is_banned then
		return {-1, redis.call("TTL", ban_key)}
	end

	-- increment the counter
	local current = redis.call("INCR", counter_key)

	-- if this is the first request in the window, set the expiration
	if current == 1 then
		redis.call("EXPIRE", counter_key, window)
	end

	-- if we've exceeded the limit, set ban and return 0
	if current > limit then
		redis.call("SET", ban_key, 1, "EX", ban_duration)
		return {0, ban_duration}
	end

	-- return 1 and remaining requests
	return {1, limit - current}
`

type RedisLimiter struct {
	Client      *redis.Client
	Limit       int
	Window      time.Duration
	BanDuration time.Duration
}

func NewRedisLimiter(client *redis.Client, limit int, window, ban time.Duration) *RedisLimiter {
	return &RedisLimiter{
		Client:      client,
		Limit:       limit,
		Window:      window,
		BanDuration: ban,
	}
}

func (r *RedisLimiter) Check(key string) (int, int, error) {
	ctx := context.Background()

	counterKey := fmt.Sprintf("rl_counter:%s", key)
	banKey := fmt.Sprintf("rl_ban:%s", key)

	result, err := r.Client.Eval(
		ctx,
		rateLimitScript,
		[]string{counterKey, banKey},
		r.Limit,
		int(r.Window.Seconds()),
		int(r.BanDuration.Seconds()),
	).Result()

	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute rate limit script: %w", err)
	}

	results, ok := result.([]interface{})
	if !ok || len(results) != 2 {
		return 0, 0, fmt.Errorf("unexpected result format from Redis script")
	}

	allowed, _ := results[0].(int64)
	ttl, _ := results[1].(int64)

	return int(allowed), int(ttl), nil
}
