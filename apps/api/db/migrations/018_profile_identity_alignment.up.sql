UPDATE personas
SET persona_type = 'visible'
WHERE persona_type IS DISTINCT FROM 'visible';

DROP INDEX IF EXISTS personas_one_ghost_per_user_idx;

CREATE UNIQUE INDEX IF NOT EXISTS personas_one_profile_per_user_idx
  ON personas(user_id);

ALTER TABLE posts
  ALTER COLUMN persona_id DROP NOT NULL;

UPDATE posts p
SET persona_id = pe.id
FROM personas pe
WHERE p.author_user_id = pe.user_id
  AND p.persona_id IS NULL;

ALTER TABLE posts
  ALTER COLUMN persona_id SET NOT NULL;

ALTER TABLE posts
  DROP CONSTRAINT IF EXISTS posts_posting_mode_persona_check;

ALTER TABLE posts
  ADD CONSTRAINT posts_posting_mode_persona_check CHECK (
    posting_mode IN ('public', 'anonymous') AND persona_id IS NOT NULL
  );

ALTER TABLE comments
  ADD COLUMN IF NOT EXISTS author_user_id UUID REFERENCES users(id) ON DELETE CASCADE;

UPDATE comments c
SET author_user_id = pe.user_id
FROM personas pe
WHERE c.persona_id = pe.id
  AND c.author_user_id IS NULL;

ALTER TABLE comments
  ALTER COLUMN author_user_id SET NOT NULL;

ALTER TABLE anonymous_thread_identities
  ADD COLUMN IF NOT EXISTS anonymous_avatar_key TEXT NOT NULL DEFAULT 'mask_01';

DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conname = 'anonymous_thread_identities_thread_id_persona_id_key'
  ) THEN
    ALTER TABLE anonymous_thread_identities
      DROP CONSTRAINT anonymous_thread_identities_thread_id_persona_id_key;
  END IF;
END $$;

ALTER TABLE anonymous_thread_identities
  ADD CONSTRAINT anonymous_thread_identities_thread_id_user_id_key UNIQUE (thread_id, user_id);

CREATE INDEX IF NOT EXISTS anonymous_thread_identities_user_id_idx
  ON anonymous_thread_identities(user_id);
