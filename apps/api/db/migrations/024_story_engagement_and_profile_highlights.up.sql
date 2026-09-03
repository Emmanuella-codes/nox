CREATE TABLE IF NOT EXISTS story_item_views (
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  story_item_id UUID NOT NULL REFERENCES story_items(id) ON DELETE CASCADE,
  viewer_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  viewer_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (story_item_id, viewer_persona_id)
);

CREATE INDEX IF NOT EXISTS story_item_views_story_id_idx
  ON story_item_views(story_id, created_at DESC);

CREATE TABLE IF NOT EXISTS story_item_reactions (
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  story_item_id UUID NOT NULL REFERENCES story_items(id) ON DELETE CASCADE,
  reactor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reactor_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  reaction_type TEXT NOT NULL CHECK (reaction_type IN ('like', 'fire', 'heart', 'laugh')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (story_item_id, reactor_persona_id)
);

CREATE INDEX IF NOT EXISTS story_item_reactions_story_id_idx
  ON story_item_reactions(story_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS profile_story_highlights (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  position INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (owner_persona_id, story_id),
  UNIQUE (owner_persona_id, position)
);

CREATE INDEX IF NOT EXISTS profile_story_highlights_owner_persona_id_idx
  ON profile_story_highlights(owner_persona_id, position);

ALTER TABLE messages
  ADD COLUMN IF NOT EXISTS story_id UUID REFERENCES stories(id) ON DELETE SET NULL,
  ADD COLUMN IF NOT EXISTS story_item_id UUID REFERENCES story_items(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS messages_story_id_created_at_idx
  ON messages(story_id, created_at DESC, id DESC);

ALTER TABLE notifications
  ADD COLUMN IF NOT EXISTS story_item_id UUID REFERENCES story_items(id) ON DELETE CASCADE;

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
      'story_highlight_removed',
      'story_reaction'
    )
  );

CREATE UNIQUE INDEX IF NOT EXISTS notifications_story_reaction_unique_idx
  ON notifications(recipient_persona_id, actor_persona_id, story_item_id, notification_type)
  WHERE actor_persona_id IS NOT NULL AND story_item_id IS NOT NULL AND notification_type = 'story_reaction';
