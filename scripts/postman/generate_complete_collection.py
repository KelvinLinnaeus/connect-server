#!/usr/bin/env python3
"""
Generate Complete Production-ready Postman Collection for Connect Server API
This script generates all 120+ endpoints with proper structure, auth, tests, and examples.
"""

import json
import uuid as uuid_lib
from typing import Dict, List, Any, Optional

BASE_URL = "{{base_url}}"

# Script templates
AUTH_PRE_REQUEST = """const token = pm.environment.get('token');
if (!token) {
    console.log('Warning: No auth token found. Please run the Login request first.');
}"""

SUCCESS_TEST = """pm.test("Status code is successful", function () {
    pm.expect(pm.response.code).to.be.oneOf([200, 201]);
});

pm.test("Response has valid structure", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.be.an('object');
});"""

LOGIN_TEST = """pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has tokens", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('access_token');
    pm.expect(jsonData.data).to.have.property('refresh_token');
    pm.environment.set('token', jsonData.data.access_token);
    pm.environment.set('refresh_token', jsonData.data.refresh_token);
    if (jsonData.data.user && jsonData.data.user.id) {
        pm.environment.set('user_id', jsonData.data.user.id);
    }
    console.log('Tokens stored successfully');
});"""

PAGINATED_TEST = """pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has pagination meta", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('meta');
    pm.expect(jsonData.meta).to.have.property('total');
    pm.expect(jsonData.meta).to.have.property('page');
    pm.expect(jsonData.meta).to.have.property('limit');
});"""

CREATE_TEST = """pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("Resource created successfully", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData.status).to.eql('success');
    pm.expect(jsonData.data).to.have.property('id');
});"""


class PostmanGenerator:
    def __init__(self):
        self.collection = {
            "info": {
                "_postman_id": str(uuid_lib.uuid4()),
                "name": "Connect Server API",
                "description": """# Connect Server API - Complete Production Collection

## Getting Started
1. Import `postman_env.json` for environment variables
2. Run "Login" in Authentication folder to get tokens
3. All authenticated endpoints auto-use stored token

## Environment Variables
- base_url, token, refresh_token, user_id, space_id
- test_community_id, test_group_id, test_post_id, test_event_id, test_conversation_id

## Coverage: 120+ endpoints across 12 modules
- Health, Auth, Users, Posts, Sessions, Communities
- Groups, Messaging, Notifications, Events
- Announcements, Mentorship, Analytics""",
                "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
            },
            "item": [],
            "variable": [{"key": "base_url", "value": "http://localhost:8080", "type": "string"}]
        }

    def req(self, name, method, path, desc="", auth=False, body=None, query=None, test=None, pre=None):
        """Create request."""
        url = {
            "raw": f"{BASE_URL}{path}",
            "host": ["{{base_url}}"],
            "path": [p for p in path.split("/") if p]
        }
        if query:
            url["query"] = query

        headers = [
            {"key": "Content-Type", "value": "application/json"},
            {"key": "Accept", "value": "application/json"}
        ]
        if auth:
            headers.append({"key": "Authorization", "value": "Bearer {{token}}"})

        request = {"method": method, "header": headers, "url": url}

        if body and method in ["POST", "PUT", "PATCH"]:
            request["body"] = {
                "mode": "raw",
                "raw": json.dumps(body, indent=2),
                "options": {"raw": {"language": "json"}}
            }

        events = []
        if auth and not pre:
            pre = AUTH_PRE_REQUEST
        if pre:
            events.append({"listen": "prerequest", "script": {"type": "text/javascript", "exec": pre.split("\n")}})

        if not test:
            test = SUCCESS_TEST
        if test:
            events.append({"listen": "test", "script": {"type": "text/javascript", "exec": test.split("\n")}})

        item = {"name": name, "request": request, "response": []}
        if desc:
            item["request"]["description"] = desc
        if events:
            item["event"] = events
        return item

    def folder(self, name, desc, items):
        """Create folder."""
        return {"name": name, "description": desc, "item": items}

    def gen_health(self):
        return self.folder("Health", "Service health check", [
            self.req("Health Check", "GET", "/health", "Check API server and database health", auth=False,
                     test="""pm.test("Status code is 200", () => pm.response.to.have.status(200));
pm.test("Service is healthy", () => {
    const data = pm.response.json();
    pm.expect(data.status).to.eql('ok');
    pm.expect(data.db).to.eql('ok');
});""")
        ])

    def gen_auth(self):
        return self.folder("Authentication", "Login, refresh, logout", [
            self.req("Login", "POST", "/api/users/login", "Authenticate and receive tokens",
                     body={"email": "test@example.com", "password": "password123"}, test=LOGIN_TEST),
            self.req("Refresh Token", "POST", "/api/users/refresh", "Get new access token",
                     body={"refresh_token": "{{refresh_token}}"}, test="""pm.test("Status 200", () => pm.response.to.have.status(200));
pm.test("Has new token", () => {
    const data = pm.response.json();
    pm.expect(data.data).to.have.property('access_token');
    pm.environment.set('token', data.data.access_token);
});"""),
            self.req("Logout", "POST", "/api/users/logout", "Logout and invalidate session", auth=True)
        ])

    def gen_users(self):
        return self.folder("Users", "User management and profiles", [
            self.req("Create User", "POST", "/api/users", "Register new user",
                     body={"space_id": "{{space_id}}", "username": "johndoe", "email": "john@example.com",
                           "password": "SecurePass123", "full_name": "John Doe", "level": "undergraduate",
                           "department": "CS", "year": 3, "interests": ["programming", "ai"]}, test=CREATE_TEST),
            self.req("Get User by ID", "GET", "/api/users/:id", "Get user info by ID", auth=True),
            self.req("Get User by Username", "GET", "/api/users/username/:username", "Get user by username",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Update User", "PUT", "/api/users/:id", "Update user profile", auth=True,
                     body={"full_name": "John Doe Updated", "bio": "Software engineer", "interests": ["ml", "ai"]}),
            self.req("Update Password", "PUT", "/api/users/:id/password", "Change password", auth=True,
                     body={"old_password": "OldPass123", "new_password": "NewPass456"}),
            self.req("Deactivate User", "DELETE", "/api/users/:id", "Deactivate account", auth=True),
            self.req("Search Users", "GET", "/api/users/search", "Search users",
                     query=[{"key": "q", "value": "john"}, {"key": "space_id", "value": "{{space_id}}"}],
                     test=PAGINATED_TEST)
        ])

    def gen_posts(self):
        store_post = """pm.test("Status 201", () => pm.response.to.have.status(201));
pm.test("Post created", () => {
    const data = pm.response.json();
    if (data.data && data.data.id) pm.environment.set('test_post_id', data.data.id);
});"""
        return self.folder("Posts & Comments", "Posts, comments, likes, feeds", [
            self.req("Create Post", "POST", "/api/posts", "Create new post", auth=True,
                     body={"space_id": "{{space_id}}", "content": "My first post!", "tags": ["intro"],
                           "visibility": "public"}, test=store_post),
            self.req("Get Post", "GET", "/api/posts/:id", "Get post by ID (increments views)"),
            self.req("Delete Post", "DELETE", "/api/posts/:id", "Delete post", auth=True),
            self.req("Get User Feed", "GET", "/api/posts/feed", "Personalized feed", auth=True,
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get User Posts", "GET", "/api/posts/user/:user_id", "Posts by user",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Community Posts", "GET", "/api/posts/community/:community_id", "Posts in community",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Group Posts", "GET", "/api/posts/group/:group_id", "Posts in group",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Trending Posts", "GET", "/api/posts/trending", "Trending posts",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get User Liked Posts", "GET", "/api/posts/liked", "Posts liked by user", auth=True,
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Search Posts", "GET", "/api/posts/search", "Basic search",
                     query=[{"key": "q", "value": "tech"}, {"key": "space_id", "value": "{{space_id}}"},
                            {"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Advanced Search Posts", "GET", "/api/posts/advanced-search", "Advanced search",
                     query=[{"key": "q", "value": "programming"}, {"key": "space_id", "value": "{{space_id}}"},
                            {"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Post Comments", "GET", "/api/posts/:id/comments", "Get comments"),
            self.req("Create Comment", "POST", "/api/posts/:id/comments", "Add comment", auth=True,
                     body={"content": "Great post!"}),
            self.req("Toggle Comment Like", "POST", "/api/comments/:id/like", "Like/unlike comment", auth=True),
            self.req("Get Post Likes", "GET", "/api/posts/:id/likes", "Users who liked post"),
            self.req("Toggle Post Like", "POST", "/api/posts/:id/like", "Like/unlike post", auth=True),
            self.req("Create Repost", "POST", "/api/posts/:id/repost", "Repost to profile", auth=True),
            self.req("Pin Post", "PUT", "/api/posts/:id/pin", "Pin post to top", auth=True)
        ])

    def gen_sessions(self):
        return self.folder("Sessions", "Auth session management", [
            self.req("Get Session", "GET", "/api/sessions/:id", "Get session by ID", auth=True)
        ])

    def gen_communities(self):
        store_comm = """pm.test("Status 201", () => pm.response.to.have.status(201));
pm.test("Community created", () => {
    const data = pm.response.json();
    if (data.data && data.data.id) pm.environment.set('test_community_id', data.data.id);
});"""
        return self.folder("Communities", "Community management and membership", [
            self.req("Create Community", "POST", "/api/communities", "Create community", auth=True,
                     body={"space_id": "{{space_id}}", "name": "Tech Enthusiasts", "slug": "tech-enthusiasts",
                           "description": "Tech community", "category": "Technology",
                           "tags": ["tech", "innovation"], "is_private": False}, test=store_comm),
            self.req("List Communities", "GET", "/api/communities", "List all communities",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Community", "GET", "/api/communities/:id", "Get community by ID"),
            self.req("Get Community by Slug", "GET", "/api/communities/slug/:slug", "Get by slug"),
            self.req("Update Community", "PUT", "/api/communities/:id", "Update community", auth=True,
                     body={"name": "Tech Enthusiasts Updated", "description": "Updated desc"}),
            self.req("Search Communities", "GET", "/api/communities/search", "Search communities",
                     query=[{"key": "q", "value": "tech"}, {"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Community Categories", "GET", "/api/communities/categories", "List categories"),
            self.req("Join Community", "POST", "/api/communities/:id/join", "Join as member", auth=True),
            self.req("Leave Community", "POST", "/api/communities/:id/leave", "Leave community", auth=True),
            self.req("Get User Communities", "GET", "/api/users/communities", "User's communities", auth=True),
            self.req("Get Community Members", "GET", "/api/communities/:id/members", "List members",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}]),
            self.req("Get Community Moderators", "GET", "/api/communities/:id/moderators", "List moderators"),
            self.req("Get Community Admins", "GET", "/api/communities/:id/admins", "List admins"),
            self.req("Add Community Moderator", "POST", "/api/communities/:id/moderators", "Add moderator", auth=True,
                     body={"user_id": "{{user_id}}"}),
            self.req("Remove Community Moderator", "DELETE", "/api/communities/:id/moderators/:userId",
                     "Remove moderator", auth=True)
        ])

    def gen_groups(self):
        store_grp = """pm.test("Status 201", () => pm.response.to.have.status(201));
pm.test("Group created", () => {
    const data = pm.response.json();
    if (data.data && data.data.id) pm.environment.set('test_group_id', data.data.id);
});"""
        return self.folder("Groups", "Project groups, roles, applications", [
            self.req("Create Group", "POST", "/api/groups", "Create project group", auth=True,
                     body={"space_id": "{{space_id}}", "name": "AI Research", "slug": "ai-research",
                           "description": "AI research group", "is_private": False,
                           "tags": ["ai", "research"]}, test=store_grp),
            self.req("List Groups", "GET", "/api/groups", "List all groups",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Group", "GET", "/api/groups/:id", "Get group by ID"),
            self.req("Update Group", "PUT", "/api/groups/:id", "Update group", auth=True,
                     body={"name": "AI Research Updated", "description": "Updated desc"}),
            self.req("Search Groups", "GET", "/api/groups/search", "Search groups",
                     query=[{"key": "q", "value": "ai"}, {"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Join Group", "POST", "/api/groups/:id/join", "Join group", auth=True),
            self.req("Leave Group", "POST", "/api/groups/:id/leave", "Leave group", auth=True),
            self.req("Get User Groups", "GET", "/api/users/groups", "User's groups", auth=True),
            self.req("Get Group Join Requests", "GET", "/api/groups/:id/join-requests", "List join requests", auth=True),
            self.req("Get Project Roles", "GET", "/api/groups/:id/roles", "List project roles"),
            self.req("Create Project Role", "POST", "/api/groups/:id/roles", "Create role", auth=True,
                     body={"title": "Frontend Developer", "description": "React developer needed",
                           "requirements": "2+ years React", "slots": 2}),
            self.req("Apply for Project Role", "POST", "/api/roles/:roleId/apply", "Apply for role", auth=True,
                     body={"cover_letter": "I'm interested in this role"}),
            self.req("Get Role Applications", "GET", "/api/groups/:id/applications", "List applications", auth=True),
            self.req("Add Group Admin", "POST", "/api/groups/:id/admins", "Add admin", auth=True,
                     body={"user_id": "{{user_id}}"}),
            self.req("Remove Group Admin", "DELETE", "/api/groups/:id/admins/:userId", "Remove admin", auth=True),
            self.req("Add Group Moderator", "POST", "/api/groups/:id/moderators", "Add moderator", auth=True,
                     body={"user_id": "{{user_id}}"}),
            self.req("Remove Group Moderator", "DELETE", "/api/groups/:id/moderators/:userId", "Remove mod", auth=True),
            self.req("Update Group Member Role", "PUT", "/api/groups/:id/members/:userId/role", "Update member role",
                     auth=True, body={"role": "moderator"})
        ])

    def gen_messaging(self):
        store_conv = """pm.test("Status 201", () => pm.response.to.have.status(201));
pm.test("Conversation created", () => {
    const data = pm.response.json();
    if (data.data && data.data.id) pm.environment.set('test_conversation_id', data.data.id);
});"""
        return self.folder("Messaging", "Conversations and messages", [
            self.req("Create Conversation", "POST", "/api/conversations", "Create group conversation", auth=True,
                     body={"participant_ids": ["{{user_id}}"], "name": "Project Discussion"}, test=store_conv),
            self.req("Get User Conversations", "GET", "/api/conversations", "List conversations", auth=True),
            self.req("Get Conversation", "GET", "/api/conversations/:id", "Get conversation details", auth=True),
            self.req("Get or Create Direct Conversation", "POST", "/api/conversations/direct",
                     "Get/create DM", auth=True, body={"participant_id": "{{user_id}}"}),
            self.req("Leave Conversation", "POST", "/api/conversations/:id/leave", "Leave conversation", auth=True),
            self.req("Update Participant Settings", "PUT", "/api/conversations/:id/settings", "Update settings",
                     auth=True, body={"muted": False, "notifications_enabled": True}),
            self.req("Get Conversation Participants", "GET", "/api/conversations/:id/participants",
                     "List participants", auth=True),
            self.req("Add Conversation Participants", "POST", "/api/conversations/:id/participants",
                     "Add participants", auth=True, body={"user_ids": ["{{user_id}}"]}),
            self.req("Send Message", "POST", "/api/conversations/:id/messages", "Send message", auth=True,
                     body={"content": "Hello!", "message_type": "text"}),
            self.req("Get Conversation Messages", "GET", "/api/conversations/:id/messages", "Get messages", auth=True,
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "50"}]),
            self.req("Mark Messages as Read", "POST", "/api/conversations/:id/read", "Mark as read", auth=True),
            self.req("Get Unread Count", "GET", "/api/conversations/:id/unread", "Get unread count", auth=True),
            self.req("Get Message", "GET", "/api/messages/:id", "Get message by ID", auth=True),
            self.req("Delete Message", "DELETE", "/api/messages/:id", "Delete message", auth=True),
            self.req("Add Message Reaction", "POST", "/api/messages/:id/reactions", "Add reaction", auth=True,
                     body={"emoji": "👍"}),
            self.req("Remove Message Reaction", "DELETE", "/api/messages/:id/reactions/:emoji",
                     "Remove reaction", auth=True)
        ])

    def gen_notifications(self):
        return self.folder("Notifications", "User notifications", [
            self.req("Create Notification", "POST", "/api/notifications", "Create notification", auth=True,
                     body={"recipient_id": "{{user_id}}", "type": "mention", "content": "You were mentioned",
                           "entity_type": "post", "entity_id": "{{test_post_id}}"}),
            self.req("Get User Notifications", "GET", "/api/notifications", "List notifications", auth=True,
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Mark as Read", "PUT", "/api/notifications/:id/read", "Mark notification as read", auth=True),
            self.req("Mark All as Read", "PUT", "/api/notifications/read-all", "Mark all as read", auth=True),
            self.req("Delete Notification", "DELETE", "/api/notifications/:id", "Delete notification", auth=True),
            self.req("Get Unread Count", "GET", "/api/notifications/unread-count", "Get unread count", auth=True)
        ])

    def gen_events(self):
        store_evt = """pm.test("Status 201", () => pm.response.to.have.status(201));
pm.test("Event created", () => {
    const data = pm.response.json();
    if (data.data && data.data.id) pm.environment.set('test_event_id', data.data.id);
});"""
        return self.folder("Events", "Event management and registration", [
            self.req("Create Event", "POST", "/api/events", "Create event", auth=True,
                     body={"space_id": "{{space_id}}", "title": "Tech Meetup", "description": "Monthly tech meetup",
                           "category": "Networking", "location": "Campus Hall", "event_type": "in_person",
                           "start_time": "2025-12-01T18:00:00Z", "end_time": "2025-12-01T20:00:00Z",
                           "max_attendees": 50}, test=store_evt),
            self.req("List Events", "GET", "/api/events", "List all events",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Event", "GET", "/api/events/:id", "Get event by ID"),
            self.req("Update Event", "PUT", "/api/events/:id", "Update event", auth=True,
                     body={"title": "Tech Meetup Updated", "description": "Updated description"}),
            self.req("Update Event Status", "PUT", "/api/events/:id/status", "Update status", auth=True,
                     body={"status": "published"}),
            self.req("Get Upcoming Events", "GET", "/api/events/upcoming", "List upcoming events"),
            self.req("Search Events", "GET", "/api/events/search", "Search events",
                     query=[{"key": "q", "value": "tech"}, {"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Event Categories", "GET", "/api/events/categories", "List event categories"),
            self.req("Get Event Attendees", "GET", "/api/events/:id/attendees", "List attendees"),
            self.req("Get Event Co-organizers", "GET", "/api/events/:id/co-organizers", "List co-organizers"),
            self.req("Register for Event", "POST", "/api/events/:id/register", "Register", auth=True),
            self.req("Unregister from Event", "POST", "/api/events/:id/unregister", "Unregister", auth=True),
            self.req("Add Event Co-organizer", "POST", "/api/events/:id/co-organizers", "Add co-organizer", auth=True,
                     body={"user_id": "{{user_id}}"}),
            self.req("Remove Event Co-organizer", "DELETE", "/api/events/:id/co-organizers/:user_id",
                     "Remove co-organizer", auth=True),
            self.req("Mark Event Attendance", "POST", "/api/events/:id/attendance/:user_id", "Mark attendance",
                     auth=True, body={"attended": True}),
            self.req("Get User Events", "GET", "/api/users/events", "User's events", auth=True)
        ])

    def gen_announcements(self):
        return self.folder("Announcements", "System announcements", [
            self.req("Create Announcement", "POST", "/api/announcements", "Create announcement", auth=True,
                     body={"space_id": "{{space_id}}", "title": "Welcome!", "content": "Welcome to our platform",
                           "priority": "normal", "target_audience": "all"}),
            self.req("List Announcements", "GET", "/api/announcements", "List announcements",
                     query=[{"key": "page", "value": "1"}, {"key": "limit", "value": "20"}], test=PAGINATED_TEST),
            self.req("Get Announcement", "GET", "/api/announcements/:id", "Get announcement by ID"),
            self.req("Update Announcement", "PUT", "/api/announcements/:id", "Update announcement", auth=True,
                     body={"title": "Updated Title", "content": "Updated content"}),
            self.req("Update Announcement Status", "PUT", "/api/announcements/:id/status", "Update status", auth=True,
                     body={"status": "published"})
        ])

    def gen_mentorship(self):
        return self.folder("Mentorship", "Mentors, tutors, sessions", [
            # Mentors
            self.req("Search Mentors", "GET", "/api/mentorship/mentors/search", "Search mentors",
                     query=[{"key": "q", "value": "python"}, {"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Mentor Profile", "GET", "/api/mentorship/mentors/profile/:id", "Get mentor profile"),
            self.req("Get Mentor Reviews", "GET", "/api/mentorship/mentors/:id/reviews", "Get reviews"),
            self.req("Create Mentor Profile", "POST", "/api/mentorship/mentors/profile", "Create profile", auth=True,
                     body={"space_id": "{{space_id}}", "bio": "Experienced mentor", "industries": ["Technology"],
                           "specialties": ["Python", "AI"], "availability": "weekends"}),
            self.req("Update Mentor Availability", "PUT", "/api/mentorship/mentors/profile/:id/availability",
                     "Update availability", auth=True, body={"availability": "weekdays"}),
            # Mentor applications
            self.req("Create Mentor Application", "POST", "/api/mentorship/mentors/applications/",
                     "Apply as mentor", auth=True,
                     body={"space_id": "{{space_id}}", "bio": "I want to mentor", "industries": ["Tech"],
                           "specialties": ["Python"], "reason": "Give back to community"}),
            self.req("Get Mentor Application", "GET", "/api/mentorship/mentors/applications/:id", "Get application"),
            self.req("Get Pending Mentor Applications", "GET", "/api/mentorship/mentors/applications/pending",
                     "List pending"),
            self.req("Update Mentor Application", "PUT", "/api/mentorship/mentors/applications/:id",
                     "Update application", auth=True, body={"bio": "Updated bio"}),
            # Tutors
            self.req("Search Tutors", "GET", "/api/mentorship/tutors/search", "Search tutors",
                     query=[{"key": "q", "value": "math"}, {"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Tutor Profile", "GET", "/api/mentorship/tutors/profile/:id", "Get tutor profile"),
            self.req("Get Tutor Reviews", "GET", "/api/mentorship/tutors/:id/reviews", "Get reviews"),
            self.req("Create Tutor Profile", "POST", "/api/mentorship/tutors/profile", "Create profile", auth=True,
                     body={"space_id": "{{space_id}}", "bio": "Math tutor", "subjects": ["Calculus", "Algebra"],
                           "hourly_rate": 25.00, "availability": "evenings"}),
            self.req("Update Tutor Availability", "PUT", "/api/mentorship/tutors/profile/:id/availability",
                     "Update availability", auth=True, body={"availability": "weekends"}),
            # Tutor applications
            self.req("Create Tutor Application", "POST", "/api/mentorship/tutors/applications/",
                     "Apply as tutor", auth=True,
                     body={"space_id": "{{space_id}}", "bio": "I want to tutor", "subjects": ["Math"],
                           "reason": "Help students"}),
            self.req("Get Tutor Application", "GET", "/api/mentorship/tutors/applications/:id", "Get application"),
            self.req("Get Pending Tutor Applications", "GET", "/api/mentorship/tutors/applications/pending",
                     "List pending"),
            self.req("Update Tutor Application", "PUT", "/api/mentorship/tutors/applications/:id",
                     "Update application", auth=True, body={"bio": "Updated bio"}),
            # Mentoring sessions
            self.req("Create Mentoring Session", "POST", "/api/mentorship/mentoring/sessions", "Create session",
                     auth=True,
                     body={"mentor_id": "{{user_id}}", "scheduled_at": "2025-12-01T15:00:00Z",
                           "duration_minutes": 60, "topic": "Career guidance"}),
            self.req("Get User Mentoring Sessions", "GET", "/api/mentorship/mentoring/sessions",
                     "List user sessions", auth=True),
            self.req("Get Mentoring Session", "GET", "/api/mentorship/mentoring/sessions/:id", "Get session"),
            self.req("Update Mentoring Session Status", "PUT", "/api/mentorship/mentoring/sessions/:id/status",
                     "Update status", auth=True, body={"status": "confirmed"}),
            self.req("Add Mentoring Session Meeting Link", "PUT",
                     "/api/mentorship/mentoring/sessions/:id/meeting-link", "Add meeting link", auth=True,
                     body={"meeting_link": "https://zoom.us/j/123456"}),
            self.req("Rate Mentoring Session", "POST", "/api/mentorship/mentoring/sessions/:id/rate",
                     "Rate session", auth=True, body={"rating": 5, "review": "Great session!"}),
            # Tutoring sessions
            self.req("Create Tutoring Session", "POST", "/api/mentorship/tutoring/sessions", "Create session",
                     auth=True,
                     body={"tutor_id": "{{user_id}}", "scheduled_at": "2025-12-01T16:00:00Z",
                           "duration_minutes": 90, "subject": "Calculus"}),
            self.req("Get User Tutoring Sessions", "GET", "/api/mentorship/tutoring/sessions",
                     "List user sessions", auth=True),
            self.req("Get Tutoring Session", "GET", "/api/mentorship/tutoring/sessions/:id", "Get session"),
            self.req("Update Tutoring Session Status", "PUT", "/api/mentorship/tutoring/sessions/:id/status",
                     "Update status", auth=True, body={"status": "confirmed"}),
            self.req("Add Tutoring Session Meeting Link", "PUT", "/api/mentorship/tutoring/sessions/:id/meeting-link",
                     "Add meeting link", auth=True, body={"meeting_link": "https://meet.google.com/xyz"}),
            self.req("Rate Tutoring Session", "POST", "/api/mentorship/tutoring/sessions/:id/rate",
                     "Rate session", auth=True, body={"rating": 5, "review": "Excellent tutor!"})
        ])

    def gen_analytics(self):
        return self.folder("Analytics", "Reports, metrics, statistics", [
            # Reports
            self.req("Create Report", "POST", "/api/analytics/reports", "Create content report", auth=True,
                     body={"entity_type": "post", "entity_id": "{{test_post_id}}", "reason": "spam",
                           "description": "This is spam"}),
            self.req("Get Report", "GET", "/api/analytics/reports/:id", "Get report by ID"),
            self.req("Get Reports by Content", "GET", "/api/analytics/reports/by-content", "Get reports for content",
                     query=[{"key": "entity_type", "value": "post"}, {"key": "entity_id", "value": "{{test_post_id}}"}]),
            self.req("Get Pending Reports", "GET", "/api/analytics/reports/pending", "List pending reports"),
            self.req("Update Report", "PUT", "/api/analytics/reports/:id", "Update report status", auth=True,
                     body={"status": "reviewed", "resolution": "removed"}),
            # Moderation
            self.req("Get Moderation Queue", "GET", "/api/analytics/moderation/queue", "Get moderation queue"),
            self.req("Get Content Moderation Stats", "GET", "/api/analytics/moderation/stats", "Get mod stats"),
            # Metrics
            self.req("Get System Metrics", "GET", "/api/analytics/metrics/system", "Get system metrics"),
            self.req("Get Space Stats", "GET", "/api/analytics/metrics/space", "Get space statistics",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Historical Metrics", "GET", "/api/analytics/metrics/historical", "Get historical metrics",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "days", "value": "30"}]),
            # Engagement
            self.req("Get Engagement Metrics", "GET", "/api/analytics/engagement/metrics", "Get engagement metrics",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get User Activity Stats", "GET", "/api/analytics/activity/stats", "Get activity stats",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get User Growth", "GET", "/api/analytics/users/growth", "Get user growth metrics",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "days", "value": "30"}]),
            self.req("Get User Engagement Ranking", "GET", "/api/analytics/users/ranking", "Get user ranking",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "limit", "value": "10"}]),
            # Top content
            self.req("Get Top Posts", "GET", "/api/analytics/top/posts", "Get top posts",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "limit", "value": "10"}]),
            self.req("Get Top Communities", "GET", "/api/analytics/top/communities", "Get top communities",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "limit", "value": "10"}]),
            self.req("Get Top Groups", "GET", "/api/analytics/top/groups", "Get top groups",
                     query=[{"key": "space_id", "value": "{{space_id}}"}, {"key": "limit", "value": "10"}]),
            # Mentorship analytics
            self.req("Get Mentoring Stats", "GET", "/api/analytics/mentorship/mentoring", "Get mentoring stats",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Tutoring Stats", "GET", "/api/analytics/mentorship/tutoring", "Get tutoring stats",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Popular Industries", "GET", "/api/analytics/mentorship/industries", "Get popular industries",
                     query=[{"key": "space_id", "value": "{{space_id}}"}]),
            self.req("Get Popular Subjects", "GET", "/api/analytics/mentorship/subjects", "Get popular subjects",
                     query=[{"key": "space_id", "value": "{{space_id}}"}])
        ])

    def generate(self):
        """Generate complete collection."""
        self.collection["item"] = [
            self.gen_health(),
            self.gen_auth(),
            self.gen_users(),
            self.gen_posts(),
            self.gen_sessions(),
            self.gen_communities(),
            self.gen_groups(),
            self.gen_messaging(),
            self.gen_notifications(),
            self.gen_events(),
            self.gen_announcements(),
            self.gen_mentorship(),
            self.gen_analytics()
        ]
        return self.collection


def main():
    print("Generating Complete Connect Server Postman Collection...")
    print("=" * 70)

    gen = PostmanGenerator()
    collection = gen.generate()

    # Write collection
    import os
    os.makedirs("db", exist_ok=True)
    with open("db/postman_collection.json", 'w') as f:
        json.dump(collection, f, indent=2)

    # Count endpoints
    total_requests = sum(len(folder.get('item', [])) for folder in collection['item'])

    print(f"\n✅ Collection generated: db/postman_collection.json")
    print(f"   Folders: {len(collection['item'])}")
    print(f"   Total endpoints: {total_requests}")
    print("\nCollection ready for import into Postman!")


if __name__ == "__main__":
    main()
