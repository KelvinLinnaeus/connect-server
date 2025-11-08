package analytics

import (
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)






type CreateReportRequest struct {
	SpaceID     uuid.UUID               `json:"space_id" binding:"required"`
	ContentType string                  `json:"content_type" binding:"required,oneof=post comment event announcement message"`
	ContentID   uuid.UUID               `json:"content_id" binding:"required"`
	Reason      string                  `json:"reason" binding:"required,min=5"`
	Description *string                 `json:"description,omitempty"`
	Priority    *string                 `json:"priority,omitempty" binding:"omitempty,oneof=low medium high urgent"`
	ReporterID  uuid.UUID               
}


type ReportResponse struct {
	ID              uuid.UUID               `json:"id"`
	SpaceID         uuid.UUID               `json:"space_id"`
	ReporterID      uuid.UUID               `json:"reporter_id"`
	ContentType     string                  `json:"content_type"`
	ContentID       uuid.UUID               `json:"content_id"`
	Reason          string                  `json:"reason"`
	Description     *string                 `json:"description,omitempty"`
	Status          string                  `json:"status"`
	Priority        *string                 `json:"priority,omitempty"`
	ReviewedBy      *uuid.UUID              `json:"reviewed_by,omitempty"`
	ReviewedAt      *time.Time              `json:"reviewed_at,omitempty"`
	ModerationNotes *string                 `json:"moderation_notes,omitempty"`
	ActionsTaken    *pqtype.NullRawMessage  `json:"actions_taken,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}


type ReportDetailResponse struct {
	ID               uuid.UUID               `json:"id"`
	SpaceID          uuid.UUID               `json:"space_id"`
	ReporterID       uuid.UUID               `json:"reporter_id"`
	ReporterUsername string                  `json:"reporter_username"`
	ReporterFullName string                  `json:"reporter_full_name"`
	ContentType      string                  `json:"content_type"`
	ContentID        uuid.UUID               `json:"content_id"`
	Reason           string                  `json:"reason"`
	Description      *string                 `json:"description,omitempty"`
	Status           string                  `json:"status"`
	Priority         *string                 `json:"priority,omitempty"`
	ReviewedBy       *uuid.UUID              `json:"reviewed_by,omitempty"`
	ReviewerUsername *string                 `json:"reviewer_username,omitempty"`
	ReviewedAt       *time.Time              `json:"reviewed_at,omitempty"`
	ModerationNotes  *string                 `json:"moderation_notes,omitempty"`
	ActionsTaken     *pqtype.NullRawMessage  `json:"actions_taken,omitempty"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}


type UpdateReportRequest struct {
	Status          string                  `json:"status" binding:"required,oneof=pending approved rejected"`
	ModerationNotes *string                 `json:"moderation_notes,omitempty"`
	ActionsTaken    *pqtype.NullRawMessage  `json:"actions_taken,omitempty"`
}


type ContentModerationStatsResponse struct {
	TotalReports    int64 `json:"total_reports"`
	PendingReports  int64 `json:"pending_reports"`
	ApprovedReports int64 `json:"approved_reports"`
	RejectedReports int64 `json:"rejected_reports"`
	UrgentReports   int64 `json:"urgent_reports"`
}






type SystemMetricsResponse struct {
	TotalUsers                int64 `json:"total_users"`
	ActiveUsers               int64 `json:"active_users"`
	NewUsersToday             int64 `json:"new_users_today"`
	DailyPosts                int64 `json:"daily_posts"`
	TotalGroups               int64 `json:"total_groups"`
	TotalCommunities          int64 `json:"total_communities"`
	TotalEvents               int64 `json:"total_events"`
	PendingTutoringSessions   int64 `json:"pending_tutoring_sessions"`
	PendingMentoringSessions  int64 `json:"pending_mentoring_sessions"`
	PendingReports            int64 `json:"pending_reports"`
	PendingTutorApplications  int64 `json:"pending_tutor_applications"`
	PendingMentorApplications int64 `json:"pending_mentor_applications"`
}


type SpaceStatsResponse struct {
	Name           string     `json:"name"`
	Slug           string     `json:"slug"`
	UserCount      int64      `json:"user_count"`
	PostCount      int64      `json:"post_count"`
	CommunityCount int64      `json:"community_count"`
	GroupCount     int64      `json:"group_count"`
	FirstUserDate  *time.Time `json:"first_user_date,omitempty"`
}






type EngagementMetricResponse struct {
	Date          time.Time `json:"date"`
	PostCount     int64     `json:"post_count"`
	TotalLikes    int64     `json:"total_likes"`
	TotalComments int64     `json:"total_comments"`
	TotalViews    int64     `json:"total_views"`
}


type UserActivityStatResponse struct {
	Action string    `json:"action"`
	Count  int64     `json:"count"`
	Date   time.Time `json:"date"`
}


type UserGrowthResponse struct {
	Date     time.Time `json:"date"`
	NewUsers int64     `json:"new_users"`
}


type UserEngagementRankingResponse struct {
	ID              uuid.UUID `json:"id"`
	Username        string    `json:"username"`
	FullName        string    `json:"full_name"`
	Avatar          *string   `json:"avatar,omitempty"`
	PostCount       int64     `json:"post_count"`
	FollowersCount  int64     `json:"followers_count"`
	FollowingCount  int64     `json:"following_count"`
	EngagementScore int32     `json:"engagement_score"`
}






type TopPostResponse struct {
	ID              uuid.UUID               `json:"id"`
	AuthorID        uuid.UUID               `json:"author_id"`
	Username        string                  `json:"username"`
	FullName        string                  `json:"full_name"`
	SpaceID         uuid.UUID               `json:"space_id"`
	CommunityID     *uuid.UUID              `json:"community_id,omitempty"`
	GroupID         *uuid.UUID              `json:"group_id,omitempty"`
	Content         string                  `json:"content"`
	Media           *pqtype.NullRawMessage  `json:"media,omitempty"`
	Tags            []string                `json:"tags"`
	LikesCount      int32                   `json:"likes_count"`
	CommentsCount   int32                   `json:"comments_count"`
	ViewsCount      int32                   `json:"views_count"`
	EngagementScore int32                   `json:"engagement_score"`
	CreatedAt       time.Time               `json:"created_at"`
}


type TopCommunityResponse struct {
	ID              uuid.UUID `json:"id"`
	SpaceID         uuid.UUID `json:"space_id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description,omitempty"`
	Category        string    `json:"category"`
	CoverImage      *string   `json:"cover_image,omitempty"`
	MemberCount     int32     `json:"member_count"`
	PostCount       int32     `json:"post_count"`
	EngagementScore int32     `json:"engagement_score"`
	IsPublic        bool      `json:"is_public"`
	CreatedAt       time.Time `json:"created_at"`
}


type TopGroupResponse struct {
	ID              uuid.UUID  `json:"id"`
	SpaceID         uuid.UUID  `json:"space_id"`
	CommunityID     *uuid.UUID `json:"community_id,omitempty"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Category        string     `json:"category"`
	Avatar          *string    `json:"avatar,omitempty"`
	MemberCount     int32      `json:"member_count"`
	PostCount       int32      `json:"post_count"`
	EngagementScore int32      `json:"engagement_score"`
	Visibility      string     `json:"visibility"`
	CreatedAt       time.Time  `json:"created_at"`
}






type MentoringStatsResponse struct {
	TotalSessions     int64   `json:"total_sessions"`
	CompletedSessions int64   `json:"completed_sessions"`
	PendingSessions   int64   `json:"pending_sessions"`
	AverageRating     float64 `json:"average_rating"`
	RatedSessions     int64   `json:"rated_sessions"`
}


type TutoringStatsResponse struct {
	TotalSessions     int64   `json:"total_sessions"`
	CompletedSessions int64   `json:"completed_sessions"`
	PendingSessions   int64   `json:"pending_sessions"`
	AverageRating     float64 `json:"average_rating"`
	RatedSessions     int64   `json:"rated_sessions"`
}


type PopularIndustryResponse struct {
	Industry      string  `json:"industry"`
	SessionCount  int64   `json:"session_count"`
	AverageRating float64 `json:"average_rating"`
}


type PopularSubjectResponse struct {
	Subject       string  `json:"subject"`
	SessionCount  int64   `json:"session_count"`
	AverageRating float64 `json:"average_rating"`
}
