# Admin System Endpoints Documentation

## Overview
This document outlines all admin endpoints needed for the Connect platform admin panel. The admin system provides comprehensive tools for managing users, content, communities, applications, and system settings.

## Authentication Requirements
All admin endpoints require:
- Valid JWT token in Authorization header
- User must have `admin` or `super_admin` role
- Some endpoints require specific permissions

## Endpoint Categories

### 1. Admin Authentication

#### POST /api/admin/auth/login
**Description:** Admin-specific login endpoint with enhanced security

**Request:**
```json
{
  "email": "admin@university.edu",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "data": {
    "access_token": "jwt_token",
    "refresh_token": "refresh_token",
    "user": {
      "id": "uuid",
      "email": "admin@university.edu",
      "full_name": "Admin User",
      "role": "admin",
      "permissions": ["user_management", "content_management", ...]
    }
  }
}
```

---

### 2. User Management

#### GET /api/admin/users
**Description:** Get all users with filtering, pagination, and search

**Permissions:** `user_management`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)
- `search`: Search by name, email, or username
- `role`: Filter by role
- `status`: Filter by status (active, suspended, banned)
- `sort`: Sort field (created_at, name, etc.)
- `order`: Sort order (asc, desc)

**Response:**
```json
{
  "data": {
    "users": [{
      "id": "uuid",
      "username": "student1",
      "email": "student1@university.edu",
      "full_name": "John Doe",
      "role": "user",
      "status": "active",
      "verified": true,
      "created_at": "2024-01-01T00:00:00Z",
      "last_login": "2024-01-15T10:30:00Z",
      "posts_count": 45,
      "followers_count": 120
    }],
    "total": 1500,
    "page": 1,
    "limit": 20
  }
}
```

#### PUT /api/admin/users/:id/suspend
**Description:** Suspend a user account

**Permissions:** `user_management`

**Request:**
```json
{
  "reason": "Violation of community guidelines",
  "duration_days": 7,
  "notes": "First offense"
}
```

**Response:**
```json
{
  "data": {
    "success": true,
    "message": "User suspended successfully",
    "suspended_until": "2024-01-22T00:00:00Z"
  }
}
```

#### PUT /api/admin/users/:id/ban
**Description:** Permanently ban a user account

**Permissions:** `user_management`

**Request:**
```json
{
  "reason": "Repeated violations",
  "notes": "Multiple warnings issued"
}
```

#### PUT /api/admin/users/:id/unsuspend
**Description:** Remove suspension from user account

**Permissions:** `user_management`

#### PUT /api/admin/users/:id/reset-password
**Description:** Force password reset for user

**Permissions:** `user_management`

**Response:**
```json
{
  "data": {
    "reset_token": "temporary_token",
    "expires_at": "2024-01-16T00:00:00Z"
  }
}
```

---

### 3. Content Moderation

#### GET /api/admin/reports
**Description:** Get all content reports

**Permissions:** `content_management`

**Query Parameters:**
- `status`: pending, resolved, escalated
- `content_type`: post, comment, message
- `page`, `limit`

**Response:**
```json
{
  "data": {
    "reports": [{
      "id": "uuid",
      "content_type": "post",
      "content_id": "uuid",
      "reported_by": {
        "id": "uuid",
        "full_name": "Reporter Name"
      },
      "reason": "spam",
      "description": "This is spam content",
      "status": "pending",
      "created_at": "2024-01-15T10:00:00Z"
    }],
    "total": 50
  }
}
```

#### PUT /api/admin/reports/:id/resolve
**Description:** Resolve a content report

**Permissions:** `content_management`

**Request:**
```json
{
  "action": "remove_content",
  "notes": "Content removed due to spam",
  "notify_user": true
}
```

#### PUT /api/admin/reports/:id/escalate
**Description:** Escalate report to higher admin level

**Permissions:** `content_management`

#### DELETE /api/admin/content/posts/:id
**Description:** Admin delete any post

**Permissions:** `content_management`

**Request:**
```json
{
  "reason": "Violates community guidelines",
  "notify_user": true
}
```

---

### 4. Announcements Management

#### GET /api/admin/announcements
**Description:** Get all announcements (including drafts)

**Permissions:** `content_management`

#### POST /api/admin/announcements
**Description:** Create new announcement

**Permissions:** `content_management`

**Request:**
```json
{
  "space_id": "uuid",
  "title": "Important Update",
  "content": "Announcement content",
  "type": "general",
  "target_audience": ["students", "staff"],
  "priority": "high",
  "status": "published",
  "scheduled_for": "2024-01-20T09:00:00Z",
  "expires_at": "2024-01-30T00:00:00Z"
}
```

#### PUT /api/admin/announcements/:id
**Description:** Update announcement

#### DELETE /api/admin/announcements/:id
**Description:** Delete announcement

---

### 5. Events Management

#### GET /api/admin/events
**Description:** Get all events (including pending approval)

**Permissions:** `content_management`

#### POST /api/admin/events
**Description:** Create new event

**Permissions:** `content_management`

**Request:**
```json
{
  "space_id": "uuid",
  "title": "Campus Festival",
  "description": "Annual campus festival",
  "location": "Main Campus",
  "start_time": "2024-02-01T10:00:00Z",
  "end_time": "2024-02-01T18:00:00Z",
  "capacity": 500,
  "registration_required": true,
  "tags": ["festival", "entertainment"]
}
```

#### PUT /api/admin/events/:id
**Description:** Update event

#### DELETE /api/admin/events/:id
**Description:** Delete event

#### PUT /api/admin/events/:id/approve
**Description:** Approve pending event

#### PUT /api/admin/events/:id/reject
**Description:** Reject pending event

---

### 6. Communities & Groups Management

#### GET /api/admin/communities
**Description:** Get all communities

**Permissions:** `community_management`

#### PUT /api/admin/communities/:id/suspend
**Description:** Suspend a community

**Permissions:** `community_management`

**Request:**
```json
{
  "reason": "Inappropriate content",
  "duration_days": 7
}
```

#### GET /api/admin/groups
**Description:** Get all groups with status filter

**Permissions:** `community_management`

**Query Parameters:**
- `status`: pending, approved, rejected

#### PUT /api/admin/groups/:id/approve
**Description:** Approve pending group

**Permissions:** `community_management`

#### PUT /api/admin/groups/:id/reject
**Description:** Reject pending group

**Permissions:** `community_management`

**Request:**
```json
{
  "reason": "Does not meet guidelines"
}
```

#### PUT /api/admin/groups/:id/suspend
**Description:** Suspend a group

---

### 7. Tutoring & Mentorship Management

#### GET /api/admin/applications/tutors
**Description:** Get tutor applications

**Permissions:** `tutoring_management`

**Query Parameters:**
- `status`: pending, approved, rejected

**Response:**
```json
{
  "data": {
    "applications": [{
      "id": "uuid",
      "user": {
        "id": "uuid",
        "full_name": "Jane Doe",
        "email": "jane@university.edu"
      },
      "subjects": ["Mathematics", "Physics"],
      "experience": "2 years tutoring experience",
      "qualifications": "BSc Mathematics",
      "status": "pending",
      "created_at": "2024-01-10T00:00:00Z"
    }]
  }
}
```

#### PUT /api/admin/applications/tutors/:id/approve
**Description:** Approve tutor application

**Permissions:** `tutoring_management`

**Request:**
```json
{
  "notes": "Qualified applicant"
}
```

#### PUT /api/admin/applications/tutors/:id/reject
**Description:** Reject tutor application

**Request:**
```json
{
  "reason": "Insufficient qualifications",
  "notes": "Please reapply after gaining more experience"
}
```

#### GET /api/admin/applications/mentors
**Description:** Get mentor applications

#### PUT /api/admin/applications/mentors/:id/approve
**Description:** Approve mentor application

#### PUT /api/admin/applications/mentors/:id/reject
**Description:** Reject mentor application

---

### 8. Analytics & Reports

#### GET /api/admin/analytics/overview
**Description:** Get system-wide analytics overview

**Permissions:** `analytics`

**Response:**
```json
{
  "data": {
    "total_users": 5000,
    "active_users": 3500,
    "new_users_this_month": 250,
    "total_posts": 15000,
    "total_communities": 45,
    "total_groups": 120,
    "engagement_rate": 0.72
  }
}
```

#### GET /api/admin/analytics/users
**Description:** Get user analytics

**Permissions:** `analytics`

**Query Parameters:**
- `period`: day, week, month, year
- `start_date`, `end_date`

**Response:**
```json
{
  "data": {
    "user_growth": [{
      "date": "2024-01-01",
      "new_users": 15,
      "active_users": 1200
    }],
    "user_retention": 0.85,
    "top_users": [{
      "id": "uuid",
      "full_name": "John Doe",
      "posts_count": 150,
      "engagement_score": 95
    }]
  }
}
```

#### GET /api/admin/analytics/content
**Description:** Get content analytics

**Permissions:** `analytics`

#### GET /api/admin/analytics/engagement
**Description:** Get engagement metrics

**Permissions:** `analytics`

#### POST /api/admin/analytics/export
**Description:** Export analytics data

**Permissions:** `analytics`

**Request:**
```json
{
  "data_type": "users",
  "format": "csv",
  "date_range": {
    "start": "2024-01-01",
    "end": "2024-01-31"
  }
}
```

---

### 9. Admin User Management (Super Admin Only)

#### GET /api/admin/admins
**Description:** Get all admin users

**Permissions:** `admin_management` (super_admin only)

**Response:**
```json
{
  "data": {
    "admins": [{
      "id": "uuid",
      "email": "admin@university.edu",
      "full_name": "Admin User",
      "role": "admin",
      "permissions": ["user_management", "content_management"],
      "status": "active",
      "created_at": "2024-01-01T00:00:00Z",
      "last_login": "2024-01-15T10:00:00Z"
    }]
  }
}
```

#### POST /api/admin/admins
**Description:** Create new admin user

**Permissions:** `admin_management` (super_admin only)

**Request:**
```json
{
  "email": "newadmin@university.edu",
  "full_name": "New Admin",
  "password": "securepassword",
  "role": "admin",
  "permissions": ["user_management", "content_management"]
}
```

#### PUT /api/admin/admins/:id
**Description:** Update admin user

**Permissions:** `admin_management` (super_admin only)

#### PUT /api/admin/admins/:id/deactivate
**Description:** Deactivate admin user

**Permissions:** `admin_management` (super_admin only)

---

### 10. System Settings (Super Admin Only)

#### GET /api/admin/settings
**Description:** Get system settings

**Permissions:** `system_settings` (super_admin only)

**Response:**
```json
{
  "data": {
    "maintenance_mode": false,
    "registration_enabled": true,
    "email_verification_required": true,
    "max_upload_size_mb": 10,
    "session_timeout_minutes": 60,
    "password_min_length": 8,
    "system_notice": "Welcome to Connect!"
  }
}
```

#### PUT /api/admin/settings
**Description:** Update system settings

**Permissions:** `system_settings` (super_admin only)

**Request:**
```json
{
  "maintenance_mode": true,
  "registration_enabled": false,
  "system_notice": "System maintenance in progress"
}
```

---

### 11. Space Activities (NEW)

#### GET /api/admin/spaces
**Description:** Get all spaces with activity metrics

**Permissions:** `admin` or `super_admin`

**Response:**
```json
{
  "data": {
    "spaces": [{
      "id": "uuid",
      "name": "University of Technology",
      "slug": "utech",
      "description": "Main campus space",
      "created_at": "2024-01-01T00:00:00Z",
      "stats": {
        "total_users": 5000,
        "active_users": 3500,
        "total_posts": 15000,
        "total_communities": 45,
        "total_groups": 120
      }
    }]
  }
}
```

#### GET /api/admin/spaces/:id/activities
**Description:** Get recent activities in a space

**Permissions:** `admin` or `super_admin`

**Query Parameters:**
- `activity_type`: user_joined, post_created, community_created, etc.
- `limit`: default 50
- `offset`: default 0

**Response:**
```json
{
  "data": {
    "activities": [{
      "id": "uuid",
      "type": "user_joined",
      "actor": {
        "id": "uuid",
        "full_name": "John Doe"
      },
      "description": "New user joined the space",
      "metadata": {},
      "created_at": "2024-01-15T10:00:00Z"
    }]
  }
}
```

#### POST /api/admin/spaces
**Description:** Create new space

**Permissions:** `super_admin`

**Request:**
```json
{
  "name": "New University",
  "slug": "new-uni",
  "description": "Description of the space"
}
```

#### PUT /api/admin/spaces/:id
**Description:** Update space

**Permissions:** `super_admin`

#### DELETE /api/admin/spaces/:id
**Description:** Delete space (soft delete)

**Permissions:** `super_admin`

---

## Database Schema Updates Needed

### New Tables

#### admin_users
```sql
CREATE TABLE admin_users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) UNIQUE NOT NULL,
  role VARCHAR(50) NOT NULL, -- 'admin', 'super_admin'
  permissions TEXT[], -- Array of permissions
  status VARCHAR(20) DEFAULT 'active', -- 'active', 'inactive'
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

#### user_suspensions
```sql
CREATE TABLE user_suspensions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) NOT NULL,
  suspended_by UUID REFERENCES admin_users(id) NOT NULL,
  reason TEXT NOT NULL,
  notes TEXT,
  suspended_at TIMESTAMP DEFAULT NOW(),
  suspended_until TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW()
);
```

#### content_reports
```sql
CREATE TABLE content_reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  space_id UUID REFERENCES spaces(id) NOT NULL,
  content_type VARCHAR(20) NOT NULL, -- 'post', 'comment', 'message'
  content_id UUID NOT NULL,
  reported_by UUID REFERENCES users(id) NOT NULL,
  reason VARCHAR(100) NOT NULL,
  description TEXT,
  status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'resolved', 'escalated'
  resolved_by UUID REFERENCES admin_users(id),
  resolution_notes TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  resolved_at TIMESTAMP
);
```

#### system_settings
```sql
CREATE TABLE system_settings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  key VARCHAR(100) UNIQUE NOT NULL,
  value JSONB NOT NULL,
  updated_by UUID REFERENCES admin_users(id),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

#### audit_logs
```sql
CREATE TABLE audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  admin_user_id UUID REFERENCES admin_users(id) NOT NULL,
  action VARCHAR(100) NOT NULL,
  resource_type VARCHAR(50) NOT NULL,
  resource_id UUID,
  details JSONB,
  ip_address INET,
  created_at TIMESTAMP DEFAULT NOW()
);
```

#### space_activities
```sql
CREATE TABLE space_activities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  space_id UUID REFERENCES spaces(id) NOT NULL,
  activity_type VARCHAR(50) NOT NULL,
  actor_id UUID REFERENCES users(id),
  description TEXT NOT NULL,
  metadata JSONB,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_space_activities_space_id ON space_activities(space_id);
CREATE INDEX idx_space_activities_created_at ON space_activities(created_at);
```

### Schema Modifications

#### users table
Add columns:
```sql
ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'active';
ALTER TABLE users ADD COLUMN IF NOT EXISTS suspended_until TIMESTAMP;
```

## Error Codes

| Code | Message | HTTP Status |
|------|---------|-------------|
| unauthorized | Not authorized | 401 |
| forbidden | Insufficient permissions | 403 |
| not_found | Resource not found | 404 |
| validation_error | Validation failed | 400 |
| already_exists | Resource already exists | 409 |
| internal_error | Internal server error | 500 |

## Implementation Priority

1. **High Priority** (Core functionality)
   - Admin authentication
   - User management (suspend, ban)
   - Content reports
   - Space activities

2. **Medium Priority** (Extended features)
   - Announcements CRUD
   - Events CRUD
   - Tutor/Mentor approval
   - Analytics overview

3. **Low Priority** (Nice to have)
   - Advanced analytics
   - Data export
   - Admin management
   - System settings

## Testing Requirements

Each endpoint must have:
- Unit tests for service layer
- Integration tests for handlers
- Authentication/authorization tests
- Input validation tests
- Error handling tests

Minimum test coverage: 80%
