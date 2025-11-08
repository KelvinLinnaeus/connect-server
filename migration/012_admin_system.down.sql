-- Rollback Admin System Migration

-- Drop triggers
DROP TRIGGER IF EXISTS trigger_log_post_created ON posts;
DROP TRIGGER IF EXISTS trigger_log_community_created ON communities;
DROP TRIGGER IF EXISTS trigger_log_group_created ON groups;

-- Drop function
DROP FUNCTION IF EXISTS log_space_activity();

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS space_activities;
DROP TABLE IF EXISTS system_settings;
DROP TABLE IF EXISTS user_suspensions;

-- Remove columns and constraints from users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users DROP COLUMN IF EXISTS suspended_until;
