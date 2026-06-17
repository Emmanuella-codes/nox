DROP TABLE IF EXISTS post_media_assets;

UPDATE media_assets
SET media_kind = 'video'
WHERE media_kind <> 'video';

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
  ADD CONSTRAINT media_assets_media_kind_check CHECK (media_kind IN ('video'));
