package mentorship

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

type Service struct {
	store db.Store
	db    *sql.DB
}

func NewService(store db.Store) *Service {
	// Try to get DB from store if it's a SQLStore
	var database *sql.DB
	if sqlStore, ok := store.(*db.SQLStore); ok {
		database = sqlStore.DB()
	}
	return &Service{
		store: store,
		db:    database,
	}
}

// ===== Mentor Profile Operations =====

func (s *Service) CreateMentorProfile(ctx context.Context, req CreateMentorProfileRequest) (*MentorProfileResponse, error) {
	params := db.CreateMentorProfileParams{
		UserID:       req.UserID,
		SpaceID:      req.SpaceID,
		Industry:     req.Industry,
		Company:      sqlNullString(req.Company),
		Position:     sqlNullString(req.Position),
		Experience:   req.Experience,
		Specialties:  req.Specialties,
		Description:  sqlNullString(req.Description),
		Availability: sqlNullRawMessage(req.Availability),
	}

	profile, err := s.store.CreateMentorProfile(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create mentor profile: %w", err)
	}

	return mentorProfileToResponse(profile), nil
}

func (s *Service) GetMentorProfile(ctx context.Context, profileID uuid.UUID) (*MentorProfileResponse, error) {
	profileRow, err := s.store.GetMentorProfile(ctx, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mentor profile not found")
		}
		return nil, fmt.Errorf("failed to get mentor profile: %w", err)
	}

	// Convert GetMentorProfileRow to MentorProfile
	profile := db.MentorProfile{
		ID:            profileRow.ID,
		UserID:        profileRow.UserID,
		SpaceID:       profileRow.SpaceID,
		Industry:      profileRow.Industry,
		Company:       profileRow.Company,
		Position:      profileRow.Position,
		Experience:    profileRow.Experience,
		Specialties:   profileRow.Specialties,
		Rating:        profileRow.Rating,
		ReviewCount:   profileRow.ReviewCount,
		TotalSessions: profileRow.TotalSessions,
		Availability:  profileRow.Availability,
		Description:   profileRow.Description,
		Verified:      profileRow.Verified,
		IsAvailable:   profileRow.IsAvailable,
		CreatedAt:     profileRow.CreatedAt,
		UpdatedAt:     profileRow.UpdatedAt,
	}

	return mentorProfileToResponse(profile), nil
}

func (s *Service) UpdateMentorAvailability(ctx context.Context, profileID uuid.UUID, req UpdateMentorAvailabilityRequest) (*MentorProfileResponse, error) {
	params := db.UpdateMentorAvailabilityParams{
		IsAvailable: sql.NullBool{Bool: req.IsAvailable, Valid: true},
		UserID:      profileID,
	}

	profile, err := s.store.UpdateMentorAvailability(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update mentor availability: %w", err)
	}

	return mentorProfileToResponse(profile), nil
}

func (s *Service) SearchMentors(ctx context.Context, params SearchMentorsParams) ([]MentorSearchResponse, error) {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	industry := ""
	if params.Industry != nil {
		industry = *params.Industry
	} 

	minRating := 0.0
	if params.MinRating != nil {
		minRating = *params.MinRating
	}

	mentors, err := s.store.SearchMentors(ctx, db.SearchMentorsParams{
		SpaceID:     params.SpaceID,
		Industry:    industry,
		Specialties: params.Specialties,
		Experience:  0,
		Column5:     minRating,
		Limit:       params.Limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search mentors: %w", err)
	}

	result := make([]MentorSearchResponse, len(mentors))
	for i, mentor := range mentors {
		result[i] = MentorSearchResponse{
			ID:          mentor.ID,
			UserID:      mentor.UserID,
			Username:    mentor.Username,
			FullName:    mentor.FullName,
			Avatar:      nullStringToPtr(mentor.Avatar),
			Industry:    mentor.Industry,
			Company:     nullStringToPtr(mentor.Company),
			Position:    nullStringToPtr(mentor.Position),
			Experience:  mentor.Experience,
			Specialties: mentor.Specialties,
			Rating:      nullFloat64ToPtr(mentor.Rating),
			ReviewCount: nullInt32ToPtr(mentor.ReviewCount),
			IsAvailable: nullBoolToPtr(mentor.IsAvailable),
		}
	}

	return result, nil
}

func (s *Service) GetMentorReviews(ctx context.Context, mentorID uuid.UUID) ([]MentorReviewResponse, error) {
	params := db.GetMentorReviewsParams{
		MentorID: mentorID,
		Limit:    100,
		Offset:   0,
	}

	reviews, err := s.store.GetMentorReviews(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get mentor reviews: %w", err)
	}

	result := make([]MentorReviewResponse, len(reviews))
	for i, review := range reviews {
		rating := int32(0)
		if review.Rating.Valid {
			rating = review.Rating.Int32
		}

		result[i] = MentorReviewResponse{
			ID:           uuid.Nil, // Not available in query result
			MentorID:     mentorID,
			MenteeID:     uuid.Nil, // Not available in query result
			MenteeName:   review.MenteeFullName,
			MenteeAvatar: nullStringToPtr(review.MenteeAvatar),
			SessionID:    uuid.Nil, // Not available in query result
			Topic:        "",       // Not available in query result
			Rating:       rating,
			Review:       nullStringToPtr(review.Review),
			CreatedAt:    nullTimeToPtr(review.CreatedAt),
		}
	}

	return result, nil
}

// ===== Tutor Profile Operations =====

func (s *Service) CreateTutorProfile(ctx context.Context, req CreateTutorProfileRequest) (*TutorProfileResponse, error) {
	params := db.CreateTutorProfileParams{
		UserID:         req.UserID,
		SpaceID:        req.SpaceID,
		Subjects:       req.Subjects,
		HourlyRate:     sqlNullString(req.HourlyRate),
		Description:    sqlNullString(req.Description),
		Experience:     sqlNullString(req.Experience),
		Qualifications: sqlNullString(req.Qualifications),
		Availability:   sqlNullRawMessage(req.Availability),
	}

	profile, err := s.store.CreateTutorProfile(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create tutor profile: %w", err)
	}

	return tutorProfileToResponse(profile), nil
}

func (s *Service) GetTutorProfile(ctx context.Context, profileID uuid.UUID) (*TutorProfileResponse, error) {
	profileRow, err := s.store.GetTutorProfile(ctx, profileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tutor profile not found")
		}
		return nil, fmt.Errorf("failed to get tutor profile: %w", err)
	}

	// Convert GetTutorProfileRow to TutorProfile
	profile := db.TutorProfile{
		ID:             profileRow.ID,
		UserID:         profileRow.UserID,
		SpaceID:        profileRow.SpaceID,
		Subjects:       profileRow.Subjects,
		HourlyRate:     profileRow.HourlyRate,
		Rating:         profileRow.Rating,
		ReviewCount:    profileRow.ReviewCount,
		TotalSessions:  profileRow.TotalSessions,
		Availability:   profileRow.Availability,
		Experience:     profileRow.Experience,
		Qualifications: profileRow.Qualifications,
		Description:    profileRow.Description,
		Verified:       profileRow.Verified,
		IsAvailable:    profileRow.IsAvailable,
		CreatedAt:      profileRow.CreatedAt,
		UpdatedAt:      profileRow.UpdatedAt,
	}

	return tutorProfileToResponse(profile), nil
}

func (s *Service) UpdateTutorAvailability(ctx context.Context, profileID uuid.UUID, req UpdateTutorAvailabilityRequest) (*TutorProfileResponse, error) {
	params := db.UpdateTutorAvailabilityParams{
		IsAvailable: sql.NullBool{Bool: req.IsAvailable, Valid: true},
		UserID:      profileID,
	}

	profile, err := s.store.UpdateTutorAvailability(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update tutor availability: %w", err)
	}

	return tutorProfileToResponse(profile), nil
}

func (s *Service) SearchTutors(ctx context.Context, params SearchTutorsParams) ([]TutorSearchResponse, error) {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	maxRate := ""
	if params.MaxRate != nil {
		maxRate = *params.MaxRate
	}

	minRating := 0.0
	if params.MinRating != nil {
		minRating = *params.MinRating
	}

	var hourlyRate sql.NullString
	if maxRate != "" {
		hourlyRate = sql.NullString{String: maxRate, Valid: true}
	}

	tutors, err := s.store.SearchTutors(ctx, db.SearchTutorsParams{
		SpaceID:      params.SpaceID,
		Subjects:     params.Subjects,
		Availability: pqtype.NullRawMessage{},
		HourlyRate:   hourlyRate,
		Column5:      minRating,
		Limit:        params.Limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search tutors: %w", err)
	}

	result := make([]TutorSearchResponse, len(tutors))
	for i, tutor := range tutors {
		result[i] = TutorSearchResponse{
			ID:             tutor.ID,
			UserID:         tutor.UserID,
			Username:       tutor.Username,
			FullName:       tutor.FullName,
			Avatar:         nullStringToPtr(tutor.Avatar),
			Subjects:       tutor.Subjects,
			HourlyRate:     nullStringToPtr(tutor.HourlyRate),
			Rating:         nullFloat64ToPtr(tutor.Rating),
			ReviewCount:    nullInt32ToPtr(tutor.ReviewCount),
			Experience:     nullStringToPtr(tutor.Experience),
			Qualifications: nullStringToPtr(tutor.Qualifications),
			IsAvailable:    nullBoolToPtr(tutor.IsAvailable),
		}
	}

	return result, nil
}

func (s *Service) GetTutorReviews(ctx context.Context, tutorID uuid.UUID) ([]TutorReviewResponse, error) {
	params := db.GetTutorReviewsParams{
		TutorID: tutorID,
		Limit:   100,
		Offset:  0,
	}

	reviews, err := s.store.GetTutorReviews(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get tutor reviews: %w", err)
	}

	result := make([]TutorReviewResponse, len(reviews))
	for i, review := range reviews {
		rating := int32(0)
		if review.Rating.Valid {
			rating = review.Rating.Int32
		}

		result[i] = TutorReviewResponse{
			ID:            uuid.Nil, // Not available in query result
			TutorID:       tutorID,
			StudentID:     uuid.Nil, // Not available in query result
			StudentName:   review.StudentFullName,
			StudentAvatar: nullStringToPtr(review.StudentAvatar),
			SessionID:     uuid.Nil, // Not available in query result
			Subject:       "",       // Not available in query result
			Rating:        rating,
			Review:        nullStringToPtr(review.Review),
			CreatedAt:     nullTimeToPtr(review.CreatedAt),
		}
	}

	return result, nil
}

// ===== Mentoring Session Operations =====

func (s *Service) CreateMentoringSession(ctx context.Context, req CreateMentoringSessionRequest) (*MentoringSessionResponse, error) {
	params := db.CreateMentoringSessionParams{
		MentorID:    req.MentorID,
		MenteeID:    req.MenteeID,
		SpaceID:     req.SpaceID,
		Topic:       req.Topic,
		ScheduledAt: req.ScheduledAt,
		Duration:    req.Duration,
		MenteeNotes: sqlNullString(req.MenteeNotes),
	}

	session, err := s.store.CreateMentoringSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create mentoring session: %w", err)
	}

	return mentoringSessionToResponse(session), nil
}

func (s *Service) GetMentoringSession(ctx context.Context, sessionID uuid.UUID) (*MentoringSessionDetailResponse, error) {
	session, err := s.store.GetMentoringSession(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mentoring session not found")
		}
		return nil, fmt.Errorf("failed to get mentoring session: %w", err)
	}

	return &MentoringSessionDetailResponse{
		ID:           session.ID,
		MentorID:     session.MentorID,
		MentorName:   session.MentorFullName,
		MentorAvatar: nullStringToPtr(session.MentorAvatar),
		MenteeID:     session.MenteeID,
		MenteeName:   session.MenteeFullName,
		MenteeAvatar: nullStringToPtr(session.MenteeAvatar),
		SpaceID:      session.SpaceID,
		Topic:        session.Topic,
		Status:       nullStringToPtr(session.Status),
		ScheduledAt:  session.ScheduledAt,
		Duration:     session.Duration,
		MenteeNotes:  nullStringToPtr(session.MenteeNotes),
		MentorNotes:  nullStringToPtr(session.MentorNotes),
		MeetingLink:  nullStringToPtr(session.MeetingLink),
		Rating:       nullInt32ToPtr(session.Rating),
		Review:       nullStringToPtr(session.Review),
		CreatedAt:    nullTimeToPtr(session.CreatedAt),
	}, nil
}

func (s *Service) GetUserMentoringSessions(ctx context.Context, userID, spaceID uuid.UUID, page, limit int32) ([]MentoringSessionDetailResponse, error) {
	if limit == 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	sessions, err := s.store.GetUserMentoringSessions(ctx, db.GetUserMentoringSessionsParams{
		MentorID: userID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user mentoring sessions: %w", err)
	}

	result := make([]MentoringSessionDetailResponse, len(sessions))
	for i, session := range sessions {
		result[i] = MentoringSessionDetailResponse{
			ID:           session.ID,
			MentorID:     session.MentorID,
			MentorName:   session.MentorFullName,
			MentorAvatar: nullStringToPtr(session.MentorAvatar),
			MenteeID:     session.MenteeID,
			MenteeName:   session.MenteeFullName,
			MenteeAvatar: nullStringToPtr(session.MenteeAvatar),
			SpaceID:      session.SpaceID,
			Topic:        session.Topic,
			Status:       nullStringToPtr(session.Status),
			ScheduledAt:  session.ScheduledAt,
			Duration:     session.Duration,
			MenteeNotes:  nullStringToPtr(session.MenteeNotes),
			MentorNotes:  nullStringToPtr(session.MentorNotes),
			MeetingLink:  nullStringToPtr(session.MeetingLink),
			Rating:       nullInt32ToPtr(session.Rating),
			Review:       nullStringToPtr(session.Review),
			CreatedAt:    nullTimeToPtr(session.CreatedAt),
		}
	}

	return result, nil
}

func (s *Service) UpdateMentoringSessionStatus(ctx context.Context, sessionID uuid.UUID, status string) error {
	_, err := s.store.UpdateMentoringSessionStatus(ctx, db.UpdateMentoringSessionStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to update mentoring session status: %w", err)
	}
	return nil
}

func (s *Service) AddMentoringSessionMeetingLink(ctx context.Context, sessionID uuid.UUID, meetingLink string) error {
	err := s.store.AddMentoringSessionMeetingLink(ctx, db.AddMentoringSessionMeetingLinkParams{
		MeetingLink: sql.NullString{String: meetingLink, Valid: true},
		ID:          sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to add meeting link: %w", err)
	}
	return nil
}

func (s *Service) RateMentoringSession(ctx context.Context, sessionID uuid.UUID, req RateMentoringSessionRequest) error {
	session, err := s.store.RateMentoringSession(ctx, db.RateMentoringSessionParams{
		Rating: sql.NullInt32{Int32: req.Rating, Valid: true},
		Review: sqlNullString(req.Review),
		ID:     sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to rate mentoring session: %w", err)
	}
	_ = session
	return nil
}

// ===== Tutoring Session Operations =====

func (s *Service) CreateTutoringSession(ctx context.Context, req CreateTutoringSessionRequest) (*TutoringSessionResponse, error) {
	params := db.CreateTutoringSessionParams{
		TutorID:      req.TutorID,
		StudentID:    req.StudentID,
		SpaceID:      req.SpaceID,
		Subject:      req.Subject,
		ScheduledAt:  req.ScheduledAt,
		Duration:     req.Duration,
		HourlyRate:   sqlNullString(req.HourlyRate),
		StudentNotes: sqlNullString(req.StudentNotes),
	}

	session, err := s.store.CreateTutoringSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create tutoring session: %w", err)
	}

	return tutoringSessionToResponse(session), nil
}

func (s *Service) GetTutoringSession(ctx context.Context, sessionID uuid.UUID) (*TutoringSessionDetailResponse, error) {
	session, err := s.store.GetTutoringSession(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tutoring session not found")
		}
		return nil, fmt.Errorf("failed to get tutoring session: %w", err)
	}

	return &TutoringSessionDetailResponse{
		ID:            session.ID,
		TutorID:       session.TutorID,
		TutorName:     session.TutorFullName,
		TutorAvatar:   nullStringToPtr(session.TutorAvatar),
		StudentID:     session.StudentID,
		StudentName:   session.StudentFullName,
		StudentAvatar: nullStringToPtr(session.StudentAvatar),
		SpaceID:       session.SpaceID,
		Subject:       session.Subject,
		Status:        nullStringToPtr(session.Status),
		ScheduledAt:   session.ScheduledAt,
		Duration:      session.Duration,
		HourlyRate:    nullStringToPtr(session.HourlyRate),
		TotalAmount:   nullStringToPtr(session.TotalAmount),
		StudentNotes:  nullStringToPtr(session.StudentNotes),
		TutorNotes:    nullStringToPtr(session.TutorNotes),
		MeetingLink:   nullStringToPtr(session.MeetingLink),
		Rating:        nullInt32ToPtr(session.Rating),
		Review:        nullStringToPtr(session.Review),
		CreatedAt:     nullTimeToPtr(session.CreatedAt),
	}, nil
}

func (s *Service) GetUserTutoringSessions(ctx context.Context, userID, spaceID uuid.UUID, page, limit int32) ([]TutoringSessionDetailResponse, error) {
	if limit == 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	sessions, err := s.store.GetUserTutoringSessions(ctx, db.GetUserTutoringSessionsParams{
		TutorID: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user tutoring sessions: %w", err)
	}

	result := make([]TutoringSessionDetailResponse, len(sessions))
	for i, session := range sessions {
		result[i] = TutoringSessionDetailResponse{
			ID:            session.ID,
			TutorID:       session.TutorID,
			TutorName:     session.TutorFullName,
			TutorAvatar:   nullStringToPtr(session.TutorAvatar),
			StudentID:     session.StudentID,
			StudentName:   session.StudentFullName,
			StudentAvatar: nullStringToPtr(session.StudentAvatar),
			SpaceID:       session.SpaceID,
			Subject:       session.Subject,
			Status:        nullStringToPtr(session.Status),
			ScheduledAt:   session.ScheduledAt,
			Duration:      session.Duration,
			HourlyRate:    nullStringToPtr(session.HourlyRate),
			TotalAmount:   nullStringToPtr(session.TotalAmount),
			StudentNotes:  nullStringToPtr(session.StudentNotes),
			TutorNotes:    nullStringToPtr(session.TutorNotes),
			MeetingLink:   nullStringToPtr(session.MeetingLink),
			Rating:        nullInt32ToPtr(session.Rating),
			Review:        nullStringToPtr(session.Review),
			CreatedAt:     nullTimeToPtr(session.CreatedAt),
		}
	}

	return result, nil
}

func (s *Service) UpdateSessionStatus(ctx context.Context, sessionID uuid.UUID, status string) error {
	_, err := s.store.UpdateSessionStatus(ctx, db.UpdateSessionStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to update tutoring session status: %w", err)
	}
	return nil
}

func (s *Service) AddSessionMeetingLink(ctx context.Context, sessionID uuid.UUID, meetingLink string) error {
	err := s.store.AddSessionMeetingLink(ctx, db.AddSessionMeetingLinkParams{
		MeetingLink: sql.NullString{String: meetingLink, Valid: true},
		ID:          sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to add meeting link: %w", err)
	}
	return nil
}

func (s *Service) RateTutoringSession(ctx context.Context, sessionID uuid.UUID, req RateTutoringSessionRequest) error {
	session, err := s.store.RateTutoringSession(ctx, db.RateTutoringSessionParams{
		Rating: sql.NullInt32{Int32: req.Rating, Valid: true},
		Review: sqlNullString(req.Review),
		ID:     sessionID,
	})
	if err != nil {
		return fmt.Errorf("failed to rate tutoring session: %w", err)
	}
	_ = session
	return nil
}

// ===== Mentor Application Operations =====

func (s *Service) CreateMentorApplication(ctx context.Context, req CreateMentorApplicationRequest) (*MentorApplicationResponse, error) {
	params := db.CreateMentorApplicationParams{
		ApplicantID:          req.UserID,
		SpaceID:              req.SpaceID,
		Industry:             req.Industry,
		Company:              sqlNullString(req.Company),
		Position:             sqlNullString(req.Position),
		Experience:           req.Experience,
		Specialties:          req.Specialties,
		Achievements:         sql.NullString{},
		MentorshipExperience: sql.NullString{},
		Availability:         []byte("{}"),
		Motivation:           sqlNullString(&req.Motivation),
		ApproachDescription:  sql.NullString{},
		LinkedinProfile:      sql.NullString{},
		Portfolio:            sql.NullString{},
	}

	application, err := s.store.CreateMentorApplication(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create mentor application: %w", err)
	}

	return mentorApplicationToResponse(application), nil
}

func (s *Service) GetMentorApplication(ctx context.Context, applicationID uuid.UUID) (*MentorApplicationResponse, error) {
	appRow, err := s.store.GetMentorApplication(ctx, applicationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mentor application not found")
		}
		return nil, fmt.Errorf("failed to get mentor application: %w", err)
	}

	// Convert Row to Application type
	application := db.MentorApplication{
		ID:                   appRow.ID,
		ApplicantID:          appRow.ApplicantID,
		SpaceID:              appRow.SpaceID,
		Industry:             appRow.Industry,
		Company:              appRow.Company,
		Position:             appRow.Position,
		Experience:           appRow.Experience,
		Specialties:          appRow.Specialties,
		Achievements:         appRow.Achievements,
		MentorshipExperience: appRow.MentorshipExperience,
		Availability:         appRow.Availability,
		Motivation:           appRow.Motivation,
		ApproachDescription:  appRow.ApproachDescription,
		LinkedinProfile:      appRow.LinkedinProfile,
		Portfolio:            appRow.Portfolio,
		Status:               sql.NullString{},
		SubmittedAt:          sql.NullTime{},
		ReviewedAt:           sql.NullTime{},
		ReviewedBy:           uuid.NullUUID{},
		ReviewerNotes:        sql.NullString{},
	}

	return mentorApplicationToResponse(application), nil
}

func (s *Service) UpdateMentorApplication(ctx context.Context, applicationID, reviewerID uuid.UUID, req UpdateMentorApplicationRequest) (*MentorApplicationResponse, error) {
	params := db.UpdateMentorApplicationParams{
		Status:        sql.NullString{String: req.Status, Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: reviewerID, Valid: true},
		ReviewerNotes: sqlNullString(req.ReviewComments),
		ID:            applicationID,
	}

	application, err := s.store.UpdateMentorApplication(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update mentor application: %w", err)
	}

	return mentorApplicationToResponse(application), nil
}

func (s *Service) GetPendingMentorApplications(ctx context.Context, spaceID uuid.UUID) ([]MentorApplicationResponse, error) {
	appRows, err := s.store.GetPendingMentorApplications(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending mentor applications: %w", err)
	}

	result := make([]MentorApplicationResponse, len(appRows))
	for i, appRow := range appRows {
		// Convert Row to Application type
		app := db.MentorApplication{
			ID:                   appRow.ID,
			ApplicantID:          appRow.ApplicantID,
			SpaceID:              appRow.SpaceID,
			Industry:             appRow.Industry,
			Company:              appRow.Company,
			Position:             appRow.Position,
			Experience:           appRow.Experience,
			Specialties:          appRow.Specialties,
			Achievements:         appRow.Achievements,
			MentorshipExperience: appRow.MentorshipExperience,
			Availability:         appRow.Availability,
			Motivation:           appRow.Motivation,
			ApproachDescription:  appRow.ApproachDescription,
			LinkedinProfile:      appRow.LinkedinProfile,
			Portfolio:            appRow.Portfolio,
			Status:               sql.NullString{},
			SubmittedAt:          sql.NullTime{},
			ReviewedAt:           sql.NullTime{},
			ReviewedBy:           uuid.NullUUID{},
			ReviewerNotes:        sql.NullString{},
		}
		result[i] = *mentorApplicationToResponse(app)
	}

	return result, nil
}

// ===== Tutor Application Operations =====

func (s *Service) CreateTutorApplication(ctx context.Context, req CreateTutorApplicationRequest) (*TutorApplicationResponse, error) {
	params := db.CreateTutorApplicationParams{
		ApplicantID:      req.UserID,
		SpaceID:          req.SpaceID,
		Subjects:         req.Subjects,
		HourlyRate:       sql.NullString{},
		Availability:     []byte("{}"),
		Experience:       sqlNullString(&req.Experience),
		Qualifications:   sqlNullString(&req.Qualifications),
		TeachingStyle:    sql.NullString{},
		Motivation:       sqlNullString(&req.Motivation),
		ReferenceLetters: sql.NullString{},
	}

	application, err := s.store.CreateTutorApplication(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create tutor application: %w", err)
	}

	return tutorApplicationToResponse(application), nil
}

func (s *Service) GetTutorApplication(ctx context.Context, applicationID uuid.UUID) (*TutorApplicationResponse, error) {
	appRow, err := s.store.GetTutorApplication(ctx, applicationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tutor application not found")
		}
		return nil, fmt.Errorf("failed to get tutor application: %w", err)
	}

	// Convert Row to Application type
	application := db.TutorApplication{
		ID:               appRow.ID,
		ApplicantID:      appRow.ApplicantID,
		SpaceID:          appRow.SpaceID,
		Subjects:         appRow.Subjects,
		HourlyRate:       appRow.HourlyRate,
		Availability:     appRow.Availability,
		Experience:       appRow.Experience,
		Qualifications:   appRow.Qualifications,
		TeachingStyle:    appRow.TeachingStyle,
		Motivation:       appRow.Motivation,
		ReferenceLetters: appRow.ReferenceLetters,
		Status:           appRow.Status,
		SubmittedAt:      appRow.SubmittedAt,
		ReviewedAt:       appRow.ReviewedAt,
		ReviewedBy:       appRow.ReviewedBy,
		ReviewerNotes:    sql.NullString{},
	}

	return tutorApplicationToResponse(application), nil
}

func (s *Service) UpdateTutorApplication(ctx context.Context, applicationID, reviewerID uuid.UUID, req UpdateTutorApplicationRequest) (*TutorApplicationResponse, error) {
	params := db.UpdateTutorApplicationParams{
		Status:        sql.NullString{String: req.Status, Valid: true},
		ReviewedBy:    uuid.NullUUID{UUID: reviewerID, Valid: true},
		ReviewerNotes: sqlNullString(req.ReviewComments),
		ID:            applicationID,
	}

	application, err := s.store.UpdateTutorApplication(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update tutor application: %w", err)
	}

	return tutorApplicationToResponse(application), nil
}

func (s *Service) GetPendingTutorApplications(ctx context.Context, spaceID uuid.UUID) ([]TutorApplicationResponse, error) {
	applications, err := s.store.GetPendingTutorApplications(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending tutor applications: %w", err)
	}

	result := make([]TutorApplicationResponse, len(applications))
	for i, appRow := range applications {
		// Convert Row to Application type
		app := db.TutorApplication{
			ID:               appRow.ID,
			ApplicantID:      appRow.ApplicantID,
			SpaceID:          appRow.SpaceID,
			Subjects:         appRow.Subjects,
			HourlyRate:       appRow.HourlyRate,
			Availability:     appRow.Availability,
			Experience:       appRow.Experience,
			Qualifications:   appRow.Qualifications,
			TeachingStyle:    appRow.TeachingStyle,
			Motivation:       appRow.Motivation,
			ReferenceLetters: appRow.ReferenceLetters,
			Status:           sql.NullString{},
			SubmittedAt:      sql.NullTime{},
			ReviewedAt:       sql.NullTime{},
			ReviewedBy:       uuid.NullUUID{},
			ReviewerNotes:    sql.NullString{},
		}
		result[i] = *tutorApplicationToResponse(app)
	}

	return result, nil
}

// ===== Helper Functions =====

func mentorProfileToResponse(profile db.MentorProfile) *MentorProfileResponse {
	return &MentorProfileResponse{
		ID:            profile.ID,
		UserID:        profile.UserID,
		SpaceID:       profile.SpaceID,
		Industry:      profile.Industry,
		Company:       nullStringToPtr(profile.Company),
		Position:      nullStringToPtr(profile.Position),
		Experience:    profile.Experience,
		Specialties:   profile.Specialties,
		Rating:        nullFloat64ToPtr(profile.Rating),
		ReviewCount:   nullInt32ToPtr(profile.ReviewCount),
		TotalSessions: nullInt32ToPtr(profile.TotalSessions),
		Availability:  nullRawMessageToPtr(profile.Availability),
		Description:   nullStringToPtr(profile.Description),
		Verified:      nullBoolToPtr(profile.Verified),
		IsAvailable:   nullBoolToPtr(profile.IsAvailable),
		CreatedAt:     nullTimeToPtr(profile.CreatedAt),
		UpdatedAt:     nullTimeToPtr(profile.UpdatedAt),
	}
}

func tutorProfileToResponse(profile db.TutorProfile) *TutorProfileResponse {
	return &TutorProfileResponse{
		ID:             profile.ID,
		UserID:         profile.UserID,
		SpaceID:        profile.SpaceID,
		Subjects:       profile.Subjects,
		HourlyRate:     nullStringToPtr(profile.HourlyRate),
		Rating:         nullFloat64ToPtr(profile.Rating),
		ReviewCount:    nullInt32ToPtr(profile.ReviewCount),
		TotalSessions:  nullInt32ToPtr(profile.TotalSessions),
		Description:    nullStringToPtr(profile.Description),
		Availability:   nullRawMessageToPtr(profile.Availability),
		Experience:     nullStringToPtr(profile.Experience),
		Qualifications: nullStringToPtr(profile.Qualifications),
		Verified:       nullBoolToPtr(profile.Verified),
		IsAvailable:    nullBoolToPtr(profile.IsAvailable),
		CreatedAt:      nullTimeToPtr(profile.CreatedAt),
		UpdatedAt:      nullTimeToPtr(profile.UpdatedAt),
	}
}

func mentoringSessionToResponse(session db.MentoringSession) *MentoringSessionResponse {
	return &MentoringSessionResponse{
		ID:          session.ID,
		MentorID:    session.MentorID,
		MenteeID:    session.MenteeID,
		SpaceID:     session.SpaceID,
		Topic:       session.Topic,
		Status:      nullStringToPtr(session.Status),
		ScheduledAt: session.ScheduledAt,
		Duration:    session.Duration,
		MenteeNotes: nullStringToPtr(session.MenteeNotes),
		MentorNotes: nullStringToPtr(session.MentorNotes),
		MeetingLink: nullStringToPtr(session.MeetingLink),
		Rating:      nullInt32ToPtr(session.Rating),
		Review:      nullStringToPtr(session.Review),
		CreatedAt:   nullTimeToPtr(session.CreatedAt),
	}
}

func tutoringSessionToResponse(session db.TutoringSession) *TutoringSessionResponse {
	return &TutoringSessionResponse{
		ID:           session.ID,
		TutorID:      session.TutorID,
		StudentID:    session.StudentID,
		SpaceID:      session.SpaceID,
		Subject:      session.Subject,
		Status:       nullStringToPtr(session.Status),
		ScheduledAt:  session.ScheduledAt,
		Duration:     session.Duration,
		HourlyRate:   nullStringToPtr(session.HourlyRate),
		TotalAmount:  nullStringToPtr(session.TotalAmount),
		StudentNotes: nullStringToPtr(session.StudentNotes),
		TutorNotes:   nullStringToPtr(session.TutorNotes),
		MeetingLink:  nullStringToPtr(session.MeetingLink),
		Rating:       nullInt32ToPtr(session.Rating),
		Review:       nullStringToPtr(session.Review),
		CreatedAt:    nullTimeToPtr(session.CreatedAt),
	}
}

func mentorApplicationToResponse(app db.MentorApplication) *MentorApplicationResponse {
	return &MentorApplicationResponse{
		ID:             app.ID,
		UserID:         app.ApplicantID,
		SpaceID:        app.SpaceID,
		Industry:       app.Industry,
		Company:        nullStringToPtr(app.Company),
		Position:       nullStringToPtr(app.Position),
		Experience:     app.Experience,
		Specialties:    app.Specialties,
		Motivation:     nullStringToStringValue(app.Motivation),
		Status:         nullStringToPtr(app.Status),
		ReviewedBy:     nullUUIDToPtr(app.ReviewedBy),
		ReviewComments: nullStringToPtr(app.ReviewerNotes),
		CreatedAt:      nullTimeToPtr(app.SubmittedAt),
		ReviewedAt:     nullTimeToPtr(app.ReviewedAt),
	}
}

func tutorApplicationToResponse(app db.TutorApplication) *TutorApplicationResponse {
	return &TutorApplicationResponse{
		ID:             app.ID,
		UserID:         app.ApplicantID,
		SpaceID:        app.SpaceID,
		Subjects:       app.Subjects,
		Experience:     nullStringToStringValue(app.Experience),
		Qualifications: nullStringToStringValue(app.Qualifications),
		Motivation:     nullStringToStringValue(app.Motivation),
		Status:         nullStringToPtr(app.Status),
		ReviewedBy:     nullUUIDToPtr(app.ReviewedBy),
		ReviewComments: nullStringToPtr(app.ReviewerNotes),
		CreatedAt:      nullTimeToPtr(app.SubmittedAt),
		ReviewedAt:     nullTimeToPtr(app.ReviewedAt),
	}
}

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullRawMessage(r []byte) pqtype.NullRawMessage {
	if r == nil || len(r) == 0 {
		return pqtype.NullRawMessage{Valid: false}
	}
	return pqtype.NullRawMessage{RawMessage: r, Valid: true}
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func nullInt32ToPtr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func nullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

func nullBoolToPtr(nb sql.NullBool) *bool {
	if !nb.Valid {
		return nil
	}
	return &nb.Bool
}

func nullTimeToPtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func nullUUIDToPtr(nu uuid.NullUUID) *uuid.UUID {
	if !nu.Valid {
		return nil
	}
	return &nu.UUID
}

func nullRawMessageToPtr(nr pqtype.NullRawMessage) *pqtype.NullRawMessage {
	if !nr.Valid {
		return nil
	}
	return &nr
}

func nullStringToStringValue(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// GetMentorProfileByUserID gets a mentor profile by user ID
func (s *Service) GetMentorProfileByUserID(ctx context.Context, userID uuid.UUID) (*MentorProfileResponse, error) {
	profile, err := s.store.GetMentorProfile(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mentor profile not found")
		}
		return nil, fmt.Errorf("failed to get mentor profile: %w", err)
	}

	return &MentorProfileResponse{
		ID:            profile.ID,
		UserID:        profile.UserID,
		SpaceID:       profile.SpaceID,
		Industry:      profile.Industry,
		Company:       nullStringToPtr(profile.Company),
		Position:      nullStringToPtr(profile.Position),
		Experience:    profile.Experience,
		Specialties:   profile.Specialties,
		Rating:        nullFloat64ToPtr(profile.Rating),
		ReviewCount:   nullInt32ToPtr(profile.ReviewCount),
		TotalSessions: nullInt32ToPtr(profile.TotalSessions),
		Availability:  nullRawMessageToPtr(profile.Availability),
		Description:   nullStringToPtr(profile.Description),
		Verified:      nullBoolToPtr(profile.Verified),
		IsAvailable:   nullBoolToPtr(profile.IsAvailable),
		CreatedAt:     nullTimeToPtr(profile.CreatedAt),
		UpdatedAt:     nullTimeToPtr(profile.UpdatedAt),
	}, nil
}

// GetTutorProfileByUserID gets a tutor profile by user ID
func (s *Service) GetTutorProfileByUserID(ctx context.Context, userID uuid.UUID) (*TutorProfileResponse, error) {
	profile, err := s.store.GetTutorProfile(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tutor profile not found")
		}
		return nil, fmt.Errorf("failed to get tutor profile: %w", err)
	}

	return &TutorProfileResponse{
		ID:             profile.ID,
		UserID:         profile.UserID,
		SpaceID:        profile.SpaceID,
		Subjects:       profile.Subjects,
		HourlyRate:     nullStringToPtr(profile.HourlyRate),
		Rating:         nullFloat64ToPtr(profile.Rating),
		ReviewCount:    nullInt32ToPtr(profile.ReviewCount),
		TotalSessions:  nullInt32ToPtr(profile.TotalSessions),
		Description:    nullStringToPtr(profile.Description),
		Availability:   nullRawMessageToPtr(profile.Availability),
		Experience:     nullStringToPtr(profile.Experience),
		Qualifications: nullStringToPtr(profile.Qualifications),
		Verified:       nullBoolToPtr(profile.Verified),
		IsAvailable:    nullBoolToPtr(profile.IsAvailable),
		CreatedAt:      nullTimeToPtr(profile.CreatedAt),
		UpdatedAt:      nullTimeToPtr(profile.UpdatedAt),
	}, nil
}

// GetMentorApplicationByUserID gets a mentor application by user ID
func (s *Service) GetMentorApplicationByUserID(ctx context.Context, userID uuid.UUID, spaceId uuid.UUID) (*MentorApplicationResponse, error) {
	application, err := s.store.GetUserMentorApplicationStatus(ctx, db.GetUserMentorApplicationStatusParams{
		ApplicantID: userID,
		SpaceID:     spaceId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no mentor application found")
		}
		return nil, fmt.Errorf("failed to get mentor application: %w", err)
	}

	resp := &MentorApplicationResponse{
		Status: &application.String,
	}

	return resp, nil
}

// GetTutorApplicationByUserID gets a tutor application by user ID
func (s *Service) GetTutorApplicationByUserID(ctx context.Context, userID uuid.UUID, spaceID uuid.UUID) (*TutorApplicationResponse, error) {
	application, err := s.store.GetUserTutorApplicationStatus(ctx, db.GetUserTutorApplicationStatusParams{
		ApplicantID: userID,
		SpaceID:     spaceID,
	})

	

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tutor application found")
		}
		return nil, fmt.Errorf("failed to get tutor application: %w", err)
	}

	resp := &TutorApplicationResponse{
		Status: &application.String,
	}

	return resp, nil
}

// ===== Recommendation Operations =====

// GetRecommendedTutors returns recommended tutors for a user based on matching criteria
func (s *Service) GetRecommendedTutors(ctx context.Context, spaceID, userID uuid.UUID, limit int32) ([]TutorSearchResponse, error) {
	if limit == 0 {
		limit = 5
	}

	tutors, err := s.store.GetRecommendedTutors(ctx, db.GetRecommendedTutorsParams{
		SpaceID: spaceID,
		UserID:  userID,
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended tutors: %w", err)
	}

	result := make([]TutorSearchResponse, len(tutors))
	for i, tutor := range tutors {
		result[i] = TutorSearchResponse{
			ID:             tutor.ID,
			UserID:         tutor.UserID,
			Username:       tutor.Username,
			FullName:       tutor.FullName,
			Avatar:         nullStringToPtr(tutor.Avatar),
			Subjects:       tutor.Subjects,
			HourlyRate:     nullStringToPtr(tutor.HourlyRate),
			Rating:         nullFloat64ToPtr(tutor.Rating),
			ReviewCount:    nullInt32ToPtr(tutor.ReviewCount),
			Experience:     nullStringToPtr(tutor.Experience),
			Qualifications: nullStringToPtr(tutor.Qualifications),
			IsAvailable:    nullBoolToPtr(tutor.IsAvailable),
		}
	}

	return result, nil
}

// GetRecommendedMentors returns recommended mentors for a user based on matching criteria
func (s *Service) GetRecommendedMentors(ctx context.Context, spaceID, userID uuid.UUID, limit int32) ([]MentorSearchResponse, error) {
	if limit == 0 {
		limit = 5
	}

	mentors, err := s.store.GetRecommendedMentors(ctx, db.GetRecommendedMentorsParams{
		SpaceID: spaceID,
		UserID:  userID,
		Limit:   limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recommended mentors: %w", err)
	}

	result := make([]MentorSearchResponse, len(mentors))
	for i, mentor := range mentors {
		result[i] = MentorSearchResponse{
			ID:          mentor.ID,
			UserID:      mentor.UserID,
			Username:    mentor.Username,
			FullName:    mentor.FullName,
			Avatar:      nullStringToPtr(mentor.Avatar),
			Industry:    mentor.Industry,
			Company:     nullStringToPtr(mentor.Company),
			Position:    nullStringToPtr(mentor.Position),
			Experience:  mentor.Experience,
			Specialties: mentor.Specialties,
			Rating:      nullFloat64ToPtr(mentor.Rating),
			ReviewCount: nullInt32ToPtr(mentor.ReviewCount),
			IsAvailable: nullBoolToPtr(mentor.IsAvailable),
		}
	}

	return result, nil
}
