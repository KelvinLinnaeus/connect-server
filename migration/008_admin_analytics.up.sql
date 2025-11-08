-- UNIVYN Database Migration
-- Version: 008_admin_analytics UP
-- Description: Create analytics and reporting tables
-- Note: Admin users managed via users.roles field, user activities tracked via user_sessions table

BEGIN;

-- Create reports table
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    reporter_id UUID NOT NULL,
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,
    reason VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    priority VARCHAR(20) DEFAULT 'medium',
    reviewed_by UUID,
    reviewed_at TIMESTAMPTZ,
    moderation_notes TEXT,
    actions_taken JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Note: system_metrics table removed - all metrics are derived from existing tables in real-time
-- See analytics.sql GetSystemMetrics query for real-time metric calculations from users, posts, groups, etc.

COMMIT;