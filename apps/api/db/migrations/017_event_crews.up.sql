CREATE TABLE IF NOT EXISTS event_crews (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  conversation_id UUID NOT NULL UNIQUE REFERENCES conversations(id) ON DELETE CASCADE,
  owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  owner_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  name TEXT NOT NULL CHECK (char_length(name) BETWEEN 1 AND 80),
  join_code TEXT NOT NULL UNIQUE CHECK (char_length(join_code) = 6),
  visibility TEXT NOT NULL DEFAULT 'invite_code' CHECK (visibility IN ('private', 'invite_code')),
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'ended')),
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS event_crew_members (
  crew_id UUID NOT NULL REFERENCES event_crews(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('owner', 'member')),
  location_sharing_enabled BOOLEAN NOT NULL DEFAULT false,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  left_at TIMESTAMPTZ,
  PRIMARY KEY (crew_id, persona_id)
);

CREATE TABLE IF NOT EXISTS event_crew_locations (
  crew_id UUID NOT NULL REFERENCES event_crews(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  latitude DOUBLE PRECISION NOT NULL CHECK (latitude BETWEEN -90 AND 90),
  longitude DOUBLE PRECISION NOT NULL CHECK (longitude BETWEEN -180 AND 180),
  accuracy_meters DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (accuracy_meters >= 0),
  battery_level INT CHECK (battery_level BETWEEN 0 AND 100),
  recorded_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL,
  PRIMARY KEY (crew_id, persona_id)
);

CREATE INDEX IF NOT EXISTS event_crews_event_id_idx ON event_crews(event_id);
CREATE INDEX IF NOT EXISTS event_crews_conversation_id_idx ON event_crews(conversation_id);
CREATE INDEX IF NOT EXISTS event_crews_join_code_idx ON event_crews(join_code);
CREATE INDEX IF NOT EXISTS event_crew_members_user_id_idx ON event_crew_members(user_id) WHERE left_at IS NULL;
CREATE INDEX IF NOT EXISTS event_crew_members_persona_id_idx ON event_crew_members(persona_id) WHERE left_at IS NULL;
CREATE INDEX IF NOT EXISTS event_crew_locations_expires_at_idx ON event_crew_locations(expires_at);
