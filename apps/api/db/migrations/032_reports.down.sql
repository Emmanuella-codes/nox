DROP INDEX IF EXISTS idx_moderation_actions_report;
ALTER TABLE moderation_actions DROP COLUMN IF EXISTS report_id;
DROP INDEX IF EXISTS idx_reports_target;
DROP INDEX IF EXISTS idx_reports_assigned;
DROP INDEX IF EXISTS idx_reports_queue;
DROP INDEX IF EXISTS idx_reports_open_reporter_target;
DROP TABLE reports;
