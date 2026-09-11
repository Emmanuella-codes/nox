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

func (p *AdminPipe) ListFeaturedSetsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListFeaturedSetsParams) *shared.PipeRes[[]models.AdminFeaturedSet] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminFeaturedSet](shared.CreatePipeMessage(message))
	}
	items, err := p.adminRepo.ListFeaturedSets(ctx, params)
	if err != nil {
		logInternalError(err, "featured_sets.list")
		return shared.PipeError[[]models.AdminFeaturedSet](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.featured_sets.list", nil)
	return shared.PipeSuccess(messages.Admin_Featured_Sets_Loaded, &items)
}

func (p *AdminPipe) FeatureSetPipe(ctx context.Context, actorID, setID uuid.UUID, dto dtos.FeatureSetDTO) *shared.PipeRes[models.AdminFeaturedSet] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminFeaturedSet](shared.CreatePipeMessage(message))
	}
	item, err := p.adminRepo.FeatureSet(ctx, setID, actorID, dto.ExpiresAt)
	if err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[models.AdminFeaturedSet](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "featured_sets.feature")
		return shared.PipeError[models.AdminFeaturedSet](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.featured_sets.feature", map[string]any{"set_id": setID.String()})
	return shared.PipeSuccess(messages.Admin_Set_Featured, item)
}

func (p *AdminPipe) UnfeatureSetPipe(ctx context.Context, actorID, setID uuid.UUID) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	if err := p.adminRepo.UnfeatureSet(ctx, setID); err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[any](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "featured_sets.unfeature")
		return shared.PipeError[any](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.featured_sets.unfeature", map[string]any{"set_id": setID.String()})
	return shared.PipeSuccess[any](messages.Admin_Set_Unfeatured, nil)
}

func (p *AdminPipe) ListStoryContributionsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListStoryContributionsParams) *shared.PipeRes[[]models.AdminStoryContribution] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminStoryContribution](shared.CreatePipeMessage(message))
	}
	items, err := p.adminRepo.ListStoryContributions(ctx, params)
	if err != nil {
		logInternalError(err, "story_contributions.list")
		return shared.PipeError[[]models.AdminStoryContribution](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.story_contributions.list", nil)
	return shared.PipeSuccess(messages.Admin_Contributions_Loaded, &items)
}

func (p *AdminPipe) ReviewStoryContributionPipe(ctx context.Context, actorID, requestID uuid.UUID, dto dtos.ReviewStoryContributionDTO) *shared.PipeRes[models.AdminStoryContribution] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminStoryContribution](shared.CreatePipeMessage(message))
	}
	dto.Action, dto.Reason = strings.ToLower(strings.TrimSpace(dto.Action)), strings.TrimSpace(dto.Reason)
	if (dto.Action != "reject" && dto.Action != "remove") || dto.Reason == "" {
		return shared.PipeError[models.AdminStoryContribution](messages.Invalid_Payload)
	}
	item, err := p.adminRepo.ReviewStoryContribution(ctx, adminrepo.ReviewStoryContributionParams{RequestID: requestID, Action: dto.Action, Reason: dto.Reason, AdminID: actorID})
	if err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[models.AdminStoryContribution](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "story_contributions.review")
		return shared.PipeError[models.AdminStoryContribution](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.story_contributions.review", map[string]any{"request_id": requestID.String(), "action": dto.Action, "reason": dto.Reason})
	return shared.PipeSuccess(messages.Admin_Contribution_Reviewed, item)
}

func (p *AdminPipe) ListAdminHighlightsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListAdminHighlightsParams) *shared.PipeRes[[]models.AdminHighlight] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[[]models.AdminHighlight](shared.CreatePipeMessage(message))
	}
	items, err := p.adminRepo.ListAdminHighlights(ctx, params)
	if err != nil {
		logInternalError(err, "highlights.list")
		return shared.PipeError[[]models.AdminHighlight](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.highlights.list", nil)
	return shared.PipeSuccess(messages.Admin_Highlights_Loaded, &items)
}

func (p *AdminPipe) RemoveAdminHighlightPipe(ctx context.Context, actorID, highlightID uuid.UUID, highlightType, reason string) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return shared.PipeError[any](messages.Invalid_Payload)
	}
	if err := p.adminRepo.RemoveAdminHighlight(ctx, highlightType, highlightID); err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[any](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "highlights.remove")
		return shared.PipeError[any](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.highlights.remove", map[string]any{"highlight_id": highlightID.String(), "highlight_type": highlightType, "reason": reason})
	return shared.PipeSuccess[any](messages.Admin_Highlight_Removed, nil)
}

func (p *AdminPipe) GetReportScopedPrivateContentPipe(ctx context.Context, actorID, reportID uuid.UUID, reason string) *shared.PipeRes[models.AdminContentRecord] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminContentRecord](shared.CreatePipeMessage(message))
	}
	if strings.TrimSpace(reason) == "" {
		return shared.PipeError[models.AdminContentRecord](messages.Invalid_Payload)
	}
	report, err := p.adminRepo.FindReportByID(ctx, reportID)
	if err != nil {
		if err == adminrepo.ErrReportNotFound {
			return shared.PipeError[models.AdminContentRecord](messages.Report_Not_Found)
		}
		logInternalError(err, "private_content.report")
		return shared.PipeError[models.AdminContentRecord](messages.Internal_Error)
	}
	var entity models.ModerationEntityType
	switch report.TargetType {
	case models.ReportTargetMessage:
		return p.getMessageEvidencePipe(ctx, actorID, reportID, report.TargetID, reason)
	case models.ReportTargetStory:
		entity = models.ModerationEntityStory
	case models.ReportTargetStoryItem:
		entity = models.ModerationEntityStoryItem
	default:
		return shared.PipeError[models.AdminContentRecord](messages.Invalid_Payload)
	}
	item, err := p.adminRepo.FindAdminContent(ctx, entity, report.TargetID)
	if err != nil {
		if err == adminrepo.ErrModerationEntityNotFound {
			return shared.PipeError[models.AdminContentRecord](messages.Moderation_Entity_Not_Found)
		}
		logInternalError(err, "private_content.load")
		return shared.PipeError[models.AdminContentRecord](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.private_content.read", map[string]any{"report_id": reportID.String(), "target_type": report.TargetType, "target_id": report.TargetID.String(), "reason": reason})
	return shared.PipeSuccess(messages.Admin_Private_Content_Loaded, item)
}
