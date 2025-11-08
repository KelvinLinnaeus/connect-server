-- UNIVYN Database Migration
-- Version: 008_admin_analytics DOWN
-- Description: Drop analytics and reporting tables

BEGIN;

-- Note: system_metrics table removed (metrics derived from existing tables)
DROP TABLE IF EXISTS reports;

COMMIT;