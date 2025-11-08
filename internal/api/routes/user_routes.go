package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
)

// SetupUserRoutes configures user-related routes
func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, tokenMaker auth.Maker) {
	users := router.Group("/users")
	{
		// Public routes
		users.POST("", userHandler.CreateUser)                         // Sign up
		users.GET("/search", userHandler.SearchUsers)                  // Search users (can be public or auth required)
		users.GET("/username/:username", userHandler.GetUserByUsername) // Get user by username

		// Protected routes
		authUsers := users.Group("")
		authUsers.Use(middleware.AuthMiddleware(tokenMaker))
		{
			authUsers.GET("/:id", userHandler.GetUser)                     // Get user by ID
			authUsers.PUT("/:id", userHandler.UpdateUser)                  // Update user
			authUsers.PUT("/:id/password", userHandler.UpdatePassword)     // Update password
			authUsers.DELETE("/:id", userHandler.DeactivateUser)           // Deactivate user
			authUsers.GET("/suggested", userHandler.GetSuggestedUsers)     // Get suggested users

			// Follow/Unfollow routes
			authUsers.POST("/:id/follow", userHandler.FollowUser)          // Follow a user
			authUsers.DELETE("/:id/follow", userHandler.UnfollowUser)      // Unfollow a user
			authUsers.GET("/:id/following/status", userHandler.CheckIfFollowing) // Check if following
			authUsers.GET("/:id/followers", userHandler.GetFollowers)      // Get user's followers
			authUsers.GET("/:id/following", userHandler.GetFollowing)      // Get users followed by user

			// authUsers.GET("/:id/stats", userHandler.GetUserStats)       // Get user stats
		}
	}
}