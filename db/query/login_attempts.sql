-- name: CreateLoginAttempt :one
INSERT INTO login_attempts (
    username,
    ip_address,
    user_agent,
    attempt_result,
    user_id,
    space_id,
    session_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetRecentLoginAttemptsByUsername :many
SELECT * FROM login_attempts
WHERE username = $1
  AND attempted_at > $2
ORDER BY attempted_at DESC;

-- name: GetRecentLoginAttemptsByIP :many
SELECT * FROM login_attempts
WHERE ip_address = $1
  AND attempted_at > $2
ORDER BY attempted_at DESC;

-- name: GetRecentFailedLoginAttemptsByUsername :many
SELECT * FROM login_attempts
WHERE username = $1
  AND attempt_result IN ('failed_password', 'failed_user_not_found')
  AND attempted_at > $2
ORDER BY attempted_at DESC;

-- name: GetRecentFailedLoginAttemptsByIP :many
SELECT * FROM login_attempts
WHERE ip_address = $1
  AND attempt_result IN ('failed_password', 'failed_user_not_found')
  AND attempted_at > $2
ORDER BY attempted_at DESC;

-- name: CountRecentFailedLoginAttemptsByUsername :one
SELECT COUNT(*) FROM login_attempts
WHERE username = $1
  AND attempt_result IN ('failed_password', 'failed_user_not_found')
  AND attempted_at > $2;

-- name: CountRecentFailedLoginAttemptsByIP :one
SELECT COUNT(*) FROM login_attempts
WHERE ip_address = $1
  AND attempt_result IN ('failed_password', 'failed_user_not_found')
  AND attempted_at > $2;

-- name: CleanupOldLoginAttempts :exec
DELETE FROM login_attempts
WHERE attempted_at < $1;

-- Note: Account lockouts are now managed via users table fields (is_locked, locked_until, failed_login_attempts)

-- name: UpdateUserLockStatus :one
UPDATE users
SET is_locked = $2,
    locked_until = $3,
    failed_login_attempts = $4,
    last_failed_login = $5
WHERE id = $1
RETURNING *;

-- name: IncrementFailedLoginAttempts :one
UPDATE users
SET failed_login_attempts = failed_login_attempts + 1,
    last_failed_login = NOW()
WHERE id = $1
RETURNING *;

-- name: ResetFailedLoginAttempts :one
UPDATE users
SET failed_login_attempts = 0,
    last_failed_login = NULL
WHERE id = $1
RETURNING *;

-- name: GetLockedUsers :many
SELECT id, username, email, is_locked, locked_until, failed_login_attempts, last_failed_login
FROM users
WHERE is_locked = TRUE
ORDER BY locked_until DESC;

-- name: UnlockExpiredAccounts :exec
UPDATE users
SET is_locked = FALSE,
    locked_until = NULL,
    failed_login_attempts = 0
WHERE is_locked = TRUE
  AND locked_until IS NOT NULL
  AND locked_until < NOW();

-- name: GetLoginAttemptsWithSessions :many
SELECT
    la.*,
    us.refresh_token,
    us.is_blocked as session_is_blocked,
    us.expires_at as session_expires_at
FROM login_attempts la
LEFT JOIN user_sessions us ON la.session_id = us.id
WHERE la.user_id = $1
  AND la.attempted_at > $2
ORDER BY la.attempted_at DESC
LIMIT $3;
