-- UNIVYN Database Migration
-- Version: 009_additional_enhancements UP
-- Description: Create additional enhancement tables
-- Note: User preferences stored in users.settings field (JSONB)

BEGIN;

-- Create user sessions table (tracks user sessions and activities)
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    space_id UUID NOT NULL,
    username varchar NOT NULL,
    refresh_token varchar NOT NULL,
    user_agent varchar NOT NULL,
    ip_address VARCHAR(50),
    is_blocked boolean NOT NULL DEFAULT false,
    last_activity TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

-- Create files table
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    user_id UUID NOT NULL,
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255),
    mime_type VARCHAR(100),
    size BIGINT,
    path TEXT NOT NULL,
    bucket VARCHAR(100),
    is_public BOOLEAN DEFAULT false,
    metadata JSONB,
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create trending topics table
CREATE TABLE trending_topics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    post_count INTEGER DEFAULT 0,
    trend_score FLOAT DEFAULT 0,
    period VARCHAR(20) DEFAULT 'daily',
    recorded_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

-- Create email queue table
CREATE TABLE email_queue (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    recipient_email VARCHAR(100) NOT NULL,
    subject VARCHAR(200) NOT NULL,
    template_name VARCHAR(100),
    template_data JSONB,
    status VARCHAR(20) DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    error_message TEXT,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

COMMIT;