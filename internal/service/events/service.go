package events

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/google/uuid"
)

type Service struct {
	store db.Store
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}


func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest) (*EventResponse, error) {
	
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	
	params := db.CreateEventParams{
		SpaceID:              req.SpaceID,
		Title:                req.Title,
		Description:          sqlNullString(req.Description),
		Category:             req.Category,
		Location:             sqlNullString(req.Location),
		VenueDetails:         sqlNullString(req.VenueDetails),
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		Timezone:             sqlNullString(req.Timezone),
		Organizer:            uuid.NullUUID{UUID: req.OrganizerID, Valid: true},
		Tags:                 req.Tags,
		ImageUrl:             sqlNullString(req.ImageURL),
		MaxAttendees:         sqlNullInt32(req.MaxAttendees),
		RegistrationRequired: sqlNullBool(req.RegistrationRequired),
		RegistrationDeadline: sqlNullTime(req.RegistrationDeadline),
		IsPublic:             sqlNullBool(req.IsPublic),
	}

	event, err := s.store.CreateEvent(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return &EventResponse{
		ID:                   event.ID,
		SpaceID:              event.SpaceID,
		Title:                event.Title,
		Description:          nullStringToPtr(event.Description),
		Category:             event.Category,
		Location:             nullStringToPtr(event.Location),
		VenueDetails:         nullStringToPtr(event.VenueDetails),
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		Timezone:             nullStringToPtr(event.Timezone),
		OrganizerID:          nullUUIDToPtr(event.Organizer),
		Tags:                 event.Tags,
		ImageURL:             nullStringToPtr(event.ImageUrl),
		MaxAttendees:         nullInt32ToPtr(event.MaxAttendees),
		CurrentAttendees:     nullInt32ToPtr(event.CurrentAttendees),
		RegistrationRequired: nullBoolToPtr(event.RegistrationRequired),
		RegistrationDeadline: nullTimeToPtr(event.RegistrationDeadline),
		Status:               nullStringToPtr(event.Status),
		IsPublic:             nullBoolToPtr(event.IsPublic),
		CreatedAt:            nullTimeToPtr(event.CreatedAt),
		UpdatedAt:            nullTimeToPtr(event.UpdatedAt),
	}, nil
}


func (s *Service) GetEventByID(ctx context.Context, eventID, userID uuid.UUID) (*EventDetailResponse, error) {
	event, err := s.store.GetEventByID(ctx, db.GetEventByIDParams{
		UserID: userID,
		ID:     eventID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("event not found")
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &EventDetailResponse{
		ID:                    event.ID,
		SpaceID:               event.SpaceID,
		Title:                 event.Title,
		Description:           nullStringToPtr(event.Description),
		Category:              event.Category,
		Location:              nullStringToPtr(event.Location),
		VenueDetails:          nullStringToPtr(event.VenueDetails),
		StartDate:             event.StartDate,
		EndDate:               event.EndDate,
		Timezone:              nullStringToPtr(event.Timezone),
		Tags:                  event.Tags,
		ImageURL:              nullStringToPtr(event.ImageUrl),
		MaxAttendees:          nullInt32ToPtr(event.MaxAttendees),
		RegistrationRequired:  nullBoolToPtr(event.RegistrationRequired),
		RegistrationDeadline:  nullTimeToPtr(event.RegistrationDeadline),
		Status:                nullStringToPtr(event.Status),
		IsPublic:              nullBoolToPtr(event.IsPublic),
		CreatedAt:             nullTimeToPtr(event.CreatedAt),
		UpdatedAt:             nullTimeToPtr(event.UpdatedAt),
		OrganizerID:           nullUUIDToPtr(event.Organizer),
		OrganizerUsername:     event.OrganizerUsername,
		OrganizerFullName:     event.OrganizerFullName,
		OrganizerAvatar:       nullStringToPtr(event.OrganizerAvatar),
		CurrentAttendeesCount: event.CurrentAttendeesCount,
		IsRegistered:          event.IsRegistered,
		UserAttendanceStatus:  nullStringToPtr(event.UserAttendanceStatus),
	}, nil
}


func (s *Service) ListEvents(ctx context.Context, params ListEventsParams) ([]EventListResponse, error) {
	
	if params.Limit == 0 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	
	startDate := time.Now()
	if params.StartDate != nil {
		startDate = *params.StartDate
	}

	
	category := ""
	if params.Category != nil {
		category = *params.Category
	}

	
	var sortBy interface{} = "recent"
	if params.Sort != nil {
		sortBy = *params.Sort
	}

	events, err := s.store.ListEvents(ctx, db.ListEventsParams{
		UserID:    params.UserID,
		SpaceID:   params.SpaceID,
		StartDate: startDate,
		Category:  category,
		Column5:   sortBy,
		Limit:     params.Limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	result := make([]EventListResponse, len(events))
	for i, event := range events {
		result[i] = EventListResponse{
			ID:                    event.ID,
			SpaceID:               event.SpaceID,
			Title:                 event.Title,
			Description:           nullStringToPtr(event.Description),
			Category:              event.Category,
			Location:              nullStringToPtr(event.Location),
			StartDate:             event.StartDate,
			EndDate:               event.EndDate,
			Tags:                  event.Tags,
			ImageURL:              nullStringToPtr(event.ImageUrl),
			MaxAttendees:          nullInt32ToPtr(event.MaxAttendees),
			RegistrationRequired:  nullBoolToPtr(event.RegistrationRequired),
			IsPublic:              nullBoolToPtr(event.IsPublic),
			OrganizerUsername:     event.OrganizerUsername,
			OrganizerFullName:     event.OrganizerFullName,
			CurrentAttendeesCount: event.CurrentAttendeesCount,
			IsRegistered:          event.IsRegistered,
		}
	}

	return result, nil
}


func (s *Service) GetUpcomingEvents(ctx context.Context, spaceID uuid.UUID) ([]EventListResponse, error) {
	events, err := s.store.GetUpcomingEvents(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	result := make([]EventListResponse, len(events))
	for i, event := range events {
		result[i] = EventListResponse{
			ID:                    event.ID,
			SpaceID:               event.SpaceID,
			Title:                 event.Title,
			Description:           nullStringToPtr(event.Description),
			Category:              event.Category,
			Location:              nullStringToPtr(event.Location),
			StartDate:             event.StartDate,
			EndDate:               event.EndDate,
			Tags:                  event.Tags,
			ImageURL:              nullStringToPtr(event.ImageUrl),
			MaxAttendees:          nullInt32ToPtr(event.MaxAttendees),
			RegistrationRequired:  nullBoolToPtr(event.RegistrationRequired),
			IsPublic:              nullBoolToPtr(event.IsPublic),
			OrganizerUsername:     event.OrganizerUsername,
			OrganizerFullName:     event.OrganizerFullName,
			CurrentAttendeesCount: event.CurrentAttendeesCount,
		}
	}

	return result, nil
}


func (s *Service) GetUserEvents(ctx context.Context, userID, spaceID uuid.UUID, page, limit int32) ([]UserEventResponse, error) {
	if limit == 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	events, err := s.store.GetUserEvents(ctx, db.GetUserEventsParams{
		UserID:  userID,
		SpaceID: spaceID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user events: %w", err)
	}

	result := make([]UserEventResponse, len(events))
	for i, event := range events {
		result[i] = UserEventResponse{
			ID:                event.ID,
			SpaceID:           event.SpaceID,
			Title:             event.Title,
			Description:       nullStringToPtr(event.Description),
			Category:          event.Category,
			Location:          nullStringToPtr(event.Location),
			StartDate:         event.StartDate,
			EndDate:           event.EndDate,
			Tags:              event.Tags,
			ImageURL:          nullStringToPtr(event.ImageUrl),
			OrganizerUsername: event.OrganizerUsername,
			OrganizerFullName: event.OrganizerFullName,
			AttendanceStatus:  nullStringToPtr(event.AttendanceStatus),
		}
	}

	return result, nil
}


func (s *Service) SearchEvents(ctx context.Context, params SearchEventsParams) ([]EventListResponse, error) {
	
	searchPattern := "%" + params.Query + "%"

	events, err := s.store.SearchEvents(ctx, db.SearchEventsParams{
		UserID:  params.UserID,
		SpaceID: params.SpaceID,
		Title:   searchPattern,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search events: %w", err)
	}

	result := make([]EventListResponse, len(events))
	for i, event := range events {
		result[i] = EventListResponse{
			ID:                    event.ID,
			SpaceID:               event.SpaceID,
			Title:                 event.Title,
			Description:           nullStringToPtr(event.Description),
			Category:              event.Category,
			Location:              nullStringToPtr(event.Location),
			StartDate:             event.StartDate,
			EndDate:               event.EndDate,
			Tags:                  event.Tags,
			ImageURL:              nullStringToPtr(event.ImageUrl),
			MaxAttendees:          nullInt32ToPtr(event.MaxAttendees),
			RegistrationRequired:  nullBoolToPtr(event.RegistrationRequired),
			IsPublic:              nullBoolToPtr(event.IsPublic),
			OrganizerUsername:     event.OrganizerUsername,
			OrganizerFullName:     event.OrganizerFullName,
			CurrentAttendeesCount: event.CurrentAttendeesCount,
			IsRegistered:          event.IsRegistered,
		}
	}

	return result, nil
}


func (s *Service) GetEventCategories(ctx context.Context, spaceID uuid.UUID) ([]string, error) {
	categories, err := s.store.GetEventCategories(ctx, spaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event categories: %w", err)
	}
	return categories, nil
}


func (s *Service) RegisterForEvent(ctx context.Context, eventID, userID uuid.UUID) (*RegistrationResponse, error) {
	attendee, err := s.store.RegisterForEvent(ctx, db.RegisterForEventParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to register for event: %w", err)
	}

	return &RegistrationResponse{
		EventID:      attendee.EventID,
		UserID:       attendee.UserID,
		Status:       nullStringToPtr(attendee.Status),
		RegisteredAt: nullTimeToPtr(attendee.RegisteredAt),
		Message:      "Successfully registered for event",
	}, nil
}


func (s *Service) UnregisterFromEvent(ctx context.Context, eventID, userID uuid.UUID) error {
	err := s.store.UnregisterFromEvent(ctx, db.UnregisterFromEventParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to unregister from event: %w", err)
	}
	return nil
}


func (s *Service) GetEventAttendees(ctx context.Context, eventID uuid.UUID) ([]AttendeeResponse, error) {
	attendees, err := s.store.GetEventAttendees(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get event attendees: %w", err)
	}

	result := make([]AttendeeResponse, len(attendees))
	for i, attendee := range attendees {
		result[i] = AttendeeResponse{
			ID:           attendee.ID,
			EventID:      attendee.EventID,
			UserID:       attendee.UserID,
			Username:     attendee.Username,
			FullName:     attendee.FullName,
			Avatar:       nullStringToPtr(attendee.Avatar),
			Department:   nullStringToPtr(attendee.Department),
			Level:        nullStringToPtr(attendee.Level),
			Status:       nullStringToPtr(attendee.Status),
			Role:         nullStringToPtr(attendee.Role),
			RegisteredAt: nullTimeToPtr(attendee.RegisteredAt),
			AttendedAt:   nullTimeToPtr(attendee.AttendedAt),
			Notes:        nullStringToPtr(attendee.Notes),
		}
	}

	return result, nil
}


func (s *Service) MarkEventAttendance(ctx context.Context, eventID, userID uuid.UUID) error {
	err := s.store.MarkEventAttendance(ctx, db.MarkEventAttendanceParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark attendance: %w", err)
	}
	return nil
}


func (s *Service) AddEventCoOrganizer(ctx context.Context, eventID, userID uuid.UUID) (*CoOrganizerResponse, error) {
	coOrganizer, err := s.store.AddEventCoOrganizer(ctx, db.AddEventCoOrganizerParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add co-organizer: %w", err)
	}

	
	
	return &CoOrganizerResponse{
		ID:           coOrganizer.ID,
		EventID:      coOrganizer.EventID,
		UserID:       coOrganizer.UserID,
		RegisteredAt: nullTimeToPtr(coOrganizer.RegisteredAt),
	}, nil
}


func (s *Service) GetEventCoOrganizers(ctx context.Context, eventID uuid.UUID) ([]CoOrganizerResponse, error) {
	coOrganizers, err := s.store.GetEventCoOrganizers(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get co-organizers: %w", err)
	}

	result := make([]CoOrganizerResponse, len(coOrganizers))
	for i, coOrg := range coOrganizers {
		result[i] = CoOrganizerResponse{
			ID:           coOrg.ID,
			EventID:      coOrg.EventID,
			UserID:       coOrg.UserID,
			Username:     coOrg.Username,
			FullName:     coOrg.FullName,
			Avatar:       nullStringToPtr(coOrg.Avatar),
			RegisteredAt: nullTimeToPtr(coOrg.RegisteredAt),
		}
	}

	return result, nil
}


func (s *Service) RemoveEventCoOrganizer(ctx context.Context, eventID, userID uuid.UUID) error {
	err := s.store.RemoveEventCoOrganizer(ctx, db.RemoveEventCoOrganizerParams{
		EventID: eventID,
		UserID:  userID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove co-organizer: %w", err)
	}
	return nil
}


func (s *Service) UpdateEvent(ctx context.Context, eventID uuid.UUID, req UpdateEventRequest) (*EventResponse, error) {
	
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	params := db.UpdateEventParams{
		Title:                req.Title,
		Description:          sqlNullString(req.Description),
		Category:             req.Category,
		Location:             sqlNullString(req.Location),
		VenueDetails:         sqlNullString(req.VenueDetails),
		StartDate:            req.StartDate,
		EndDate:              req.EndDate,
		Timezone:             sqlNullString(req.Timezone),
		Tags:                 req.Tags,
		ImageUrl:             sqlNullString(req.ImageURL),
		MaxAttendees:         sqlNullInt32(req.MaxAttendees),
		RegistrationRequired: sqlNullBool(req.RegistrationRequired),
		RegistrationDeadline: sqlNullTime(req.RegistrationDeadline),
		IsPublic:             sqlNullBool(req.IsPublic),
		ID:                   eventID,
	}

	event, err := s.store.UpdateEvent(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	return &EventResponse{
		ID:                   event.ID,
		SpaceID:              event.SpaceID,
		Title:                event.Title,
		Description:          nullStringToPtr(event.Description),
		Category:             event.Category,
		Location:             nullStringToPtr(event.Location),
		VenueDetails:         nullStringToPtr(event.VenueDetails),
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		Timezone:             nullStringToPtr(event.Timezone),
		OrganizerID:          nullUUIDToPtr(event.Organizer),
		Tags:                 event.Tags,
		ImageURL:             nullStringToPtr(event.ImageUrl),
		MaxAttendees:         nullInt32ToPtr(event.MaxAttendees),
		CurrentAttendees:     nullInt32ToPtr(event.CurrentAttendees),
		RegistrationRequired: nullBoolToPtr(event.RegistrationRequired),
		RegistrationDeadline: nullTimeToPtr(event.RegistrationDeadline),
		Status:               nullStringToPtr(event.Status),
		IsPublic:             nullBoolToPtr(event.IsPublic),
		CreatedAt:            nullTimeToPtr(event.CreatedAt),
		UpdatedAt:            nullTimeToPtr(event.UpdatedAt),
	}, nil
}


func (s *Service) UpdateEventStatus(ctx context.Context, eventID uuid.UUID, status string) (*EventResponse, error) {
	event, err := s.store.UpdateEventStatus(ctx, db.UpdateEventStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     eventID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update event status: %w", err)
	}

	return &EventResponse{
		ID:                   event.ID,
		SpaceID:              event.SpaceID,
		Title:                event.Title,
		Description:          nullStringToPtr(event.Description),
		Category:             event.Category,
		Location:             nullStringToPtr(event.Location),
		VenueDetails:         nullStringToPtr(event.VenueDetails),
		StartDate:            event.StartDate,
		EndDate:              event.EndDate,
		Timezone:             nullStringToPtr(event.Timezone),
		OrganizerID:          nullUUIDToPtr(event.Organizer),
		Tags:                 event.Tags,
		ImageURL:             nullStringToPtr(event.ImageUrl),
		MaxAttendees:         nullInt32ToPtr(event.MaxAttendees),
		CurrentAttendees:     nullInt32ToPtr(event.CurrentAttendees),
		RegistrationRequired: nullBoolToPtr(event.RegistrationRequired),
		RegistrationDeadline: nullTimeToPtr(event.RegistrationDeadline),
		Status:               nullStringToPtr(event.Status),
		IsPublic:             nullBoolToPtr(event.IsPublic),
		CreatedAt:            nullTimeToPtr(event.CreatedAt),
		UpdatedAt:            nullTimeToPtr(event.UpdatedAt),
	}, nil
}



func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func sqlNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func sqlNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func sqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
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
