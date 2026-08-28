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
	cloudinaryclient "github.com/emmanuella-codes/nox/shared/cloudinary/client"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InitiateUploadResponse struct {
	MediaAsset *models.MediaAsset `json:"media_asset"`
	UploadURL  string             `json:"upload_url"`
	StorageKey string             `json:"storage_key"`
}

type CloudinaryPostUploadResponse = cloudinaryclient.UploadSignature

type MediaCleanupResponse struct {
	DeletedCount int64 `json:"deleted_count"`
}

// InitiateSetVideoUploadPipe validates ownership and creates a pending set video asset.
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
	if persona.UserID != userID || !canOwnSetMedia(persona) {
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

// InitiateStoryVideoUploadPipe validates ownership and creates a pending story video asset.
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

// InitiatePostMediaUploadPipe signs a direct upload for post and messaging media assets.
func (p *MediaPipe) InitiatePostMediaUploadPipe(ctx context.Context, userID uuid.UUID, dto dtos.InitiatePostMediaUploadDTO) *shared.PipeRes[CloudinaryPostUploadResponse] {
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if p.cloudinaryClient == nil || !p.cloudinaryClient.Configured() || !validPostMedia(dto.MediaKind, dto.MimeType, dto.SizeBytes, 0) {
		return shared.PipeError[CloudinaryPostUploadResponse](messages.Invalid_Media)
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OwnerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[CloudinaryPostUploadResponse](messages.Persona_Not_Found)
		}
		return pipeInternalError[CloudinaryPostUploadResponse](err, "media.post_upload_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return shared.PipeError[CloudinaryPostUploadResponse](messages.Forbidden)
	}

	resourceType := cloudinaryResourceType(dto.MediaKind)
	publicID := cloudinaryclient.PostPublicID(dto.OwnerPersonaID)
	upload := p.cloudinaryClient.SignUpload(resourceType, publicID)

	return shared.PipeSuccess(messages.Media_Upload_Initiated, &upload)
}

// ConfirmPostMediaUploadPipe stores a ready media asset after upload completion.
func (p *MediaPipe) ConfirmPostMediaUploadPipe(ctx context.Context, userID uuid.UUID, dto dtos.ConfirmPostMediaUploadDTO) *shared.PipeRes[models.MediaAsset] {
	dto.PublicID = strings.TrimSpace(dto.PublicID)
	dto.SecureURL = strings.TrimSpace(dto.SecureURL)
	dto.ThumbnailURL = strings.TrimSpace(dto.ThumbnailURL)
	dto.MimeType = strings.TrimSpace(dto.MimeType)
	if !validPostMedia(dto.MediaKind, dto.MimeType, dto.SizeBytes, dto.DurationSeconds) || dto.PublicID == "" || dto.SecureURL == "" {
		return shared.PipeError[models.MediaAsset](messages.Invalid_Media)
	}
	if dto.MediaKind == models.ImageMediaKind && dto.DurationSeconds <= 0 {
		dto.DurationSeconds = 1
	}
	if (dto.MediaKind == models.VideoMediaKind || dto.MediaKind == models.AudioMediaKind) && dto.DurationSeconds <= 0 {
		return shared.PipeError[models.MediaAsset](messages.Invalid_Media)
	}
	persona, err := p.personaRepo.FindPersonaByID(ctx, dto.OwnerPersonaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return shared.PipeError[models.MediaAsset](messages.Persona_Not_Found)
		}
		return pipeInternalError[models.MediaAsset](err, "media.post_confirm_persona")
	}
	if persona.UserID != userID || persona.PersonaType != models.VisiblePersonaType {
		return shared.PipeError[models.MediaAsset](messages.Forbidden)
	}
	asset, err := p.mediaRepo.CreatePostMediaAsset(ctx, userID, dto)
	if err != nil {
		return pipeInternalError[models.MediaAsset](err, "media.post_confirm_create")
	}
	return shared.PipeSuccess(messages.Media_Asset_Created, asset)
}

// CompleteMediaProcessingPipe marks one pending set media asset as ready.
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

// CompleteStoryMediaProcessingPipe marks one pending story media asset as ready.
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

// FailMediaProcessingPipe marks one pending media asset as failed.
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

// CleanupOrphanedMediaAssetsPipe removes old failed or pending media that is no longer referenced.
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
