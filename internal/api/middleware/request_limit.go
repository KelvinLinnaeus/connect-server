package middleware

import (
	"fmt"
	"net/http"

	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	
	DefaultMaxBodySize int64 = 1 << 20 

	
	LargeMaxBodySize int64 = 10 << 20 

	
	MaxHeaderSize int64 = 8 << 10 
)



func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		
		contentLength := c.Request.ContentLength
		if contentLength > maxSize {
			log.Warn().
				Int64("content_length", contentLength).
				Int64("max_size", maxSize).
				Str("method", c.Request.Method).
				Str("path", c.Request.URL.Path).
				Str("remote_addr", c.ClientIP()).
				Msg("Request rejected: Content-Length exceeds maximum allowed size")

			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge,
				util.NewErrorResponse(http.StatusRequestEntityTooLarge,
					fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize)))
			return
		}

		
		
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()

		
		if c.Errors.Last() != nil {
			err := c.Errors.Last().Err
			if err.Error() == "http: request body too large" {
				log.Warn().
					Err(err).
					Int64("max_size", maxSize).
					Str("method", c.Request.Method).
					Str("path", c.Request.URL.Path).
					Str("remote_addr", c.ClientIP()).
					Msg("Request body exceeded maximum size during read")

				c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge,
					util.NewErrorResponse(http.StatusRequestEntityTooLarge,
						fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize)))
			}
		}
	}
}


func DefaultRequestSizeLimitMiddleware() gin.HandlerFunc {
	return RequestSizeLimitMiddleware(DefaultMaxBodySize)
}


func LargeRequestSizeLimitMiddleware() gin.HandlerFunc {
	return RequestSizeLimitMiddleware(LargeMaxBodySize)
}
