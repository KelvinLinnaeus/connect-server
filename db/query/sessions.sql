-- name: CreateSession :one
INSERT INTO user_sessions (
    id,
    user_id,
    username,
    refresh_token,
    user_agent,
    is_blocked,
    space_id,
    ip_address,
    last_activity,
    expires_at
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;


-- name: GetSession :one
SELECT *
FROM user_sessions
WHERE id = $1
LIMIT 1;

