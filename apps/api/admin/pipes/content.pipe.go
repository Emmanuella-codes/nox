package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *AdminPipe) ListContentPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListAdminContentParams) *shared.PipeRes[models.AdminContentPage] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminContentPage](shared.CreatePipeMessage(message))
	}
	if !validModerationEntity(params.EntityType) {
		return shared.PipeError[models.AdminContentPage](messages.Invalid_Moderation_Status)
	}
	page, err := p.adminRepo.ListAdminContent(ctx, params)
	if err != nil {
		logInternalError(err, "content.list")
		return shared.PipeError[models.AdminContentPage](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.content.list", map[string]any{"entity_type": params.EntityType})
	return shared.PipeSuccess(messages.Admin_Content_Loaded, page)
}

func (p *AdminPipe) GetContentPipe(ctx context.Context, actorID uuid.UUID, entityType models.ModerationEntityType, entityID uuid.UUID) *shared.PipeRes[models.AdminContentRecord] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminContentRecord](shared.CreatePipeMessage(message))
	}
	if !validModerationEntity(entityType) {
		return shared.PipeError[models.AdminContentRecord](messages.Invalid_Moderation_Status)
	}
	item, err := p.adminRepo.FindAdminContent(ctx, entityType, entityID)
	if err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[models.AdminContentRecord](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "content.get")
		return shared.PipeError[models.AdminContentRecord](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.content.get", map[string]any{"entity_type": entityType, "entity_id": entityID.String()})
	return shared.PipeSuccess(messages.Admin_Content_Loaded, item)
}
