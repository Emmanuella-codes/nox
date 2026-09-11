DROP TABLE IF EXISTS moderation_actions;

ALTER TABLE events
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE sets
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE story_items
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE stories
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE comments
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE posts
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;

ALTER TABLE personas
DROP COLUMN IF EXISTS moderated_by_user_id,
DROP COLUMN IF EXISTS moderated_at,
DROP COLUMN IF EXISTS moderation_reason,
DROP COLUMN IF EXISTS moderation_status;
