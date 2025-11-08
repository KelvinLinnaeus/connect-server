package spaces

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)


func stringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}


type Service struct {
	store db.Store
}


func NewService(store db.Store) *Service {
	return &Service{
		store: store,
	}
}


func (s *Service) CreateSpace(ctx context.Context, req CreateSpaceRequest) (*SpaceResponse, error) {
	
	if err := ValidateName(req.Name); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}
	if err := ValidateSlug(req.Slug); err != nil {
		return nil, fmt.Errorf("%w: %v", util.ErrBadRequest, err)
	}

	
	existingSpace, err := s.store.GetSpaceBySlug(ctx, req.Slug)
	if err == nil && existingSpace.ID != uuid.Nil {
		return nil, fmt.Errorf("%w: space with this slug already exists", util.ErrConflict)
	}

	
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


func (s *Service) GetSpace(ctx context.Context, id uuid.UUID) (*SpaceResponse, error) {
	space, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	return ToSpaceResponse(&space), nil
}


func (s *Service) GetSpaceBySlug(ctx context.Context, slug string) (*SpaceResponse, error) {
	space, err := s.store.GetSpaceBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	return ToSpaceResponse(&space), nil
}


func (s *Service) ListSpaces(ctx context.Context, page, limit int32) (*PaginatedSpacesResponse, error) {
	
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

	
	var spaceResponses []*SpaceResponse
	for i := range spaces {
		spaceResponses = append(spaceResponses, ToSpaceResponse(&spaces[i]))
	}

	return &PaginatedSpacesResponse{
		Spaces: spaceResponses,
		Total:  int64(len(spaceResponses)), 
		Page:   page,
		Limit:  limit,
	}, nil
}


func (s *Service) UpdateSpace(ctx context.Context, id uuid.UUID, req UpdateSpaceRequest) (*SpaceResponse, error) {
	
	_, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	
	if req.Slug != nil {
		existingSpace, err := s.store.GetSpaceBySlug(ctx, *req.Slug)
		if err == nil && existingSpace.ID != id {
			return nil, fmt.Errorf("%w: space with this slug already exists", util.ErrConflict)
		}
	}

	
	var settings pqtype.NullRawMessage
	if req.Settings != nil {
		settings = *req.Settings
	}

	
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


func (s *Service) DeleteSpace(ctx context.Context, id uuid.UUID) error {
	
	_, err := s.store.GetSpace(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: space not found", util.ErrNotFound)
	}

	
	if err := s.store.DeleteSpace(ctx, id); err != nil {
		return fmt.Errorf("failed to delete space: %w", err)
	}

	return nil
}
