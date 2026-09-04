package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *AdminPipe) MePipe(ctx context.Context, userID uuid.UUID) *shared.PipeRes[MeResponse] {
	identity, err := p.adminRepo.FindIdentityByUserID(ctx, userID)
	if err != nil {
		logInternalError(err, "me.find_identity")
		return shared.PipeError[MeResponse](messages.Internal_Error)
	}
	if identity == nil || !identity.IsActive {
		return shared.PipeError[MeResponse](messages.Admin_Access_Denied)
	}

	p.audit(ctx, userID, "admin.auth.me", map[string]any{"role": identity.Role})
	return shared.PipeSuccess(messages.Admin_Profile_Loaded, meResponse(identity))
}
