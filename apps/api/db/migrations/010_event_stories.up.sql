CREATE TABLE IF NOT EXISTS stories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  owner_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  contribution_mode TEXT NOT NULL CHECK (contribution_mode IN ('public', 'followers')),
  total_duration_seconds INT NOT NULL DEFAULT 0 CHECK (total_duration_seconds >= 0 AND total_duration_seconds <= 300),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS story_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  media_asset_id UUID NOT NULL UNIQUE REFERENCES media_assets(id) ON DELETE RESTRICT,
  contributor_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  contributor_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  posting_mode TEXT NOT NULL CHECK (posting_mode IN ('public', 'anonymous')),
  anonymous_label TEXT,
  duration_seconds INT NOT NULL CHECK (duration_seconds > 0),
  position INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (story_id, position)
);

CREATE TABLE IF NOT EXISTS event_highlight_stories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  story_id UUID NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
  added_by_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  position INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (event_id, story_id),
  UNIQUE (event_id, position)
);

CREATE INDEX IF NOT EXISTS stories_event_id_idx ON stories(event_id);
CREATE INDEX IF NOT EXISTS stories_owner_persona_id_idx ON stories(owner_persona_id);
CREATE INDEX IF NOT EXISTS story_items_story_id_idx ON story_items(story_id);
CREATE INDEX IF NOT EXISTS story_items_contributor_persona_id_idx ON story_items(contributor_persona_id);
CREATE INDEX IF NOT EXISTS event_highlight_stories_event_id_idx ON event_highlight_stories(event_id);
