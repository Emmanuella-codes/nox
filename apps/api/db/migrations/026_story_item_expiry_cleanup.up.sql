ALTER TABLE story_items
  ADD COLUMN IF NOT EXISTS expires_at TIMESTAMPTZ;

UPDATE story_items
SET expires_at = created_at + interval '24 hours'
WHERE expires_at IS NULL;

ALTER TABLE story_items
  ALTER COLUMN expires_at SET NOT NULL,
  ALTER COLUMN expires_at SET DEFAULT (now() + interval '24 hours');

CREATE INDEX IF NOT EXISTS story_items_expires_at_idx ON story_items(expires_at);
CREATE INDEX IF NOT EXISTS story_items_story_id_expires_at_idx ON story_items(story_id, expires_at);

UPDATE stories s
SET expires_at = COALESCE((
  SELECT MAX(si.expires_at)
  FROM story_items si
  WHERE si.story_id = s.id
), s.created_at + interval '24 hours');
