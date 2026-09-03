DROP INDEX IF EXISTS message_attachments_message_id_position_idx;
DROP TABLE IF EXISTS message_attachments;

ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_message_type_check;
ALTER TABLE messages
  ADD CONSTRAINT messages_message_type_check CHECK (message_type IN ('text', 'image', 'video', 'system'));

ALTER TABLE media_assets DROP CONSTRAINT IF EXISTS media_assets_media_kind_check;
ALTER TABLE media_assets
  ADD CONSTRAINT media_assets_media_kind_check CHECK (media_kind IN ('image', 'video'));

ALTER TABLE direct_conversations
  ADD CONSTRAINT direct_conversations_check CHECK (persona_a_id <> persona_b_id);
