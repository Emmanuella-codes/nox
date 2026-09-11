package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminCrewSummary struct {
	ID             uuid.UUID  `json:"id"`
	EventID        uuid.UUID  `json:"event_id"`
	Name           string     `json:"name"`
	OwnerPersonaID uuid.UUID  `json:"owner_persona_id"`
	Status         CrewStatus `json:"status"`
	MemberCount    int        `json:"member_count"`
	SharingCount   int        `json:"sharing_count"`
	ExpiresAt      time.Time  `json:"expires_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AdminCrewMember struct {
	CrewID                 uuid.UUID      `json:"crew_id"`
	UserID                 uuid.UUID      `json:"user_id"`
	PersonaID              uuid.UUID      `json:"persona_id"`
	Role                   CrewMemberRole `json:"role"`
	LocationSharingEnabled bool           `json:"location_sharing_enabled"`
	JoinedAt               time.Time      `json:"joined_at"`
	LeftAt                 *time.Time     `json:"left_at,omitempty"`
}

type AdminCrewDetail struct {
	Crew    AdminCrewSummary  `json:"crew"`
	Members []AdminCrewMember `json:"members"`
}

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
