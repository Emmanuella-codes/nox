package pipes

import (
	"context"
	"errors"
	"time"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type AdminUserResponse struct {
	UserID    uuid.UUID        `json:"user_id"`
	Fullname  string           `json:"fullname"`
	Email     string           `json:"email"`
	Role      models.AdminRole `json:"role"`
	IsActive  bool             `json:"is_active"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (p *AdminPipe) ListAdminUsersPipe(ctx context.Context, actorID uuid.UUID) *shared.PipeRes[[]AdminUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]AdminUserResponse](shared.CreatePipeMessage(message))
	}

	identities, err := p.adminRepo.ListIdentities(ctx)
	if err != nil {
		logInternalError(err, "admin_users.list")
		return shared.PipeError[[]AdminUserResponse](messages.Internal_Error)
	}

	responses := make([]AdminUserResponse, 0, len(identities))
	for _, identity := range identities {
		responses = append(responses, adminUserResponse(identity))
	}

	p.audit(ctx, actorID, "admin.users.list", nil)
	return shared.PipeSuccess(messages.Admin_Users_Loaded, &responses)
}

func (p *AdminPipe) GetAdminUserPipe(ctx context.Context, actorID, targetUserID uuid.UUID) *shared.PipeRes[AdminUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}

	identity, err := p.adminRepo.FindIdentityByUserID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "admin_users.get")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}
	if identity == nil {
		return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
	}

	response := adminUserResponse(*identity)
	p.audit(ctx, actorID, "admin.users.get", map[string]any{"target_user_id": targetUserID.String()})
	return shared.PipeSuccess(messages.Admin_User_Loaded, &response)
}

func (p *AdminPipe) CreateAdminUserPipe(ctx context.Context, actorID uuid.UUID, dto dtos.CreateAdminUserDTO) *shared.PipeRes[AdminUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}
	if !models.ValidAdminRole(dto.Role) {
		return shared.PipeError[AdminUserResponse](messages.Invalid_Admin_Role)
	}

	targetUserID, err := uuid.Parse(dto.UserID)
	if err != nil {
		return shared.PipeError[AdminUserResponse](messages.Invalid_Payload)
	}

	user, err := p.userRepo.FindUserByID(ctx, targetUserID.String())
	if err != nil {
		logInternalError(err, "admin_users.find_user")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}
	if user == nil {
		return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
	}

	identity, err := p.adminRepo.CreateMembership(ctx, targetUserID, dto.Role)
	if err != nil {
		if errors.Is(err, adminrepo.ErrMembershipAlreadyExists) {
			return shared.PipeError[AdminUserResponse](messages.Admin_User_Exists)
		}
		logInternalError(err, "admin_users.create")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}

	response := adminUserResponse(*identity)
	p.audit(ctx, actorID, "admin.users.create", map[string]any{"target_user_id": targetUserID.String(), "role": dto.Role})
	return shared.PipeSuccess(messages.Admin_User_Created, &response)
}

func (p *AdminPipe) UpdateAdminUserRolePipe(ctx context.Context, actorID, targetUserID uuid.UUID, dto dtos.UpdateAdminRoleDTO) *shared.PipeRes[AdminUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}
	if !models.ValidAdminRole(dto.Role) {
		return shared.PipeError[AdminUserResponse](messages.Invalid_Admin_Role)
	}

	current, err := p.adminRepo.FindMembershipByUserID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "admin_users.find_membership")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}
	if current == nil {
		return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
	}
	if message := p.protectLastSuperAdmin(ctx, current, dto.Role, current.IsActive); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}

	identity, err := p.adminRepo.UpdateMembershipRole(ctx, targetUserID, dto.Role)
	if err != nil {
		if errors.Is(err, adminrepo.ErrMembershipNotFound) {
			return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
		}
		logInternalError(err, "admin_users.update_role")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}

	response := adminUserResponse(*identity)
	p.audit(ctx, actorID, "admin.users.update_role", map[string]any{"target_user_id": targetUserID.String(), "role": dto.Role})
	return shared.PipeSuccess(messages.Admin_User_Role_Updated, &response)
}

func (p *AdminPipe) UpdateAdminUserStatusPipe(ctx context.Context, actorID, targetUserID uuid.UUID, dto dtos.UpdateAdminStatusDTO) *shared.PipeRes[AdminUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}

	current, err := p.adminRepo.FindMembershipByUserID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "admin_users.find_membership")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}
	if current == nil {
		return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
	}
	if message := p.protectLastSuperAdmin(ctx, current, current.Role, dto.IsActive); message != "" {
		return shared.PipeError[AdminUserResponse](shared.CreatePipeMessage(message))
	}

	identity, err := p.adminRepo.UpdateMembershipStatus(ctx, targetUserID, dto.IsActive)
	if err != nil {
		if errors.Is(err, adminrepo.ErrMembershipNotFound) {
			return shared.PipeError[AdminUserResponse](messages.Admin_User_Not_Found)
		}
		logInternalError(err, "admin_users.update_status")
		return shared.PipeError[AdminUserResponse](messages.Internal_Error)
	}

	response := adminUserResponse(*identity)
	p.audit(ctx, actorID, "admin.users.update_status", map[string]any{"target_user_id": targetUserID.String(), "is_active": dto.IsActive})
	return shared.PipeSuccess(messages.Admin_User_Status_Updated, &response)
}

func adminUserResponse(identity models.AdminIdentity) AdminUserResponse {
	return AdminUserResponse{
		UserID:    identity.UserID,
		Fullname:  identity.Fullname,
		Email:     identity.Email,
		Role:      identity.Role,
		IsActive:  identity.IsActive,
		CreatedAt: identity.CreatedAt,
		UpdatedAt: identity.UpdatedAt,
	}
}
