package models

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type AdminRole string
type AdminScope string

const (
	AdminRoleSuperAdmin AdminRole = "super_admin"
	AdminRoleModerator  AdminRole = "moderator"
	AdminRoleSupport    AdminRole = "support"
	AdminRoleOps        AdminRole = "ops"
	AdminRoleFinance    AdminRole = "finance"
)

const (
	AdminScopeStaff    AdminScope = "staff"
	AdminScopeModerate AdminScope = "moderate"
	AdminScopeSupport  AdminScope = "support"
	AdminScopeOps      AdminScope = "ops"
	AdminScopeFinance  AdminScope = "finance"
)

var adminRoleScopes = map[AdminRole][]AdminScope{
	AdminRoleSuperAdmin: {AdminScopeStaff, AdminScopeModerate, AdminScopeSupport, AdminScopeOps, AdminScopeFinance},
	AdminRoleModerator:  {AdminScopeModerate},
	AdminRoleSupport:    {AdminScopeSupport},
	AdminRoleOps:        {AdminScopeOps},
	AdminRoleFinance:    {AdminScopeFinance},
}

type AdminMembership struct {
	UserID    uuid.UUID `json:"user_id"`
	Role      AdminRole `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminIdentity struct {
	UserID    uuid.UUID    `json:"user_id"`
	Fullname  string       `json:"fullname"`
	Email     string       `json:"email"`
	Role      AdminRole    `json:"role"`
	Scopes    []AdminScope `json:"scopes"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func ValidAdminRole(role AdminRole) bool {
	_, ok := adminRoleScopes[role]
	return ok
}

func AdminScopesForRole(role AdminRole) []AdminScope {
	scopes := adminRoleScopes[role]
	if len(scopes) == 0 {
		return []AdminScope{}
	}

	cloned := make([]AdminScope, len(scopes))
	copy(cloned, scopes)
	return cloned
}

func AdminRoleHasScope(role AdminRole, scope AdminScope) bool {
	return slices.Contains(adminRoleScopes[role], scope)
}
