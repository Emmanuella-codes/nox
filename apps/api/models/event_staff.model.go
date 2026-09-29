package models

import (
	"github.com/google/uuid"
	"time"
)

type EventStaffMembership struct {
	ID         uuid.UUID  `json:"id"`
	EventID    uuid.UUID  `json:"event_id"`
	UserID     uuid.UUID  `json:"user_id"`
	Status     string     `json:"status"`
	Roles      []string   `json:"roles"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
}
type EventStaffShift struct {
	ID           uuid.UUID `json:"id"`
	MembershipID uuid.UUID `json:"membership_id"`
	StartsAt     time.Time `json:"starts_at"`
	EndsAt       time.Time `json:"ends_at"`
	Status       string    `json:"status"`
}
