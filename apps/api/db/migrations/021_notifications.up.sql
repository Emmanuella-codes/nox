CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipient_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recipient_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  actor_persona_id UUID REFERENCES personas(id) ON DELETE CASCADE,
  actor_posting_mode TEXT NOT NULL DEFAULT 'public' CHECK (actor_posting_mode IN ('public', 'anonymous')),
  actor_anonymous_handle TEXT NOT NULL DEFAULT '',
  actor_anonymous_avatar_key TEXT NOT NULL DEFAULT '',
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  read_at TIMESTAMPTZ,
  notification_type TEXT NOT NULL CHECK (notification_type IN ('follow', 'like', 'comment', 'repost', 'mention', 'message_direct', 'message_group')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS notifications_recipient_user_id_created_at_idx
  ON notifications(recipient_user_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS notifications_recipient_persona_id_created_at_idx
  ON notifications(recipient_persona_id, created_at DESC, id DESC);

CREATE UNIQUE INDEX IF NOT EXISTS notifications_message_unique_idx
  ON notifications(recipient_persona_id, message_id, notification_type)
  WHERE message_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS notifications_comment_unique_idx
  ON notifications(recipient_persona_id, comment_id, notification_type)
  WHERE comment_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS notifications_like_unique_idx
  ON notifications(recipient_persona_id, actor_persona_id, post_id, notification_type)
  WHERE message_id IS NULL AND comment_id IS NULL AND post_id IS NOT NULL AND actor_persona_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS notifications_follow_unique_idx
  ON notifications(recipient_persona_id, actor_persona_id, notification_type)
  WHERE message_id IS NULL AND comment_id IS NULL AND post_id IS NULL AND conversation_id IS NULL AND actor_persona_id IS NOT NULL;
