CREATE TABLE event_guest_settings (
  event_id UUID PRIMARY KEY REFERENCES events(id) ON DELETE CASCADE,
  capacity INT NOT NULL CHECK (capacity > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE event_guest_quotas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, label TEXT NOT NULL,
  capacity INT NOT NULL CHECK (capacity > 0), reserved_count INT NOT NULL DEFAULT 0 CHECK (reserved_count >= 0),
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','paused','closed')),
  UNIQUE(event_id, owner_user_id, label)
);
CREATE TABLE event_guest_entries (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  quota_id UUID NOT NULL REFERENCES event_guest_quotas(id) ON DELETE RESTRICT, added_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  name TEXT NOT NULL, email TEXT, phone TEXT, matched_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  status TEXT NOT NULL DEFAULT 'invited' CHECK (status IN ('invited','claimed','revoked','declined','checked_in','expired')),
  invitation_hash TEXT NOT NULL UNIQUE, invitation_expires_at TIMESTAMPTZ NOT NULL,
  claimed_at TIMESTAMPTZ, checked_in_at TIMESTAMPTZ, checked_in_by UUID REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (email IS NOT NULL OR phone IS NOT NULL)
);
CREATE INDEX event_guest_entries_event_status_idx ON event_guest_entries(event_id,status);
CREATE INDEX event_guest_entries_invitation_hash_idx ON event_guest_entries(invitation_hash);
