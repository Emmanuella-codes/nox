package pipes

import (
	"context"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
)

func (p *AdminPipe) ListAdminMediaPipe(ctx context.Context, actorID uuid.UUID, params adminrepo.ListAdminMediaParams) *shared.PipeRes[models.AdminMediaPage] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminMediaPage](shared.CreatePipeMessage(message))
	}
	page, err := p.adminRepo.ListAdminMedia(ctx, params)
	if err != nil {
		logInternalError(err, "media.list")
		return shared.PipeError[models.AdminMediaPage](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.media.list", nil)
	return shared.PipeSuccess(messages.Admin_Media_Loaded, page)
}

func (p *AdminPipe) GetAdminMediaPipe(ctx context.Context, actorID, mediaID uuid.UUID) *shared.PipeRes[models.AdminMediaAsset] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminMediaAsset](shared.CreatePipeMessage(message))
	}
	item, err := p.adminRepo.FindAdminMedia(ctx, mediaID)
	if err != nil {
		if err == adminrepo.ErrMediaNotFound {
			return shared.PipeError[models.AdminMediaAsset](messages.Media_Not_Found)
		}
		logInternalError(err, "media.get")
		return shared.PipeError[models.AdminMediaAsset](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.media.get", map[string]any{"media_asset_id": mediaID.String()})
	return shared.PipeSuccess(messages.Admin_Media_Loaded, item)
}

func (p *AdminPipe) RetryMediaPipe(ctx context.Context, actorID, mediaID uuid.UUID) *shared.PipeRes[models.AdminMediaAsset] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminMediaAsset](shared.CreatePipeMessage(message))
	}
	current, err := p.adminRepo.FindAdminMedia(ctx, mediaID)
	if err != nil {
		if err == adminrepo.ErrMediaNotFound {
			return shared.PipeError[models.AdminMediaAsset](messages.Media_Not_Found)
		}
		logInternalError(err, "media.retry_find")
		return shared.PipeError[models.AdminMediaAsset](messages.Internal_Error)
	}
	if current.ProcessingStatus != models.PendingMediaStatus && current.ProcessingStatus != models.FailedMediaStatus {
		return shared.PipeError[models.AdminMediaAsset](messages.Invalid_Payload)
	}
	item, err := p.adminRepo.RetryMedia(ctx, mediaID, actorID)
	if err != nil {
		if err == adminrepo.ErrMediaNotFound {
			return shared.PipeError[models.AdminMediaAsset](messages.Media_Not_Found)
		}
		logInternalError(err, "media.retry")
		return shared.PipeError[models.AdminMediaAsset](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.media.retry", map[string]any{"media_asset_id": mediaID.String()})
	return shared.PipeSuccess(messages.Admin_Media_Retried, item)
}

func (p *AdminPipe) CorrectMediaPipe(ctx context.Context, actorID, mediaID uuid.UUID, dto dtos.CorrectMediaDTO) *shared.PipeRes[models.AdminMediaAsset] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminMediaAsset](shared.CreatePipeMessage(message))
	}
	dto.Status = models.MediaProcessingStatus(strings.ToLower(strings.TrimSpace(string(dto.Status))))
	dto.Reason = strings.TrimSpace(dto.Reason)
	if !validMediaProcessingStatus(dto.Status) || (dto.Status == models.ReadyMediaStatus && (strings.TrimSpace(dto.PlaybackURL) == "" || dto.SizeBytes <= 0 || dto.Duration <= 0)) || (dto.Status == models.FailedMediaStatus && dto.Reason == "") {
		return shared.PipeError[models.AdminMediaAsset](messages.Invalid_Payload)
	}
	item, err := p.adminRepo.CorrectMedia(ctx, adminrepo.CorrectMediaParams{MediaAssetID: mediaID, Status: dto.Status, PlaybackURL: strings.TrimSpace(dto.PlaybackURL), ThumbnailURL: strings.TrimSpace(dto.ThumbnailURL), MimeType: strings.TrimSpace(dto.MimeType), Duration: dto.Duration, SizeBytes: dto.SizeBytes, Reason: dto.Reason, AdminID: actorID})
	if err != nil {
		if err == adminrepo.ErrMediaNotFound {
			return shared.PipeError[models.AdminMediaAsset](messages.Media_Not_Found)
		}
		logInternalError(err, "media.correct")
		return shared.PipeError[models.AdminMediaAsset](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.media.correct", map[string]any{"media_asset_id": mediaID.String(), "status": dto.Status, "reason": dto.Reason})
	return shared.PipeSuccess(messages.Admin_Media_Corrected, item)
}

func (p *AdminPipe) ModerateMediaPipe(ctx context.Context, actorID, mediaID uuid.UUID, dto dtos.ModerateMediaDTO) *shared.PipeRes[models.AdminMediaAsset] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.AdminMediaAsset](shared.CreatePipeMessage(message))
	}
	dto.Reason = strings.TrimSpace(dto.Reason)
	if !validModerationStatus(dto.Status) || (dto.Status != models.ModerationStatusActive && dto.Reason == "") {
		return shared.PipeError[models.AdminMediaAsset](messages.Invalid_Payload)
	}
	item, err := p.adminRepo.ModerateMedia(ctx, adminrepo.ModerateMediaParams{MediaAssetID: mediaID, Status: dto.Status, Reason: dto.Reason, AdminID: actorID})
	if err != nil {
		if err == adminrepo.ErrMediaNotFound {
			return shared.PipeError[models.AdminMediaAsset](messages.Media_Not_Found)
		}
		logInternalError(err, "media.moderate")
		return shared.PipeError[models.AdminMediaAsset](messages.Internal_Error)
	}
	p.audit(ctx, actorID, "admin.media.moderate", map[string]any{"media_asset_id": mediaID.String(), "status": dto.Status, "reason": dto.Reason})
	return shared.PipeSuccess(messages.Admin_Media_Moderated, item)
}

func (p *AdminPipe) CleanupAdminMediaPipe(ctx context.Context, actorID uuid.UUID, olderThan time.Duration, limit int) *shared.PipeRes[models.MediaCleanupResult] {
	if message := p.requireActiveAdmin(ctx, actorID); message != "" {
		return shared.PipeError[models.MediaCleanupResult](shared.CreatePipeMessage(message))
	}
	if olderThan <= 0 {
		olderThan = 24 * time.Hour
	}
	deleted, err := p.adminRepo.CleanupOrphanedMedia(ctx, time.Now().Add(-olderThan), limit)
	if err != nil {
		logInternalError(err, "media.cleanup")
		return shared.PipeError[models.MediaCleanupResult](messages.Internal_Error)
	}
	result := &models.MediaCleanupResult{DeletedCount: deleted}
	p.audit(ctx, actorID, "admin.media.cleanup", map[string]any{"deleted_count": deleted, "older_than": olderThan.String()})
	return shared.PipeSuccess(messages.Admin_Media_Cleanup_Completed, result)
}

func validMediaProcessingStatus(status models.MediaProcessingStatus) bool {
	return status == models.PendingMediaStatus || status == models.ReadyMediaStatus || status == models.FailedMediaStatus
}
