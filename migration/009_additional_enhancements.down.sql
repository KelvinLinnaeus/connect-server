-- UNIVYN Database Migration
-- Version: 009_additional_enhancements DOWN
-- Description: Drop additional enhancement tables

BEGIN;

DROP TABLE IF EXISTS email_queue;
DROP TABLE IF EXISTS trending_topics;
DROP TABLE IF EXISTS files;
DROP TABLE IF EXISTS user_sessions;

COMMIT;