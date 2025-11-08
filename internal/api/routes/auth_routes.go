package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
)


func SetupAuthRoutes(router *gin.RouterGroup, authHandler *handlers.AuthHandler, tokenMaker auth.Maker) {
	
	users := router.Group("/users")
	{
		users.POST("/login", authHandler.Login)
		users.POST("/refresh", authHandler.RefreshToken)
		
		
		authUsers := users.Group("")
		authUsers.Use(middleware.AuthMiddleware(tokenMaker))
		{
			authUsers.POST("/logout", authHandler.Logout)
		}
	}

	
	
	
	
	
	
	
}