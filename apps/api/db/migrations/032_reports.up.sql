CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  reporter_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  reporter_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  target_type TEXT NOT NULL CHECK (target_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event')),
  target_id UUID NOT NULL,
  reason TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'reviewing', 'resolved', 'dismissed')),
  assigned_admin_id UUID REFERENCES users(id) ON DELETE SET NULL,
  resolution_note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_reports_open_reporter_target
  ON reports (reporter_user_id, target_type, target_id)
  WHERE status IN ('open', 'reviewing');
CREATE INDEX idx_reports_queue ON reports (status, created_at ASC);
CREATE INDEX idx_reports_assigned ON reports (assigned_admin_id, status, created_at DESC);
CREATE INDEX idx_reports_target ON reports (target_type, target_id, created_at DESC);

ALTER TABLE moderation_actions
ADD COLUMN report_id UUID REFERENCES reports(id) ON DELETE SET NULL;

CREATE INDEX idx_moderation_actions_report ON moderation_actions (report_id, created_at DESC);
