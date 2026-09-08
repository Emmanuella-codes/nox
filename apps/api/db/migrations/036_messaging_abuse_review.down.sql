DROP INDEX IF EXISTS message_report_evidence_expiry_idx;
DROP TABLE IF EXISTS message_report_evidence;

ALTER TABLE moderation_actions DROP CONSTRAINT IF EXISTS moderation_actions_entity_type_check;
ALTER TABLE moderation_actions
  ADD CONSTRAINT moderation_actions_entity_type_check CHECK (entity_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event'));

ALTER TABLE reports DROP CONSTRAINT IF EXISTS reports_target_type_check;
ALTER TABLE reports
  ADD CONSTRAINT reports_target_type_check CHECK (target_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event'));
