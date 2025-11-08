-- Migration: Add login attempts tracking for account lockout protection
-- Purpose: Track failed login attempts to implement brute-force protection
-- Version: 011
-- Created: 2025-11-05

-- Table to track login attempts (both successful and failed)
-- Links to user_sessions for successful logins
CREATE TABLE IF NOT EXISTS login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    attempt_result VARCHAR(20) NOT NULL, -- 'success', 'failed_password', 'failed_user_not_found', 'account_locked'
    attempted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID, -- NULL if user not found or invalid
    space_id UUID,
    session_id UUID, -- References user_sessions for successful logins, NULL for failed attempts

    -- Foreign key constraints
    CONSTRAINT fk_login_attempts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_login_attempts_space FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    CONSTRAINT fk_login_attempts_session FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE SET NULL,
    CONSTRAINT check_attempt_result CHECK (attempt_result IN ('success', 'failed_password', 'failed_user_not_found', 'account_locked', 'account_disabled'))
);

-- Indexes for efficient lookups
CREATE INDEX idx_login_attempts_username ON login_attempts(username, attempted_at DESC);
CREATE INDEX idx_login_attempts_ip ON login_attempts(ip_address, attempted_at DESC);
CREATE INDEX idx_login_attempts_user_id ON login_attempts(user_id, attempted_at DESC);
CREATE INDEX idx_login_attempts_attempted_at ON login_attempts(attempted_at DESC);

-- Note: Account lockouts managed via users table fields (is_locked, locked_until, failed_login_attempts)

-- Add account status fields to users table if not exist
-- Using ALTER TABLE ... ADD COLUMN IF NOT EXISTS (PostgreSQL 9.6+)
-- This syntax is parseable by SQLC unlike DO blocks
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_locked BOOLEAN DEFAULT FALSE NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS locked_until TIMESTAMP;
ALTER TABLE users ADD COLUMN IF NOT EXISTS failed_login_attempts INT DEFAULT 0 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_failed_login TIMESTAMP;

-- Create index on users.is_locked for efficient filtering
CREATE INDEX IF NOT EXISTS idx_users_is_locked ON users(is_locked) WHERE is_locked = TRUE;
CREATE INDEX IF NOT EXISTS idx_users_locked_until ON users(locked_until) WHERE locked_until IS NOT NULL;

-- Function to automatically unlock expired lockouts
CREATE OR REPLACE FUNCTION unlock_expired_accounts()
RETURNS void AS $$
BEGIN
    UPDATE users
    SET is_locked = FALSE,
        locked_until = NULL,
        failed_login_attempts = 0
    WHERE is_locked = TRUE
      AND locked_until IS NOT NULL
      AND locked_until < NOW();
END;
$$ LANGUAGE plpgsql;

-- Function to clean up old login attempts (keep only last 90 days)
CREATE OR REPLACE FUNCTION cleanup_old_login_attempts()
RETURNS void AS $$
BEGIN
    DELETE FROM login_attempts
    WHERE attempted_at < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- Comments for documentation
COMMENT ON TABLE login_attempts IS 'Tracks all login attempts for security audit and brute-force protection';
COMMENT ON COLUMN users.is_locked IS 'Whether the account is currently locked due to security reasons';
COMMENT ON COLUMN users.locked_until IS 'Timestamp until which the account is locked (NULL if not locked or permanently locked)';
COMMENT ON COLUMN users.failed_login_attempts IS 'Counter for consecutive failed login attempts';
COMMENT ON COLUMN users.last_failed_login IS 'Timestamp of the last failed login attempt';
