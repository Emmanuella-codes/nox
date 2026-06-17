DROP INDEX IF EXISTS stories_expires_at_idx;

ALTER TABLE stories
  DROP COLUMN IF EXISTS expires_at;
