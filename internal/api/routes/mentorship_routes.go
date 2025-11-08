package routes

import (
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
)

func SetupMentorshipRoutes(
	r *gin.RouterGroup,
	mentorshipHandler *handlers.MentorshipHandler,
	tokenMaker auth.Maker,
	rateLimitDefault int,
) {
	mentorship := r.Group("/mentorship")
	mentorship.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	{
		
		
		
		mentors := mentorship.Group("/mentors")
		{
			
			mentors.GET("/search", mentorshipHandler.SearchMentors)
			mentors.GET("/profile/:id", mentorshipHandler.GetMentorProfile)
			mentors.GET("/:id/reviews", mentorshipHandler.GetMentorReviews)

			
			mentorsAuth := mentors.Group("")
			mentorsAuth.Use(middleware.AuthMiddleware(tokenMaker))
			{
				mentorsAuth.POST("/profile", mentorshipHandler.CreateMentorProfile)
				mentorsAuth.PUT("/profile/:id/availability", mentorshipHandler.UpdateMentorAvailability)
				mentorsAuth.GET("/my-profile", mentorshipHandler.GetMyMentorProfile)
				mentorsAuth.GET("/recommended", mentorshipHandler.GetRecommendedMentors)
			}

			
			applications := mentors.Group("/applications")
			{
				
				applications.GET("/:id", mentorshipHandler.GetMentorApplication)
				applications.GET("/pending", mentorshipHandler.GetPendingMentorApplications)

				
				applicationsAuth := applications.Group("")
				applicationsAuth.Use(middleware.AuthMiddleware(tokenMaker))
				{
					applicationsAuth.POST("", mentorshipHandler.CreateMentorApplication)
					applicationsAuth.PUT("/:id", mentorshipHandler.UpdateMentorApplication)
					applicationsAuth.GET("/my-application", mentorshipHandler.GetMyMentorApplication)
				}
			}
		}

		
		
		
		tutors := mentorship.Group("/tutors")
		{
			
			tutors.GET("/search", mentorshipHandler.SearchTutors)
			tutors.GET("/profile/:id", mentorshipHandler.GetTutorProfile)
			tutors.GET("/:id/reviews", mentorshipHandler.GetTutorReviews)

			
			tutorsAuth := tutors.Group("")
			tutorsAuth.Use(middleware.AuthMiddleware(tokenMaker))
			{
				tutorsAuth.POST("/profile", mentorshipHandler.CreateTutorProfile)
				tutorsAuth.PUT("/profile/:id/availability", mentorshipHandler.UpdateTutorAvailability)
				tutorsAuth.GET("/my-profile", mentorshipHandler.GetMyTutorProfile)
				tutorsAuth.GET("/recommended", mentorshipHandler.GetRecommendedTutors)
			}

			
			applications := tutors.Group("/applications")
			{
				
				applications.GET("/:id", mentorshipHandler.GetTutorApplication)
				applications.GET("/pending", mentorshipHandler.GetPendingTutorApplications)

				
				applicationsAuth := applications.Group("")
				applicationsAuth.Use(middleware.AuthMiddleware(tokenMaker))
				{
					applicationsAuth.POST("", mentorshipHandler.CreateTutorApplication)
					applicationsAuth.PUT("/:id", mentorshipHandler.UpdateTutorApplication)
					applicationsAuth.GET("/my-application", mentorshipHandler.GetMyTutorApplication)
				}
			}
		}

		
		
		
		mentoringSessions := mentorship.Group("/mentoring/sessions")
		{
			
			mentoringSessions.GET("/:id", mentorshipHandler.GetMentoringSession)

			
			mentoringAuth := mentoringSessions.Group("")
			mentoringAuth.Use(middleware.AuthMiddleware(tokenMaker))
			{
				mentoringAuth.POST("", mentorshipHandler.CreateMentoringSession)
				mentoringAuth.GET("", mentorshipHandler.GetUserMentoringSessions)
				mentoringAuth.PUT("/:id/status", mentorshipHandler.UpdateMentoringSessionStatus)
				mentoringAuth.PUT("/:id/meeting-link", mentorshipHandler.AddMentoringSessionMeetingLink)
				mentoringAuth.POST("/:id/rate", mentorshipHandler.RateMentoringSession)
			}
		}

		
		
		
		tutoringSessions := mentorship.Group("/tutoring/sessions")
		{
			
			tutoringSessions.GET("/:id", mentorshipHandler.GetTutoringSession)

			
			tutoringAuth := tutoringSessions.Group("")
			tutoringAuth.Use(middleware.AuthMiddleware(tokenMaker))
			{
				tutoringAuth.POST("", mentorshipHandler.CreateTutoringSession)
				tutoringAuth.GET("", mentorshipHandler.GetUserTutoringSessions)
				tutoringAuth.PUT("/:id/status", mentorshipHandler.UpdateTutoringSessionStatus)
				tutoringAuth.PUT("/:id/meeting-link", mentorshipHandler.AddTutoringSessionMeetingLink)
				tutoringAuth.POST("/:id/rate", mentorshipHandler.RateTutoringSession)
			}
		}
	}
}
