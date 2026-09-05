ALTER TABLE personas
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE posts
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE comments
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE stories
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE story_items
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE sets
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE events
ADD COLUMN moderation_status TEXT NOT NULL DEFAULT 'active'
CHECK (moderation_status IN ('active', 'hidden', 'removed')),
ADD COLUMN moderation_reason TEXT NOT NULL DEFAULT '',
ADD COLUMN moderated_at TIMESTAMPTZ,
ADD COLUMN moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL;

CREATE TABLE moderation_actions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  entity_type TEXT NOT NULL CHECK (entity_type IN ('persona', 'post', 'comment', 'story', 'story_item', 'set', 'event')),
  entity_id UUID NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('active', 'hidden', 'removed')),
  reason TEXT NOT NULL DEFAULT '',
  admin_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_moderation_actions_entity_created_at
  ON moderation_actions (entity_type, entity_id, created_at DESC);

CREATE INDEX idx_moderation_actions_admin_created_at
  ON moderation_actions (admin_user_id, created_at DESC);
