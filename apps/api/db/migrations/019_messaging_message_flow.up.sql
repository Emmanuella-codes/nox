DO $$
BEGIN
  IF EXISTS (
    SELECT 1
    FROM pg_constraint
    WHERE conrelid = 'direct_conversations'::regclass
      AND conname = 'direct_conversations_check'
  ) THEN
    ALTER TABLE direct_conversations
      DROP CONSTRAINT direct_conversations_check;
  END IF;
END $$;

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  SELECT con.conname INTO constraint_name
  FROM pg_constraint con
  WHERE con.conrelid = 'media_assets'::regclass
    AND pg_get_constraintdef(con.oid) LIKE '%media_kind IN (%';

  IF constraint_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE media_assets DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END IF;
END $$;

ALTER TABLE media_assets
  ADD CONSTRAINT media_assets_media_kind_check CHECK (media_kind IN ('image', 'video', 'audio'));

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  SELECT con.conname INTO constraint_name
  FROM pg_constraint con
  WHERE con.conrelid = 'messages'::regclass
    AND pg_get_constraintdef(con.oid) LIKE '%message_type IN (%';

  IF constraint_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE messages DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END IF;
END $$;

ALTER TABLE messages
  ADD CONSTRAINT messages_message_type_check CHECK (message_type IN ('text', 'image', 'video', 'audio', 'system'));

CREATE TABLE IF NOT EXISTS message_attachments (
  message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  media_asset_id UUID NOT NULL UNIQUE REFERENCES media_assets(id) ON DELETE RESTRICT,
  position INT NOT NULL DEFAULT 0 CHECK (position >= 0 AND position < 5),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (message_id, media_asset_id)
);

INSERT INTO message_attachments (message_id, media_asset_id, position)
SELECT id, media_asset_id, 0
FROM messages
WHERE media_asset_id IS NOT NULL
ON CONFLICT (message_id, media_asset_id) DO NOTHING;

CREATE INDEX IF NOT EXISTS message_attachments_message_id_position_idx
  ON message_attachments(message_id, position);
