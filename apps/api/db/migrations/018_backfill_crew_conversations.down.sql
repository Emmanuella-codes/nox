DROP INDEX IF EXISTS event_crews_conversation_id_idx;

ALTER TABLE event_crews DROP CONSTRAINT IF EXISTS event_crews_conversation_id_fkey;
ALTER TABLE event_crews DROP CONSTRAINT IF EXISTS event_crews_conversation_id_key;
ALTER TABLE event_crews DROP COLUMN IF EXISTS conversation_id;
