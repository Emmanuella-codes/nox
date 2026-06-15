package pipes

import (
	"context"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/media/messages"
	"github.com/emmanuella-codes/nox/models"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InitiateUploadResponse struct {
	MediaAsset *models.MediaAsset `json:"media_asset"`
	UploadURL  string             `json:"upload_url"`
	StorageKey string             `json:"storage_key"`
}

type MediaCleanupResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}

func (p *MediaPipe) InitiateSetVideoUploadPipe(ctx context.Context, userID uuid.UUID, dto dtos.InitiateSetVideoUploadDTO) *shared.PipeRes[InitiateUploadResponse] {
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if !validSetVideoMime(dto.MimeType) || dto.SizeBytes <= 0 {
		return shared.PipeError[InitiateUploadResponse](messages.Invalid_Media)
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OwnerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[InitiateUploadResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[InitiateUploadResponse](err, "media.upload_persona")
	}
	if persona.UserID != userID || !isDJPersona(persona) {
		return shared.PipeError[InitiateUploadResponse](messages.Forbidden)
	}
	storageKey := setVideoStorageKey(dto.OwnerPersonaID.String())
	playbackURL := p.playbackURL(storageKey)
	asset, err := p.mediaRepo.CreatePendingMediaAsset(ctx, userID, storageKey, playbackURL, dto)
	if err != nil {
		return pipeInternalError[InitiateUploadResponse](err, "media.upload_create")
	}
	return shared.PipeSuccess(messages.Media_Upload_Initiated, &InitiateUploadResponse{
		MediaAsset: asset,
		UploadURL:  p.uploadURL(storageKey),
		StorageKey: storageKey,
	})
}

func (p *MediaPipe) InitiateStoryVideoUploadPipe(ctx context.Context, userID uuid.UUID, dto dtos.InitiateStoryVideoUploadDTO) *shared.PipeRes[InitiateUploadResponse] {
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if !validSetVideoMime(dto.MimeType) || dto.SizeBytes <= 0 {
		return shared.PipeError[InitiateUploadResponse](messages.Invalid_Media)
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OwnerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[InitiateUploadResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[InitiateUploadResponse](err, "media.story_upload_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return shared.PipeError[InitiateUploadResponse](messages.Forbidden)
	}
	storageKey := storyVideoStorageKey(dto.OwnerPersonaID.String())
	playbackURL := p.playbackURL(storageKey)
	asset, err := p.mediaRepo.CreatePendingStoryMediaAsset(ctx, userID, storageKey, playbackURL, dto)
	if err != nil {
		return pipeInternalError[InitiateUploadResponse](err, "media.story_upload_create")
	}
	return shared.PipeSuccess(messages.Media_Upload_Initiated, &InitiateUploadResponse{
		MediaAsset: asset,
		UploadURL:  p.uploadURL(storageKey),
		StorageKey: storageKey,
	})
}

func (p *MediaPipe) CompleteMediaProcessingPipe(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) *shared.PipeRes[models.MediaAsset] {
	dto.PlaybackURL = strings.TrimSpace(dto.PlaybackURL)
	dto.ThumbnailURL = strings.TrimSpace(dto.ThumbnailURL)
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if !validSetVideo(dto.MimeType, dto.DurationSeconds) || dto.SizeBytes <= 0 || dto.PlaybackURL == "" {
		return shared.PipeError[models.MediaAsset](messages.Invalid_Media)
	}
	asset, err := p.mediaRepo.MarkMediaAssetReady(ctx, mediaAssetID, dto)
	if err != nil {
		if err == pgx.ErrNoRows {
			return shared.PipeError[models.MediaAsset](messages.Media_Not_Found)
		}
		return pipeInternalError[models.MediaAsset](err, "media.processing_ready")
	}
	return shared.PipeSuccess(messages.Media_Processing_Updated, asset)
}

func (p *MediaPipe) CompleteStoryMediaProcessingPipe(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) *shared.PipeRes[models.MediaAsset] {
	dto.PlaybackURL = strings.TrimSpace(dto.PlaybackURL)
	dto.ThumbnailURL = strings.TrimSpace(dto.ThumbnailURL)
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if !validStoryVideo(dto.MimeType, dto.DurationSeconds) || dto.SizeBytes <= 0 || dto.PlaybackURL == "" {
		return shared.PipeError[models.MediaAsset](messages.Invalid_Media)
	}
	asset, err := p.mediaRepo.MarkMediaAssetReady(ctx, mediaAssetID, dto)
	if err != nil {
		if err == pgx.ErrNoRows {
			return shared.PipeError[models.MediaAsset](messages.Media_Not_Found)
		}
		return pipeInternalError[models.MediaAsset](err, "media.story_processing_ready")
	}
	return shared.PipeSuccess(messages.Media_Processing_Updated, asset)
}

func (p *MediaPipe) FailMediaProcessingPipe(ctx context.Context, mediaAssetID uuid.UUID) *shared.PipeRes[models.MediaAsset] {
	asset, err := p.mediaRepo.MarkMediaAssetFailed(ctx, mediaAssetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return shared.PipeError[models.MediaAsset](messages.Media_Not_Found)
		}
		return pipeInternalError[models.MediaAsset](err, "media.processing_failed")
	}
	return shared.PipeSuccess(messages.Media_Processing_Updated, asset)
}

func (p *MediaPipe) CleanupOrphanedMediaAssetsPipe(ctx context.Context, olderThan time.Duration, limit int) *shared.PipeRes[MediaCleanupResponse] {
	if olderThan <= 0 {
		olderThan = 24 * time.Hour
	}
	deletedCount, err := p.mediaRepo.DeleteOrphanedMediaAssets(ctx, time.Now().Add(-olderThan), limit)
	if err != nil {
		return pipeInternalError[MediaCleanupResponse](err, "media.cleanup_orphans")
	}
	return shared.PipeSuccess(messages.Media_Cleanup_Completed, &MediaCleanupResponse{DeletedCount: deletedCount})
}
