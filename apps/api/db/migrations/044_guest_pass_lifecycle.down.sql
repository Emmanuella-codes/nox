ALTER TABLE event_guest_entries DROP COLUMN IF EXISTS no_show_at;
ALTER TABLE event_guest_entries DROP COLUMN IF EXISTS declined_at;
ALTER TABLE event_guest_entries DROP COLUMN IF EXISTS revoked_at;
ALTER TABLE event_guest_entries DROP COLUMN IF EXISTS credential_hash;
ALTER TABLE event_guest_entries DROP CONSTRAINT IF EXISTS event_guest_entries_status_check;
ALTER TABLE event_guest_entries ADD CONSTRAINT event_guest_entries_status_check CHECK (status IN ('invited','claimed','revoked','declined','checked_in','expired'));
ALTER TABLE event_guest_entries ALTER COLUMN invitation_hash SET NOT NULL;
