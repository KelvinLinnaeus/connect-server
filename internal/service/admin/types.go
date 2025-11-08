package admin

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SuspendUserRequest struct {
	UserID       uuid.UUID
	SuspendedBy  uuid.UUID
	Reason       string
	Notes        string
	DurationDays int
	IsPermanent  bool
}

type GetReportsRequest struct {
	SpaceID     uuid.UUID
	Status      string
	ContentType string
	Limit       int32
	Offset      int32
}

type ContentReportResponse struct {
	ID               uuid.UUID  `json:"id"`
	SpaceID          uuid.UUID  `json:"space_id"`
	ContentType      string     `json:"content_type"`
	ContentID        uuid.UUID  `json:"content_id"`
	Reason           string     `json:"reason"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	ResolutionAction *string    `json:"resolution_action,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	ResolvedAt       *time.Time `json:"resolved_at,omitempty"`
}

type SpaceActivityResponse struct {
	ID           uuid.UUID  `json:"id"`
	ActivityType string     `json:"activity_type"`
	ActorID      *uuid.UUID `json:"actor_id,omitempty"`
	ActorName    string     `json:"actor_name"`
	Description  string     `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
}

type DashboardStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	NewUsersMonth    int64 `json:"new_users_month"`
	TotalPosts       int64 `json:"total_posts"`
	TotalCommunities int64 `json:"total_communities"`
	TotalGroups      int64 `json:"total_groups"`
	PendingReports   int64 `json:"pending_reports"`
	SuspensionsMonth int64 `json:"suspensions_month"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Avatar     string    `json:"avatar"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Roles      []string  `json:"roles"`
	Department string    `json:"department"`
}

type TutorApplicationResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}

type MentorApplicationResponse struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"user_id"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	ReviewedAt *time.Time `json:"reviewed_at,omitempty"`
}

type GroupResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetAllMentorApplicationsResponse struct {
	ID          uuid.UUID `json:"id"`
	ApplicantID uuid.UUID `json:"applicant_id"`
	Industry    string    `json:"industry"`
	Experience  int32     `json:"experience"`
	Specialties []string  `json:"specialties"`
	Status      string    `json:"status"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	SubmittedAt time.Time `json:"submitted_at"`
	FullName    string    `json:"full_name"`
	UserID      uuid.UUID `json:"user_id"`
}

type GetAllTutorApplicationsResponse struct {
	ApplicantID uuid.UUID `json:"applicant_id"`
	ID          uuid.UUID `json:"id"`
	Subjects    []string  `json:"subjects"`
	HourlyRate  string    `json:"hourly_rate"`
	Status      string    `json:"status"`
	SubmittedAt time.Time `json:"submitted_at"`
	FullName    string    `json:"full_name"`
	UserID      uuid.UUID `json:"user_id"`
}



type CreateCommunityRequest struct {
	SpaceID     uuid.UUID `json:"space_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	CoverImage  string    `json:"cover_image"`
	IsPublic    bool      `json:"is_public"`
	Settings    []byte    `json:"setting"`
}

type UpdateCommunityRequest struct {
	Name        string
	Description string
	Category    string
	CoverImage  string
	IsPublic    bool
	Settings    []byte
}



type CreateAnnouncementRequest struct {
	SpaceID        uuid.UUID  `json:"space_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	TargetAudience []string   `json:"target_audience"`
	Priority       string     `json:"priority"`
	ScheduledFor   *time.Time `json:"scheduled_for"`
	ExpiresAt      *time.Time `json:"expires_at"`
	Attachments    []byte     `json:"attachments"`
	IsPinned       bool       `json:"is_pinned"`
}

type UpdateAnnouncementRequest struct {
	Title          string
	Content        string
	Type           string
	TargetAudience []string
	Priority       string
	ScheduledFor   *time.Time
	ExpiresAt      *time.Time
	Attachments    []byte
	IsPinned       bool
}



type CreateEventRequest struct {
	SpaceID              uuid.UUID  `json:"space_id"`
	Title                string     `json:"title"`
	Description          string     `json:"description"`
	Category             string     `json:"category"`
	Location             string     `json:"location"`
	VenueDetails         string     `json:"venue_details"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	Timezone             string     `json:"time_zone"`
	Tags                 []string   `json:"tags"`
	ImageURL             string     `json:"image_url"`
	MaxAttendees         int32      `json:"max_attendees"`
	RegistrationRequired bool       `json:"registration_required"`
	RegistrationDeadline *time.Time `json:"registration_deadline"`
	IsPublic             bool       `json:"is_public"`
}

type UpdateEventRequest struct {
	Title                string
	Description          string
	Category             string
	Location             string
	VenueDetails         string
	StartDate            time.Time
	EndDate              time.Time
	Timezone             string
	Tags                 []string
	ImageURL             string
	MaxAttendees         int32
	RegistrationRequired bool
	RegistrationDeadline *time.Time
	IsPublic             bool
}



type CreateUserRequest struct {
	SpaceID    uuid.UUID
	Username   string
	Email      string
	Password   string
	FullName   string
	Roles      []string
	Status     string
	Department string
	Level      string
	Verified   bool
}

type UpdateUserRequest struct {
	FullName   *string
	Email      *string
	Roles      []string
	Status     *string
	Department *string
	Level      *string
	Verified   *bool
}

type GetAllCommunitiesResponse struct {
	ID                uuid.UUID       `json:"id"`
	SpaceID           uuid.UUID       `json:"space_id"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Category          string          `json:"category"`
	CoverImage        string          `json:"cover_image"`
	MemberCount       int32           `json:"member_count"`
	Status            string          `json:"status"`
	PostCount         int32           `json:"post_count"`
	IsPublic          bool            `json:"is_public"`
	CreatedBy         uuid.UUID       `json:"created_by"`
	Settings          json.RawMessage `json:"settings"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	CreatedByUsername string          `json:"created_by_username"`
	CreatedByFullName string          `json:"created_by_full_name"`
	ActualMemberCount int64           `json:"actual_member_count"`
	ActualPostCount   int64           `json:"actual_post_count"`
}
