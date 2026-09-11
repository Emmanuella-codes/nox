ALTER TABLE media_assets
  ADD COLUMN processing_attempts INT NOT NULL DEFAULT 0,
  ADD COLUMN last_processing_error TEXT NOT NULL DEFAULT '',
  ADD COLUMN last_processing_started_at TIMESTAMPTZ;

CREATE TABLE media_moderation (
  media_asset_id UUID PRIMARY KEY REFERENCES media_assets(id) ON DELETE CASCADE,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'hidden', 'removed')),
  reason TEXT NOT NULL DEFAULT '',
  moderated_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
  moderated_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_media_moderation_status ON media_moderation (status);
CREATE INDEX idx_media_assets_stuck ON media_assets (processing_status, updated_at);
