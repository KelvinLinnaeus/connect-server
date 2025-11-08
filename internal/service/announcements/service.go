package announcements

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
}

func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

// CreateAnnouncement creates a new announcement
func (s *Service) CreateAnnouncement(ctx context.Context, req CreateAnnouncementRequest) (*AnnouncementResponse, error) {
	params := db.CreateAnnouncementParams{
		SpaceID:        req.SpaceID,
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		TargetAudience: req.TargetAudience,
		Priority:       sqlNullString(req.Priority),
		AuthorID:       uuid.NullUUID{UUID: req.AuthorID, Valid: true},
		ScheduledFor:   sqlNullTime(req.ScheduledFor),
		ExpiresAt:      sqlNullTime(req.ExpiresAt),
		Attachments:    sqlNullRawMessage(req.Attachments),
		IsPinned:       sqlNullBool(req.IsPinned),
	}

	announcement, err := s.store.CreateAnnouncement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	return &AnnouncementResponse{
		ID:             announcement.ID,
		SpaceID:        announcement.SpaceID,
		Title:          announcement.Title,
		Content:        announcement.Content,
		Type:           announcement.Type,
		TargetAudience: announcement.TargetAudience,
		Priority:       nullStringToPtr(announcement.Priority),
		Status:         nullStringToPtr(announcement.Status),
		AuthorID:       nullUUIDToPtr(announcement.AuthorID),
		ScheduledFor:   nullTimeToPtr(announcement.ScheduledFor),
		ExpiresAt:      nullTimeToPtr(announcement.ExpiresAt),
		Attachments:    nullRawMessageToPtr(announcement.Attachments),
		IsPinned:       nullBoolToPtr(announcement.IsPinned),
		CreatedAt:      nullTimeToPtr(announcement.CreatedAt),
		UpdatedAt:      nullTimeToPtr(announcement.UpdatedAt),
	}, nil
}

// GetAnnouncementByID retrieves an announcement by ID
func (s *Service) GetAnnouncementByID(ctx context.Context, announcementID uuid.UUID) (*AnnouncementDetailResponse, error) {
	announcement, err := s.store.GetAnnouncementByID(ctx, announcementID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("announcement not found")
		}
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return &AnnouncementDetailResponse{
		ID:             announcement.ID,
		SpaceID:        announcement.SpaceID,
		Title:          announcement.Title,
		Content:        announcement.Content,
		Type:           announcement.Type,
		TargetAudience: announcement.TargetAudience,
		Priority:       nullStringToPtr(announcement.Priority),
		Status:         nullStringToPtr(announcement.Status),
		ScheduledFor:   nullTimeToPtr(announcement.ScheduledFor),
		ExpiresAt:      nullTimeToPtr(announcement.ExpiresAt),
		Attachments:    nullRawMessageToPtr(announcement.Attachments),
		IsPinned:       nullBoolToPtr(announcement.IsPinned),
		CreatedAt:      nullTimeToPtr(announcement.CreatedAt),
		UpdatedAt:      nullTimeToPtr(announcement.UpdatedAt),
		AuthorID:       nullUUIDToPtr(announcement.AuthorID),
		AuthorUsername: announcement.AuthorUsername,
		AuthorFullName: announcement.AuthorFullName,
		AuthorAvatar:   nullStringToPtr(announcement.AuthorAvatar),
	}, nil
}

// ListAnnouncements retrieves announcements with filtering and pagination
func (s *Service) ListAnnouncements(ctx context.Context, params ListAnnouncementsParams) ([]AnnouncementListResponse, error) {
	// Default pagination
	if params.Limit == 0 {
		params.Limit = 20
	}
	offset := (params.Page - 1) * params.Limit

	announcements, err := s.store.ListAnnouncements(ctx, db.ListAnnouncementsParams{
		SpaceID:        params.SpaceID,
		TargetAudience: params.TargetAudience,
		Limit:          params.Limit,
		Offset:         offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list announcements: %w", err)
	}

	result := make([]AnnouncementListResponse, len(announcements))
	for i, announcement := range announcements {
		result[i] = AnnouncementListResponse{
			ID:             announcement.ID,
			SpaceID:        announcement.SpaceID,
			Title:          announcement.Title,
			Content:        announcement.Content,
			Type:           announcement.Type,
			TargetAudience: announcement.TargetAudience,
			Priority:       nullStringToPtr(announcement.Priority),
			IsPinned:       nullBoolToPtr(announcement.IsPinned),
			CreatedAt:      nullTimeToPtr(announcement.CreatedAt),
			AuthorUsername: announcement.AuthorUsername,
			AuthorFullName: announcement.AuthorFullName,
		}
	}

	return result, nil
}

// UpdateAnnouncement updates announcement details
func (s *Service) UpdateAnnouncement(ctx context.Context, announcementID uuid.UUID, req UpdateAnnouncementRequest) (*AnnouncementResponse, error) {
	params := db.UpdateAnnouncementParams{
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		TargetAudience: req.TargetAudience,
		Priority:       sqlNullString(req.Priority),
		ScheduledFor:   sqlNullTime(req.ScheduledFor),
		ExpiresAt:      sqlNullTime(req.ExpiresAt),
		Attachments:    sqlNullRawMessage(req.Attachments),
		IsPinned:       sqlNullBool(req.IsPinned),
		ID:             announcementID,
	}

	announcement, err := s.store.UpdateAnnouncement(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement: %w", err)
	}

	return &AnnouncementResponse{
		ID:             announcement.ID,
		SpaceID:        announcement.SpaceID,
		Title:          announcement.Title,
		Content:        announcement.Content,
		Type:           announcement.Type,
		TargetAudience: announcement.TargetAudience,
		Priority:       nullStringToPtr(announcement.Priority),
		Status:         nullStringToPtr(announcement.Status),
		AuthorID:       nullUUIDToPtr(announcement.AuthorID),
		ScheduledFor:   nullTimeToPtr(announcement.ScheduledFor),
		ExpiresAt:      nullTimeToPtr(announcement.ExpiresAt),
		Attachments:    nullRawMessageToPtr(announcement.Attachments),
		IsPinned:       nullBoolToPtr(announcement.IsPinned),
		CreatedAt:      nullTimeToPtr(announcement.CreatedAt),
		UpdatedAt:      nullTimeToPtr(announcement.UpdatedAt),
	}, nil
}

// UpdateAnnouncementStatus updates the status of an announcement
func (s *Service) UpdateAnnouncementStatus(ctx context.Context, announcementID uuid.UUID, status string) (*AnnouncementResponse, error) {
	announcement, err := s.store.UpdateAnnouncementStatus(ctx, db.UpdateAnnouncementStatusParams{
		Status: sql.NullString{String: status, Valid: true},
		ID:     announcementID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement status: %w", err)
	}

	return &AnnouncementResponse{
		ID:             announcement.ID,
		SpaceID:        announcement.SpaceID,
		Title:          announcement.Title,
		Content:        announcement.Content,
		Type:           announcement.Type,
		TargetAudience: announcement.TargetAudience,
		Priority:       nullStringToPtr(announcement.Priority),
		Status:         nullStringToPtr(announcement.Status),
		AuthorID:       nullUUIDToPtr(announcement.AuthorID),
		ScheduledFor:   nullTimeToPtr(announcement.ScheduledFor),
		ExpiresAt:      nullTimeToPtr(announcement.ExpiresAt),
		Attachments:    nullRawMessageToPtr(announcement.Attachments),
		IsPinned:       nullBoolToPtr(announcement.IsPinned),
		CreatedAt:      nullTimeToPtr(announcement.CreatedAt),
		UpdatedAt:      nullTimeToPtr(announcement.UpdatedAt),
	}, nil
}

// Helper functions

func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
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

func sqlNullRawMessage(r *pqtype.NullRawMessage) pqtype.NullRawMessage {
	if r == nil {
		return pqtype.NullRawMessage{Valid: false}
	}
	return *r
}

func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
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
