CREATE TABLE event_staff_memberships (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT, invited_by UUID REFERENCES users(id) ON DELETE SET NULL,
  status TEXT NOT NULL DEFAULT 'invited' CHECK (status IN ('invited','active','declined','suspended','revoked')),
  accepted_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(event_id,user_id)
);
CREATE TABLE event_staff_roles (
  membership_id UUID NOT NULL REFERENCES event_staff_memberships(id) ON DELETE CASCADE,
  role TEXT NOT NULL CHECK (role IN ('owner','manager','safety_lead','safety_officer','gate','finance','read_only')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY(membership_id,role)
);
CREATE TABLE event_staff_shifts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), membership_id UUID NOT NULL REFERENCES event_staff_memberships(id) ON DELETE CASCADE,
  starts_at TIMESTAMPTZ NOT NULL, ends_at TIMESTAMPTZ NOT NULL, status TEXT NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled','active','completed','cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), CHECK (ends_at > starts_at)
);
CREATE TABLE event_staff_audit_logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  membership_id UUID REFERENCES event_staff_memberships(id) ON DELETE SET NULL, actor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  action TEXT NOT NULL, reason TEXT, metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX event_staff_memberships_event_idx ON event_staff_memberships(event_id,status);
CREATE INDEX event_staff_shifts_membership_idx ON event_staff_shifts(membership_id,starts_at,ends_at);
CREATE INDEX event_staff_audit_event_idx ON event_staff_audit_logs(event_id,created_at DESC);
CREATE OR REPLACE FUNCTION establish_event_owner() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE membership UUID;
BEGIN
  INSERT INTO event_staff_memberships(event_id,user_id,status) VALUES(NEW.id,NEW.organizer_id,'active') RETURNING id INTO membership;
  INSERT INTO event_staff_roles(membership_id,role) VALUES(membership,'owner');
  RETURN NEW;
END $$;
CREATE TRIGGER events_establish_owner AFTER INSERT ON events FOR EACH ROW EXECUTE FUNCTION establish_event_owner();
INSERT INTO event_staff_memberships(event_id,user_id,status) SELECT id,organizer_id,CASE WHEN status='archived' THEN 'revoked' ELSE 'active' END FROM events ON CONFLICT(event_id,user_id) DO NOTHING;
INSERT INTO event_staff_roles(membership_id,role) SELECT m.id,'owner' FROM event_staff_memberships m WHERE NOT EXISTS (SELECT 1 FROM event_staff_roles r WHERE r.membership_id=m.id AND r.role='owner');
