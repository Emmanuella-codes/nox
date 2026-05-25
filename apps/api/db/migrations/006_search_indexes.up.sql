CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX IF NOT EXISTS personas_handle_trgm_idx ON personas USING gin (handle gin_trgm_ops);
CREATE INDEX IF NOT EXISTS personas_display_name_trgm_idx ON personas USING gin (display_name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS personas_bio_trgm_idx ON personas USING gin (bio gin_trgm_ops);
CREATE INDEX IF NOT EXISTS personas_genre_tags_idx ON personas USING gin (genre_tags);

CREATE INDEX IF NOT EXISTS posts_body_trgm_idx ON posts USING gin (body gin_trgm_ops);
CREATE INDEX IF NOT EXISTS posts_location_trgm_idx ON posts USING gin (location gin_trgm_ops);

CREATE INDEX IF NOT EXISTS events_title_trgm_idx ON events USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS events_venue_trgm_idx ON events USING gin (venue gin_trgm_ops);
CREATE INDEX IF NOT EXISTS events_location_trgm_idx ON events USING gin (location gin_trgm_ops);
CREATE INDEX IF NOT EXISTS events_description_trgm_idx ON events USING gin (description gin_trgm_ops);
CREATE INDEX IF NOT EXISTS events_genre_tags_idx ON events USING gin (genre_tags);
