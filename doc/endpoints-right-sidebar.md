# RightSidebar Backend Endpoints Documentation

## Overview
This document describes the backend API endpoints required by the RightSidebar component in the frontend application.

---

## 1. Suggested Users

### Endpoint
```
GET /api/users/suggested
```

### Description
Returns a paginated list of suggested users for the authenticated user to follow.

### Authentication
**Required**: Yes (Bearer token)

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `space_id` | UUID | Yes | Space ID to scope the suggestions |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10, max: 50) |

### Request Example
```http
GET /api/users/suggested?space_id=123e4567-e89b-12d3-a456-426614174000&page=1&limit=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Response Example
```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "john_doe",
      "full_name": "John Doe",
      "avatar": "https://example.com/avatars/john.jpg",
      "bio": "Computer Science student passionate about AI",
      "verified": true,
      "followers_count": 245,
      "following_count": 180,
      "is_following": false
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "username": "jane_smith",
      "full_name": "Jane Smith",
      "avatar": null,
      "bio": "Math major | Chess enthusiast",
      "verified": false,
      "followers_count": 89,
      "following_count": 120,
      "is_following": false
    }
  ]
}
```

### Algorithm (v1.0)
Current implementation uses a simple algorithm:
1. Get recent users from the same space (last 30 days of activity)
2. Exclude users the current user already follows
3. Exclude the current user
4. Order by: followers_count DESC, created_at DESC
5. Return paginated results

**Future improvements:**
- Mutual connections analysis
- Shared interests/departments
- Activity similarity scores
- ML-based recommendations

### Error Responses
```json
// Missing space_id
{
  "error": {
    "code": "missing_space_id",
    "message": "space_id query parameter is required"
  }
}

// Invalid space_id format
{
  "error": {
    "code": "invalid_space_id",
    "message": "Invalid space ID format"
  }
}

// Unauthorized
{
  "error": {
    "code": "unauthorized",
    "message": "Authentication required"
  }
}
```

---

## 2. Trending Topics

### Endpoint
```
GET /api/topics/trending
```

### Description
Returns trending topics/hashtags aggregated from recent posts, ordered by relevance and activity.

### Authentication
**Required**: No (public endpoint, but can be scoped)

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `space_id` | UUID | Yes | Space ID to scope the trends |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 10, max: 50) |
| `timeframe` | string | No | Time period: `24h`, `7d`, `30d` (default: `24h`) |

### Request Example
```http
GET /api/topics/trending?space_id=123e4567-e89b-12d3-a456-426614174000&timeframe=24h&limit=10
```

### Response Example
```json
{
  "data": [
    {
      "id": "exam-week",
      "name": "#ExamWeek",
      "category": "Campus Life",
      "posts_count": 1247,
      "trend_score": 95.5,
      "created_at": "2025-11-06T10:00:00Z"
    },
    {
      "id": "ml-workshop",
      "name": "#MLWorkshop",
      "category": "Events",
      "posts_count": 458,
      "trend_score": 82.3,
      "created_at": "2025-11-06T09:30:00Z"
    },
    {
      "id": "basketball",
      "name": "#Basketball",
      "category": "Sports",
      "posts_count": 324,
      "trend_score": 71.8,
      "created_at": "2025-11-05T18:00:00Z"
    }
  ]
}
```

### Algorithm (v1.0)
Trending score calculation:
```
trend_score = (posts_count * recency_factor * engagement_factor) / time_decay

Where:
- recency_factor: Higher for recent topics (exponential decay)
- engagement_factor: Based on likes, comments, reposts
- time_decay: Reduces score for older topics
```

Current implementation:
1. Extract hashtags and tags from posts in the specified timeframe
2. Count occurrences per topic
3. Calculate engagement metrics (likes + comments + reposts)
4. Apply time decay (more recent = higher score)
5. Categorize topics based on common keywords
6. Sort by trend_score DESC
7. Return paginated results

### Error Responses
```json
// Missing space_id
{
  "error": {
    "code": "missing_space_id",
    "message": "space_id query parameter is required"
  }
}

// Invalid timeframe
{
  "error": {
    "code": "invalid_timeframe",
    "message": "timeframe must be one of: 24h, 7d, 30d"
  }
}
```

---

## 3. Campus Highlights (Announcements)

### Endpoint
```
GET /api/announcements
```

### Description
Returns campus announcements and highlights. This endpoint is already fully implemented.

### Authentication
**Required**: No (public endpoint)

### Query Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `space_id` | UUID | Yes | Space ID to scope announcements |
| `status` | string | No | Filter by status: `published`, `draft`, `archived` |
| `page` | integer | No | Page number (default: 1) |
| `limit` | integer | No | Items per page (default: 20, max: 100) |

### Request Example
```http
GET /api/announcements?space_id=123e4567-e89b-12d3-a456-426614174000&status=published&limit=5
```

### Response Example
```json
{
  "data": [
    {
      "id": "ann-001",
      "space_id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "Campus Library Extended Hours for Exams",
      "content": "The library will be open 24/7 starting next week to support students during the exam period.",
      "type": "info",
      "priority": "high",
      "status": "published",
      "published_at": "2025-11-06T08:00:00Z",
      "created_at": "2025-11-05T15:00:00Z"
    }
  ]
}
```

### Notes
- This endpoint is fully functional and requires no changes
- Frontend already uses this endpoint successfully
- Supports filtering by status to show only published announcements

---

## Implementation Status

| Endpoint | Status | Tests | Documentation |
|----------|--------|-------|---------------|
| `GET /api/announcements` | ✅ Complete | ✅ Yes | ✅ Yes |
| `GET /api/users/suggested` | ⚠️  To Implement | ❌ No | ✅ Yes |
| `GET /api/topics/trending` | ⚠️  To Implement | ❌ No | ✅ Yes |

---

## Database Schema Notes

### For Suggested Users
Uses existing `users` table with:
- `followers_count` (computed or cached)
- `following_count` (computed or cached)
- `created_at` for sorting
- Requires join with `follows` table to check `is_following` status

### For Trending Topics
Requires aggregation from:
- `posts` table (content, tags, created_at)
- `post_interactions` table (likes, comments, reposts)
- Real-time calculation (no dedicated topics table needed for v1.0)

---

## Frontend Integration

### Environment Variable
```bash
VITE_DEFAULT_SPACE_ID=<your-space-uuid>
```

This is automatically injected by the frontend API client using `getValidatedSpaceId()`.

### API Client Pattern
```typescript
import { getValidatedSpaceId } from '@/lib/apiClient';

export const getSuggestedUsers = async (params?: PaginationParams) => {
  const spaceId = getValidatedSpaceId();
  const response = await apiClient.get('/users/suggested', {
    params: { space_id: spaceId, ...params },
  });
  return response.data.data;
};
```

---

## Testing Guidelines

### Backend Tests
Located in: `backend/internal/api/handlers/*_test.go`

Test cases required:
1. **Happy path**: Valid request returns suggested users
2. **Pagination**: Page and limit parameters work correctly
3. **Authentication**: Requires valid token
4. **Missing space_id**: Returns 400 error
5. **Invalid space_id**: Returns 400 error
6. **Empty results**: Returns empty array when no suggestions

### Frontend Tests
Located in: `frontend/src/components/layout/__tests__/RightSidebar.test.tsx`

Test cases required:
1. **Loading state**: Shows skeletons while fetching
2. **Data display**: Shows users/topics when loaded
3. **Error handling**: Shows error message on API failure
4. **Follow action**: Calls follow API and updates UI
5. **Empty state**: Shows "No suggestions" message

---

## Performance Considerations

### Caching Strategy
- **Suggested Users**: Cache for 10 minutes (low volatility)
- **Trending Topics**: Cache for 5 minutes (moderate volatility)
- **Announcements**: Cache for 5 minutes (moderate volatility)

### Optimization Ideas
- Add Redis caching layer for trending topics calculation
- Pre-compute suggested users for active users
- Use database indexes on frequently queried fields
- Implement rate limiting per user/IP

---

## Future Enhancements

### Suggested Users (v2.0)
- [ ] ML-based recommendations using collaborative filtering
- [ ] Mutual connections ("followed by users you follow")
- [ ] Interest-based matching (shared communities, groups)
- [ ] Department/major similarity
- [ ] Activity similarity (posting patterns, engagement)

### Trending Topics (v2.0)
- [ ] Dedicated `topics` table with pre-computed trends
- [ ] Real-time trending using streaming aggregation
- [ ] Geographic/demographic trending filters
- [ ] Topic categorization using NLP
- [ ] Trending images/media content

### General
- [ ] A/B testing framework for recommendation algorithms
- [ ] Analytics dashboard for trend performance
- [ ] User feedback loop (dismiss/not interested)
- [ ] Personalized trending based on user interests

---

**Last Updated**: 2025-11-06
**Version**: 1.0
**Maintainer**: Backend Team
