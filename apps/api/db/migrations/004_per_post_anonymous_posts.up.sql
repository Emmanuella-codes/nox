ALTER TABLE posts
  ADD COLUMN IF NOT EXISTS author_user_id UUID REFERENCES users(id) ON DELETE CASCADE;

UPDATE posts p
SET author_user_id = pe.user_id
FROM personas pe
WHERE p.persona_id = pe.id
  AND p.author_user_id IS NULL;

ALTER TABLE posts
  ALTER COLUMN author_user_id SET NOT NULL;

ALTER TABLE posts
  ADD COLUMN IF NOT EXISTS posting_mode TEXT;

UPDATE posts
SET posting_mode = 'public'
WHERE posting_mode IS NULL;

ALTER TABLE posts
  ALTER COLUMN posting_mode SET NOT NULL;

ALTER TABLE posts
  DROP CONSTRAINT IF EXISTS posts_posting_mode_check;

ALTER TABLE posts
  ADD CONSTRAINT posts_posting_mode_check CHECK (posting_mode IN ('public', 'anonymous'));

ALTER TABLE posts
  ALTER COLUMN persona_id DROP NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'posts_posting_mode_persona_check'
  ) THEN
    ALTER TABLE posts
      ADD CONSTRAINT posts_posting_mode_persona_check CHECK (
        (posting_mode = 'public' AND persona_id IS NOT NULL)
        OR (posting_mode = 'anonymous' AND persona_id IS NULL)
      );
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS posts_author_user_id_idx ON posts(author_user_id);
