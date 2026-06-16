package pipes

import (
	"context"
	"testing"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	personadtos "github.com/emmanuella-codes/nox/persona/dtos"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	setdtos "github.com/emmanuella-codes/nox/set/dtos"
	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/google/uuid"
)

func TestCreateSetPipeAllowsDJWithReadyShortVideo(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	mediaID := uuid.New()
	pipe := NewSetPipe(&setTestRepo{}, &setTestMediaRepo{
		assets: map[uuid.UUID]*models.MediaAsset{
			mediaID: {
				ID:               mediaID,
				OwnerUserID:      userID,
				OwnerPersonaID:   personaID,
				MediaKind:        models.VideoMediaKind,
				ProcessingStatus: models.ReadyMediaStatus,
				DurationSeconds:  899,
			},
		},
	}, &setTestPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID, PersonaType: models.VisiblePersonaType, Category: models.DJPersonaCategory},
		},
	})

	res := pipe.CreateSetPipe(context.Background(), userID, setdtos.CreateSetDTO{
		PersonaID:    personaID,
		MediaAssetID: mediaID,
		Title:        "Late night set",
		GenreTags:    []string{"amapiano"},
	})
	if !res.Success {
		t.Fatalf("expected set create success, got %q", res.Message)
	}
	if res.Data.DurationSeconds != 899 {
		t.Fatalf("expected duration from media asset, got %d", res.Data.DurationSeconds)
	}
}

func TestCreateSetPipeAllowsOrganizerWithReadyShortVideo(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	mediaID := uuid.New()
	pipe := NewSetPipe(&setTestRepo{}, &setTestMediaRepo{
		assets: map[uuid.UUID]*models.MediaAsset{
			mediaID: {
				ID:               mediaID,
				OwnerUserID:      userID,
				OwnerPersonaID:   personaID,
				MediaKind:        models.VideoMediaKind,
				ProcessingStatus: models.ReadyMediaStatus,
				DurationSeconds:  600,
			},
		},
	}, &setTestPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID, PersonaType: models.VisiblePersonaType, Category: models.OrganizerPersonaCategory},
		},
	})

	res := pipe.CreateSetPipe(context.Background(), userID, setdtos.CreateSetDTO{
		PersonaID:    personaID,
		MediaAssetID: mediaID,
		Title:        "Organizer set",
		GenreTags:    []string{"afro-house"},
	})
	if !res.Success {
		t.Fatalf("expected set create success, got %q", res.Message)
	}
}

func TestCreateSetPipeRejectsOverlongVideo(t *testing.T) {
	userID := uuid.New()
	personaID := uuid.New()
	mediaID := uuid.New()
	pipe := NewSetPipe(&setTestRepo{}, &setTestMediaRepo{
		assets: map[uuid.UUID]*models.MediaAsset{
			mediaID: {
				ID:               mediaID,
				OwnerUserID:      userID,
				OwnerPersonaID:   personaID,
				MediaKind:        models.VideoMediaKind,
				ProcessingStatus: models.ReadyMediaStatus,
				DurationSeconds:  901,
			},
		},
	}, &setTestPersonaRepo{
		personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID, PersonaType: models.VisiblePersonaType, Category: models.DJPersonaCategory},
		},
	})

	res := pipe.CreateSetPipe(context.Background(), userID, setdtos.CreateSetDTO{
		PersonaID:    personaID,
		MediaAssetID: mediaID,
		Title:        "Too long",
		GenreTags:    []string{"afro-house"},
	})
	if res.Message != messages.Invalid_Set {
		t.Fatalf("expected invalid set, got %q", res.Message)
	}
}

type setTestRepo struct {
	sets map[uuid.UUID]*models.Set
}

func (r *setTestRepo) CreateSet(ctx context.Context, authorUserID uuid.UUID, durationSeconds int, dto setdtos.CreateSetDTO) (*models.Set, error) {
	return &models.Set{
		ID:              uuid.New(),
		AuthorUserID:    authorUserID,
		PersonaID:       dto.PersonaID,
		MediaAssetID:    dto.MediaAssetID,
		Title:           dto.Title,
		Description:     dto.Description,
		GenreTags:       dto.GenreTags,
		DurationSeconds: durationSeconds,
	}, nil
}

func (r *setTestRepo) FindSetByID(ctx context.Context, setID uuid.UUID) (*models.Set, error) {
	if r.sets == nil {
		return nil, nil
	}
	return r.sets[setID], nil
}

func (r *setTestRepo) FindSets(ctx context.Context, limit int, offset int) ([]*models.Set, error) {
	return nil, nil
}

func (r *setTestRepo) FindSetsWithFilters(ctx context.Context, genreTag string, sort string, limit int, offset int) ([]*models.Set, error) {
	return nil, nil
}

func (r *setTestRepo) FindSetsByPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Set, error) {
	return nil, nil
}

func (r *setTestRepo) DeleteSet(ctx context.Context, setID uuid.UUID) error {
	return nil
}

func (r *setTestRepo) LikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error {
	return nil
}

func (r *setTestRepo) UnlikeSet(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) error {
	return nil
}

func (r *setTestRepo) HasSetLike(ctx context.Context, personaID uuid.UUID, setID uuid.UUID) (bool, error) {
	return false, nil
}

func (r *setTestRepo) FindLikedSetIDs(ctx context.Context, personaID uuid.UUID, setIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	return map[uuid.UUID]bool{}, nil
}

func (r *setTestRepo) IncrementPlayCount(ctx context.Context, setID uuid.UUID) error {
	return nil
}

func (r *setTestRepo) CreateSetComment(ctx context.Context, personaID uuid.UUID, setID uuid.UUID, body string, parentID uuid.UUID) (*models.SetComment, error) {
	return nil, nil
}

func (r *setTestRepo) FindSetComments(ctx context.Context, setID uuid.UUID, limit int, offset int) ([]*models.SetComment, error) {
	return nil, nil
}

type setTestMediaRepo struct {
	assets map[uuid.UUID]*models.MediaAsset
}

func (r *setTestMediaRepo) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	return nil, nil
}

func (r *setTestMediaRepo) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	return nil, nil
}

func (r *setTestMediaRepo) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error) {
	return nil, nil
}

func (r *setTestMediaRepo) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	return r.assets[mediaAssetID], nil
}

func (r *setTestMediaRepo) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	return nil, nil
}

func (r *setTestMediaRepo) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	return nil, nil
}

func (r *setTestMediaRepo) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return 0, nil
}

type setTestPersonaRepo struct {
	personas map[uuid.UUID]*models.Persona
}

func (r *setTestPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto personadtos.CreatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}

func (r *setTestPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona := r.personas[personaID]
	if persona == nil {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

func (r *setTestPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	return nil, nil
}

func (r *setTestPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	return nil, persona_repo.ErrPersonaNotFound
}

func (r *setTestPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto personadtos.UpdatePersonaDTO) (*models.Persona, error) {
	return nil, nil
}
