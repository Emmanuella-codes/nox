DROP INDEX IF EXISTS notifications_set_id_idx;
ALTER TABLE notifications
  DROP COLUMN IF EXISTS set_comment_id,
  DROP COLUMN IF EXISTS set_id;
