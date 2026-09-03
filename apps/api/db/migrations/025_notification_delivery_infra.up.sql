CREATE TABLE IF NOT EXISTS notification_devices (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  install_id TEXT NOT NULL,
  platform TEXT NOT NULL CHECK (platform IN ('ios', 'android', 'web')),
  push_token TEXT NOT NULL,
  app_version TEXT NOT NULL DEFAULT '',
  last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  disabled_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, install_id),
  UNIQUE (push_token)
);

CREATE INDEX IF NOT EXISTS notification_devices_user_id_idx
  ON notification_devices(user_id, last_seen_at DESC);

CREATE TABLE IF NOT EXISTS notification_preferences (
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  notification_type TEXT NOT NULL,
  push_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (persona_id, notification_type)
);

CREATE TABLE IF NOT EXISTS notification_outbox (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
  recipient_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  recipient_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  channel TEXT NOT NULL CHECK (channel IN ('push')),
  status TEXT NOT NULL CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'dead', 'skipped')),
  payload JSONB NOT NULL,
  attempt_count INT NOT NULL DEFAULT 0,
  next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_error TEXT NOT NULL DEFAULT '',
  worker_id TEXT,
  claimed_at TIMESTAMPTZ,
  sent_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (notification_id, channel)
);

CREATE INDEX IF NOT EXISTS notification_outbox_status_next_attempt_idx
  ON notification_outbox(status, next_attempt_at ASC, created_at ASC);
