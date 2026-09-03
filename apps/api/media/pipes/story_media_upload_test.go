package pipes

import (
	"context"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/config"
	media_dtos "github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	persona_dtos "github.com/emmanuella-codes/nox/persona/dtos"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestInitiateStoryMediaUploadPipeAcceptsImages(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	repo := &storyMediaUploadRepo{}
	pipe := NewMediaPipe(
		repo,
		&storyMediaPersonaRepo{persona: &models.Persona{
			ID:          personaID,
			UserID:      userID,
			PersonaType: models.VisiblePersonaType,
		}},
		&config.Config{MediaUploadBaseURL: "https://upload.test", MediaPublicBaseURL: "https://cdn.test"},
	)

	res := pipe.InitiateStoryMediaUploadPipe(context.Background(), userID, media_dtos.InitiateStoryMediaUploadDTO{
		OwnerPersonaID: personaID,
		MediaKind:      models.ImageMediaKind,
		MimeType:       "image/jpeg",
		SizeBytes:      1024,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.createdStoryDTO == nil || repo.createdStoryDTO.MediaKind != models.ImageMediaKind {
		t.Fatalf("expected image upload dto, got %#v", repo.createdStoryDTO)
	}
	if res.Data == nil || res.Data.MediaAsset == nil || res.Data.MediaAsset.MediaKind != models.ImageMediaKind {
		t.Fatalf("expected image media asset, got %#v", res.Data)
	}
	if res.Data.MediaAsset.DurationSeconds != storyImageDurationSeconds {
		t.Fatalf("expected %d second image duration, got %d", storyImageDurationSeconds, res.Data.MediaAsset.DurationSeconds)
	}
}

func TestCompleteStoryMediaProcessingPipeNormalizesImageDuration(t *testing.T) {
	mediaAssetID := uuid.New()
	userID := uuid.New()
	personaID := uuid.New()
	repo := &storyMediaUploadRepo{
		foundAsset: &models.MediaAsset{
			ID:               mediaAssetID,
			OwnerUserID:      userID,
			OwnerPersonaID:   personaID,
			MediaKind:        models.ImageMediaKind,
			ProcessingStatus: models.PendingMediaStatus,
		},
	}
	pipe := NewMediaPipe(repo, &storyMediaPersonaRepo{}, &config.Config{})

	res := pipe.CompleteStoryMediaProcessingPipe(context.Background(), mediaAssetID, media_dtos.CompleteStoryMediaProcessingDTO{
		PlaybackURL:  "https://cdn.test/story.jpg",
		ThumbnailURL: "",
		MimeType:     "image/jpeg",
		SizeBytes:    2048,
	})
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if repo.readyDTO == nil {
		t.Fatal("expected ready dto to be captured")
	}
	if repo.readyDTO.DurationSeconds != storyImageDurationSeconds {
		t.Fatalf("expected normalized image duration %d, got %d", storyImageDurationSeconds, repo.readyDTO.DurationSeconds)
	}
	if res.Data == nil || res.Data.DurationSeconds != storyImageDurationSeconds {
		t.Fatalf("expected ready asset duration %d, got %#v", storyImageDurationSeconds, res.Data)
	}
}

type storyMediaUploadRepo struct {
	createdStoryDTO *media_dtos.InitiateStoryMediaUploadDTO
	foundAsset      *models.MediaAsset
	readyDTO        *media_dtos.CompleteMediaProcessingDTO
}

func (r *storyMediaUploadRepo) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreateMediaAsset")
}

func (r *storyMediaUploadRepo) CreatePostMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.ConfirmPostMediaUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePostMediaAsset")
}

func (r *storyMediaUploadRepo) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePendingMediaAsset")
}

func (r *storyMediaUploadRepo) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateStoryMediaUploadDTO) (*models.MediaAsset, error) {
	r.createdStoryDTO = &dto
	return &models.MediaAsset{
		ID:               uuid.New(),
		OwnerUserID:      ownerUserID,
		OwnerPersonaID:   dto.OwnerPersonaID,
		MediaKind:        dto.MediaKind,
		PlaybackURL:      playbackURL,
		MimeType:         dto.MimeType,
		DurationSeconds:  storyMediaDuration(dto.MediaKind, 1),
		SizeBytes:        dto.SizeBytes,
		ProcessingStatus: models.PendingMediaStatus,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

func (r *storyMediaUploadRepo) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	if r.foundAsset == nil || r.foundAsset.ID != mediaAssetID {
		return nil, pgx.ErrNoRows
	}
	return r.foundAsset, nil
}

func (r *storyMediaUploadRepo) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto media_dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	r.readyDTO = &dto
	return &models.MediaAsset{
		ID:               mediaAssetID,
		OwnerUserID:      uuid.New(),
		OwnerPersonaID:   uuid.New(),
		MediaKind:        models.ImageMediaKind,
		PlaybackURL:      dto.PlaybackURL,
		ThumbnailURL:     dto.ThumbnailURL,
		MimeType:         dto.MimeType,
		DurationSeconds:  dto.DurationSeconds,
		SizeBytes:        dto.SizeBytes,
		ProcessingStatus: models.ReadyMediaStatus,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}, nil
}

func (r *storyMediaUploadRepo) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	panic("unexpected call to MarkMediaAssetFailed")
}

func (r *storyMediaUploadRepo) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	panic("unexpected call to DeleteOrphanedMediaAssets")
}

type storyMediaPersonaRepo struct {
	persona *models.Persona
}

func (r *storyMediaPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto persona_dtos.CreatePersonaDTO) (*models.Persona, error) {
	panic("unexpected call to CreatePersona")
}

func (r *storyMediaPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	if r.persona == nil || r.persona.ID != personaID {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return r.persona, nil
}

func (r *storyMediaPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	panic("unexpected call to FindPersonasByUserID")
}

func (r *storyMediaPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	panic("unexpected call to FindPersonaByHandle")
}

func (r *storyMediaPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto persona_dtos.UpdatePersonaDTO) (*models.Persona, error) {
	panic("unexpected call to UpdatePersona")
}

var _ media_repo.MediaRepository = (*storyMediaUploadRepo)(nil)
var _ persona_repo.PersonaRepository = (*storyMediaPersonaRepo)(nil)
