




package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createSpace = `-- name: CreateSpace :one
INSERT INTO spaces(name, slug, description, type, logo, location, website, contact_email, phone_number) VALUES(
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, name, slug, description, type, logo, cover_image, location, website, contact_email, phone_number, status, settings, created_at, updated_at
`

type CreateSpaceParams struct {
	Name         string         `json:"name"`
	Slug         string         `json:"slug"`
	Description  sql.NullString `json:"description"`
	Type         sql.NullString `json:"type"`
	Logo         sql.NullString `json:"logo"`
	Location     sql.NullString `json:"location"`
	Website      sql.NullString `json:"website"`
	ContactEmail sql.NullString `json:"contact_email"`
	PhoneNumber  sql.NullString `json:"phone_number"`
}

func (q *Queries) CreateSpace(ctx context.Context, arg CreateSpaceParams) (Space, error) {
	row := q.db.QueryRowContext(ctx, createSpace,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.Type,
		arg.Logo,
		arg.Location,
		arg.Website,
		arg.ContactEmail,
		arg.PhoneNumber,
	)
	var i Space
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Type,
		&i.Logo,
		&i.CoverImage,
		&i.Location,
		&i.Website,
		&i.ContactEmail,
		&i.PhoneNumber,
		&i.Status,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteSpace = `-- name: DeleteSpace :exec
DELETE FROM spaces WHERE id = $1
`

func (q *Queries) DeleteSpace(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteSpace, id)
	return err
}

const getSpace = `-- name: GetSpace :one
SELECT id, name, slug, description, type, logo, cover_image, location, website, contact_email, phone_number, status, settings, created_at, updated_at FROM spaces
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSpace(ctx context.Context, id uuid.UUID) (Space, error) {
	row := q.db.QueryRowContext(ctx, getSpace, id)
	var i Space
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Type,
		&i.Logo,
		&i.CoverImage,
		&i.Location,
		&i.Website,
		&i.ContactEmail,
		&i.PhoneNumber,
		&i.Status,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSpaceBySlug = `-- name: GetSpaceBySlug :one
SELECT id, name, slug, description, type, logo, cover_image, location, website, contact_email, phone_number, status, settings, created_at, updated_at FROM spaces
WHERE slug = $1 LIMIT 1
`

func (q *Queries) GetSpaceBySlug(ctx context.Context, slug string) (Space, error) {
	row := q.db.QueryRowContext(ctx, getSpaceBySlug, slug)
	var i Space
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Type,
		&i.Logo,
		&i.CoverImage,
		&i.Location,
		&i.Website,
		&i.ContactEmail,
		&i.PhoneNumber,
		&i.Status,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listSpaces = `-- name: ListSpaces :many
SELECT id, name, slug, description, type, logo, cover_image, location, website, contact_email, phone_number, status, settings, created_at, updated_at FROM spaces
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListSpacesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListSpaces(ctx context.Context, arg ListSpacesParams) ([]Space, error) {
	rows, err := q.db.QueryContext(ctx, listSpaces, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Space{}
	for rows.Next() {
		var i Space
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.Description,
			&i.Type,
			&i.Logo,
			&i.CoverImage,
			&i.Location,
			&i.Website,
			&i.ContactEmail,
			&i.PhoneNumber,
			&i.Status,
			&i.Settings,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateSpace = `-- name: UpdateSpace :one
UPDATE spaces
SET
    name = COALESCE($1, name),
    slug = COALESCE($2, slug),
    description = COALESCE($3, description),
    type = COALESCE($4, type),
    logo = COALESCE($5, logo),
    cover_image = COALESCE($6, cover_image),
    location = COALESCE($7, location),
    website = COALESCE($8, website),
    contact_email = COALESCE($9, contact_email),
    phone_number = COALESCE($10, phone_number),
    status = COALESCE($11, status),
    settings = COALESCE($12, settings),
    updated_at = NOW()
WHERE id = $13
RETURNING id, name, slug, description, type, logo, cover_image, location, website, contact_email, phone_number, status, settings, created_at, updated_at
`

type UpdateSpaceParams struct {
	Name         sql.NullString        `json:"name"`
	Slug         sql.NullString        `json:"slug"`
	Description  sql.NullString        `json:"description"`
	Type         sql.NullString        `json:"type"`
	Logo         sql.NullString        `json:"logo"`
	CoverImage   sql.NullString        `json:"cover_image"`
	Location     sql.NullString        `json:"location"`
	Website      sql.NullString        `json:"website"`
	ContactEmail sql.NullString        `json:"contact_email"`
	PhoneNumber  sql.NullString        `json:"phone_number"`
	Status       sql.NullString        `json:"status"`
	Settings     pqtype.NullRawMessage `json:"settings"`
	ID           uuid.UUID             `json:"id"`
}

func (q *Queries) UpdateSpace(ctx context.Context, arg UpdateSpaceParams) (Space, error) {
	row := q.db.QueryRowContext(ctx, updateSpace,
		arg.Name,
		arg.Slug,
		arg.Description,
		arg.Type,
		arg.Logo,
		arg.CoverImage,
		arg.Location,
		arg.Website,
		arg.ContactEmail,
		arg.PhoneNumber,
		arg.Status,
		arg.Settings,
		arg.ID,
	)
	var i Space
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.Description,
		&i.Type,
		&i.Logo,
		&i.CoverImage,
		&i.Location,
		&i.Website,
		&i.ContactEmail,
		&i.PhoneNumber,
		&i.Status,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
