# Admin Dashboard API Implementation Plan

## Overview
This document outlines the implementation plan for all missing admin dashboard endpoints to replace mock data with live API integration.

## Implementation Status

### ✅ Already Implemented (18 endpoints)
1. GET /api/admin/users
2. DELETE /api/admin/users/:id
3. PUT /api/admin/users/:id/suspend
4. PUT /api/admin/users/:id/unsuspend
5. PUT /api/admin/users/:id/ban
6. GET /api/admin/reports
7. PUT /api/admin/reports/:id/resolve
8. PUT /api/admin/reports/:id/escalate
9. GET /api/admin/applications/tutors
10. PUT /api/admin/applications/tutors/:id/approve
11. PUT /api/admin/applications/tutors/:id/reject
12. GET /api/admin/applications/mentors
13. PUT /api/admin/applications/mentors/:id/approve
14. PUT /api/admin/applications/mentors/:id/reject
15. GET /api/admin/groups
16. PUT /api/admin/groups/:id/approve
17. PUT /api/admin/groups/:id/reject
18. DELETE /api/admin/groups/:id
19. GET /api/admin/spaces/:id/activities
20. GET /api/admin/dashboard/stats

### 🔨 To Be Implemented (Priority 1 - Critical)

#### System Settings (2 endpoints)
**Status:** SQLC queries exist, need handlers

1. **GET /api/admin/settings**
   - Query: `GetAllSystemSettings` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: SystemSettings.tsx (currently uses mock)

2. **PUT /api/admin/settings/:key**
   - Query: `UpsertSystemSetting` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: SystemSettings.tsx (currently uses mock)

#### Admin Management (6 endpoints)
**Status:** SQLC queries exist, need handlers

3. **GET /api/admin/admins**
   - Query: `GetAllAdminUsers` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx (currently disconnected)

4. **POST /api/admin/admins**
   - Query: `AdminCreateUser` ✅ added (needs sqlc generate)
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx

5. **PUT /api/admin/admins/:id**
   - Query: `AdminUpdateUser` ✅ added (needs sqlc generate)
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx

6. **DELETE /api/admin/admins/:id**
   - Query: `DeleteUser` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx

7. **PUT /api/admin/admins/:id/permissions**
   - Query: `UpdateUserRole` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx

8. **PUT /api/admin/admins/:id/status**
   - Query: `UpdateUserAccountStatus` ✅ exists
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: AdminManagement.tsx

#### Communities Management (6 endpoints)
**Status:** Partial SQLC queries exist, need admin variants

9. **GET /api/admin/communities**
   - Query: `ListAllCommunitiesAdmin` ✅ added (needs sqlc generate)
   - Handler: ❌ needs implementation
   - Service: ❌ needs implementation
   - Frontend: CommunitiesGroups.tsx (currently uses mockCommunitiesApi)

10. **POST /api/admin/communities**
    - Query: `CreateCommunity` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

11. **PUT /api/admin/communities/:id**
    - Query: `UpdateCommunity` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

12. **DELETE /api/admin/communities/:id**
    - Query: `DeleteCommunity` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

13. **PUT /api/admin/communities/:id/assign-moderator**
    - Query: `AddCommunityModerator` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

14. **PUT /api/admin/communities/:id/status**
    - Query: `UpdateCommunityStatus` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

#### User Management Completion (4 endpoints)

15. **POST /api/admin/users**
    - Query: `CreateUser` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: UserManagement.tsx (Add User modal)

16. **PUT /api/admin/users/:id**
    - Query: `UpdateUser` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: UserManagement.tsx (Edit User modal)

17. **PUT /api/admin/users/:id/reset-password**
    - Query: `ResetUserPassword` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: UserManagement.tsx (Reset Password action)

18. **GET /api/admin/users/:id/activity**
    - Query: Needs new query ❌
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: UserManagement.tsx (View Activity action)

### 🔨 To Be Implemented (Priority 2 - High)

#### Announcements Management (6 endpoints)

19. **GET /api/admin/announcements**
    - Query: `ListAllAnnouncementsAdmin` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx (Announcements tab)

20. **POST /api/admin/announcements**
    - Query: `CreateAnnouncement` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

21. **PUT /api/admin/announcements/:id**
    - Query: `UpdateAnnouncement` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

22. **DELETE /api/admin/announcements/:id**
    - Query: `DeleteAnnouncement` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

23. **PUT /api/admin/announcements/:id/publish**
    - Query: `UpdateAnnouncementStatus` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

24. **PUT /api/admin/announcements/:id/schedule**
    - Query: `UpdateAnnouncement` ✅ exists (can set scheduled_for)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

#### Events Management (6 endpoints)

25. **GET /api/admin/events**
    - Query: `ListAllEventsAdmin` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx (Events tab)

26. **POST /api/admin/events**
    - Query: `CreateEvent` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

27. **PUT /api/admin/events/:id**
    - Query: `UpdateEvent` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

28. **DELETE /api/admin/events/:id**
    - Query: `DeleteEvent` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

29. **PUT /api/admin/events/:id/cancel**
    - Query: `UpdateEventStatus` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

30. **GET /api/admin/events/:id/registrations**
    - Query: `GetEventAttendees` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: ContentManagement.tsx

#### Analytics (4 endpoints)

31. **GET /api/admin/analytics/user-growth**
    - Query: `GetUserGrowthData` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Analytics.tsx (currently uses mock userGrowthData)

32. **GET /api/admin/analytics/engagement**
    - Query: `GetContentGrowthData` ✅ exists (can be adapted)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Analytics.tsx (currently uses mock engagementData)

33. **GET /api/admin/analytics/group-activity**
    - Query: `GetActivityStats` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Analytics.tsx (currently uses mock groupActivityData)

34. **GET /api/admin/analytics/category-distribution**
    - Query: Needs aggregation query ❌
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Analytics.tsx (currently uses mock categoryData)

#### Notifications (4 endpoints)

35. **GET /api/admin/notifications**
    - Query: `GetAdminNotifications` ✅ added (needs sqlc generate)
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Notifications.tsx (currently uses mockNotifications)

36. **PUT /api/admin/notifications/:id/read**
    - Query: `MarkAsRead` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Notifications.tsx

37. **DELETE /api/admin/notifications/:id**
    - Query: `DeleteNotification` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Notifications.tsx

38. **PUT /api/admin/notifications/read-all**
    - Query: `MarkAllAsRead` ✅ exists
    - Handler: ❌ needs implementation
    - Service: ❌ needs implementation
    - Frontend: Notifications.tsx

### 🔨 To Be Implemented (Priority 3 - Nice to Have)

#### Export Functionality (3 endpoints)

39. **GET /api/admin/users/export**
    - Format: CSV
    - Handler: ❌ needs implementation
    - Frontend: UserManagement.tsx (Export CSV button)

40. **GET /api/admin/groups/export**
    - Format: CSV
    - Handler: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

41. **GET /api/admin/communities/export**
    - Format: CSV
    - Handler: ❌ needs implementation
    - Frontend: CommunitiesGroups.tsx

#### Bulk Operations (2 endpoints)

42. **POST /api/admin/users/bulk-action**
    - Actions: suspend, activate, delete
    - Handler: ❌ needs implementation
    - Frontend: UserManagement.tsx (bulk selection)

43. **POST /api/admin/reports/bulk-resolve**
    - Handler: ❌ needs implementation
    - Frontend: ContentManagement.tsx

## Implementation Steps

### Step 1: SQLC Generation
Run `make sqlc` or `sqlc generate` to generate Go code for new queries:
- ListAllCommunitiesAdmin
- ListAllAnnouncementsAdmin
- ListAllEventsAdmin
- DeleteCommunity
- DeleteAnnouncement
- DeleteEvent
- GetAdminNotifications
- AdminCreateUser
- AdminUpdateUser
- ResetUserPassword
- UpdateCommunityStatus
- GetEventWithRegistrations

### Step 2: Create Service Methods
Extend `backend/internal/service/admin/service.go` with new methods for all endpoints.

### Step 3: Create HTTP Handlers
Extend `backend/internal/api/handlers/admin.go` with handler functions.

### Step 4: Register Routes
Update `backend/internal/api/routes/admin_routes.go` with new route definitions.

### Step 5: Add Authorization
Apply role-based middleware to routes:
- Super admin only: /api/admin/admins/*, /api/admin/settings/*
- Admin only: All other admin routes

### Step 6: Frontend Integration
Update `frontend/src/services/adminApi.ts` with all new API calls.

### Step 7: Replace Mock Data
Update frontend components to use live APIs:
- CommunitiesGroups.tsx → replace mockCommunitiesApi
- ContentManagement.tsx → replace mock announcements/events
- Analytics.tsx → replace mock chart data
- Notifications.tsx → replace mockNotifications
- SystemSettings.tsx → replace mockSettings
- AdminManagement.tsx → connect to real API

### Step 8: Add UX Enhancements
- Loading skeletons for all data fetching
- Toast notifications for all mutations
- Error boundaries
- Optimistic updates where appropriate

### Step 9: Testing
- Backend unit tests for new endpoints
- Integration tests
- Frontend build verification
- End-to-end testing

## File Structure

### Backend
```
backend/
├── db/
│   ├── query/
│   │   ├── admins.sql (extended)
│   │   ├── communities.sql (extended)
│   │   ├── events.sql (extended)
│   │   └── notification.sql (extended)
│   └── sqlc/
│       └── *.sql.go (regenerated)
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── admin.go (extended)
│   │   └── routes/
│   │       └── admin_routes.go (extended)
│   └── service/
│       └── admin/
│           ├── service.go (extended)
│           └── types.go (extended)
└── test/
    └── api/
        └── admin_test.go (extended)
```

### Frontend
```
frontend/
├── src/
│   ├── services/
│   │   └── adminApi.ts (completely rewritten)
│   ├── pages/
│   │   └── admin/
│   │       ├── UserManagement.tsx (updated)
│   │       ├── ContentManagement.tsx (updated)
│   │       ├── CommunitiesGroups.tsx (updated)
│   │       ├── Analytics.tsx (updated)
│   │       ├── Notifications.tsx (updated)
│   │       ├── SystemSettings.tsx (updated)
│   │       └── AdminManagement.tsx (updated)
│   └── data/
│       └── mockAdminCommunitiesData.ts (deprecated/removed)
```

## Estimated Implementation Time

- **Priority 1 (Critical):** 18 endpoints × 30 min = 9 hours
- **Priority 2 (High):** 20 endpoints × 25 min = 8.3 hours
- **Priority 3 (Nice to Have):** 5 endpoints × 20 min = 1.7 hours
- **Frontend Integration:** 4 hours
- **Testing & Validation:** 3 hours
- **Total:** ~26 hours

## Success Criteria

- [ ] All 43 new endpoints implemented and tested
- [ ] SQLC queries regenerated successfully
- [ ] Backend builds without errors
- [ ] Frontend builds without errors
- [ ] All mock data removed from frontend
- [ ] All admin tabs display live data
- [ ] CRUD operations work end-to-end
- [ ] Loading states and error handling in place
- [ ] Toast notifications for all user actions
- [ ] Backend tests pass
- [ ] No TypeScript errors
- [ ] No Go compilation errors

## Dependencies

1. **sqlc** - Must be installable or Docker available
2. **PostgreSQL** - Database with all migrations applied
3. **Go 1.21+** - Backend compilation
4. **Node.js 18+** - Frontend build

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| SQLC generation fails | High | Use Docker sqlc or manual implementation |
| Type mismatches | Medium | Careful type alignment, use TypeScript generators |
| Breaking changes | High | Maintain backward compatibility, version APIs |
| Performance issues | Medium | Add pagination, implement caching |

---

**Last Updated:** 2025-11-06
**Status:** Ready for Implementation
