ALTER TABLE set_comments
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS set_comment_likes (
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  comment_id UUID NOT NULL REFERENCES set_comments(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (persona_id, comment_id)
);

CREATE INDEX IF NOT EXISTS set_comment_likes_comment_id_idx ON set_comment_likes(comment_id);
