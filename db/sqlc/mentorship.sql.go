




package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const addMentoringSessionMeetingLink = `-- name: AddMentoringSessionMeetingLink :exec
UPDATE mentoring_sessions 
SET meeting_link = $1, updated_at = NOW()
WHERE id = $2
`

type AddMentoringSessionMeetingLinkParams struct {
	MeetingLink sql.NullString `json:"meeting_link"`
	ID          uuid.UUID      `json:"id"`
}

func (q *Queries) AddMentoringSessionMeetingLink(ctx context.Context, arg AddMentoringSessionMeetingLinkParams) error {
	_, err := q.db.ExecContext(ctx, addMentoringSessionMeetingLink, arg.MeetingLink, arg.ID)
	return err
}

const addSessionMeetingLink = `-- name: AddSessionMeetingLink :exec
UPDATE tutoring_sessions 
SET meeting_link = $1, updated_at = NOW()
WHERE id = $2
`

type AddSessionMeetingLinkParams struct {
	MeetingLink sql.NullString `json:"meeting_link"`
	ID          uuid.UUID      `json:"id"`
}

func (q *Queries) AddSessionMeetingLink(ctx context.Context, arg AddSessionMeetingLinkParams) error {
	_, err := q.db.ExecContext(ctx, addSessionMeetingLink, arg.MeetingLink, arg.ID)
	return err
}

const createMentorApplication = `-- name: CreateMentorApplication :one
INSERT INTO mentor_applications (
    applicant_id, space_id, industry, company, position, experience,
    specialties, achievements, mentorship_experience, availability,
    motivation, approach_description, linkedin_profile, portfolio
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id, applicant_id, space_id, industry, company, position, experience, specialties, achievements, mentorship_experience, availability, motivation, approach_description, linkedin_profile, portfolio, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type CreateMentorApplicationParams struct {
	ApplicantID          uuid.UUID       `json:"applicant_id"`
	SpaceID              uuid.UUID       `json:"space_id"`
	Industry             string          `json:"industry"`
	Company              sql.NullString  `json:"company"`
	Position             sql.NullString  `json:"position"`
	Experience           int32           `json:"experience"`
	Specialties          []string        `json:"specialties"`
	Achievements         sql.NullString  `json:"achievements"`
	MentorshipExperience sql.NullString  `json:"mentorship_experience"`
	Availability         json.RawMessage `json:"availability"`
	Motivation           sql.NullString  `json:"motivation"`
	ApproachDescription  sql.NullString  `json:"approach_description"`
	LinkedinProfile      sql.NullString  `json:"linkedin_profile"`
	Portfolio            sql.NullString  `json:"portfolio"`
}

func (q *Queries) CreateMentorApplication(ctx context.Context, arg CreateMentorApplicationParams) (MentorApplication, error) {
	row := q.db.QueryRowContext(ctx, createMentorApplication,
		arg.ApplicantID,
		arg.SpaceID,
		arg.Industry,
		arg.Company,
		arg.Position,
		arg.Experience,
		pq.Array(arg.Specialties),
		arg.Achievements,
		arg.MentorshipExperience,
		arg.Availability,
		arg.Motivation,
		arg.ApproachDescription,
		arg.LinkedinProfile,
		arg.Portfolio,
	)
	var i MentorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Achievements,
		&i.MentorshipExperience,
		&i.Availability,
		&i.Motivation,
		&i.ApproachDescription,
		&i.LinkedinProfile,
		&i.Portfolio,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const createMentorProfile = `-- name: CreateMentorProfile :one
INSERT INTO mentor_profiles (
    user_id, space_id, industry, company, position, experience,
    specialties, description, availability, verified
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (user_id) 
DO UPDATE SET
    industry = EXCLUDED.industry,
    company = EXCLUDED.company,
    position = EXCLUDED.position,
    experience = EXCLUDED.experience,
    specialties = EXCLUDED.specialties,
    description = EXCLUDED.description,
    availability = EXCLUDED.availability,
    verified = EXCLUDED.verified,
    updated_at = NOW()
RETURNING id, user_id, space_id, industry, company, position, experience, specialties, rating, review_count, total_sessions, availability, description, verified, is_available, created_at, updated_at
`

type CreateMentorProfileParams struct {
	UserID       uuid.UUID             `json:"user_id"`
	SpaceID      uuid.UUID             `json:"space_id"`
	Industry     string                `json:"industry"`
	Company      sql.NullString        `json:"company"`
	Position     sql.NullString        `json:"position"`
	Experience   int32                 `json:"experience"`
	Specialties  []string              `json:"specialties"`
	Description  sql.NullString        `json:"description"`
	Availability pqtype.NullRawMessage `json:"availability"`
	Verified     sql.NullBool          `json:"verified"`
}

func (q *Queries) CreateMentorProfile(ctx context.Context, arg CreateMentorProfileParams) (MentorProfile, error) {
	row := q.db.QueryRowContext(ctx, createMentorProfile,
		arg.UserID,
		arg.SpaceID,
		arg.Industry,
		arg.Company,
		arg.Position,
		arg.Experience,
		pq.Array(arg.Specialties),
		arg.Description,
		arg.Availability,
		arg.Verified,
	)
	var i MentorProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Availability,
		&i.Description,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createMentoringSession = `-- name: CreateMentoringSession :one
INSERT INTO mentoring_sessions (
    mentor_id, mentee_id, space_id, topic, scheduled_at,
    duration, mentee_notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, mentor_id, mentee_id, space_id, topic, status, scheduled_at, duration, mentee_notes, mentor_notes, meeting_link, rating, review, created_at, updated_at
`

type CreateMentoringSessionParams struct {
	MentorID    uuid.UUID      `json:"mentor_id"`
	MenteeID    uuid.UUID      `json:"mentee_id"`
	SpaceID     uuid.UUID      `json:"space_id"`
	Topic       string         `json:"topic"`
	ScheduledAt time.Time      `json:"scheduled_at"`
	Duration    int32          `json:"duration"`
	MenteeNotes sql.NullString `json:"mentee_notes"`
}

func (q *Queries) CreateMentoringSession(ctx context.Context, arg CreateMentoringSessionParams) (MentoringSession, error) {
	row := q.db.QueryRowContext(ctx, createMentoringSession,
		arg.MentorID,
		arg.MenteeID,
		arg.SpaceID,
		arg.Topic,
		arg.ScheduledAt,
		arg.Duration,
		arg.MenteeNotes,
	)
	var i MentoringSession
	err := row.Scan(
		&i.ID,
		&i.MentorID,
		&i.MenteeID,
		&i.SpaceID,
		&i.Topic,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.MenteeNotes,
		&i.MentorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createTutorApplication = `-- name: CreateTutorApplication :one

INSERT INTO tutor_applications (
    applicant_id, space_id, subjects, hourly_rate, availability,
    experience, qualifications, teaching_style, motivation, reference_letters
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id, applicant_id, space_id, subjects, hourly_rate, availability, experience, qualifications, teaching_style, motivation, reference_letters, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type CreateTutorApplicationParams struct {
	ApplicantID      uuid.UUID       `json:"applicant_id"`
	SpaceID          uuid.UUID       `json:"space_id"`
	Subjects         []string        `json:"subjects"`
	HourlyRate       sql.NullString  `json:"hourly_rate"`
	Availability     json.RawMessage `json:"availability"`
	Experience       sql.NullString  `json:"experience"`
	Qualifications   sql.NullString  `json:"qualifications"`
	TeachingStyle    sql.NullString  `json:"teaching_style"`
	Motivation       sql.NullString  `json:"motivation"`
	ReferenceLetters sql.NullString  `json:"reference_letters"`
}


func (q *Queries) CreateTutorApplication(ctx context.Context, arg CreateTutorApplicationParams) (TutorApplication, error) {
	row := q.db.QueryRowContext(ctx, createTutorApplication,
		arg.ApplicantID,
		arg.SpaceID,
		pq.Array(arg.Subjects),
		arg.HourlyRate,
		arg.Availability,
		arg.Experience,
		arg.Qualifications,
		arg.TeachingStyle,
		arg.Motivation,
		arg.ReferenceLetters,
	)
	var i TutorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.TeachingStyle,
		&i.Motivation,
		&i.ReferenceLetters,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const createTutorProfile = `-- name: CreateTutorProfile :one
INSERT INTO tutor_profiles (
    user_id, space_id, subjects, hourly_rate, description,
    availability, experience, qualifications, verified
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id) 
DO UPDATE SET
    subjects = EXCLUDED.subjects,
    hourly_rate = EXCLUDED.hourly_rate,
    description = EXCLUDED.description,
    availability = EXCLUDED.availability,
    experience = EXCLUDED.experience,
    qualifications = EXCLUDED.qualifications,
    verified = EXCLUDED.verified,
    updated_at = NOW()
RETURNING id, user_id, space_id, subjects, hourly_rate, rating, review_count, total_sessions, description, availability, experience, qualifications, verified, is_available, created_at, updated_at
`

type CreateTutorProfileParams struct {
	UserID         uuid.UUID             `json:"user_id"`
	SpaceID        uuid.UUID             `json:"space_id"`
	Subjects       []string              `json:"subjects"`
	HourlyRate     sql.NullString        `json:"hourly_rate"`
	Description    sql.NullString        `json:"description"`
	Availability   pqtype.NullRawMessage `json:"availability"`
	Experience     sql.NullString        `json:"experience"`
	Qualifications sql.NullString        `json:"qualifications"`
	Verified       sql.NullBool          `json:"verified"`
}

func (q *Queries) CreateTutorProfile(ctx context.Context, arg CreateTutorProfileParams) (TutorProfile, error) {
	row := q.db.QueryRowContext(ctx, createTutorProfile,
		arg.UserID,
		arg.SpaceID,
		pq.Array(arg.Subjects),
		arg.HourlyRate,
		arg.Description,
		arg.Availability,
		arg.Experience,
		arg.Qualifications,
		arg.Verified,
	)
	var i TutorProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Description,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createTutoringSession = `-- name: CreateTutoringSession :one
INSERT INTO tutoring_sessions (
    tutor_id, student_id, space_id, subject, scheduled_at,
    duration, hourly_rate, total_amount, student_notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, tutor_id, student_id, space_id, subject, status, scheduled_at, duration, hourly_rate, total_amount, student_notes, tutor_notes, meeting_link, rating, review, created_at, updated_at
`

type CreateTutoringSessionParams struct {
	TutorID      uuid.UUID      `json:"tutor_id"`
	StudentID    uuid.UUID      `json:"student_id"`
	SpaceID      uuid.UUID      `json:"space_id"`
	Subject      string         `json:"subject"`
	ScheduledAt  time.Time      `json:"scheduled_at"`
	Duration     int32          `json:"duration"`
	HourlyRate   sql.NullString `json:"hourly_rate"`
	TotalAmount  sql.NullString `json:"total_amount"`
	StudentNotes sql.NullString `json:"student_notes"`
}

func (q *Queries) CreateTutoringSession(ctx context.Context, arg CreateTutoringSessionParams) (TutoringSession, error) {
	row := q.db.QueryRowContext(ctx, createTutoringSession,
		arg.TutorID,
		arg.StudentID,
		arg.SpaceID,
		arg.Subject,
		arg.ScheduledAt,
		arg.Duration,
		arg.HourlyRate,
		arg.TotalAmount,
		arg.StudentNotes,
	)
	var i TutoringSession
	err := row.Scan(
		&i.ID,
		&i.TutorID,
		&i.StudentID,
		&i.SpaceID,
		&i.Subject,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.HourlyRate,
		&i.TotalAmount,
		&i.StudentNotes,
		&i.TutorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllMentorApplications = `-- name: GetAllMentorApplications :many
SELECT 
    mt.id,
    mt.applicant_id,
    mt.industry,
    mt.experience,
    mt.specialties,
    mt.status,
    mt.company,
    mt.position,
    mt.submitted_at,
    u.full_name,
    u.id as user_id
FROM mentor_applications mt
JOIN users u ON mt.applicant_id = u.id
ORDER BY submitted_at DESC
LIMIT $1 OFFSET $2
`

type GetAllMentorApplicationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAllMentorApplicationsRow struct {
	ID          uuid.UUID      `json:"id"`
	ApplicantID uuid.UUID      `json:"applicant_id"`
	Industry    string         `json:"industry"`
	Experience  int32          `json:"experience"`
	Specialties []string       `json:"specialties"`
	Status      sql.NullString `json:"status"`
	Company     sql.NullString `json:"company"`
	Position    sql.NullString `json:"position"`
	SubmittedAt sql.NullTime   `json:"submitted_at"`
	FullName    string         `json:"full_name"`
	UserID      uuid.UUID      `json:"user_id"`
}

func (q *Queries) GetAllMentorApplications(ctx context.Context, arg GetAllMentorApplicationsParams) ([]GetAllMentorApplicationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMentorApplications, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllMentorApplicationsRow{}
	for rows.Next() {
		var i GetAllMentorApplicationsRow
		if err := rows.Scan(
			&i.ID,
			&i.ApplicantID,
			&i.Industry,
			&i.Experience,
			pq.Array(&i.Specialties),
			&i.Status,
			&i.Company,
			&i.Position,
			&i.SubmittedAt,
			&i.FullName,
			&i.UserID,
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

const getAllTutorApplications = `-- name: GetAllTutorApplications :many

SELECT
    ta.applicant_id,
    ta.id,
    ta.subjects,
    ta.hourly_rate,
    ta.status,
    ta.submitted_at,
    u.full_name,
    u.id as user_id
 FROM tutor_applications ta
 JOIN users u ON ta.applicant_id = u.id
ORDER BY submitted_at DESC
LIMIT $1 OFFSET $2
`

type GetAllTutorApplicationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetAllTutorApplicationsRow struct {
	ApplicantID uuid.UUID      `json:"applicant_id"`
	ID          uuid.UUID      `json:"id"`
	Subjects    []string       `json:"subjects"`
	HourlyRate  sql.NullString `json:"hourly_rate"`
	Status      sql.NullString `json:"status"`
	SubmittedAt sql.NullTime   `json:"submitted_at"`
	FullName    string         `json:"full_name"`
	UserID      uuid.UUID      `json:"user_id"`
}


func (q *Queries) GetAllTutorApplications(ctx context.Context, arg GetAllTutorApplicationsParams) ([]GetAllTutorApplicationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllTutorApplications, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllTutorApplicationsRow{}
	for rows.Next() {
		var i GetAllTutorApplicationsRow
		if err := rows.Scan(
			&i.ApplicantID,
			&i.ID,
			pq.Array(&i.Subjects),
			&i.HourlyRate,
			&i.Status,
			&i.SubmittedAt,
			&i.FullName,
			&i.UserID,
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

const getMentorApplication = `-- name: GetMentorApplication :one
SELECT 
    ma.id, ma.applicant_id, ma.space_id, ma.industry, ma.company, ma.position, ma.experience, ma.specialties, ma.achievements, ma.mentorship_experience, ma.availability, ma.motivation, ma.approach_description, ma.linkedin_profile, ma.portfolio, ma.status, ma.submitted_at, ma.reviewed_at, ma.reviewed_by, ma.reviewer_notes,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM mentor_applications ma
JOIN users u ON ma.applicant_id = u.id
WHERE ma.id = $1
`

type GetMentorApplicationRow struct {
	ID                   uuid.UUID       `json:"id"`
	ApplicantID          uuid.UUID       `json:"applicant_id"`
	SpaceID              uuid.UUID       `json:"space_id"`
	Industry             string          `json:"industry"`
	Company              sql.NullString  `json:"company"`
	Position             sql.NullString  `json:"position"`
	Experience           int32           `json:"experience"`
	Specialties          []string        `json:"specialties"`
	Achievements         sql.NullString  `json:"achievements"`
	MentorshipExperience sql.NullString  `json:"mentorship_experience"`
	Availability         json.RawMessage `json:"availability"`
	Motivation           sql.NullString  `json:"motivation"`
	ApproachDescription  sql.NullString  `json:"approach_description"`
	LinkedinProfile      sql.NullString  `json:"linkedin_profile"`
	Portfolio            sql.NullString  `json:"portfolio"`
	Status               sql.NullString  `json:"status"`
	SubmittedAt          sql.NullTime    `json:"submitted_at"`
	ReviewedAt           sql.NullTime    `json:"reviewed_at"`
	ReviewedBy           uuid.NullUUID   `json:"reviewed_by"`
	ReviewerNotes        sql.NullString  `json:"reviewer_notes"`
	Username             string          `json:"username"`
	FullName             string          `json:"full_name"`
	Avatar               sql.NullString  `json:"avatar"`
	Email                string          `json:"email"`
	Department           sql.NullString  `json:"department"`
	Level                sql.NullString  `json:"level"`
}

func (q *Queries) GetMentorApplication(ctx context.Context, id uuid.UUID) (GetMentorApplicationRow, error) {
	row := q.db.QueryRowContext(ctx, getMentorApplication, id)
	var i GetMentorApplicationRow
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Achievements,
		&i.MentorshipExperience,
		&i.Availability,
		&i.Motivation,
		&i.ApproachDescription,
		&i.LinkedinProfile,
		&i.Portfolio,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
		&i.Username,
		&i.FullName,
		&i.Avatar,
		&i.Email,
		&i.Department,
		&i.Level,
	)
	return i, err
}

const getMentorProfile = `-- name: GetMentorProfile :one
SELECT 
    mp.id, mp.user_id, mp.space_id, mp.industry, mp.company, mp.position, mp.experience, mp.specialties, mp.rating, mp.review_count, mp.total_sessions, mp.availability, mp.description, mp.verified, mp.is_available, mp.created_at, mp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.user_id = $1
`

type GetMentorProfileRow struct {
	ID            uuid.UUID             `json:"id"`
	UserID        uuid.UUID             `json:"user_id"`
	SpaceID       uuid.UUID             `json:"space_id"`
	Industry      string                `json:"industry"`
	Company       sql.NullString        `json:"company"`
	Position      sql.NullString        `json:"position"`
	Experience    int32                 `json:"experience"`
	Specialties   []string              `json:"specialties"`
	Rating        sql.NullFloat64       `json:"rating"`
	ReviewCount   sql.NullInt32         `json:"review_count"`
	TotalSessions sql.NullInt32         `json:"total_sessions"`
	Availability  pqtype.NullRawMessage `json:"availability"`
	Description   sql.NullString        `json:"description"`
	Verified      sql.NullBool          `json:"verified"`
	IsAvailable   sql.NullBool          `json:"is_available"`
	CreatedAt     sql.NullTime          `json:"created_at"`
	UpdatedAt     sql.NullTime          `json:"updated_at"`
	Username      string                `json:"username"`
	FullName      string                `json:"full_name"`
	Avatar        sql.NullString        `json:"avatar"`
	UserVerified  sql.NullBool          `json:"user_verified"`
	Department    sql.NullString        `json:"department"`
	Level         sql.NullString        `json:"level"`
}

func (q *Queries) GetMentorProfile(ctx context.Context, userID uuid.UUID) (GetMentorProfileRow, error) {
	row := q.db.QueryRowContext(ctx, getMentorProfile, userID)
	var i GetMentorProfileRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Availability,
		&i.Description,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.FullName,
		&i.Avatar,
		&i.UserVerified,
		&i.Department,
		&i.Level,
	)
	return i, err
}

const getMentorReviews = `-- name: GetMentorReviews :many
SELECT
    ms.rating,
    ms.review,
    ms.created_at,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE ms.mentor_id = $1 AND ms.rating IS NOT NULL
ORDER BY ms.created_at DESC
LIMIT $2 OFFSET $3
`

type GetMentorReviewsParams struct {
	MentorID uuid.UUID `json:"mentor_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type GetMentorReviewsRow struct {
	Rating         sql.NullInt32  `json:"rating"`
	Review         sql.NullString `json:"review"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	MenteeUsername string         `json:"mentee_username"`
	MenteeFullName string         `json:"mentee_full_name"`
	MenteeAvatar   sql.NullString `json:"mentee_avatar"`
}

func (q *Queries) GetMentorReviews(ctx context.Context, arg GetMentorReviewsParams) ([]GetMentorReviewsRow, error) {
	rows, err := q.db.QueryContext(ctx, getMentorReviews, arg.MentorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetMentorReviewsRow{}
	for rows.Next() {
		var i GetMentorReviewsRow
		if err := rows.Scan(
			&i.Rating,
			&i.Review,
			&i.CreatedAt,
			&i.MenteeUsername,
			&i.MenteeFullName,
			&i.MenteeAvatar,
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

const getMentoringSession = `-- name: GetMentoringSession :one
SELECT 
    ms.id, ms.mentor_id, ms.mentee_id, ms.space_id, ms.topic, ms.status, ms.scheduled_at, ms.duration, ms.mentee_notes, ms.mentor_notes, ms.meeting_link, ms.rating, ms.review, ms.created_at, ms.updated_at,
    mentor.username as mentor_username,
    mentor.full_name as mentor_full_name,
    mentor.avatar as mentor_avatar,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentor ON ms.mentor_id = mentor.id
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE ms.id = $1
`

type GetMentoringSessionRow struct {
	ID             uuid.UUID      `json:"id"`
	MentorID       uuid.UUID      `json:"mentor_id"`
	MenteeID       uuid.UUID      `json:"mentee_id"`
	SpaceID        uuid.UUID      `json:"space_id"`
	Topic          string         `json:"topic"`
	Status         sql.NullString `json:"status"`
	ScheduledAt    time.Time      `json:"scheduled_at"`
	Duration       int32          `json:"duration"`
	MenteeNotes    sql.NullString `json:"mentee_notes"`
	MentorNotes    sql.NullString `json:"mentor_notes"`
	MeetingLink    sql.NullString `json:"meeting_link"`
	Rating         sql.NullInt32  `json:"rating"`
	Review         sql.NullString `json:"review"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	MentorUsername string         `json:"mentor_username"`
	MentorFullName string         `json:"mentor_full_name"`
	MentorAvatar   sql.NullString `json:"mentor_avatar"`
	MenteeUsername string         `json:"mentee_username"`
	MenteeFullName string         `json:"mentee_full_name"`
	MenteeAvatar   sql.NullString `json:"mentee_avatar"`
}

func (q *Queries) GetMentoringSession(ctx context.Context, id uuid.UUID) (GetMentoringSessionRow, error) {
	row := q.db.QueryRowContext(ctx, getMentoringSession, id)
	var i GetMentoringSessionRow
	err := row.Scan(
		&i.ID,
		&i.MentorID,
		&i.MenteeID,
		&i.SpaceID,
		&i.Topic,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.MenteeNotes,
		&i.MentorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.MentorUsername,
		&i.MentorFullName,
		&i.MentorAvatar,
		&i.MenteeUsername,
		&i.MenteeFullName,
		&i.MenteeAvatar,
	)
	return i, err
}

const getPendingMentorApplications = `-- name: GetPendingMentorApplications :many
SELECT 
    ma.id, ma.applicant_id, ma.space_id, ma.industry, ma.company, ma.position, ma.experience, ma.specialties, ma.achievements, ma.mentorship_experience, ma.availability, ma.motivation, ma.approach_description, ma.linkedin_profile, ma.portfolio, ma.status, ma.submitted_at, ma.reviewed_at, ma.reviewed_by, ma.reviewer_notes,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM mentor_applications ma
JOIN users u ON ma.applicant_id = u.id
WHERE ma.space_id = $1 AND ma.status = 'pending'
ORDER BY ma.submitted_at DESC
`

type GetPendingMentorApplicationsRow struct {
	ID                   uuid.UUID       `json:"id"`
	ApplicantID          uuid.UUID       `json:"applicant_id"`
	SpaceID              uuid.UUID       `json:"space_id"`
	Industry             string          `json:"industry"`
	Company              sql.NullString  `json:"company"`
	Position             sql.NullString  `json:"position"`
	Experience           int32           `json:"experience"`
	Specialties          []string        `json:"specialties"`
	Achievements         sql.NullString  `json:"achievements"`
	MentorshipExperience sql.NullString  `json:"mentorship_experience"`
	Availability         json.RawMessage `json:"availability"`
	Motivation           sql.NullString  `json:"motivation"`
	ApproachDescription  sql.NullString  `json:"approach_description"`
	LinkedinProfile      sql.NullString  `json:"linkedin_profile"`
	Portfolio            sql.NullString  `json:"portfolio"`
	Status               sql.NullString  `json:"status"`
	SubmittedAt          sql.NullTime    `json:"submitted_at"`
	ReviewedAt           sql.NullTime    `json:"reviewed_at"`
	ReviewedBy           uuid.NullUUID   `json:"reviewed_by"`
	ReviewerNotes        sql.NullString  `json:"reviewer_notes"`
	Username             string          `json:"username"`
	FullName             string          `json:"full_name"`
	Avatar               sql.NullString  `json:"avatar"`
	Email                string          `json:"email"`
	Department           sql.NullString  `json:"department"`
	Level                sql.NullString  `json:"level"`
}

func (q *Queries) GetPendingMentorApplications(ctx context.Context, spaceID uuid.UUID) ([]GetPendingMentorApplicationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPendingMentorApplications, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPendingMentorApplicationsRow{}
	for rows.Next() {
		var i GetPendingMentorApplicationsRow
		if err := rows.Scan(
			&i.ID,
			&i.ApplicantID,
			&i.SpaceID,
			&i.Industry,
			&i.Company,
			&i.Position,
			&i.Experience,
			pq.Array(&i.Specialties),
			&i.Achievements,
			&i.MentorshipExperience,
			&i.Availability,
			&i.Motivation,
			&i.ApproachDescription,
			&i.LinkedinProfile,
			&i.Portfolio,
			&i.Status,
			&i.SubmittedAt,
			&i.ReviewedAt,
			&i.ReviewedBy,
			&i.ReviewerNotes,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.Email,
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

const getPendingTutorApplications = `-- name: GetPendingTutorApplications :many
SELECT 
    ta.id, ta.applicant_id, ta.space_id, ta.subjects, ta.hourly_rate, ta.availability, ta.experience, ta.qualifications, ta.teaching_style, ta.motivation, ta.reference_letters, ta.status, ta.submitted_at, ta.reviewed_at, ta.reviewed_by, ta.reviewer_notes,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM tutor_applications ta
JOIN users u ON ta.applicant_id = u.id
WHERE ta.space_id = $1 AND ta.status = 'pending'
ORDER BY ta.submitted_at DESC
`

type GetPendingTutorApplicationsRow struct {
	ID               uuid.UUID       `json:"id"`
	ApplicantID      uuid.UUID       `json:"applicant_id"`
	SpaceID          uuid.UUID       `json:"space_id"`
	Subjects         []string        `json:"subjects"`
	HourlyRate       sql.NullString  `json:"hourly_rate"`
	Availability     json.RawMessage `json:"availability"`
	Experience       sql.NullString  `json:"experience"`
	Qualifications   sql.NullString  `json:"qualifications"`
	TeachingStyle    sql.NullString  `json:"teaching_style"`
	Motivation       sql.NullString  `json:"motivation"`
	ReferenceLetters sql.NullString  `json:"reference_letters"`
	Status           sql.NullString  `json:"status"`
	SubmittedAt      sql.NullTime    `json:"submitted_at"`
	ReviewedAt       sql.NullTime    `json:"reviewed_at"`
	ReviewedBy       uuid.NullUUID   `json:"reviewed_by"`
	ReviewerNotes    sql.NullString  `json:"reviewer_notes"`
	Username         string          `json:"username"`
	FullName         string          `json:"full_name"`
	Avatar           sql.NullString  `json:"avatar"`
	Email            string          `json:"email"`
	Department       sql.NullString  `json:"department"`
	Level            sql.NullString  `json:"level"`
}

func (q *Queries) GetPendingTutorApplications(ctx context.Context, spaceID uuid.UUID) ([]GetPendingTutorApplicationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPendingTutorApplications, spaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPendingTutorApplicationsRow{}
	for rows.Next() {
		var i GetPendingTutorApplicationsRow
		if err := rows.Scan(
			&i.ID,
			&i.ApplicantID,
			&i.SpaceID,
			pq.Array(&i.Subjects),
			&i.HourlyRate,
			&i.Availability,
			&i.Experience,
			&i.Qualifications,
			&i.TeachingStyle,
			&i.Motivation,
			&i.ReferenceLetters,
			&i.Status,
			&i.SubmittedAt,
			&i.ReviewedAt,
			&i.ReviewedBy,
			&i.ReviewerNotes,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.Email,
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

const getRecommendedMentors = `-- name: GetRecommendedMentors :many
SELECT
    mp.id, mp.user_id, mp.space_id, mp.industry, mp.company, mp.position, mp.experience, mp.specialties, mp.rating, mp.review_count, mp.total_sessions, mp.availability, mp.description, mp.verified, mp.is_available, mp.created_at, mp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    (SELECT AVG(rating) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND rating IS NOT NULL) as avg_rating,
    (SELECT COUNT(*) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND status = 'completed') as completed_sessions
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.space_id = $1
  AND mp.is_available = true
  AND mp.user_id != $2
  AND (
    -- Match by department
    u.department = (SELECT department FROM users WHERE id = $2)
    OR
    -- Match by industry relevant to user's major
    mp.industry IN (
      SELECT CASE
        WHEN major LIKE '%Computer%' OR major LIKE '%Engineering%' THEN 'Technology'
        WHEN major LIKE '%Business%' OR major LIKE '%Finance%' THEN 'Finance'
        WHEN major LIKE '%Art%' OR major LIKE '%Design%' THEN 'Creative'
        ELSE 'General'
      END
      FROM users WHERE id = $2
    )
    OR
    -- Match by specialties overlapping with interests
    mp.specialties && (SELECT interests FROM users WHERE id = $2)
  )
ORDER BY
    -- Prioritize verified mentors
    u.verified DESC,
    -- Then by rating
    (SELECT AVG(rating) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND rating IS NOT NULL) DESC NULLS LAST,
    -- Then by experience
    mp.experience DESC,
    -- Then by completed sessions
    (SELECT COUNT(*) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND status = 'completed') DESC,
    -- Finally by availability status
    mp.is_available DESC
LIMIT $3
`

type GetRecommendedMentorsParams struct {
	SpaceID uuid.UUID `json:"space_id"`
	UserID  uuid.UUID `json:"user_id"`
	Limit   int32     `json:"limit"`
}

type GetRecommendedMentorsRow struct {
	ID                uuid.UUID             `json:"id"`
	UserID            uuid.UUID             `json:"user_id"`
	SpaceID           uuid.UUID             `json:"space_id"`
	Industry          string                `json:"industry"`
	Company           sql.NullString        `json:"company"`
	Position          sql.NullString        `json:"position"`
	Experience        int32                 `json:"experience"`
	Specialties       []string              `json:"specialties"`
	Rating            sql.NullFloat64       `json:"rating"`
	ReviewCount       sql.NullInt32         `json:"review_count"`
	TotalSessions     sql.NullInt32         `json:"total_sessions"`
	Availability      pqtype.NullRawMessage `json:"availability"`
	Description       sql.NullString        `json:"description"`
	Verified          sql.NullBool          `json:"verified"`
	IsAvailable       sql.NullBool          `json:"is_available"`
	CreatedAt         sql.NullTime          `json:"created_at"`
	UpdatedAt         sql.NullTime          `json:"updated_at"`
	Username          string                `json:"username"`
	FullName          string                `json:"full_name"`
	Avatar            sql.NullString        `json:"avatar"`
	UserVerified      sql.NullBool          `json:"user_verified"`
	Department        sql.NullString        `json:"department"`
	Level             sql.NullString        `json:"level"`
	AvgRating         float64               `json:"avg_rating"`
	CompletedSessions int64                 `json:"completed_sessions"`
}

func (q *Queries) GetRecommendedMentors(ctx context.Context, arg GetRecommendedMentorsParams) ([]GetRecommendedMentorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRecommendedMentors, arg.SpaceID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRecommendedMentorsRow{}
	for rows.Next() {
		var i GetRecommendedMentorsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SpaceID,
			&i.Industry,
			&i.Company,
			&i.Position,
			&i.Experience,
			pq.Array(&i.Specialties),
			&i.Rating,
			&i.ReviewCount,
			&i.TotalSessions,
			&i.Availability,
			&i.Description,
			&i.Verified,
			&i.IsAvailable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.UserVerified,
			&i.Department,
			&i.Level,
			&i.AvgRating,
			&i.CompletedSessions,
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

const getRecommendedTutors = `-- name: GetRecommendedTutors :many

SELECT
    tp.id, tp.user_id, tp.space_id, tp.subjects, tp.hourly_rate, tp.rating, tp.review_count, tp.total_sessions, tp.description, tp.availability, tp.experience, tp.qualifications, tp.verified, tp.is_available, tp.created_at, tp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    (SELECT AVG(rating) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND rating IS NOT NULL) as avg_rating,
    (SELECT COUNT(*) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND status = 'completed') as completed_sessions
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.space_id = $1
  AND tp.is_available = true
  AND tp.user_id != $2
  AND (
    -- Match by department
    u.department = (SELECT department FROM users WHERE id = $2)
    OR
    -- Match by level
    u.level = (SELECT level FROM users WHERE id = $2)
    OR
    -- Match by subjects (if user has any in their profile)
    tp.subjects && (SELECT interests FROM users WHERE id = $2)
  )
ORDER BY
    -- Prioritize verified tutors
    u.verified DESC,
    -- Then by rating
    (SELECT AVG(rating) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND rating IS NOT NULL) DESC NULLS LAST,
    -- Then by completed sessions
    (SELECT COUNT(*) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND status = 'completed') DESC,
    -- Finally by availability status
    tp.is_available DESC
LIMIT $3
`

type GetRecommendedTutorsParams struct {
	SpaceID uuid.UUID `json:"space_id"`
	UserID  uuid.UUID `json:"user_id"`
	Limit   int32     `json:"limit"`
}

type GetRecommendedTutorsRow struct {
	ID                uuid.UUID             `json:"id"`
	UserID            uuid.UUID             `json:"user_id"`
	SpaceID           uuid.UUID             `json:"space_id"`
	Subjects          []string              `json:"subjects"`
	HourlyRate        sql.NullString        `json:"hourly_rate"`
	Rating            sql.NullFloat64       `json:"rating"`
	ReviewCount       sql.NullInt32         `json:"review_count"`
	TotalSessions     sql.NullInt32         `json:"total_sessions"`
	Description       sql.NullString        `json:"description"`
	Availability      pqtype.NullRawMessage `json:"availability"`
	Experience        sql.NullString        `json:"experience"`
	Qualifications    sql.NullString        `json:"qualifications"`
	Verified          sql.NullBool          `json:"verified"`
	IsAvailable       sql.NullBool          `json:"is_available"`
	CreatedAt         sql.NullTime          `json:"created_at"`
	UpdatedAt         sql.NullTime          `json:"updated_at"`
	Username          string                `json:"username"`
	FullName          string                `json:"full_name"`
	Avatar            sql.NullString        `json:"avatar"`
	UserVerified      sql.NullBool          `json:"user_verified"`
	Department        sql.NullString        `json:"department"`
	Level             sql.NullString        `json:"level"`
	AvgRating         float64               `json:"avg_rating"`
	CompletedSessions int64                 `json:"completed_sessions"`
}


func (q *Queries) GetRecommendedTutors(ctx context.Context, arg GetRecommendedTutorsParams) ([]GetRecommendedTutorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRecommendedTutors, arg.SpaceID, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetRecommendedTutorsRow{}
	for rows.Next() {
		var i GetRecommendedTutorsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SpaceID,
			pq.Array(&i.Subjects),
			&i.HourlyRate,
			&i.Rating,
			&i.ReviewCount,
			&i.TotalSessions,
			&i.Description,
			&i.Availability,
			&i.Experience,
			&i.Qualifications,
			&i.Verified,
			&i.IsAvailable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.UserVerified,
			&i.Department,
			&i.Level,
			&i.AvgRating,
			&i.CompletedSessions,
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

const getTutorApplication = `-- name: GetTutorApplication :one
SELECT 
    ta.id, ta.applicant_id, ta.space_id, ta.subjects, ta.hourly_rate, ta.availability, ta.experience, ta.qualifications, ta.teaching_style, ta.motivation, ta.reference_letters, ta.status, ta.submitted_at, ta.reviewed_at, ta.reviewed_by, ta.reviewer_notes,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM tutor_applications ta
JOIN users u ON ta.applicant_id = u.id
WHERE ta.id = $1
`

type GetTutorApplicationRow struct {
	ID               uuid.UUID       `json:"id"`
	ApplicantID      uuid.UUID       `json:"applicant_id"`
	SpaceID          uuid.UUID       `json:"space_id"`
	Subjects         []string        `json:"subjects"`
	HourlyRate       sql.NullString  `json:"hourly_rate"`
	Availability     json.RawMessage `json:"availability"`
	Experience       sql.NullString  `json:"experience"`
	Qualifications   sql.NullString  `json:"qualifications"`
	TeachingStyle    sql.NullString  `json:"teaching_style"`
	Motivation       sql.NullString  `json:"motivation"`
	ReferenceLetters sql.NullString  `json:"reference_letters"`
	Status           sql.NullString  `json:"status"`
	SubmittedAt      sql.NullTime    `json:"submitted_at"`
	ReviewedAt       sql.NullTime    `json:"reviewed_at"`
	ReviewedBy       uuid.NullUUID   `json:"reviewed_by"`
	ReviewerNotes    sql.NullString  `json:"reviewer_notes"`
	Username         string          `json:"username"`
	FullName         string          `json:"full_name"`
	Avatar           sql.NullString  `json:"avatar"`
	Email            string          `json:"email"`
	Department       sql.NullString  `json:"department"`
	Level            sql.NullString  `json:"level"`
}

func (q *Queries) GetTutorApplication(ctx context.Context, id uuid.UUID) (GetTutorApplicationRow, error) {
	row := q.db.QueryRowContext(ctx, getTutorApplication, id)
	var i GetTutorApplicationRow
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.TeachingStyle,
		&i.Motivation,
		&i.ReferenceLetters,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
		&i.Username,
		&i.FullName,
		&i.Avatar,
		&i.Email,
		&i.Department,
		&i.Level,
	)
	return i, err
}

const getTutorApplicationsByStatus = `-- name: GetTutorApplicationsByStatus :many
SELECT id, applicant_id, space_id, subjects, hourly_rate, availability, experience, qualifications, teaching_style, motivation, reference_letters, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes FROM tutor_applications
WHERE status = $1
ORDER BY submitted_at DESC
LIMIT $2 OFFSET $3
`

type GetTutorApplicationsByStatusParams struct {
	Status sql.NullString `json:"status"`
	Limit  int32          `json:"limit"`
	Offset int32          `json:"offset"`
}

func (q *Queries) GetTutorApplicationsByStatus(ctx context.Context, arg GetTutorApplicationsByStatusParams) ([]TutorApplication, error) {
	rows, err := q.db.QueryContext(ctx, getTutorApplicationsByStatus, arg.Status, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []TutorApplication{}
	for rows.Next() {
		var i TutorApplication
		if err := rows.Scan(
			&i.ID,
			&i.ApplicantID,
			&i.SpaceID,
			pq.Array(&i.Subjects),
			&i.HourlyRate,
			&i.Availability,
			&i.Experience,
			&i.Qualifications,
			&i.TeachingStyle,
			&i.Motivation,
			&i.ReferenceLetters,
			&i.Status,
			&i.SubmittedAt,
			&i.ReviewedAt,
			&i.ReviewedBy,
			&i.ReviewerNotes,
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

const getTutorProfile = `-- name: GetTutorProfile :one
SELECT 
    tp.id, tp.user_id, tp.space_id, tp.subjects, tp.hourly_rate, tp.rating, tp.review_count, tp.total_sessions, tp.description, tp.availability, tp.experience, tp.qualifications, tp.verified, tp.is_available, tp.created_at, tp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.user_id = $1
`

type GetTutorProfileRow struct {
	ID             uuid.UUID             `json:"id"`
	UserID         uuid.UUID             `json:"user_id"`
	SpaceID        uuid.UUID             `json:"space_id"`
	Subjects       []string              `json:"subjects"`
	HourlyRate     sql.NullString        `json:"hourly_rate"`
	Rating         sql.NullFloat64       `json:"rating"`
	ReviewCount    sql.NullInt32         `json:"review_count"`
	TotalSessions  sql.NullInt32         `json:"total_sessions"`
	Description    sql.NullString        `json:"description"`
	Availability   pqtype.NullRawMessage `json:"availability"`
	Experience     sql.NullString        `json:"experience"`
	Qualifications sql.NullString        `json:"qualifications"`
	Verified       sql.NullBool          `json:"verified"`
	IsAvailable    sql.NullBool          `json:"is_available"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	UpdatedAt      sql.NullTime          `json:"updated_at"`
	Username       string                `json:"username"`
	FullName       string                `json:"full_name"`
	Avatar         sql.NullString        `json:"avatar"`
	UserVerified   sql.NullBool          `json:"user_verified"`
	Department     sql.NullString        `json:"department"`
	Level          sql.NullString        `json:"level"`
}

func (q *Queries) GetTutorProfile(ctx context.Context, userID uuid.UUID) (GetTutorProfileRow, error) {
	row := q.db.QueryRowContext(ctx, getTutorProfile, userID)
	var i GetTutorProfileRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Description,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.FullName,
		&i.Avatar,
		&i.UserVerified,
		&i.Department,
		&i.Level,
	)
	return i, err
}

const getTutorReviews = `-- name: GetTutorReviews :many
SELECT 
    ts.rating,
    ts.review,
    ts.created_at,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users student ON ts.student_id = student.id
WHERE ts.tutor_id = $1 AND ts.rating IS NOT NULL
ORDER BY ts.created_at DESC
LIMIT $2 OFFSET $3
`

type GetTutorReviewsParams struct {
	TutorID uuid.UUID `json:"tutor_id"`
	Limit   int32     `json:"limit"`
	Offset  int32     `json:"offset"`
}

type GetTutorReviewsRow struct {
	Rating          sql.NullInt32  `json:"rating"`
	Review          sql.NullString `json:"review"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	StudentUsername string         `json:"student_username"`
	StudentFullName string         `json:"student_full_name"`
	StudentAvatar   sql.NullString `json:"student_avatar"`
}

func (q *Queries) GetTutorReviews(ctx context.Context, arg GetTutorReviewsParams) ([]GetTutorReviewsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTutorReviews, arg.TutorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTutorReviewsRow{}
	for rows.Next() {
		var i GetTutorReviewsRow
		if err := rows.Scan(
			&i.Rating,
			&i.Review,
			&i.CreatedAt,
			&i.StudentUsername,
			&i.StudentFullName,
			&i.StudentAvatar,
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

const getTutoringSession = `-- name: GetTutoringSession :one
SELECT 
    ts.id, ts.tutor_id, ts.student_id, ts.space_id, ts.subject, ts.status, ts.scheduled_at, ts.duration, ts.hourly_rate, ts.total_amount, ts.student_notes, ts.tutor_notes, ts.meeting_link, ts.rating, ts.review, ts.created_at, ts.updated_at,
    tutor.username as tutor_username,
    tutor.full_name as tutor_full_name,
    tutor.avatar as tutor_avatar,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users tutor ON ts.tutor_id = tutor.id
JOIN users student ON ts.student_id = student.id
WHERE ts.id = $1
`

type GetTutoringSessionRow struct {
	ID              uuid.UUID      `json:"id"`
	TutorID         uuid.UUID      `json:"tutor_id"`
	StudentID       uuid.UUID      `json:"student_id"`
	SpaceID         uuid.UUID      `json:"space_id"`
	Subject         string         `json:"subject"`
	Status          sql.NullString `json:"status"`
	ScheduledAt     time.Time      `json:"scheduled_at"`
	Duration        int32          `json:"duration"`
	HourlyRate      sql.NullString `json:"hourly_rate"`
	TotalAmount     sql.NullString `json:"total_amount"`
	StudentNotes    sql.NullString `json:"student_notes"`
	TutorNotes      sql.NullString `json:"tutor_notes"`
	MeetingLink     sql.NullString `json:"meeting_link"`
	Rating          sql.NullInt32  `json:"rating"`
	Review          sql.NullString `json:"review"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	TutorUsername   string         `json:"tutor_username"`
	TutorFullName   string         `json:"tutor_full_name"`
	TutorAvatar     sql.NullString `json:"tutor_avatar"`
	StudentUsername string         `json:"student_username"`
	StudentFullName string         `json:"student_full_name"`
	StudentAvatar   sql.NullString `json:"student_avatar"`
}

func (q *Queries) GetTutoringSession(ctx context.Context, id uuid.UUID) (GetTutoringSessionRow, error) {
	row := q.db.QueryRowContext(ctx, getTutoringSession, id)
	var i GetTutoringSessionRow
	err := row.Scan(
		&i.ID,
		&i.TutorID,
		&i.StudentID,
		&i.SpaceID,
		&i.Subject,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.HourlyRate,
		&i.TotalAmount,
		&i.StudentNotes,
		&i.TutorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TutorUsername,
		&i.TutorFullName,
		&i.TutorAvatar,
		&i.StudentUsername,
		&i.StudentFullName,
		&i.StudentAvatar,
	)
	return i, err
}

const getUserMentorApplicationStatus = `-- name: GetUserMentorApplicationStatus :one
SELECT status FROM mentor_applications
WHERE applicant_id = $1 AND space_id = $2
ORDER BY submitted_at DESC
LIMIT 1
`

type GetUserMentorApplicationStatusParams struct {
	ApplicantID uuid.UUID `json:"applicant_id"`
	SpaceID     uuid.UUID `json:"space_id"`
}

func (q *Queries) GetUserMentorApplicationStatus(ctx context.Context, arg GetUserMentorApplicationStatusParams) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getUserMentorApplicationStatus, arg.ApplicantID, arg.SpaceID)
	var status sql.NullString
	err := row.Scan(&status)
	return status, err
}

const getUserMentorApplicationStatusById = `-- name: GetUserMentorApplicationStatusById :one
SELECT status FROM mentor_applications
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserMentorApplicationStatusById(ctx context.Context, id uuid.UUID) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getUserMentorApplicationStatusById, id)
	var status sql.NullString
	err := row.Scan(&status)
	return status, err
}

const getUserMentoringSessions = `-- name: GetUserMentoringSessions :many
SELECT 
    ms.id, ms.mentor_id, ms.mentee_id, ms.space_id, ms.topic, ms.status, ms.scheduled_at, ms.duration, ms.mentee_notes, ms.mentor_notes, ms.meeting_link, ms.rating, ms.review, ms.created_at, ms.updated_at,
    mentor.username as mentor_username,
    mentor.full_name as mentor_full_name,
    mentor.avatar as mentor_avatar,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentor ON ms.mentor_id = mentor.id
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE (ms.mentor_id = $1 OR ms.mentee_id = $1)
ORDER BY ms.scheduled_at DESC
LIMIT $2 OFFSET $3
`

type GetUserMentoringSessionsParams struct {
	MentorID uuid.UUID `json:"mentor_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type GetUserMentoringSessionsRow struct {
	ID             uuid.UUID      `json:"id"`
	MentorID       uuid.UUID      `json:"mentor_id"`
	MenteeID       uuid.UUID      `json:"mentee_id"`
	SpaceID        uuid.UUID      `json:"space_id"`
	Topic          string         `json:"topic"`
	Status         sql.NullString `json:"status"`
	ScheduledAt    time.Time      `json:"scheduled_at"`
	Duration       int32          `json:"duration"`
	MenteeNotes    sql.NullString `json:"mentee_notes"`
	MentorNotes    sql.NullString `json:"mentor_notes"`
	MeetingLink    sql.NullString `json:"meeting_link"`
	Rating         sql.NullInt32  `json:"rating"`
	Review         sql.NullString `json:"review"`
	CreatedAt      sql.NullTime   `json:"created_at"`
	UpdatedAt      sql.NullTime   `json:"updated_at"`
	MentorUsername string         `json:"mentor_username"`
	MentorFullName string         `json:"mentor_full_name"`
	MentorAvatar   sql.NullString `json:"mentor_avatar"`
	MenteeUsername string         `json:"mentee_username"`
	MenteeFullName string         `json:"mentee_full_name"`
	MenteeAvatar   sql.NullString `json:"mentee_avatar"`
}

func (q *Queries) GetUserMentoringSessions(ctx context.Context, arg GetUserMentoringSessionsParams) ([]GetUserMentoringSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserMentoringSessions, arg.MentorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserMentoringSessionsRow{}
	for rows.Next() {
		var i GetUserMentoringSessionsRow
		if err := rows.Scan(
			&i.ID,
			&i.MentorID,
			&i.MenteeID,
			&i.SpaceID,
			&i.Topic,
			&i.Status,
			&i.ScheduledAt,
			&i.Duration,
			&i.MenteeNotes,
			&i.MentorNotes,
			&i.MeetingLink,
			&i.Rating,
			&i.Review,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.MentorUsername,
			&i.MentorFullName,
			&i.MentorAvatar,
			&i.MenteeUsername,
			&i.MenteeFullName,
			&i.MenteeAvatar,
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

const getUserTutorApplicationStatus = `-- name: GetUserTutorApplicationStatus :one
SELECT status FROM tutor_applications
WHERE applicant_id = $1 AND space_id = $2
ORDER BY submitted_at DESC
LIMIT 1
`

type GetUserTutorApplicationStatusParams struct {
	ApplicantID uuid.UUID `json:"applicant_id"`
	SpaceID     uuid.UUID `json:"space_id"`
}

func (q *Queries) GetUserTutorApplicationStatus(ctx context.Context, arg GetUserTutorApplicationStatusParams) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getUserTutorApplicationStatus, arg.ApplicantID, arg.SpaceID)
	var status sql.NullString
	err := row.Scan(&status)
	return status, err
}

const getUserTutorApplicationStatusById = `-- name: GetUserTutorApplicationStatusById :one
SELECT status FROM tutor_applications
WHERE id =$1
LIMIT 1
`

func (q *Queries) GetUserTutorApplicationStatusById(ctx context.Context, id uuid.UUID) (sql.NullString, error) {
	row := q.db.QueryRowContext(ctx, getUserTutorApplicationStatusById, id)
	var status sql.NullString
	err := row.Scan(&status)
	return status, err
}

const getUserTutoringSessions = `-- name: GetUserTutoringSessions :many
SELECT 
    ts.id, ts.tutor_id, ts.student_id, ts.space_id, ts.subject, ts.status, ts.scheduled_at, ts.duration, ts.hourly_rate, ts.total_amount, ts.student_notes, ts.tutor_notes, ts.meeting_link, ts.rating, ts.review, ts.created_at, ts.updated_at,
    tutor.username as tutor_username,
    tutor.full_name as tutor_full_name,
    tutor.avatar as tutor_avatar,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users tutor ON ts.tutor_id = tutor.id
JOIN users student ON ts.student_id = student.id
WHERE (ts.tutor_id = $1 OR ts.student_id = $1)
ORDER BY ts.scheduled_at DESC
LIMIT $2 OFFSET $3
`

type GetUserTutoringSessionsParams struct {
	TutorID uuid.UUID `json:"tutor_id"`
	Limit   int32     `json:"limit"`
	Offset  int32     `json:"offset"`
}

type GetUserTutoringSessionsRow struct {
	ID              uuid.UUID      `json:"id"`
	TutorID         uuid.UUID      `json:"tutor_id"`
	StudentID       uuid.UUID      `json:"student_id"`
	SpaceID         uuid.UUID      `json:"space_id"`
	Subject         string         `json:"subject"`
	Status          sql.NullString `json:"status"`
	ScheduledAt     time.Time      `json:"scheduled_at"`
	Duration        int32          `json:"duration"`
	HourlyRate      sql.NullString `json:"hourly_rate"`
	TotalAmount     sql.NullString `json:"total_amount"`
	StudentNotes    sql.NullString `json:"student_notes"`
	TutorNotes      sql.NullString `json:"tutor_notes"`
	MeetingLink     sql.NullString `json:"meeting_link"`
	Rating          sql.NullInt32  `json:"rating"`
	Review          sql.NullString `json:"review"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
	TutorUsername   string         `json:"tutor_username"`
	TutorFullName   string         `json:"tutor_full_name"`
	TutorAvatar     sql.NullString `json:"tutor_avatar"`
	StudentUsername string         `json:"student_username"`
	StudentFullName string         `json:"student_full_name"`
	StudentAvatar   sql.NullString `json:"student_avatar"`
}

func (q *Queries) GetUserTutoringSessions(ctx context.Context, arg GetUserTutoringSessionsParams) ([]GetUserTutoringSessionsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserTutoringSessions, arg.TutorID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserTutoringSessionsRow{}
	for rows.Next() {
		var i GetUserTutoringSessionsRow
		if err := rows.Scan(
			&i.ID,
			&i.TutorID,
			&i.StudentID,
			&i.SpaceID,
			&i.Subject,
			&i.Status,
			&i.ScheduledAt,
			&i.Duration,
			&i.HourlyRate,
			&i.TotalAmount,
			&i.StudentNotes,
			&i.TutorNotes,
			&i.MeetingLink,
			&i.Rating,
			&i.Review,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.TutorUsername,
			&i.TutorFullName,
			&i.TutorAvatar,
			&i.StudentUsername,
			&i.StudentFullName,
			&i.StudentAvatar,
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

const rateMentoringSession = `-- name: RateMentoringSession :one
UPDATE mentoring_sessions 
SET rating = $1, review = $2, updated_at = NOW()
WHERE id = $3 AND mentee_id = $4
RETURNING id, mentor_id, mentee_id, space_id, topic, status, scheduled_at, duration, mentee_notes, mentor_notes, meeting_link, rating, review, created_at, updated_at
`

type RateMentoringSessionParams struct {
	Rating   sql.NullInt32  `json:"rating"`
	Review   sql.NullString `json:"review"`
	ID       uuid.UUID      `json:"id"`
	MenteeID uuid.UUID      `json:"mentee_id"`
}

func (q *Queries) RateMentoringSession(ctx context.Context, arg RateMentoringSessionParams) (MentoringSession, error) {
	row := q.db.QueryRowContext(ctx, rateMentoringSession,
		arg.Rating,
		arg.Review,
		arg.ID,
		arg.MenteeID,
	)
	var i MentoringSession
	err := row.Scan(
		&i.ID,
		&i.MentorID,
		&i.MenteeID,
		&i.SpaceID,
		&i.Topic,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.MenteeNotes,
		&i.MentorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const rateTutoringSession = `-- name: RateTutoringSession :one
UPDATE tutoring_sessions 
SET rating = $1, review = $2, updated_at = NOW()
WHERE id = $3 AND student_id = $4
RETURNING id, tutor_id, student_id, space_id, subject, status, scheduled_at, duration, hourly_rate, total_amount, student_notes, tutor_notes, meeting_link, rating, review, created_at, updated_at
`

type RateTutoringSessionParams struct {
	Rating    sql.NullInt32  `json:"rating"`
	Review    sql.NullString `json:"review"`
	ID        uuid.UUID      `json:"id"`
	StudentID uuid.UUID      `json:"student_id"`
}

func (q *Queries) RateTutoringSession(ctx context.Context, arg RateTutoringSessionParams) (TutoringSession, error) {
	row := q.db.QueryRowContext(ctx, rateTutoringSession,
		arg.Rating,
		arg.Review,
		arg.ID,
		arg.StudentID,
	)
	var i TutoringSession
	err := row.Scan(
		&i.ID,
		&i.TutorID,
		&i.StudentID,
		&i.SpaceID,
		&i.Subject,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.HourlyRate,
		&i.TotalAmount,
		&i.StudentNotes,
		&i.TutorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const searchMentors = `-- name: SearchMentors :many
SELECT
    mp.id, mp.user_id, mp.space_id, mp.industry, mp.company, mp.position, mp.experience, mp.specialties, mp.rating, mp.review_count, mp.total_sessions, mp.availability, mp.description, mp.verified, mp.is_available, mp.created_at, mp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified AS user_verified,
    u.department,
    u.level,
    COALESCE(
        (SELECT AVG(rating)
         FROM mentoring_sessions
         WHERE mentor_id = mp.user_id
           AND rating IS NOT NULL), 0
    ) AS avg_rating,
    COALESCE(
        (SELECT COUNT(*)
         FROM mentoring_sessions
         WHERE mentor_id = mp.user_id
           AND status = 'completed'), 0
    ) AS completed_sessions
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.space_id = $1
  AND mp.is_available = true
  OR (mp.industry = $2 OR $2 IS NULL)
  OR (mp.specialties @> $3 OR $3 IS NULL)
  OR (mp.experience >= $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'rating' THEN
        COALESCE((SELECT AVG(rating)
                  FROM mentoring_sessions
                  WHERE mentor_id = mp.user_id
                    AND rating IS NOT NULL), 0)
    END DESC NULLS LAST,
    CASE WHEN $5 = 'experience' THEN mp.experience END DESC,
    COALESCE((SELECT COUNT(*)
              FROM mentoring_sessions
              WHERE mentor_id = mp.user_id
                AND status = 'completed'), 0) DESC
LIMIT $6 OFFSET $7
`

type SearchMentorsParams struct {
	SpaceID     uuid.UUID   `json:"space_id"`
	Industry    string      `json:"industry"`
	Specialties []string    `json:"specialties"`
	Experience  int32       `json:"experience"`
	Column5     interface{} `json:"column_5"`
	Limit       int32       `json:"limit"`
	Offset      int32       `json:"offset"`
}

type SearchMentorsRow struct {
	ID                uuid.UUID             `json:"id"`
	UserID            uuid.UUID             `json:"user_id"`
	SpaceID           uuid.UUID             `json:"space_id"`
	Industry          string                `json:"industry"`
	Company           sql.NullString        `json:"company"`
	Position          sql.NullString        `json:"position"`
	Experience        int32                 `json:"experience"`
	Specialties       []string              `json:"specialties"`
	Rating            sql.NullFloat64       `json:"rating"`
	ReviewCount       sql.NullInt32         `json:"review_count"`
	TotalSessions     sql.NullInt32         `json:"total_sessions"`
	Availability      pqtype.NullRawMessage `json:"availability"`
	Description       sql.NullString        `json:"description"`
	Verified          sql.NullBool          `json:"verified"`
	IsAvailable       sql.NullBool          `json:"is_available"`
	CreatedAt         sql.NullTime          `json:"created_at"`
	UpdatedAt         sql.NullTime          `json:"updated_at"`
	Username          string                `json:"username"`
	FullName          string                `json:"full_name"`
	Avatar            sql.NullString        `json:"avatar"`
	UserVerified      sql.NullBool          `json:"user_verified"`
	Department        sql.NullString        `json:"department"`
	Level             sql.NullString        `json:"level"`
	AvgRating         interface{}           `json:"avg_rating"`
	CompletedSessions interface{}           `json:"completed_sessions"`
}

func (q *Queries) SearchMentors(ctx context.Context, arg SearchMentorsParams) ([]SearchMentorsRow, error) {
	rows, err := q.db.QueryContext(ctx, searchMentors,
		arg.SpaceID,
		arg.Industry,
		pq.Array(arg.Specialties),
		arg.Experience,
		arg.Column5,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchMentorsRow{}
	for rows.Next() {
		var i SearchMentorsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SpaceID,
			&i.Industry,
			&i.Company,
			&i.Position,
			&i.Experience,
			pq.Array(&i.Specialties),
			&i.Rating,
			&i.ReviewCount,
			&i.TotalSessions,
			&i.Availability,
			&i.Description,
			&i.Verified,
			&i.IsAvailable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.UserVerified,
			&i.Department,
			&i.Level,
			&i.AvgRating,
			&i.CompletedSessions,
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

const searchTutors = `-- name: SearchTutors :many
SELECT
    tp.id, tp.user_id, tp.space_id, tp.subjects, tp.hourly_rate, tp.rating, tp.review_count, tp.total_sessions, tp.description, tp.availability, tp.experience, tp.qualifications, tp.verified, tp.is_available, tp.created_at, tp.updated_at,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    COALESCE(
        (SELECT AVG(rating)
         FROM tutoring_sessions
         WHERE tutor_id = tp.user_id
           AND rating IS NOT NULL), 0
    ) as avg_rating,
    COALESCE(
        (SELECT COUNT(*)
         FROM tutoring_sessions
         WHERE tutor_id = tp.user_id
           AND status = 'completed'), 0
    ) as completed_sessions
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.space_id = $1
  AND tp.is_available = true
  AND (tp.subjects @> $2 OR $2 IS NULL)
  AND (tp.availability @> $3 OR $3 IS NULL)
  AND (tp.hourly_rate <= $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'rating' THEN
        COALESCE((SELECT AVG(rating)
                  FROM tutoring_sessions
                  WHERE tutor_id = tp.user_id
                    AND rating IS NOT NULL), 0)
    END DESC NULLS LAST,
    CASE WHEN $5 = 'experience' THEN
        COALESCE((SELECT COUNT(*)
                  FROM tutoring_sessions
                  WHERE tutor_id = tp.user_id
                    AND status = 'completed'), 0)
    END DESC,
    tp.hourly_rate ASC
LIMIT $6 OFFSET $7
`

type SearchTutorsParams struct {
	SpaceID      uuid.UUID             `json:"space_id"`
	Subjects     []string              `json:"subjects"`
	Availability pqtype.NullRawMessage `json:"availability"`
	HourlyRate   sql.NullString        `json:"hourly_rate"`
	Column5      interface{}           `json:"column_5"`
	Limit        int32                 `json:"limit"`
	Offset       int32                 `json:"offset"`
}

type SearchTutorsRow struct {
	ID                uuid.UUID             `json:"id"`
	UserID            uuid.UUID             `json:"user_id"`
	SpaceID           uuid.UUID             `json:"space_id"`
	Subjects          []string              `json:"subjects"`
	HourlyRate        sql.NullString        `json:"hourly_rate"`
	Rating            sql.NullFloat64       `json:"rating"`
	ReviewCount       sql.NullInt32         `json:"review_count"`
	TotalSessions     sql.NullInt32         `json:"total_sessions"`
	Description       sql.NullString        `json:"description"`
	Availability      pqtype.NullRawMessage `json:"availability"`
	Experience        sql.NullString        `json:"experience"`
	Qualifications    sql.NullString        `json:"qualifications"`
	Verified          sql.NullBool          `json:"verified"`
	IsAvailable       sql.NullBool          `json:"is_available"`
	CreatedAt         sql.NullTime          `json:"created_at"`
	UpdatedAt         sql.NullTime          `json:"updated_at"`
	Username          string                `json:"username"`
	FullName          string                `json:"full_name"`
	Avatar            sql.NullString        `json:"avatar"`
	UserVerified      sql.NullBool          `json:"user_verified"`
	Department        sql.NullString        `json:"department"`
	Level             sql.NullString        `json:"level"`
	AvgRating         interface{}           `json:"avg_rating"`
	CompletedSessions interface{}           `json:"completed_sessions"`
}

func (q *Queries) SearchTutors(ctx context.Context, arg SearchTutorsParams) ([]SearchTutorsRow, error) {
	rows, err := q.db.QueryContext(ctx, searchTutors,
		arg.SpaceID,
		pq.Array(arg.Subjects),
		arg.Availability,
		arg.HourlyRate,
		arg.Column5,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchTutorsRow{}
	for rows.Next() {
		var i SearchTutorsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SpaceID,
			pq.Array(&i.Subjects),
			&i.HourlyRate,
			&i.Rating,
			&i.ReviewCount,
			&i.TotalSessions,
			&i.Description,
			&i.Availability,
			&i.Experience,
			&i.Qualifications,
			&i.Verified,
			&i.IsAvailable,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.UserVerified,
			&i.Department,
			&i.Level,
			&i.AvgRating,
			&i.CompletedSessions,
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

const updateMentorApplication = `-- name: UpdateMentorApplication :one
UPDATE mentor_applications 
SET 
    status = $1,
    reviewed_at = NOW(),
    reviewed_by = $2,
    reviewer_notes = $3
WHERE id = $4
RETURNING id, applicant_id, space_id, industry, company, position, experience, specialties, achievements, mentorship_experience, availability, motivation, approach_description, linkedin_profile, portfolio, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type UpdateMentorApplicationParams struct {
	Status        sql.NullString `json:"status"`
	ReviewedBy    uuid.NullUUID  `json:"reviewed_by"`
	ReviewerNotes sql.NullString `json:"reviewer_notes"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateMentorApplication(ctx context.Context, arg UpdateMentorApplicationParams) (MentorApplication, error) {
	row := q.db.QueryRowContext(ctx, updateMentorApplication,
		arg.Status,
		arg.ReviewedBy,
		arg.ReviewerNotes,
		arg.ID,
	)
	var i MentorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Achievements,
		&i.MentorshipExperience,
		&i.Availability,
		&i.Motivation,
		&i.ApproachDescription,
		&i.LinkedinProfile,
		&i.Portfolio,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const updateMentorApplicationStatus = `-- name: UpdateMentorApplicationStatus :one
UPDATE mentor_applications
SET
    status = $1,
    reviewed_by = $2,
    reviewer_notes = $3,
    reviewed_at = NOW()
WHERE id = $4
RETURNING id, applicant_id, space_id, industry, company, position, experience, specialties, achievements, mentorship_experience, availability, motivation, approach_description, linkedin_profile, portfolio, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type UpdateMentorApplicationStatusParams struct {
	Status        sql.NullString `json:"status"`
	ReviewedBy    uuid.NullUUID  `json:"reviewed_by"`
	ReviewerNotes sql.NullString `json:"reviewer_notes"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateMentorApplicationStatus(ctx context.Context, arg UpdateMentorApplicationStatusParams) (MentorApplication, error) {
	row := q.db.QueryRowContext(ctx, updateMentorApplicationStatus,
		arg.Status,
		arg.ReviewedBy,
		arg.ReviewerNotes,
		arg.ID,
	)
	var i MentorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Achievements,
		&i.MentorshipExperience,
		&i.Availability,
		&i.Motivation,
		&i.ApproachDescription,
		&i.LinkedinProfile,
		&i.Portfolio,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const updateMentorAvailability = `-- name: UpdateMentorAvailability :one
UPDATE mentor_profiles 
SET is_available = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING id, user_id, space_id, industry, company, position, experience, specialties, rating, review_count, total_sessions, availability, description, verified, is_available, created_at, updated_at
`

type UpdateMentorAvailabilityParams struct {
	IsAvailable sql.NullBool `json:"is_available"`
	UserID      uuid.UUID    `json:"user_id"`
}

func (q *Queries) UpdateMentorAvailability(ctx context.Context, arg UpdateMentorAvailabilityParams) (MentorProfile, error) {
	row := q.db.QueryRowContext(ctx, updateMentorAvailability, arg.IsAvailable, arg.UserID)
	var i MentorProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		&i.Industry,
		&i.Company,
		&i.Position,
		&i.Experience,
		pq.Array(&i.Specialties),
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Availability,
		&i.Description,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateMentoringSessionStatus = `-- name: UpdateMentoringSessionStatus :one
UPDATE mentoring_sessions 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING id, mentor_id, mentee_id, space_id, topic, status, scheduled_at, duration, mentee_notes, mentor_notes, meeting_link, rating, review, created_at, updated_at
`

type UpdateMentoringSessionStatusParams struct {
	Status sql.NullString `json:"status"`
	ID     uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateMentoringSessionStatus(ctx context.Context, arg UpdateMentoringSessionStatusParams) (MentoringSession, error) {
	row := q.db.QueryRowContext(ctx, updateMentoringSessionStatus, arg.Status, arg.ID)
	var i MentoringSession
	err := row.Scan(
		&i.ID,
		&i.MentorID,
		&i.MenteeID,
		&i.SpaceID,
		&i.Topic,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.MenteeNotes,
		&i.MentorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateSessionStatus = `-- name: UpdateSessionStatus :one
UPDATE tutoring_sessions 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING id, tutor_id, student_id, space_id, subject, status, scheduled_at, duration, hourly_rate, total_amount, student_notes, tutor_notes, meeting_link, rating, review, created_at, updated_at
`

type UpdateSessionStatusParams struct {
	Status sql.NullString `json:"status"`
	ID     uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateSessionStatus(ctx context.Context, arg UpdateSessionStatusParams) (TutoringSession, error) {
	row := q.db.QueryRowContext(ctx, updateSessionStatus, arg.Status, arg.ID)
	var i TutoringSession
	err := row.Scan(
		&i.ID,
		&i.TutorID,
		&i.StudentID,
		&i.SpaceID,
		&i.Subject,
		&i.Status,
		&i.ScheduledAt,
		&i.Duration,
		&i.HourlyRate,
		&i.TotalAmount,
		&i.StudentNotes,
		&i.TutorNotes,
		&i.MeetingLink,
		&i.Rating,
		&i.Review,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateTutorApplication = `-- name: UpdateTutorApplication :one
UPDATE tutor_applications 
SET 
    status = $1,
    reviewed_at = NOW(),
    reviewed_by = $2,
    reviewer_notes = $3
WHERE id = $4
RETURNING id, applicant_id, space_id, subjects, hourly_rate, availability, experience, qualifications, teaching_style, motivation, reference_letters, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type UpdateTutorApplicationParams struct {
	Status        sql.NullString `json:"status"`
	ReviewedBy    uuid.NullUUID  `json:"reviewed_by"`
	ReviewerNotes sql.NullString `json:"reviewer_notes"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateTutorApplication(ctx context.Context, arg UpdateTutorApplicationParams) (TutorApplication, error) {
	row := q.db.QueryRowContext(ctx, updateTutorApplication,
		arg.Status,
		arg.ReviewedBy,
		arg.ReviewerNotes,
		arg.ID,
	)
	var i TutorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.TeachingStyle,
		&i.Motivation,
		&i.ReferenceLetters,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const updateTutorApplicationStatus = `-- name: UpdateTutorApplicationStatus :one
UPDATE tutor_applications
SET
    status = $1,
    reviewed_by = $2,
    reviewer_notes = $3,
    reviewed_at = NOW()
WHERE id = $4
RETURNING id, applicant_id, space_id, subjects, hourly_rate, availability, experience, qualifications, teaching_style, motivation, reference_letters, status, submitted_at, reviewed_at, reviewed_by, reviewer_notes
`

type UpdateTutorApplicationStatusParams struct {
	Status        sql.NullString `json:"status"`
	ReviewedBy    uuid.NullUUID  `json:"reviewed_by"`
	ReviewerNotes sql.NullString `json:"reviewer_notes"`
	ID            uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateTutorApplicationStatus(ctx context.Context, arg UpdateTutorApplicationStatusParams) (TutorApplication, error) {
	row := q.db.QueryRowContext(ctx, updateTutorApplicationStatus,
		arg.Status,
		arg.ReviewedBy,
		arg.ReviewerNotes,
		arg.ID,
	)
	var i TutorApplication
	err := row.Scan(
		&i.ID,
		&i.ApplicantID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.TeachingStyle,
		&i.Motivation,
		&i.ReferenceLetters,
		&i.Status,
		&i.SubmittedAt,
		&i.ReviewedAt,
		&i.ReviewedBy,
		&i.ReviewerNotes,
	)
	return i, err
}

const updateTutorAvailability = `-- name: UpdateTutorAvailability :one
UPDATE tutor_profiles 
SET is_available = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING id, user_id, space_id, subjects, hourly_rate, rating, review_count, total_sessions, description, availability, experience, qualifications, verified, is_available, created_at, updated_at
`

type UpdateTutorAvailabilityParams struct {
	IsAvailable sql.NullBool `json:"is_available"`
	UserID      uuid.UUID    `json:"user_id"`
}

func (q *Queries) UpdateTutorAvailability(ctx context.Context, arg UpdateTutorAvailabilityParams) (TutorProfile, error) {
	row := q.db.QueryRowContext(ctx, updateTutorAvailability, arg.IsAvailable, arg.UserID)
	var i TutorProfile
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SpaceID,
		pq.Array(&i.Subjects),
		&i.HourlyRate,
		&i.Rating,
		&i.ReviewCount,
		&i.TotalSessions,
		&i.Description,
		&i.Availability,
		&i.Experience,
		&i.Qualifications,
		&i.Verified,
		&i.IsAvailable,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
