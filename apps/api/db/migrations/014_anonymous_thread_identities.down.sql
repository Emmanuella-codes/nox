DROP TABLE IF EXISTS anonymous_thread_identities;

ALTER TABLE comments
  DROP COLUMN IF EXISTS posting_mode;
