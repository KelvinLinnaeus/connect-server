package users

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	SpaceID    uuid.UUID `json:"space_id" binding:"required"`
	Username   string    `json:"username" binding:"required,min=3,max=30"`
	Email      string    `json:"email" binding:"required,email"`
	Password   string    `json:"password" binding:"required,min=8"`
	FullName   string    `json:"full_name" binding:"required"`
	Level      *string   `json:"level,omitempty"`
	Department *string   `json:"department,omitempty"`
	Major      *string   `json:"major,omitempty"`
	Year       *int32    `json:"year,omitempty"`
	Interests  []string  `json:"interests,omitempty"`
}

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	FullName   string   `json:"full_name"`
	Bio        *string  `json:"bio,omitempty"`
	Avatar     *string  `json:"avatar,omitempty"`
	Level      *string  `json:"level,omitempty"`
	Department *string  `json:"department,omitempty"`
	Major      *string  `json:"major,omitempty"`
	Year       *int32   `json:"year,omitempty"`
	Interests  []string `json:"interests,omitempty"`
}

// UpdatePasswordRequest represents the request to update password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UserResponse represents the user data returned to clients (excludes password)
type UserResponse struct {
	ID             uuid.UUID  `json:"id"`
	SpaceID        uuid.UUID  `json:"space_id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name"`
	Avatar         *string    `json:"avatar"`
	Bio            *string    `json:"bio"`
	Verified       *bool      `json:"verified"`
	Roles          []string   `json:"roles"`
	Level          *string    `json:"level"`
	Department     *string    `json:"department"`
	Major          *string    `json:"major"`
	Year           *int32     `json:"year"`
	Interests      []string   `json:"interests"`
	FollowersCount *int32     `json:"followers_count"`
	FollowingCount *int32     `json:"following_count"`
	MentorStatus   *string    `json:"mentor_status"`
	TutorStatus    *string    `json:"tutor_status"`
	Status         *string    `json:"status"`
	Settings       *string    `json:"settings"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	SpaceName      *string    `json:"space_name,omitempty"`
	SpaceSlug      *string    `json:"space_slug,omitempty"`
}

// SearchUsersRequest represents search parameters
type SearchUsersRequest struct {
	Query   string    `json:"query" form:"q" binding:"required"`
	SpaceID uuid.UUID `json:"space_id" form:"space_id" binding:"required"`
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page  int `form:"page" binding:"min=1"`
	Limit int `form:"limit" binding:"min=1,max=100"`
}

// GetDefaultPagination returns default pagination values
func GetDefaultPagination() PaginationParams {
	return PaginationParams{
		Page:  1,
		Limit: 20,
	}
}

// Offset calculates the offset for database queries
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// UserFollowResponse represents a user in follow/follower lists
type UserFollowResponse struct {
	ID             uuid.UUID  `json:"id"`
	Username       string     `json:"username"`
	FullName       string     `json:"full_name"`
	Avatar         *string    `json:"avatar"`
	Bio            *string    `json:"bio"`
	Verified       *bool      `json:"verified"`
	FollowersCount *int32     `json:"followers_count"`
	FollowingCount *int32     `json:"following_count"`
	FollowedAt     *time.Time `json:"followed_at"`
}
