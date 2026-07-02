DROP INDEX IF EXISTS messages_sender_persona_id_idx;
DROP INDEX IF EXISTS messages_conversation_id_created_at_idx;
DROP INDEX IF EXISTS direct_conversations_persona_b_id_idx;
DROP INDEX IF EXISTS conversation_members_persona_id_idx;
DROP INDEX IF EXISTS conversation_members_user_id_idx;

ALTER TABLE conversation_members DROP CONSTRAINT IF EXISTS conversation_members_last_read_message_id_fkey;
ALTER TABLE conversations DROP CONSTRAINT IF EXISTS conversations_last_message_id_fkey;

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS direct_conversations;
DROP TABLE IF EXISTS conversation_members;
DROP TABLE IF EXISTS conversations;
