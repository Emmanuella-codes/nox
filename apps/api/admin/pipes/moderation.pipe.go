package pipes

import (
	"context"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *AdminPipe) ModerateContentPipe(ctx context.Context, actorID uuid.UUID, entityType models.ModerationEntityType, entityID uuid.UUID, dto dtos.ModerateContentDTO) *shared.PipeRes[models.ModerationState] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.ModerationState](shared.CreatePipeMessage(message))
	}
	if !validModerationEntity(entityType) || !validModerationStatus(dto.Status) {
		return shared.PipeError[models.ModerationState](messages.Invalid_Moderation_Status)
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if dto.Status != models.ModerationStatusActive && dto.Reason == "" {
		return shared.PipeError[models.ModerationState](messages.Invalid_Payload)
	}
	state, err := p.adminRepo.Moderate(ctx, adminrepo.ModerateParams{
		EntityType: entityType,
		EntityID:   entityID,
		Status:     dto.Status,
		Reason:     dto.Reason,
		AdminID:    actorID,
	})
	if err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[models.ModerationState](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "moderation.update")
		return shared.PipeError[models.ModerationState](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.moderation.update", map[string]any{
		"entity_type": entityType, "entity_id": entityID.String(), "status": dto.Status, "reason": dto.Reason,
	})
	return shared.PipeSuccess(messages.Moderation_Updated, state)
}

func (p *AdminPipe) ListModerationActionsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListModerationActionsParams) *shared.PipeRes[[]models.ModerationAction] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.ModerationAction](shared.CreatePipeMessage(message))
	}
	actions, err := p.adminRepo.ListModerationActions(ctx, params)
	if err != nil {
		logInternalError(err, "moderation.list")
		return shared.PipeError[[]models.ModerationAction](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.moderation.list", nil)
	return shared.PipeSuccess(messages.Moderation_Actions_Loaded, &actions)
}

func (p *AdminPipe) requireActiveAdmin(ctx context.Context, userID uuid.UUID) string {
	membership, err := p.adminRepo.FindMembershipByUserID(ctx, userID)
	if err != nil {
		logInternalError(err, "moderation.find_membership")
		return string(messages.Internal_Error)
	}
	if membership == nil || !membership.IsActive {
		return string(messages.Admin_Access_Denied)
	}
	return ""
}

func validModerationEntity(entityType models.ModerationEntityType) bool {
	switch entityType {
	case models.ModerationEntityPersona, models.ModerationEntityPost, models.ModerationEntityComment,
		models.ModerationEntityStory, models.ModerationEntityStoryItem, models.ModerationEntitySet, models.ModerationEntityEvent:
		return true
	default:
		return false
	}
}

func validModerationStatus(status models.ModerationStatus) bool {
	return status == models.ModerationStatusActive || status == models.ModerationStatusHidden || status == models.ModerationStatusRemoved
}
