package handlers

import (
	"net/http"
	"strconv"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/messaging"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// MessagingHandler handles messaging-related HTTP requests
type MessagingHandler struct {
	messagingService *messaging.Service
}

// NewMessagingHandler creates a new messaging handler
func NewMessagingHandler(messagingService *messaging.Service) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
	}
}

// CreateConversation handles POST /api/conversations
func (h *MessagingHandler) CreateConversation(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	var req messaging.CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	// Ensure creator is in participants
	creatorID, _ := uuid.Parse(authPayload.UserID)
	hasCreator := false
	for _, pid := range req.ParticipantIDs {
		if pid == creatorID {
			hasCreator = true
			break
		}
	}
	if !hasCreator {
		req.ParticipantIDs = append(req.ParticipantIDs, creatorID)
	}
	
	conversation, err := h.messagingService.CreateConversation(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(conversation))
}

// GetConversation handles GET /api/conversations/:id
func (h *MessagingHandler) GetConversation(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	conversation, err := h.messagingService.GetConversationByID(c.Request.Context(), conversationID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(conversation))
}

// GetUserConversations handles GET /api/conversations
func (h *MessagingHandler) GetUserConversations(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	conversations, err := h.messagingService.GetUserConversations(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(conversations))
}

// GetOrCreateDirectConversation handles POST /api/conversations/direct
func (h *MessagingHandler) GetOrCreateDirectConversation(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	var req struct {
		SpaceID    uuid.UUID `json:"space_id" binding:"required"`
		RecipientID uuid.UUID `json:"recipient_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	conversationID, err := h.messagingService.GetOrCreateDirectConversation(
		c.Request.Context(),
		req.SpaceID,
		userID,
		req.RecipientID,
	)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"conversation_id": conversationID}))
}

// SendMessage handles POST /api/conversations/:id/messages
func (h *MessagingHandler) SendMessage(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	var req messaging.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	senderID, _ := uuid.Parse(authPayload.UserID)
	req.ConversationID = conversationID
	req.SenderID = senderID
	
	message, err := h.messagingService.SendMessage(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(message))
}

// GetConversationMessages handles GET /api/conversations/:id/messages
func (h *MessagingHandler) GetConversationMessages(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	page, limit := parsePagination(c)
	
	params := messaging.GetConversationMessagesParams{
		ConversationID: conversationID,
		Page:           int32(page),
		Limit:          int32(limit),
	}

	messages, err := h.messagingService.GetConversationMessages(c.Request.Context(), params)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(messages))
}

// GetMessage handles GET /api/messages/:id
func (h *MessagingHandler) GetMessage(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid message ID format"))
		return
	}
	
	message, err := h.messagingService.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(message))
}

// DeleteMessage handles DELETE /api/messages/:id
func (h *MessagingHandler) DeleteMessage(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid message ID format"))
		return
	}
	
	senderID, _ := uuid.Parse(authPayload.UserID)
	
	err = h.messagingService.DeleteMessage(c.Request.Context(), messageID, senderID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Message deleted successfully"}))
}

// MarkMessagesAsRead handles POST /api/conversations/:id/read
func (h *MessagingHandler) MarkMessagesAsRead(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	err = h.messagingService.MarkMessagesAsRead(c.Request.Context(), conversationID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Messages marked as read"}))
}

// GetUnreadCount handles GET /api/conversations/:id/unread
func (h *MessagingHandler) GetUnreadCount(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	count, err := h.messagingService.GetUnreadMessageCount(c.Request.Context(), conversationID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"unread_count": count}))
}

// GetConversationParticipants handles GET /api/conversations/:id/participants
func (h *MessagingHandler) GetConversationParticipants(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	participants, err := h.messagingService.GetConversationParticipants(c.Request.Context(), conversationID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(participants))
}

// AddConversationParticipants handles POST /api/conversations/:id/participants
func (h *MessagingHandler) AddConversationParticipants(c *gin.Context) {
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	var req struct {
		UserIDs []uuid.UUID `json:"user_ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	err = h.messagingService.AddConversationParticipants(c.Request.Context(), conversationID, req.UserIDs)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Participants added successfully"}))
}

// LeaveConversation handles POST /api/conversations/:id/leave
func (h *MessagingHandler) LeaveConversation(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	err = h.messagingService.LeaveConversation(c.Request.Context(), conversationID, userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Left conversation successfully"}))
}

// UpdateParticipantSettings handles PUT /api/conversations/:id/settings
func (h *MessagingHandler) UpdateParticipantSettings(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	conversationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid conversation ID format"))
		return
	}
	
	var req messaging.UpdateParticipantSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	userID, _ := uuid.Parse(authPayload.UserID)
	
	err = h.messagingService.UpdateParticipantSettings(c.Request.Context(), conversationID, userID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Settings updated successfully"}))
}

// AddMessageReaction handles POST /api/messages/:id/reactions
func (h *MessagingHandler) AddMessageReaction(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid message ID format"))
		return
	}
	
	var req messaging.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	req.UserID, _ = uuid.Parse(authPayload.UserID)
	
	err = h.messagingService.AddMessageReaction(c.Request.Context(), messageID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Reaction added successfully"}))
}

// RemoveMessageReaction handles DELETE /api/messages/:id/reactions/:emoji
func (h *MessagingHandler) RemoveMessageReaction(c *gin.Context) {
	messageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid message ID format"))
		return
	}
	
	emoji := c.Param("emoji")
	if emoji == "" {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("missing_emoji", "Emoji is required"))
		return
	}
	
	err = h.messagingService.RemoveMessageReaction(c.Request.Context(), messageID, emoji)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, util.NewSuccessResponse(gin.H{"message": "Reaction removed successfully"}))
}

// Helper function to parse pagination
func parseMessagePagination(c *gin.Context) (int32, int32) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	
	return int32(page), int32(limit)
}
