package notifications

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)


type CreateNotificationRequest struct {
	ToUserID       uuid.UUID              `json:"to_user_id" binding:"required"`
	FromUserID     *uuid.UUID             `json:"from_user_id,omitempty"`
	Type           string                 `json:"type" binding:"required"` 
	Title          *string                `json:"title,omitempty"`
	Message        *string                `json:"message,omitempty"`
	RelatedID      *uuid.UUID             `json:"related_id,omitempty"` 
	Metadata       *pqtype.NullRawMessage `json:"metadata,omitempty"`
	Priority       string                 `json:"priority,omitempty"` 
	ActionRequired bool                   `json:"action_required,omitempty"`
}


type NotificationResponse struct {
	ID             uuid.UUID              `json:"id"`
	ToUserID       uuid.UUID              `json:"to_user_id"`
	FromUserID     *uuid.UUID             `json:"from_user_id,omitempty"`
	Type           string                 `json:"type"`
	Title          *string                `json:"title,omitempty"`
	Message        *string                `json:"message,omitempty"`
	RelatedID      *uuid.UUID             `json:"related_id,omitempty"`
	Metadata       *pqtype.NullRawMessage `json:"metadata,omitempty"`
	IsRead         bool                   `json:"is_read"`
	Priority       string                 `json:"priority"`
	ActionRequired bool                   `json:"action_required"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
}


type NotificationWithUserResponse struct {
	ID             uuid.UUID              `json:"id"`
	ToUserID       uuid.UUID              `json:"to_user_id"`
	FromUserID     *uuid.UUID             `json:"from_user_id,omitempty"`
	FromUsername   *string                `json:"from_username,omitempty"`
	FromFullName   *string                `json:"from_full_name,omitempty"`
	FromAvatar     *string                `json:"from_avatar,omitempty"`
	Type           string                 `json:"type"`
	Title          *string                `json:"title,omitempty"`
	Message        *string                `json:"message,omitempty"`
	RelatedID      *uuid.UUID             `json:"related_id,omitempty"`
	Metadata       *pqtype.NullRawMessage `json:"metadata,omitempty"`
	IsRead         bool                   `json:"is_read"`
	Priority       string                 `json:"priority"`
	ActionRequired bool                   `json:"action_required"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
}


type GetNotificationsParams struct {
	UserID   uuid.UUID
	IsRead   *bool 
	Page     int32
	Limit    int32
	Offset   int32
	Priority *string 
}
