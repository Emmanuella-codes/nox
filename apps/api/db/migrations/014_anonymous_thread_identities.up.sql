ALTER TABLE comments
  ADD COLUMN IF NOT EXISTS posting_mode TEXT NOT NULL DEFAULT 'public'
  CHECK (posting_mode IN ('public', 'anonymous'));

CREATE TABLE IF NOT EXISTS anonymous_thread_identities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  thread_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  anonymous_handle TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (thread_id, persona_id),
  UNIQUE (thread_id, anonymous_handle)
);

CREATE INDEX IF NOT EXISTS anonymous_thread_identities_persona_id_idx
  ON anonymous_thread_identities(persona_id);
