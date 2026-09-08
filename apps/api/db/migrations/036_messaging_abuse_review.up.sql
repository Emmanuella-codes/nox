DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  SELECT con.conname INTO constraint_name
  FROM pg_constraint con
  WHERE con.conrelid = 'reports'::regclass
    AND pg_get_constraintdef(con.oid) LIKE '%target_type IN%';
  IF constraint_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE reports DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END IF;
END $$;

ALTER TABLE reports
  ADD CONSTRAINT reports_target_type_check CHECK (target_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event', 'message'));

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  SELECT con.conname INTO constraint_name
  FROM pg_constraint con
  WHERE con.conrelid = 'moderation_actions'::regclass
    AND pg_get_constraintdef(con.oid) LIKE '%entity_type IN%';
  IF constraint_name IS NOT NULL THEN
    EXECUTE format('ALTER TABLE moderation_actions DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END IF;
END $$;

ALTER TABLE moderation_actions
  ADD CONSTRAINT moderation_actions_entity_type_check CHECK (entity_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event', 'message'));

CREATE TABLE message_report_evidence (
  report_id UUID PRIMARY KEY REFERENCES reports(id) ON DELETE CASCADE,
  message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
  conversation_id UUID NOT NULL,
  sender_user_id UUID NOT NULL,
  sender_persona_id UUID NOT NULL,
  body TEXT NOT NULL,
  message_type TEXT NOT NULL,
  attachments JSONB NOT NULL DEFAULT '[]'::jsonb,
  message_created_at TIMESTAMPTZ NOT NULL,
  message_edited_at TIMESTAMPTZ,
  message_deleted_at TIMESTAMPTZ,
  captured_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT now() + INTERVAL '90 days'
);

CREATE INDEX message_report_evidence_expiry_idx ON message_report_evidence (expires_at);
