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

func (p *AdminPipe) ListAdminHashtagsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListAdminHashtagsParams) *shared.PipeRes[models.AdminHashtagPage] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminHashtagPage](shared.CreatePipeMessage(message))
	}
	page, err := p.adminRepo.ListAdminHashtags(ctx, params)
	if err != nil {
		logInternalError(err, "hashtags.list")
		return shared.PipeError[models.AdminHashtagPage](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.hashtags.list", nil)
	return shared.PipeSuccess(messages.Admin_Hashtags_Loaded, page)
}

func (p *AdminPipe) ModerateHashtagPipe(ctx context.Context, actorID uuid.UUID, tag string, dto dtos.ModerateHashtagDTO) *shared.PipeRes[models.AdminHashtag] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminHashtag](shared.CreatePipeMessage(message))
	}
	if !validModerationStatus(dto.Status) {
		return shared.PipeError[models.AdminHashtag](messages.Invalid_Moderation_Status)
	}
	item, err := p.adminRepo.ModerateHashtag(ctx, adminrepo.ModerateHashtagParams{Tag: tag, Status: dto.Status, Reason: strings.TrimSpace(dto.Reason), AdminID: actorID})
	if err != nil {
		if err == adminrepo.ErrHashtagNotFound {
			return shared.PipeError[models.AdminHashtag](messages.Hashtag_Not_Found)
		}
		logInternalError(err, "hashtags.moderate")
		return shared.PipeError[models.AdminHashtag](messages.Internal_Error)
	}
	deleted := p.invalidateSearchCache(ctx)
	p.audit(ctx, actorID, "admin.hashtags.moderate", map[string]any{"tag": tag, "status": dto.Status, "cache_keys_deleted": deleted})
	return shared.PipeSuccess(messages.Admin_Hashtag_Updated, item)
}

func (p *AdminPipe) SuppressHashtagPipe(ctx context.Context, actorID uuid.UUID, tag string, dto dtos.SuppressHashtagDTO) *shared.PipeRes[models.AdminHashtag] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminHashtag](shared.CreatePipeMessage(message))
	}
	item, err := p.adminRepo.SuppressHashtag(ctx, adminrepo.SuppressHashtagParams{Tag: tag, Reason: strings.TrimSpace(dto.Reason), AdminID: actorID, ExpiresAt: dto.ExpiresAt})
	if err != nil {
		if err == adminrepo.ErrHashtagNotFound {
			return shared.PipeError[models.AdminHashtag](messages.Hashtag_Not_Found)
		}
		logInternalError(err, "hashtags.suppress")
		return shared.PipeError[models.AdminHashtag](messages.Internal_Error)
	}
	deleted := p.invalidateSearchCache(ctx)
	p.audit(ctx, actorID, "admin.hashtags.suppress", map[string]any{"tag": tag, "cache_keys_deleted": deleted})
	return shared.PipeSuccess(messages.Admin_Hashtag_Suppressed, item)
}

func (p *AdminPipe) UnsuppressHashtagPipe(ctx context.Context, actorID uuid.UUID, tag string) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	if err := p.adminRepo.UnsuppressHashtag(ctx, tag); err != nil {
		if err == adminrepo.ErrHashtagNotFound {
			return shared.PipeError[any](messages.Hashtag_Not_Found)
		}
		logInternalError(err, "hashtags.unsuppress")
		return shared.PipeError[any](messages.Internal_Error)
	}
	deleted := p.invalidateSearchCache(ctx)
	p.audit(ctx, actorID, "admin.hashtags.unsuppress", map[string]any{"tag": tag, "cache_keys_deleted": deleted})
	return shared.PipeSuccess[any](messages.Admin_Hashtag_Unsuppressed, nil)
}

func (p *AdminPipe) ListSearchSuppressionsPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListSearchSuppressionsParams) *shared.PipeRes[models.SearchSuppressionPage] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.SearchSuppressionPage](shared.CreatePipeMessage(message))
	}
	page, err := p.adminRepo.ListSearchSuppressions(ctx, params)
	if err != nil {
		logInternalError(err, "search.suppressions.list")
		return shared.PipeError[models.SearchSuppressionPage](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.search.suppressions.list", nil)
	return shared.PipeSuccess(messages.Admin_Search_Suppressions_Loaded, page)
}

func (p *AdminPipe) CreateSearchSuppressionPipe(ctx context.Context, actorID uuid.UUID, dto dtos.CreateSearchSuppressionDTO) *shared.PipeRes[models.SearchSuppression] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.SearchSuppression](shared.CreatePipeMessage(message))
	}
	item, err := p.adminRepo.CreateSearchSuppression(ctx, adminrepo.CreateSearchSuppressionParams{Query: dto.Query, Reason: strings.TrimSpace(dto.Reason), AdminID: actorID, ExpiresAt: dto.ExpiresAt})
	if err != nil {
		logInternalError(err, "search.suppressions.create")
		return shared.PipeError[models.SearchSuppression](messages.Internal_Error)
	}
	deleted := p.invalidateSearchCache(ctx)
	p.audit(ctx, actorID, "admin.search.suppressions.create", map[string]any{"query": item.NormalizedQuery, "cache_keys_deleted": deleted})
	return shared.PipeSuccess(messages.Admin_Search_Suppression_Created, item)
}

func (p *AdminPipe) DeleteSearchSuppressionPipe(ctx context.Context, actorID, id uuid.UUID) *shared.PipeRes[any] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[any](shared.CreatePipeMessage(message))
	}
	if err := p.adminRepo.DeleteSearchSuppression(ctx, id); err != nil {
		if err == adminrepo.ErrSearchSuppressionNotFound {
			return shared.PipeError[any](messages.Search_Suppression_Not_Found)
		}
		logInternalError(err, "search.suppressions.delete")
		return shared.PipeError[any](messages.Internal_Error)
	}
	deleted := p.invalidateSearchCache(ctx)
	p.audit(ctx, actorID, "admin.search.suppressions.delete", map[string]any{"suppression_id": id.String(), "cache_keys_deleted": deleted})
	return shared.PipeSuccess[any](messages.Admin_Search_Suppression_Deleted, nil)
}

func (p *AdminPipe) ReindexSearchPipe(ctx context.Context, actorID uuid.UUID) *shared.PipeRes[models.SearchReindexResult] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.SearchReindexResult](shared.CreatePipeMessage(message))
	}
	result := &models.SearchReindexResult{CacheKeysDeleted: p.invalidateSearchCache(ctx)}
	p.audit(ctx, actorID, "admin.search.reindex", map[string]any{"cache_keys_deleted": result.CacheKeysDeleted})
	return shared.PipeSuccess(messages.Admin_Search_Reindexed, result)
}
