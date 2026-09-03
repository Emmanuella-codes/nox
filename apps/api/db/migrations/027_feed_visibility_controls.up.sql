CREATE TABLE IF NOT EXISTS user_blocks (
  blocker_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  blocked_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (blocker_user_id, blocked_user_id),
  CHECK (blocker_user_id <> blocked_user_id)
);

CREATE INDEX IF NOT EXISTS user_blocks_blocked_user_id_idx ON user_blocks(blocked_user_id);

CREATE TABLE IF NOT EXISTS user_mutes (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  muted_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, muted_user_id),
  CHECK (user_id <> muted_user_id)
);

CREATE INDEX IF NOT EXISTS user_mutes_muted_user_id_idx ON user_mutes(muted_user_id);

CREATE TABLE IF NOT EXISTS discovery_suppressions (
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_type TEXT NOT NULL,
  target_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, target_type, target_id),
  CHECK (target_type IN ('persona', 'post', 'event', 'set'))
);

CREATE INDEX IF NOT EXISTS discovery_suppressions_lookup_idx
  ON discovery_suppressions(user_id, target_type, created_at DESC);
