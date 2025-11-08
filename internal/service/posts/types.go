package posts

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreatePostRequest represents request to create a new post
type CreatePostRequest struct {
	AuthorID     uuid.UUID              `json:"author_id"` // Set from auth context, not from request
	SpaceID      uuid.UUID              `json:"space_id" binding:"required"`
	CommunityID  *uuid.UUID             `json:"community_id,omitempty"`
	GroupID      *uuid.UUID             `json:"group_id,omitempty"`
	ParentPostID *uuid.UUID             `json:"parent_post_id,omitempty"`
	QuotedPostID *uuid.UUID             `json:"quoted_post_id,omitempty"`
	Content      string                 `json:"content" binding:"required,min=1,max=5000"`
	Media        *pqtype.NullRawMessage `json:"media,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	Visibility   string                 `json:"visibility,omitempty"` // public, followers, private
}

// UpdatePostRequest represents request to update a post
type UpdatePostRequest struct {
	Content    *string  `json:"content,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Visibility *string  `json:"visibility,omitempty"`
}

// CreateCommentRequest represents request to create a comment
type CreateCommentRequest struct {
	PostID          uuid.UUID  `json:"post_id"` // Set from URL param
	AuthorID        uuid.UUID  `json:"author_id"` // Set from auth context
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
	Content         string     `json:"content" binding:"required,min=1,max=1000"`
}

// CreateRepostRequest represents request to create a repost
type CreateRepostParams struct {
	AuthorID     uuid.UUID  `json:"author_id"` // Set from auth context
	SpaceID      uuid.UUID  `json:"space_id"` // Set from auth context
	QuotedPostID *uuid.UUID `json:"quoted_post_id"` // Set from URL param
	Content      string     `json:"content,omitempty"`
	Visibility   string     `json:"visibility,omitempty"`
}

// PostResponse represents a post in API responses
type PostResponse struct {
	ID               uuid.UUID              `json:"id"`
	AuthorID         uuid.UUID              `json:"author_id"`
	SpaceID          uuid.UUID              `json:"space_id"`
	CommunityID      *uuid.UUID             `json:"community_id,omitempty"`
	GroupID          *uuid.UUID             `json:"group_id,omitempty"`
	ParentPostID     *uuid.UUID             `json:"parent_post_id,omitempty"`
	QuotedPostID     *uuid.UUID             `json:"quoted_post_id,omitempty"`
	Content          string                 `json:"content"`
	Media            *pqtype.NullRawMessage `json:"media,omitempty"`
	Tags             []string               `json:"tags"`
	LikesCount       int32                  `json:"likes_count"`
	CommentsCount    int32                  `json:"comments_count"`
	RepostsCount     int32                  `json:"reposts_count"`
	QuotesCount      int32                  `json:"quotes_count"`
	ViewsCount       int32                  `json:"views_count"`
	IsPinned         bool                   `json:"is_pinned"`
	Visibility       string                 `json:"visibility"`
	Status           string                 `json:"status"`
	CreatedAt        *time.Time             `json:"created_at,omitempty"`
	UpdatedAt        *time.Time             `json:"updated_at,omitempty"`
	Username         *string                `json:"username,omitempty"`
	FullName         *string                `json:"full_name,omitempty"`
	AuthorAvatar     *string                `json:"author_avatar,omitempty"`
	CommunityName    *string                `json:"community_name,omitempty"`
	GroupName        *string                `json:"group_name,omitempty"`
	IsLiked          *bool                  `json:"is_liked,omitempty"`
	AuthorVerified   *bool                  `json:"author_verified,omitempty"`
	QuotedContent    *string                `json:"quoted_content,omitempty"`
	QuotedAuthorID   *uuid.UUID             `json:"quoted_author_id,omitempty"`
	QuotedUsername   *string                `json:"quoted_username,omitempty"`
	QuotedFullName   *string                `json:"quoted_full_name,omitempty"`
	RelevanceScore   *float32               `json:"relevance_score,omitempty"`
	EngagementScore  *int32                 `json:"engagement_score,omitempty"`
}

// CommentResponse represents a comment in API responses
type CommentResponse struct {
	ID              uuid.UUID  `json:"id"`
	PostID          uuid.UUID  `json:"post_id"`
	AuthorID        uuid.UUID  `json:"author_id"`
	ParentCommentID *uuid.UUID `json:"parent_comment_id,omitempty"`
	Content         string     `json:"content"`
	LikesCount      int32      `json:"likes_count"`
	Status          string     `json:"status"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
	Username        *string    `json:"username,omitempty"`
	FullName        *string    `json:"full_name,omitempty"`
	Avatar          *string    `json:"avatar,omitempty"`
	Depth           *int32     `json:"depth,omitempty"`
}

// UserLikeResponse represents a user who liked a post
type UserLikeResponse struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	FullName  string     `json:"full_name"`
	Avatar    *string    `json:"avatar,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

// ListPostsParams represents parameters for listing posts
type ListPostsParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	SortBy string `json:"sort_by,omitempty"` // created_at, likes_count, etc.
}

// SearchPostsParams represents parameters for searching posts
type SearchPostsParams struct {
	Query   string    `json:"query" binding:"required"`
	SpaceID uuid.UUID `json:"space_id" binding:"required"`
	Page    int       `json:"page"`
	Limit   int       `json:"limit"`
}

// TrendingTopicResponse represents a trending topic/hashtag
type TrendingTopicResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Category   string `json:"category"`
	PostsCount int64  `json:"posts_count"`
	TrendScore int32  `json:"trend_score"`
}
