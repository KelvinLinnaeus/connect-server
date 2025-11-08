package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupPostRoutes sets up post-related routes
func SetupPostRoutes(r *gin.RouterGroup, postHandler *handlers.PostHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	posts := r.Group("/posts")
	posts.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes (optional auth)
		posts.GET("/search", postHandler.SearchPosts)
		posts.GET("/advanced-search", postHandler.AdvancedSearchPosts)
		posts.GET("/trending", postHandler.GetTrendingPosts)
		posts.GET("/:id", postHandler.GetPost)
		posts.GET("/:id/comments", postHandler.GetPostComments)
		posts.GET("/:id/likes", postHandler.GetPostLikes)
		posts.GET("/user/:user_id", postHandler.GetUserPosts)
		posts.GET("/community/:community_id", postHandler.GetCommunityPosts)
		posts.GET("/group/:group_id", postHandler.GetGroupPosts)

		// Protected routes - require authentication
		postsAuth := posts.Group("")
		postsAuth.Use(middleware.AuthMiddleware(tokenMaker))
		{
			postsAuth.POST("", postHandler.CreatePost)
			postsAuth.DELETE("/:id", postHandler.DeletePost)
			postsAuth.GET("/feed", postHandler.GetUserFeed)
			postsAuth.GET("/liked", postHandler.GetUserLikedPosts)
			postsAuth.POST("/:id/comments", postHandler.CreateComment)
			postsAuth.POST("/:id/repost", postHandler.CreateRepost)
			postsAuth.POST("/:id/like", postHandler.TogglePostLike)
			postsAuth.PUT("/:id/pin", postHandler.PinPost)
		}
	}

	// Topics routes (related to posts)
	topics := r.Group("/topics")
	topics.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public route
		topics.GET("/trending", postHandler.GetTrendingTopics) // Get trending topics/hashtags
	}

	// Comments routes
	comments := r.Group("/comments")
	comments.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	comments.Use(middleware.AuthMiddleware(tokenMaker))
	{
		comments.POST("/:id/like", postHandler.ToggleCommentLike)
	}
}
