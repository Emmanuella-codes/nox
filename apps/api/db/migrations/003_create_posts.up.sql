CREATE UNIQUE INDEX IF NOT EXISTS personas_one_ghost_per_user_idx
  ON personas(user_id)
  WHERE persona_type = 'ghost';

CREATE TABLE IF NOT EXISTS posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
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
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS posts_persona_id_idx ON posts(persona_id);
CREATE INDEX IF NOT EXISTS posts_created_at_idx ON posts(created_at DESC);
