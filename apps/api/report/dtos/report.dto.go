package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type CreateReportDTO struct {
	ReporterPersonaID uuid.UUID               `json:"reporter_persona_id" validate:"required"`
	TargetType        models.ReportTargetType `json:"target_type" validate:"required"`
	TargetID          uuid.UUID               `json:"target_id" validate:"required"`
	Reason            string                  `json:"reason" validate:"required,max=80"`
	Description       string                  `json:"description" validate:"max=1000"`
}
