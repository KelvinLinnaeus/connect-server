package messaging

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateConversationRequest represents the request to create a new conversation
type CreateConversationRequest struct {
	SpaceID          uuid.UUID              `json:"space_id" binding:"required"`
	Name             *string                `json:"name,omitempty"`
	Avatar           *string                `json:"avatar,omitempty"`
	Description      *string                `json:"description,omitempty"`
	ConversationType string                 `json:"conversation_type" binding:"required"` // direct, group, channel
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
	ParticipantIDs   []uuid.UUID            `json:"participant_ids" binding:"required,min=1"`
}

// SendMessageRequest represents the request to send a message
type SendMessageRequest struct {
	ConversationID uuid.UUID              `json:"conversation_id" binding:"required"`
	SenderID       uuid.UUID              `json:"sender_id"`
	RecipientID    *uuid.UUID             `json:"recipient_id,omitempty"`
	Content        string                 `json:"content" binding:"required,min=1"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	MessageType    string                 `json:"message_type,omitempty"` // text, image, file, audio, video
	ReplyToID      *uuid.UUID             `json:"reply_to_id,omitempty"`
}

// UpdateParticipantSettingsRequest represents participant preferences
type UpdateParticipantSettingsRequest struct {
	NotificationsEnabled *bool                  `json:"notifications_enabled,omitempty"`
	CustomSettings       *pqtype.NullRawMessage `json:"custom_settings,omitempty"`
}

// AddReactionRequest represents adding a reaction to a message
type AddReactionRequest struct {
	Emoji  string    `json:"emoji" binding:"required"`
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// ConversationResponse represents a conversation
type ConversationResponse struct {
	ID               uuid.UUID              `json:"id"`
	SpaceID          uuid.UUID              `json:"space_id"`
	Name             *string                `json:"name,omitempty"`
	Avatar           *string                `json:"avatar,omitempty"`
	Description      *string                `json:"description,omitempty"`
	ConversationType string                 `json:"conversation_type"`
	LastMessageID    *uuid.UUID             `json:"last_message_id,omitempty"`
	LastMessageAt    *time.Time             `json:"last_message_at,omitempty"`
	IsActive         bool                   `json:"is_active"`
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	UnreadCount      int64                  `json:"unread_count"`
}

// ConversationDetailResponse represents detailed conversation info
type ConversationDetailResponse struct {
	ID                   uuid.UUID              `json:"id"`
	SpaceID              uuid.UUID              `json:"space_id"`
	Name                 *string                `json:"name,omitempty"`
	Avatar               *string                `json:"avatar,omitempty"`
	Description          *string                `json:"description,omitempty"`
	ConversationType     string                 `json:"conversation_type"`
	LastMessageID        *uuid.UUID             `json:"last_message_id,omitempty"`
	LastMessageAt        *time.Time             `json:"last_message_at,omitempty"`
	IsActive             bool                   `json:"is_active"`
	Settings             *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt            *time.Time             `json:"created_at,omitempty"`
	UpdatedAt            *time.Time             `json:"updated_at,omitempty"`
	UnreadCount          int64                  `json:"unread_count"`
	UserRole             *string                `json:"user_role,omitempty"`
	NotificationsEnabled bool                   `json:"notifications_enabled"`
	LastMessageContent   *string                `json:"last_message_content,omitempty"`
	LastMessageTime      *time.Time             `json:"last_message_time,omitempty"`
	LastSenderUsername   *string                `json:"last_sender_username,omitempty"`
	LastSenderFullName   *string                `json:"last_sender_full_name,omitempty"`
}

// MessageResponse represents a message
type MessageResponse struct {
	ID             uuid.UUID              `json:"id"`
	ConversationID uuid.UUID              `json:"conversation_id"`
	SenderID       uuid.UUID              `json:"sender_id"`
	RecipientID    *uuid.UUID             `json:"recipient_id,omitempty"`
	Content        string                 `json:"content"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	MessageType    string                 `json:"message_type"`
	IsRead         bool                   `json:"is_read"`
	ReadAt         *time.Time             `json:"read_at,omitempty"`
	Reactions      *pqtype.NullRawMessage `json:"reactions,omitempty"`
	ReplyToID      *uuid.UUID             `json:"reply_to_id,omitempty"`
	Status         string                 `json:"status"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	SenderUsername string                 `json:"sender_username"`
	SenderFullName string                 `json:"sender_full_name"`
	SenderAvatar   *string                `json:"sender_avatar,omitempty"`
}

// MessageDetailResponse represents a message with reply context
type MessageDetailResponse struct {
	ID             uuid.UUID              `json:"id"`
	ConversationID uuid.UUID              `json:"conversation_id"`
	SenderID       uuid.UUID              `json:"sender_id"`
	RecipientID    *uuid.UUID             `json:"recipient_id,omitempty"`
	Content        string                 `json:"content"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	MessageType    string                 `json:"message_type"`
	IsRead         bool                   `json:"is_read"`
	ReadAt         *time.Time             `json:"read_at,omitempty"`
	Reactions      *pqtype.NullRawMessage `json:"reactions,omitempty"`
	ReplyToID      *uuid.UUID             `json:"reply_to_id,omitempty"`
	Status         string                 `json:"status"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	SenderUsername string                 `json:"sender_username"`
	SenderFullName string                 `json:"sender_full_name"`
	SenderAvatar   *string                `json:"sender_avatar,omitempty"`
	ReplyContent   *string                `json:"reply_content,omitempty"`
	ReplyUsername  *string                `json:"reply_username,omitempty"`
}

// ParticipantResponse represents a conversation participant
type ParticipantResponse struct {
	ID                   uuid.UUID  `json:"id"`
	Username             string     `json:"username"`
	FullName             string     `json:"full_name"`
	Avatar               *string    `json:"avatar,omitempty"`
	Verified             bool       `json:"verified"`
	Role                 string     `json:"role"`
	JoinedAt             *time.Time `json:"joined_at,omitempty"`
	IsActive             bool       `json:"is_active"`
	NotificationsEnabled bool       `json:"notifications_enabled"`
}

// GetConversationMessagesParams represents parameters for getting messages
type GetConversationMessagesParams struct {
	ConversationID uuid.UUID
	Page           int32
	Limit          int32
}
