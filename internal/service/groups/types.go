package groups

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)


type CreateGroupRequest struct {
	SpaceID          uuid.UUID              `json:"space_id" binding:"required"`
	CommunityID      *uuid.UUID             `json:"community_id,omitempty"`
	Name             string                 `json:"name" binding:"required,min=1,max=100"`
	Description      *string                `json:"description,omitempty"`
	Category         string                 `json:"category" binding:"required"`
	GroupType        string                 `json:"group_type" binding:"required"` 
	Avatar           *string                `json:"avatar,omitempty"`
	Banner           *string                `json:"banner,omitempty"`
	AllowInvites     *bool                  `json:"allow_invites,omitempty"`
	AllowMemberPosts *bool                  `json:"allow_member_posts,omitempty"`
	CreatedBy        uuid.UUID              `json:"created_by"`
	Tags             []string               `json:"tags,omitempty"`
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
}


type UpdateGroupRequest struct {
	Name             string                 `json:"name" binding:"required,min=1,max=100"`
	Description      *string                `json:"description,omitempty"`
	Category         string                 `json:"category" binding:"required"`
	Avatar           *string                `json:"avatar,omitempty"`
	Banner           *string                `json:"banner,omitempty"`
	AllowInvites     *bool                  `json:"allow_invites,omitempty"`
	AllowMemberPosts *bool                  `json:"allow_member_posts,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
}


type CreateProjectRoleRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    *string  `json:"description,omitempty"`
	SlotsTotal     int32    `json:"slots_total" binding:"required,min=1"`
	Requirements   *string  `json:"requirements,omitempty"`
	SkillsRequired []string `json:"skills_required,omitempty"`
}


type ApplyForRoleRequest struct {
	Message *string `json:"message,omitempty"`
}


type AddGroupAdminRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Permissions []string  `json:"permissions,omitempty"`
}


type AddGroupModeratorRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Permissions []string  `json:"permissions,omitempty"`
}


type UpdateMemberRoleRequest struct {
	Role        string   `json:"role" binding:"required"`
	Permissions []string `json:"permissions,omitempty"`
}


type GroupResponse struct {
	ID               uuid.UUID              `json:"id"`
	SpaceID          uuid.UUID              `json:"space_id"`
	CommunityID      *uuid.UUID             `json:"community_id,omitempty"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	Category         string                 `json:"category"`
	GroupType        string                 `json:"group_type"`
	Avatar           *string                `json:"avatar,omitempty"`
	Banner           *string                `json:"banner,omitempty"`
	MemberCount      int32                  `json:"member_count"`
	PostCount        int32                  `json:"post_count"`
	Status           string                 `json:"status"`
	Visibility       string                 `json:"visibility"`
	AllowInvites     bool                   `json:"allow_invites"`
	AllowMemberPosts bool                   `json:"allow_member_posts"`
	CreatedBy        *uuid.UUID             `json:"created_by,omitempty"`
	Tags             []string               `json:"tags"`
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
}


type GroupDetailResponse struct {
	ID                uuid.UUID              `json:"id"`
	SpaceID           uuid.UUID              `json:"space_id"`
	CommunityID       *uuid.UUID             `json:"community_id,omitempty"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description,omitempty"`
	Category          string                 `json:"category"`
	GroupType         string                 `json:"group_type"`
	Avatar            *string                `json:"avatar,omitempty"`
	Banner            *string                `json:"banner,omitempty"`
	MemberCount       int32                  `json:"member_count"`
	PostCount         int32                  `json:"post_count"`
	Status            string                 `json:"status"`
	Visibility        string                 `json:"visibility"`
	AllowInvites      bool                   `json:"allow_invites"`
	AllowMemberPosts  bool                   `json:"allow_member_posts"`
	CreatedBy         *uuid.UUID             `json:"created_by,omitempty"`
	Tags              []string               `json:"tags"`
	Settings          *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	CommunityName     *string                `json:"community_name,omitempty"`
	CreatedByUsername string                 `json:"created_by_username"`
	CreatedByFullName string                 `json:"created_by_full_name"`
	IsMember          bool                   `json:"is_member"`
	UserRole          *string                `json:"user_role,omitempty"`
	ActualMemberCount int64                  `json:"actual_member_count"`
	ActualPostCount   int64                  `json:"actual_post_count"`
}


type GroupListResponse struct {
	ID                uuid.UUID              `json:"id"`
	SpaceID           uuid.UUID              `json:"space_id"`
	CommunityID       *uuid.UUID             `json:"community_id,omitempty"`
	Name              string                 `json:"name"`
	Description       *string                `json:"description,omitempty"`
	Category          string                 `json:"category"`
	GroupType         string                 `json:"group_type"`
	Avatar            *string                `json:"avatar,omitempty"`
	Banner            *string                `json:"banner,omitempty"`
	MemberCount       int32                  `json:"member_count"`
	PostCount         int32                  `json:"post_count"`
	Status            string                 `json:"status"`
	Visibility        string                 `json:"visibility"`
	AllowInvites      bool                   `json:"allow_invites"`
	AllowMemberPosts  bool                   `json:"allow_member_posts"`
	CreatedBy         *uuid.UUID             `json:"created_by,omitempty"`
	Tags              []string               `json:"tags"`
	Settings          *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt         *time.Time             `json:"created_at,omitempty"`
	UpdatedAt         *time.Time             `json:"updated_at,omitempty"`
	CommunityName     *string                `json:"community_name,omitempty"`
	IsMember          bool                   `json:"is_member"`
	UserRole          *string                `json:"user_role,omitempty"`
	ActualMemberCount int64                  `json:"actual_member_count"`
}


type UserGroupResponse struct {
	ID               uuid.UUID              `json:"id"`
	SpaceID          uuid.UUID              `json:"space_id"`
	CommunityID      *uuid.UUID             `json:"community_id,omitempty"`
	Name             string                 `json:"name"`
	Description      *string                `json:"description,omitempty"`
	Category         string                 `json:"category"`
	GroupType        string                 `json:"group_type"`
	Avatar           *string                `json:"avatar,omitempty"`
	Banner           *string                `json:"banner,omitempty"`
	MemberCount      int32                  `json:"member_count"`
	PostCount        int32                  `json:"post_count"`
	Status           string                 `json:"status"`
	Visibility       string                 `json:"visibility"`
	AllowInvites     bool                   `json:"allow_invites"`
	AllowMemberPosts bool                   `json:"allow_member_posts"`
	CreatedBy        *uuid.UUID             `json:"created_by,omitempty"`
	Tags             []string               `json:"tags"`
	Settings         *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	CommunityName    *string                `json:"community_name,omitempty"`
	UserRole         string                 `json:"user_role"`
	JoinedAt         *time.Time             `json:"joined_at,omitempty"`
}


type GroupMemberResponse struct {
	ID          uuid.UUID  `json:"id"`
	Username    string     `json:"username"`
	FullName    string     `json:"full_name"`
	Avatar      *string    `json:"avatar,omitempty"`
	Level       *string    `json:"level,omitempty"`
	Department  *string    `json:"department,omitempty"`
	Verified    bool       `json:"verified"`
	Role        string     `json:"role"`
	JoinedAt    *time.Time `json:"joined_at,omitempty"`
	Permissions []string   `json:"permissions"`
}


type GroupMembershipResponse struct {
	ID          uuid.UUID  `json:"id"`
	GroupID     uuid.UUID  `json:"group_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Role        string     `json:"role"`
	JoinedAt    *time.Time `json:"joined_at,omitempty"`
	InvitedBy   *uuid.UUID `json:"invited_by,omitempty"`
	Permissions []string   `json:"permissions"`
}


type ProjectRoleResponse struct {
	ID             uuid.UUID  `json:"id"`
	GroupID        uuid.UUID  `json:"group_id"`
	Name           string     `json:"name"`
	Description    *string    `json:"description,omitempty"`
	SlotsTotal     int32      `json:"slots_total"`
	SlotsFilled    int32      `json:"slots_filled"`
	Requirements   *string    `json:"requirements,omitempty"`
	SkillsRequired []string   `json:"skills_required"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}


type RoleApplicationResponse struct {
	ID          uuid.UUID  `json:"id"`
	RoleID      uuid.UUID  `json:"role_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Message     *string    `json:"message,omitempty"`
	Status      string     `json:"status"`
	AppliedAt   *time.Time `json:"applied_at,omitempty"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	ReviewedBy  *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewNotes *string    `json:"review_notes,omitempty"`
	Username    string     `json:"username"`
	FullName    string     `json:"full_name"`
	Avatar      *string    `json:"avatar,omitempty"`
	RoleName    string     `json:"role_name"`
}


type ListGroupsParams struct {
	UserID  uuid.UUID
	SpaceID uuid.UUID
	SortBy  string 
	Page    int32
	Limit   int32
}


type SearchGroupsParams struct {
	UserID  uuid.UUID
	SpaceID uuid.UUID
	Query   string
}
