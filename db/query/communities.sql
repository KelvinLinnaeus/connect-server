-- Communities Management Queries

-- name: CreateCommunity :one
INSERT INTO communities (
    space_id, name, description, category, cover_image, is_public, created_by, settings
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetCommunityByID :one
SELECT 
    c.*,
    u.username as created_by_username,
    u.full_name as created_by_full_name,
    EXISTS(SELECT 1 FROM community_members cm2 WHERE cm2.community_id = c.id AND cm2.user_id = $1) as is_member,
    cm.role as user_role,
    (SELECT COUNT(*) FROM community_members cm3 WHERE cm3.community_id = c.id) as actual_member_count,
    (SELECT COUNT(*) FROM posts WHERE community_id = c.id AND status = 'active') as actual_post_count
FROM communities c
JOIN users u ON c.created_by = u.id
LEFT JOIN community_members cm ON c.id = cm.community_id AND cm.user_id = $1
WHERE c.id = $2;

-- name: GetCommunityBySlug :one
SELECT 
    c.*,
    u.username as created_by_username,
    u.full_name as created_by_full_name,
    EXISTS(SELECT 1 FROM community_members cm2 WHERE cm2.community_id = c.id AND cm2.user_id = $1) as is_member,
    cm.role as user_role
FROM communities c
JOIN users u ON c.created_by = u.id
LEFT JOIN community_members cm ON c.id = cm.community_id AND cm.user_id = $1
WHERE c.space_id = $2 AND LOWER(c.name) = LOWER($3);

-- name: UpdateCommunity :one
UPDATE communities 
SET 
    name = $1,
    description = $2,
    cover_image = $3,
    category = $4,
    is_public = $5,
    settings = $6,
    updated_at = NOW()
WHERE id = $7
RETURNING *;

-- name: ListCommunities :many
SELECT 
    c.*,
    cm.role as user_role,
    EXISTS(SELECT 1 FROM community_members cm2 WHERE cm2.community_id = c.id AND cm2.user_id = $1) as is_member,
    (SELECT COUNT(*) FROM community_members cm3 WHERE cm3.community_id = c.id) as actual_member_count
FROM communities c
LEFT JOIN community_members cm ON c.id = cm.community_id AND cm.user_id = $1
WHERE c.space_id = $2
ORDER BY 
    CASE WHEN $3 = 'members' THEN c.member_count END DESC,
    CASE WHEN $3 = 'posts' THEN c.post_count END DESC,
    CASE WHEN $3 = 'recent' THEN c.created_at END DESC,
    c.name ASC
LIMIT $4 OFFSET $5;

-- name: SearchCommunities :many
SELECT 
    c.*,
    EXISTS(SELECT 1 FROM community_members cm2 WHERE cm2.community_id = c.id AND cm2.user_id = $1) as is_member,
    (SELECT COUNT(*) FROM community_members cm3 WHERE cm3.community_id = c.id) as actual_member_count
FROM communities c
WHERE c.space_id = $2
  AND (c.name ILIKE $3 OR c.description ILIKE $3 OR c.category ILIKE $3)
  AND (c.is_public = true OR EXISTS(SELECT 1 FROM community_members cm4 WHERE cm4.community_id = c.id AND cm4.user_id = $1))
ORDER BY c.member_count DESC
LIMIT 50;

-- name: JoinCommunity :one
INSERT INTO community_members (community_id, user_id, role)
VALUES ($1, $2, 'member')
ON CONFLICT (community_id, user_id) 
DO UPDATE SET role = 'member', joined_at = NOW()
RETURNING *;

-- name: LeaveCommunity :exec
DELETE FROM community_members 
WHERE community_id = $1 AND user_id = $2;

-- name: GetCommunityMembers :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    u.level,
    u.department,
    u.verified,
    cm.role,
    cm.joined_at
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.community_id = $1 AND u.status = 'active'
ORDER BY
    cm.joined_at;




-- name: AddCommunityModerator :one
INSERT INTO community_members (community_id, user_id, permissions, role)
VALUES ($1, $2, $3, 'moderator')
RETURNING *;

-- name: RemoveCommunityModerator :exec
DELETE FROM community_members WHERE community_id = $1 AND user_id = $2 AND role = 'moderator';

-- name: IsCommunityAdmin :one
SELECT EXISTS(
    SELECT 1 FROM community_members cm 
    WHERE cm.community_id = $1 AND cm.user_id = $2 AND cm.role = 'admin'
) as is_admin;

-- name: IsCommunityModerator :one
SELECT EXISTS(
    SELECT 1 FROM community_members cm 
    WHERE cm.community_id = $1 AND cm.user_id = $2 AND role = 'moderator'
) as is_moderator;

-- name: GetCommunityAdmins :many
SELECT
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    ca.permissions
FROM community_members ca
JOIN users u ON ca.user_id = u.id
WHERE ca.community_id = $1 AND ca.role = 'admin' AND u.status = 'active';

-- name: GetCommunityModerators :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    cm.permissions
FROM community_members cm
JOIN users u ON cm.user_id = u.id
WHERE cm.community_id = $1 AND u.status = 'active' AND role='moderator';

-- name: UpdateCommunityStats :exec
UPDATE communities 
SET 
    member_count = (SELECT COUNT(*) FROM community_members cm5 WHERE cm5.community_id = $1),
    post_count = (SELECT COUNT(*) FROM posts WHERE community_id = $1 AND status = 'active'),
    updated_at = NOW()
WHERE id = $1;

-- name: GetCommunityCategories :many
SELECT DISTINCT category 
FROM communities 
WHERE space_id = $1 
ORDER BY category;

-- name: GetUserCommunities :many
SELECT
    c.*,
    cm.role as user_role,
    cm.joined_at
FROM communities c
JOIN community_members cm ON c.id = cm.community_id
WHERE cm.user_id = $1 AND c.space_id = $2
ORDER BY cm.joined_at DESC;

-- Admin-specific Community Queries

-- name: ListAllCommunitiesAdmin :many
SELECT
    c.*,
    u.username as created_by_username,
    u.full_name as created_by_full_name,
    (SELECT COUNT(*) FROM community_members cm WHERE cm.community_id = c.id) as actual_member_count,
    (SELECT COUNT(*) FROM posts WHERE community_id = c.id AND status = 'active') as actual_post_count
FROM communities c
LEFT JOIN users u ON c.created_by = u.id
WHERE c.space_id = $1
  OR (c.category = $2 OR $2 = '')
  OR (c.status = $3 OR $3 = '')
ORDER BY c.created_at DESC
LIMIT $4 OFFSET $5;

-- name: DeleteCommunity :exec
DELETE FROM communities WHERE id = $1;

-- name: UpdateCommunityStatus :one
UPDATE communities
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;