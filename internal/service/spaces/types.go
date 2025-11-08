package spaces

import (
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// CreateSpaceRequest represents the request to create a space
type CreateSpaceRequest struct {
	Name         string  `json:"name" binding:"required,min=2,max=100"`
	Slug         string  `json:"slug" binding:"required,min=2,max=100"`
	Description  *string `json:"description,omitempty"`
	Type         *string `json:"type,omitempty"`
	Logo         *string `json:"logo,omitempty"`
	Location     *string `json:"location,omitempty"`
	Website      *string `json:"website,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty"`
	PhoneNumber  *string `json:"phone_number,omitempty"`
}

// UpdateSpaceRequest represents the request to update a space
type UpdateSpaceRequest struct {
	Name         *string                `json:"name,omitempty"`
	Slug         *string                `json:"slug,omitempty"`
	Description  *string                `json:"description,omitempty"`
	Type         *string                `json:"type,omitempty"`
	Logo         *string                `json:"logo,omitempty"`
	CoverImage   *string                `json:"cover_image,omitempty"`
	Location     *string                `json:"location,omitempty"`
	Website      *string                `json:"website,omitempty"`
	ContactEmail *string                `json:"contact_email,omitempty"`
	PhoneNumber  *string                `json:"phone_number,omitempty"`
	Status       *string                `json:"status,omitempty"`
	Settings     *pqtype.NullRawMessage `json:"settings,omitempty"`
}

// SpaceResponse represents the response for a space
type SpaceResponse struct {
	ID           uuid.UUID             `json:"id"`
	Name         string                `json:"name"`
	Slug         string                `json:"slug"`
	Description  *string               `json:"description,omitempty"`
	Type         *string               `json:"type,omitempty"`
	Logo         *string               `json:"logo,omitempty"`
	CoverImage   *string               `json:"cover_image,omitempty"`
	Location     *string               `json:"location,omitempty"`
	Website      *string               `json:"website,omitempty"`
	ContactEmail *string               `json:"contact_email,omitempty"`
	PhoneNumber  *string               `json:"phone_number,omitempty"`
	Status       *string               `json:"status,omitempty"`
	Settings     *pqtype.NullRawMessage `json:"settings,omitempty"`
	CreatedAt    *time.Time            `json:"created_at,omitempty"`
	UpdatedAt    *time.Time            `json:"updated_at,omitempty"`
}

// PaginatedSpacesResponse represents paginated spaces response
type PaginatedSpacesResponse struct {
	Spaces []*SpaceResponse `json:"spaces"`
	Total  int64            `json:"total"`
	Page   int32            `json:"page"`
	Limit  int32            `json:"limit"`
}

// ToSpaceResponse converts a database space to a response
func ToSpaceResponse(space *db.Space) *SpaceResponse {
	if space == nil {
		return nil
	}

	response := &SpaceResponse{
		ID:   space.ID,
		Name: space.Name,
		Slug: space.Slug,
	}

	if space.Description.Valid {
		response.Description = &space.Description.String
	}
	if space.Type.Valid {
		response.Type = &space.Type.String
	}
	if space.Logo.Valid {
		response.Logo = &space.Logo.String
	}
	if space.CoverImage.Valid {
		response.CoverImage = &space.CoverImage.String
	}
	if space.Location.Valid {
		response.Location = &space.Location.String
	}
	if space.Website.Valid {
		response.Website = &space.Website.String
	}
	if space.ContactEmail.Valid {
		response.ContactEmail = &space.ContactEmail.String
	}
	if space.PhoneNumber.Valid {
		response.PhoneNumber = &space.PhoneNumber.String
	}
	if space.Status.Valid {
		response.Status = &space.Status.String
	}
	if space.Settings.Valid {
		response.Settings = &space.Settings
	}
	if space.CreatedAt.Valid {
		response.CreatedAt = &space.CreatedAt.Time
	}
	if space.UpdatedAt.Valid {
		response.UpdatedAt = &space.UpdatedAt.Time
	}

	return response
}
