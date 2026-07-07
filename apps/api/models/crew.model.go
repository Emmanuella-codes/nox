package models

import (
	"time"

	"github.com/google/uuid"
)

type CrewVisibility string
type CrewStatus string
type CrewMemberRole string

const (
	PrivateCrewVisibility    CrewVisibility = "private"
	InviteCodeCrewVisibility CrewVisibility = "invite_code"
)

const (
	ActiveCrewStatus CrewStatus = "active"
	EndedCrewStatus  CrewStatus = "ended"
)

const (
	OwnerCrewMemberRole  CrewMemberRole = "owner"
	MemberCrewMemberRole CrewMemberRole = "member"
)

type EventCrew struct {
	ID             uuid.UUID      `json:"id"`
	EventID        uuid.UUID      `json:"event_id"`
	ConversationID uuid.UUID      `json:"conversation_id"`
	OwnerUserID    uuid.UUID      `json:"-"`
	OwnerPersonaID uuid.UUID      `json:"owner_persona_id"`
	Name           string         `json:"name"`
	JoinCode       string         `json:"join_code"`
	Visibility     CrewVisibility `json:"visibility"`
	Status         CrewStatus     `json:"status"`
	ExpiresAt      time.Time      `json:"expires_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type EventCrewMember struct {
	CrewID                 uuid.UUID      `json:"crew_id"`
	UserID                 uuid.UUID      `json:"-"`
	PersonaID              uuid.UUID      `json:"persona_id"`
	Role                   CrewMemberRole `json:"role"`
	LocationSharingEnabled bool           `json:"location_sharing_enabled"`
	JoinedAt               time.Time      `json:"joined_at"`
	LeftAt                 *time.Time     `json:"left_at,omitempty"`
}

type EventCrewLocation struct {
	CrewID         uuid.UUID `json:"crew_id"`
	UserID         uuid.UUID `json:"-"`
	PersonaID      uuid.UUID `json:"persona_id"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	AccuracyMeters float64   `json:"accuracy_meters"`
	BatteryLevel   *int      `json:"battery_level,omitempty"`
	RecordedAt     time.Time `json:"recorded_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}
