CREATE TABLE ticket_orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE RESTRICT, reservation_id UUID NOT NULL UNIQUE REFERENCES ticket_reservations(id) ON DELETE RESTRICT,
  subtotal_minor BIGINT NOT NULL CHECK (subtotal_minor >= 0), provider_fee_minor BIGINT NOT NULL CHECK (provider_fee_minor >= 0), platform_fee_minor BIGINT NOT NULL CHECK (platform_fee_minor >= 0), total_minor BIGINT NOT NULL CHECK (total_minor >= 0), currency TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','paid','failed','expired','refund_pending','refunded','refund_failed')),
  provider TEXT NOT NULL DEFAULT 'paystack', provider_reference TEXT NOT NULL UNIQUE, expires_at TIMESTAMPTZ NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE ticket_payments (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), order_id UUID NOT NULL REFERENCES ticket_orders(id) ON DELETE RESTRICT, provider_event_id TEXT NOT NULL UNIQUE, provider_reference TEXT NOT NULL, status TEXT NOT NULL, amount_minor BIGINT NOT NULL, currency TEXT NOT NULL, payload JSONB NOT NULL DEFAULT '{}', processed_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE INDEX ticket_orders_user_status_idx ON ticket_orders(user_id,status,created_at DESC);
CREATE INDEX ticket_payments_reference_idx ON ticket_payments(provider_reference);
CREATE TABLE tickets (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), order_id UUID NOT NULL REFERENCES ticket_orders(id) ON DELETE RESTRICT, event_id UUID NOT NULL REFERENCES events(id) ON DELETE RESTRICT, tier_id UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE RESTRICT, status TEXT NOT NULL DEFAULT 'issued', credential TEXT NOT NULL UNIQUE, created_at TIMESTAMPTZ NOT NULL DEFAULT now());
CREATE TABLE attendee_entitlements (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE RESTRICT, event_id UUID NOT NULL REFERENCES events(id) ON DELETE RESTRICT, user_id UUID REFERENCES users(id) ON DELETE SET NULL, status TEXT NOT NULL DEFAULT 'unassigned', created_at TIMESTAMPTZ NOT NULL DEFAULT now());
