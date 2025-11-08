package communities

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateCommunityRequest represents the request to create a new community
type CreateCommunityRequest struct {
	SpaceID     uuid.UUID              `json:"space_id" binding:"required"`
	Name        string                 `json:"name" binding:"required,min=1,max=100"`
	Description *string                `json:"description,omitempty"`
	Category    string                 `json:"category" binding:"required"`
	CoverImage  *string                `json:"cover_image,omitempty"`
	IsPublic    *bool                  `json:"is_public,omitempty"`
	CreatedBy   uuid.UUID              `json:"created_by"`
	Settings    *pqtype.NullRawMessage `json:"settings,omitempty"`
}

// UpdateCommunityRequest represents the request to update a community
type UpdateCommunityRequest struct {
	Name        string                 `json:"name" binding:"required,min=1,max=100"`
	Description *string                `json:"description,omitempty"`
	CoverImage  *string                `json:"cover_image,omitempty"`
	Category    string                 `json:"category" binding:"required"`
	IsPublic    *bool                  `json:"is_public,omitempty"`
	Settings    *pqtype.NullRawMessage `json:"settings,omitempty"`
}

// AddModeratorRequest represents the request to add a moderator
type AddModeratorRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Permissions []string  `json:"permissions,omitempty"`
}

// CommunityResponse represents a basic community response
type CommunityResponse struct {
	ID          uuid.UUID              `json:"id"`
	SpaceID     uuid.UUID              `json:"space_id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description,omitempty"`
	Category    string                 `json:"category"`
	CoverImage  *string                `json:"cover_image,omitempty"`
	MemberCount int32                  `json:"member_count"`
	Status      string                 `json:"status"`
	PostCount   int32                  `json:"post_count"`
	IsPublic    bool                   `json:"is_public"`
	CreatedBy   *uuid.UUID             `json:"created_by,omitempty"`
	Settings    *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt   *time.Time             `json:"created_at,omitempty"`
	UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
}

// CommunityDetailResponse represents detailed community information
type CommunityDetailResponse struct {
	ID                uuid.UUID              `json:"id"`
	SpaceID           uuid.UUID              `json:"space_id"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description,omitempty"`
	Category          string                 `json:"category"`
	CoverImage        *string                `json:"cover_image,omitempty"`
	MemberCount       int32                  `json:"member_count"`
	Status            string                 `json:"status"`
	PostCount         int32                  `json:"post_count"`
	IsPublic          bool                   `json:"is_public"`
	CreatedBy         *uuid.UUID             `json:"created_by,omitempty"`
	Settings          *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	CreatedByUsername *string                `json:"created_by_username,omitempty"`
	CreatedByFullName *string                `json:"created_by_full_name,omitempty"`
	IsMember          bool                   `json:"is_member"`
	UserRole          *string                `json:"user_role,omitempty"`
	ActualMemberCount *int64                 `json:"actual_member_count,omitempty"`
	ActualPostCount   *int64                 `json:"actual_post_count,omitempty"`
}

// CommunityListResponse represents a community in list view
type CommunityListResponse struct {
	ID                uuid.UUID              `json:"id"`
	SpaceID           uuid.UUID              `json:"space_id"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description,omitempty"`
	Category          string                 `json:"category"`
	CoverImage        *string                `json:"cover_image,omitempty"`
	MemberCount       int32                  `json:"member_count"`
	Status            string                 `json:"status"`
	PostCount         int32                  `json:"post_count"`
	IsPublic          bool                   `json:"is_public"`
	CreatedBy         *uuid.UUID             `json:"created_by,omitempty"`
	Settings          *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	UserRole          *string                `json:"user_role,omitempty"`
	IsMember          bool                   `json:"is_member"`
	ActualMemberCount int64                  `json:"actual_member_count"`
}

// UserCommunityResponse represents a user's community membership
type UserCommunityResponse struct {
	ID          uuid.UUID  `json:"id"`
	SpaceID     uuid.UUID  `json:"space_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Category    string     `json:"category"`
	CoverImage  *string    `json:"cover_image,omitempty"`
	MemberCount int32      `json:"member_count"`
	Status      string     `json:"status"`
	PostCount   int32      `json:"post_count"`
	IsPublic    bool       `json:"is_public"`
	UserRole    string     `json:"user_role"`
	JoinedAt    *time.Time `json:"joined_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// CommunityMemberResponse represents a community member
type CommunityMemberResponse struct {
	ID         uuid.UUID  `json:"id"`
	Username   string     `json:"username"`
	FullName   string     `json:"full_name"`
	Avatar     *string    `json:"avatar,omitempty"`
	Level      *string    `json:"level,omitempty"`
	Department *string    `json:"department,omitempty"`
	Verified   bool       `json:"verified"`
	Role       string     `json:"role"`
	JoinedAt   *time.Time `json:"joined_at,omitempty"`
}

// CommunityModeratorResponse represents a community moderator
type CommunityModeratorResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Avatar      *string   `json:"avatar,omitempty"`
	Permissions []string  `json:"permissions"`
}

// CommunityAdminResponse represents a community admin
type CommunityAdminResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Avatar      *string   `json:"avatar,omitempty"`
	Permissions []string  `json:"permissions"`
}

// CommunityMembershipResponse represents a membership action result
type CommunityMembershipResponse struct {
	ID          uuid.UUID  `json:"id"`
	CommunityID uuid.UUID  `json:"community_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Role        string     `json:"role"`
	Permissions []string   `json:"permissions"`
	JoinedAt    *time.Time `json:"joined_at,omitempty"`
}

// ListCommunitiesParams represents parameters for listing communities
type ListCommunitiesParams struct {
	UserID  uuid.UUID `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
	SortBy  string    `json:"sort_by"` // "members", "posts", "recent"
	Page    int32     `json:"page"`
	Limit   int32     `json:"limit"`
}

// SearchCommunitiesParams represents parameters for searching communities
type SearchCommunitiesParams struct {
	UserID  uuid.UUID
	SpaceID uuid.UUID
	Query   string
}
