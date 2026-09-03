package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type PersonaTargetDTO struct {
	PersonaID       uuid.UUID `json:"persona_id" validate:"required"`
	TargetPersonaID uuid.UUID `json:"target_persona_id" validate:"required"`
}

type DiscoverySuppressionDTO struct {
	PersonaID  uuid.UUID                             `json:"persona_id" validate:"required"`
	TargetType models.DiscoverySuppressionTargetType `json:"target_type" validate:"required"`
	TargetID   uuid.UUID                             `json:"target_id" validate:"required"`
}
