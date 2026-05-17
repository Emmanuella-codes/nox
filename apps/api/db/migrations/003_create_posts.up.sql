CREATE UNIQUE INDEX IF NOT EXISTS personas_one_ghost_per_user_idx
  ON personas(user_id)
  WHERE persona_type = 'ghost';

CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  author_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID REFERENCES personas(id) ON DELETE SET NULL,
  posting_mode TEXT NOT NULL CHECK (posting_mode IN ('public', 'anonymous')),
  event_id UUID,
  body TEXT NOT NULL,
  post_type TEXT NOT NULL CHECK (post_type IN ('text', 'image', 'set', 'event_tag')),
  media_url TEXT,
  media_type TEXT CHECK (media_type IN ('image', 'youtube', 'soundcloud')),
  location TEXT,
  like_count INT NOT NULL DEFAULT 0,
  comment_count INT NOT NULL DEFAULT 0,
  repost_count INT NOT NULL DEFAULT 0,
  is_repost BOOLEAN NOT NULL DEFAULT false,
  repost_of UUID REFERENCES posts(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT posts_posting_mode_persona_check CHECK (
    (posting_mode = 'public' AND persona_id IS NOT NULL)
    OR (posting_mode = 'anonymous' AND persona_id IS NULL)
  )
);

CREATE INDEX IF NOT EXISTS posts_author_user_id_idx ON posts(author_user_id);
CREATE INDEX IF NOT EXISTS posts_persona_id_idx ON posts(persona_id);
CREATE INDEX IF NOT EXISTS posts_created_at_idx ON posts(created_at DESC);
