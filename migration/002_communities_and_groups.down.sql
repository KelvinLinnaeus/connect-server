-- UNIVYN Database Migration
-- Version: 002_communities_and_groups DOWN
-- Description: Drop community and group structures

BEGIN;

DROP TABLE IF EXISTS group_join_requests;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS community_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS communities;

COMMIT;