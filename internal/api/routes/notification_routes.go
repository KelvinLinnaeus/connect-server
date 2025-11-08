package routes

import (
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
)


func SetupNotificationRoutes(r *gin.RouterGroup, notificationHandler *handlers.NotificationHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	
	notifications := r.Group("/notifications")
	notifications.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	notifications.Use(middleware.AuthMiddleware(tokenMaker))
	{
		
		notifications.POST("", notificationHandler.CreateNotification)
		
		
		notifications.GET("", notificationHandler.GetUserNotifications)
		notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
		notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
		notifications.DELETE("/:id", notificationHandler.DeleteNotification)
		notifications.GET("/unread-count", notificationHandler.GetUnreadCount)
	}
}
