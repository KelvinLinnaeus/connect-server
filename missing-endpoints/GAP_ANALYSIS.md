# Admin Dashboard - Frontend/Backend Gap Analysis

## Executive Summary

**Total Frontend Admin Sections:** 10
**Existing Backend Endpoints:** 18 (across 12 routes)
**Required New Endpoints:** ~35
**Integration Status:** ~40% Complete

---

## Detailed Gap Analysis by Section

### 1. User Management ✅ MOSTLY COMPLETE

**Frontend Requirements:**
- User listing with pagination ✅
- Search by username/email ⚠️ (exists but not used)
- Filter by role (student, tutor, mentor, lecturer) ⚠️ (needs enhancement)
- Suspend user ✅
- Unsuspend user ✅
- Ban user ✅
- Delete user ✅
- Reset password ❌ MISSING
- View user activity ❌ MISSING
- Bulk actions (suspend, activate, delete) ❌ MISSING
- Export users as CSV ❌ MISSING (frontend calls it but backend missing)
- User statistics cards ✅ (via dashboard stats)
- Add user modal ❌ MISSING
- Edit user modal ❌ MISSING

**Existing Backend Endpoints:**
- ✅ GET /api/admin/users (paginated list)
- ✅ DELETE /api/admin/users/:id
- ✅ PUT /api/admin/users/:id/suspend
- ✅ PUT /api/admin/users/:id/unsuspend
- ✅ PUT /api/admin/users/:id/ban

**Missing Backend Endpoints:**
- ❌ POST /api/admin/users (create user)
- ❌ PUT /api/admin/users/:id (update user)
- ❌ PUT /api/admin/users/:id/reset-password
- ❌ GET /api/admin/users/:id/activity
- ❌ POST /api/admin/users/bulk-action (bulk operations)
- ❌ GET /api/admin/users/export (CSV export)

**Status:** 60% Complete

---

### 2. Content Management ⚠️ PARTIAL

#### Sub-section: Content Moderation ✅ COMPLETE
**Frontend Requirements:**
- List reported content with filters ✅
- Approve content ✅
- Reject content ✅
- Warn user ✅
- Bulk review ⚠️ (frontend has it but not connected)

**Existing Backend Endpoints:**
- ✅ GET /api/admin/reports
- ✅ PUT /api/admin/reports/:id/resolve
- ✅ PUT /api/admin/reports/:id/escalate

**Missing:**
- ❌ POST /api/admin/reports/bulk-resolve

**Status:** 90% Complete

#### Sub-section: Announcements ❌ COMPLETELY MISSING
**Frontend Requirements:**
- List announcements ❌
- Create announcement ❌
- Edit announcement ❌
- Delete announcement ❌
- Publish/unpublish ❌
- Schedule announcement ❌

**Missing Backend Endpoints:**
- ❌ GET /api/admin/announcements
- ❌ POST /api/admin/announcements
- ❌ PUT /api/admin/announcements/:id
- ❌ DELETE /api/admin/announcements/:id
- ❌ PUT /api/admin/announcements/:id/publish
- ❌ PUT /api/admin/announcements/:id/schedule

**Status:** 0% Complete

#### Sub-section: Events ❌ COMPLETELY MISSING
**Frontend Requirements:**
- List events ❌
- Create event ❌
- Edit event ❌
- Delete event ❌
- Cancel event ❌
- Manage registrations ❌

**Missing Backend Endpoints:**
- ❌ GET /api/admin/events
- ❌ POST /api/admin/events
- ❌ PUT /api/admin/events/:id
- ❌ DELETE /api/admin/events/:id
- ❌ PUT /api/admin/events/:id/cancel
- ❌ GET /api/admin/events/:id/registrations

**Status:** 0% Complete

**Overall Content Management Status:** 30% Complete

---

### 3. Communities & Groups ⚠️ PARTIAL

#### Sub-section: Communities Management ❌ MOSTLY MISSING
**Frontend Requirements:**
- List communities by category ❌
- Create community ❌
- Update community ❌
- Delete community ❌
- Assign moderators ❌
- Export communities ❌

**Current State:**
- Uses `mockCommunitiesApi` with hardcoded data
- Categories: Academic, Department, Level, Hostel, Faculty

**Missing Backend Endpoints:**
- ❌ GET /api/admin/communities
- ❌ POST /api/admin/communities
- ❌ PUT /api/admin/communities/:id
- ❌ DELETE /api/admin/communities/:id
- ❌ PUT /api/admin/communities/:id/assign-moderator
- ❌ GET /api/admin/communities/export

**Status:** 0% Complete

#### Sub-section: Groups Management ✅ MOSTLY COMPLETE
**Frontend Requirements:**
- List groups with filters ✅
- Approve group ✅
- Reject group ✅
- Delete group ✅
- Suspend group ⚠️ (frontend has dialog but uses approve/reject)
- Reactivate group ⚠️ (same as suspend)
- View group details ⚠️ (needs full group info endpoint)
- Export groups ❌

**Existing Backend Endpoints:**
- ✅ GET /api/admin/groups
- ✅ PUT /api/admin/groups/:id/approve
- ✅ PUT /api/admin/groups/:id/reject
- ✅ DELETE /api/admin/groups/:id

**Missing Backend Endpoints:**
- ❌ GET /api/admin/groups/:id (full details)
- ❌ PUT /api/admin/groups/:id/suspend
- ❌ PUT /api/admin/groups/:id/reactivate
- ❌ GET /api/admin/groups/export

**Status:** 75% Complete

**Overall Communities & Groups Status:** 40% Complete

---

### 4. Tutoring & Mentorship ✅ COMPLETE

**Frontend Requirements:**
- List tutor applications ✅
- Approve tutor application ✅
- Reject tutor application ✅
- List mentor applications ✅
- Approve mentor application ✅
- Reject mentor application ✅
- Statistics (Total, Pending, Approved, Rejected, Active) ✅

**Existing Backend Endpoints:**
- ✅ GET /api/admin/applications/tutors
- ✅ PUT /api/admin/applications/tutors/:id/approve
- ✅ PUT /api/admin/applications/tutors/:id/reject
- ✅ GET /api/admin/applications/mentors
- ✅ PUT /api/admin/applications/mentors/:id/approve
- ✅ PUT /api/admin/applications/mentors/:id/reject

**Status:** 100% Complete ✅

---

### 5. Analytics & Reports (Content Moderation Queue)

#### Analytics Tab ❌ COMPLETELY MISSING REAL DATA
**Frontend Requirements:**
- User growth trends (6 months) ❌
- Engagement metrics (posts, comments, shares) ❌
- Group activity analysis ❌
- Category distribution ❌
- Time range filters ❌

**Current State:**
- Uses completely mock data in `Analytics.tsx` (lines 22-55)
- No backend integration

**Missing Backend Endpoints:**
- ❌ GET /api/admin/analytics/user-growth
- ❌ GET /api/admin/analytics/engagement
- ❌ GET /api/admin/analytics/group-activity
- ❌ GET /api/admin/analytics/category-distribution

**Note:** Backend has `GetAdminDashboardStats`, `GetUserGrowthData`, `GetContentGrowthData` in SQLC but no exposed HTTP endpoints

**Status:** 0% Complete

#### Reports Tab ✅ COMPLETE
**Frontend Requirements:**
- List content reports ✅
- Filter by type, priority, status ✅
- Resolve report ✅
- View full post details ⚠️ (needs post content endpoint)
- Warn user ⚠️ (uses resolve with action)

**Existing Backend Endpoints:**
- ✅ GET /api/admin/reports
- ✅ PUT /api/admin/reports/:id/resolve

**Status:** 90% Complete

**Overall Analytics & Reports Status:** 45% Complete

---

### 6. Space Activities ✅ COMPLETE

**Frontend Requirements:**
- List space activities ✅
- Filter by activity type ✅
- Space stats (users, posts, communities, groups) ✅
- Real-time activity stream ✅

**Existing Backend Endpoints:**
- ✅ GET /api/admin/spaces/:id/activities
- ✅ GET /api/admin/dashboard/stats (provides space stats)

**Status:** 100% Complete ✅

---

### 7. Notifications ❌ COMPLETELY MISSING

**Frontend Requirements:**
- List notifications by type ❌
- Mark notification as read ❌
- Delete notification ❌
- Filter by priority ❌
- Notification types: report, application, system, user_activity ❌

**Current State:**
- Uses `mockNotifications` array (6 hardcoded items)
- No backend integration

**Missing Backend Endpoints:**
- ❌ GET /api/admin/notifications
- ❌ PUT /api/admin/notifications/:id/read
- ❌ DELETE /api/admin/notifications/:id
- ❌ PUT /api/admin/notifications/read-all

**Status:** 0% Complete

---

### 8. Admin Management (Super Admin Only) ❌ COMPLETELY MISSING

**Frontend Requirements:**
- List admins ❌
- Create admin ❌
- Edit admin ❌
- Delete admin ❌
- Manage permissions ❌
- Activate/deactivate admin ❌

**Current State:**
- No backend integration (frontend exists but disconnected)

**Missing Backend Endpoints:**
- ❌ GET /api/admin/admins
- ❌ POST /api/admin/admins
- ❌ PUT /api/admin/admins/:id
- ❌ DELETE /api/admin/admins/:id
- ❌ PUT /api/admin/admins/:id/permissions
- ❌ PUT /api/admin/admins/:id/status

**Note:** Backend has `GetAllAdminUsers` SQLC query but no HTTP endpoint

**Status:** 0% Complete

---

### 9. System Settings (Super Admin Only) ❌ COMPLETELY MISSING

**Frontend Requirements:**
- Get all settings ❌
- Update settings ❌
- Categories: general, user, community, email, notification, theme, security ❌

**Current State:**
- Uses `mockSettings` object with 7+ setting categories
- No backend integration

**Missing Backend Endpoints:**
- ❌ GET /api/admin/settings
- ❌ PUT /api/admin/settings/:key

**Note:** Backend has `system_settings` table and SQLC queries (`GetSystemSetting`, `UpsertSystemSetting`) but no HTTP endpoints

**Status:** 0% Complete

---

### 10. Dashboard (Main) ⚠️ PARTIAL

**Frontend Requirements:**
- Overview statistics ✅
- Recent activity ✅
- Quick actions ⚠️

**Existing Backend Endpoints:**
- ✅ GET /api/admin/dashboard/stats (7 metrics)
- ✅ GET /api/admin/spaces/:id/activities (recent activity)

**Status:** 80% Complete

---

## Priority Implementation Plan

### Phase 1: High Priority (Required for Basic Functionality)
1. **User Management Completion**
   - POST /api/admin/users (create user)
   - PUT /api/admin/users/:id (edit user)
   - PUT /api/admin/users/:id/reset-password
   - POST /api/admin/users/bulk-action

2. **Communities Management (CRITICAL - Currently 100% Mock)**
   - GET /api/admin/communities
   - POST /api/admin/communities
   - PUT /api/admin/communities/:id
   - DELETE /api/admin/communities/:id
   - PUT /api/admin/communities/:id/assign-moderator

3. **Admin Management (Super Admin Feature)**
   - GET /api/admin/admins
   - POST /api/admin/admins
   - PUT /api/admin/admins/:id
   - DELETE /api/admin/admins/:id
   - PUT /api/admin/admins/:id/permissions

4. **System Settings**
   - GET /api/admin/settings
   - PUT /api/admin/settings/:key

### Phase 2: Medium Priority (Enhance Functionality)
1. **Announcements**
   - Full CRUD endpoints for announcements
   - Publish/schedule functionality

2. **Events**
   - Full CRUD endpoints for events
   - Registration management

3. **Analytics**
   - Expose existing SQLC analytics queries as HTTP endpoints
   - User growth, engagement, category distribution

4. **Notifications**
   - Admin notification system
   - Read/unread tracking

### Phase 3: Low Priority (Nice to Have)
1. **Export Functionality**
   - CSV export for users, groups, communities

2. **Bulk Operations**
   - Bulk content moderation

3. **Enhanced Monitoring**
   - User activity logs
   - Detailed post content viewing

---

## Summary Statistics

| Section | Completion | Missing Endpoints | Priority |
|---------|-----------|------------------|----------|
| User Management | 60% | 6 | High |
| Content Moderation | 90% | 1 | Low |
| Announcements | 0% | 6 | Medium |
| Events | 0% | 6 | Medium |
| Communities | 0% | 6 | **CRITICAL** |
| Groups | 75% | 4 | Medium |
| Tutoring/Mentorship | 100% | 0 | ✅ |
| Analytics | 0% | 4 | Medium |
| Reports | 90% | 0 | Low |
| Space Activities | 100% | 0 | ✅ |
| Notifications | 0% | 4 | Medium |
| Admin Management | 0% | 6 | High |
| System Settings | 0% | 2 | High |

**Total Missing Endpoints: ~45**
**Estimated Implementation Time: 3-5 days for Phase 1**

---

## Type Alignment Issues

### Frontend TypeScript → Backend Go Type Mismatches

1. **User Roles**
   - Frontend: `"student" | "tutor" | "mentor" | "lecturer"`
   - Backend: `TEXT[]` array in users table
   - **Action:** Ensure role validation in backend

2. **Group Status**
   - Frontend: `"active" | "inactive" | "suspended" | "pending"`
   - Backend: May differ - needs verification
   - **Action:** Align status enums

3. **Community Categories**
   - Frontend: `"Academic" | "Department" | "Level" | "Hostel" | "Faculty"`
   - Backend: Not yet defined
   - **Action:** Define category enum in backend

4. **Report Priority**
   - Frontend: `"urgent" | "high" | "medium" | "low"`
   - Backend: `VARCHAR(20)` with same values ✅
   - **Status:** ALIGNED

5. **Notification Types**
   - Frontend: `"report" | "application" | "system" | "user_activity"`
   - Backend: Not yet implemented
   - **Action:** Define notification type enum

---

## Database Schema Gaps

### Tables Needed but Missing:

1. **announcements**
   - id, space_id, title, content, author_id
   - status (draft, published, scheduled, expired)
   - priority, scheduled_for, expires_at
   - created_at, updated_at

2. **events**
   - id, space_id, title, description, author_id
   - start_date, end_date, location
   - max_attendees, registration_required
   - status (draft, published, cancelled, completed)
   - created_at, updated_at

3. **event_registrations**
   - id, event_id, user_id
   - status (registered, cancelled, attended)
   - registered_at

4. **admin_notifications**
   - id, admin_user_id, type, title, message
   - priority, is_read, related_resource_type, related_resource_id
   - created_at

5. **communities** (if not exists)
   - id, space_id, name, description
   - category, moderator_ids
   - member_count, post_count
   - created_at, updated_at

### Tables That Exist but Need Enhancement:

1. **users**
   - ✅ Already has: roles[], status, suspended_until
   - Needs: Ensure proper indexing on roles

2. **groups**
   - Need to verify: Has status field?
   - May need: suspended_until, suspension_reason

---

## Frontend API Client Issues

### Files to Update:

1. **frontend/src/services/adminApi.ts**
   - Currently has partial implementations
   - Many functions return mock data or are placeholders
   - Needs: Complete implementation for all endpoints

2. **frontend/src/data/mockAdminCommunitiesData.ts**
   - Contains mock API implementations
   - **Action:** Replace with real API calls

3. **frontend/src/pages/admin/Analytics.tsx**
   - Uses inline mock data
   - **Action:** Connect to real analytics endpoints

4. **frontend/src/pages/admin/Notifications.tsx**
   - Uses inline mock notifications
   - **Action:** Connect to real notifications API

5. **frontend/src/pages/admin/SystemSettings.tsx**
   - Uses mock settings object
   - **Action:** Connect to real settings API

---

## Authentication & Authorization Gaps

### Current State:
- ✅ JWT authentication on all admin routes
- ✅ Middleware for role checking exists
- ❌ Not applied to routes yet

### Required Actions:
1. Apply `RequireAdmin` middleware to all admin routes
2. Apply `RequireRole("super_admin")` to:
   - Admin Management endpoints
   - System Settings endpoints
3. Verify frontend permission checks match backend enforcement

---

## Next Steps

1. ✅ Complete this gap analysis
2. **Create migration files for missing tables**
3. **Implement missing SQLC queries**
4. **Create backend handlers/services for missing endpoints**
5. **Update frontend API client**
6. **Connect frontend components to real APIs**
7. **Add loading states and error handling**
8. **Write backend tests**
9. **Validate frontend build**
10. **Validate backend build**

---

**Last Updated:** 2025-11-06
**Status:** Ready for Implementation
