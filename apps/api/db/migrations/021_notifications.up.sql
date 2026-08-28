CREATE TABLE IF NOT EXISTS notifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  recipient_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recipient_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  actor_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  conversation_id UUID REFERENCES conversations(id) ON DELETE CASCADE,
  message_id UUID REFERENCES messages(id) ON DELETE CASCADE,
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
  is_read BOOLEAN NOT NULL DEFAULT FALSE,
  read_at TIMESTAMPTZ,
  notification_type TEXT NOT NULL CHECK (notification_type IN ('follow', 'like', 'comment', 'repost', 'mention', 'message_direct', 'message_group')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (recipient_persona_id, message_id, notification_type)
);

CREATE INDEX IF NOT EXISTS notifications_recipient_user_id_created_at_idx
  ON notifications(recipient_user_id, created_at DESC, id DESC);

CREATE INDEX IF NOT EXISTS notifications_recipient_persona_id_created_at_idx
  ON notifications(recipient_persona_id, created_at DESC, id DESC);
