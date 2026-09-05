ALTER TABLE moderation_actions
ADD COLUMN previous_status TEXT CHECK (previous_status IN ('active', 'hidden', 'removed'));
