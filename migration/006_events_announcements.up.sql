-- UNIVYN Database Migration
-- Version: 006_events_announcements UP
-- Description: Create events, announcements, and past questions

BEGIN;

-- Create events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    title VARCHAR(150) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    location VARCHAR(100),
    venue_details TEXT,
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    timezone VARCHAR(50),
    organizer UUID,
    tags TEXT[] DEFAULT '{}',
    image_url TEXT,
    max_attendees INTEGER,
    current_attendees INTEGER DEFAULT 0,
    registration_required BOOLEAN DEFAULT false,
    registration_deadline TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'published',
    is_public BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (organizer) REFERENCES users(id) ON DELETE SET NULL
);


-- Create event attendees table
CREATE TABLE event_attendees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status VARCHAR(20) DEFAULT 'registered',
    role VARCHAR(30) DEFAULT 'attendee',
    registered_at TIMESTAMPTZ DEFAULT NOW(),
    attended_at TIMESTAMPTZ,
    notes TEXT,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(event_id, user_id)
);

-- Create announcements table
CREATE TABLE announcements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    title VARCHAR(150) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    target_audience TEXT[] DEFAULT '{}',
    priority VARCHAR(20) DEFAULT 'medium',
    status VARCHAR(20) DEFAULT 'published',
    author_id UUID,
    scheduled_for TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    attachments JSONB,
    is_pinned BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Create past questions table
CREATE TABLE past_questions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    course_code VARCHAR(50) NOT NULL,
    course_title VARCHAR(150) NOT NULL,
    department VARCHAR(100),
    academic_year VARCHAR(20) NOT NULL,
    semester VARCHAR(20),
    instructor VARCHAR(100),
    exam_type VARCHAR(50),
    file_url TEXT NOT NULL,
    file_name VARCHAR(200),
    file_size INTEGER,
    uploaded_by UUID NOT NULL,
    verified BOOLEAN DEFAULT false,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    download_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE,
    FOREIGN KEY (uploaded_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (verified_by) REFERENCES users(id) ON DELETE SET NULL
);

COMMIT;