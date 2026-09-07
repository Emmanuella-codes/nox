CREATE TABLE set_features (
  set_id UUID PRIMARY KEY REFERENCES sets(id) ON DELETE CASCADE,
  position INT NOT NULL,
  featured_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ,
  UNIQUE (position)
);

CREATE INDEX idx_set_features_created_at ON set_features(created_at DESC);
