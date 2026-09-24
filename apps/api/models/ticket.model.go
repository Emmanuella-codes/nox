package models

import (
	"github.com/google/uuid"
	"time"
)

type TicketWave struct {
	ID           uuid.UUID     `json:"id"`
	EventID      uuid.UUID     `json:"event_id"`
	Name         string        `json:"name"`
	DisplayOrder int           `json:"display_order"`
	SaleStartAt  *time.Time    `json:"sale_start_at,omitempty"`
	SaleEndAt    *time.Time    `json:"sale_end_at,omitempty"`
	ReleaseMode  string        `json:"release_mode"`
	Status       string        `json:"status"`
	Tiers        []*TicketTier `json:"tiers"`
}
type TicketTier struct {
	ID                   uuid.UUID `json:"id"`
	WaveID               uuid.UUID `json:"wave_id"`
	Name                 string    `json:"name"`
	Description          string    `json:"description"`
	PriceMinor           int64     `json:"price_minor"`
	Currency             string    `json:"currency"`
	Inventory            int       `json:"inventory"`
	RemainingInventory   int       `json:"remaining_inventory"`
	AttendeesPerPurchase int       `json:"attendees_per_purchase"`
	MaxPurchasesPerUser  int       `json:"max_purchases_per_user"`
	Status               string    `json:"status"`
}
type TicketReservation struct {
	ID               uuid.UUID `json:"id"`
	TierID           uuid.UUID `json:"tier_id"`
	UserID           uuid.UUID `json:"-"`
	Quantity         int       `json:"quantity"`
	AttendeeQuantity int       `json:"attendee_quantity"`
	Status           string    `json:"status"`
	ExpiresAt        time.Time `json:"expires_at"`
	OrderID          uuid.UUID `json:"order_id,omitempty"`
}

type TicketOrder struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"-"`
	EventID           uuid.UUID `json:"event_id"`
	ReservationID     uuid.UUID `json:"reservation_id"`
	SubtotalMinor     int64     `json:"subtotal_minor"`
	ProviderFeeMinor  int64     `json:"provider_fee_minor"`
	PlatformFeeMinor  int64     `json:"platform_fee_minor"`
	TotalMinor        int64     `json:"total_minor"`
	Currency          string    `json:"currency"`
	Status            string    `json:"status"`
	ProviderReference string    `json:"provider_reference"`
	AuthorizationURL  string    `json:"authorization_url,omitempty"`
	ExpiresAt         time.Time `json:"expires_at"`
}

type Ticket struct {
	ID         uuid.UUID `json:"id"`
	OrderID    uuid.UUID `json:"order_id"`
	EventID    uuid.UUID `json:"event_id"`
	TierID     uuid.UUID `json:"tier_id"`
	Status     string    `json:"status"`
	Credential string    `json:"credential"`
}
type AttendeeEntitlement struct {
	ID       uuid.UUID  `json:"id"`
	TicketID uuid.UUID  `json:"ticket_id"`
	EventID  uuid.UUID  `json:"event_id"`
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Status   string     `json:"status"`
}
