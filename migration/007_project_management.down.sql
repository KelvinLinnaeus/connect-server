-- UNIVYN Database Migration
-- Version: 007_project_management DOWN
-- Description: Drop project roles and applications system

BEGIN;

DROP TABLE IF EXISTS group_applications;
DROP TABLE IF EXISTS group_roles;

COMMIT;