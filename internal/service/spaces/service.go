package spaces

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// stringToNullString converts *string to sql.NullString
func stringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

// Service handles space business logic
type Service struct {
	store db.Store
}

// NewService creates a new space service
func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}

// CreateSpace creates a new space with validation
func (s *Service) CreateSpace(ctx context.Context, req CreateSpaceRequest) (*SpaceResponse, error) {
	// Validate input
	if err := ValidateName(req.Name); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}
	if err := ValidateSlug(req.Slug); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}

	// Check if space with same slug already exists
	existingSpace, err := s.store.GetSpaceBySlug(ctx, req.Slug)
	if err == nil && existingSpace.ID != uuid.Nil {
		return nil, fmt.Errorf("%w: space with this slug already exists", util.ErrConflict)
	}

	// Create space
	space, err := s.store.CreateSpace(ctx, db.CreateSpaceParams{
		Name:         req.Name,
		Slug:         req.Slug,
		Description:  stringToNullString(req.Description),
		Type:         stringToNullString(req.Type),
		Logo:         stringToNullString(req.Logo),
		Location:     stringToNullString(req.Location),
		Website:      stringToNullString(req.Website),
		ContactEmail: stringToNullString(req.ContactEmail),
		PhoneNumber:  stringToNullString(req.PhoneNumber),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create space: %w", err)
	}

	return ToSpaceResponse(&space), nil
}

// GetSpace retrieves a space by ID
func (s *Service) GetSpace(ctx context.Context, id uuid.UUID) (*SpaceResponse, error) {
	space, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	return ToSpaceResponse(&space), nil
}

// GetSpaceBySlug retrieves a space by slug
func (s *Service) GetSpaceBySlug(ctx context.Context, slug string) (*SpaceResponse, error) {
	space, err := s.store.GetSpaceBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	return ToSpaceResponse(&space), nil
}

// ListSpaces retrieves all spaces with pagination
func (s *Service) ListSpaces(ctx context.Context, page, limit int32) (*PaginatedSpacesResponse, error) {
	// Validate and set defaults
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	spaces, err := s.store.ListSpaces(ctx, db.ListSpacesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list spaces: %w", err)
	}

	// Convert to response format
	var spaceResponses []*SpaceResponse
	for i := range spaces {
		spaceResponses = append(spaceResponses, ToSpaceResponse(&spaces[i]))
	}

	return &PaginatedSpacesResponse{
		Spaces: spaceResponses,
		Total:  int64(len(spaceResponses)), // In production, you'd want a count query
		Page:   page,
		Limit:  limit,
	}, nil
}

// UpdateSpace updates a space
func (s *Service) UpdateSpace(ctx context.Context, id uuid.UUID, req UpdateSpaceRequest) (*SpaceResponse, error) {
	// Verify space exists
	_, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	// If slug is being updated, check uniqueness
	if req.Slug != nil {
		existingSpace, err := s.store.GetSpaceBySlug(ctx, *req.Slug)
		if err == nil && existingSpace.ID != id {
			return nil, fmt.Errorf("%w: space with this slug already exists", util.ErrConflict)
		}
	}

	// Handle Settings field (it's already a pointer to NullRawMessage)
	var settings pqtype.NullRawMessage
	if req.Settings != nil {
		settings = *req.Settings
	}

	// Update space
	space, err := s.store.UpdateSpace(ctx, db.UpdateSpaceParams{
		ID:           id,
		Name:         stringToNullString(req.Name),
		Slug:         stringToNullString(req.Slug),
		Description:  stringToNullString(req.Description),
		Type:         stringToNullString(req.Type),
		Logo:         stringToNullString(req.Logo),
		CoverImage:   stringToNullString(req.CoverImage),
		Location:     stringToNullString(req.Location),
		Website:      stringToNullString(req.Website),
		ContactEmail: stringToNullString(req.ContactEmail),
		PhoneNumber:  stringToNullString(req.PhoneNumber),
		Status:       stringToNullString(req.Status),
		Settings:     settings,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update space: %w", err)
	}

	return ToSpaceResponse(&space), nil
}

// DeleteSpace deletes a space
func (s *Service) DeleteSpace(ctx context.Context, id uuid.UUID) error {
	// Verify space exists
	_, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	// Delete space
	if err := s.store.DeleteSpace(ctx, id); err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}

	return nil
}
