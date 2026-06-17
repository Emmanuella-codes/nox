DROP INDEX IF EXISTS sets_genre_tags_idx;
DROP INDEX IF EXISTS sets_created_at_idx;
DROP INDEX IF EXISTS sets_persona_id_idx;
DROP TABLE IF EXISTS sets;

DROP INDEX IF EXISTS media_assets_processing_status_idx;
DROP INDEX IF EXISTS media_assets_owner_persona_idx;
DROP TABLE IF EXISTS media_assets;

DROP INDEX IF EXISTS personas_category_idx;
ALTER TABLE personas DROP COLUMN IF EXISTS category;
