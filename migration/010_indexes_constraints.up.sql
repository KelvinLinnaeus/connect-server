-- UNIVYN Database Migration
-- Version: 010_indexes_constraints UP
-- Description: Create indexes and additional constraints for performance
BEGIN;
-- Spaces indexes
CREATE INDEX idx_spaces_type ON spaces(type);
CREATE INDEX idx_spaces_is_active ON spaces(status);
-- Users indexes
CREATE INDEX idx_users_space_id ON users(space_id);
CREATE INDEX idx_users_mentor_status ON users(mentor_status);
CREATE INDEX idx_users_tutor_status ON users(tutor_status);
CREATE INDEX idx_users_space_email ON users(space_id, email);
-- Communities indexes
CREATE INDEX idx_communities_space_id ON communities(space_id);
CREATE INDEX idx_communities_category ON communities(category);
CREATE INDEX idx_communities_created_by ON communities(created_by);
CREATE INDEX idx_communities_is_public ON communities(is_public);
-- Groups indexes
CREATE INDEX idx_groups_space_id ON groups(space_id);
CREATE INDEX idx_groups_community_id ON groups(community_id);
CREATE INDEX idx_groups_category ON groups(category);
CREATE INDEX idx_groups_group_type ON groups(group_type);
CREATE INDEX idx_groups_status ON groups(status);
CREATE INDEX idx_groups_created_by ON groups(created_by);
-- Community membership indexes
CREATE INDEX idx_community_members_community_id ON community_members(community_id);
CREATE INDEX idx_community_members_user_id ON community_members(user_id);
-- Group membership indexes
CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_role ON group_members(role);
-- Posts indexes
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_space_id ON posts(space_id);
CREATE INDEX idx_posts_community_id ON posts(community_id);
CREATE INDEX idx_posts_group_id ON posts(group_id);
CREATE INDEX idx_posts_parent_post_id ON posts(parent_post_id);
CREATE INDEX idx_posts_quoted_post_id ON posts(quoted_post_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
CREATE INDEX idx_posts_space_created ON posts(space_id, created_at);
CREATE INDEX idx_posts_status ON posts(status);
CREATE INDEX idx_posts_is_pinned ON posts(is_pinned);
-- Comments indexes
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id);
CREATE INDEX idx_comments_created_at ON comments(created_at);
-- Likes indexes
CREATE INDEX idx_likes_user_id ON likes(user_id);
CREATE INDEX idx_likes_post_id ON likes(post_id);
CREATE INDEX idx_likes_comment_id ON likes(comment_id);
-- Follows indexes
CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
CREATE INDEX idx_follows_space_id ON follows(space_id);
-- Notifications indexes
CREATE INDEX idx_notifications_to_user_id ON notifications(to_user_id);
CREATE INDEX idx_notifications_from_user_id ON notifications(from_user_id);
CREATE INDEX idx_notifications_type ON notifications(type);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_priority ON notifications(priority);
-- Tutor/Mentor applications indexes
CREATE INDEX idx_tutor_applications_applicant_id ON tutor_applications(applicant_id);
CREATE INDEX idx_tutor_applications_space_id ON tutor_applications(space_id);
CREATE INDEX idx_tutor_applications_status ON tutor_applications(status);
CREATE INDEX idx_tutor_applications_submitted_at ON tutor_applications(submitted_at);
CREATE INDEX idx_mentor_applications_applicant_id ON mentor_applications(applicant_id);
CREATE INDEX idx_mentor_applications_space_id ON mentor_applications(space_id);
CREATE INDEX idx_mentor_applications_status ON mentor_applications(status);
CREATE INDEX idx_mentor_applications_industry ON mentor_applications(industry);
CREATE INDEX idx_mentor_applications_submitted_at ON mentor_applications(submitted_at);
-- Tutor/Mentor profiles indexes
CREATE INDEX idx_tutor_profiles_space_id ON tutor_profiles(space_id);
CREATE INDEX idx_tutor_profiles_subjects ON tutor_profiles USING GIN(subjects);
CREATE INDEX idx_tutor_profiles_rating ON tutor_profiles(rating);
CREATE INDEX idx_tutor_profiles_is_available ON tutor_profiles(is_available);
CREATE INDEX idx_mentor_profiles_space_id ON mentor_profiles(space_id);
CREATE INDEX idx_mentor_profiles_industry ON mentor_profiles(industry);
CREATE INDEX idx_mentor_profiles_rating ON mentor_profiles(rating);
CREATE INDEX idx_mentor_profiles_is_available ON mentor_profiles(is_available);
-- Session indexes
CREATE INDEX idx_tutoring_sessions_tutor_id ON tutoring_sessions(tutor_id);
CREATE INDEX idx_tutoring_sessions_student_id ON tutoring_sessions(student_id);
CREATE INDEX idx_tutoring_sessions_space_id ON tutoring_sessions(space_id);
CREATE INDEX idx_tutoring_sessions_status ON tutoring_sessions(status);
CREATE INDEX idx_tutoring_sessions_scheduled_at ON tutoring_sessions(scheduled_at);
CREATE INDEX idx_tutoring_sessions_subject ON tutoring_sessions(subject);
CREATE INDEX idx_mentoring_sessions_mentor_id ON mentoring_sessions(mentor_id);
CREATE INDEX idx_mentoring_sessions_mentee_id ON mentoring_sessions(mentee_id);
CREATE INDEX idx_mentoring_sessions_space_id ON mentoring_sessions(space_id);
CREATE INDEX idx_mentoring_sessions_status ON mentoring_sessions(status);
CREATE INDEX idx_mentoring_sessions_scheduled_at ON mentoring_sessions(scheduled_at);
CREATE INDEX idx_mentoring_sessions_topic ON mentoring_sessions(topic);
-- Messaging indexes
CREATE INDEX idx_conversations_space_id ON conversations(space_id);
CREATE INDEX idx_conversations_last_message_at ON conversations(last_message_at);
CREATE INDEX idx_conversations_conversation_type ON conversations(conversation_type);
CREATE INDEX idx_conversations_is_active ON conversations(is_active);
CREATE INDEX idx_conversation_participants_conversation_id ON conversation_participants(conversation_id);
CREATE INDEX idx_conversation_participants_user_id ON conversation_participants(user_id);
CREATE INDEX idx_conversation_participants_role ON conversation_participants(role);
CREATE INDEX idx_conversation_participants_is_active ON conversation_participants(is_active);
CREATE INDEX idx_conversation_participants_user_active ON conversation_participants(user_id, is_active);
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_sender_id ON messages(sender_id);
CREATE INDEX idx_messages_recipient_id ON messages(recipient_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
CREATE INDEX idx_messages_conversation_created ON messages(conversation_id, created_at);
CREATE INDEX idx_messages_is_read ON messages(is_read);
-- Events indexes
CREATE INDEX idx_events_space_id ON events(space_id);
CREATE INDEX idx_events_organizer ON events(organizer);
CREATE INDEX idx_events_category ON events(category);
CREATE INDEX idx_events_start_date ON events(start_date);
CREATE INDEX idx_events_status ON events(status);
CREATE INDEX idx_events_is_public ON events(is_public);
CREATE INDEX idx_event_attendees_event_id ON event_attendees(event_id);
CREATE INDEX idx_event_attendees_user_id ON event_attendees(user_id);
CREATE INDEX idx_event_attendees_status ON event_attendees(status);
-- Announcements indexes
CREATE INDEX idx_announcements_space_id ON announcements(space_id);
CREATE INDEX idx_announcements_type ON announcements(type);
CREATE INDEX idx_announcements_priority ON announcements(priority);
CREATE INDEX idx_announcements_status ON announcements(status);
CREATE INDEX idx_announcements_scheduled_for ON announcements(scheduled_for);
CREATE INDEX idx_announcements_author_id ON announcements(author_id);
CREATE INDEX idx_announcements_is_pinned ON announcements(is_pinned);
-- Past questions indexes
CREATE INDEX idx_past_questions_space_id ON past_questions(space_id);
CREATE INDEX idx_past_questions_course_code ON past_questions(course_code);
CREATE INDEX idx_past_questions_department ON past_questions(department);
CREATE INDEX idx_past_questions_academic_year ON past_questions(academic_year);
CREATE INDEX idx_past_questions_uploaded_by ON past_questions(uploaded_by);
CREATE INDEX idx_past_questions_verified ON past_questions(verified);
-- Project roles indexes
CREATE INDEX idx_group_roles_group_id ON group_roles(group_id);
CREATE INDEX idx_group_roles_name ON group_roles(name);
CREATE INDEX idx_group_applications_role_id ON group_applications(role_id);
CREATE INDEX idx_group_applications_user_id ON group_applications(user_id);
CREATE INDEX idx_group_applications_status ON group_applications(status);
-- Reports indexes
CREATE INDEX idx_reports_space_id ON reports(space_id);
CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_content_type ON reports(content_type);
CREATE INDEX idx_reports_content_id ON reports(content_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_priority ON reports(priority);
CREATE INDEX idx_reports_created_at ON reports(created_at);
-- User sessions indexes (renamed from sessions)
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_refresh_token ON user_sessions(refresh_token);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
-- Note: system_metrics indexes removed (table removed - metrics derived from existing tables)
-- Files indexes
CREATE INDEX idx_files_space_id ON files(space_id);
CREATE INDEX idx_files_user_id ON files(user_id);
CREATE INDEX idx_files_mime_type ON files(mime_type);
CREATE INDEX idx_files_uploaded_at ON files(uploaded_at);
-- Trending topics indexes
CREATE INDEX idx_trending_topics_space_id ON trending_topics(space_id);
CREATE INDEX idx_trending_topics_name ON trending_topics(name);
CREATE INDEX idx_trending_topics_trend_score ON trending_topics(trend_score);
CREATE INDEX idx_trending_topics_recorded_at ON trending_topics(recorded_at);
-- Email queue indexes
CREATE INDEX idx_email_queue_space_id ON email_queue(space_id);
CREATE INDEX idx_email_queue_status ON email_queue(status);
CREATE INDEX idx_email_queue_created_at ON email_queue(created_at);
COMMIT;

