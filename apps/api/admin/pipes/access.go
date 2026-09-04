package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	"github.com/google/uuid"
)

func (p *AdminPipe) requireSuperAdmin(ctx context.Context, userID uuid.UUID) string {
	membership, err := p.adminRepo.FindMembershipByUserID(ctx, userID)
	if err != nil {
		logInternalError(err, "access.find_membership")
		return string(messages.Internal_Error)
	}
	if membership == nil || !membership.IsActive || membership.Role != models.AdminRoleSuperAdmin {
		return string(messages.Admin_Access_Denied)
	}

	return ""
}

func (p *AdminPipe) protectLastSuperAdmin(ctx context.Context, current *models.AdminMembership, nextRole models.AdminRole, nextActive bool) string {
	if current == nil || current.Role != models.AdminRoleSuperAdmin || !current.IsActive {
		return ""
	}
	if nextRole == models.AdminRoleSuperAdmin && nextActive {
		return ""
	}

	count, err := p.adminRepo.CountActiveMembershipsByRole(ctx, models.AdminRoleSuperAdmin)
	if err != nil {
		logInternalError(err, "access.count_active_super_admins")
		return string(messages.Internal_Error)
	}
	if count <= 1 {
		return string(messages.Last_Super_Admin_Required)
	}

	return ""
}
