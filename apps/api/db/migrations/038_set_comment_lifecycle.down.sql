DROP TABLE IF EXISTS set_comment_likes;
ALTER TABLE set_comments
  DROP COLUMN IF EXISTS updated_at,
  DROP COLUMN IF EXISTS deleted_at;
