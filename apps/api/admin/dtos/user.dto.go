package dtos

import "github.com/emmanuella-codes/nox/models"

type UpdateUserStatusDTO struct {
	Status models.UserStatus `json:"status" validate:"required"`
}
