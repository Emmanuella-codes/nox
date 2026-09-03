package models

import (
	"time"

	"github.com/google/uuid"
)

type DiscoverySuppressionTargetType string

const (
	PersonaSuppressionTargetType DiscoverySuppressionTargetType = "persona"
	PostSuppressionTargetType    DiscoverySuppressionTargetType = "post"
	EventSuppressionTargetType   DiscoverySuppressionTargetType = "event"
	SetSuppressionTargetType     DiscoverySuppressionTargetType = "set"
)

type UserBlock struct {
	BlockerUserID uuid.UUID `json:"blocker_user_id"`
	BlockedUserID uuid.UUID `json:"blocked_user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserMute struct {
	UserID      uuid.UUID `json:"user_id"`
	MutedUserID uuid.UUID `json:"muted_user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type DiscoverySuppression struct {
	UserID     uuid.UUID                      `json:"user_id"`
	TargetType DiscoverySuppressionTargetType `json:"target_type"`
	TargetID   uuid.UUID                      `json:"target_id"`
	CreatedAt  time.Time                      `json:"created_at"`
}
