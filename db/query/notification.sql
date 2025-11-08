-- name: CreateNotification :one
INSERT INTO notifications (
    to_user_id, from_user_id, type, title, message, related_id, metadata, priority, action_required
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetUserNotifications :many
SELECT n.*,
       u.username as from_username,
       u.full_name as from_full_name,
       u.avatar as from_avatar
FROM notifications n
LEFT JOIN users u ON n.from_user_id = u.id
WHERE n.to_user_id = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetNotification :one
SELECT * FROM notifications
WHERE id = $1
LIMIT 1;

-- name: MarkAsRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1;

-- name: MarkAllAsRead :exec
UPDATE notifications
SET is_read = true
WHERE to_user_id = $1 AND is_read = false;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1;

-- name: GetUnreadCount :one
SELECT COUNT(*) as count
FROM notifications
WHERE to_user_id = $1 AND is_read = false;

-- Admin-specific Notification Queries

-- name: GetAdminNotifications :many
SELECT n.*,
       u.username as from_username,
       u.full_name as from_full_name,
       u.avatar as from_avatar
FROM notifications n
LEFT JOIN users u ON n.from_user_id = u.id
WHERE n.to_user_id = $1
  AND (n.type = $2 OR $2 = '')
  AND (n.priority = $3 OR $3 = '')
  AND (n.is_read = $4 OR $4 IS NULL)
ORDER BY n.created_at DESC
LIMIT $5 OFFSET $6;