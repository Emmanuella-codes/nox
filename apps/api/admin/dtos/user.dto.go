package dtos

import "github.com/emmanuella-codes/nox/models"

type ModerateContentDTO struct {
	Status models.ModerationStatus `json:"status" validate:"required"`
	Reason string                  `json:"reason" validate:"max=500"`
}

type UpdateUserStatusDTO struct {
	Status models.UserStatus `json:"status" validate:"required"`
}
