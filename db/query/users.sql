-- User Management Queries

-- name: CreateUser :one
INSERT INTO users (
    space_id, username, email, password, full_name, 
    roles, level, department, major, year, interests, settings, phone_number
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetUserByID :one
SELECT 
    u.*,
    s.name as space_name,
    s.slug as space_slug
FROM users u
JOIN spaces s ON u.space_id = s.id
WHERE u.id = $1 AND u.status = 'active';

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1  AND  status= 'active';

-- name: GetUserByUsername :one
SELECT * FROM users 
WHERE username = $1  AND status = 'active';

-- name: UpdateUser :one
UPDATE users 
SET 
    full_name = $1,
    bio = $2,
    avatar = $3,
    level = $4,
    department = $5,
    major = $6,
    year = $7,
    interests = $8,
    settings = $9,
    updated_at = NOW()
WHERE id = $10 AND status = 'active'
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users 
SET password = $1, updated_at = NOW() 
WHERE id = $2;

-- name: UpdateUserLastActive :exec
UPDATE users SET last_active = NOW() WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users SET status = 'inactive', updated_at = NOW() WHERE id = $1;

-- name: SearchUsers :many
SELECT 
    id,
    username,
    full_name,
    avatar,
    bio,
    level,
    department,
    major,
    verified,
    followers_count,
    following_count
FROM users 
WHERE space_id = $1 
  AND status = 'active'
  AND (username ILIKE $2 OR full_name ILIKE $2 OR department ILIKE $2 OR bio ILIKE $2)
ORDER BY 
    CASE 
        WHEN username = $2 THEN 1
        WHEN full_name = $2 THEN 2
        ELSE 3
    END,
    followers_count DESC
LIMIT 50;

-- name: GetSuggestedUsers :many
SELECT
    u.id,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar,
    u.bio,
    u.verified,
    u.department,
    u.level,
    u.followers_count,
    u.following_count,
    COUNT(DISTINCT f.id) as mutual_followers
FROM users u
LEFT JOIN follows f ON u.id = f.following_id AND f.follower_id IN (
    SELECT following_id FROM follows fl WHERE fl.follower_id = $1
)
WHERE u.space_id = $2
  AND u.id != $1
  AND u.status = 'active'
  AND NOT EXISTS (SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = u.id)
  AND (
    u.department = (SELECT department FROM users u2 WHERE u2.id = $1 AND department IS NOT NULL)
    OR u.level = (SELECT level FROM users u3 WHERE u3.id = $1 AND level IS NOT NULL)
    OR (u.interests && (SELECT interests FROM users u4 WHERE u4.id = $1) AND array_length(u.interests, 1) > 0)
    OR ((SELECT department FROM users u5 WHERE u5.id = $1) IS NULL AND (SELECT level FROM users u6 WHERE u6.id = $1) IS NULL)
  )
GROUP BY u.id
ORDER BY mutual_followers DESC, u.followers_count DESC
LIMIT $3 OFFSET $4;

-- name: GetUserStats :one
SELECT 
    COUNT(*) as total_posts,
    COUNT(DISTINCT flwrs.follower_id) as total_followers,
    COUNT(DISTINCT flwng.following_id) as total_following,
    COUNT(DISTINCT gm.group_id) as total_groups,
    COUNT(DISTINCT cm.community_id) as total_communities,
    COUNT(DISTINCT ts.id) as total_tutoring_sessions,
    COUNT(DISTINCT ms.id) as total_mentoring_sessions
FROM users u
LEFT JOIN posts p ON u.id = p.author_id AND p.status = 'active'
LEFT JOIN follows flwrs ON u.id = flwrs.following_id
LEFT JOIN follows flwng ON u.id = flwng.follower_id
LEFT JOIN group_members gm ON u.id = gm.user_id
LEFT JOIN community_members cm ON u.id = cm.user_id
LEFT JOIN tutoring_sessions ts ON u.id = ts.tutor_id OR u.id = ts.student_id
LEFT JOIN mentoring_sessions ms ON u.id = ms.mentor_id OR u.id = ms.mentee_id
WHERE u.id = $1
GROUP BY u.id;

-- name: GetUsersByRole :many
SELECT 
    id,
    username,
    full_name,
    avatar,
    level,
    department,
    verified,
    mentor_status,
    tutor_status
FROM users 
WHERE space_id = $1 
  AND status = 'active'
  AND $2 = ANY(roles)
ORDER BY full_name;

-- name: UpdateMentorStatus :one
UPDATE users 
SET mentor_status = $1, updated_at = NOW() 
WHERE id = $2 
RETURNING *;

-- name: UpdateTutorStatus :one
UPDATE users 
SET tutor_status = $1, updated_at = NOW() 
WHERE id = $2 
RETURNING *;

-- name: GetUsersWithPendingApplications :many
SELECT 
    u.*,
    ta.id as tutor_application_id,
    ma.id as mentor_application_id
FROM users u
LEFT JOIN tutor_applications ta ON u.id = ta.applicant_id AND ta.status = 'pending'
LEFT JOIN mentor_applications ma ON u.id = ma.applicant_id AND ma.status = 'pending'
WHERE u.space_id = $1 
  AND (ta.id IS NOT NULL OR ma.id IS NOT NULL)
  AND u.status = 'active';

-- -- name: GetUserNotifications :many
-- SELECT *
-- FROM notifications 
-- WHERE to_user_id = $1 
-- ORDER BY created_at DESC 
-- LIMIT $2 OFFSET $3;

-- name: MarkNotificationsAsRead :exec
UPDATE notifications 
SET is_read = true, read_at = NOW() 
WHERE to_user_id = $1 AND is_read = false;

-- name: GetUnreadNotificationCount :one
SELECT COUNT(*) FROM notifications
WHERE to_user_id = $1 AND is_read = false;

-- Note: User preferences are now stored in users.settings JSONB field
-- Note: User activities are now tracked via user_sessions table

-- name: GetUserSettings :one
SELECT settings FROM users
WHERE id = $1;

-- name: UpdateUserSettings :one
UPDATE users
SET settings = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserSessionActivity :many
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY last_activity DESC
LIMIT $2 OFFSET $3;

-- name: AdvancedSearchUsers :many
SELECT 
    u.*,
    ts_rank_cd(to_tsvector('english', 
        COALESCE(u.username, '') || ' ' || 
        COALESCE(u.full_name, '') || ' ' || 
        COALESCE(u.bio, '') || ' ' || 
        COALESCE(u.department, '') || ' ' || 
        COALESCE(u.major, '')
    ), plainto_tsquery('english', $1)) as relevance_score
FROM users u
WHERE u.space_id = $2 
  AND u.status = 'active'
  AND (
    to_tsvector('english', 
        COALESCE(u.username, '') || ' ' || 
        COALESCE(u.full_name, '') || ' ' || 
        COALESCE(u.bio, '') || ' ' || 
        COALESCE(u.department, '') || ' ' || 
        COALESCE(u.major, '')
    ) @@ plainto_tsquery('english', $1)
    OR u.username ILIKE '%' || $1 || '%'
    OR u.full_name ILIKE '%' || $1 || '%'
  )
ORDER BY relevance_score DESC, u.followers_count DESC
LIMIT $3 OFFSET $4;

-- name: GetUserEngagementAnalytics :one
SELECT
    COUNT(DISTINCT p.id) as posts_created,
    COUNT(DISTINCT l.id) as likes_given,
    (SELECT COUNT(*) FROM likes l2
     JOIN posts p2 ON l2.post_id = p2.id
     WHERE p2.author_id = $1) as likes_received,
    COUNT(DISTINCT c.id) as comments_made,
    COUNT(DISTINCT cm.community_id) as communities_joined,
    COUNT(DISTINCT gm.group_id) as groups_joined,
    COUNT(DISTINCT ea.event_id) as events_attended,
    COUNT(DISTINCT ts.id) + COUNT(DISTINCT ms.id) as sessions_completed
FROM users u
LEFT JOIN posts p ON u.id = p.author_id AND p.status = 'active'
LEFT JOIN likes l ON u.id = l.user_id
LEFT JOIN comments c ON u.id = c.author_id AND c.status = 'active'
LEFT JOIN community_members cm ON u.id = cm.user_id
LEFT JOIN group_members gm ON u.id = gm.user_id
LEFT JOIN event_attendees ea ON u.id = ea.user_id AND ea.status = 'attended'
LEFT JOIN tutoring_sessions ts ON (u.id = ts.tutor_id OR u.id = ts.student_id) AND ts.status = 'completed'
LEFT JOIN mentoring_sessions ms ON (u.id = ms.mentor_id OR u.id = ms.mentee_id) AND ms.status = 'completed'
WHERE u.id = $1
GROUP BY u.id;

-- Follow System Queries

-- name: FollowUser :one
INSERT INTO follows (follower_id, following_id, space_id)
VALUES ($1, $2, $3)
ON CONFLICT (follower_id, following_id) DO NOTHING
RETURNING *;

-- name: UnfollowUser :exec
DELETE FROM follows
WHERE follower_id = $1 AND following_id = $2;

-- name: CheckIfFollowing :one
SELECT EXISTS(
    SELECT 1 FROM follows
    WHERE follower_id = $1 AND following_id = $2
) as is_following;

-- name: GetUserFollowers :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    u.bio,
    u.verified,
    u.followers_count,
    u.following_count,
    f.created_at as followed_at
FROM follows f
JOIN users u ON f.follower_id = u.id
WHERE f.following_id = $1
  AND u.status = 'active'
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserFollowing :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    u.bio,
    u.verified,
    u.followers_count,
    u.following_count,
    f.created_at as followed_at
FROM follows f
JOIN users u ON f.following_id = u.id
WHERE f.follower_id = $1
  AND u.status = 'active'
ORDER BY f.created_at DESC
LIMIT $2 OFFSET $3;

-- name: IncrementFollowersCount :exec
UPDATE users
SET followers_count = followers_count + 1,
    updated_at = NOW()
WHERE id = $1;

-- name: DecrementFollowersCount :exec
UPDATE users
SET followers_count = GREATEST(0, followers_count - 1),
    updated_at = NOW()
WHERE id = $1;

-- name: IncrementFollowingCount :exec
UPDATE users
SET following_count = following_count + 1,
    updated_at = NOW()
WHERE id = $1;

-- name: DecrementFollowingCount :exec
UPDATE users
SET following_count = GREATEST(0, following_count - 1),
    updated_at = NOW()
WHERE id = $1;

-- Admin User Management Queries

-- name: ListUsers :many
SELECT
    id,
    username,
    full_name,
    email,
    avatar,
    status,
    created_at,
    roles, 
    department
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;