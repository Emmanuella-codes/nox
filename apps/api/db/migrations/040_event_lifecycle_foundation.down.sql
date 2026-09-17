DROP TABLE IF EXISTS event_lifecycle_history;
ALTER TABLE events
  DROP COLUMN IF EXISTS updated_at,
  DROP COLUMN IF EXISTS archived_at,
  DROP COLUMN IF EXISTS closed_at,
  DROP COLUMN IF EXISTS ended_at,
  DROP COLUMN IF EXISTS published_at,
  DROP COLUMN IF EXISTS event_mode_status,
  DROP COLUMN IF EXISTS ticket_sale_status,
  DROP COLUMN IF EXISTS status,
  DROP COLUMN IF EXISTS event_end_date;
