package routes

import (
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
)


func SetupSessionRoutes(r *gin.RouterGroup, sessionHandler *handlers.SessionHandler, tokenMaker auth.Maker) {
	sessions := r.Group("/sessions")
	{
		
		sessions.GET("/:id", middleware.AuthMiddleware(tokenMaker), sessionHandler.GetSession)
	}
}
