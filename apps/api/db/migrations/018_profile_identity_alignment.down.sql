DROP INDEX IF EXISTS anonymous_thread_identities_user_id_idx;

ALTER TABLE anonymous_thread_identities
  DROP CONSTRAINT IF EXISTS anonymous_thread_identities_thread_id_user_id_key;

ALTER TABLE anonymous_thread_identities
  DROP COLUMN IF EXISTS anonymous_avatar_key;

ALTER TABLE comments
  DROP COLUMN IF EXISTS author_user_id;

ALTER TABLE posts
  DROP CONSTRAINT IF EXISTS posts_posting_mode_persona_check;

ALTER TABLE posts
  ADD CONSTRAINT posts_posting_mode_persona_check CHECK (
    (posting_mode = 'public' AND persona_id IS NOT NULL)
    OR (posting_mode = 'anonymous' AND persona_id IS NULL)
  );

DROP INDEX IF EXISTS personas_one_profile_per_user_idx;

CREATE UNIQUE INDEX IF NOT EXISTS personas_one_ghost_per_user_idx
  ON personas(user_id)
  WHERE persona_type = 'ghost';
