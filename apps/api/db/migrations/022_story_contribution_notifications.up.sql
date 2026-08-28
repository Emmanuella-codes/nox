CREATE TABLE IF NOT EXISTS story_contribution_requests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  media_asset_id UUID NOT NULL REFERENCES media_assets(id) ON DELETE RESTRICT,
  contributor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  contributor_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('pending', 'accepted', 'rejected')),
  reviewed_by_persona_id UUID REFERENCES personas(id) ON DELETE SET NULL,
  story_item_id UUID REFERENCES story_items(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  reviewed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS story_contribution_requests_pending_media_asset_idx
  ON story_contribution_requests(media_asset_id)
  WHERE status = 'pending';

CREATE INDEX IF NOT EXISTS story_contribution_requests_story_id_created_at_idx
  ON story_contribution_requests(story_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS story_contribution_requests_contributor_persona_id_idx
  ON story_contribution_requests(contributor_persona_id);

ALTER TABLE notifications
  ADD COLUMN IF NOT EXISTS event_id UUID REFERENCES events(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS story_id UUID REFERENCES stories(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS story_contribution_request_id UUID REFERENCES story_contribution_requests(id) ON DELETE CASCADE;

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

CREATE INDEX IF NOT EXISTS notifications_story_id_created_at_idx
  ON notifications(story_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS notifications_event_id_created_at_idx
  ON notifications(event_id, created_at DESC, id DESC);

CREATE UNIQUE INDEX IF NOT EXISTS notifications_story_contribution_request_unique_idx
  ON notifications(recipient_persona_id, story_contribution_request_id, notification_type)
  WHERE story_contribution_request_id IS NOT NULL;
