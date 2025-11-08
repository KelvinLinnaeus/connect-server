package notifications

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/live"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/sqlc-dev/pqtype"
)

// Service handles notification business logic
type Service struct {
	store       db.Store
	liveService *live.Service
}

// NewService creates a new notification service
func NewService(store db.Store, liveService *live.Service) *Service {
	return &Service{
		store:       store,
		liveService: liveService,
	}
}

// CreateNotification creates a new notification
func (s *Service) CreateNotification(ctx context.Context, req CreateNotificationRequest) (*NotificationResponse, error) {
	var fromUserID, relatedID uuid.NullUUID
	var title, message, priority sql.NullString
	var metadata pqtype.NullRawMessage
	var actionRequired sql.NullBool
	
	if req.FromUserID != nil {
		fromUserID = uuid.NullUUID{UUID: *req.FromUserID, Valid: true}
	}
	if req.RelatedID != nil {
		relatedID = uuid.NullUUID{UUID: *req.RelatedID, Valid: true}
	}
	if req.Title != nil {
		title = sql.NullString{String: *req.Title, Valid: true}
	}
	if req.Message != nil {
		message = sql.NullString{String: *req.Message, Valid: true}
	}
	if req.Priority != "" {
		priority = sql.NullString{String: req.Priority, Valid: true}
	} else {
		priority = sql.NullString{String: "normal", Valid: true}
	}
	if req.Metadata != nil {
		metadata = *req.Metadata
	}
	actionRequired = sql.NullBool{Bool: req.ActionRequired, Valid: true}
	
	notification, err := s.store.CreateNotification(ctx, db.CreateNotificationParams{
		ToUserID:       req.ToUserID,
		FromUserID:     fromUserID,
		Type:           req.Type,
		Title:          title,
		Message:        message,
		RelatedID:      relatedID,
		Metadata:       metadata,
		Priority:       priority,
		ActionRequired: actionRequired,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	response := s.toNotificationResponse(notification)

	// Publish real-time event for notification creation
	if s.liveService != nil {
		notificationPayload := map[string]interface{}{
			"id":              notification.ID.String(),
			"to_user_id":      notification.ToUserID.String(),
			"type":            notification.Type,
			"is_read":         response.IsRead,
			"priority":        response.Priority,
			"action_required": response.ActionRequired,
		}
		if response.FromUserID != nil {
			notificationPayload["from_user_id"] = response.FromUserID.String()
		}
		if response.Title != nil {
			notificationPayload["title"] = *response.Title
		}
		if response.Message != nil {
			notificationPayload["message"] = *response.Message
		}
		if response.RelatedID != nil {
			notificationPayload["related_id"] = response.RelatedID.String()
		}
		if response.Metadata != nil {
			notificationPayload["metadata"] = response.Metadata
		}
		if response.CreatedAt != nil {
			notificationPayload["created_at"] = response.CreatedAt.Unix()
		}

		if err := s.liveService.PublishNotificationCreated(ctx, req.ToUserID, notificationPayload); err != nil {
			log.Error().Err(err).Msg("Failed to publish notification.created event")
		}
	}

	return response, nil
}

// GetUserNotifications gets notifications for a user
func (s *Service) GetUserNotifications(ctx context.Context, params GetNotificationsParams) ([]NotificationWithUserResponse, error) {
	notifications, err := s.store.GetUserNotifications(ctx, db.GetUserNotificationsParams{
		ToUserID: params.UserID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	responses := make([]NotificationWithUserResponse, len(notifications))
	for i, n := range notifications {
		responses[i] = s.rowToNotificationResponse(n)
	}

	return responses, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *Service) MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) error {
	if err := s.store.MarkAsRead(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *Service) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.store.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

// DeleteNotification deletes a notification
func (s *Service) DeleteNotification(ctx context.Context, notificationID uuid.UUID) error {
	if err := s.store.DeleteNotification(ctx, notificationID); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// GetUnreadCount gets unread notification count
func (s *Service) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.store.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

// Helper conversion functions

func (s *Service) rowToNotificationResponse(row db.GetUserNotificationsRow) NotificationWithUserResponse {
	resp := NotificationWithUserResponse{
		ID:       row.ID,
		ToUserID: row.ToUserID,
		Type:     row.Type,
	}

	if row.FromUserID.Valid {
		resp.FromUserID = &row.FromUserID.UUID
	}
	if row.Title.Valid {
		resp.Title = &row.Title.String
	}
	if row.Message.Valid {
		resp.Message = &row.Message.String
	}
	if row.RelatedID.Valid {
		resp.RelatedID = &row.RelatedID.UUID
	}
	if row.Metadata.Valid {
		resp.Metadata = &row.Metadata
	}
	if row.IsRead.Valid {
		resp.IsRead = row.IsRead.Bool
	}
	if row.Priority.Valid {
		resp.Priority = row.Priority.String
	}
	if row.ActionRequired.Valid {
		resp.ActionRequired = row.ActionRequired.Bool
	}
	if row.CreatedAt.Valid {
		resp.CreatedAt = &row.CreatedAt.Time
	}
	if row.FromUsername.Valid {
		resp.FromUsername = &row.FromUsername.String
	}
	if row.FromFullName.Valid {
		resp.FromFullName = &row.FromFullName.String
	}
	if row.FromAvatar.Valid {
		resp.FromAvatar = &row.FromAvatar.String
	}

	return resp
}

func (s *Service) toNotificationResponse(n db.Notification) *NotificationResponse {
	resp := &NotificationResponse{
		ID:       n.ID,
		ToUserID: n.ToUserID,
		Type:     n.Type,
	}
	
	if n.FromUserID.Valid {
		resp.FromUserID = &n.FromUserID.UUID
	}
	if n.Title.Valid {
		resp.Title = &n.Title.String
	}
	if n.Message.Valid {
		resp.Message = &n.Message.String
	}
	if n.RelatedID.Valid {
		resp.RelatedID = &n.RelatedID.UUID
	}
	if n.Metadata.Valid {
		resp.Metadata = &n.Metadata
	}
	if n.IsRead.Valid {
		resp.IsRead = n.IsRead.Bool
	}
	if n.Priority.Valid {
		resp.Priority = n.Priority.String
	}
	if n.ActionRequired.Valid {
		resp.ActionRequired = n.ActionRequired.Bool
	}
	if n.CreatedAt.Valid {
		resp.CreatedAt = &n.CreatedAt.Time
	}
	
	return resp
}
