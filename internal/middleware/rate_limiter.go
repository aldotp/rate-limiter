package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *Middleware) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		key := c.Request.Header.Get("X-Api-Key")
		if key == "" {
			key = c.Request.RemoteAddr
			m.log.Debug("No API key provided, using remote address as key:", zap.String("key", key))
		}

		allowed, ttl, err := m.rateLimit.Check(key)
		if err != nil {
			m.log.Error("Rate limiter check failed", zap.Error(err), zap.String("key", key))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "internal server error",
				"code":    http.StatusInternalServerError,
			})
			return
		}

		switch allowed {
		case -1:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("already banned, remaining %ds", ttl),
				"code":    http.StatusTooManyRequests,
			})
			return
		case 0:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("rate limit exceeded, banned for %ds", ttl),
				"code":    http.StatusTooManyRequests,
			})
			return
		case 1:
			c.Next()
			return
		}
	}
}
