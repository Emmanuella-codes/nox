DROP INDEX IF EXISTS notifications_story_contribution_request_unique_idx;
DROP INDEX IF EXISTS notifications_event_id_created_at_idx;
DROP INDEX IF EXISTS notifications_story_id_created_at_idx;

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  FOR constraint_name IN
    SELECT con.conname
    FROM pg_constraint con
    WHERE con.conrelid = 'notifications'::regclass
      AND con.contype = 'c'
      AND pg_get_constraintdef(con.oid) LIKE '%notification_type%'
  LOOP
    EXECUTE format('ALTER TABLE notifications DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END LOOP;
END $$;

ALTER TABLE notifications
  ADD CONSTRAINT notifications_notification_type_check CHECK (
    notification_type IN (
      'follow',
      'like',
      'comment',
      'repost',
      'mention',
      'message_direct',
      'message_group'
    )
  );

ALTER TABLE notifications
  DROP COLUMN IF EXISTS story_contribution_request_id,
  DROP COLUMN IF EXISTS story_id,
  DROP COLUMN IF EXISTS event_id;

DROP INDEX IF EXISTS story_contribution_requests_contributor_persona_id_idx;
DROP INDEX IF EXISTS story_contribution_requests_story_id_created_at_idx;
DROP INDEX IF EXISTS story_contribution_requests_pending_media_asset_idx;

DROP TABLE IF EXISTS story_contribution_requests;
