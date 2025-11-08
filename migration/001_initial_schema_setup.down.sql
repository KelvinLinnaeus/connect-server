-- UNIVYN Database Migration
-- Version: 001_initial_schema_setup DOWN
-- Description: Drop core tables and relationships

BEGIN;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS spaces;
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;