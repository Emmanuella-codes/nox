CREATE TABLE IF NOT EXISTS conversations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_type TEXT NOT NULL CHECK (conversation_type IN ('direct', 'group')),
  title TEXT NOT NULL DEFAULT '',
  created_by UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  last_message_id UUID,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS conversation_members (
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('member', 'admin')),
  last_read_message_id UUID,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  left_at TIMESTAMPTZ,
  PRIMARY KEY (conversation_id, persona_id)
);

CREATE TABLE IF NOT EXISTS direct_conversations (
  conversation_id UUID PRIMARY KEY REFERENCES conversations(id) ON DELETE CASCADE,
  persona_a_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  persona_b_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  CHECK (persona_a_id <> persona_b_id),
  UNIQUE (persona_a_id, persona_b_id)
);

CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  sender_persona_id UUID NOT NULL REFERENCES personas(id) ON DELETE CASCADE,
  body TEXT NOT NULL DEFAULT '' CHECK (char_length(body) <= 2000),
  message_type TEXT NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'video', 'system')),
  media_asset_id UUID REFERENCES media_assets(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  edited_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ,
  CHECK (deleted_at IS NOT NULL OR body <> '' OR media_asset_id IS NOT NULL)
);

ALTER TABLE conversations
  ADD CONSTRAINT conversations_last_message_id_fkey
  FOREIGN KEY (last_message_id) REFERENCES messages(id) ON DELETE SET NULL;

ALTER TABLE conversation_members
  ADD CONSTRAINT conversation_members_last_read_message_id_fkey
  FOREIGN KEY (last_read_message_id) REFERENCES messages(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS conversation_members_user_id_idx ON conversation_members(user_id) WHERE left_at IS NULL;
CREATE INDEX IF NOT EXISTS conversation_members_persona_id_idx ON conversation_members(persona_id) WHERE left_at IS NULL;
CREATE INDEX IF NOT EXISTS direct_conversations_persona_b_id_idx ON direct_conversations(persona_b_id);
CREATE INDEX IF NOT EXISTS messages_conversation_id_created_at_idx ON messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS messages_sender_persona_id_idx ON messages(sender_persona_id);
