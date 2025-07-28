package redis_rate_limiter

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aldotp/rate-limiter/internal/bootstrap"
	"github.com/stretchr/testify/assert"
)

func setupTestLimiter() *RedisLimiter {
	os.Setenv("REDIS_LIMIT", "100")
	os.Setenv("REDIS_WINDOW", "5s")
	os.Setenv("REDIS_BAN_DURATION", "1m")

	f, err := bootstrap.NewBootstrap(context.Background()).BuildDependencies()
	if err != nil {
		log.Fatal(err)
	}

	f.RedisClient.FlushDB(f.RedisClient.Context())
	return NewRedisLimiter(f.RedisClient, f.Config.RateLimiter.Limit, f.Config.RateLimiter.Window, f.Config.RateLimiter.BanDuration)
}

func TestAllowWithinLimit(t *testing.T) {
	limiter := setupTestLimiter()
	key := "client-within-limit"

	t.Run("should allow multiple requests within limit", func(t *testing.T) {
		for i := 1; i <= 100; i++ {
			allowed, ttl, err := limiter.Check(key)
			log.Println("allowed:", allowed, "ttl:", ttl)
			assert.NoError(t, err)
			assert.Equal(t, 1, allowed, "expected request to be allowed")
		}
	})
}

func TestExceedLimit(t *testing.T) {
	limiter := setupTestLimiter()
	key := "client-exceed-limit"

	t.Run("should deny request after limit is exceeded", func(t *testing.T) {
		for i := 1; i <= 100; i++ {
			allowed, ttl, err := limiter.Check(key)
			log.Println("allowed:", allowed, "ttl:", ttl)
			assert.NoError(t, err)
			assert.Equal(t, 1, allowed, "expected request to be allowed")
		}

		// The next request should be denied and return TTL
		allowed, ttl, err := limiter.Check(key)
		log.Println("allowed:", allowed)
		fmt.Println("ttl:", ttl)
		assert.NoError(t, err)
		assert.Equal(t, 0, allowed, "expected request to be denied after exceeding limit")
		assert.True(t, ttl > 0, "expected a positive TTL value")
	})
}

func TestBannedClient(t *testing.T) {
	limiter := setupTestLimiter()
	key := "banned-client"

	t.Run("should return banned status after repeated violations", func(t *testing.T) {
		// Trigger ban
		for i := 1; i <= 101; i++ {
			allowed, ttl, err := limiter.Check(key)
			log.Println("allowed:", allowed)
			fmt.Println("ttl:", ttl)
			assert.NoError(t, err)
		}

		// Check if banned
		allowed, ttl, err := limiter.Check(key)
		log.Println("allowed:", allowed)
		fmt.Println("ttl:", ttl)
		assert.NoError(t, err)
		assert.Equal(t, -1, allowed, "expected client to be banned")
		assert.True(t, ttl > 0, "expected TTL on ban to be positive")
	})
}
