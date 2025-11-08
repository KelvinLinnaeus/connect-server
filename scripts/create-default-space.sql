-- Create default space if none exists
-- This script creates a default space for development and testing

-- Check if any spaces exist, if not create a default one
INSERT INTO spaces (id, name, slug, description, settings, status, is_active)
SELECT
    gen_random_uuid(),
    'Default University',
    'default-university',
    'Default space for university community',
    '{"theme": "default", "features": {"communities": true, "events": true, "messaging": true}}'::jsonb,
    'active',
    true
WHERE NOT EXISTS (SELECT 1 FROM spaces LIMIT 1);

-- Display the space ID for configuration
SELECT
    id as space_id,
    name,
    slug,
    'Add this space_id to your .env.local as VITE_DEFAULT_SPACE_ID=' || id as configuration_instruction
FROM spaces
ORDER BY created_at ASC
LIMIT 1;
