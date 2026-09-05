package dtos

import (
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

type ModerateContentDTO struct {
	Status models.ModerationStatus `json:"status" validate:"required"`
	Reason string                  `json:"reason" validate:"max=500"`
}

type UpdateReportDTO struct {
	Status          models.ReportStatus `json:"status" validate:"required"`
	AssignedAdminID *uuid.UUID          `json:"assigned_admin_id"`
	ResolutionNote  string              `json:"resolution_note" validate:"max=1000"`
}

type ReportActionDTO struct {
	Action     string `json:"action" validate:"required,max=20"`
	Reason     string `json:"reason" validate:"max=500"`
	Resolution string `json:"resolution" validate:"max=1000"`
}

type UpdateUserStatusDTO struct {
	Status models.UserStatus `json:"status" validate:"required"`
}
