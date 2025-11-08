-- Events and Announcements Queries

-- name: CreateEvent :one
INSERT INTO events (
    space_id, title, description, category, location, venue_details,
    start_date, end_date, timezone, organizer, tags, image_url,
    max_attendees, registration_required, registration_deadline, is_public
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: GetEventByID :one
SELECT 
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea2 WHERE ea2.event_id = e.id AND ea2.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea3 WHERE ea3.event_id = e.id AND ea3.user_id = $1) as is_registered,
    ea.status as user_attendance_status
FROM events e
JOIN users u ON e.organizer = u.id
LEFT JOIN event_attendees ea ON e.id = ea.event_id AND ea.user_id = $1
WHERE e.id = $2 AND e.status = 'published';

-- name: UpdateEvent :one
UPDATE events 
SET 
    title = $1,
    description = $2,
    category = $3,
    location = $4,
    venue_details = $5,
    start_date = $6,
    end_date = $7,
    timezone = $8,
    tags = $9,
    image_url = $10,
    max_attendees = $11,
    registration_required = $12,
    registration_deadline = $13,
    is_public = $14,
    updated_at = NOW()
WHERE id = $15
RETURNING *;

-- name: ListEvents :many
SELECT
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea4 WHERE ea4.event_id = e.id AND ea4.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea5 WHERE ea5.event_id = e.id AND ea5.user_id = $1) as is_registered
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $2
  AND e.status = 'published'
  AND (e.is_public = true OR e.organizer = $1)
  AND (e.start_date >= $3 OR $3 IS NULL)
  AND (e.category = $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'upcoming' THEN e.start_date END ASC,
    CASE WHEN $5 = 'popular' THEN (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') END DESC,
    e.created_at DESC
LIMIT $6 OFFSET $7;

-- name: SearchEvents :many
SELECT 
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea6 WHERE ea6.event_id = e.id AND ea6.status = 'registered') as current_attendees_count,
    EXISTS(SELECT 1 FROM event_attendees ea7 WHERE ea7.event_id = e.id AND ea7.user_id = $1) as is_registered
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $2 
  AND e.status = 'published'
  AND (e.title ILIKE $3 OR e.description ILIKE $3 OR e.tags @> ARRAY[$3]::text[])
  AND (e.is_public = true OR e.organizer = $1)
ORDER BY e.start_date ASC
LIMIT 50;

-- name: RegisterForEvent :one
INSERT INTO event_attendees (event_id, user_id)
VALUES ($1, $2)
ON CONFLICT (event_id, user_id) 
DO UPDATE SET status = 'registered', registered_at = NOW()
RETURNING *;

-- name: UnregisterFromEvent :exec
DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2;

-- name: GetEventAttendees :many
SELECT 
    ea.*,
    u.username,
    u.full_name,
    u.avatar,
    u.department,
    u.level
FROM event_attendees ea
JOIN users u ON ea.user_id = u.id
WHERE ea.event_id = $1 AND u.status = 'active'
ORDER BY ea.registered_at DESC;

-- name: AddEventCoOrganizer :one
INSERT INTO event_attendees (event_id, user_id, role)
VALUES ($1, $2, 'organizer')
ON CONFLICT (event_id, user_id) DO NOTHING
RETURNING *;

-- name: RemoveEventCoOrganizer :exec
DELETE FROM event_attendees WHERE event_id = $1 AND user_id = $2 AND role 'organizer';

-- name: GetEventCoOrganizers :many
SELECT 
    eco.*,
    u.username,
    u.full_name,
    u.avatar
FROM event_attendees eco
JOIN users u ON eco.user_id = u.id
WHERE eco.event_id = $1 AND eco.role = 'organizer' AND u.status = 'active';

-- name: UpdateEventStatus :one
UPDATE events 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: CreateAnnouncement :one
INSERT INTO announcements (
    space_id, title, content, type, target_audience, priority,
    author_id, scheduled_for, expires_at, attachments, is_pinned
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetAnnouncementByID :one
SELECT 
    a.*,
    u.username as author_username,
    u.full_name as author_full_name,
    u.avatar as author_avatar
FROM announcements a
JOIN users u ON a.author_id = u.id
WHERE a.id = $1;

-- name: ListAnnouncements :many
SELECT 
    a.*,
    u.username as author_username,
    u.full_name as author_full_name
FROM announcements a
JOIN users u ON a.author_id = u.id
WHERE a.space_id = $1 
  AND a.status = 'published'
  AND (a.scheduled_for <= NOW() OR a.scheduled_for IS NULL)
  AND (a.expires_at >= NOW() OR a.expires_at IS NULL)
  AND (a.target_audience @> $2 OR $2 IS NULL)
ORDER BY 
    a.is_pinned DESC,
    a.priority DESC,
    a.created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateAnnouncement :one
UPDATE announcements 
SET 
    title = $1,
    content = $2,
    type = $3,
    target_audience = $4,
    priority = $5,
    scheduled_for = $6,
    expires_at = $7,
    attachments = $8,
    is_pinned = $9,
    updated_at = NOW()
WHERE id = $10
RETURNING *;

-- name: UpdateAnnouncementStatus :one
UPDATE announcements
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- Admin-specific Announcement Queries

-- name: ListAllAnnouncementsAdmin :many
SELECT
    a.*,
    u.username as author_username,
    u.full_name as author_full_name,
    u.avatar as author_avatar
FROM announcements a
LEFT JOIN users u ON a.author_id = u.id
WHERE a.space_id = $1
  OR (a.status = $2 OR $2 = '')
  OR (a.priority = $3 OR $3 = '')
ORDER BY a.created_at DESC
LIMIT $4 OFFSET $5;

-- name: DeleteAnnouncement :exec
DELETE FROM announcements WHERE id = $1;

-- name: GetUserEvents :many
SELECT 
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    ea.status as attendance_status
FROM events e
JOIN users u ON e.organizer = u.id
JOIN event_attendees ea ON e.id = ea.event_id
WHERE ea.user_id = $1 AND e.space_id = $2
ORDER BY e.start_date DESC
LIMIT $3 OFFSET $4;

-- name: GetUpcomingEvents :many
SELECT 
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    (SELECT COUNT(*) FROM event_attendees ea8 WHERE ea8.event_id = e.id AND ea8.status = 'registered') as current_attendees_count
FROM events e
JOIN users u ON e.organizer = u.id
WHERE e.space_id = $1 
  AND e.status = 'published'
  AND e.start_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'
ORDER BY e.start_date ASC
LIMIT 10;

-- name: MarkEventAttendance :exec
UPDATE event_attendees 
SET status = 'attended', attended_at = NOW()
WHERE event_id = $1 AND user_id = $2;

-- name: GetEventCategories :many
SELECT DISTINCT category
FROM events
WHERE space_id = $1 AND status = 'published'
ORDER BY category;

-- Admin-specific Event Queries

-- name: ListAllEventsAdmin :many
SELECT
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') as current_attendees_count
FROM events e
LEFT JOIN users u ON e.organizer = u.id
WHERE e.space_id = $1
  OR (e.status = $2 OR $2 = '')
  OR (e.category = $3 OR $3 = '')
ORDER BY e.created_at DESC
LIMIT $4 OFFSET $5;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;

-- name: GetEventWithRegistrations :one
SELECT
    e.*,
    u.username as organizer_username,
    u.full_name as organizer_full_name,
    u.avatar as organizer_avatar,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'registered') as registered_count,
    (SELECT COUNT(*) FROM event_attendees ea WHERE ea.event_id = e.id AND ea.status = 'attended') as attended_count
FROM events e
LEFT JOIN users u ON e.organizer = u.id
WHERE e.id = $1;