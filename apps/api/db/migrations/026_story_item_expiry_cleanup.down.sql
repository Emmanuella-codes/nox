DROP INDEX IF EXISTS story_items_story_id_expires_at_idx;
DROP INDEX IF EXISTS story_items_expires_at_idx;

ALTER TABLE story_items
  DROP COLUMN IF EXISTS expires_at;
