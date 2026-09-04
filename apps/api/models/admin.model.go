package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminRole string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleModerator  AdminRole = "moderator"
	AdminRoleSupport    AdminRole = "support"
	AdminRoleOps        AdminRole = "ops"
	AdminRoleFinance    AdminRole = "finance"
)

var validAdminRoles = map[AdminRole]struct{}{
	AdminRoleSuperAdmin: {},
	AdminRoleModerator:  {},
	AdminRoleSupport:    {},
	AdminRoleOps:        {},
	AdminRoleFinance:    {},
}

type AdminMembership struct {
	UserID    uuid.UUID `json:"user_id"`
	Role      AdminRole `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminIdentity struct {
	UserID    uuid.UUID `json:"user_id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Role      AdminRole `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ValidAdminRole(role AdminRole) bool {
	_, ok := validAdminRoles[role]
	return ok
}
