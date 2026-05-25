DROP INDEX IF EXISTS events_genre_tags_idx;
DROP INDEX IF EXISTS events_description_trgm_idx;
DROP INDEX IF EXISTS events_location_trgm_idx;
DROP INDEX IF EXISTS events_venue_trgm_idx;
DROP INDEX IF EXISTS events_title_trgm_idx;

DROP INDEX IF EXISTS posts_location_trgm_idx;
DROP INDEX IF EXISTS posts_body_trgm_idx;

DROP INDEX IF EXISTS personas_genre_tags_idx;
DROP INDEX IF EXISTS personas_bio_trgm_idx;
DROP INDEX IF EXISTS personas_display_name_trgm_idx;
DROP INDEX IF EXISTS personas_handle_trgm_idx;
