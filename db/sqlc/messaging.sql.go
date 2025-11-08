




package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
)

const addConversationParticipants = `-- name: AddConversationParticipants :exec
INSERT INTO conversation_participants (conversation_id, user_id, role)
SELECT $1, unnest($2::uuid[]), 'member'
ON CONFLICT (conversation_id, user_id) DO NOTHING
`

type AddConversationParticipantsParams struct {
	ConversationID uuid.UUID   `json:"conversation_id"`
	Column2        []uuid.UUID `json:"column_2"`
}

func (q *Queries) AddConversationParticipants(ctx context.Context, arg AddConversationParticipantsParams) error {
	_, err := q.db.ExecContext(ctx, addConversationParticipants, arg.ConversationID, pq.Array(arg.Column2))
	return err
}

const addMessageReaction = `-- name: AddMessageReaction :exec
UPDATE messages 
SET reactions = jsonb_set(
    COALESCE(reactions, '{}'::jsonb),
    ARRAY[$2],
    to_jsonb($3)
)
WHERE id = $1
`

type AddMessageReactionParams struct {
	ID      uuid.UUID   `json:"id"`
	Column2 interface{} `json:"column_2"`
	ToJsonb interface{} `json:"to_jsonb"`
}

func (q *Queries) AddMessageReaction(ctx context.Context, arg AddMessageReactionParams) error {
	_, err := q.db.ExecContext(ctx, addMessageReaction, arg.ID, arg.Column2, arg.ToJsonb)
	return err
}

const createConversation = `-- name: CreateConversation :one

INSERT INTO conversations (space_id, name, avatar, description, conversation_type, settings)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, space_id, name, avatar, description, conversation_type, last_message_id, last_message_at, is_active, settings, created_at, updated_at
`

type CreateConversationParams struct {
	SpaceID          uuid.UUID             `json:"space_id"`
	Name             sql.NullString        `json:"name"`
	Avatar           sql.NullString        `json:"avatar"`
	Description      sql.NullString        `json:"description"`
	ConversationType sql.NullString        `json:"conversation_type"`
	Settings         pqtype.NullRawMessage `json:"settings"`
}


func (q *Queries) CreateConversation(ctx context.Context, arg CreateConversationParams) (Conversation, error) {
	row := q.db.QueryRowContext(ctx, createConversation,
		arg.SpaceID,
		arg.Name,
		arg.Avatar,
		arg.Description,
		arg.ConversationType,
		arg.Settings,
	)
	var i Conversation
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Name,
		&i.Avatar,
		&i.Description,
		&i.ConversationType,
		&i.LastMessageID,
		&i.LastMessageAt,
		&i.IsActive,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteMessage = `-- name: DeleteMessage :exec
UPDATE messages 
SET status = 'deleted', content = '[message deleted]'
WHERE id = $1 AND sender_id = $2
`

type DeleteMessageParams struct {
	ID       uuid.UUID `json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
}

func (q *Queries) DeleteMessage(ctx context.Context, arg DeleteMessageParams) error {
	_, err := q.db.ExecContext(ctx, deleteMessage, arg.ID, arg.SenderID)
	return err
}

const getConversationByID = `-- name: GetConversationByID :one
SELECT 
    c.id, c.space_id, c.name, c.avatar, c.description, c.conversation_type, c.last_message_id, c.last_message_at, c.is_active, c.settings, c.created_at, c.updated_at,
    (SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.is_read = false AND m.recipient_id = $1) as unread_count
FROM conversations c
WHERE c.id = $2 AND c.is_active = true
`

type GetConversationByIDParams struct {
	RecipientID uuid.NullUUID `json:"recipient_id"`
	ID          uuid.UUID     `json:"id"`
}

type GetConversationByIDRow struct {
	ID               uuid.UUID             `json:"id"`
	SpaceID          uuid.UUID             `json:"space_id"`
	Name             sql.NullString        `json:"name"`
	Avatar           sql.NullString        `json:"avatar"`
	Description      sql.NullString        `json:"description"`
	ConversationType sql.NullString        `json:"conversation_type"`
	LastMessageID    uuid.NullUUID         `json:"last_message_id"`
	LastMessageAt    sql.NullTime          `json:"last_message_at"`
	IsActive         sql.NullBool          `json:"is_active"`
	Settings         pqtype.NullRawMessage `json:"settings"`
	CreatedAt        sql.NullTime          `json:"created_at"`
	UpdatedAt        sql.NullTime          `json:"updated_at"`
	UnreadCount      int64                 `json:"unread_count"`
}

func (q *Queries) GetConversationByID(ctx context.Context, arg GetConversationByIDParams) (GetConversationByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getConversationByID, arg.RecipientID, arg.ID)
	var i GetConversationByIDRow
	err := row.Scan(
		&i.ID,
		&i.SpaceID,
		&i.Name,
		&i.Avatar,
		&i.Description,
		&i.ConversationType,
		&i.LastMessageID,
		&i.LastMessageAt,
		&i.IsActive,
		&i.Settings,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UnreadCount,
	)
	return i, err
}

const getConversationByParticipants = `-- name: GetConversationByParticipants :one
SELECT c.id
FROM conversations c
JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
WHERE c.space_id = $1 
  AND c.conversation_type = 'direct'
  AND cp1.user_id = $2
  AND cp2.user_id = $3
  AND c.is_active = true
LIMIT 1
`

type GetConversationByParticipantsParams struct {
	SpaceID  uuid.UUID `json:"space_id"`
	UserID   uuid.UUID `json:"user_id"`
	UserID_2 uuid.UUID `json:"user_id_2"`
}

func (q *Queries) GetConversationByParticipants(ctx context.Context, arg GetConversationByParticipantsParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getConversationByParticipants, arg.SpaceID, arg.UserID, arg.UserID_2)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getConversationMessages = `-- name: GetConversationMessages :many
SELECT 
    m.id, m.conversation_id, m.sender_id, m.recipient_id, m.content, m.attachments, m.message_type, m.is_read, m.read_at, m.reactions, m.reply_to_id, m.status, m.created_at,
    u.username as sender_username,
    u.full_name as sender_full_name,
    u.avatar as sender_avatar,
    reply_msg.content as reply_content,
    reply_user.username as reply_username
FROM messages m
JOIN users u ON m.sender_id = u.id
LEFT JOIN messages reply_msg ON m.reply_to_id = reply_msg.id
LEFT JOIN users reply_user ON reply_msg.sender_id = reply_user.id
WHERE m.conversation_id = $1
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3
`

type GetConversationMessagesParams struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	Limit          int32     `json:"limit"`
	Offset         int32     `json:"offset"`
}

type GetConversationMessagesRow struct {
	ID             uuid.UUID             `json:"id"`
	ConversationID uuid.UUID             `json:"conversation_id"`
	SenderID       uuid.UUID             `json:"sender_id"`
	RecipientID    uuid.NullUUID         `json:"recipient_id"`
	Content        sql.NullString        `json:"content"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	MessageType    sql.NullString        `json:"message_type"`
	IsRead         sql.NullBool          `json:"is_read"`
	ReadAt         sql.NullTime          `json:"read_at"`
	Reactions      pqtype.NullRawMessage `json:"reactions"`
	ReplyToID      uuid.NullUUID         `json:"reply_to_id"`
	Status         sql.NullString        `json:"status"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	SenderUsername string                `json:"sender_username"`
	SenderFullName string                `json:"sender_full_name"`
	SenderAvatar   sql.NullString        `json:"sender_avatar"`
	ReplyContent   sql.NullString        `json:"reply_content"`
	ReplyUsername  sql.NullString        `json:"reply_username"`
}

func (q *Queries) GetConversationMessages(ctx context.Context, arg GetConversationMessagesParams) ([]GetConversationMessagesRow, error) {
	rows, err := q.db.QueryContext(ctx, getConversationMessages, arg.ConversationID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetConversationMessagesRow{}
	for rows.Next() {
		var i GetConversationMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.ConversationID,
			&i.SenderID,
			&i.RecipientID,
			&i.Content,
			&i.Attachments,
			&i.MessageType,
			&i.IsRead,
			&i.ReadAt,
			&i.Reactions,
			&i.ReplyToID,
			&i.Status,
			&i.CreatedAt,
			&i.SenderUsername,
			&i.SenderFullName,
			&i.SenderAvatar,
			&i.ReplyContent,
			&i.ReplyUsername,
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

const getConversationParticipants = `-- name: GetConversationParticipants :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    u.verified,
    cp.role,
    cp.joined_at,
    cp.is_active,
    cp.notifications_enabled
FROM conversation_participants cp
JOIN users u ON cp.user_id = u.id
WHERE cp.conversation_id = $1 AND u.status = 'active'
ORDER BY cp.joined_at
`

type GetConversationParticipantsRow struct {
	ID                   uuid.UUID      `json:"id"`
	Username             string         `json:"username"`
	FullName             string         `json:"full_name"`
	Avatar               sql.NullString `json:"avatar"`
	Verified             sql.NullBool   `json:"verified"`
	Role                 sql.NullString `json:"role"`
	JoinedAt             sql.NullTime   `json:"joined_at"`
	IsActive             sql.NullBool   `json:"is_active"`
	NotificationsEnabled sql.NullBool   `json:"notifications_enabled"`
}

func (q *Queries) GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]GetConversationParticipantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getConversationParticipants, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetConversationParticipantsRow{}
	for rows.Next() {
		var i GetConversationParticipantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.FullName,
			&i.Avatar,
			&i.Verified,
			&i.Role,
			&i.JoinedAt,
			&i.IsActive,
			&i.NotificationsEnabled,
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

const getMessageByID = `-- name: GetMessageByID :one
SELECT 
    m.id, m.conversation_id, m.sender_id, m.recipient_id, m.content, m.attachments, m.message_type, m.is_read, m.read_at, m.reactions, m.reply_to_id, m.status, m.created_at,
    u.username as sender_username,
    u.full_name as sender_full_name,
    u.avatar as sender_avatar
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.id = $1
`

type GetMessageByIDRow struct {
	ID             uuid.UUID             `json:"id"`
	ConversationID uuid.UUID             `json:"conversation_id"`
	SenderID       uuid.UUID             `json:"sender_id"`
	RecipientID    uuid.NullUUID         `json:"recipient_id"`
	Content        sql.NullString        `json:"content"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	MessageType    sql.NullString        `json:"message_type"`
	IsRead         sql.NullBool          `json:"is_read"`
	ReadAt         sql.NullTime          `json:"read_at"`
	Reactions      pqtype.NullRawMessage `json:"reactions"`
	ReplyToID      uuid.NullUUID         `json:"reply_to_id"`
	Status         sql.NullString        `json:"status"`
	CreatedAt      sql.NullTime          `json:"created_at"`
	SenderUsername string                `json:"sender_username"`
	SenderFullName string                `json:"sender_full_name"`
	SenderAvatar   sql.NullString        `json:"sender_avatar"`
}

func (q *Queries) GetMessageByID(ctx context.Context, id uuid.UUID) (GetMessageByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getMessageByID, id)
	var i GetMessageByIDRow
	err := row.Scan(
		&i.ID,
		&i.ConversationID,
		&i.SenderID,
		&i.RecipientID,
		&i.Content,
		&i.Attachments,
		&i.MessageType,
		&i.IsRead,
		&i.ReadAt,
		&i.Reactions,
		&i.ReplyToID,
		&i.Status,
		&i.CreatedAt,
		&i.SenderUsername,
		&i.SenderFullName,
		&i.SenderAvatar,
	)
	return i, err
}

const getOrCreateDirectConversation = `-- name: GetOrCreateDirectConversation :one
WITH existing_conversation AS (
    SELECT c.id
    FROM conversations c
    JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
    JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
    WHERE c.space_id = $1 
      AND c.conversation_type = 'direct'
      AND cp1.user_id = $2
      AND cp2.user_id = $3
      AND c.is_active = true
    LIMIT 1
),
new_conversation AS (
    INSERT INTO conversations (space_id, conversation_type, settings)
    SELECT $1, 'direct', '{}'
    WHERE NOT EXISTS (SELECT 1 FROM existing_conversation)
    RETURNING id
)
SELECT id FROM existing_conversation
UNION ALL
SELECT id FROM new_conversation
`

type GetOrCreateDirectConversationParams struct {
	SpaceID  uuid.UUID `json:"space_id"`
	UserID   uuid.UUID `json:"user_id"`
	UserID_2 uuid.UUID `json:"user_id_2"`
}

func (q *Queries) GetOrCreateDirectConversation(ctx context.Context, arg GetOrCreateDirectConversationParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getOrCreateDirectConversation, arg.SpaceID, arg.UserID, arg.UserID_2)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUnreadMessageCount = `-- name: GetUnreadMessageCount :one
SELECT COUNT(*) FROM messages 
WHERE conversation_id = $1 AND recipient_id = $2 AND is_read = false
`

type GetUnreadMessageCountParams struct {
	ConversationID uuid.UUID     `json:"conversation_id"`
	RecipientID    uuid.NullUUID `json:"recipient_id"`
}

func (q *Queries) GetUnreadMessageCount(ctx context.Context, arg GetUnreadMessageCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getUnreadMessageCount, arg.ConversationID, arg.RecipientID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUserConversations = `-- name: GetUserConversations :many
SELECT 
    c.id, c.space_id, c.name, c.avatar, c.description, c.conversation_type, c.last_message_id, c.last_message_at, c.is_active, c.settings, c.created_at, c.updated_at,
    cp.role as user_role,
    cp.notifications_enabled,
    cp.custom_settings,
    (SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.is_read = false AND m.recipient_id = $1) as unread_count,
    last_msg.content as last_message_content,
    last_msg.created_at as last_message_time,
    last_sender.username as last_sender_username,
    last_sender.full_name as last_sender_full_name
FROM conversations c
JOIN conversation_participants cp ON c.id = cp.conversation_id
LEFT JOIN LATERAL (
    SELECT m2.content, m2.created_at, m2.sender_id
    FROM messages m2
    WHERE m2.conversation_id = c.id
    ORDER BY m2.created_at DESC
    LIMIT 1
) last_msg ON true
LEFT JOIN users last_sender ON last_msg.sender_id = last_sender.id
WHERE cp.user_id = $1 AND cp.is_active = true AND c.is_active = true
ORDER BY c.last_message_at DESC NULLS LAST
`

type GetUserConversationsRow struct {
	ID                   uuid.UUID             `json:"id"`
	SpaceID              uuid.UUID             `json:"space_id"`
	Name                 sql.NullString        `json:"name"`
	Avatar               sql.NullString        `json:"avatar"`
	Description          sql.NullString        `json:"description"`
	ConversationType     sql.NullString        `json:"conversation_type"`
	LastMessageID        uuid.NullUUID         `json:"last_message_id"`
	LastMessageAt        sql.NullTime          `json:"last_message_at"`
	IsActive             sql.NullBool          `json:"is_active"`
	Settings             pqtype.NullRawMessage `json:"settings"`
	CreatedAt            sql.NullTime          `json:"created_at"`
	UpdatedAt            sql.NullTime          `json:"updated_at"`
	UserRole             sql.NullString        `json:"user_role"`
	NotificationsEnabled sql.NullBool          `json:"notifications_enabled"`
	CustomSettings       pqtype.NullRawMessage `json:"custom_settings"`
	UnreadCount          int64                 `json:"unread_count"`
	LastMessageContent   sql.NullString        `json:"last_message_content"`
	LastMessageTime      sql.NullTime          `json:"last_message_time"`
	LastSenderUsername   sql.NullString        `json:"last_sender_username"`
	LastSenderFullName   sql.NullString        `json:"last_sender_full_name"`
}

func (q *Queries) GetUserConversations(ctx context.Context, recipientID uuid.NullUUID) ([]GetUserConversationsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserConversations, recipientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserConversationsRow{}
	for rows.Next() {
		var i GetUserConversationsRow
		if err := rows.Scan(
			&i.ID,
			&i.SpaceID,
			&i.Name,
			&i.Avatar,
			&i.Description,
			&i.ConversationType,
			&i.LastMessageID,
			&i.LastMessageAt,
			&i.IsActive,
			&i.Settings,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserRole,
			&i.NotificationsEnabled,
			&i.CustomSettings,
			&i.UnreadCount,
			&i.LastMessageContent,
			&i.LastMessageTime,
			&i.LastSenderUsername,
			&i.LastSenderFullName,
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

const leaveConversation = `-- name: LeaveConversation :exec
UPDATE conversation_participants 
SET is_active = false, left_at = NOW()
WHERE conversation_id = $1 AND user_id = $2
`

type LeaveConversationParams struct {
	ConversationID uuid.UUID `json:"conversation_id"`
	UserID         uuid.UUID `json:"user_id"`
}

func (q *Queries) LeaveConversation(ctx context.Context, arg LeaveConversationParams) error {
	_, err := q.db.ExecContext(ctx, leaveConversation, arg.ConversationID, arg.UserID)
	return err
}

const markMessagesAsRead = `-- name: MarkMessagesAsRead :exec
UPDATE messages 
SET is_read = true, read_at = NOW()
WHERE conversation_id = $1 AND recipient_id = $2 AND is_read = false
`

type MarkMessagesAsReadParams struct {
	ConversationID uuid.UUID     `json:"conversation_id"`
	RecipientID    uuid.NullUUID `json:"recipient_id"`
}

func (q *Queries) MarkMessagesAsRead(ctx context.Context, arg MarkMessagesAsReadParams) error {
	_, err := q.db.ExecContext(ctx, markMessagesAsRead, arg.ConversationID, arg.RecipientID)
	return err
}

const removeMessageReaction = `-- name: RemoveMessageReaction :exec
UPDATE messages 
SET reactions = reactions - $2
WHERE id = $1
`

type RemoveMessageReactionParams struct {
	ID        uuid.UUID             `json:"id"`
	Reactions pqtype.NullRawMessage `json:"reactions"`
}

func (q *Queries) RemoveMessageReaction(ctx context.Context, arg RemoveMessageReactionParams) error {
	_, err := q.db.ExecContext(ctx, removeMessageReaction, arg.ID, arg.Reactions)
	return err
}

const sendMessage = `-- name: SendMessage :one
INSERT INTO messages (conversation_id, sender_id, recipient_id, content, attachments, message_type, reply_to_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, conversation_id, sender_id, recipient_id, content, attachments, message_type, is_read, read_at, reactions, reply_to_id, status, created_at
`

type SendMessageParams struct {
	ConversationID uuid.UUID             `json:"conversation_id"`
	SenderID       uuid.UUID             `json:"sender_id"`
	RecipientID    uuid.NullUUID         `json:"recipient_id"`
	Content        sql.NullString        `json:"content"`
	Attachments    pqtype.NullRawMessage `json:"attachments"`
	MessageType    sql.NullString        `json:"message_type"`
	ReplyToID      uuid.NullUUID         `json:"reply_to_id"`
}

func (q *Queries) SendMessage(ctx context.Context, arg SendMessageParams) (Message, error) {
	row := q.db.QueryRowContext(ctx, sendMessage,
		arg.ConversationID,
		arg.SenderID,
		arg.RecipientID,
		arg.Content,
		arg.Attachments,
		arg.MessageType,
		arg.ReplyToID,
	)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.ConversationID,
		&i.SenderID,
		&i.RecipientID,
		&i.Content,
		&i.Attachments,
		&i.MessageType,
		&i.IsRead,
		&i.ReadAt,
		&i.Reactions,
		&i.ReplyToID,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const updateConversationLastMessage = `-- name: UpdateConversationLastMessage :exec
UPDATE conversations 
SET last_message_id = $1, last_message_at = NOW(), updated_at = NOW()
WHERE id = $2
`

type UpdateConversationLastMessageParams struct {
	LastMessageID uuid.NullUUID `json:"last_message_id"`
	ID            uuid.UUID     `json:"id"`
}

func (q *Queries) UpdateConversationLastMessage(ctx context.Context, arg UpdateConversationLastMessageParams) error {
	_, err := q.db.ExecContext(ctx, updateConversationLastMessage, arg.LastMessageID, arg.ID)
	return err
}

const updateConversationSettings = `-- name: UpdateConversationSettings :exec
UPDATE conversations 
SET settings = $1, updated_at = NOW()
WHERE id = $2
`

type UpdateConversationSettingsParams struct {
	Settings pqtype.NullRawMessage `json:"settings"`
	ID       uuid.UUID             `json:"id"`
}

func (q *Queries) UpdateConversationSettings(ctx context.Context, arg UpdateConversationSettingsParams) error {
	_, err := q.db.ExecContext(ctx, updateConversationSettings, arg.Settings, arg.ID)
	return err
}

const updateParticipantSettings = `-- name: UpdateParticipantSettings :exec
UPDATE conversation_participants 
SET notifications_enabled = $1, custom_settings = $2
WHERE conversation_id = $3 AND user_id = $4
`

type UpdateParticipantSettingsParams struct {
	NotificationsEnabled sql.NullBool          `json:"notifications_enabled"`
	CustomSettings       pqtype.NullRawMessage `json:"custom_settings"`
	ConversationID       uuid.UUID             `json:"conversation_id"`
	UserID               uuid.UUID             `json:"user_id"`
}

func (q *Queries) UpdateParticipantSettings(ctx context.Context, arg UpdateParticipantSettingsParams) error {
	_, err := q.db.ExecContext(ctx, updateParticipantSettings,
		arg.NotificationsEnabled,
		arg.CustomSettings,
		arg.ConversationID,
		arg.UserID,
	)
	return err
}
