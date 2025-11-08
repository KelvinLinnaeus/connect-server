-- UNIVYN Database Migration
-- Version: 004_mentorship_tutoring DOWN
-- Description: Drop mentorship and tutoring system tables

BEGIN;

DROP TABLE IF EXISTS mentoring_sessions;
DROP TABLE IF EXISTS tutoring_sessions;
DROP TABLE IF EXISTS mentor_profiles;
DROP TABLE IF EXISTS tutor_profiles;
DROP TABLE IF EXISTS mentor_applications;
DROP TABLE IF EXISTS tutor_applications;

COMMIT;