package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type FormatterStatus interface {
	String() string
}

func (m *Middleware) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		req := c.Request
		var reqBody []byte
		if req.Body != nil {
			bodyBytes, _ := io.ReadAll(req.Body)
			reqBody = bodyBytes
			req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		traceID, _ := c.Get("requestid")
		traceIDStr := ""
		if v, ok := traceID.(string); ok {
			traceIDStr = v
		}

		c.Set("traceId", traceIDStr)
		c.Set("srcIP", c.ClientIP())
		c.Set("port", req.URL.Port())
		c.Set("path", req.URL.Path)

		respWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = respWriter

		c.Next()

		err := c.Errors.Last()
		statusCode := c.Writer.Status()
		code := "SUCCESS"
		if err != nil {
			code = err.Error()
		}

		responseTime := time.Since(startTime)

		m.log.Info("request log",
			zap.String("trace_id", traceIDStr),
			zap.String("method", req.Method),
			zap.String("path", req.URL.Path),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("status", statusCode),
			zap.String("status_code", code),
			zap.Duration("response_time", responseTime),
			zap.ByteString("request", reqBody),
			zap.ByteString("response", respWriter.body.Bytes()),
			zap.Any("header", req.Header),
			zap.Error(err),
		)
	}
}
