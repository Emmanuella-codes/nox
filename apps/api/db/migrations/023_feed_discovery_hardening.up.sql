CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS posts_created_at_id_idx
  ON posts(created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS posts_public_persona_created_at_id_idx
  ON posts(persona_id, created_at DESC, id DESC)
  WHERE posting_mode = 'public';

CREATE INDEX IF NOT EXISTS posts_author_created_at_id_idx
  ON posts(author_user_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS sets_title_trgm_idx
  ON sets USING gin (title gin_trgm_ops);

CREATE INDEX IF NOT EXISTS sets_description_trgm_idx
  ON sets USING gin (description gin_trgm_ops);
