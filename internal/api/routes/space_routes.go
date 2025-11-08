package routes

import (
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
)


func SetupSpaceRoutes(r *gin.RouterGroup, spaceHandler *handlers.SpaceHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	spaces := r.Group("/spaces")
	spaces.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		
		spaces.GET("", spaceHandler.ListSpaces)
		spaces.GET("/:id", spaceHandler.GetSpace)
		spaces.GET("/slug/:slug", spaceHandler.GetSpaceBySlug)

		
		spacesAuth := spaces.Group("")
		spacesAuth.Use(middleware.AuthMiddleware(tokenMaker))
		{
			spacesAuth.POST("", spaceHandler.CreateSpace)
			spacesAuth.PUT("/:id", spaceHandler.UpdateSpace)
			spacesAuth.DELETE("/:id", spaceHandler.DeleteSpace)
		}
	}
}
