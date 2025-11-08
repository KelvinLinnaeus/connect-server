package routes

import (
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
)


func SetupMessagingRoutes(r *gin.RouterGroup, messagingHandler *handlers.MessagingHandler, tokenMaker auth.Maker, rateLimitDefault int) {
	
	conversations := r.Group("/conversations")
	conversations.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	conversations.Use(middleware.AuthMiddleware(tokenMaker))
	{
		
		conversations.POST("", messagingHandler.CreateConversation)
		conversations.GET("", messagingHandler.GetUserConversations)
		conversations.GET("/:id", messagingHandler.GetConversation)
		conversations.POST("/direct", messagingHandler.GetOrCreateDirectConversation)
		conversations.POST("/:id/leave", messagingHandler.LeaveConversation)
		conversations.PUT("/:id/settings", messagingHandler.UpdateParticipantSettings)
		
		
		conversations.GET("/:id/participants", messagingHandler.GetConversationParticipants)
		conversations.POST("/:id/participants", messagingHandler.AddConversationParticipants)
		
		
		conversations.POST("/:id/messages", messagingHandler.SendMessage)
		conversations.GET("/:id/messages", messagingHandler.GetConversationMessages)
		conversations.POST("/:id/read", messagingHandler.MarkMessagesAsRead)
		conversations.GET("/:id/unread", messagingHandler.GetUnreadCount)
	}
	
	
	messages := r.Group("/messages")
	messages.Use(middleware.RateLimitMiddleware(rateLimitDefault))
	messages.Use(middleware.AuthMiddleware(tokenMaker))
	{
		messages.GET("/:id", messagingHandler.GetMessage)
		messages.DELETE("/:id", messagingHandler.DeleteMessage)
		messages.POST("/:id/reactions", messagingHandler.AddMessageReaction)
		messages.DELETE("/:id/reactions/:emoji", messagingHandler.RemoveMessageReaction)
	}
}
