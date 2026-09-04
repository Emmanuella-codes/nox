package pipes

import (
	"context"
	"time"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

type ManagedUserResponse struct {
	UserID          uuid.UUID               `json:"user_id"`
	Fullname        string                  `json:"fullname"`
	Email           string                  `json:"email"`
	EmailVerified   bool                    `json:"email_verified"`
	EmailVerifiedAt *time.Time              `json:"email_verified_at,omitempty"`
	Status          models.UserStatus       `json:"status"`
	AdminMembership *ManagedAdminMembership `json:"admin_membership,omitempty"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

type ManagedAdminMembership struct {
	Role     models.AdminRole `json:"role"`
	IsActive bool             `json:"is_active"`
}

func (p *AdminPipe) ListUsersPipe(ctx context.Context, actorID uuid.UUID) *shared.PipeRes[[]ManagedUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]ManagedUserResponse](shared.CreatePipeMessage(message))
	}

	users, err := p.adminRepo.ListManagedUsers(ctx)
	if err != nil {
		logInternalError(err, "users.list")
		return shared.PipeError[[]ManagedUserResponse](messages.Internal_Error)
	}

	responses := make([]ManagedUserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, managedUserResponse(user))
	}

	p.audit(ctx, actorID, "admin.users.list", nil)
	return shared.PipeSuccess(messages.Users_Loaded, &responses)
}

func (p *AdminPipe) GetUserPipe(ctx context.Context, actorID, targetUserID uuid.UUID) *shared.PipeRes[ManagedUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[ManagedUserResponse](shared.CreatePipeMessage(message))
	}

	user, err := p.adminRepo.FindManagedUserByID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "users.get")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}
	if user == nil {
		return shared.PipeError[ManagedUserResponse](messages.User_Not_Found)
	}

	response := managedUserResponse(*user)
	p.audit(ctx, actorID, "admin.users.get", map[string]any{"target_user_id": targetUserID.String()})
	return shared.PipeSuccess(messages.User_Loaded, &response)
}

func (p *AdminPipe) UpdateUserStatusPipe(ctx context.Context, actorID, targetUserID uuid.UUID, status models.UserStatus) *shared.PipeRes[ManagedUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[ManagedUserResponse](shared.CreatePipeMessage(message))
	}
	if !models.ValidUserStatus(status) {
		return shared.PipeError[ManagedUserResponse](messages.Invalid_User_Status)
	}

	user, err := p.adminRepo.UpdateUserStatus(ctx, targetUserID, status)
	if err != nil {
		if err == adminrepo.ErrUserNotFound {
			return shared.PipeError[ManagedUserResponse](messages.User_Not_Found)
		}
		logInternalError(err, "users.update_status")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}

	if status == models.UserStatusSuspended {
		if err := p.revokeAllUserSessions(ctx, targetUserID); err != nil {
			logInternalError(err, "users.revoke_suspended_sessions")
			return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
		}
	}

	response := managedUserResponse(*user)
	p.audit(ctx, actorID, "admin.users.update_status", map[string]any{"target_user_id": targetUserID.String(), "status": status})
	return shared.PipeSuccess(messages.User_Status_Updated, &response)
}

func (p *AdminPipe) RevokeUserSessionsPipe(ctx context.Context, actorID, targetUserID uuid.UUID) *shared.PipeRes[any] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}

	user, err := p.adminRepo.FindManagedUserByID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "users.find_for_revoke")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if user == nil {
		return shared.PipeError[any](messages.User_Not_Found)
	}

	if err := p.revokeAllUserSessions(ctx, targetUserID); err != nil {
		logInternalError(err, "users.revoke_sessions")
		return shared.PipeError[any](messages.Internal_Error)
	}

	p.audit(ctx, actorID, "admin.users.revoke_sessions", map[string]any{"target_user_id": targetUserID.String()})
	return shared.PipeSuccess[any](messages.User_Sessions_Revoked, nil)
}

func (p *AdminPipe) ResendUserVerificationPipe(ctx context.Context, actorID, targetUserID uuid.UUID) *shared.PipeRes[any] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}

	user, err := p.userRepo.FindUserByID(ctx, targetUserID.String())
	if err != nil {
		logInternalError(err, "users.find_for_resend_verification")
		return shared.PipeError[any](messages.Internal_Error)
	}
	if user == nil {
		return shared.PipeError[any](messages.User_Not_Found)
	}
	if user.EmailVerified {
		return shared.PipeError[any](messages.User_Already_Verified)
	}

	if err := p.sendAppVerificationOTP(ctx, user); err != nil {
		logInternalError(err, "users.send_verification")
		return shared.PipeError[any](messages.Internal_Error)
	}

	p.audit(ctx, actorID, "admin.users.resend_verification", map[string]any{"target_user_id": targetUserID.String()})
	return shared.PipeSuccess[any](messages.User_Verification_Sent, nil)
}

func (p *AdminPipe) MarkUserEmailVerifiedPipe(ctx context.Context, actorID, targetUserID uuid.UUID) *shared.PipeRes[ManagedUserResponse] {
	if message := p.requireSuperAdmin(ctx, actorID); message != "" {
		return shared.PipeError[ManagedUserResponse](shared.CreatePipeMessage(message))
	}

	user, err := p.userRepo.FindUserByID(ctx, targetUserID.String())
	if err != nil {
		logInternalError(err, "users.find_for_mark_verified")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}
	if user == nil {
		return shared.PipeError[ManagedUserResponse](messages.User_Not_Found)
	}
	if user.EmailVerified {
		return shared.PipeError[ManagedUserResponse](messages.User_Already_Verified)
	}

	if err := p.userRepo.MarkEmailVerified(ctx, targetUserID.String()); err != nil {
		logInternalError(err, "users.mark_verified")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}
	if err := p.deleteAppEmailVerification(ctx, targetUserID.String()); err != nil {
		logInternalError(err, "users.delete_verification_state")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}

	updatedUser, err := p.adminRepo.FindManagedUserByID(ctx, targetUserID)
	if err != nil {
		logInternalError(err, "users.reload_verified")
		return shared.PipeError[ManagedUserResponse](messages.Internal_Error)
	}
	if updatedUser == nil {
		return shared.PipeError[ManagedUserResponse](messages.User_Not_Found)
	}

	response := managedUserResponse(*updatedUser)
	p.audit(ctx, actorID, "admin.users.mark_email_verified", map[string]any{"target_user_id": targetUserID.String()})
	return shared.PipeSuccess(messages.User_Email_Verified, &response)
}

func managedUserResponse(user models.AdminManagedUser) ManagedUserResponse {
	response := ManagedUserResponse{
		UserID:        user.ID,
		Fullname:      user.Fullname,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
		Status:        models.NormalizeUserStatus(user.Status),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
	if user.EmailVerifiedAt.Valid {
		response.EmailVerifiedAt = &user.EmailVerifiedAt.Time
	}
	if user.AdminRole != nil && user.AdminIsActive != nil {
		response.AdminMembership = &ManagedAdminMembership{
			Role:     *user.AdminRole,
			IsActive: *user.AdminIsActive,
		}
	}

	return response
}
