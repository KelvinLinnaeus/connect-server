package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupSpaceRoutes sets up all space-related routes
func SetupSpaceRoutes(r *gin.RouterGroup, spaceHandler *handlers.SpaceHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	spaces := r.Group("/spaces")
	spaces.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes (no auth required or optional auth)
		spaces.GET("", spaceHandler.ListSpaces)
		spaces.GET("/:id", spaceHandler.GetSpace)
		spaces.GET("/slug/:slug", spaceHandler.GetSpaceBySlug)

		// Protected routes (auth required)
		spacesAuth := spaces.Group("")
		spacesAuth.Use(middleware.AuthMiddleware(tokenMaker))
		{
			spacesAuth.POST("", spaceHandler.CreateSpace)
			spacesAuth.PUT("/:id", spaceHandler.UpdateSpace)
			spacesAuth.DELETE("/:id", spaceHandler.DeleteSpace)
		}
	}
}
