CREATE TABLE IF NOT EXISTS persona_follows (
  follower_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  following_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (follower_id, following_id),
  CHECK (follower_id <> following_id)
);

CREATE INDEX IF NOT EXISTS persona_follows_following_id_idx ON persona_follows(following_id);
CREATE INDEX IF NOT EXISTS persona_follows_follower_id_idx ON persona_follows(follower_id);
