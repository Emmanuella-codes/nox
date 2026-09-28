ALTER TABLE event_guest_entries ALTER COLUMN invitation_hash DROP NOT NULL;
ALTER TABLE event_guest_entries DROP CONSTRAINT IF EXISTS event_guest_entries_status_check;
ALTER TABLE event_guest_entries ADD CONSTRAINT event_guest_entries_status_check CHECK (status IN ('invited','claimed','revoked','declined','checked_in','expired','no_show'));
ALTER TABLE event_guest_entries ADD COLUMN IF NOT EXISTS credential_hash TEXT UNIQUE;
ALTER TABLE event_guest_entries ADD COLUMN IF NOT EXISTS revoked_at TIMESTAMPTZ;
ALTER TABLE event_guest_entries ADD COLUMN IF NOT EXISTS declined_at TIMESTAMPTZ;
ALTER TABLE event_guest_entries ADD COLUMN IF NOT EXISTS no_show_at TIMESTAMPTZ;
