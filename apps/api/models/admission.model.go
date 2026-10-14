package models

import (
	"time"

	"github.com/google/uuid"
)

type Admission struct {
	ID          uuid.UUID `json:"id"`
	AttendeeID  uuid.UUID `json:"attendee_id"`
	EventID     uuid.UUID `json:"event_id"`
	SourceType  string    `json:"source_type"`
	SourceID    uuid.UUID `json:"source_id"`
	Status      string    `json:"status"`
	DisplayName string    `json:"display_name"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Credential  string    `json:"credential,omitempty"`
}
type AdmissionSearchResult struct {
	AdmissionID uuid.UUID `json:"admission_id"`
	DisplayName string    `json:"display_name"`
	SourceType  string    `json:"source_type"`
	Status      string    `json:"status"`
	Email       string    `json:"email,omitempty"`
	Phone       string    `json:"phone,omitempty"`
}
type CheckInResult struct {
	ID            uuid.UUID `json:"id"`
	AdmissionID   uuid.UUID `json:"admission_id"`
	Method        string    `json:"method"`
	Status        string    `json:"status"`
	InvalidReason string    `json:"invalid_reason,omitempty"`
	CheckedInAt   time.Time `json:"checked_in_at"`
}
