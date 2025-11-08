package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

func SetupAnnouncementRoutes(
	r *gin.RouterGroup,
	announcementHandler *handlers.AnnouncementHandler,
	tokenMaker auth.Maker,
	rateLimitDefault int,
) {
	announcements := r.Group("/announcements")
	announcements.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes
		announcements.GET("", announcementHandler.ListAnnouncements)
		announcements.GET("/:id", announcementHandler.GetAnnouncement)

		// Protected routes
		authenticated := announcements.Group("")
		authenticated.Use(middleware.AuthMiddleware(tokenMaker))
		{
			authenticated.POST("", announcementHandler.CreateAnnouncement)
			authenticated.PUT("/:id", announcementHandler.UpdateAnnouncement)
			authenticated.PUT("/:id/status", announcementHandler.UpdateAnnouncementStatus)
		}
	}
}
