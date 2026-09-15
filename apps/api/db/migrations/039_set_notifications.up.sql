ALTER TABLE notifications
  ADD COLUMN IF NOT EXISTS set_id UUID REFERENCES sets(id) ON DELETE CASCADE,
  ADD COLUMN IF NOT EXISTS set_comment_id UUID REFERENCES set_comments(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS notifications_set_id_idx ON notifications(set_id);
