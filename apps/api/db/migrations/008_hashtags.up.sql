CREATE TABLE IF NOT EXISTS hashtags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tag TEXT NOT NULL UNIQUE,
  post_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS post_hashtags (
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  hashtag_id UUID NOT NULL REFERENCES hashtags(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (post_id, hashtag_id)
);

CREATE INDEX IF NOT EXISTS hashtags_post_count_idx ON hashtags(post_count DESC);
CREATE INDEX IF NOT EXISTS post_hashtags_hashtag_id_idx ON post_hashtags(hashtag_id);
