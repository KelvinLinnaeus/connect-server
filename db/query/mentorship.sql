-- Mentorship and Tutoring Queries

-- name: CreateTutorApplication :one
INSERT INTO tutor_applications (
    applicant_id, space_id, subjects, hourly_rate, availability,
    experience, qualifications, teaching_style, motivation, reference_letters
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: CreateMentorApplication :one
INSERT INTO mentor_applications (
    applicant_id, space_id, industry, company, position, experience,
    specialties, achievements, mentorship_experience, availability,
    motivation, approach_description, linkedin_profile, portfolio
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetTutorApplication :one
SELECT 
    ta.*,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM tutor_applications ta
JOIN users u ON ta.applicant_id = u.id
WHERE ta.id = $1;

-- name: GetMentorApplication :one
SELECT 
    ma.*,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM mentor_applications ma
JOIN users u ON ma.applicant_id = u.id
WHERE ma.id = $1;

-- name: GetPendingTutorApplications :many
SELECT 
    ta.*,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM tutor_applications ta
JOIN users u ON ta.applicant_id = u.id
WHERE ta.space_id = $1 AND ta.status = 'pending'
ORDER BY ta.submitted_at DESC;

-- name: GetPendingMentorApplications :many
SELECT 
    ma.*,
    u.username,
    u.full_name,
    u.avatar,
    u.email,
    u.department,
    u.level
FROM mentor_applications ma
JOIN users u ON ma.applicant_id = u.id
WHERE ma.space_id = $1 AND ma.status = 'pending'
ORDER BY ma.submitted_at DESC;

-- name: UpdateTutorApplication :one
UPDATE tutor_applications 
SET 
    status = $1,
    reviewed_at = NOW(),
    reviewed_by = $2,
    reviewer_notes = $3
WHERE id = $4
RETURNING *;

-- name: UpdateMentorApplication :one
UPDATE mentor_applications 
SET 
    status = $1,
    reviewed_at = NOW(),
    reviewed_by = $2,
    reviewer_notes = $3
WHERE id = $4
RETURNING *;

-- name: CreateTutorProfile :one
INSERT INTO tutor_profiles (
    user_id, space_id, subjects, hourly_rate, description,
    availability, experience, qualifications, verified
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (user_id) 
DO UPDATE SET
    subjects = EXCLUDED.subjects,
    hourly_rate = EXCLUDED.hourly_rate,
    description = EXCLUDED.description,
    availability = EXCLUDED.availability,
    experience = EXCLUDED.experience,
    qualifications = EXCLUDED.qualifications,
    verified = EXCLUDED.verified,
    updated_at = NOW()
RETURNING *;

-- name: CreateMentorProfile :one
INSERT INTO mentor_profiles (
    user_id, space_id, industry, company, position, experience,
    specialties, description, availability, verified
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT (user_id) 
DO UPDATE SET
    industry = EXCLUDED.industry,
    company = EXCLUDED.company,
    position = EXCLUDED.position,
    experience = EXCLUDED.experience,
    specialties = EXCLUDED.specialties,
    description = EXCLUDED.description,
    availability = EXCLUDED.availability,
    verified = EXCLUDED.verified,
    updated_at = NOW()
RETURNING *;

-- name: GetTutorProfile :one
SELECT 
    tp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.user_id = $1;

-- name: GetMentorProfile :one
SELECT 
    mp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.user_id = $1;

-- name: SearchTutors :many
SELECT
    tp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    COALESCE(
        (SELECT AVG(rating)
         FROM tutoring_sessions
         WHERE tutor_id = tp.user_id
           AND rating IS NOT NULL), 0
    ) as avg_rating,
    COALESCE(
        (SELECT COUNT(*)
         FROM tutoring_sessions
         WHERE tutor_id = tp.user_id
           AND status = 'completed'), 0
    ) as completed_sessions
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.space_id = $1
  AND tp.is_available = true
  AND (tp.subjects @> $2 OR $2 IS NULL)
  AND (tp.availability @> $3 OR $3 IS NULL)
  AND (tp.hourly_rate <= $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'rating' THEN
        COALESCE((SELECT AVG(rating)
                  FROM tutoring_sessions
                  WHERE tutor_id = tp.user_id
                    AND rating IS NOT NULL), 0)
    END DESC NULLS LAST,
    CASE WHEN $5 = 'experience' THEN
        COALESCE((SELECT COUNT(*)
                  FROM tutoring_sessions
                  WHERE tutor_id = tp.user_id
                    AND status = 'completed'), 0)
    END DESC,
    tp.hourly_rate ASC
LIMIT $6 OFFSET $7;



-- name: SearchMentors :many
SELECT
    mp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified AS user_verified,
    u.department,
    u.level,
    COALESCE(
        (SELECT AVG(rating)
         FROM mentoring_sessions
         WHERE mentor_id = mp.user_id
           AND rating IS NOT NULL), 0
    ) AS avg_rating,
    COALESCE(
        (SELECT COUNT(*)
         FROM mentoring_sessions
         WHERE mentor_id = mp.user_id
           AND status = 'completed'), 0
    ) AS completed_sessions
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.space_id = $1
  AND mp.is_available = true
  OR (mp.industry = $2 OR $2 IS NULL)
  OR (mp.specialties @> $3 OR $3 IS NULL)
  OR (mp.experience >= $4 OR $4 IS NULL)
ORDER BY
    CASE WHEN $5 = 'rating' THEN
        COALESCE((SELECT AVG(rating)
                  FROM mentoring_sessions
                  WHERE mentor_id = mp.user_id
                    AND rating IS NOT NULL), 0)
    END DESC NULLS LAST,
    CASE WHEN $5 = 'experience' THEN mp.experience END DESC,
    COALESCE((SELECT COUNT(*)
              FROM mentoring_sessions
              WHERE mentor_id = mp.user_id
                AND status = 'completed'), 0) DESC
LIMIT $6 OFFSET $7;


-- name: CreateTutoringSession :one
INSERT INTO tutoring_sessions (
    tutor_id, student_id, space_id, subject, scheduled_at,
    duration, hourly_rate, total_amount, student_notes
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: CreateMentoringSession :one
INSERT INTO mentoring_sessions (
    mentor_id, mentee_id, space_id, topic, scheduled_at,
    duration, mentee_notes
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetTutoringSession :one
SELECT 
    ts.*,
    tutor.username as tutor_username,
    tutor.full_name as tutor_full_name,
    tutor.avatar as tutor_avatar,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users tutor ON ts.tutor_id = tutor.id
JOIN users student ON ts.student_id = student.id
WHERE ts.id = $1;

-- name: GetMentoringSession :one
SELECT 
    ms.*,
    mentor.username as mentor_username,
    mentor.full_name as mentor_full_name,
    mentor.avatar as mentor_avatar,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentor ON ms.mentor_id = mentor.id
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE ms.id = $1;

-- name: GetUserTutoringSessions :many
SELECT 
    ts.*,
    tutor.username as tutor_username,
    tutor.full_name as tutor_full_name,
    tutor.avatar as tutor_avatar,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users tutor ON ts.tutor_id = tutor.id
JOIN users student ON ts.student_id = student.id
WHERE (ts.tutor_id = $1 OR ts.student_id = $1)
ORDER BY ts.scheduled_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserMentoringSessions :many
SELECT 
    ms.*,
    mentor.username as mentor_username,
    mentor.full_name as mentor_full_name,
    mentor.avatar as mentor_avatar,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentor ON ms.mentor_id = mentor.id
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE (ms.mentor_id = $1 OR ms.mentee_id = $1)
ORDER BY ms.scheduled_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSessionStatus :one
UPDATE tutoring_sessions 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateMentoringSessionStatus :one
UPDATE mentoring_sessions 
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: AddSessionMeetingLink :exec
UPDATE tutoring_sessions 
SET meeting_link = $1, updated_at = NOW()
WHERE id = $2;

-- name: AddMentoringSessionMeetingLink :exec
UPDATE mentoring_sessions 
SET meeting_link = $1, updated_at = NOW()
WHERE id = $2;

-- name: RateTutoringSession :one
UPDATE tutoring_sessions 
SET rating = $1, review = $2, updated_at = NOW()
WHERE id = $3 AND student_id = $4
RETURNING *;

-- name: RateMentoringSession :one
UPDATE mentoring_sessions 
SET rating = $1, review = $2, updated_at = NOW()
WHERE id = $3 AND mentee_id = $4
RETURNING *;

-- name: UpdateTutorAvailability :one
UPDATE tutor_profiles 
SET is_available = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING *;

-- name: UpdateMentorAvailability :one
UPDATE mentor_profiles 
SET is_available = $1, updated_at = NOW()
WHERE user_id = $2
RETURNING *;

-- name: GetTutorReviews :many
SELECT 
    ts.rating,
    ts.review,
    ts.created_at,
    student.username as student_username,
    student.full_name as student_full_name,
    student.avatar as student_avatar
FROM tutoring_sessions ts
JOIN users student ON ts.student_id = student.id
WHERE ts.tutor_id = $1 AND ts.rating IS NOT NULL
ORDER BY ts.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetMentorReviews :many
SELECT
    ms.rating,
    ms.review,
    ms.created_at,
    mentee.username as mentee_username,
    mentee.full_name as mentee_full_name,
    mentee.avatar as mentee_avatar
FROM mentoring_sessions ms
JOIN users mentee ON ms.mentee_id = mentee.id
WHERE ms.mentor_id = $1 AND ms.rating IS NOT NULL
ORDER BY ms.created_at DESC
LIMIT $2 OFFSET $3;

-- Admin Application Management Queries

-- name: GetAllTutorApplications :many
SELECT
    ta.applicant_id,
    ta.id,
    ta.subjects,
    ta.hourly_rate,
    ta.status,
    ta.submitted_at,
    u.full_name,
    u.id as user_id
 FROM tutor_applications ta
 JOIN users u ON ta.applicant_id = u.id
ORDER BY submitted_at DESC
LIMIT $1 OFFSET $2;

-- name: GetTutorApplicationsByStatus :many
SELECT * FROM tutor_applications
WHERE status = $1
ORDER BY submitted_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateTutorApplicationStatus :one
UPDATE tutor_applications
SET
    status = $1,
    reviewed_by = $2,
    reviewer_notes = $3,
    reviewed_at = NOW()
WHERE id = $4
RETURNING *;

-- name: GetAllMentorApplications :many
SELECT 
    mt.id,
    mt.applicant_id,
    mt.industry,
    mt.experience,
    mt.specialties,
    mt.status,
    mt.company,
    mt.position,
    mt.submitted_at,
    u.full_name,
    u.id as user_id
FROM mentor_applications mt
JOIN users u ON mt.applicant_id = u.id
ORDER BY submitted_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdateMentorApplicationStatus :one
UPDATE mentor_applications
SET
    status = $1,
    reviewed_by = $2,
    reviewer_notes = $3,
    reviewed_at = NOW()
WHERE id = $4
RETURNING *;

-- Recommendation Queries

-- name: GetRecommendedTutors :many
SELECT
    tp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    (SELECT AVG(rating) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND rating IS NOT NULL) as avg_rating,
    (SELECT COUNT(*) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND status = 'completed') as completed_sessions
FROM tutor_profiles tp
JOIN users u ON tp.user_id = u.id
WHERE tp.space_id = $1
  AND tp.is_available = true
  AND tp.user_id != $2
  AND (
    -- Match by department
    u.department = (SELECT department FROM users WHERE id = $2)
    OR
    -- Match by level
    u.level = (SELECT level FROM users WHERE id = $2)
    OR
    -- Match by subjects (if user has any in their profile)
    tp.subjects && (SELECT interests FROM users WHERE id = $2)
  )
ORDER BY
    -- Prioritize verified tutors
    u.verified DESC,
    -- Then by rating
    (SELECT AVG(rating) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND rating IS NOT NULL) DESC NULLS LAST,
    -- Then by completed sessions
    (SELECT COUNT(*) FROM tutoring_sessions WHERE tutor_id = tp.user_id AND status = 'completed') DESC,
    -- Finally by availability status
    tp.is_available DESC
LIMIT $3;

-- name: GetRecommendedMentors :many
SELECT
    mp.*,
    u.username,
    u.full_name,
    u.avatar,
    u.verified as user_verified,
    u.department,
    u.level,
    (SELECT AVG(rating) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND rating IS NOT NULL) as avg_rating,
    (SELECT COUNT(*) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND status = 'completed') as completed_sessions
FROM mentor_profiles mp
JOIN users u ON mp.user_id = u.id
WHERE mp.space_id = $1
  AND mp.is_available = true
  AND mp.user_id != $2
  AND (
    -- Match by department
    u.department = (SELECT department FROM users WHERE id = $2)
    OR
    -- Match by industry relevant to user's major
    mp.industry IN (
      SELECT CASE
        WHEN major LIKE '%Computer%' OR major LIKE '%Engineering%' THEN 'Technology'
        WHEN major LIKE '%Business%' OR major LIKE '%Finance%' THEN 'Finance'
        WHEN major LIKE '%Art%' OR major LIKE '%Design%' THEN 'Creative'
        ELSE 'General'
      END
      FROM users WHERE id = $2
    )
    OR
    -- Match by specialties overlapping with interests
    mp.specialties && (SELECT interests FROM users WHERE id = $2)
  )
ORDER BY
    -- Prioritize verified mentors
    u.verified DESC,
    -- Then by rating
    (SELECT AVG(rating) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND rating IS NOT NULL) DESC NULLS LAST,
    -- Then by experience
    mp.experience DESC,
    -- Then by completed sessions
    (SELECT COUNT(*) FROM mentoring_sessions WHERE mentor_id = mp.user_id AND status = 'completed') DESC,
    -- Finally by availability status
    mp.is_available DESC
LIMIT $3;

-- name: GetUserTutorApplicationStatus :one
SELECT status FROM tutor_applications
WHERE applicant_id = $1 AND space_id = $2
ORDER BY submitted_at DESC
LIMIT 1;

-- name: GetUserMentorApplicationStatus :one
SELECT status FROM mentor_applications
WHERE applicant_id = $1 AND space_id = $2
ORDER BY submitted_at DESC
LIMIT 1;

-- name: GetUserMentorApplicationStatusById :one
SELECT status FROM mentor_applications
WHERE id = $1
LIMIT 1;

-- name: GetUserTutorApplicationStatusById :one
SELECT status FROM tutor_applications
WHERE id =$1
LIMIT 1;

