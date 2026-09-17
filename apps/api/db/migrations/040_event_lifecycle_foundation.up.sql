ALTER TABLE events
  ADD COLUMN IF NOT EXISTS event_end_date TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'pre_event', 'live', 'closing', 'ended', 'follow_up', 'closed', 'archived')),
  ADD COLUMN IF NOT EXISTS ticket_sale_status TEXT NOT NULL DEFAULT 'not_released' CHECK (ticket_sale_status IN ('not_released', 'on_sale', 'sold_out', 'closed')),
  ADD COLUMN IF NOT EXISTS event_mode_status TEXT NOT NULL DEFAULT 'not_started' CHECK (event_mode_status IN ('not_started', 'pre_event', 'live', 'closing', 'ended')),
  ADD COLUMN IF NOT EXISTS published_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS ended_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS closed_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS archived_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

UPDATE events SET event_end_date = event_date + INTERVAL '6 hours' WHERE event_end_date IS NULL;
ALTER TABLE events ALTER COLUMN event_end_date SET NOT NULL;

CREATE TABLE event_lifecycle_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  previous_status TEXT NOT NULL,
  next_status TEXT NOT NULL,
  actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  reason TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX event_lifecycle_history_event_idx ON event_lifecycle_history(event_id, created_at DESC);
CREATE INDEX events_lifecycle_due_idx ON events(status, event_date, event_end_date);
