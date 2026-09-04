package dtos

import "github.com/emmanuella-codes/nox/models"

type CreateAdminUserDTO struct {
	UserID string           `json:"user_id" validate:"required,uuid"`
	Role   models.AdminRole `json:"role" validate:"required"`
}

type UpdateAdminRoleDTO struct {
	Role models.AdminRole `json:"role" validate:"required"`
}

type UpdateAdminStatusDTO struct {
	IsActive bool `json:"is_active"`
}
