-- Messaging System Queries

-- name: CreateConversation :one
INSERT INTO conversations (space_id, name, avatar, description, conversation_type, settings)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetConversationByID :one
SELECT 
    c.*,
    (SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.is_read = false AND m.recipient_id = $1) as unread_count
FROM conversations c
WHERE c.id = $2 AND c.is_active = true;

-- name: GetOrCreateDirectConversation :one
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
SELECT id FROM new_conversation;

-- name: AddConversationParticipants :exec
INSERT INTO conversation_participants (conversation_id, user_id, role)
SELECT $1, unnest($2::uuid[]), 'member'
ON CONFLICT (conversation_id, user_id) DO NOTHING;

-- name: GetUserConversations :many
SELECT 
    c.*,
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
ORDER BY c.last_message_at DESC NULLS LAST;

-- name: GetConversationParticipants :many
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
ORDER BY cp.joined_at;

-- name: SendMessage :one
INSERT INTO messages (conversation_id, sender_id, recipient_id, content, attachments, message_type, reply_to_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateConversationLastMessage :exec
UPDATE conversations 
SET last_message_id = $1, last_message_at = NOW(), updated_at = NOW()
WHERE id = $2;

-- name: GetConversationMessages :many
SELECT 
    m.*,
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
LIMIT $2 OFFSET $3;

-- name: MarkMessagesAsRead :exec
UPDATE messages 
SET is_read = true, read_at = NOW()
WHERE conversation_id = $1 AND recipient_id = $2 AND is_read = false;

-- name: GetUnreadMessageCount :one
SELECT COUNT(*) FROM messages 
WHERE conversation_id = $1 AND recipient_id = $2 AND is_read = false;

-- name: AddMessageReaction :exec
UPDATE messages 
SET reactions = jsonb_set(
    COALESCE(reactions, '{}'::jsonb),
    ARRAY[$2],
    to_jsonb($3)
)
WHERE id = $1;

-- name: RemoveMessageReaction :exec
UPDATE messages 
SET reactions = reactions - $2
WHERE id = $1;

-- name: LeaveConversation :exec
UPDATE conversation_participants 
SET is_active = false, left_at = NOW()
WHERE conversation_id = $1 AND user_id = $2;

-- name: UpdateConversationSettings :exec
UPDATE conversations 
SET settings = $1, updated_at = NOW()
WHERE id = $2;

-- name: UpdateParticipantSettings :exec
UPDATE conversation_participants 
SET notifications_enabled = $1, custom_settings = $2
WHERE conversation_id = $3 AND user_id = $4;

-- name: GetConversationByParticipants :one
SELECT c.id
FROM conversations c
JOIN conversation_participants cp1 ON c.id = cp1.conversation_id
JOIN conversation_participants cp2 ON c.id = cp2.conversation_id
WHERE c.space_id = $1 
  AND c.conversation_type = 'direct'
  AND cp1.user_id = $2
  AND cp2.user_id = $3
  AND c.is_active = true
LIMIT 1;

-- name: DeleteMessage :exec
UPDATE messages 
SET status = 'deleted', content = '[message deleted]'
WHERE id = $1 AND sender_id = $2;

-- name: GetMessageByID :one
SELECT 
    m.*,
    u.username as sender_username,
    u.full_name as sender_full_name,
    u.avatar as sender_avatar
FROM messages m
JOIN users u ON m.sender_id = u.id
WHERE m.id = $1;