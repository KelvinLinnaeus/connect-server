package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				
				log.Error().
					Interface("error", err).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Bytes("stack", debug.Stack()).
					Msg("Panic recovered")

				
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					util.NewErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Internal server error: %v", err)),
				)
			}
		}()

		c.Next()
	}
}
