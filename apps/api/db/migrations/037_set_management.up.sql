CREATE TABLE set_plays (
  set_id UUID NOT NULL REFERENCES sets(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  play_day DATE NOT NULL DEFAULT CURRENT_DATE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (set_id, user_id, play_day)
);

CREATE INDEX set_plays_user_created_at_idx ON set_plays(user_id, created_at DESC);
