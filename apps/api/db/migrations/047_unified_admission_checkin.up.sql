CREATE TABLE event_attendees (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE SET NULL, display_name TEXT NOT NULL, email TEXT, phone TEXT,
  status TEXT NOT NULL DEFAULT 'eligible' CHECK (status IN ('eligible','restricted','removed')), created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE event_admissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), attendee_id UUID NOT NULL REFERENCES event_attendees(id) ON DELETE CASCADE,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE, source_type TEXT NOT NULL CHECK (source_type IN ('paid_ticket','guest','complimentary','staff')),
  source_id UUID NOT NULL, status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','revoked','expired','refunded','used')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(source_type,source_id)
);
CREATE TABLE admission_credentials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), admission_id UUID NOT NULL UNIQUE REFERENCES event_admissions(id) ON DELETE CASCADE,
  credential_hash TEXT NOT NULL UNIQUE, status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','used','revoked','expired')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), used_at TIMESTAMPTZ
);
CREATE TABLE event_attendee_restrictions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, reason TEXT NOT NULL, status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','cleared')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE event_check_ins (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  admission_id UUID NOT NULL REFERENCES event_admissions(id) ON DELETE CASCADE, staff_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  method TEXT NOT NULL CHECK (method IN ('qr','manual','override')), reason TEXT, invalid_reason TEXT,
  checked_in_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(admission_id)
);
CREATE TABLE event_admission_audits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  admission_id UUID REFERENCES event_admissions(id) ON DELETE SET NULL, attendee_id UUID REFERENCES event_attendees(id) ON DELETE SET NULL,
  actor_user_id UUID REFERENCES users(id) ON DELETE SET NULL, action TEXT NOT NULL, reason TEXT, metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX event_attendees_event_idx ON event_attendees(event_id,status);
CREATE INDEX event_admissions_event_idx ON event_admissions(event_id,status);
CREATE INDEX event_admission_audits_event_idx ON event_admission_audits(event_id,created_at DESC);
