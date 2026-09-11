DROP INDEX IF EXISTS idx_media_assets_stuck;
DROP INDEX IF EXISTS idx_media_moderation_status;
DROP TABLE IF EXISTS media_moderation;
ALTER TABLE media_assets
  DROP COLUMN IF EXISTS last_processing_started_at,
  DROP COLUMN IF EXISTS last_processing_error,
  DROP COLUMN IF EXISTS processing_attempts;
