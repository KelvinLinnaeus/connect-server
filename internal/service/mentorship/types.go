package mentorship

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)



type CreateMentorProfileRequest struct {
	SpaceID      uuid.UUID       `json:"space_id" binding:"required"`
	Industry     string          `json:"industry" binding:"required"`
	Company      *string         `json:"company,omitempty"`
	Position     *string         `json:"position,omitempty"`
	Experience   int32           `json:"experience" binding:"required,min=0"`
	Specialties  []string        `json:"specialties" binding:"required,min=1"`
	Description  *string         `json:"description,omitempty"`
	Availability json.RawMessage `json:"availability,omitempty"`
	UserID       uuid.UUID       
}

type UpdateMentorAvailabilityRequest struct {
	IsAvailable  bool            `json:"is_available"`
	Availability json.RawMessage `json:"availability,omitempty"`
}

type SearchMentorsParams struct {
	SpaceID     uuid.UUID
	Industry    *string
	Specialties []string
	MinRating   *float64
	Page        int32
	Limit       int32
}

type MentorProfileResponse struct {
	ID            uuid.UUID              `json:"id"`
	UserID        uuid.UUID              `json:"user_id"`
	SpaceID       uuid.UUID              `json:"space_id"`
	Industry      string                 `json:"industry"`
	Company       *string                `json:"company,omitempty"`
	Position      *string                `json:"position,omitempty"`
	Experience    int32                  `json:"experience"`
	Specialties   []string               `json:"specialties"`
	Rating        *float64               `json:"rating,omitempty"`
	ReviewCount   *int32                 `json:"review_count,omitempty"`
	TotalSessions *int32                 `json:"total_sessions,omitempty"`
	Availability  *pqtype.NullRawMessage `json:"availability,omitempty"`
	Description   *string                `json:"description,omitempty"`
	Verified      *bool                  `json:"verified,omitempty"`
	IsAvailable   *bool                  `json:"is_available,omitempty"`
	CreatedAt     *time.Time             `json:"created_at,omitempty"`
	UpdatedAt     *time.Time             `json:"updated_at,omitempty"`
}

type MentorSearchResponse struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	FullName    string    `json:"full_name"`
	Avatar      *string   `json:"avatar,omitempty"`
	Industry    string    `json:"industry"`
	Company     *string   `json:"company,omitempty"`
	Position    *string   `json:"position,omitempty"`
	Experience  int32     `json:"experience"`
	Specialties []string  `json:"specialties"`
	Rating      *float64  `json:"rating,omitempty"`
	ReviewCount *int32    `json:"review_count,omitempty"`
	IsAvailable *bool     `json:"is_available,omitempty"`
}



type CreateTutorProfileRequest struct {
	SpaceID        uuid.UUID       `json:"space_id" binding:"required"`
	Subjects       []string        `json:"subjects" binding:"required,min=1"`
	HourlyRate     *string         `json:"hourly_rate,omitempty"`
	Description    *string         `json:"description,omitempty"`
	Experience     *string         `json:"experience,omitempty"`
	Qualifications *string         `json:"qualifications,omitempty"`
	Availability   json.RawMessage `json:"availability,omitempty"`
	UserID         uuid.UUID       
}

type UpdateTutorAvailabilityRequest struct {
	IsAvailable  bool            `json:"is_available"`
	Availability json.RawMessage `json:"availability,omitempty"`
}

type SearchTutorsParams struct {
	SpaceID   uuid.UUID
	Subjects  []string
	MinRating *float64
	MaxRate   *string
	Page      int32
	Limit     int32
}

type TutorProfileResponse struct {
	ID             uuid.UUID              `json:"id"`
	UserID         uuid.UUID              `json:"user_id"`
	SpaceID        uuid.UUID              `json:"space_id"`
	Subjects       []string               `json:"subjects"`
	HourlyRate     *string                `json:"hourly_rate,omitempty"`
	Rating         *float64               `json:"rating,omitempty"`
	ReviewCount    *int32                 `json:"review_count,omitempty"`
	TotalSessions  *int32                 `json:"total_sessions,omitempty"`
	Description    *string                `json:"description,omitempty"`
	Availability   *pqtype.NullRawMessage `json:"availability,omitempty"`
	Experience     *string                `json:"experience,omitempty"`
	Qualifications *string                `json:"qualifications,omitempty"`
	Verified       *bool                  `json:"verified,omitempty"`
	IsAvailable    *bool                  `json:"is_available,omitempty"`
	CreatedAt      *time.Time             `json:"created_at,omitempty"`
	UpdatedAt      *time.Time             `json:"updated_at,omitempty"`
}

type TutorSearchResponse struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Username       string    `json:"username"`
	FullName       string    `json:"full_name"`
	Avatar         *string   `json:"avatar,omitempty"`
	Subjects       []string  `json:"subjects"`
	HourlyRate     *string   `json:"hourly_rate,omitempty"`
	Rating         *float64  `json:"rating,omitempty"`
	ReviewCount    *int32    `json:"review_count,omitempty"`
	Experience     *string   `json:"experience,omitempty"`
	Qualifications *string   `json:"qualifications,omitempty"`
	IsAvailable    *bool     `json:"is_available,omitempty"`
}



type CreateMentoringSessionRequest struct {
	MentorID    uuid.UUID `json:"mentor_id" binding:"required"`
	SpaceID     uuid.UUID `json:"space_id" binding:"required"`
	Topic       string    `json:"topic" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	Duration    int32     `json:"duration" binding:"required,min=15"`
	MenteeNotes *string   `json:"mentee_notes,omitempty"`
	MenteeID    uuid.UUID 
}

type UpdateMentoringSessionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=scheduled confirmed in-progress completed cancelled"`
}

type AddMeetingLinkRequest struct {
	MeetingLink string `json:"meeting_link" binding:"required,url"`
}

type RateMentoringSessionRequest struct {
	Rating int32   `json:"rating" binding:"required,min=1,max=5"`
	Review *string `json:"review,omitempty"`
}

type MentoringSessionResponse struct {
	ID          uuid.UUID  `json:"id"`
	MentorID    uuid.UUID  `json:"mentor_id"`
	MenteeID    uuid.UUID  `json:"mentee_id"`
	SpaceID     uuid.UUID  `json:"space_id"`
	Topic       string     `json:"topic"`
	Status      *string    `json:"status,omitempty"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	Duration    int32      `json:"duration"`
	MenteeNotes *string    `json:"mentee_notes,omitempty"`
	MentorNotes *string    `json:"mentor_notes,omitempty"`
	MeetingLink *string    `json:"meeting_link,omitempty"`
	Rating      *int32     `json:"rating,omitempty"`
	Review      *string    `json:"review,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

type MentoringSessionDetailResponse struct {
	ID           uuid.UUID  `json:"id"`
	MentorID     uuid.UUID  `json:"mentor_id"`
	MentorName   string     `json:"mentor_name"`
	MentorAvatar *string    `json:"mentor_avatar,omitempty"`
	MenteeID     uuid.UUID  `json:"mentee_id"`
	MenteeName   string     `json:"mentee_name"`
	MenteeAvatar *string    `json:"mentee_avatar,omitempty"`
	SpaceID      uuid.UUID  `json:"space_id"`
	Topic        string     `json:"topic"`
	Status       *string    `json:"status,omitempty"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Duration     int32      `json:"duration"`
	MenteeNotes  *string    `json:"mentee_notes,omitempty"`
	MentorNotes  *string    `json:"mentor_notes,omitempty"`
	MeetingLink  *string    `json:"meeting_link,omitempty"`
	Rating       *int32     `json:"rating,omitempty"`
	Review       *string    `json:"review,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}



type CreateTutoringSessionRequest struct {
	TutorID      uuid.UUID `json:"tutor_id" binding:"required"`
	SpaceID      uuid.UUID `json:"space_id" binding:"required"`
	Subject      string    `json:"subject" binding:"required"`
	ScheduledAt  time.Time `json:"scheduled_at" binding:"required"`
	Duration     int32     `json:"duration" binding:"required,min=15"`
	HourlyRate   *string   `json:"hourly_rate,omitempty"`
	StudentNotes *string   `json:"student_notes,omitempty"`
	StudentID    uuid.UUID 
}

type UpdateTutoringSessionStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=scheduled confirmed in-progress completed cancelled"`
}

type RateTutoringSessionRequest struct {
	Rating int32   `json:"rating" binding:"required,min=1,max=5"`
	Review *string `json:"review,omitempty"`
}

type TutoringSessionResponse struct {
	ID           uuid.UUID  `json:"id"`
	TutorID      uuid.UUID  `json:"tutor_id"`
	StudentID    uuid.UUID  `json:"student_id"`
	SpaceID      uuid.UUID  `json:"space_id"`
	Subject      string     `json:"subject"`
	Status       *string    `json:"status,omitempty"`
	ScheduledAt  time.Time  `json:"scheduled_at"`
	Duration     int32      `json:"duration"`
	HourlyRate   *string    `json:"hourly_rate,omitempty"`
	TotalAmount  *string    `json:"total_amount,omitempty"`
	StudentNotes *string    `json:"student_notes,omitempty"`
	TutorNotes   *string    `json:"tutor_notes,omitempty"`
	MeetingLink  *string    `json:"meeting_link,omitempty"`
	Rating       *int32     `json:"rating,omitempty"`
	Review       *string    `json:"review,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

type TutoringSessionDetailResponse struct {
	ID            uuid.UUID  `json:"id"`
	TutorID       uuid.UUID  `json:"tutor_id"`
	TutorName     string     `json:"tutor_name"`
	TutorAvatar   *string    `json:"tutor_avatar,omitempty"`
	StudentID     uuid.UUID  `json:"student_id"`
	StudentName   string     `json:"student_name"`
	StudentAvatar *string    `json:"student_avatar,omitempty"`
	SpaceID       uuid.UUID  `json:"space_id"`
	Subject       string     `json:"subject"`
	Status        *string    `json:"status,omitempty"`
	ScheduledAt   time.Time  `json:"scheduled_at"`
	Duration      int32      `json:"duration"`
	HourlyRate    *string    `json:"hourly_rate,omitempty"`
	TotalAmount   *string    `json:"total_amount,omitempty"`
	StudentNotes  *string    `json:"student_notes,omitempty"`
	TutorNotes    *string    `json:"tutor_notes,omitempty"`
	MeetingLink   *string    `json:"meeting_link,omitempty"`
	Rating        *int32     `json:"rating,omitempty"`
	Review        *string    `json:"review,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}



type CreateMentorApplicationRequest struct {
	SpaceID     uuid.UUID `json:"space_id" binding:"required"`
	Industry    string    `json:"industry" binding:"required"`
	Company     *string   `json:"company,omitempty"`
	Position    *string   `json:"position,omitempty"`
	Experience  int32     `json:"experience" binding:"required,min=0"`
	Specialties []string  `json:"specialties" binding:"required,min=1"`
	Motivation  string    `json:"motivation" binding:"required,min=50"`
	UserID      uuid.UUID 
}

type UpdateMentorApplicationRequest struct {
	Status         string  `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComments *string `json:"review_comments,omitempty"`
}

type MentorApplicationResponse struct {
	ID             uuid.UUID  `json:"id,omitempty"`
	UserID         uuid.UUID  `json:"user_id,omitempty"`
	SpaceID        uuid.UUID  `json:"space_id,omitempty"`
	Industry       string     `json:"industry,omitempty"`
	Company        *string    `json:"company,omitempty"`
	Position       *string    `json:"position,omitempty"`
	Experience     int32      `json:"experience,omitempty"`
	Specialties    []string   `json:"specialties,omitempty"`
	Motivation     string     `json:"motivation,omitempty"`
	Status         *string    `json:"status,omitempty"`
	ReviewedBy     *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewComments *string    `json:"review_comments,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
}

type CreateTutorApplicationRequest struct {
	SpaceID        uuid.UUID `json:"space_id" binding:"required"`
	Subjects       []string  `json:"subjects" binding:"required,min=1"`
	Experience     string    `json:"experience" binding:"required"`
	Qualifications string    `json:"qualifications" binding:"required"`
	Motivation     string    `json:"motivation" binding:"required,min=50"`
	UserID         uuid.UUID 
}

type UpdateTutorApplicationRequest struct {
	Status         string  `json:"status" binding:"required,oneof=approved rejected"`
	ReviewComments *string `json:"review_comments,omitempty"`
}

type TutorApplicationResponse struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	SpaceID        uuid.UUID  `json:"space_id"`
	Subjects       []string   `json:"subjects"`
	Experience     string     `json:"experience"`
	Qualifications string     `json:"qualifications"`
	Motivation     string     `json:"motivation"`
	Status         *string    `json:"status,omitempty"`
	ReviewedBy     *uuid.UUID `json:"reviewed_by,omitempty"`
	ReviewComments *string    `json:"review_comments,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
}



type MentorReviewResponse struct {
	ID           uuid.UUID  `json:"id"`
	MentorID     uuid.UUID  `json:"mentor_id"`
	MenteeID     uuid.UUID  `json:"mentee_id"`
	MenteeName   string     `json:"mentee_name"`
	MenteeAvatar *string    `json:"mentee_avatar,omitempty"`
	SessionID    uuid.UUID  `json:"session_id"`
	Topic        string     `json:"topic"`
	Rating       int32      `json:"rating"`
	Review       *string    `json:"review,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

type TutorReviewResponse struct {
	ID            uuid.UUID  `json:"id"`
	TutorID       uuid.UUID  `json:"tutor_id"`
	StudentID     uuid.UUID  `json:"student_id"`
	StudentName   string     `json:"student_name"`
	StudentAvatar *string    `json:"student_avatar,omitempty"`
	SessionID     uuid.UUID  `json:"session_id"`
	Subject       string     `json:"subject"`
	Rating        int32      `json:"rating"`
	Review        *string    `json:"review,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
}
