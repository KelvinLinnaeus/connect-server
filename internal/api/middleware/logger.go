package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)


func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		
		c.Next()

		
		latency := time.Since(start)

		
		statusCode := c.Writer.Status()

		
		if raw != "" {
			path = path + "?" + raw
		}

		
		logger := log.Info()

		if statusCode >= 500 {
			logger = log.Error()
		} else if statusCode >= 400 {
			logger = log.Warn()
		}

		logger.
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP Request")
	}
}
