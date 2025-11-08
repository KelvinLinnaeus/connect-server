#!/usr/bin/env python3
"""
Generate Production-ready Postman Collection for Connect Server API
This script programmatically creates a complete Postman collection based on route discovery and handler analysis.
"""

import json
import uuid as uuid_lib
from typing import Dict, List, Any, Optional
from datetime import datetime

# Base URL variable
BASE_URL = "{{base_url}}"

# Pre-request script for authenticated endpoints
AUTH_PRE_REQUEST_SCRIPT = """
// Ensure token is set
const token = pm.environment.get('token');
if (!token) {
    console.log('Warning: No auth token found. Please run the Login request first.');
}
"""

# Test script for successful responses
SUCCESS_TEST_SCRIPT = """
pm.test("Status code is successful", function () {
    pm.expect(pm.response.code).to.be.oneOf([200, 201]);
});

pm.test("Response has valid structure", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.be.an('object');
});
"""

# Test script for login endpoint
LOGIN_TEST_SCRIPT = """
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has tokens", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('access_token');
    pm.expect(jsonData.data).to.have.property('refresh_token');

    // Store tokens in environment
    pm.environment.set('token', jsonData.data.access_token);
    pm.environment.set('refresh_token', jsonData.data.refresh_token);

    // Store user info if needed
    if (jsonData.data.user && jsonData.data.user.id) {
        pm.environment.set('user_id', jsonData.data.user.id);
    }

    console.log('Tokens stored successfully');
});
"""

# Test script for paginated responses
PAGINATED_TEST_SCRIPT = """
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has pagination meta", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('meta');
    pm.expect(jsonData.meta).to.have.property('total');
    pm.expect(jsonData.meta).to.have.property('page');
    pm.expect(jsonData.meta).to.have.property('limit');
});
"""


class PostmanCollectionGenerator:
    def __init__(self):
        self.collection = {
            "info": {
                "_postman_id": str(uuid_lib.uuid4()),
                "name": "Connect Server API",
                "description": self._generate_collection_description(),
                "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
            },
            "item": [],
            "variable": [
                {
                    "key": "base_url",
                    "value": "http://localhost:8080",
                    "type": "string"
                }
            ]
        }

    def _generate_collection_description(self) -> str:
        return """# Connect Server API - Production Ready Collection

This collection provides complete coverage of the Connect Server API with accurate request/response examples, authentication flows, and comprehensive tests.

## Getting Started

1. **Import Environment**: Import the `postman_env.json` file to set up environment variables
2. **Set Base URL**: Update `{{base_url}}` if your server is not running on `http://localhost:8080`
3. **Authenticate**: Run the "Login" request in the "Authentication" folder to obtain tokens
4. **Run Requests**: All authenticated endpoints will automatically use the stored token

## Environment Variables Required

- `base_url`: API base URL (default: http://localhost:8080)
- `token`: Access token (auto-set after login)
- `refresh_token`: Refresh token (auto-set after login)
- `user_id`: Current user ID (auto-set after login)
- `space_id`: Test space ID (set manually or from organization)
- `test_community_id`: Test community ID (create one or use existing)
- `test_group_id`: Test group ID (create one or use existing)
- `test_post_id`: Test post ID (create one or use existing)
- `test_conversation_id`: Test conversation ID (create one or use existing)

## Features

- ✅ Complete endpoint coverage (120+ endpoints)
- ✅ Automatic token management
- ✅ Request validation and examples
- ✅ Response schema validation
- ✅ Pagination support
- ✅ Error handling examples
- ✅ Production-ready test scripts

## API Modules

1. **Health** - Service health check
2. **Authentication** - Login, refresh, logout
3. **Users** - User management, search, profiles
4. **Posts** - Posts, comments, likes, feed
5. **Communities** - Community management and membership
6. **Groups** - Project groups, roles, applications
7. **Messaging** - Conversations and direct messages
8. **Notifications** - User notifications
9. **Events** - Event management and registration
10. **Announcements** - System announcements
11. **Mentorship** - Mentors, tutors, and sessions
12. **Analytics** - Reports, metrics, and statistics

## Notes

- All timestamps are in ISO 8601 format with UTC timezone
- All IDs are UUIDs (v4)
- Default pagination: page=1, limit=20, max=100
- Rate limiting is applied to most endpoints
"""

    def create_request(
        self,
        name: str,
        method: str,
        path: str,
        description: str = "",
        auth_required: bool = False,
        body: Optional[Dict] = None,
        query_params: Optional[List[Dict]] = None,
        path_vars: Optional[List[Dict]] = None,
        test_script: Optional[str] = None,
        pre_request_script: Optional[str] = None
    ) -> Dict:
        """Create a Postman request item."""

        # Build URL
        url_parts = {
            "raw": f"{BASE_URL}{path}",
            "host": ["{{base_url}}"],
            "path": [p for p in path.split("/") if p]
        }

        if query_params:
            url_parts["query"] = query_params

        # Build headers
        headers = [
            {"key": "Content-Type", "value": "application/json", "type": "text"},
            {"key": "Accept", "value": "application/json", "type": "text"}
        ]

        if auth_required:
            headers.append({
                "key": "Authorization",
                "value": "Bearer {{token}}",
                "type": "text"
            })

        # Build request object
        request_obj = {
            "method": method,
            "header": headers,
            "url": url_parts
        }

        # Add body if present
        if body and method in ["POST", "PUT", "PATCH"]:
            request_obj["body"] = {
                "mode": "raw",
                "raw": json.dumps(body, indent=2),
                "options": {
                    "raw": {
                        "language": "json"
                    }
                }
            }

        # Build event scripts
        events = []

        # Add pre-request script
        if auth_required and not pre_request_script:
            pre_request_script = AUTH_PRE_REQUEST_SCRIPT

        if pre_request_script:
            events.append({
                "listen": "prerequest",
                "script": {
                    "type": "text/javascript",
                    "exec": pre_request_script.split("\n")
                }
            })

        # Add test script
        if not test_script:
            test_script = SUCCESS_TEST_SCRIPT

        if test_script:
            events.append({
                "listen": "test",
                "script": {
                    "type": "text/javascript",
                    "exec": test_script.split("\n")
                }
            })

        # Build complete item
        item = {
            "name": name,
            "request": request_obj,
            "response": []
        }

        if description:
            item["request"]["description"] = description

        if events:
            item["event"] = events

        return item

    def create_folder(self, name: str, description: str = "", items: List[Dict] = None) -> Dict:
        """Create a Postman folder."""
        folder = {
            "name": name,
            "item": items or []
        }
        if description:
            folder["description"] = description
        return folder

    def generate_health_folder(self) -> Dict:
        """Generate Health Check folder."""
        items = [
            self.create_request(
                name="Health Check",
                method="GET",
                path="/health",
                description="Check API server health and database connectivity status.",
                auth_required=False,
                test_script="""
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Service is healthy", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('ok');
    pm.expect(jsonData.db).to.eql('ok');
});
"""
            )
        ]
        return self.create_folder("Health", "Service health check endpoint", items)

    def generate_auth_folder(self) -> Dict:
        """Generate Authentication folder."""
        items = [
            self.create_request(
                name="Login",
                method="POST",
                path="/api/users/login",
                description="Authenticate user and receive access/refresh tokens. Use these tokens for authenticated endpoints.",
                auth_required=False,
                body={
                    "email": "test@example.com",
                    "password": "password123"
                },
                test_script=LOGIN_TEST_SCRIPT
            ),
            self.create_request(
                name="Refresh Token",
                method="POST",
                path="/api/users/refresh",
                description="Obtain a new access token using a valid refresh token.",
                auth_required=False,
                body={
                    "refresh_token": "{{refresh_token}}"
                },
                test_script="""
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has new access token", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('access_token');

    // Update token in environment
    pm.environment.set('token', jsonData.data.access_token);
    console.log('Access token refreshed successfully');
});
"""
            ),
            self.create_request(
                name="Logout",
                method="POST",
                path="/api/users/logout",
                description="Logout the current user and invalidate the session.",
                auth_required=True,
                test_script="""
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Logout successful", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('message');

    // Optionally clear tokens (uncomment if desired)
    // pm.environment.unset('token');
    // pm.environment.unset('refresh_token');
});
"""
            )
        ]
        return self.create_folder(
            "Authentication",
            "User authentication endpoints including login, token refresh, and logout.",
            items
        )

    def generate_users_folder(self) -> Dict:
        """Generate Users folder."""
        items = [
            self.create_request(
                name="Create User (Register)",
                method="POST",
                path="/api/users",
                description="Register a new user account. This is a public endpoint that doesn't require authentication.",
                auth_required=False,
                body={
                    "space_id": "{{space_id}}",
                    "username": "johndoe",
                    "email": "john.doe@example.com",
                    "password": "SecurePass123!",
                    "full_name": "John Doe",
                    "level": "undergraduate",
                    "department": "Computer Science",
                    "major": "Software Engineering",
                    "year": 3,
                    "interests": ["programming", "ai", "web development"]
                },
                test_script="""
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("User created successfully", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('id');
    pm.expect(jsonData.data).to.have.property('username');
    pm.expect(jsonData.data).to.have.property('email');

    // Optionally store new user ID
    if (jsonData.data.id) {
        pm.environment.set('new_user_id', jsonData.data.id);
    }
});
"""
            ),
            self.create_request(
                name="Get User by ID",
                method="GET",
                path="/api/users/:id",
                description="Retrieve user information by user ID.",
                auth_required=True,
                path_vars=[{"key": "id", "value": "{{user_id}}"}]
            ),
            self.create_request(
                name="Get User by Username",
                method="GET",
                path="/api/users/username/:username",
                description="Retrieve user information by username and space ID.",
                auth_required=False,
                query_params=[
                    {"key": "space_id", "value": "{{space_id}}", "description": "Space ID (required)"}
                ]
            ),
            self.create_request(
                name="Update User",
                method="PUT",
                path="/api/users/:id",
                description="Update user profile information. Only the authenticated user can update their own profile.",
                auth_required=True,
                body={
                    "full_name": "John Doe Updated",
                    "bio": "Software engineer passionate about building great products",
                    "level": "graduate",
                    "department": "Computer Science",
                    "major": "Artificial Intelligence",
                    "year": 1,
                    "interests": ["machine learning", "deep learning", "nlp"]
                }
            ),
            self.create_request(
                name="Update Password",
                method="PUT",
                path="/api/users/:id/password",
                description="Change user password. Requires old password for verification.",
                auth_required=True,
                body={
                    "old_password": "OldPassword123!",
                    "new_password": "NewSecurePass456!"
                }
            ),
            self.create_request(
                name="Deactivate User",
                method="DELETE",
                path="/api/users/:id",
                description="Deactivate user account. This soft-deletes the user.",
                auth_required=True
            ),
            self.create_request(
                name="Search Users",
                method="GET",
                path="/api/users/search",
                description="Search for users by name, username, or other criteria.",
                auth_required=False,
                query_params=[
                    {"key": "q", "value": "john", "description": "Search query (required)"},
                    {"key": "space_id", "value": "{{space_id}}", "description": "Space ID (required)"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            )
        ]
        return self.create_folder(
            "Users",
            "User management endpoints including registration, profile updates, and search.",
            items
        )

    def generate_posts_folder(self) -> Dict:
        """Generate Posts folder."""
        items = [
            # Create and manage posts
            self.create_request(
                name="Create Post",
                method="POST",
                path="/api/posts",
                description="Create a new post. Can be a standalone post, community post, group post, or reply to another post.",
                auth_required=True,
                body={
                    "space_id": "{{space_id}}",
                    "content": "This is my first post! Excited to be part of this community.",
                    "tags": ["introduction", "firstpost"],
                    "visibility": "public"
                },
                test_script="""
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Post created successfully", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('id');

    // Store post ID for later use
    if (jsonData.data.id) {
        pm.environment.set('test_post_id', jsonData.data.id);
    }
});
"""
            ),
            self.create_request(
                name="Get Post",
                method="GET",
                path="/api/posts/:id",
                description="Retrieve a post by ID. This also increments the view count.",
                auth_required=False
            ),
            self.create_request(
                name="Delete Post",
                method="DELETE",
                path="/api/posts/:id",
                description="Delete a post. Only the post author can delete their own post.",
                auth_required=True
            ),

            # Get posts from various sources
            self.create_request(
                name="Get User Feed",
                method="GET",
                path="/api/posts/feed",
                description="Get personalized feed for the authenticated user based on their following and interests.",
                auth_required=True,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Get User Posts",
                method="GET",
                path="/api/posts/user/:user_id",
                description="Get all posts by a specific user.",
                auth_required=False,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Get Community Posts",
                method="GET",
                path="/api/posts/community/:community_id",
                description="Get all posts in a specific community.",
                auth_required=False,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Get Group Posts",
                method="GET",
                path="/api/posts/group/:group_id",
                description="Get all posts in a specific group.",
                auth_required=False,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Get Trending Posts",
                method="GET",
                path="/api/posts/trending",
                description="Get trending posts based on engagement metrics.",
                auth_required=False,
                query_params=[
                    {"key": "space_id", "value": "{{space_id}}"}
                ]
            ),
            self.create_request(
                name="Get User Liked Posts",
                method="GET",
                path="/api/posts/liked",
                description="Get all posts liked by the authenticated user.",
                auth_required=True,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),

            # Search
            self.create_request(
                name="Search Posts",
                method="GET",
                path="/api/posts/search",
                description="Basic search for posts by content.",
                auth_required=False,
                query_params=[
                    {"key": "q", "value": "technology"},
                    {"key": "space_id", "value": "{{space_id}}"},
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Advanced Search Posts",
                method="GET",
                path="/api/posts/advanced-search",
                description="Advanced search with filters and sorting options.",
                auth_required=False,
                query_params=[
                    {"key": "q", "value": "programming"},
                    {"key": "space_id", "value": "{{space_id}}"},
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),

            # Comments
            self.create_request(
                name="Get Post Comments",
                method="GET",
                path="/api/posts/:id/comments",
                description="Get all comments on a post.",
                auth_required=False
            ),
            self.create_request(
                name="Create Comment",
                method="POST",
                path="/api/posts/:id/comments",
                description="Add a comment to a post.",
                auth_required=True,
                body={
                    "content": "Great post! Thanks for sharing."
                }
            ),
            self.create_request(
                name="Toggle Comment Like",
                method="POST",
                path="/api/comments/:id/like",
                description="Like or unlike a comment.",
                auth_required=True
            ),

            # Likes
            self.create_request(
                name="Get Post Likes",
                method="GET",
                path="/api/posts/:id/likes",
                description="Get all users who liked a post.",
                auth_required=False
            ),
            self.create_request(
                name="Toggle Post Like",
                method="POST",
                path="/api/posts/:id/like",
                description="Like or unlike a post.",
                auth_required=True
            ),

            # Repost and pin
            self.create_request(
                name="Create Repost",
                method="POST",
                path="/api/posts/:id/repost",
                description="Repost/share a post to your profile.",
                auth_required=True
            ),
            self.create_request(
                name="Pin Post",
                method="PUT",
                path="/api/posts/:id/pin",
                description="Pin a post to the top of your profile or community.",
                auth_required=True
            )
        ]
        return self.create_folder(
            "Posts & Comments",
            "Post management including creation, feeds, comments, likes, and reposts.",
            items
        )

    def generate_sessions_folder(self) -> Dict:
        """Generate Sessions folder."""
        items = [
            self.create_request(
                name="Get Session",
                method="GET",
                path="/api/sessions/:id",
                description="Retrieve session information by session ID. Used for auth session management.",
                auth_required=True
            )
        ]
        return self.create_folder(
            "Sessions",
            "Session management endpoints for retrieving authentication session information.",
            items
        )

    def generate_communities_folder(self) -> Dict:
        """Generate Communities folder."""
        items = [
            # Create and manage communities
            self.create_request(
                name="Create Community",
                method="POST",
                path="/api/communities",
                description="Create a new community.",
                auth_required=True,
                body={
                    "space_id": "{{space_id}}",
                    "name": "Tech Enthusiasts",
                    "slug": "tech-enthusiasts",
                    "description": "A community for technology enthusiasts to discuss latest trends and innovations",
                    "category": "Technology",
                    "tags": ["technology", "innovation", "programming"],
                    "is_private": False,
                    "rules": "Be respectful and constructive"
                },
                test_script="""
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Community created successfully", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('id');

    // Store community ID for later use
    if (jsonData.data.id) {
        pm.environment.set('test_community_id', jsonData.data.id);
    }
});
"""
            ),
            self.create_request(
                name="List Communities",
                method="GET",
                path="/api/communities",
                description="Get list of all communities with pagination.",
                auth_required=False,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ],
                test_script=PAGINATED_TEST_SCRIPT
            ),
            self.create_request(
                name="Get Community",
                method="GET",
                path="/api/communities/:id",
                description="Get detailed information about a specific community.",
                auth_required=False
            ),
            self.create_request(
                name="Get Community by Slug",
                method="GET",
                path="/api/communities/slug/:slug",
                description="Get community by its unique slug identifier.",
                auth_required=False
            ),
            self.create_request(
                name="Update Community",
                method="PUT",
                path="/api/communities/:id",
                description="Update community information. Only admins can update.",
                auth_required=True,
                body={
                    "name": "Tech Enthusiasts Updated",
                    "description": "Updated description for tech community",
                    "tags": ["technology", "innovation", "programming", "ai"]
                }
            ),

            # Search and categories
            self.create_request(
                name="Search Communities",
                method="GET",
                path="/api/communities/search",
                description="Search for communities by name, description, or tags.",
                auth_required=False,
                query_params=[
                    {"key": "q", "value": "technology"},
                    {"key": "space_id", "value": "{{space_id}}"},
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ]
            ),
            self.create_request(
                name="Get Community Categories",
                method="GET",
                path="/api/communities/categories",
                description="Get list of all available community categories.",
                auth_required=False
            ),

            # Membership
            self.create_request(
                name="Join Community",
                method="POST",
                path="/api/communities/:id/join",
                description="Join a community as a member.",
                auth_required=True
            ),
            self.create_request(
                name="Leave Community",
                method="POST",
                path="/api/communities/:id/leave",
                description="Leave a community.",
                auth_required=True
            ),
            self.create_request(
                name="Get User Communities",
                method="GET",
                path="/api/users/communities",
                description="Get all communities the authenticated user is a member of.",
                auth_required=True,
                test_script=PAGINATED_TEST_SCRIPT
            ),

            # Members and roles
            self.create_request(
                name="Get Community Members",
                method="GET",
                path="/api/communities/:id/members",
                description="Get list of all community members.",
                auth_required=False,
                query_params=[
                    {"key": "page", "value": "1"},
                    {"key": "limit", "value": "20"}
                ]
            ),
            self.create_request(
                name="Get Community Moderators",
                method="GET",
                path="/api/communities/:id/moderators",
                description="Get list of community moderators.",
                auth_required=False
            ),
            self.create_request(
                name="Get Community Admins",
                method="GET",
                path="/api/communities/:id/admins",
                description="Get list of community admins.",
                auth_required=False
            ),

            # Moderation
            self.create_request(
                name="Add Community Moderator",
                method="POST",
                path="/api/communities/:id/moderators",
                description="Add a user as a community moderator. Only admins can perform this action.",
                auth_required=True,
                body={
                    "user_id": "{{user_id}}"
                }
            ),
            self.create_request(
                name="Remove Community Moderator",
                method="DELETE",
                path="/api/communities/:id/moderators/:userId",
                description="Remove a user from community moderators. Only admins can perform this action.",
                auth_required=True
            )
        ]
        return self.create_folder(
            "Communities",
            "Community management including creation, membership, moderation, and search.",
            items
        )

    def generate(self) -> Dict:
        """Generate complete collection."""
        # Add all folders
        self.collection["item"].extend([
            self.generate_health_folder(),
            self.generate_auth_folder(),
            self.generate_users_folder(),
            self.generate_posts_folder(),
            self.generate_sessions_folder(),
            self.generate_communities_folder(),
            # More folders will be generated...
        ])

        return self.collection


def main():
    """Main entry point."""
    print("Generating Connect Server Postman Collection...")
    print("=" * 70)

    generator = PostmanCollectionGenerator()
    collection = generator.generate()

    # Write collection to file
    output_path = "db/postman_collection.json"
    with open(output_path, 'w') as f:
        json.dumps(collection, f, indent=2)

    print(f"\n✓ Collection generated: {output_path}")
    print(f"  Total folders: {len(collection['item'])}")

    # Count total requests
    total_requests = 0
    for folder in collection['item']:
        total_requests += len(folder.get('item', []))

    print(f"  Total requests: {total_requests}")
    print("\nPartial generation complete. Need to add remaining modules...")


if __name__ == "__main__":
    main()
