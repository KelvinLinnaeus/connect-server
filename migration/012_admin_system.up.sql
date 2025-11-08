-- Admin System Migration
-- Creates tables and schema for comprehensive admin functionality
-- Note: Admin users are managed via users table with roles field

-- User suspensions table
CREATE TABLE IF NOT EXISTS user_suspensions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    suspended_by UUID REFERENCES users(id) NOT NULL,
    reason TEXT NOT NULL,
    notes TEXT,
    suspended_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    suspended_until TIMESTAMPTZ,
    is_permanent BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_suspensions_user_id ON user_suspensions(user_id);
CREATE INDEX idx_user_suspensions_suspended_until ON user_suspensions(suspended_until);

-- Note: Content reports functionality uses existing 'reports' table from migration 008
-- No need to create a duplicate table

-- System settings table
CREATE TABLE IF NOT EXISTS system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) UNIQUE NOT NULL,
    value JSONB NOT NULL,
    description TEXT,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_system_settings_key ON system_settings(key);

-- Audit logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_user_id UUID REFERENCES users(id) ON DELETE SET NULL NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    details JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_admin_user_id ON audit_logs(admin_user_id);
CREATE INDEX idx_audit_logs_resource_type ON audit_logs(resource_type);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- Space activities table (for the new Space Activities tab)
CREATE TABLE IF NOT EXISTS space_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    space_id UUID REFERENCES spaces(id) ON DELETE CASCADE NOT NULL,
    activity_type VARCHAR(50) NOT NULL,
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
    actor_name VARCHAR(255),
    description TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_space_activities_space_id ON space_activities(space_id);
CREATE INDEX idx_space_activities_activity_type ON space_activities(activity_type);
CREATE INDEX idx_space_activities_created_at ON space_activities(created_at DESC);
CREATE INDEX idx_space_activities_actor_id ON space_activities(actor_id);

-- Add suspended_until column and CHECK constraint to existing status column
-- Note: users.status already exists from migration 001, we're just adding suspended_until and a constraint
ALTER TABLE users ADD COLUMN IF NOT EXISTS suspended_until TIMESTAMPTZ;

-- Add CHECK constraint to existing status column (will fail silently if exists)
-- PostgreSQL doesn't support IF NOT EXISTS for constraints, so we'll use DO block only for this
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.constraint_column_usage
                   WHERE table_name='users' AND constraint_name='users_status_check') THEN
        ALTER TABLE users ADD CONSTRAINT users_status_check
            CHECK (status IN ('active', 'suspended', 'banned', 'pending', 'inactive'));
    END IF;
END $$;

-- Insert default system settings
INSERT INTO system_settings (key, value, description) VALUES
    ('maintenance_mode', 'false'::jsonb, 'Enable/disable maintenance mode'),
    ('registration_enabled', 'true'::jsonb, 'Allow new user registrations'),
    ('email_verification_required', 'true'::jsonb, 'Require email verification for new users'),
    ('max_upload_size_mb', '10'::jsonb, 'Maximum file upload size in MB'),
    ('session_timeout_minutes', '60'::jsonb, 'Session timeout in minutes'),
    ('password_min_length', '8'::jsonb, 'Minimum password length'),
    ('system_notice', '""'::jsonb, 'System-wide notice message')
ON CONFLICT (key) DO NOTHING;

-- Create function to log space activities
CREATE OR REPLACE FUNCTION log_space_activity()
RETURNS TRIGGER AS $$
DECLARE
    v_actor_id UUID;
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Determine the actor_id based on the table
        -- Posts table uses author_id, communities and groups use created_by
        IF TG_TABLE_NAME = 'posts' THEN
            v_actor_id := NEW.author_id;
        ELSIF TG_TABLE_NAME IN ('communities', 'groups') THEN
            v_actor_id := NEW.created_by;
        ELSE
            v_actor_id := NULL;
        END IF;

        INSERT INTO space_activities (space_id, activity_type, actor_id, actor_name, description, metadata)
        VALUES (
            NEW.space_id,
            TG_ARGV[0],
            v_actor_id,
            (SELECT full_name FROM users WHERE id = v_actor_id),
            TG_ARGV[1],
            jsonb_build_object('resource_id', NEW.id)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for automatic activity logging
-- Log new posts
DROP TRIGGER IF EXISTS trigger_log_post_created ON posts;
CREATE TRIGGER trigger_log_post_created
    AFTER INSERT ON posts
    FOR EACH ROW
    EXECUTE FUNCTION log_space_activity('post_created', 'New post created');

-- Log new users (we'll add space_id to context in future)
-- Log new communities
DROP TRIGGER IF EXISTS trigger_log_community_created ON communities;
CREATE TRIGGER trigger_log_community_created
    AFTER INSERT ON communities
    FOR EACH ROW
    EXECUTE FUNCTION log_space_activity('community_created', 'New community created');

-- Log new groups
DROP TRIGGER IF EXISTS trigger_log_group_created ON groups;
CREATE TRIGGER trigger_log_group_created
    AFTER INSERT ON groups
    FOR EACH ROW
    EXECUTE FUNCTION log_space_activity('group_created', 'New group created');

-- Comments
COMMENT ON TABLE user_suspensions IS 'Tracks user account suspensions and bans';
COMMENT ON TABLE system_settings IS 'System-wide configuration settings';
COMMENT ON TABLE audit_logs IS 'Audit trail for admin actions';
COMMENT ON TABLE space_activities IS 'Activity log for space-level events';
