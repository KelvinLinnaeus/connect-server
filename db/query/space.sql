-- name: CreateSpace :one
INSERT INTO spaces(name, slug, description, type, logo, location, website, contact_email, phone_number) VALUES(
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetSpace :one
SELECT * FROM spaces
WHERE id = $1 LIMIT 1;

-- name: GetSpaceBySlug :one
SELECT * FROM spaces
WHERE slug = $1 LIMIT 1;

-- name: ListSpaces :many
SELECT * FROM spaces
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateSpace :one
UPDATE spaces
SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    type = COALESCE(sqlc.narg('type'), type),
    logo = COALESCE(sqlc.narg('logo'), logo),
    cover_image = COALESCE(sqlc.narg('cover_image'), cover_image),
    location = COALESCE(sqlc.narg('location'), location),
    website = COALESCE(sqlc.narg('website'), website),
    contact_email = COALESCE(sqlc.narg('contact_email'), contact_email),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    status = COALESCE(sqlc.narg('status'), status),
    settings = COALESCE(sqlc.narg('settings'), settings),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSpace :exec
DELETE FROM spaces WHERE id = $1;