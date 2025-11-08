package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupSessionRoutes sets up session-related routes
func SetupSessionRoutes(r *gin.RouterGroup, sessionHandler *handlers.SessionHandler, tokenMaker auth.Maker) {
	sessions := r.Group("/sessions")
	{
		// Protected routes - require authentication
		sessions.GET("/:id", middleware.AuthMiddleware(tokenMaker), sessionHandler.GetSession)
	}
}
