# Admin Dashboard - Backend Integration Map

This document tracks all admin dashboard pages, their required endpoints, and implementation status.

## Admin Dashboard Pages

### 1. User Management (`/admin/users`)

**Current Status:** Partially Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/users` | GET | ✓ Implemented | Get paginated list of users |
| `/api/admin/users/:id` | DELETE | ✓ Implemented | Delete a user |
| `/api/admin/users/:id/suspend` | PUT | ✓ Implemented | Suspend a user |
| `/api/admin/users/:id/unsuspend` | PUT | ✓ Implemented | Unsuspend a user |
| `/api/admin/users/:id/ban` | PUT | ✓ Implemented | Ban a user |
| `/api/admin/users/:id` | PUT | ✗ Missing | Update user details |
| `/api/admin/users/:id/roles` | PUT | ✗ Missing | Update user roles |

**Frontend Integration Status:** ✓ Real API (admin.api.ts) - Integrated

---

### 2. Tutoring & Mentorship (`/admin/tutoring-mentorship`)

**Current Status:** Partially Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/applications/tutors` | GET | ✓ Implemented | Get tutor applications |
| `/api/admin/applications/tutors/:id/approve` | PUT | ✓ Implemented | Approve tutor application |
| `/api/admin/applications/tutors/:id/reject` | PUT | ✓ Implemented | Reject tutor application |
| `/api/admin/applications/mentors` | GET | ✓ Implemented | Get mentor applications |
| `/api/admin/applications/mentors/:id/approve` | PUT | ✓ Implemented | Approve mentor application |
| `/api/admin/applications/mentors/:id/reject` | PUT | ✓ Implemented | Reject mentor application |

**Frontend Integration Status:** ✓ Real API (admin.api.ts) - Integrated

---

### 3. Communities & Groups (`/admin/communities-groups`)

**Current Status:** Partially Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/groups` | GET | ✓ Implemented | Get groups list |
| `/api/admin/groups/:id/approve` | PUT | ✓ Implemented | Approve a group |
| `/api/admin/groups/:id/reject` | PUT | ✓ Implemented | Reject a group |
| `/api/admin/groups/:id` | DELETE | ✓ Implemented | Delete a group |
| `/api/communities` | GET | ✓ Implemented | Get communities list |
| `/api/communities/:id` | PUT | ✗ Missing | Update community |
| `/api/communities/:id` | DELETE | ✗ Missing | Delete community |

**Frontend Integration Status:** ✓ Real API (admin.api.ts) - Integrated

---

### 4. Reports (`/admin/reports`)

**Current Status:** Partially Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/reports` | GET | ✓ Implemented | Get content reports |
| `/api/admin/reports/:id/resolve` | PUT | ✓ Implemented | Resolve a report |
| `/api/admin/reports/:id/escalate` | PUT | ✓ Implemented | Escalate a report |

**Frontend Integration Status:** ✓ Real API (admin.api.ts) - Integrated

---

### 5. Content Management (`/admin/content`)

**Current Status:** Not Implemented ✗

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/events` | GET | ✓ Implemented | Get events list |
| `/api/events` | POST | ✓ Implemented | Create event |
| `/api/events/:id` | PUT | ✓ Implemented | Update event |
| `/api/events/:id` | DELETE | ✓ Implemented | Delete event |
| `/api/announcements` | GET | ✗ Missing | Get announcements |
| `/api/announcements` | POST | ✗ Missing | Create announcement |
| `/api/announcements/:id` | PUT | ✗ Missing | Update announcement |
| `/api/announcements/:id` | DELETE | ✗ Missing | Delete announcement |
| `/api/admin/posts` | GET | ✗ Missing | Get all posts for moderation |
| `/api/admin/posts/:id` | DELETE | ✗ Missing | Delete post |
| `/api/admin/posts/:id/hide` | PUT | ✗ Missing | Hide/unhide post |

**Frontend Integration Status:** ✓ Real API (admin.api.ts) - Integrated

---

### 6. Space Activities (`/admin/activities`)

**Current Status:** Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/spaces/:id/activities` | GET | ✓ Implemented | Get space activities log |

**Frontend Integration Status:** Real API (admin.api.ts) ✓

---

### 7. Dashboard Stats (`/admin`)

**Current Status:** Implemented ✓

**Required Endpoints:**

| Endpoint | Method | Status | Notes |
|----------|--------|--------|-------|
| `/api/admin/dashboard/stats` | GET | ✓ Implemented | Get dashboard statistics |
| `/api/analytics/metrics/space` | GET | ✓ Implemented | Get space metrics |

**Frontend Integration Status:** Real API (admin.api.ts) ✓

---

## Implementation Status Summary

### ✓ Completed (High Priority Core Functionality)
1. ✅ **User Management** - Integrated with real API (admin.api.ts)
   - List, suspend, unsuspend, ban, delete users working
   - Create/edit user disabled (backend not implemented)

2. ✅ **Tutoring & Mentorship** - Integrated with real API (admin.api.ts)
   - Tutor applications: list, approve, reject working
   - Mentor applications: list, approve, reject working

3. ✅ **Communities & Groups** - Partially integrated with real API (admin.api.ts)
   - Groups tab: list, approve, reject, delete working
   - Communities tab: still using mock data (backend not implemented)
   - Suspend/reactivate disabled (backend not implemented)

4. ✅ **Reports** - Integrated with real API (admin.api.ts)
   - List, approve, reject, warn user working
   - Resolve and escalate working

### 🚧 Pending (Medium Priority Content Management)
5. ⏳ **Announcements CRUD** - Backend endpoints not implemented
6. ⏳ **Post moderation** - Backend endpoints not implemented

### ✅ Already Working (Low Priority)
7. ✅ **Space Activities** - Already using real API ✓
8. ✅ **Dashboard Stats** - Already using real API ✓

---

## Missing Backend Endpoints to Implement

### 1. User Management
- `PUT /api/admin/users/:id` - Update user details (name, email, etc.)
- `PUT /api/admin/users/:id/roles` - Update user roles

### 2. Content Management
- `GET /api/announcements` - List announcements
- `POST /api/announcements` - Create announcement
- `PUT /api/announcements/:id` - Update announcement
- `DELETE /api/announcements/:id` - Delete announcement
- `GET /api/admin/posts` - List all posts with filters
- `DELETE /api/admin/posts/:id` - Delete post as admin
- `PUT /api/admin/posts/:id/hide` - Hide/unhide post

### 3. Communities
- `PUT /api/communities/:id` - Update community
- `DELETE /api/communities/:id` - Delete community

---

## Frontend TypeScript Types to Sync

Types that need to match backend Go structs:

1. `User` - Match with backend User model
2. `TutorApplication` - Match with backend TutorApplication
3. `MentorApplication` - Match with backend MentorApplication
4. `Group` - Match with backend Group
5. `Community` - Match with backend Community
6. `ContentReport` - Match with backend Report
7. `Announcement` - New type needed
8. `Post` - Match with backend Post

---

## Testing Requirements

Each new endpoint needs:
- Unit tests for handler functions
- Integration tests with database
- Error case handling tests
- Authorization/permission tests
