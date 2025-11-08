package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoutes(
	r *gin.RouterGroup,
	analyticsHandler *handlers.AnalyticsHandler,
	tokenMaker auth.Maker,
	rateLimitDefault int,
) {
	analytics := r.Group("/analytics")
	analytics.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		// ========================================================================
		// Content Moderation & Reporting Routes
		// ========================================================================
		reports := analytics.Group("/reports")
		{
			// Public routes
			reports.GET("/:id", analyticsHandler.GetReport)
			reports.GET("/by-content", analyticsHandler.GetReportsByContent)
			reports.GET("/pending", analyticsHandler.GetPendingReports)

			// Protected routes
			reportsAuth := reports.Group("")
			reportsAuth.Use(middleware.AuthMiddleware(tokenMaker))
			{
				reportsAuth.POST("", analyticsHandler.CreateReport)
				reportsAuth.PUT("/:id", analyticsHandler.UpdateReport)
			}
		}

		moderation := analytics.Group("/moderation")
		{
			moderation.GET("/queue", analyticsHandler.GetModerationQueue)
			moderation.GET("/stats", analyticsHandler.GetContentModerationStats)
		}

		// ========================================================================
		// System & Space Metrics Routes
		// ========================================================================
		metrics := analytics.Group("/metrics")
		{
			metrics.GET("/system", analyticsHandler.GetSystemMetrics)
			metrics.GET("/space", analyticsHandler.GetSpaceStats)
		}

		// ========================================================================
		// Engagement & Activity Routes
		// ========================================================================
		engagement := analytics.Group("/engagement")
		{
			engagement.GET("/metrics", analyticsHandler.GetEngagementMetrics)
		}

		activity := analytics.Group("/activity")
		{
			activity.GET("/stats", analyticsHandler.GetUserActivityStats)
		}

		users := analytics.Group("/users")
		{
			users.GET("/growth", analyticsHandler.GetUserGrowth)
			users.GET("/ranking", analyticsHandler.GetUserEngagementRanking)
		}

		// ========================================================================
		// Top Content Routes
		// ========================================================================
		top := analytics.Group("/top")
		{
			top.GET("/posts", analyticsHandler.GetTopPosts)
			top.GET("/communities", analyticsHandler.GetTopCommunities)
			top.GET("/groups", analyticsHandler.GetTopGroups)
		}

		// ========================================================================
		// Mentorship Analytics Routes
		// ========================================================================
		mentorshipAnalytics := analytics.Group("/mentorship")
		{
			mentorshipAnalytics.GET("/mentoring", analyticsHandler.GetMentoringStats)
			mentorshipAnalytics.GET("/tutoring", analyticsHandler.GetTutoringStats)
			mentorshipAnalytics.GET("/industries", analyticsHandler.GetPopularIndustries)
			mentorshipAnalytics.GET("/subjects", analyticsHandler.GetPopularSubjects)
		}
	}
}
