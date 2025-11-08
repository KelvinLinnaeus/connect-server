package announcements

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)



type CreateAnnouncementRequest struct {
	SpaceID        uuid.UUID              `json:"space_id" binding:"required"`
	Title          string                 `json:"title" binding:"required,min=3,max=200"`
	Content        string                 `json:"content" binding:"required,min=1"`
	Type           string                 `json:"type" binding:"required"` 
	TargetAudience []string               `json:"target_audience,omitempty"` 
	Priority       *string                `json:"priority,omitempty"` 
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"` 
	IsPinned       *bool                  `json:"is_pinned,omitempty"`
	AuthorID       uuid.UUID              
}

type UpdateAnnouncementRequest struct {
	Title          string                 `json:"title" binding:"required,min=3,max=200"`
	Content        string                 `json:"content" binding:"required,min=1"`
	Type           string                 `json:"type" binding:"required"`
	TargetAudience []string               `json:"target_audience,omitempty"`
	Priority       *string                `json:"priority,omitempty"`
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	IsPinned       *bool                  `json:"is_pinned,omitempty"`
}

type UpdateAnnouncementStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published archived"`
}

type ListAnnouncementsParams struct {
	SpaceID        uuid.UUID
	TargetAudience []string
	Page           int32
	Limit          int32
}



type AnnouncementResponse struct {
	ID             uuid.UUID              `json:"id"`
	SpaceID        uuid.UUID              `json:"space_id"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	Type           string                 `json:"type"`
	TargetAudience []string               `json:"target_audience,omitempty"`
	Priority       *string                `json:"priority,omitempty"`
	Status         *string                `json:"status,omitempty"`
	AuthorID       *uuid.UUID             `json:"author_id,omitempty"`
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	IsPinned       *bool                  `json:"is_pinned,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
}

type AnnouncementDetailResponse struct {
	ID             uuid.UUID              `json:"id"`
	SpaceID        uuid.UUID              `json:"space_id"`
	Title          string                 `json:"title"`
	Content        string                 `json:"content"`
	Type           string                 `json:"type"`
	TargetAudience []string               `json:"target_audience,omitempty"`
	Priority       *string                `json:"priority,omitempty"`
	Status         *string                `json:"status,omitempty"`
	ScheduledFor   *time.Time             `json:"scheduled_for,omitempty"`
	ExpiresAt      *time.Time             `json:"expires_at,omitempty"`
	Attachments    *pqtype.NullRawMessage `json:"attachments,omitempty"`
	IsPinned       *bool                  `json:"is_pinned,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`

	
	AuthorID       *uuid.UUID `json:"author_id,omitempty"`
	AuthorUsername string     `json:"author_username"`
	AuthorFullName string     `json:"author_full_name"`
	AuthorAvatar   *string    `json:"author_avatar,omitempty"`
}

type AnnouncementListResponse struct {
	ID             uuid.UUID  `json:"id"`
	SpaceID        uuid.UUID  `json:"space_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	TargetAudience []string   `json:"target_audience,omitempty"`
	Priority       *string    `json:"priority,omitempty"`
	IsPinned       *bool      `json:"is_pinned,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`

	
	AuthorUsername string `json:"author_username"`
	AuthorFullName string `json:"author_full_name"`
}
