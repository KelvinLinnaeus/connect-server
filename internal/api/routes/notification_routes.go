package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

// SetupNotificationRoutes sets up all notification-related routes
func SetupNotificationRoutes(r *gin.RouterGroup, notificationHandler *handlers.NotificationHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	// All notification routes require authentication
	notifications := r.Group("/notifications")
	notifications.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	notifications.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// Working endpoint
		notifications.POST("", notificationHandler.CreateNotification)
		
		// Placeholder endpoints (require SQLC query implementation)
		notifications.GET("", notificationHandler.GetUserNotifications)
		notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
	}
}
