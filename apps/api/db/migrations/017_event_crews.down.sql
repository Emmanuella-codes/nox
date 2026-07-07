DROP INDEX IF EXISTS event_crew_locations_expires_at_idx;
DROP INDEX IF EXISTS event_crew_members_persona_id_idx;
DROP INDEX IF EXISTS event_crew_members_user_id_idx;
DROP INDEX IF EXISTS event_crews_join_code_idx;
DROP INDEX IF EXISTS event_crews_conversation_id_idx;
DROP INDEX IF EXISTS event_crews_event_id_idx;

DROP TABLE IF EXISTS event_crew_locations;
DROP TABLE IF EXISTS event_crew_members;
DROP TABLE IF EXISTS event_crews;
