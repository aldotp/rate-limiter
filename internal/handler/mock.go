package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Ping(c *gin.Context) {
	m.Called(c)
}

var _ HandlerInterface = (*MockHandler)(nil)
