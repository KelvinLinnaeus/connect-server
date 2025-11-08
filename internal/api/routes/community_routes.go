package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupCommunityRoutes sets up all community-related routes
func SetupCommunityRoutes(r *gin.RouterGroup, communityHandler *handlers.CommunityHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	communities := r.Group("/communities")
	communities.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes (no auth required or optional auth)
		communities.GET("/search", communityHandler.SearchCommunities)
		communities.GET("/categories", communityHandler.GetCommunityCategories)
		communities.GET("", communityHandler.ListCommunities)
		communities.GET("/:id", communityHandler.GetCommunity)
		communities.GET("/slug/:slug", communityHandler.GetCommunityBySlug)
		communities.GET("/:id/members", communityHandler.GetCommunityMembers)
		communities.GET("/:id/moderators", communityHandler.GetCommunityModerators)
		communities.GET("/:id/admins", communityHandler.GetCommunityAdmins)
		
		// Protected routes (auth required)
		communitiesAuth := communities.Group("")
		communitiesAuth.Use(middleware.AuthMiddleware(tokenMaker))
		{
			communitiesAuth.POST("", communityHandler.CreateCommunity)
			communitiesAuth.PUT("/:id", communityHandler.UpdateCommunity)
			communitiesAuth.POST("/:id/join", communityHandler.JoinCommunity)
			communitiesAuth.POST("/:id/leave", communityHandler.LeaveCommunity)
			communitiesAuth.POST("/:id/moderators", communityHandler.AddCommunityModerator)
			communitiesAuth.DELETE("/:id/moderators/:userId", communityHandler.RemoveCommunityModerator)
		}
	}
	
	// User-specific community routes
	users := r.Group("/users")
	users.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	users.Use(middleware.AuthMiddleware(tokenMaker))
	{
		users.GET("/communities", communityHandler.GetUserCommunities)
	}
}
