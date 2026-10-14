CREATE TABLE ticket_waves (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name TEXT NOT NULL, display_order INT NOT NULL DEFAULT 0, sale_start_at TIMESTAMPTZ, sale_end_at TIMESTAMPTZ,
  release_mode TEXT NOT NULL DEFAULT 'scheduled' CHECK (release_mode IN ('scheduled','after_previous_sold_out','manual')),
  status TEXT NOT NULL DEFAULT 'not_released' CHECK (status IN ('not_released','scheduled','on_sale','sold_out','closed','cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE ticket_tiers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), wave_id UUID NOT NULL REFERENCES ticket_waves(id) ON DELETE CASCADE,
  name TEXT NOT NULL, description TEXT NOT NULL DEFAULT '', price_minor BIGINT NOT NULL CHECK (price_minor >= 0), currency TEXT NOT NULL DEFAULT 'NGN',
  inventory INT NOT NULL CHECK (inventory >= 0), remaining_inventory INT NOT NULL CHECK (remaining_inventory >= 0), attendees_per_purchase INT NOT NULL CHECK (attendees_per_purchase > 0),
  max_purchases_per_user INT NOT NULL CHECK (max_purchases_per_user > 0), status TEXT NOT NULL DEFAULT 'not_released' CHECK (status IN ('not_released','on_sale','sold_out','closed','cancelled')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE ticket_reservations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), tier_id UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE RESTRICT, user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  quantity INT NOT NULL CHECK (quantity > 0), attendee_quantity INT NOT NULL CHECK (attendee_quantity > 0), status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active','converted','expired','cancelled')),
  idempotency_key TEXT NOT NULL, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(user_id, idempotency_key)
);
CREATE INDEX ticket_waves_event_order_idx ON ticket_waves(event_id, display_order);
CREATE INDEX ticket_reservations_expiry_idx ON ticket_reservations(status, expires_at);
