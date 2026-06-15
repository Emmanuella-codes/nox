ALTER TABLE stories
  ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ NOT NULL DEFAULT (now() + interval '24 hours');

CREATE INDEX IF NOT EXISTS stories_expires_at_idx ON stories(expires_at);
