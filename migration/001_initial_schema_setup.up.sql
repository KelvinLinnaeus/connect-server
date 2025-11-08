-- UNIVYN Database Migration
-- Version: 001_initial_schema_setup UP
-- Description: Create core tables and relationships

BEGIN;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create spaces table (foundation of multi-tenant architecture)
CREATE TABLE spaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(150) UNIQUE NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    type VARCHAR(50),
    logo TEXT,
    cover_image TEXT,
    location VARCHAR(50),
    website VARCHAR(150),
    contact_email VARCHAR(100),
    phone_number VARCHAR(10),
    status VARCHAR(20) DEFAULT 'active',
    settings JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    avatar TEXT,
    bio TEXT,
    verified BOOLEAN DEFAULT false,
    roles TEXT[] DEFAULT '{}',
    level VARCHAR(50),
    department VARCHAR(100),
    major VARCHAR(100),
    year INTEGER,
    interests TEXT[] DEFAULT '{}',
    followers_count INTEGER DEFAULT 0,
    following_count INTEGER DEFAULT 0,
    mentor_status VARCHAR(20) DEFAULT 'pending',
    tutor_status VARCHAR(20) DEFAULT 'pending',
    status VARCHAR(20) DEFAULT 'active',
    settings JSONB,
    phone_number VARCHAR(10) NOT NULL,
    additional_phone_number VARCHAR(10),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (space_id) REFERENCES spaces(id) ON DELETE CASCADE
);

COMMIT;