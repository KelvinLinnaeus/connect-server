package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupGroupRoutes sets up all group-related routes
func SetupGroupRoutes(r *gin.RouterGroup, groupHandler *handlers.GroupHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	groups := r.Group("/groups")
	groups.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes (no auth required or optional auth)
		groups.GET("/search", groupHandler.SearchGroups)
		groups.GET("", groupHandler.ListGroups)
		groups.GET("/:id", groupHandler.GetGroup)
		groups.GET("/:id/roles", groupHandler.GetProjectRoles)
		
		// Protected routes (auth required)
		groupsAuth := groups.Group("")
		groupsAuth.Use(middleware.AuthMiddleware(tokenMaker))
		{
			groupsAuth.POST("", groupHandler.CreateGroup)
			groupsAuth.PUT("/:id", groupHandler.UpdateGroup)
			groupsAuth.POST("/:id/join", groupHandler.JoinGroup)
			groupsAuth.POST("/:id/leave", groupHandler.LeaveGroup)
			groupsAuth.GET("/:id/join-requests", groupHandler.GetGroupJoinRequests)
			groupsAuth.POST("/:id/admins", groupHandler.AddGroupAdmin)
			groupsAuth.DELETE("/:id/admins/:userId", groupHandler.RemoveGroupAdmin)
			groupsAuth.POST("/:id/moderators", groupHandler.AddGroupModerator)
			groupsAuth.DELETE("/:id/moderators/:userId", groupHandler.RemoveGroupModerator)
			groupsAuth.PUT("/:id/members/:userId/role", groupHandler.UpdateGroupMemberRole)
			groupsAuth.POST("/:id/roles", groupHandler.CreateProjectRole)
			groupsAuth.GET("/:id/applications", groupHandler.GetRoleApplications)
		}
	}
	
	// User-specific group routes
	users := r.Group("/users")
	users.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	users.Use(middleware.AuthMiddleware(tokenMaker))
	{
		users.GET("/groups", groupHandler.GetUserGroups)
	}
	
	// Project role application routes
	roles := r.Group("/roles")
	roles.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	roles.Use(middleware.AuthMiddleware(tokenMaker))
	{
		roles.POST("/:roleId/apply", groupHandler.ApplyForProjectRole)
	}
}
