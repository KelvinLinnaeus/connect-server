package routes

import (
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(router *gin.RouterGroup, adminHandler *handlers.AdminHandler, tokenMaker auth.Maker) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(tokenMaker))
	{
		// User Management
		admin.GET("/users", adminHandler.GetUsers)
		admin.DELETE("/users/:id", adminHandler.DeleteUser)
		admin.PUT("/users/:id/suspend", adminHandler.SuspendUser)
		admin.PUT("/users/:id/unsuspend", adminHandler.UnsuspendUser)
		admin.PUT("/users/:id/ban", adminHandler.BanUser)

		// Content Moderation
		admin.GET("/reports", adminHandler.GetReports)
		admin.PUT("/reports/:id/resolve", adminHandler.ResolveReport)
		admin.PUT("/reports/:id/escalate", adminHandler.EscalateReport)

		// Application Management
		admin.GET("/applications/tutors", adminHandler.GetTutorApplications)
		admin.PUT("/applications/tutors/:id/approve", adminHandler.ApproveTutorApplication)
		admin.PUT("/applications/tutors/:id/reject", adminHandler.RejectTutorApplication)
		admin.GET("/applications/mentors", adminHandler.GetMentorApplications)
		admin.PUT("/applications/mentors/:id/approve", adminHandler.ApproveMentorApplication)
		admin.PUT("/applications/mentors/:id/reject", adminHandler.RejectMentorApplication)

		// Groups Management
		admin.GET("/groups", adminHandler.GetGroups)
		admin.PUT("/groups/:id/approve", adminHandler.ApproveGroup)
		admin.PUT("/groups/:id/reject", adminHandler.RejectGroup)
		admin.DELETE("/groups/:id", adminHandler.DeleteGroup)

		// Space Activities
		admin.GET("/spaces/:id/activities", adminHandler.GetSpaceActivities)

		// Dashboard
		admin.GET("/dashboard/stats", adminHandler.GetDashboardStats)

		// System Settings
		admin.GET("/settings", adminHandler.GetSettings)
		admin.PUT("/settings/:key", adminHandler.UpdateSetting)

		// Analytics
		admin.GET("/analytics/user-growth", adminHandler.GetUserGrowth)
		admin.GET("/analytics/engagement", adminHandler.GetEngagementMetrics)
		admin.GET("/analytics/activity", adminHandler.GetActivityAnalytics)

		// Admin Management
		admin.GET("/admins", adminHandler.GetAdmins)
		admin.PUT("/admins/:id/role", adminHandler.UpdateAdminRole)
		admin.PUT("/admins/:id/status", adminHandler.UpdateAdminStatus)

		// Notifications
		admin.GET("/notifications", adminHandler.GetNotifications)
		admin.PUT("/notifications/:id/read", adminHandler.MarkNotificationRead)
		admin.DELETE("/notifications/:id", adminHandler.DeleteNotification)
		admin.PUT("/notifications/read-all", adminHandler.MarkAllNotificationsRead)

		// Communities Management
		admin.GET("/communities", adminHandler.GetCommunities)
		admin.POST("/communities", adminHandler.CreateCommunity)
		admin.PUT("/communities/:id", adminHandler.UpdateCommunity)
		admin.DELETE("/communities/:id", adminHandler.DeleteCommunity)
		admin.PUT("/communities/:id/status", adminHandler.UpdateCommunityStatus)
		admin.POST("/communities/:id/moderators", adminHandler.AssignCommunityModerator)

		// Announcements Management
		admin.GET("/announcements", adminHandler.GetAnnouncements)
		admin.POST("/announcements", adminHandler.CreateAnnouncement)
		admin.PUT("/announcements/:id", adminHandler.UpdateAnnouncement)
		admin.DELETE("/announcements/:id", adminHandler.DeleteAnnouncement)
		admin.PUT("/announcements/:id/status", adminHandler.UpdateAnnouncementStatus)

		// Events Management
		admin.GET("/events", adminHandler.GetEvents)
		admin.POST("/events", adminHandler.CreateEvent)
		admin.PUT("/events/:id", adminHandler.UpdateEvent)
		admin.DELETE("/events/:id", adminHandler.DeleteEvent)
		admin.PUT("/events/:id/status", adminHandler.UpdateEventStatus)
		admin.GET("/events/:id/registrations", adminHandler.GetEventRegistrations)

		// Enhanced User Management
		admin.POST("/users", adminHandler.CreateUser)
		admin.PUT("/users/:id", adminHandler.UpdateUser)
		admin.POST("/users/:id/reset-password", adminHandler.ResetUserPassword)

		// Data Export
		admin.GET("/export/:dataType", adminHandler.ExportData)
	}
}
