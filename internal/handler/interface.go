package handler

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	Ping(c *gin.Context)
}
