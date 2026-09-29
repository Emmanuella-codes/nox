DROP TRIGGER IF EXISTS events_establish_owner ON events;
DROP FUNCTION IF EXISTS establish_event_owner();
DROP TABLE IF EXISTS event_staff_audit_logs;
DROP TABLE IF EXISTS event_staff_shifts;
DROP TABLE IF EXISTS event_staff_roles;
DROP TABLE IF EXISTS event_staff_memberships;
