package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aldotp/rate-limiter/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	h := handler.NewHandler()
	h.Ping(c)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.JSONEq(t, `{"message":"pong"}`, w.Body.String())
}
