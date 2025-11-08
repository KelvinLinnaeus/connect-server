package middleware

import (
	"fmt"
	"net/http"

	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultMaxBodySize is 1MB - suitable for most API requests
	DefaultMaxBodySize int64 = 1 << 20 // 1 MB

	// LargeMaxBodySize is 10MB - for endpoints that may need larger payloads (e.g., file uploads)
	LargeMaxBodySize int64 = 10 << 20 // 10 MB

	// MaxHeaderSize limits the size of request headers
	MaxHeaderSize int64 = 8 << 10 // 8 KB
)

// RequestSizeLimitMiddleware limits the size of incoming request bodies
// This prevents memory exhaustion attacks and resource abuse
func RequestSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip size check for GET, HEAD, and OPTIONS requests (they typically have no body)
		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// Check Content-Length header first (if present)
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
				util.NewErrorResponse("request_too_large",
					fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize)))
			return
		}

		// Limit the request body reader to prevent reading more than maxSize
		// This protects against clients that don't set Content-Length or set it incorrectly
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()

		// Check if MaxBytesReader triggered an error
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
					util.NewErrorResponse("request_too_large",
						fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize)))
			}
		}
	}
}

// DefaultRequestSizeLimitMiddleware applies the default 1MB size limit
func DefaultRequestSizeLimitMiddleware() gin.HandlerFunc {
	return RequestSizeLimitMiddleware(DefaultMaxBodySize)
}

// LargeRequestSizeLimitMiddleware applies a larger 10MB size limit for special endpoints
func LargeRequestSizeLimitMiddleware() gin.HandlerFunc {
	return RequestSizeLimitMiddleware(LargeMaxBodySize)
}
