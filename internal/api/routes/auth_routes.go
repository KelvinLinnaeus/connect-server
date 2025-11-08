package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
)

// SetupAuthRoutes configures authentication and session routes
func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, tokenMaker auth.Maker) {
	// Auth routes under /users
	users := router.Group("/users")
	{
		users.POST("/login", authHandler.Login)
		users.POST("/refresh", authHandler.RefreshToken)
		
		// Protected auth routes
		authUsers := users.Group("")
		authUsers.Use(middleware.AuthMiddleware(tokenMaker))
		{
			authUsers.POST("/logout", authHandler.Logout)
		}
	}

	// Session routes
	// sessions := router.Group("/sessions")
	// sessions.Use(middleware.AuthMiddleware(tokenMaker))
	// {
	// 	sessions.GET("/:id", authHandler.GetSession)
	// 	// sessions.GET("/me", authHandler.GetCurrentUserSessions)  // List user's sessions
	// }
}