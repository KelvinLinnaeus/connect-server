-- Revert log_space_activity trigger function to original version
-- This should only be used if rolling back this specific fix

CREATE OR REPLACE FUNCTION log_space_activity()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO space_activities (space_id, activity_type, actor_id, actor_name, description, metadata)
        VALUES (
            NEW.space_id,
            TG_ARGV[0],
            NEW.author_id,
            (SELECT full_name FROM users WHERE id = NEW.author_id),
            TG_ARGV[1],
            jsonb_build_object('resource_id', NEW.id)
        );
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
