package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

func SetupEventRoutes(
	r *gin.RouterGroup,
	eventHandler *handlers.EventHandler,
	tokenMaker auth.Maker,
	rateLimitDefault int,
) {
	events := r.Group("/events")
	events.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// Public routes (optional auth for personalized data)
		events.GET("", eventHandler.ListEvents)
		events.GET("/:id", eventHandler.GetEvent)
		events.GET("/upcoming", eventHandler.GetUpcomingEvents)
		events.GET("/search", eventHandler.SearchEvents)
		events.GET("/categories", eventHandler.GetEventCategories)
		events.GET("/:id/attendees", eventHandler.GetEventAttendees)
		events.GET("/:id/co-organizers", eventHandler.GetEventCoOrganizers)

		// Protected routes
		authenticated := events.Group("")
		authenticated.Use(middleware.AuthMiddleware(tokenMaker))
		{
			// Event management
			authenticated.POST("", eventHandler.CreateEvent)
			authenticated.PUT("/:id", eventHandler.UpdateEvent)
			authenticated.PUT("/:id/status", eventHandler.UpdateEventStatus)

			// Registration
			authenticated.POST("/:id/register", eventHandler.RegisterForEvent)
			authenticated.POST("/:id/unregister", eventHandler.UnregisterFromEvent)

			// Co-organizers
			authenticated.POST("/:id/co-organizers", eventHandler.AddEventCoOrganizer)
			authenticated.DELETE("/:id/co-organizers/:user_id", eventHandler.RemoveEventCoOrganizer)

			// Attendance
			authenticated.POST("/:id/attendance/:user_id", eventHandler.MarkEventAttendance)
		}
	}

	// User-specific event routes
	users := r.Group("/users")
	users.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	users.Use(middleware.AuthMiddleware(tokenMaker))
	{
		users.GET("/events", eventHandler.GetUserEvents)
	}
}
