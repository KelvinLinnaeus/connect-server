-- Fix log_space_activity trigger function to handle different actor_id field names
-- This recreates the function to correctly handle posts (author_id), communities and groups (created_by)

CREATE OR REPLACE FUNCTION log_space_activity()
RETURNS TRIGGER AS $$
DECLARE
    v_actor_id UUID;
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Determine the actor_id based on the table
        -- Posts table uses author_id, communities and groups use created_by
        IF TG_TABLE_NAME = 'posts' THEN
            v_actor_id := NEW.author_id;
        ELSIF TG_TABLE_NAME IN ('communities', 'groups') THEN
            v_actor_id := NEW.created_by;
        ELSE
            v_actor_id := NULL;
        END IF;

        INSERT INTO space_activities (space_id, activity_type, actor_id, actor_name, description, metadata)
        VALUES (
            NEW.space_id,
            TG_ARGV[0],
            v_actor_id,
            (SELECT full_name FROM users WHERE id = v_actor_id),
            TG_ARGV[1],
            jsonb_build_object('resource_id', NEW.id)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
