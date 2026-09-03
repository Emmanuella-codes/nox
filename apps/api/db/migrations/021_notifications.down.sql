DROP INDEX IF EXISTS notifications_follow_unique_idx;
DROP INDEX IF EXISTS notifications_like_unique_idx;
DROP INDEX IF EXISTS notifications_comment_unique_idx;
DROP INDEX IF EXISTS notifications_message_unique_idx;
DROP INDEX IF EXISTS notifications_recipient_persona_id_created_at_idx;
DROP INDEX IF EXISTS notifications_recipient_user_id_created_at_idx;
DROP TABLE IF EXISTS notifications;
