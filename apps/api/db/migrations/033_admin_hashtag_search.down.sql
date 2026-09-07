DROP INDEX IF EXISTS idx_search_suppressions_expires_at;
DROP INDEX IF EXISTS idx_hashtag_suppressions_expires_at;
DROP INDEX IF EXISTS idx_hashtag_moderation_status;
DROP TABLE IF EXISTS search_suppressions;
DROP TABLE IF EXISTS hashtag_suppressions;
DROP TABLE IF EXISTS hashtag_moderation;
