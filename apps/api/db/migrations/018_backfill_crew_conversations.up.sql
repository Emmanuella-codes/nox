ALTER TABLE event_crews ADD COLUMN IF NOT EXISTS conversation_id UUID;

DO $$
DECLARE
  crew_row RECORD;
  created_conversation_id UUID;
BEGIN
  FOR crew_row IN
    SELECT id, owner_user_id, owner_persona_id, name, conversation_id
    FROM event_crews
  LOOP
    IF crew_row.conversation_id IS NULL THEN
      INSERT INTO conversations (conversation_type, title, created_by)
      VALUES ('group', crew_row.name, crew_row.owner_persona_id)
      RETURNING id INTO created_conversation_id;

      INSERT INTO conversation_members (conversation_id, user_id, persona_id, role)
      VALUES (created_conversation_id, crew_row.owner_user_id, crew_row.owner_persona_id, 'admin')
      ON CONFLICT (conversation_id, persona_id) DO UPDATE SET left_at = NULL;

      UPDATE event_crews
      SET conversation_id = created_conversation_id
      WHERE id = crew_row.id;

      INSERT INTO conversation_members (conversation_id, user_id, persona_id, role)
      SELECT
        created_conversation_id,
        user_id,
        persona_id,
        CASE WHEN role = 'owner' THEN 'admin' ELSE 'member' END
      FROM event_crew_members
      WHERE crew_id = crew_row.id AND left_at IS NULL
      ON CONFLICT (conversation_id, persona_id) DO UPDATE SET left_at = NULL;
    END IF;
  END LOOP;
END $$;

ALTER TABLE event_crews ALTER COLUMN conversation_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'event_crews_conversation_id_key'
  ) THEN
    ALTER TABLE event_crews ADD CONSTRAINT event_crews_conversation_id_key UNIQUE (conversation_id);
  END IF;
END $$;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'event_crews_conversation_id_fkey'
  ) THEN
    ALTER TABLE event_crews
      ADD CONSTRAINT event_crews_conversation_id_fkey
      FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE;
  END IF;
END $$;

CREATE INDEX IF NOT EXISTS event_crews_conversation_id_idx ON event_crews(conversation_id);
