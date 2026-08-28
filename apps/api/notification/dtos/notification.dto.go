package dtos

import "github.com/google/uuid"

type MarkNotificationReadDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}

type MarkAllNotificationsReadDTO struct {
	PersonaID uuid.UUID `json:"persona_id" validate:"required"`
}
