package middleware

import (
	"github.com/gin-gonic/gin"
)

type MiddlewareInterface interface {
	RateLimiter() gin.HandlerFunc
	Log() gin.HandlerFunc
}
