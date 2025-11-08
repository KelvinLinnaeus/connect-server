-- UNIVYN Database Migration
-- Version: 003_social_features DOWN
-- Description: Drop social interaction tables

BEGIN;

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;

COMMIT;