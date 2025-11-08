-- UNIVYN Database Migration
-- Version: 006_events_announcements DOWN
-- Description: Drop events, announcements, and past questions

BEGIN;

DROP TABLE IF EXISTS past_questions;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS event_attendees;
DROP TABLE IF EXISTS event_co_organizers;
DROP TABLE IF EXISTS events;

COMMIT;