ALTER TABLE posts
  DROP CONSTRAINT IF EXISTS posts_posting_mode_persona_check;

ALTER TABLE posts
  DROP CONSTRAINT IF EXISTS posts_posting_mode_check;

DROP INDEX IF EXISTS posts_author_user_id_idx;

ALTER TABLE posts
  DROP COLUMN IF EXISTS posting_mode;

ALTER TABLE posts
  DROP COLUMN IF EXISTS author_user_id;
