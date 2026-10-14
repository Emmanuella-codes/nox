package models

import (
	"github.com/google/uuid"
	"time"
)

type GuestQuota struct {
	ID            uuid.UUID `json:"id"`
	EventID       uuid.UUID `json:"event_id"`
	OwnerUserID   uuid.UUID `json:"owner_user_id"`
	Label         string    `json:"label"`
	Capacity      int       `json:"capacity"`
	ReservedCount int       `json:"reserved_count"`
	Status        string    `json:"status"`
}

type GuestEntry struct {
	ID            uuid.UUID  `json:"id"`
	EventID       uuid.UUID  `json:"event_id"`
	QuotaID       uuid.UUID  `json:"quota_id"`
	Name          string     `json:"name"`
	Email         string     `json:"email,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Status        string     `json:"status"`
	MatchedUserID *uuid.UUID `json:"-"`
	InvitationURL string     `json:"invitation_url,omitempty"`
	Credential    string     `json:"credential,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type GuestListReport struct {
	EventID  uuid.UUID      `json:"event_id"`
	Capacity int            `json:"capacity"`
	Reserved int            `json:"reserved"`
	Quotas   []*GuestQuota  `json:"quotas"`
	Statuses map[string]int `json:"statuses"`
}
