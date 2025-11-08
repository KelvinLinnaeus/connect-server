




package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const addEventCoOrganizer = `-- name: AddEventCoOrganizer :one
INSERT INTO event_attendees (event_id, user_id, role)
VALUES ($1, $2, 'organizer')
ON CONFLICT (event_id, user_id) DO NOTHING
RETURNING id, event_id, user_id, status, role, registered_at, attended_at, notes
`

type AddEventCoOrganizerParams struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) AddEventCoOrganizer(ctx context.Context, arg AddEventCoOrganizerParams) (EventAttendee, error) {
	row := q.db.QueryRowContext(ctx, addEventCoOrganizer, arg.EventID, arg.UserID)
	var i EventAttendee
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.UserID,
		&i.Status,
		&i.Role,
		&i.RegisteredAt,
		&i.AttendedAt,
		&i.Notes,
	)
	return i, err
}

const createAnnouncement = `-- name: CreateAnnouncement :one
INSERT INTO announcements (
    space_id, title, content, type, target_audience, priority,
    author_id, scheduled_for, expires_at, attachments, is_pinned
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, space_id, title, content, type, target_audience, priority, status, author_id, scheduled_for, expires_at, attachments, is_pinned, created_at, updated_at
`

type CreateAnnouncementParams struct {
	SpaceID        uuid.UUID             `json:"space_id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	TargetAudience []string              `json:"target_audience"`
	Priority       sql.NullString        `json:"priority"`
	AuthorID       uuid.NullUUID         `json:"author_id"`
	ScheduledFor   sql.NullTime          `json:"scheduled_for"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	IsPinned       sql.NullBool          `json:"is_pinned"`
}

func (q *Queries) CreateAnnouncement(ctx context.Context, arg CreateAnnouncementParams) (Announcement, error) {
	row := q.db.QueryRowContext(ctx, createAnnouncement,
		arg.SpaceID,
		arg.Title,
		arg.Content,
		arg.Type,
		pq.Array(arg.TargetAudience),
		arg.Priority,
		arg.AuthorID,
		arg.ScheduledFor,
		arg.ExpiresAt,
		arg.Attachments,
		arg.IsPinned,
	)
	var i Announcement
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Content,
		&i.Type,
		pq.Array(&i.TargetAudience),
		&i.Priority,
		&i.Status,
		&i.AuthorID,
		&i.ScheduledFor,
		&i.ExpiresAt,
		&i.Attachments,
		&i.IsPinned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createEvent = `-- name: CreateEvent :one

INSERT INTO events (
    space_id, title, description, category, location, venue_details,
    start_date, end_date, timezone, organizer, tags, image_url,
    max_attendees, registration_required, registration_deadline, is_public
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, space_id, title, description, category, location, venue_details, start_date, end_date, timezone, organizer, tags, image_url, max_attendees, current_attendees, registration_required, registration_deadline, status, is_public, created_at, updated_at
`

type CreateEventParams struct {
	SpaceID              uuid.UUID      `json:"space_id"`
	Title                string         `json:"title"`
	Description          sql.NullString `json:"description"`
	Category             string         `json:"category"`
	Location             sql.NullString `json:"location"`
	VenueDetails         sql.NullString `json:"venue_details"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	Timezone             sql.NullString `json:"timezone"`
	Organizer            uuid.NullUUID  `json:"organizer"`
	Tags                 []string       `json:"tags"`
	ImageUrl             sql.NullString `json:"image_url"`
	MaxAttendees         sql.NullInt32  `json:"max_attendees"`
	RegistrationRequired sql.NullBool   `json:"registration_required"`
	RegistrationDeadline sql.NullTime   `json:"registration_deadline"`
	IsPublic             sql.NullBool   `json:"is_public"`
}


func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.SpaceID,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.Location,
		arg.VenueDetails,
		arg.StartDate,
		arg.EndDate,
		arg.Timezone,
		arg.Organizer,
		pq.Array(arg.Tags),
		arg.ImageUrl,
		arg.MaxAttendees,
		arg.RegistrationRequired,
		arg.RegistrationDeadline,
		arg.IsPublic,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Location,
		&i.VenueDetails,
		&i.StartDate,
		&i.EndDate,
		&i.Timezone,
		&i.Organizer,
		pq.Array(&i.Tags),
		&i.ImageUrl,
		&i.MaxAttendees,
		&i.CurrentAttendees,
		&i.RegistrationRequired,
		&i.RegistrationDeadline,
		&i.Status,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAnnouncement = `-- name: DeleteAnnouncement :exec
DELETE FROM announcements WHERE id = $1
`

func (q *Queries) DeleteAnnouncement(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAnnouncement, id)
	return err
}

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1
`

func (q *Queries) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteEvent, id)
	return err
}

const getAnnouncementByID = `-- name: GetAnnouncementByID :one
SELECT 
    a.id, a.space_id, a.title, a.content, a.type, a.target_audience, a.priority, a.status, a.author_id, a.scheduled_for, a.expires_at, a.attachments, a.is_pinned, a.created_at, a.updated_at,
    u.username as author_username,
    u.full_name as author_full_name,
    u.avatar as author_avatar
FROM announcements a
JOIN users u ON a.author_id = u.id
WHERE a.id = $1
`

type GetAnnouncementByIDRow struct {
	ID             uuid.UUID             `json:"id"`
	SpaceID        uuid.UUID             `json:"space_id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	TargetAudience []string              `json:"target_audience"`
	Priority       sql.NullString        `json:"priority"`
	Status         sql.NullString        `json:"status"`
	AuthorID       uuid.NullUUID         `json:"author_id"`
	ScheduledFor   sql.NullTime          `json:"scheduled_for"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	IsPinned       sql.NullBool          `json:"is_pinned"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	UpdatedAt      sql.NullTime          `json:"updated_at"`
	AuthorUsername string                `json:"author_username"`
	AuthorFullName string                `json:"author_full_name"`
	AuthorAvatar   sql.NullString        `json:"author_avatar"`
}

func (q *Queries) GetAnnouncementByID(ctx context.Context, id uuid.UUID) (GetAnnouncementByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getAnnouncementByID, id)
	var i GetAnnouncementByIDRow
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Content,
		&i.Type,
		pq.Array(&i.TargetAudience),
		&i.Priority,
		&i.Status,
		&i.AuthorID,
		&i.ScheduledFor,
		&i.ExpiresAt,
		&i.Attachments,
		&i.IsPinned,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AuthorUsername,
		&i.AuthorFullName,
		&i.AuthorAvatar,
	)
	return i, err
}

const getEventAttendees = `-- name: GetEventAttendees :many
SELECT 
    ea.id, ea.event_id, ea.user_id, ea.status, ea.role, ea.registered_at, ea.attended_at, ea.notes,
    u.username,
    u.full_name,
    u.avatar,
    u.department,
    u.level
FROM event_attendees ea
JOIN users u ON ea.user_id = u.id
WHERE ea.event_id = $1 AND u.status = 'active'
ORDER BY ea.registered_at DESC
`

type GetEventAttendeesRow struct {
	ID           uuid.UUID      `json:"id"`
	EventID      uuid.UUID      `json:"event_id"`
	UserID       uuid.UUID      `json:"user_id"`
	Status       sql.NullString `json:"status"`
	Role         sql.NullString `json:"role"`
	RegisteredAt sql.NullTime   `json:"registered_at"`
	AttendedAt   sql.NullTime   `json:"attended_at"`
	Notes        sql.NullString `json:"notes"`
	Username     string         `json:"username"`
	FullName     string         `json:"full_name"`
	Avatar       sql.NullString `json:"avatar"`
	Department   sql.NullString `json:"department"`
	Level        sql.NullString `json:"level"`
}

func (q *Queries) GetEventAttendees(ctx context.Context, eventID uuid.UUID) ([]GetEventAttendeesRow, error) {
	rows, err := q.db.QueryContext(ctx, getEventAttendees, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEventAttendeesRow{}
	for rows.Next() {
		var i GetEventAttendeesRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.UserID,
			&i.Status,
			&i.Role,
			&i.RegisteredAt,
			&i.AttendedAt,
			&i.Notes,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.Department,
			&i.Level,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventByID = `-- name: GetEventByID :one
SELECT 
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea2 WHERE ea2.event_id = e.id AND ea2.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea3 WHERE ea3.event_id = e.id AND ea3.user_id = $1) as is_registered,
    ea.status as user_attendance_status
FROM events e
JOIN users u ON e.organizer = u.id
LEFT JOIN event_attendees ea ON e.id = ea.event_id AND ea.user_id = $1
WHERE e.id = $2 AND e.status = 'published'
`

type GetEventByIDParams struct {
	UserID uuid.UUID `json:"user_id"`
	ID     uuid.UUID `json:"id"`
}

type GetEventByIDRow struct {
	ID                    uuid.UUID      `json:"id"`
	SpaceID               uuid.UUID      `json:"space_id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	Category              string         `json:"category"`
	Location              sql.NullString `json:"location"`
	VenueDetails          sql.NullString `json:"venue_details"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
	Timezone              sql.NullString `json:"timezone"`
	Organizer             uuid.NullUUID  `json:"organizer"`
	Tags                  []string       `json:"tags"`
	ImageUrl              sql.NullString `json:"image_url"`
	MaxAttendees          sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees      sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired  sql.NullBool   `json:"registration_required"`
	RegistrationDeadline  sql.NullTime   `json:"registration_deadline"`
	Status                sql.NullString `json:"status"`
	IsPublic              sql.NullBool   `json:"is_public"`
	CreatedAt             sql.NullTime   `json:"created_at"`
	UpdatedAt             sql.NullTime   `json:"updated_at"`
	OrganizerUsername     string         `json:"organizer_username"`
	OrganizerFullName     string         `json:"organizer_full_name"`
	OrganizerAvatar       sql.NullString `json:"organizer_avatar"`
	CurrentAttendeesCount int64          `json:"current_attendees_count"`
	IsRegistered          bool           `json:"is_registered"`
	UserAttendanceStatus  sql.NullString `json:"user_attendance_status"`
}

func (q *Queries) GetEventByID(ctx context.Context, arg GetEventByIDParams) (GetEventByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getEventByID, arg.UserID, arg.ID)
	var i GetEventByIDRow
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Location,
		&i.VenueDetails,
		&i.StartDate,
		&i.EndDate,
		&i.Timezone,
		&i.Organizer,
		pq.Array(&i.Tags),
		&i.ImageUrl,
		&i.MaxAttendees,
		&i.CurrentAttendees,
		&i.RegistrationRequired,
		&i.RegistrationDeadline,
		&i.Status,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrganizerUsername,
		&i.OrganizerFullName,
		&i.OrganizerAvatar,
		&i.CurrentAttendeesCount,
		&i.IsRegistered,
		&i.UserAttendanceStatus,
	)
	return i, err
}

const getEventCategories = `-- name: GetEventCategories :many
SELECT DISTINCT category
FROM events
WHERE space_id = $1 AND status = 'published'
ORDER BY category
`

func (q *Queries) GetEventCategories(ctx context.Context, spaceID uuid.UUID) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getEventCategories, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventCoOrganizers = `-- name: GetEventCoOrganizers :many
SELECT 
    eco.id, eco.event_id, eco.user_id, eco.status, eco.role, eco.registered_at, eco.attended_at, eco.notes,
    u.username,
    u.full_name,
    u.avatar
FROM event_attendees eco
JOIN users u ON eco.user_id = u.id
WHERE eco.event_id = $1 AND eco.role = 'organizer' AND u.status = 'active'
`

type GetEventCoOrganizersRow struct {
	ID           uuid.UUID      `json:"id"`
	EventID      uuid.UUID      `json:"event_id"`
	UserID       uuid.UUID      `json:"user_id"`
	Status       sql.NullString `json:"status"`
	Role         sql.NullString `json:"role"`
	RegisteredAt sql.NullTime   `json:"registered_at"`
	AttendedAt   sql.NullTime   `json:"attended_at"`
	Notes        sql.NullString `json:"notes"`
	Username     string         `json:"username"`
	FullName     string         `json:"full_name"`
	Avatar       sql.NullString `json:"avatar"`
}

func (q *Queries) GetEventCoOrganizers(ctx context.Context, eventID uuid.UUID) ([]GetEventCoOrganizersRow, error) {
	rows, err := q.db.QueryContext(ctx, getEventCoOrganizers, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetEventCoOrganizersRow{}
	for rows.Next() {
		var i GetEventCoOrganizersRow
		if err := rows.Scan(
			&i.ID,
			&i.EventID,
			&i.UserID,
			&i.Status,
			&i.Role,
			&i.RegisteredAt,
			&i.AttendedAt,
			&i.Notes,
			&i.Username,
			&i.FullName,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventWithRegistrations = `-- name: GetEventWithRegistrations :one
SELECT
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') as registered_count,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'attended') as attended_count
FROM events e
LEFT JOIN users u ON e.organizer = u.id
WHERE e.id = $1
`

type GetEventWithRegistrationsRow struct {
	ID                   uuid.UUID      `json:"id"`
	SpaceID              uuid.UUID      `json:"space_id"`
	Title                string         `json:"title"`
	Description          sql.NullString `json:"description"`
	Category             string         `json:"category"`
	Location             sql.NullString `json:"location"`
	VenueDetails         sql.NullString `json:"venue_details"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	Timezone             sql.NullString `json:"timezone"`
	Organizer            uuid.NullUUID  `json:"organizer"`
	Tags                 []string       `json:"tags"`
	ImageUrl             sql.NullString `json:"image_url"`
	MaxAttendees         sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees     sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired sql.NullBool   `json:"registration_required"`
	RegistrationDeadline sql.NullTime   `json:"registration_deadline"`
	Status               sql.NullString `json:"status"`
	IsPublic             sql.NullBool   `json:"is_public"`
	CreatedAt            sql.NullTime   `json:"created_at"`
	UpdatedAt            sql.NullTime   `json:"updated_at"`
	OrganizerUsername    sql.NullString `json:"organizer_username"`
	OrganizerFullName    sql.NullString `json:"organizer_full_name"`
	OrganizerAvatar      sql.NullString `json:"organizer_avatar"`
	RegisteredCount      int64          `json:"registered_count"`
	AttendedCount        int64          `json:"attended_count"`
}

func (q *Queries) GetEventWithRegistrations(ctx context.Context, id uuid.UUID) (GetEventWithRegistrationsRow, error) {
	row := q.db.QueryRowContext(ctx, getEventWithRegistrations, id)
	var i GetEventWithRegistrationsRow
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Location,
		&i.VenueDetails,
		&i.StartDate,
		&i.EndDate,
		&i.Timezone,
		&i.Organizer,
		pq.Array(&i.Tags),
		&i.ImageUrl,
		&i.MaxAttendees,
		&i.CurrentAttendees,
		&i.RegistrationRequired,
		&i.RegistrationDeadline,
		&i.Status,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OrganizerUsername,
		&i.OrganizerFullName,
		&i.OrganizerAvatar,
		&i.RegisteredCount,
		&i.AttendedCount,
	)
	return i, err
}

const getUpcomingEvents = `-- name: GetUpcomingEvents :many
SELECT 
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea8 WHERE ea8.event_id = e.id AND ea8.status = 'registered') as current_attendees_count
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $1 
  AND e.status = 'published'
  AND e.start_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'
ORDER BY e.start_date ASC
LIMIT 10
`

type GetUpcomingEventsRow struct {
	ID                    uuid.UUID      `json:"id"`
	SpaceID               uuid.UUID      `json:"space_id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	Category              string         `json:"category"`
	Location              sql.NullString `json:"location"`
	VenueDetails          sql.NullString `json:"venue_details"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
	Timezone              sql.NullString `json:"timezone"`
	Organizer             uuid.NullUUID  `json:"organizer"`
	Tags                  []string       `json:"tags"`
	ImageUrl              sql.NullString `json:"image_url"`
	MaxAttendees          sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees      sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired  sql.NullBool   `json:"registration_required"`
	RegistrationDeadline  sql.NullTime   `json:"registration_deadline"`
	Status                sql.NullString `json:"status"`
	IsPublic              sql.NullBool   `json:"is_public"`
	CreatedAt             sql.NullTime   `json:"created_at"`
	UpdatedAt             sql.NullTime   `json:"updated_at"`
	OrganizerUsername     string         `json:"organizer_username"`
	OrganizerFullName     string         `json:"organizer_full_name"`
	CurrentAttendeesCount int64          `json:"current_attendees_count"`
}

func (q *Queries) GetUpcomingEvents(ctx context.Context, spaceID uuid.UUID) ([]GetUpcomingEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUpcomingEvents, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUpcomingEventsRow{}
	for rows.Next() {
		var i GetUpcomingEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Location,
			&i.VenueDetails,
			&i.StartDate,
			&i.EndDate,
			&i.Timezone,
			&i.Organizer,
			pq.Array(&i.Tags),
			&i.ImageUrl,
			&i.MaxAttendees,
			&i.CurrentAttendees,
			&i.RegistrationRequired,
			&i.RegistrationDeadline,
			&i.Status,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrganizerUsername,
			&i.OrganizerFullName,
			&i.CurrentAttendeesCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserEvents = `-- name: GetUserEvents :many
SELECT 
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    ea.status as attendance_status
FROM events e
JOIN users u ON e.organizer = u.id
JOIN event_attendees ea ON e.id = ea.event_id
WHERE ea.user_id = $1 AND e.space_id = $2
ORDER BY e.start_date DESC
LIMIT $3 OFFSET $4
`

type GetUserEventsParams struct {
	UserID  uuid.UUID `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
	Limit   int32     `json:"limit"`
	Offset  int32     `json:"offset"`
}

type GetUserEventsRow struct {
	ID                   uuid.UUID      `json:"id"`
	SpaceID              uuid.UUID      `json:"space_id"`
	Title                string         `json:"title"`
	Description          sql.NullString `json:"description"`
	Category             string         `json:"category"`
	Location             sql.NullString `json:"location"`
	VenueDetails         sql.NullString `json:"venue_details"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	Timezone             sql.NullString `json:"timezone"`
	Organizer            uuid.NullUUID  `json:"organizer"`
	Tags                 []string       `json:"tags"`
	ImageUrl             sql.NullString `json:"image_url"`
	MaxAttendees         sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees     sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired sql.NullBool   `json:"registration_required"`
	RegistrationDeadline sql.NullTime   `json:"registration_deadline"`
	Status               sql.NullString `json:"status"`
	IsPublic             sql.NullBool   `json:"is_public"`
	CreatedAt            sql.NullTime   `json:"created_at"`
	UpdatedAt            sql.NullTime   `json:"updated_at"`
	OrganizerUsername    string         `json:"organizer_username"`
	OrganizerFullName    string         `json:"organizer_full_name"`
	AttendanceStatus     sql.NullString `json:"attendance_status"`
}

func (q *Queries) GetUserEvents(ctx context.Context, arg GetUserEventsParams) ([]GetUserEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserEvents,
		arg.UserID,
		arg.SpaceID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserEventsRow{}
	for rows.Next() {
		var i GetUserEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Location,
			&i.VenueDetails,
			&i.StartDate,
			&i.EndDate,
			&i.Timezone,
			&i.Organizer,
			pq.Array(&i.Tags),
			&i.ImageUrl,
			&i.MaxAttendees,
			&i.CurrentAttendees,
			&i.RegistrationRequired,
			&i.RegistrationDeadline,
			&i.Status,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrganizerUsername,
			&i.OrganizerFullName,
			&i.AttendanceStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAllAnnouncementsAdmin = `-- name: ListAllAnnouncementsAdmin :many

SELECT
    a.id, a.space_id, a.title, a.content, a.type, a.target_audience, a.priority, a.status, a.author_id, a.scheduled_for, a.expires_at, a.attachments, a.is_pinned, a.created_at, a.updated_at,
    u.username as author_username,
    u.full_name as author_full_name,
    u.avatar as author_avatar
FROM announcements a
LEFT JOIN users u ON a.author_id = u.id
WHERE a.space_id = $1
  OR (a.status = $2 OR $2 = '')
  OR (a.priority = $3 OR $3 = '')
ORDER BY a.created_at DESC
LIMIT $4 OFFSET $5
`

type ListAllAnnouncementsAdminParams struct {
	SpaceID  uuid.UUID      `json:"space_id"`
	Status   sql.NullString `json:"status"`
	Priority sql.NullString `json:"priority"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type ListAllAnnouncementsAdminRow struct {
	ID             uuid.UUID             `json:"id"`
	SpaceID        uuid.UUID             `json:"space_id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	TargetAudience []string              `json:"target_audience"`
	Priority       sql.NullString        `json:"priority"`
	Status         sql.NullString        `json:"status"`
	AuthorID       uuid.NullUUID         `json:"author_id"`
	ScheduledFor   sql.NullTime          `json:"scheduled_for"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	IsPinned       sql.NullBool          `json:"is_pinned"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	UpdatedAt      sql.NullTime          `json:"updated_at"`
	AuthorUsername sql.NullString        `json:"author_username"`
	AuthorFullName sql.NullString        `json:"author_full_name"`
	AuthorAvatar   sql.NullString        `json:"author_avatar"`
}


func (q *Queries) ListAllAnnouncementsAdmin(ctx context.Context, arg ListAllAnnouncementsAdminParams) ([]ListAllAnnouncementsAdminRow, error) {
	rows, err := q.db.QueryContext(ctx, listAllAnnouncementsAdmin,
		arg.SpaceID,
		arg.Status,
		arg.Priority,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAllAnnouncementsAdminRow{}
	for rows.Next() {
		var i ListAllAnnouncementsAdminRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Content,
			&i.Type,
			pq.Array(&i.TargetAudience),
			&i.Priority,
			&i.Status,
			&i.AuthorID,
			&i.ScheduledFor,
			&i.ExpiresAt,
			&i.Attachments,
			&i.IsPinned,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorUsername,
			&i.AuthorFullName,
			&i.AuthorAvatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAllEventsAdmin = `-- name: ListAllEventsAdmin :many

SELECT
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') as current_attendees_count
FROM events e
LEFT JOIN users u ON e.organizer = u.id
WHERE e.space_id = $1
  OR (e.status = $2 OR $2 = '')
  OR (e.category = $3 OR $3 = '')
ORDER BY e.created_at DESC
LIMIT $4 OFFSET $5
`

type ListAllEventsAdminParams struct {
	SpaceID  uuid.UUID      `json:"space_id"`
	Status   sql.NullString `json:"status"`
	Category string         `json:"category"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type ListAllEventsAdminRow struct {
	ID                    uuid.UUID      `json:"id"`
	SpaceID               uuid.UUID      `json:"space_id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	Category              string         `json:"category"`
	Location              sql.NullString `json:"location"`
	VenueDetails          sql.NullString `json:"venue_details"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
	Timezone              sql.NullString `json:"timezone"`
	Organizer             uuid.NullUUID  `json:"organizer"`
	Tags                  []string       `json:"tags"`
	ImageUrl              sql.NullString `json:"image_url"`
	MaxAttendees          sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees      sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired  sql.NullBool   `json:"registration_required"`
	RegistrationDeadline  sql.NullTime   `json:"registration_deadline"`
	Status                sql.NullString `json:"status"`
	IsPublic              sql.NullBool   `json:"is_public"`
	CreatedAt             sql.NullTime   `json:"created_at"`
	UpdatedAt             sql.NullTime   `json:"updated_at"`
	OrganizerUsername     sql.NullString `json:"organizer_username"`
	OrganizerFullName     sql.NullString `json:"organizer_full_name"`
	OrganizerAvatar       sql.NullString `json:"organizer_avatar"`
	CurrentAttendeesCount int64          `json:"current_attendees_count"`
}


func (q *Queries) ListAllEventsAdmin(ctx context.Context, arg ListAllEventsAdminParams) ([]ListAllEventsAdminRow, error) {
	rows, err := q.db.QueryContext(ctx, listAllEventsAdmin,
		arg.SpaceID,
		arg.Status,
		arg.Category,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAllEventsAdminRow{}
	for rows.Next() {
		var i ListAllEventsAdminRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Location,
			&i.VenueDetails,
			&i.StartDate,
			&i.EndDate,
			&i.Timezone,
			&i.Organizer,
			pq.Array(&i.Tags),
			&i.ImageUrl,
			&i.MaxAttendees,
			&i.CurrentAttendees,
			&i.RegistrationRequired,
			&i.RegistrationDeadline,
			&i.Status,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrganizerUsername,
			&i.OrganizerFullName,
			&i.OrganizerAvatar,
			&i.CurrentAttendeesCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAnnouncements = `-- name: ListAnnouncements :many
SELECT 
    a.id, a.space_id, a.title, a.content, a.type, a.target_audience, a.priority, a.status, a.author_id, a.scheduled_for, a.expires_at, a.attachments, a.is_pinned, a.created_at, a.updated_at,
    u.username as author_username,
    u.full_name as author_full_name
FROM announcements a
JOIN users u ON a.author_id = u.id
WHERE a.space_id = $1 
  AND a.status = 'published'
  AND (a.scheduled_for <= NOW() OR a.scheduled_for IS NULL)
  AND (a.expires_at >= NOW() OR a.expires_at IS NULL)
  AND (a.target_audience @> $2 OR $2 IS NULL)
ORDER BY 
    a.is_pinned DESC,
    a.priority DESC,
    a.created_at DESC
LIMIT $3 OFFSET $4
`

type ListAnnouncementsParams struct {
	SpaceID        uuid.UUID `json:"space_id"`
	TargetAudience []string  `json:"target_audience"`
	Limit          int32     `json:"limit"`
	Offset         int32     `json:"offset"`
}

type ListAnnouncementsRow struct {
	ID             uuid.UUID             `json:"id"`
	SpaceID        uuid.UUID             `json:"space_id"`
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	TargetAudience []string              `json:"target_audience"`
	Priority       sql.NullString        `json:"priority"`
	Status         sql.NullString        `json:"status"`
	AuthorID       uuid.NullUUID         `json:"author_id"`
	ScheduledFor   sql.NullTime          `json:"scheduled_for"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	IsPinned       sql.NullBool          `json:"is_pinned"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	UpdatedAt      sql.NullTime          `json:"updated_at"`
	AuthorUsername string                `json:"author_username"`
	AuthorFullName string                `json:"author_full_name"`
}

func (q *Queries) ListAnnouncements(ctx context.Context, arg ListAnnouncementsParams) ([]ListAnnouncementsRow, error) {
	rows, err := q.db.QueryContext(ctx, listAnnouncements,
		arg.SpaceID,
		pq.Array(arg.TargetAudience),
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAnnouncementsRow{}
	for rows.Next() {
		var i ListAnnouncementsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Content,
			&i.Type,
			pq.Array(&i.TargetAudience),
			&i.Priority,
			&i.Status,
			&i.AuthorID,
			&i.ScheduledFor,
			&i.ExpiresAt,
			&i.Attachments,
			&i.IsPinned,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AuthorUsername,
			&i.AuthorFullName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEvents = `-- name: ListEvents :many
SELECT
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea4 WHERE ea4.event_id = e.id AND ea4.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea5 WHERE ea5.event_id = e.id AND ea5.user_id = $1) as is_registered
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $2
  AND e.status = 'published'
  AND (e.is_public = true OR e.organizer = $1)
  AND (e.start_date >= $3 OR $3 IS NULL)
  AND (e.category = $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'upcoming' THEN e.start_date END ASC,
    CASE WHEN $5 = 'popular' THEN (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') END DESC,
    e.created_at DESC
LIMIT $6 OFFSET $7
`

type ListEventsParams struct {
	UserID    uuid.UUID   `json:"user_id"`
	SpaceID   uuid.UUID   `json:"space_id"`
	StartDate time.Time   `json:"start_date"`
	Category  string      `json:"category"`
	Column5   interface{} `json:"column_5"`
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
}

type ListEventsRow struct {
	ID                    uuid.UUID      `json:"id"`
	SpaceID               uuid.UUID      `json:"space_id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	Category              string         `json:"category"`
	Location              sql.NullString `json:"location"`
	VenueDetails          sql.NullString `json:"venue_details"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
	Timezone              sql.NullString `json:"timezone"`
	Organizer             uuid.NullUUID  `json:"organizer"`
	Tags                  []string       `json:"tags"`
	ImageUrl              sql.NullString `json:"image_url"`
	MaxAttendees          sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees      sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired  sql.NullBool   `json:"registration_required"`
	RegistrationDeadline  sql.NullTime   `json:"registration_deadline"`
	Status                sql.NullString `json:"status"`
	IsPublic              sql.NullBool   `json:"is_public"`
	CreatedAt             sql.NullTime   `json:"created_at"`
	UpdatedAt             sql.NullTime   `json:"updated_at"`
	OrganizerUsername     string         `json:"organizer_username"`
	OrganizerFullName     string         `json:"organizer_full_name"`
	CurrentAttendeesCount int64          `json:"current_attendees_count"`
	IsRegistered          bool           `json:"is_registered"`
}

func (q *Queries) ListEvents(ctx context.Context, arg ListEventsParams) ([]ListEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, listEvents,
		arg.UserID,
		arg.SpaceID,
		arg.StartDate,
		arg.Category,
		arg.Column5,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListEventsRow{}
	for rows.Next() {
		var i ListEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Location,
			&i.VenueDetails,
			&i.StartDate,
			&i.EndDate,
			&i.Timezone,
			&i.Organizer,
			pq.Array(&i.Tags),
			&i.ImageUrl,
			&i.MaxAttendees,
			&i.CurrentAttendees,
			&i.RegistrationRequired,
			&i.RegistrationDeadline,
			&i.Status,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrganizerUsername,
			&i.OrganizerFullName,
			&i.CurrentAttendeesCount,
			&i.IsRegistered,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markEventAttendance = `-- name: MarkEventAttendance :exec
UPDATE event_attendees 
SET status = 'attended', attended_at = NOW()
WHERE event_id = $1 AND user_id = $2
`

type MarkEventAttendanceParams struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) MarkEventAttendance(ctx context.Context, arg MarkEventAttendanceParams) error {
	_, err := q.db.ExecContext(ctx, markEventAttendance, arg.EventID, arg.UserID)
	return err
}

const registerForEvent = `-- name: RegisterForEvent :one
INSERT INTO event_attendees (event_id, user_id)
VALUES ($1, $2)
ON CONFLICT (event_id, user_id) 
DO UPDATE SET status = 'registered', registered_at = NOW()
RETURNING id, event_id, user_id, status, role, registered_at, attended_at, notes
`

type RegisterForEventParams struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) RegisterForEvent(ctx context.Context, arg RegisterForEventParams) (EventAttendee, error) {
	row := q.db.QueryRowContext(ctx, registerForEvent, arg.EventID, arg.UserID)
	var i EventAttendee
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.UserID,
		&i.Status,
		&i.Role,
		&i.RegisteredAt,
		&i.AttendedAt,
		&i.Notes,
	)
	return i, err
}

const removeEventCoOrganizer = `-- name: RemoveEventCoOrganizer :exec
DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2 AND role 'organizer'
`

type RemoveEventCoOrganizerParams struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) RemoveEventCoOrganizer(ctx context.Context, arg RemoveEventCoOrganizerParams) error {
	_, err := q.db.ExecContext(ctx, removeEventCoOrganizer, arg.EventID, arg.UserID)
	return err
}

const searchEvents = `-- name: SearchEvents :many
SELECT 
    e.id, e.space_id, e.title, e.description, e.category, e.location, e.venue_details, e.start_date, e.end_date, e.timezone, e.organizer, e.tags, e.image_url, e.max_attendees, e.current_attendees, e.registration_required, e.registration_deadline, e.status, e.is_public, e.created_at, e.updated_at,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea6 WHERE ea6.event_id = e.id AND ea6.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea7 WHERE ea7.event_id = e.id AND ea7.user_id = $1) as is_registered
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $2 
  AND e.status = 'published'
  AND (e.title ILIKE $3 OR e.description ILIKE $3 OR e.tags @> ARRAY[$3]::text[])
  AND (e.is_public = true OR e.organizer = $1)
ORDER BY e.start_date ASC
LIMIT 50
`

type SearchEventsParams struct {
	UserID  uuid.UUID `json:"user_id"`
	SpaceID uuid.UUID `json:"space_id"`
	Title   string    `json:"title"`
}

type SearchEventsRow struct {
	ID                    uuid.UUID      `json:"id"`
	SpaceID               uuid.UUID      `json:"space_id"`
	Title                 string         `json:"title"`
	Description           sql.NullString `json:"description"`
	Category              string         `json:"category"`
	Location              sql.NullString `json:"location"`
	VenueDetails          sql.NullString `json:"venue_details"`
	StartDate             time.Time      `json:"start_date"`
	EndDate               time.Time      `json:"end_date"`
	Timezone              sql.NullString `json:"timezone"`
	Organizer             uuid.NullUUID  `json:"organizer"`
	Tags                  []string       `json:"tags"`
	ImageUrl              sql.NullString `json:"image_url"`
	MaxAttendees          sql.NullInt32  `json:"max_attendees"`
	CurrentAttendees      sql.NullInt32  `json:"current_attendees"`
	RegistrationRequired  sql.NullBool   `json:"registration_required"`
	RegistrationDeadline  sql.NullTime   `json:"registration_deadline"`
	Status                sql.NullString `json:"status"`
	IsPublic              sql.NullBool   `json:"is_public"`
	CreatedAt             sql.NullTime   `json:"created_at"`
	UpdatedAt             sql.NullTime   `json:"updated_at"`
	OrganizerUsername     string         `json:"organizer_username"`
	OrganizerFullName     string         `json:"organizer_full_name"`
	CurrentAttendeesCount int64          `json:"current_attendees_count"`
	IsRegistered          bool           `json:"is_registered"`
}

func (q *Queries) SearchEvents(ctx context.Context, arg SearchEventsParams) ([]SearchEventsRow, error) {
	rows, err := q.db.QueryContext(ctx, searchEvents, arg.UserID, arg.SpaceID, arg.Title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchEventsRow{}
	for rows.Next() {
		var i SearchEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Title,
			&i.Description,
			&i.Category,
			&i.Location,
			&i.VenueDetails,
			&i.StartDate,
			&i.EndDate,
			&i.Timezone,
			&i.Organizer,
			pq.Array(&i.Tags),
			&i.ImageUrl,
			&i.MaxAttendees,
			&i.CurrentAttendees,
			&i.RegistrationRequired,
			&i.RegistrationDeadline,
			&i.Status,
			&i.IsPublic,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OrganizerUsername,
			&i.OrganizerFullName,
			&i.CurrentAttendeesCount,
			&i.IsRegistered,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const unregisterFromEvent = `-- name: UnregisterFromEvent :exec
DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2
`

type UnregisterFromEventParams struct {
	EventID uuid.UUID `json:"event_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) UnregisterFromEvent(ctx context.Context, arg UnregisterFromEventParams) error {
	_, err := q.db.ExecContext(ctx, unregisterFromEvent, arg.EventID, arg.UserID)
	return err
}

const updateAnnouncement = `-- name: UpdateAnnouncement :one
UPDATE announcements 
SET 
    title = $1,
    content = $2,
    type = $3,
    target_audience = $4,
    priority = $5,
    scheduled_for = $6,
    expires_at = $7,
    attachments = $8,
    is_pinned = $9,
    updated_at = NOW()
WHERE id = $10
RETURNING id, space_id, title, content, type, target_audience, priority, status, author_id, scheduled_for, expires_at, attachments, is_pinned, created_at, updated_at
`

type UpdateAnnouncementParams struct {
	Title          string                `json:"title"`
	Content        string                `json:"content"`
	Type           string                `json:"type"`
	TargetAudience []string              `json:"target_audience"`
	Priority       sql.NullString        `json:"priority"`
	ScheduledFor   sql.NullTime          `json:"scheduled_for"`
	ExpiresAt      sql.NullTime          `json:"expires_at"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	IsPinned       sql.NullBool          `json:"is_pinned"`
	ID             uuid.UUID             `json:"id"`
}

func (q *Queries) UpdateAnnouncement(ctx context.Context, arg UpdateAnnouncementParams) (Announcement, error) {
	row := q.db.QueryRowContext(ctx, updateAnnouncement,
		arg.Title,
		arg.Content,
		arg.Type,
		pq.Array(arg.TargetAudience),
		arg.Priority,
		arg.ScheduledFor,
		arg.ExpiresAt,
		arg.Attachments,
		arg.IsPinned,
		arg.ID,
	)
	var i Announcement
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Content,
		&i.Type,
		pq.Array(&i.TargetAudience),
		&i.Priority,
		&i.Status,
		&i.AuthorID,
		&i.ScheduledFor,
		&i.ExpiresAt,
		&i.Attachments,
		&i.IsPinned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAnnouncementStatus = `-- name: UpdateAnnouncementStatus :one
UPDATE announcements
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING id, space_id, title, content, type, target_audience, priority, status, author_id, scheduled_for, expires_at, attachments, is_pinned, created_at, updated_at
`

type UpdateAnnouncementStatusParams struct {
	Status sql.NullString `json:"status"`
	ID     uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateAnnouncementStatus(ctx context.Context, arg UpdateAnnouncementStatusParams) (Announcement, error) {
	row := q.db.QueryRowContext(ctx, updateAnnouncementStatus, arg.Status, arg.ID)
	var i Announcement
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Content,
		&i.Type,
		pq.Array(&i.TargetAudience),
		&i.Priority,
		&i.Status,
		&i.AuthorID,
		&i.ScheduledFor,
		&i.ExpiresAt,
		&i.Attachments,
		&i.IsPinned,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateEvent = `-- name: UpdateEvent :one
UPDATE events 
SET 
    title = $1,
    description = $2,
    category = $3,
    location = $4,
    venue_details = $5,
    start_date = $6,
    end_date = $7,
    timezone = $8,
    tags = $9,
    image_url = $10,
    max_attendees = $11,
    registration_required = $12,
    registration_deadline = $13,
    is_public = $14,
    updated_at = NOW()
WHERE id = $15
RETURNING id, space_id, title, description, category, location, venue_details, start_date, end_date, timezone, organizer, tags, image_url, max_attendees, current_attendees, registration_required, registration_deadline, status, is_public, created_at, updated_at
`

type UpdateEventParams struct {
	Title                string         `json:"title"`
	Description          sql.NullString `json:"description"`
	Category             string         `json:"category"`
	Location             sql.NullString `json:"location"`
	VenueDetails         sql.NullString `json:"venue_details"`
	StartDate            time.Time      `json:"start_date"`
	EndDate              time.Time      `json:"end_date"`
	Timezone             sql.NullString `json:"timezone"`
	Tags                 []string       `json:"tags"`
	ImageUrl             sql.NullString `json:"image_url"`
	MaxAttendees         sql.NullInt32  `json:"max_attendees"`
	RegistrationRequired sql.NullBool   `json:"registration_required"`
	RegistrationDeadline sql.NullTime   `json:"registration_deadline"`
	IsPublic             sql.NullBool   `json:"is_public"`
	ID                   uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateEvent(ctx context.Context, arg UpdateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, updateEvent,
		arg.Title,
		arg.Description,
		arg.Category,
		arg.Location,
		arg.VenueDetails,
		arg.StartDate,
		arg.EndDate,
		arg.Timezone,
		pq.Array(arg.Tags),
		arg.ImageUrl,
		arg.MaxAttendees,
		arg.RegistrationRequired,
		arg.RegistrationDeadline,
		arg.IsPublic,
		arg.ID,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Location,
		&i.VenueDetails,
		&i.StartDate,
		&i.EndDate,
		&i.Timezone,
		&i.Organizer,
		pq.Array(&i.Tags),
		&i.ImageUrl,
		&i.MaxAttendees,
		&i.CurrentAttendees,
		&i.RegistrationRequired,
		&i.RegistrationDeadline,
		&i.Status,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateEventStatus = `-- name: UpdateEventStatus :one
UPDATE events 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING id, space_id, title, description, category, location, venue_details, start_date, end_date, timezone, organizer, tags, image_url, max_attendees, current_attendees, registration_required, registration_deadline, status, is_public, created_at, updated_at
`

type UpdateEventStatusParams struct {
	Status sql.NullString `json:"status"`
	ID     uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateEventStatus(ctx context.Context, arg UpdateEventStatusParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, updateEventStatus, arg.Status, arg.ID)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Title,
		&i.Description,
		&i.Category,
		&i.Location,
		&i.VenueDetails,
		&i.StartDate,
		&i.EndDate,
		&i.Timezone,
		&i.Organizer,
		pq.Array(&i.Tags),
		&i.ImageUrl,
		&i.MaxAttendees,
		&i.CurrentAttendees,
		&i.RegistrationRequired,
		&i.RegistrationDeadline,
		&i.Status,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
