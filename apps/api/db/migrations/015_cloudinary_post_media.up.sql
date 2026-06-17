DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  FOR constraint_name IN
    SELECT con.conname
    FROM pg_constraint con
    WHERE con.conrelid = 'media_assets'::regclass
      AND con.contype = 'c'
      AND pg_get_constraintdef(con.oid) LIKE '%media_kind%'
  LOOP
    EXECUTE format('ALTER TABLE media_assets DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END LOOP;
END $$;

ALTER TABLE media_assets
  ADD CONSTRAINT media_assets_media_kind_check CHECK (media_kind IN ('image', 'video'));

CREATE TABLE IF NOT EXISTS post_media_assets (
  post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  media_asset_id UUID NOT NULL UNIQUE REFERENCES media_assets(id) ON DELETE RESTRICT,
  position INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (post_id, media_asset_id),
  UNIQUE (post_id, position)
);

CREATE INDEX IF NOT EXISTS post_media_assets_post_id_idx ON post_media_assets(post_id);
