




package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
    to_user_id, from_user_id, type, title, message, related_id, metadata, priority, action_required
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, to_user_id, from_user_id, type, title, message, related_id, metadata, is_read, priority, action_required, created_at
`

type CreateNotificationParams struct {
	ToUserID       uuid.UUID             `json:"to_user_id"`
	FromUserID     uuid.NullUUID         `json:"from_user_id"`
	Type           string                `json:"type"`
	Title          sql.NullString        `json:"title"`
	Message        sql.NullString        `json:"message"`
	RelatedID      uuid.NullUUID         `json:"related_id"`
	Metadata       pqtype.NullRawMessage `json:"metadata"`
	Priority       sql.NullString        `json:"priority"`
	ActionRequired sql.NullBool          `json:"action_required"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRowContext(ctx, createNotification,
		arg.ToUserID,
		arg.FromUserID,
		arg.Type,
		arg.Title,
		arg.Message,
		arg.RelatedID,
		arg.Metadata,
		arg.Priority,
		arg.ActionRequired,
	)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.ToUserID,
		&i.FromUserID,
		&i.Type,
		&i.Title,
		&i.Message,
		&i.RelatedID,
		&i.Metadata,
		&i.IsRead,
		&i.Priority,
		&i.ActionRequired,
		&i.CreatedAt,
	)
	return i, err
}

const deleteNotification = `-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1
`

func (q *Queries) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteNotification, id)
	return err
}

const getAdminNotifications = `-- name: GetAdminNotifications :many

SELECT n.id, n.to_user_id, n.from_user_id, n.type, n.title, n.message, n.related_id, n.metadata, n.is_read, n.priority, n.action_required, n.created_at,
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
LIMIT $5 OFFSET $6
`

type GetAdminNotificationsParams struct {
	ToUserID uuid.UUID      `json:"to_user_id"`
	Type     string         `json:"type"`
	Priority sql.NullString `json:"priority"`
	IsRead   sql.NullBool   `json:"is_read"`
	Limit    int32          `json:"limit"`
	Offset   int32          `json:"offset"`
}

type GetAdminNotificationsRow struct {
	ID             uuid.UUID             `json:"id"`
	ToUserID       uuid.UUID             `json:"to_user_id"`
	FromUserID     uuid.NullUUID         `json:"from_user_id"`
	Type           string                `json:"type"`
	Title          sql.NullString        `json:"title"`
	Message        sql.NullString        `json:"message"`
	RelatedID      uuid.NullUUID         `json:"related_id"`
	Metadata       pqtype.NullRawMessage `json:"metadata"`
	IsRead         sql.NullBool          `json:"is_read"`
	Priority       sql.NullString        `json:"priority"`
	ActionRequired sql.NullBool          `json:"action_required"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	FromUsername   sql.NullString        `json:"from_username"`
	FromFullName   sql.NullString        `json:"from_full_name"`
	FromAvatar     sql.NullString        `json:"from_avatar"`
}


func (q *Queries) GetAdminNotifications(ctx context.Context, arg GetAdminNotificationsParams) ([]GetAdminNotificationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAdminNotifications,
		arg.ToUserID,
		arg.Type,
		arg.Priority,
		arg.IsRead,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAdminNotificationsRow{}
	for rows.Next() {
		var i GetAdminNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.ToUserID,
			&i.FromUserID,
			&i.Type,
			&i.Title,
			&i.Message,
			&i.RelatedID,
			&i.Metadata,
			&i.IsRead,
			&i.Priority,
			&i.ActionRequired,
			&i.CreatedAt,
			&i.FromUsername,
			&i.FromFullName,
			&i.FromAvatar,
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

const getNotification = `-- name: GetNotification :one
SELECT id, to_user_id, from_user_id, type, title, message, related_id, metadata, is_read, priority, action_required, created_at FROM notifications
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNotification(ctx context.Context, id uuid.UUID) (Notification, error) {
	row := q.db.QueryRowContext(ctx, getNotification, id)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.ToUserID,
		&i.FromUserID,
		&i.Type,
		&i.Title,
		&i.Message,
		&i.RelatedID,
		&i.Metadata,
		&i.IsRead,
		&i.Priority,
		&i.ActionRequired,
		&i.CreatedAt,
	)
	return i, err
}

const getUnreadCount = `-- name: GetUnreadCount :one
SELECT COUNT(*) as count
FROM notifications
WHERE to_user_id = $1 AND is_read = false
`

func (q *Queries) GetUnreadCount(ctx context.Context, toUserID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUnreadCount, toUserID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserNotifications = `-- name: GetUserNotifications :many
SELECT n.id, n.to_user_id, n.from_user_id, n.type, n.title, n.message, n.related_id, n.metadata, n.is_read, n.priority, n.action_required, n.created_at,
       u.username as from_username,
       u.full_name as from_full_name,
       u.avatar as from_avatar
FROM notifications n
LEFT JOIN users u ON n.from_user_id = u.id
WHERE n.to_user_id = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3
`

type GetUserNotificationsParams struct {
	ToUserID uuid.UUID `json:"to_user_id"`
	Limit    int32     `json:"limit"`
	Offset   int32     `json:"offset"`
}

type GetUserNotificationsRow struct {
	ID             uuid.UUID             `json:"id"`
	ToUserID       uuid.UUID             `json:"to_user_id"`
	FromUserID     uuid.NullUUID         `json:"from_user_id"`
	Type           string                `json:"type"`
	Title          sql.NullString        `json:"title"`
	Message        sql.NullString        `json:"message"`
	RelatedID      uuid.NullUUID         `json:"related_id"`
	Metadata       pqtype.NullRawMessage `json:"metadata"`
	IsRead         sql.NullBool          `json:"is_read"`
	Priority       sql.NullString        `json:"priority"`
	ActionRequired sql.NullBool          `json:"action_required"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	FromUsername   sql.NullString        `json:"from_username"`
	FromFullName   sql.NullString        `json:"from_full_name"`
	FromAvatar     sql.NullString        `json:"from_avatar"`
}

func (q *Queries) GetUserNotifications(ctx context.Context, arg GetUserNotificationsParams) ([]GetUserNotificationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserNotifications, arg.ToUserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserNotificationsRow{}
	for rows.Next() {
		var i GetUserNotificationsRow
		if err := rows.Scan(
			&i.ID,
			&i.ToUserID,
			&i.FromUserID,
			&i.Type,
			&i.Title,
			&i.Message,
			&i.RelatedID,
			&i.Metadata,
			&i.IsRead,
			&i.Priority,
			&i.ActionRequired,
			&i.CreatedAt,
			&i.FromUsername,
			&i.FromFullName,
			&i.FromAvatar,
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

const markAllAsRead = `-- name: MarkAllAsRead :exec
UPDATE notifications
SET is_read = true
WHERE to_user_id = $1 AND is_read = false
`

func (q *Queries) MarkAllAsRead(ctx context.Context, toUserID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markAllAsRead, toUserID)
	return err
}

const markAsRead = `-- name: MarkAsRead :exec
UPDATE notifications
SET is_read = true
WHERE id = $1
`

func (q *Queries) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, markAsRead, id)
	return err
}
