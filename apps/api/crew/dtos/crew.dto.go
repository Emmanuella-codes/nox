package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateCrewDTO struct {
	OwnerPersonaID uuid.UUID             `json:"owner_persona_id" validate:"required"`
	Name           string                `json:"name" validate:"required,max=80"`
	Visibility     models.CrewVisibility `json:"visibility"`
}

type JoinCrewDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	JoinCode  string    `json:"join_code" validate:"required,len=6"`
}

type UpdateLocationDTO struct {
	PersonaID      uuid.UUID `json:"persona_id" validate:"required"`
	Latitude       float64   `json:"latitude" validate:"required"`
	Longitude      float64   `json:"longitude" validate:"required"`
	AccuracyMeters float64   `json:"accuracy_meters"`
	BatteryLevel   *int      `json:"battery_level"`
}

type UpdateSharingDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
	Enabled   bool      `json:"enabled"`
}
