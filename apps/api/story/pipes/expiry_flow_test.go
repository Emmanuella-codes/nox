package pipes

import (
	"context"
	"testing"
	"time"

	media_dtos "github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

// Keeps additions closed after the configured twenty-hour cutoff.
func TestAddStoryItemPipeBlocksClosedStories(t *testing.T) {
	userID, personaID, storyID, assetID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	storyModel := &models.Story{
		ID:               storyID,
		EventID:          uuid.New(),
		OwnerUserID:      userID,
		OwnerPersonaID:   personaID,
		Title:            "late story",
		ContributionMode: models.PublicStoryContributionMode,
		ExpiresAt:        time.Now().Add(4 * time.Hour),
		CreatedAt:        time.Now().Add(-21 * time.Hour),
	}
	pipe := NewStoryPipe(
		&storyHighlightTestRepo{story: storyModel, acceptsAdditions: false},
		&storyHighlightEventRepo{},
		&storyHighlightPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			personaID: {ID: personaID, UserID: userID, Handle: "owner", DisplayName: "Owner"},
		}},
		&storyExpiryMediaRepo{asset: &models.MediaAsset{
			ID:               assetID,
			OwnerUserID:      userID,
			OwnerPersonaID:   personaID,
			MediaKind:        models.VideoMediaKind,
			DurationSeconds:  30,
			ProcessingStatus: models.ReadyMediaStatus,
		}},
		&storyHighlightFollowRepo{},
	)

	res := pipe.AddStoryItemPipe(context.Background(), userID, storyID, dtos.AddStoryItemDTO{
		ContributorPersonaID: personaID,
		MediaAssetID:         assetID,
		PostingMode:          models.PublicPostingMode,
	})
	if res.Message != messages.Story_Closed_For_Additions {
		t.Fatalf("expected %q, got %q", messages.Story_Closed_For_Additions, res.Message)
	}
}

// Returns preserved highlight items even after some of them expire from the normal story rail.
func TestProfileHighlightResponseKeepsExpiredItems(t *testing.T) {
	storyID := uuid.New()
	ownerUserID := uuid.New()
	ownerPersonaID := uuid.New()
	activeItemID := uuid.New()
	expiredItemID := uuid.New()
	activeAssetID := uuid.New()
	expiredAssetID := uuid.New()
	now := time.Now()
	repo := &storyHighlightTestRepo{
		story: &models.Story{
			ID:               storyID,
			EventID:          uuid.New(),
			OwnerUserID:      ownerUserID,
			OwnerPersonaID:   ownerPersonaID,
			Title:            "event a",
			ContributionMode: models.PublicStoryContributionMode,
			ExpiresAt:        now.Add(6 * time.Hour),
			CreatedAt:        now.Add(-5 * time.Hour),
		},
		items: []*models.StoryItem{
			{ID: expiredItemID, StoryID: storyID, MediaAssetID: expiredAssetID, ContributorUserID: ownerUserID, ContributorPersonaID: ownerPersonaID, PostingMode: models.PublicPostingMode, DurationSeconds: 10, Position: 1, ExpiresAt: now.Add(-time.Hour), CreatedAt: now.Add(-25 * time.Hour)},
			{ID: activeItemID, StoryID: storyID, MediaAssetID: activeAssetID, ContributorUserID: ownerUserID, ContributorPersonaID: ownerPersonaID, PostingMode: models.PublicPostingMode, DurationSeconds: 12, Position: 2, ExpiresAt: now.Add(5 * time.Hour), CreatedAt: now.Add(-3 * time.Hour)},
		},
		acceptsAdditions: true,
	}
	pipe := NewStoryPipe(
		repo,
		&storyHighlightEventRepo{},
		&storyHighlightPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			ownerPersonaID: {ID: ownerPersonaID, UserID: ownerUserID, Handle: "owner", DisplayName: "Owner"},
		}},
		&storyExpiryMediaRepo{assets: map[uuid.UUID]*models.MediaAsset{
			activeAssetID:  {ID: activeAssetID, OwnerUserID: ownerUserID, OwnerPersonaID: ownerPersonaID, MediaKind: models.VideoMediaKind, DurationSeconds: 12, ProcessingStatus: models.ReadyMediaStatus},
			expiredAssetID: {ID: expiredAssetID, OwnerUserID: ownerUserID, OwnerPersonaID: ownerPersonaID, MediaKind: models.VideoMediaKind, DurationSeconds: 10, ProcessingStatus: models.ReadyMediaStatus},
		}},
		&storyHighlightFollowRepo{},
	)

	res := pipe.ListProfileStoryHighlightsPipe(context.Background(), ownerPersonaID, nil, nil)
	if !res.Success {
		t.Fatalf("expected success, got %q", res.Message)
	}
	if res.Data == nil || len(*res.Data) != 0 {
		t.Fatalf("expected no highlights before storing one, got %#v", res.Data)
	}
	highlight, err := repo.AddProfileStoryHighlight(context.Background(), ownerPersonaID, storyID)
	if err != nil {
		t.Fatalf("add highlight: %v", err)
	}
	response, err := pipe.profileHighlightResponse(context.Background(), highlight, nil)
	if err != nil {
		t.Fatalf("profileHighlightResponse: %v", err)
	}
	if len(response.Story.Items) != 2 {
		t.Fatalf("expected 2 preserved items, got %d", len(response.Story.Items))
	}
	if response.Story.Items[0].ID != expiredItemID.String() || response.Story.Items[1].ID != activeItemID.String() {
		t.Fatalf("expected preserved item order, got %#v", response.Story.Items)
	}
}

type storyExpiryMediaRepo struct {
	asset  *models.MediaAsset
	assets map[uuid.UUID]*models.MediaAsset
}

// Returns stored media fixtures for story expiry tests.
func (r *storyExpiryMediaRepo) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	if r.asset != nil && r.asset.ID == mediaAssetID {
		return r.asset, nil
	}
	if asset := r.assets[mediaAssetID]; asset != nil {
		return asset, nil
	}
	return nil, story_repo.ErrStoryNotFound
}

func (r *storyExpiryMediaRepo) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreateMediaAsset")
}

func (r *storyExpiryMediaRepo) CreatePostMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.ConfirmPostMediaUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePostMediaAsset")
}

func (r *storyExpiryMediaRepo) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePendingMediaAsset")
}

func (r *storyExpiryMediaRepo) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePendingStoryMediaAsset")
}

func (r *storyExpiryMediaRepo) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto media_dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	panic("unexpected call to MarkMediaAssetReady")
}

func (r *storyExpiryMediaRepo) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	panic("unexpected call to MarkMediaAssetFailed")
}

func (r *storyExpiryMediaRepo) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	panic("unexpected call to DeleteOrphanedMediaAssets")
}

var _ media_repo.MediaRepository = (*storyExpiryMediaRepo)(nil)
var _ persona_repo.PersonaRepository = (*storyHighlightPersonaRepo)(nil)
var _ follow_repo.FollowRepository = (*storyHighlightFollowRepo)(nil)
var _ story_repo.StoryRepository = (*storyHighlightTestRepo)(nil)
