DROP INDEX IF EXISTS notifications_story_reaction_unique_idx;

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
      'message_group',
      'story_contribution_request',
      'story_contribution_accepted',
      'story_contribution_rejected',
      'story_highlight_added',
      'story_highlight_removed'
    )
  );

ALTER TABLE notifications
  DROP COLUMN IF EXISTS story_item_id;

DROP INDEX IF EXISTS messages_story_id_created_at_idx;

ALTER TABLE messages
  DROP COLUMN IF EXISTS story_item_id,
  DROP COLUMN IF EXISTS story_id;

DROP INDEX IF EXISTS profile_story_highlights_owner_persona_id_idx;
DROP TABLE IF EXISTS profile_story_highlights;

DROP INDEX IF EXISTS story_item_reactions_story_id_idx;
DROP TABLE IF EXISTS story_item_reactions;

DROP INDEX IF EXISTS story_item_views_story_id_idx;
DROP TABLE IF EXISTS story_item_views;
