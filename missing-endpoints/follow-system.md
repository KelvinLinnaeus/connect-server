# Missing Follow System Endpoints

## Overview
The frontend RightSidebar component expects follow/unfollow functionality, but the backend does not currently implement these endpoints. The frontend is calling these endpoints but they will return 404 errors.

## Required Endpoints

### 1. Follow User
**Endpoint:** `POST /api/users/:id/follow`
**Authentication:** Required
**Description:** Follow a user

**Request:**
```http
POST /api/users/:id/follow
Authorization: Bearer <token>
```

**URL Parameters:**
- `id` (required): UUID of the user to follow

**Response (201 Created):**
```json
{
  "data": {
    "follower_id": "uuid",
    "following_id": "uuid",
    "created_at": "2025-11-06T01:00:00Z"
  }
}
```

**Errors:**
- `400 Bad Request`: Invalid user ID format
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User not found
- `409 Conflict`: Already following this user

---

### 2. Unfollow User
**Endpoint:** `DELETE /api/users/:id/follow`
**Authentication:** Required
**Description:** Unfollow a user

**Request:**
```http
DELETE /api/users/:id/follow
Authorization: Bearer <token>
```

**URL Parameters:**
- `id` (required): UUID of the user to unfollow

**Response (200 OK):**
```json
{
  "data": {
    "message": "Successfully unfollowed user"
  }
}
```

**Errors:**
- `400 Bad Request`: Invalid user ID format
- `401 Unauthorized`: Not authenticated
- `404 Not Found`: User not found or not following

---

### 3. Get User Followers (Optional - for future enhancement)
**Endpoint:** `GET /api/users/:id/followers`
**Authentication:** Optional
**Description:** Get list of users following the specified user

**Request:**
```http
GET /api/users/:id/followers?page=1&limit=20
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "username": "john_doe",
      "full_name": "John Doe",
      "avatar": "https://...",
      "verified": true,
      "followers_count": 150,
      "is_following": false
    }
  ]
}
```

---

### 4. Get User Following (Optional - for future enhancement)
**Endpoint:** `GET /api/users/:id/following`
**Authentication:** Optional
**Description:** Get list of users that the specified user is following

**Request:**
```http
GET /api/users/:id/following?page=1&limit=20
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid",
      "username": "jane_smith",
      "full_name": "Jane Smith",
      "avatar": "https://...",
      "verified": false,
      "followers_count": 300,
      "is_following": true
    }
  ]
}
```

---

## Database Schema Requirements

### follows Table
The backend likely already has a `follows` table (as referenced in the GetSuggestedUsers query). If not, it should be created:

```sql
CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);
```

---

## Implementation Checklist

### Backend Implementation
- [ ] **Handler**: Create follow/unfollow methods in `backend/internal/api/handlers/user.go`
- [ ] **Service**: Create follow/unfollow service methods in `backend/internal/service/users/service.go`
- [ ] **SQL Queries**: Add queries to `backend/db/query/users.sql`:
  - `CreateFollow` (INSERT INTO follows)
  - `DeleteFollow` (DELETE FROM follows)
  - `GetUserFollowers` (SELECT followers)
  - `GetUserFollowing` (SELECT following)
  - `CheckIsFollowing` (EXISTS query)
- [ ] **Routes**: Register routes in `backend/internal/api/routes/user_routes.go`:
  ```go
  authUsers.POST("/:id/follow", userHandler.FollowUser)
  authUsers.DELETE("/:id/follow", userHandler.UnfollowUser)
  authUsers.GET("/:id/followers", userHandler.GetUserFollowers)
  authUsers.GET("/:id/following", userHandler.GetUserFollowing)
  ```
- [ ] **Validation**: Ensure users cannot follow themselves
- [ ] **Counters**: Update `users.followers_count` and `users.following_count` on follow/unfollow

### Frontend Integration
✅ **Already Implemented**: Frontend is ready and will work once backend endpoints are available
- `followUser()` in `frontend/src/api/users.api.ts`
- `unfollowUser()` in `frontend/src/api/users.api.ts`
- React Query mutations in RightSidebar component with loading states and toast notifications

---

## Current Status
- **Frontend**: ✅ Fully implemented with proper error handling
- **Backend**: ❌ Not implemented
- **Priority**: High (required for RightSidebar "Who to follow" feature)
- **Estimated Effort**: 2-3 hours

---

## Testing Checklist
Once implemented, verify:
- [ ] Can follow a user successfully
- [ ] Cannot follow the same user twice (409 Conflict)
- [ ] Cannot follow yourself (400 Bad Request)
- [ ] Can unfollow a user successfully
- [ ] Cannot unfollow a user you're not following (404 Not Found)
- [ ] Follower/following counts update correctly
- [ ] Frontend toasts display appropriate messages
- [ ] Loading spinners appear during mutation

---

## Notes
- The `GetSuggestedUsers` SQL query (line 80-102 in `backend/db/query/users.sql`) already uses the `follows` table for mutual followers calculation, confirming the table exists in the database schema
- The frontend already has proper error handling and will gracefully handle 404 errors with toast notifications
- Consider implementing rate limiting on follow/unfollow actions to prevent spam
