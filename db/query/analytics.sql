-- Analytics and Reporting Queries

-- name: GetSystemMetrics :one
SELECT
    (SELECT COUNT(*) FROM users u2 WHERE u2.space_id = $1 AND u2.status = 'active') as total_users,
    (SELECT COUNT(*) FROM users u3 WHERE u3.space_id = $1 AND u3.updated_at >= NOW() - INTERVAL '1 day') as active_users,
    (SELECT COUNT(*) FROM users u4 WHERE u4.space_id = $1 AND u4.created_at >= CURRENT_DATE) as new_users_today,
    (SELECT COUNT(*) FROM posts p2 WHERE p2.space_id = $1 AND p2.created_at >= CURRENT_DATE) as daily_posts,
    (SELECT COUNT(*) FROM groups g2 WHERE g2.space_id = $1) as total_groups,
    (SELECT COUNT(*) FROM communities c2 WHERE c2.space_id = $1) as total_communities,
    (SELECT COUNT(*) FROM events e2 WHERE e2.space_id = $1 AND e2.status = 'published') as total_events,
    (SELECT COUNT(*) FROM tutoring_sessions ts2 WHERE ts2.space_id = $1 AND ts2.status = 'pending') as pending_tutoring_sessions,
    (SELECT COUNT(*) FROM mentoring_sessions ms2 WHERE ms2.space_id = $1 AND ms2.status = 'pending') as pending_mentoring_sessions,
    (SELECT COUNT(*) FROM reports r2 WHERE r2.space_id = $1 AND r2.status = 'pending') as pending_reports,
    (SELECT COUNT(*) FROM tutor_applications ta2 WHERE ta2.space_id = $1 AND ta2.status = 'pending') as pending_tutor_applications,
    (SELECT COUNT(*) FROM mentor_applications ma2 WHERE ma2.space_id = $1 AND ma2.status = 'pending') as pending_mentor_applications;

-- name: GetUserGrowth :many
SELECT 
    DATE(created_at) as date,
    COUNT(*) as new_users
FROM users 
WHERE space_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date;

-- name: GetEngagementMetrics :many
SELECT 
    DATE(created_at) as date,
    COUNT(*) as post_count,
    SUM(likes_count) as total_likes,
    SUM(comments_count) as total_comments,
    SUM(views_count) as total_views
FROM posts 
WHERE space_id = $1 AND created_at >= NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at)
ORDER BY date;

-- name: GetTopPosts :many
SELECT 
    p.*,
    u.username,
    u.full_name,
    (p.likes_count + p.comments_count + p.views_count) as engagement_score
FROM posts p
JOIN users u ON p.author_id = u.id
WHERE p.space_id = $1 
  AND p.created_at >= NOW() - INTERVAL '7 days'
  AND p.status = 'active'
ORDER BY engagement_score DESC
LIMIT 10;

-- name: GetTopCommunities :many
SELECT 
    c.*,
    (c.member_count + c.post_count) as engagement_score
FROM communities c
WHERE c.space_id = $1
ORDER BY engagement_score DESC
LIMIT 10;

-- name: GetTopGroups :many
SELECT 
    g.*,
    (g.member_count + g.post_count) as engagement_score
FROM groups g
WHERE g.space_id = $1 AND g.status = 'active'
ORDER BY engagement_score DESC
LIMIT 10;

-- Note: User activities are now tracked via user_sessions table

-- name: GetUserActivityStats :many
SELECT
    'session_activity' as action,
    COUNT(*) as count,
    DATE(us.last_activity) as date
FROM user_sessions us
WHERE us.space_id = $1 AND us.last_activity >= NOW() - INTERVAL '7 days'
GROUP BY DATE(us.last_activity)
ORDER BY date DESC, count DESC;

-- name: CreateReport :one
INSERT INTO reports (
    space_id, reporter_id, content_type, content_id, reason, description, priority
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetReport :one
SELECT 
    r.*,
    u.username as reporter_username,
    u.full_name as reporter_full_name,
    reviewer.username as reviewer_username
FROM reports r
JOIN users u ON r.reporter_id = u.id
LEFT JOIN users reviewer ON r.reviewed_by = reviewer.id
WHERE r.id = $1;

-- name: GetPendingReports :many
SELECT 
    r.*,
    u.username as reporter_username,
    u.full_name as reporter_full_name
FROM reports r
JOIN users u ON r.reporter_id = u.id
WHERE r.space_id = $1 AND r.status = 'pending'
ORDER BY 
    CASE r.priority
        WHEN 'urgent' THEN 1
        WHEN 'high' THEN 2
        WHEN 'medium' THEN 3
        ELSE 4
    END,
    r.created_at DESC;

-- name: UpdateReport :one
UPDATE reports 
SET 
    status = $1,
    reviewed_by = $2,
    reviewed_at = NOW(),
    moderation_notes = $3,
    actions_taken = $4,
    updated_at = NOW()
WHERE id = $5
RETURNING *;

-- name: GetReportsByContent :many
SELECT *
FROM reports 
WHERE content_type = $1 AND content_id = $2
ORDER BY created_at DESC;

-- name: GetModerationQueue :many
SELECT 
    r.*,
    u.username as reporter_username,
    u.full_name as reporter_full_name,
    reviewer.username as reviewer_username
FROM reports r
JOIN users u ON r.reporter_id = u.id
LEFT JOIN users reviewer ON r.reviewed_by = reviewer.id
WHERE r.space_id = $1
ORDER BY r.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetContentModerationStats :one
SELECT 
    COUNT(*) as total_reports,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_reports,
    COUNT(*) FILTER (WHERE status = 'approved') as approved_reports,
    COUNT(*) FILTER (WHERE status = 'rejected') as rejected_reports,
    COUNT(*) FILTER (WHERE priority = 'urgent') as urgent_reports
FROM reports 
WHERE space_id = $1;

-- name: GetTutoringStats :one
SELECT
    COUNT(*) as total_sessions,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_sessions,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_sessions,
    COALESCE(AVG(rating), 0) as average_rating,
    COUNT(*) FILTER (WHERE rating IS NOT NULL) as rated_sessions
FROM tutoring_sessions
WHERE space_id = $1;

-- name: GetMentoringStats :one
SELECT
    COUNT(*) as total_sessions,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_sessions,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_sessions,
    COALESCE(AVG(rating), 0) as average_rating,
    COUNT(*) FILTER (WHERE rating IS NOT NULL) as rated_sessions
FROM mentoring_sessions
WHERE space_id = $1;

-- name: GetPopularSubjects :many
SELECT 
    subject,
    COUNT(*) as session_count,
    AVG(rating) as average_rating
FROM tutoring_sessions 
WHERE space_id = $1 AND status = 'completed'
GROUP BY subject
ORDER BY session_count DESC
LIMIT 10;

-- name: GetPopularIndustries :many
SELECT
    industry,
    COUNT(*) as session_count,
    AVG(ms.rating) as average_rating
FROM mentor_profiles mp
JOIN mentoring_sessions ms ON mp.user_id = ms.mentor_id
WHERE mp.space_id = $1 AND ms.status = 'completed'
GROUP BY industry
ORDER BY session_count DESC
LIMIT 10;

-- name: GetUserEngagementRanking :many
SELECT 
    u.id,
    u.username,
    u.full_name,
    u.avatar,
    COUNT(p.id) as post_count,
    COUNT(DISTINCT f1.follower_id) as followers_count,
    COUNT(DISTINCT f2.following_id) as following_count,
    (COUNT(p.id) + COUNT(DISTINCT f1.follower_id) * 2) as engagement_score
FROM users u
LEFT JOIN posts p ON u.id = p.author_id AND p.status = 'active'
LEFT JOIN follows f1 ON u.id = f1.following_id
LEFT JOIN follows f2 ON u.id = f2.follower_id
WHERE u.space_id = $1 AND u.status = 'active'
GROUP BY u.id
ORDER BY engagement_score DESC
LIMIT 20;

-- name: GetSpaceStats :one
SELECT
    s.name,
    s.slug,
    (SELECT COUNT(*) FROM users u5 WHERE u5.space_id = s.id AND u5.status = 'active') as user_count,
    (SELECT COUNT(*) FROM posts p3 WHERE p3.space_id = s.id AND p3.status = 'active') as post_count,
    (SELECT COUNT(*) FROM communities c3 WHERE c3.space_id = s.id) as community_count,
    (SELECT COUNT(*) FROM groups g3 WHERE g3.space_id = s.id AND g3.status = 'active') as group_count,
    (SELECT created_at FROM users u6 WHERE u6.space_id = s.id ORDER BY u6.created_at ASC LIMIT 1) as first_user_date
FROM spaces s
WHERE s.id = $1;

-- Note: RecordSystemMetrics and GetHistoricalMetrics removed
-- system_metrics table removed - all metrics derived in real-time from existing tables
-- See GetSystemMetrics above for real-time metrics calculation

-- Admin Content Report Management Queries

-- -- name: UpdateContentReportWithAction :exec
-- UPDATE reports
-- SET
--     status = $1,
--     reviewed_by = $2,
--     moderation_notes = $3,
--     actions_taken = $4,
--     reviewed_at = NOW(),
--     updated_at = NOW()
-- WHERE id = $5;

-- -- name: UpdateContentReportPriority :exec
-- UPDATE reports
-- SET priority = $1, updated_at = NOW()
-- WHERE id = $2;