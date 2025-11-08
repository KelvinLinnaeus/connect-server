-- Rollback: Remove login attempts tracking
-- Version: 011
-- Created: 2025-11-05

-- Drop functions
DROP FUNCTION IF EXISTS unlock_expired_accounts();
DROP FUNCTION IF EXISTS cleanup_old_login_attempts();

-- Drop tables (cascade to remove foreign key constraints)
DROP TABLE IF EXISTS login_attempts CASCADE;

-- Remove columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS is_locked;
ALTER TABLE users DROP COLUMN IF EXISTS locked_until;
ALTER TABLE users DROP COLUMN IF EXISTS failed_login_attempts;
ALTER TABLE users DROP COLUMN IF EXISTS last_failed_login;

-- Drop indexes (if they weren't dropped with the tables)
DROP INDEX IF EXISTS idx_users_is_locked;
DROP INDEX IF EXISTS idx_users_locked_until;
