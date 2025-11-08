-- UNIVYN Database Migration
-- Version: 004_mentorship_tutoring UP
-- Description: Create mentorship and tutoring system tables

BEGIN;

-- Create tutor applications table
CREATE TABLE tutor_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    applicant_id UUID NOT NULL,
    space_id UUID NOT NULL,
    subjects TEXT[] DEFAULT '{}',
    hourly_rate NUMERIC(10,2),
    availability JSONB NOT NULL,
    experience TEXT,
    qualifications TEXT,
    teaching_style TEXT,
    motivation TEXT,
    reference_letters TEXT, -- CHANGED: renamed from 'references'
    status VARCHAR(20) DEFAULT 'pending',
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    reviewer_notes TEXT,
    FOREIGN KEY (applicant_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create mentor applications table
CREATE TABLE mentor_applications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    applicant_id UUID NOT NULL,
    space_id UUID NOT NULL,
    industry VARCHAR(100) NOT NULL,
    company VARCHAR(100),
    position VARCHAR(100),
    experience INTEGER NOT NULL,
    specialties TEXT[] DEFAULT '{}',
    achievements TEXT,
    mentorship_experience TEXT,
    availability JSONB NOT NULL,
    motivation TEXT,
    approach_description TEXT,
    linkedin_profile VARCHAR(200),
    portfolio VARCHAR(200),
    status VARCHAR(20) DEFAULT 'pending',
    submitted_at TIMESTAMPTZ DEFAULT NOW(),
    reviewed_at TIMESTAMPTZ,
    reviewed_by UUID,
    reviewer_notes TEXT,
    FOREIGN KEY (applicant_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- Create tutor profiles table
CREATE TABLE tutor_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    space_id UUID NOT NULL,
    subjects TEXT[] DEFAULT '{}',
    hourly_rate NUMERIC(10,2),
    rating FLOAT DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    description TEXT,
    availability JSONB,
    experience TEXT,
    qualifications TEXT,
    verified BOOLEAN DEFAULT false,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- Create mentor profiles table
CREATE TABLE mentor_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    space_id UUID NOT NULL,
    industry VARCHAR(100) NOT NULL,
    company VARCHAR(100),
    position VARCHAR(100),
    experience INTEGER NOT NULL,
    specialties TEXT[] DEFAULT '{}',
    rating FLOAT DEFAULT 0,
    review_count INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    availability JSONB,
    description TEXT,
    verified BOOLEAN DEFAULT false,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    UNIQUE(user_id)
);

-- Create tutoring sessions table
CREATE TABLE tutoring_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tutor_id UUID NOT NULL,
    student_id UUID NOT NULL,
    space_id UUID NOT NULL,
    subject VARCHAR(100) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    scheduled_at TIMESTAMPTZ NOT NULL,
    duration INTEGER NOT NULL,
    hourly_rate NUMERIC(10,2),
    total_amount NUMERIC(10,2),
    student_notes TEXT,
    tutor_notes TEXT,
    meeting_link TEXT,
    rating INTEGER,
    review TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (tutor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

-- Create mentoring sessions table
CREATE TABLE mentoring_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    mentor_id UUID NOT NULL,
    mentee_id UUID NOT NULL,
    space_id UUID NOT NULL,
    topic VARCHAR(150) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    scheduled_at TIMESTAMPTZ NOT NULL,
    duration INTEGER NOT NULL,
    mentee_notes TEXT,
    mentor_notes TEXT,
    meeting_link TEXT,
    rating INTEGER,
    review TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (mentor_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (mentee_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

COMMIT;