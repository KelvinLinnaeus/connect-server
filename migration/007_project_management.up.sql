-- UNIVYN Database Migration
-- Version: 007_project_management UP
-- Description: Create project roles and applications system

BEGIN;

-- Create project roles table
CREATE TABLE group_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    group_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    slots_total INTEGER NOT NULL,
    slots_filled INTEGER DEFAULT 0,
    requirements TEXT,
    skills_required TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- Create role applications table
CREATE TABLE group_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL,
    user_id UUID NOT NULL,
    message TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    applied_at TIMESTAMPTZ DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    review_notes TEXT,
    FOREIGN KEY (role_id) REFERENCES group_roles(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(role_id, user_id)
);

COMMIT;