ALTER TABLE personas
  ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'fan'
  CHECK (category IN ('fan', 'dj', 'organizer', 'creator'));

CREATE INDEX IF NOT EXISTS personas_category_idx ON personas(category);

CREATE TABLE IF NOT EXISTS media_assets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  owner_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  media_kind TEXT NOT NULL CHECK (media_kind IN ('video')),
  storage_key TEXT NOT NULL UNIQUE,
  playback_url TEXT NOT NULL,
  thumbnail_url TEXT,
  mime_type TEXT NOT NULL,
  duration_seconds INT NOT NULL CHECK (duration_seconds > 0),
  size_bytes BIGINT NOT NULL CHECK (size_bytes > 0),
  processing_status TEXT NOT NULL CHECK (processing_status IN ('pending', 'ready', 'failed')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS media_assets_owner_persona_idx ON media_assets(owner_persona_id);
CREATE INDEX IF NOT EXISTS media_assets_processing_status_idx ON media_assets(processing_status);

CREATE TABLE IF NOT EXISTS sets (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  author_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  media_asset_id UUID NOT NULL UNIQUE REFERENCES media_assets(id) ON DELETE RESTRICT,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  genre_tags TEXT[] NOT NULL DEFAULT '{}',
  duration_seconds INT NOT NULL CHECK (duration_seconds > 0 AND duration_seconds <= 900),
  like_count INT NOT NULL DEFAULT 0,
  comment_count INT NOT NULL DEFAULT 0,
  play_count INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS sets_persona_id_idx ON sets(persona_id);
CREATE INDEX IF NOT EXISTS sets_created_at_idx ON sets(created_at DESC);
CREATE INDEX IF NOT EXISTS sets_genre_tags_idx ON sets USING gin (genre_tags);
