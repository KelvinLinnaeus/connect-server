-- Admin User Management Queries
-- Note: Admin users are managed via users table with roles field

-- name: GetAllAdminUsers :many
SELECT
    u.id,
    u.username,
    u.email,
    u.full_name,
    u.avatar,
    u.roles,
    u.status,
    u.created_at,
    u.updated_at
FROM users u
WHERE
    (u.roles && ARRAY['admin', 'super_admin']::text[])
    AND (u.status = $1 OR $1 = '')
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUserRole :one
UPDATE users
SET
    roles = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CheckAdminPermission :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE id = $1
    AND status = 'active'
    AND roles && ARRAY['admin', 'super_admin']::text[]
) as has_permission;

-- name: IsUserSuperAdmin :one
SELECT EXISTS(
    SELECT 1 FROM users
    WHERE id = $1
    AND status = 'active'
    AND 'super_admin' = ANY(roles)
) as is_super_admin;

-- User Suspension Queries

-- name: CreateUserSuspension :one
INSERT INTO user_suspensions (
    user_id,
    suspended_by,
    reason,
    notes,
    suspended_until,
    is_permanent
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetActiveSuspension :one
SELECT * FROM user_suspensions
WHERE user_id = $1
AND (suspended_until IS NULL OR suspended_until > NOW())
ORDER BY created_at DESC
LIMIT 1;

-- name: GetUserSuspensions :many
SELECT
    us.*,
    u.username,
    u.full_name,
    admin.username as suspended_by_username
FROM user_suspensions us
JOIN users u ON us.user_id = u.id
JOIN users admin ON us.suspended_by = admin.id
WHERE us.user_id = $1
ORDER BY us.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateUserAccountStatus :exec
UPDATE users
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: LiftSuspension :exec
UPDATE users
SET
    status = 'active',
    suspended_until = NULL,
    updated_at = NOW()
WHERE id = $1;

-- Content Reports Queries
-- Note: Uses existing 'reports' table from migration 008

-- name: CreateContentReport :one
INSERT INTO reports (
    space_id,
    reporter_id,
    content_type,
    content_id,
    reason,
    description,
    status,
    priority
) VALUES ($1, $2, $3, $4, $5, $6, 'pending', 'medium')
RETURNING *;

-- name: GetContentReports :many
SELECT
    r.*,
    u.username as reporter_username,
    u.full_name as reporter_name,
    reviewer.username as reviewer_username
FROM reports r
LEFT JOIN users u ON r.reporter_id = u.id
LEFT JOIN users reviewer ON r.reviewed_by = reviewer.id
WHERE
    (r.space_id = $1 OR $1 = '00000000-0000-0000-0000-000000000000'::uuid)
    AND (r.status = $2 OR $2 = '')
    AND (r.content_type = $3 OR $3 = '')
ORDER BY r.created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetContentReportByID :one
SELECT
    r.*,
    u.username as reporter_username,
    u.full_name as reporter_name
FROM reports r
LEFT JOIN users u ON r.reporter_id = u.id
WHERE r.id = $1;

-- name: UpdateContentReportStatus :one
UPDATE reports
SET
    status = $2,
    reviewed_by = $3,
    moderation_notes = $4,
    reviewed_at = CASE WHEN $2 IN ('resolved', 'dismissed') THEN NOW() ELSE reviewed_at END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetReportStats :one
SELECT
    COUNT(*) FILTER (WHERE status = 'pending') as pending_count,
    COUNT(*) FILTER (WHERE status = 'resolved') as resolved_count,
    COUNT(*) FILTER (WHERE status = 'escalated') as escalated_count,
    COUNT(*) FILTER (WHERE status = 'dismissed') as dismissed_count
FROM reports
WHERE space_id = $1;

-- System Settings Queries

-- name: GetSystemSetting :one
SELECT * FROM system_settings
WHERE key = $1
LIMIT 1;

-- name: GetAllSystemSettings :many
SELECT * FROM system_settings
ORDER BY key;

-- name: UpsertSystemSetting :one
INSERT INTO system_settings (key, value, description, updated_by)
VALUES ($1, $2, $3, $4)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING *;

-- name: DeleteSystemSetting :exec
DELETE FROM system_settings
WHERE key = $1;

-- Audit Logs Queries

-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    admin_user_id,
    action,
    resource_type,
    resource_id,
    details,
    ip_address,
    user_agent
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAuditLogs :many
SELECT
    al.*,
    u.username,
    u.full_name
FROM audit_logs al
LEFT JOIN users u ON al.admin_user_id = u.id
WHERE
    (al.admin_user_id = $1 OR $1 = '00000000-0000-0000-0000-000000000000'::uuid)
    AND (al.resource_type = $2 OR $2 = '')
    AND (al.created_at >= $3 OR $3 IS NULL)
    AND (al.created_at <= $4 OR $4 IS NULL)
ORDER BY al.created_at DESC
LIMIT $5 OFFSET $6;

-- Space Activities Queries

-- name: CreateSpaceActivity :one
INSERT INTO space_activities (
    space_id,
    activity_type,
    actor_id,
    actor_name,
    description,
    metadata
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSpaceActivities :many
SELECT
    sa.*,
    u.username,
    u.avatar
FROM space_activities sa
LEFT JOIN users u ON sa.actor_id = u.id
WHERE
    sa.space_id = $1
    AND (sa.activity_type = $2 OR $2 = '')
    AND (sa.created_at >= $3 OR $3 IS NULL)
ORDER BY sa.created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetActivityStats :one
SELECT
    COUNT(*) FILTER (WHERE activity_type = 'user_joined') as users_joined,
    COUNT(*) FILTER (WHERE activity_type = 'post_created') as posts_created,
    COUNT(*) FILTER (WHERE activity_type = 'community_created') as communities_created,
    COUNT(*) FILTER (WHERE activity_type = 'group_created') as groups_created,
    COUNT(*) FILTER (WHERE activity_type = 'event_created') as events_created
FROM space_activities
WHERE space_id = $1
AND created_at >= $2;

-- Admin Dashboard Analytics Queries

-- name: GetAdminDashboardStats :one
-- name: GetAdminStats :one
SELECT
    (SELECT COUNT(*) FROM users u WHERE u.space_id = sqlc.arg(space_id)) AS total_users,
    (SELECT COUNT(*) FROM users u WHERE u.space_id = sqlc.arg(space_id) AND u.created_at >= NOW() - INTERVAL '30 days') AS new_users_month,
    (SELECT COUNT(*) FROM posts p WHERE p.space_id = sqlc.arg(space_id) AND p.status = 'active') AS total_posts,
    (SELECT COUNT(*) FROM communities c WHERE c.space_id = sqlc.arg(space_id)) AS total_communities,
    (SELECT COUNT(*) FROM groups g WHERE g.space_id = sqlc.arg(space_id)) AS total_groups,
    (SELECT COUNT(*) FROM reports r WHERE r.space_id = sqlc.arg(space_id) AND r.status = 'pending') AS pending_reports,
    (SELECT COUNT(*) FROM user_suspensions us JOIN users u ON us.user_id = u.id WHERE u.space_id = sqlc.arg(space_id) AND us.created_at >= NOW() - INTERVAL '30 days') AS suspensions_month;


-- name: GetUserGrowthData :many
SELECT
    DATE(created_at) as date,
    COUNT(*) as new_users
FROM users
WHERE space_id = $1
AND created_at >= $2
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- name: GetContentGrowthData :many
SELECT
    DATE(created_at) as date,
    COUNT(*) FILTER (WHERE status = 'active') as active_posts,
    COUNT(*) as total_posts
FROM posts
WHERE space_id = $1
AND created_at >= $2
GROUP BY DATE(created_at)
ORDER BY date DESC;

-- User Management Queries (Admin variant with extended filtering)

-- name: SearchUsersAdmin :many
SELECT
    u.id,
    u.username,
    u.email,
    u.full_name,
    u.avatar,
    u.verified,
    u.roles,
    u.status,
    u.suspended_until,
    u.created_at,
    u.followers_count,
    u.following_count,
    (SELECT COUNT(*) FROM posts WHERE author_id = u.id AND status = 'active') as posts_count
FROM users u
WHERE
    u.space_id = $1
    AND (
        u.full_name ILIKE '%' || $2 || '%'
        OR u.username ILIKE '%' || $2 || '%'
        OR u.email ILIKE '%' || $2 || '%'
    )
    AND (u.status = $3 OR $3 = '')
    AND (u.roles && $4 OR array_length($4::text[], 1) IS NULL)
ORDER BY u.created_at DESC
LIMIT $5 OFFSET $6;

-- name: GetUserDetails :one
SELECT
    u.*,
    (SELECT COUNT(*) FROM posts WHERE author_id = u.id AND status = 'active') as posts_count,
    (SELECT COUNT(*) FROM likes WHERE user_id = u.id) as likes_given,
    (SELECT COUNT(*) FROM comments WHERE author_id = u.id AND status = 'active') as comments_count
FROM users u
WHERE u.id = $1;

-- Space Management Queries

-- name: GetAllSpacesWithStats :many
SELECT
    s.*,
    (SELECT COUNT(*) FROM users WHERE space_id = s.id) as total_users,
    (SELECT COUNT(*) FROM posts WHERE space_id = s.id AND status = 'active') as total_posts,
    (SELECT COUNT(*) FROM communities WHERE space_id = s.id) as total_communities,
    (SELECT COUNT(*) FROM groups WHERE space_id = s.id) as total_groups
FROM spaces s
ORDER BY s.created_at DESC;

-- name: GetSpaceWithStats :one
SELECT
    s.*,
    (SELECT COUNT(*) FROM users WHERE space_id = s.id) as total_users,
    (SELECT COUNT(*) FROM posts WHERE space_id = s.id AND status = 'active') as total_posts,
    (SELECT COUNT(*) FROM communities WHERE space_id = s.id) as total_communities,
    (SELECT COUNT(*) FROM groups WHERE space_id = s.id) as total_groups,
    (SELECT COUNT(*) FROM users WHERE space_id = s.id AND created_at >= NOW() - INTERVAL '7 days') as new_users_week
FROM spaces s
WHERE s.id = $1;

-- Additional Admin User Management Queries

-- name: AdminCreateUser :one
INSERT INTO users (
    space_id, username, email, password, full_name, roles, status
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AdminUpdateUser :one
UPDATE users
SET
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    email = COALESCE(sqlc.narg(email), email),
    roles = COALESCE(sqlc.narg(roles), roles),
    status = COALESCE(sqlc.narg(status), status),
    department = COALESCE(sqlc.narg(department), department),
    level = COALESCE(sqlc.narg(level), level),
    verified = COALESCE(sqlc.narg(verified), verified),
    updated_at = NOW()
WHERE id = sqlc.arg(user_id)
RETURNING *;

-- name: ResetUserPassword :exec
UPDATE users
SET
    password = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateContentReportWithAction :one
UPDATE reports
SET
    status = $2,
    reviewed_by = $3,
    moderation_notes = $4,
    actions_taken = $5,
    reviewed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateContentReportPriority :one
UPDATE reports
SET
    priority = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
