-- Groups Management Queries

-- name: CreateGroup :one
INSERT INTO groups (
    space_id, community_id, name, description, category, group_type,
    avatar, banner, allow_invites, allow_member_posts,
     created_by, tags, settings
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetGroupByID :one
SELECT 
    g.*,
    c.name as community_name,
    u.username as created_by_username,
    u.full_name as created_by_full_name,
    EXISTS(SELECT 1 FROM group_members gm2 WHERE gm2.group_id = g.id AND gm2.user_id = $1) as is_member,
    gm.role as user_role,
    (SELECT COUNT(*) FROM group_members gm3 WHERE gm3.group_id = g.id) as actual_member_count,
    (SELECT COUNT(*) FROM posts WHERE group_id = g.id AND status = 'active') as actual_post_count
FROM groups g
LEFT JOIN communities c ON g.community_id = c.id
JOIN users u ON g.created_by = u.id
LEFT JOIN group_members gm ON g.id = gm.group_id AND gm.user_id = $1
WHERE g.id = $2 AND g.status = 'active';

-- name: UpdateGroup :one
UPDATE groups 
SET 
    name = $1,
    description = $2,
    category = $3,
    avatar = $4,
    banner = $5,
    allow_invites = $6,
    allow_member_posts = $7,
    tags = $8,
    settings = $9,
    updated_at = NOW()
WHERE id = $10
RETURNING *;

-- name: ListGroups :many
SELECT 
    g.*,
    c.name as community_name,
    EXISTS(SELECT 1 FROM group_members gm2 WHERE gm2.group_id = g.id AND gm2.user_id = $1) as is_member,
    gm.role as user_role,
    (SELECT COUNT(*) FROM group_members gm3 WHERE gm3.group_id = g.id) as actual_member_count
FROM groups g
LEFT JOIN communities c ON g.community_id = c.id
LEFT JOIN group_members gm ON g.id = gm.group_id AND gm.user_id = $1
WHERE g.space_id = $2 AND g.status = 'active'
  AND (g.visibility = 'public' OR EXISTS(SELECT 1 FROM group_members gm4 WHERE gm4.group_id = g.id AND gm4.user_id = $1))
ORDER BY 
    CASE WHEN $3 = 'members' THEN g.member_count END DESC,
    CASE WHEN $3 = 'recent' THEN g.created_at END DESC,
    g.name ASC
LIMIT $4 OFFSET $5;

-- name: SearchGroups :many
SELECT 
    g.*,
    c.name as community_name,
    EXISTS(SELECT 1 FROM group_members gm2 WHERE gm2.group_id = g.id AND gm2.user_id = $1) as is_member,
    (SELECT COUNT(*) FROM group_members gm3 WHERE gm3.group_id = g.id) as actual_member_count
FROM groups g
LEFT JOIN communities c ON g.community_id = c.id
WHERE g.space_id = $2
  AND g.status = 'active'
  AND (g.name ILIKE $3 OR g.description ILIKE $3 OR g.tags @> ARRAY[$3]::text[])
  AND (g.visibility = 'public' OR EXISTS(SELECT 1 FROM group_members gm4 WHERE gm4.group_id = g.id AND gm4.user_id = $1))
ORDER BY g.member_count DESC
LIMIT 50;

-- name: JoinGroup :one
INSERT INTO group_members (group_id, user_id, role, invited_by)
VALUES ($1, $2, 'member', $3)
ON CONFLICT (group_id, user_id) 
DO UPDATE SET role = 'member', joined_at = NOW()
RETURNING *;

-- name: LeaveGroup :exec
DELETE FROM group_members WHERE group_id = $1 AND user_id = $2;

-- name: GetGroupJoinRequests :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    u.level,
    u.department,
    u.verified,
    gm.role,
    gm.joined_at,
    gm.permissions
FROM group_members gm
JOIN users u ON gm.user_id = u.id
WHERE gm.group_id = $1 AND u.status = 'active'
ORDER BY
    gm.joined_at;

-- name: AddGroupAdmin :one
INSERT INTO group_members (group_id, user_id, permissions, role)
VALUES ($1, $2, $3, 'admin')
ON CONFLICT (group_id, user_id) 
DO UPDATE SET permissions = $4, assigned_at = NOW()
RETURNING *;

-- name: RemoveGroupAdmin :exec
DELETE FROM group_members WHERE group_id = $1 AND user_id = $2 AND role='admin';

-- name: AddGroupModerator :one
INSERT INTO group_members (group_id, user_id, permissions, role)
VALUES ($1, $2, $3, 'moderator')
ON CONFLICT (group_id, user_id) 
DO UPDATE SET permissions = $4, assigned_at = NOW()
RETURNING *;

-- name: RemoveGroupModerator :exec
DELETE FROM group_members WHERE group_id = $1 AND user_id = $2 AND role 'moderator';

-- name: IsGroupAdmin :one
SELECT EXISTS(
    SELECT 1 FROM group_members ga 
    WHERE ga.group_id = $1 AND ga.user_id = $2 AND role 'admin'
) as is_admin;

-- name: IsGroupModerator :one
SELECT EXISTS(
    SELECT 1 FROM group_members gm 
    WHERE gm.group_id = $1 AND gm.user_id = $2 AND role = 'moderator'
) as is_moderator;

-- name: UpdateGroupMemberRole :exec
UPDATE group_members 
SET role = $1, permissions = $2 
WHERE group_id = $3 AND user_id = $4;

-- name: GetUserGroups :many
SELECT 
    g.*,
    c.name as community_name,
    gm.role as user_role,
    gm.joined_at
FROM groups g
LEFT JOIN communities c ON g.community_id = c.id
JOIN group_members gm ON g.id = gm.group_id
WHERE gm.user_id = $1 AND g.space_id = $2 AND g.status = 'active'
ORDER BY gm.joined_at DESC;

-- name: UpdateGroupStats :exec
UPDATE groups 
SET 
    member_count = (SELECT COUNT(*) FROM group_members gm7 WHERE gm7.group_id = $1),
    post_count = (SELECT COUNT(*) FROM posts WHERE group_id = $1 AND status = 'active'),
    updated_at = NOW()
WHERE id = $1;

-- name: CreateProjectRole :one
INSERT INTO group_roles (group_id, name, description, slots_total, requirements, skills_required)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProjectRoles :many
SELECT * FROM group_roles WHERE group_id = $1 ORDER BY created_at DESC;

-- name: ApplyForProjectRole :one
INSERT INTO group_applications (role_id, user_id, message)
VALUES ($1, $2, $3)
ON CONFLICT (role_id, user_id) 
DO UPDATE SET message = $3, applied_at = NOW()
RETURNING *;

-- name: GetRoleApplications :many
SELECT
    ra.*,
    u.username,
    u.full_name,
    u.avatar,
    pr.name as role_name
FROM group_applications ra
JOIN users u ON ra.user_id = u.id
JOIN group_roles pr ON ra.role_id = pr.id
WHERE pr.group_id = $1 AND ra.status = 'pending'
ORDER BY ra.applied_at DESC;

-- Admin Group Management Queries

-- name: GetGroupsBySpaceID :many
SELECT * FROM groups
WHERE space_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetGroupsByStatus :many
SELECT * FROM groups
WHERE space_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateGroupStatus :exec
UPDATE groups
SET status = $1, updated_at = NOW()
WHERE id = $2;

-- name: DeleteGroup :exec
DELETE FROM groups WHERE id = $1;