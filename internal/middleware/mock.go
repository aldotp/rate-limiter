package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) RateLimiter() gin.HandlerFunc {
	args := m.Called()
	if handler, ok := args.Get(0).(gin.HandlerFunc); ok && handler != nil {
		return handler
	}
	// Default fallback if nil
	return func(c *gin.Context) {
		c.Next()
	}
}

func (m *MockMiddleware) Log() gin.HandlerFunc {
	args := m.Called()
	if handler, ok := args.Get(0).(gin.HandlerFunc); ok && handler != nil {
		return handler
	}
	// Default fallback if nil
	return func(c *gin.Context) {
		c.Next()
	}
}
