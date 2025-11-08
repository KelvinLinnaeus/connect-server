-- Posts and Social Interactions Queries

-- name: CreatePost :one
INSERT INTO posts (
    author_id, space_id, community_id, group_id, parent_post_id, quoted_post_id,
    content, media, tags, visibility
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetPostByID :one
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    u.verified as author_verified,
    c.name as community_name,
    c.id as community_id,
    g.name as group_name,
    g.id as group_id,
    qp.content as quoted_content,
    qp.author_id as quoted_author_id,
    COALESCE(NULLIF(qu.username, ''), 'user_' || SUBSTRING(qu.id::text, 1, 8)) as quoted_username,
    COALESCE(NULLIF(qu.full_name, ''), 'User') as quoted_full_name,
    EXISTS(SELECT 1 FROM likes l2 WHERE l2.post_id = p.id AND l2.user_id = $1) as is_liked,
    EXISTS(SELECT 1 FROM posts pq2 WHERE pq2.quoted_post_id = p.id AND pq2.author_id = $1) as is_quoted,
    (SELECT COUNT(*) FROM likes l3 WHERE l3.post_id = p.id) as actual_likes_count,
    (SELECT COUNT(*) FROM comments c2 WHERE c2.post_id = p.id AND c2.status = 'active') as actual_comments_count,
    (SELECT COUNT(*) FROM posts p2 WHERE p2.quoted_post_id = p.id AND p2.status = 'active') as actual_quotes_count
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN communities c ON p.community_id = c.id
LEFT JOIN groups g ON p.group_id = g.id
LEFT JOIN posts qp ON p.quoted_post_id = qp.id
LEFT JOIN users qu ON qp.author_id = qu.id
WHERE p.id = $2 AND p.status = 'active';

-- name: GetUserFeed :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    u.verified as author_verified,
    c.name as community_name,
    g.name as group_name,
    EXISTS(SELECT 1 FROM likes l4 WHERE l4.post_id = p.id AND l4.user_id = $1) as is_liked,
    EXISTS(SELECT 1 FROM posts pq3 WHERE pq3.quoted_post_id = p.id AND pq3.author_id = $1) as is_quoted
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN communities c ON p.community_id = c.id
LEFT JOIN groups g ON p.group_id = g.id
WHERE p.space_id = $2
  AND p.status = 'active'
  AND (p.visibility = 'public'
       OR p.author_id = $1
       OR p.author_id IN (SELECT following_id FROM follows WHERE follower_id = $1)
       OR p.community_id IN (SELECT community_id FROM community_members WHERE user_id = $1)
       OR p.group_id IN (SELECT group_id FROM group_members WHERE user_id = $1))
ORDER BY p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetUserPosts :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    c.name as community_name,
    g.name as group_name,
    EXISTS(SELECT 1 FROM likes l5 WHERE l5.post_id = p.id AND l5.user_id = $1) as is_liked
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN communities c ON p.community_id = c.id
LEFT JOIN groups g ON p.group_id = g.id
WHERE p.author_id = $1 AND p.status = 'active'
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetCommunityPosts :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    EXISTS(SELECT 1 FROM likes l6 WHERE l6.post_id = p.id AND l6.user_id = $1) as is_liked
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.community_id = $2 AND p.status = 'active'
ORDER BY p.is_pinned DESC, p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetGroupPosts :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    EXISTS(SELECT 1 FROM likes l7 WHERE l7.post_id = p.id AND l7.user_id = $1) as is_liked
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.group_id = $2 AND p.status = 'active'
ORDER BY p.is_pinned DESC, p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetTrendingPosts :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    (p.likes_count + (p.comments_count * 2) + p.views_count + (p.quotes_count * 3)) as engagement_score
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.space_id = $1
  AND p.status = 'active'
  AND p.created_at >= NOW() - INTERVAL '7 days'
ORDER BY engagement_score DESC, p.created_at DESC
LIMIT 20;

-- name: SearchPosts :many
SELECT
    p.*,
    COALESCE(NULLIF(u.username, ''), 'user_' || SUBSTRING(u.id::text, 1, 8)) as username,
    COALESCE(NULLIF(u.full_name, ''), 'User') as full_name,
    u.avatar as author_avatar,
    c.name as community_name,
    g.name as group_name,
    ts_rank_cd(to_tsvector('english', p.content), plainto_tsquery('english', $2)) as rank
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN communities c ON p.community_id = c.id
LEFT JOIN groups g ON p.group_id = g.id
WHERE p.space_id = $1
  AND p.status = 'active'
  AND (p.content ILIKE $3 OR p.tags @> ARRAY[$2]::text[] OR to_tsvector('english', p.content) @@ plainto_tsquery('english', $2))
ORDER BY rank DESC, p.created_at DESC
LIMIT $4 OFFSET $5;

-- name: IncrementPostViews :exec
UPDATE posts SET views_count = views_count + 1 WHERE id = $1;

-- name: TogglePostLike :one
WITH like_action AS (
    INSERT INTO likes (user_id, post_id) 
    VALUES ($1, $2)
    ON CONFLICT (user_id, post_id) 
    DO NOTHING
    RETURNING true as liked
)
UPDATE posts 
SET likes_count = likes_count + 
    CASE 
        WHEN EXISTS (SELECT 1 FROM like_action) THEN 1 
        ELSE -1 
    END
WHERE posts.id = $2
RETURNING likes_count;

-- name: CreateComment :one
INSERT INTO comments (post_id, author_id, parent_comment_id, content)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPostComments :many
WITH RECURSIVE comment_tree AS (
    SELECT 
        c.*, 
        u.username,
        u.full_name,
        u.avatar,
        0 as depth,
        ARRAY[c.id] as path
    FROM comments c
    JOIN users u ON c.author_id = u.id
    WHERE c.post_id = $1 AND c.parent_comment_id IS NULL AND c.status = 'active'
    
    UNION ALL
    
    SELECT 
        c.*,
        u.username,
        u.full_name,
        u.avatar,
        ct.depth + 1 as depth,
        ct.path || c.id as path
    FROM comments c
    JOIN users u ON c.author_id = u.id
    JOIN comment_tree ct ON c.parent_comment_id = ct.id
    WHERE c.status = 'active'
)
SELECT * FROM comment_tree
ORDER BY path, created_at;

-- name: ToggleCommentLike :one
INSERT INTO likes (user_id, comment_id) 
VALUES ($1, $2)
ON CONFLICT (user_id, comment_id) 
DO NOTHING
RETURNING (xmax = 0) as liked;

-- name: DeletePost :exec
UPDATE posts SET status = 'removed', updated_at = NOW() WHERE id = $1 AND author_id = $2;

-- name: PinPost :exec
UPDATE posts SET is_pinned = $1, updated_at = NOW() WHERE id = $2;

-- name: GetPostLikes :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    l.created_at
FROM likes l
JOIN users u ON l.user_id = u.id
WHERE l.post_id = $1 AND u.status = 'active'
ORDER BY l.created_at DESC;

-- name: CreateRepost :one
INSERT INTO posts (author_id, space_id, quoted_post_id, content, visibility)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserLikedPosts :many
SELECT 
    p.*,
    u.username,
    u.full_name,
    u.avatar as author_avatar
FROM posts p
JOIN likes l ON p.id = l.post_id
JOIN users u ON p.author_id = u.id
WHERE l.user_id = $1 AND p.status = 'active'
ORDER BY l.created_at DESC
LIMIT $2 OFFSET $3;


-- name: AdvancedSearchPosts :many
SELECT 
    p.*,
    u.username,
    u.full_name,
    u.avatar as author_avatar,
    c.name as community_name,
    g.name as group_name,
    ts_rank_cd(
        to_tsvector('english', p.content || ' ' || COALESCE(array_to_string(p.tags, ' '), '')),
        plainto_tsquery('english', $1)
    ) as relevance_score
FROM posts p
JOIN users u ON p.author_id = u.id
LEFT JOIN communities c ON p.community_id = c.id
LEFT JOIN groups g ON p.group_id = g.id
WHERE p.space_id = $2
  AND p.status = 'active'
  AND (
    to_tsvector('english', p.content || ' ' || COALESCE(array_to_string(p.tags, ' '), '')) 
    @@ plainto_tsquery('english', $1)
    OR p.content ILIKE '%' || $1 || '%'
    OR p.tags @> ARRAY[$1]::text[]
  )
ORDER BY relevance_score DESC, p.created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetTrendingTopics :many
WITH hashtags AS (
  SELECT
    unnest(regexp_matches(content, '#([a-zA-Z0-9_]+)', 'g')) AS tag,
    p.id,
    p.created_at,
    COALESCE(p.likes_count, 0) + COALESCE(p.comments_count, 0) + COALESCE(p.reposts_count, 0) AS engagement
  FROM posts p
  WHERE p.space_id = $1
    AND p.status = 'active'
    AND p.created_at > NOW() - INTERVAL '24 hours'
)
SELECT
  lower(tag) AS id,
  '#' || tag AS name,
  'General' AS category,
  COUNT(DISTINCT id) AS posts_count,
  -- Trend score: posts * avg_engagement * recency_factor
  COUNT(*) * AVG(engagement)::int AS trend_score
FROM hashtags
GROUP BY lower(tag), tag
ORDER BY trend_score DESC
LIMIT $2 OFFSET $3;