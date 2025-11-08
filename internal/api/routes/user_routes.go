package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
)


func SetupUserRoutes(router *gin.RouterGroup, userHandler *handlers.UserHandler, tokenMaker auth.Maker) {
	users := router.Group("/users")
	{
		
		users.POST("", userHandler.CreateUser)                         
		users.GET("/search", userHandler.SearchUsers)                  
		users.GET("/username/:username", userHandler.GetUserByUsername) 

		
		authUsers := users.Group("")
		authUsers.Use(middleware.AuthMiddleware(tokenMaker))
		{
			authUsers.GET("/:id", userHandler.GetUser)                     
			authUsers.PUT("/:id", userHandler.UpdateUser)                  
			authUsers.PUT("/:id/password", userHandler.UpdatePassword)     
			authUsers.DELETE("/:id", userHandler.DeactivateUser)           
			authUsers.GET("/suggested", userHandler.GetSuggestedUsers)     

			
			authUsers.POST("/:id/follow", userHandler.FollowUser)          
			authUsers.DELETE("/:id/follow", userHandler.UnfollowUser)      
			authUsers.GET("/:id/following/status", userHandler.CheckIfFollowing) 
			authUsers.GET("/:id/followers", userHandler.GetFollowers)      
			authUsers.GET("/:id/following", userHandler.GetFollowing)      

			
		}
	}
}