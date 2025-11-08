package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Bytes("stack", debug.Stack()).
					Msg("Panic recovered")

				// Return error response
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					util.NewErrorResponse("internal_error", fmt.Sprintf("Internal server error: %v", err)),
				)
			}
		}()

		c.Next()
	}
}
