package events

import (
	"time"

	"github.com/google/uuid"
)

// Event Request Types

type CreateEventRequest struct {
	SpaceID              uuid.UUID  `json:"space_id" binding:"required"`
	Title                string     `json:"title" binding:"required,min=3,max=200"`
	Description          *string    `json:"description,omitempty"`
	Category             string     `json:"category" binding:"required"`
	Location             *string    `json:"location,omitempty"`
	VenueDetails         *string    `json:"venue_details,omitempty"`
	StartDate            time.Time  `json:"start_date" binding:"required"`
	EndDate              time.Time  `json:"end_date" binding:"required"`
	Timezone             *string    `json:"timezone,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	ImageURL             *string    `json:"image_url,omitempty"`
	MaxAttendees         *int32     `json:"max_attendees,omitempty"`
	RegistrationRequired *bool      `json:"registration_required,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
	IsPublic             *bool      `json:"is_public,omitempty"`
	OrganizerID          uuid.UUID  // Set from auth context, not from request body
}

type UpdateEventRequest struct {
	Title                string     `json:"title" binding:"required,min=3,max=200"`
	Description          *string    `json:"description,omitempty"`
	Category             string     `json:"category" binding:"required"`
	Location             *string    `json:"location,omitempty"`
	VenueDetails         *string    `json:"venue_details,omitempty"`
	StartDate            time.Time  `json:"start_date" binding:"required"`
	EndDate              time.Time  `json:"end_date" binding:"required"`
	Timezone             *string    `json:"timezone,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	ImageURL             *string    `json:"image_url,omitempty"`
	MaxAttendees         *int32     `json:"max_attendees,omitempty"`
	RegistrationRequired *bool      `json:"registration_required,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
	IsPublic             *bool      `json:"is_public,omitempty"`
}

type UpdateEventStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=draft published cancelled completed"`
}

type ListEventsParams struct {
	SpaceID   uuid.UUID
	UserID    uuid.UUID
	Category  *string
	StartDate *time.Time
	Sort      *string // "upcoming", "popular", or empty for recent
	Page      int32
	Limit     int32
}

type SearchEventsParams struct {
	SpaceID uuid.UUID
	UserID  uuid.UUID
	Query   string
}

type AddCoOrganizerRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

// Event Response Types

type EventResponse struct {
	ID                   uuid.UUID  `json:"id"`
	SpaceID              uuid.UUID  `json:"space_id"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	Category             string     `json:"category"`
	Location             *string    `json:"location,omitempty"`
	VenueDetails         *string    `json:"venue_details,omitempty"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	Timezone             *string    `json:"timezone,omitempty"`
	OrganizerID          *uuid.UUID `json:"organizer_id,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	ImageURL             *string    `json:"image_url,omitempty"`
	MaxAttendees         *int32     `json:"max_attendees,omitempty"`
	CurrentAttendees     *int32     `json:"current_attendees,omitempty"`
	RegistrationRequired *bool      `json:"registration_required,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
	Status               *string    `json:"status,omitempty"`
	IsPublic             *bool      `json:"is_public,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

type EventDetailResponse struct {
	ID                   uuid.UUID  `json:"id"`
	SpaceID              uuid.UUID  `json:"space_id"`
	Title                string     `json:"title"`
	Description          *string    `json:"description,omitempty"`
	Category             string     `json:"category"`
	Location             *string    `json:"location,omitempty"`
	VenueDetails         *string    `json:"venue_details,omitempty"`
	StartDate            time.Time  `json:"start_date"`
	EndDate              time.Time  `json:"end_date"`
	Timezone             *string    `json:"timezone,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	ImageURL             *string    `json:"image_url,omitempty"`
	MaxAttendees         *int32     `json:"max_attendees,omitempty"`
	RegistrationRequired *bool      `json:"registration_required,omitempty"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty"`
	Status               *string    `json:"status,omitempty"`
	IsPublic             *bool      `json:"is_public,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`

	// Organizer info
	OrganizerID       *uuid.UUID `json:"organizer_id,omitempty"`
	OrganizerUsername string     `json:"organizer_username"`
	OrganizerFullName string     `json:"organizer_full_name"`
	OrganizerAvatar   *string    `json:"organizer_avatar,omitempty"`

	// Attendee info
	CurrentAttendeesCount int64   `json:"current_attendees_count"`
	IsRegistered          bool    `json:"is_registered"`
	UserAttendanceStatus  *string `json:"user_attendance_status,omitempty"`
}

type EventListResponse struct {
	ID                    uuid.UUID  `json:"id"`
	SpaceID               uuid.UUID  `json:"space_id"`
	Title                 string     `json:"title"`
	Description           *string    `json:"description,omitempty"`
	Category              string     `json:"category"`
	Location              *string    `json:"location,omitempty"`
	StartDate             time.Time  `json:"start_date"`
	EndDate               time.Time  `json:"end_date"`
	Tags                  []string   `json:"tags,omitempty"`
	ImageURL              *string    `json:"image_url,omitempty"`
	MaxAttendees          *int32     `json:"max_attendees,omitempty"`
	RegistrationRequired  *bool      `json:"registration_required,omitempty"`
	IsPublic              *bool      `json:"is_public,omitempty"`

	// Organizer
	OrganizerUsername string `json:"organizer_username"`
	OrganizerFullName string `json:"organizer_full_name"`

	// Attendance
	CurrentAttendeesCount int64 `json:"current_attendees_count"`
	IsRegistered          bool  `json:"is_registered"`
}

type UserEventResponse struct {
	ID                uuid.UUID  `json:"id"`
	SpaceID           uuid.UUID  `json:"space_id"`
	Title             string     `json:"title"`
	Description       *string    `json:"description,omitempty"`
	Category          string     `json:"category"`
	Location          *string    `json:"location,omitempty"`
	StartDate         time.Time  `json:"start_date"`
	EndDate           time.Time  `json:"end_date"`
	Tags              []string   `json:"tags,omitempty"`
	ImageURL          *string    `json:"image_url,omitempty"`
	OrganizerUsername string     `json:"organizer_username"`
	OrganizerFullName string     `json:"organizer_full_name"`
	AttendanceStatus  *string    `json:"attendance_status,omitempty"` // registered, attended, cancelled
}

type AttendeeResponse struct {
	ID           uuid.UUID  `json:"id"`
	EventID      uuid.UUID  `json:"event_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Username     string     `json:"username"`
	FullName     string     `json:"full_name"`
	Avatar       *string    `json:"avatar,omitempty"`
	Department   *string    `json:"department,omitempty"`
	Level        *string    `json:"level,omitempty"`
	Status       *string    `json:"status,omitempty"` // registered, attended, cancelled
	Role         *string    `json:"role,omitempty"`   // organizer, attendee
	RegisteredAt *time.Time `json:"registered_at,omitempty"`
	AttendedAt   *time.Time `json:"attended_at,omitempty"`
	Notes        *string    `json:"notes,omitempty"`
}

type CoOrganizerResponse struct {
	ID           uuid.UUID  `json:"id"`
	EventID      uuid.UUID  `json:"event_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Username     string     `json:"username"`
	FullName     string     `json:"full_name"`
	Avatar       *string    `json:"avatar,omitempty"`
	RegisteredAt *time.Time `json:"registered_at,omitempty"`
}

type RegistrationResponse struct {
	EventID      uuid.UUID  `json:"event_id"`
	UserID       uuid.UUID  `json:"user_id"`
	Status       *string    `json:"status,omitempty"`
	RegisteredAt *time.Time `json:"registered_at,omitempty"`
	Message      string     `json:"message"`
}
