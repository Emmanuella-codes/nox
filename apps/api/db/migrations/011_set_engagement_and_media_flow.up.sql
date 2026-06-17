CREATE TABLE IF NOT EXISTS set_likes (
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  set_id UUID NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (persona_id, set_id)
);

CREATE INDEX IF NOT EXISTS set_likes_set_id_idx ON set_likes(set_id);

CREATE TABLE IF NOT EXISTS set_comments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  set_id UUID NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
  body TEXT NOT NULL CHECK (char_length(body) <= 280),
  parent_id UUID REFERENCES set_comments(id) ON DELETE CASCADE,
  like_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS set_comments_set_id_idx ON set_comments(set_id);
CREATE INDEX IF NOT EXISTS set_comments_parent_id_idx ON set_comments(parent_id);
