package handlers

import (
	"fmt"
	"net/http"

	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/service/notifications"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *notifications.Service
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *notifications.Service) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// CreateNotification handles POST /api/notifications
func (h *NotificationHandler) CreateNotification(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	
	var req notifications.CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}
	
	// Set from_user_id from auth payload if not provided
	if req.FromUserID == nil {
		fromUserID, _ := uuid.Parse(authPayload.UserID)
		req.FromUserID = &fromUserID
	}
	
	notification, err := h.notificationService.CreateNotification(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, util.NewSuccessResponse(notification))
}

// GetUserNotifications handles GET /api/notifications
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	// Parse pagination parameters
	limit := int32(20)
	offset := int32(0)
	if l := c.Query("limit"); l != "" {
		if parsed, err := parseInt32(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := parseInt32(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	notifs, err := h.notificationService.GetUserNotifications(c.Request.Context(), notifications.GetNotificationsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(notifs))
}

// MarkAsRead handles PUT /api/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid notification ID"))
		return
	}

	if err := h.notificationService.MarkNotificationAsRead(c.Request.Context(), notificationID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Notification marked as read"}))
}

// MarkAllAsRead handles PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "All notifications marked as read"}))
}

// DeleteNotification handles DELETE /api/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", "Invalid notification ID"))
		return
	}

	if err := h.notificationService.DeleteNotification(c.Request.Context(), notificationID); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]string{"message": "Notification deleted"}))
}

// GetUnreadCount handles GET /api/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	payload, exists := c.Get("authorization_payload")
	if !exists {
		c.JSON(http.StatusUnauthorized, util.NewErrorResponse("unauthorized", "Not authenticated"))
		return
	}
	authPayload := payload.(*auth.Payload)
	userID, _ := uuid.Parse(authPayload.UserID)

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(map[string]int64{"count": count}))
}

// Helper function to parse int32 from string
func parseInt32(s string) (int32, error) {
	var result int32
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
