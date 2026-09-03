package pipes

import (
	"context"
	"time"

	event_dtos "github.com/emmanuella-codes/nox/event/dtos"
	media_dtos "github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/models"
	persona_dtos "github.com/emmanuella-codes/nox/persona/dtos"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/google/uuid"
	"testing"
)

// TestStoryCanBeHighlightedOnEventAndProfile verifies one story can live in both highlight containers.
func TestStoryCanBeHighlightedOnEventAndProfile(t *testing.T) {
	storyID := uuid.New()
	eventID := uuid.New()
	ownerUserID := uuid.New()
	ownerPersonaID := uuid.New()
	organizerUserID := uuid.New()
	organizerPersonaID := uuid.New()

	storyModel := &models.Story{
		ID:               storyID,
		EventID:          eventID,
		OwnerUserID:      ownerUserID,
		OwnerPersonaID:   ownerPersonaID,
		Title:            "aftermovie",
		ContributionMode: models.PublicStoryContributionMode,
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}
	pipe := NewStoryPipe(
		&storyHighlightTestRepo{story: storyModel},
		&storyHighlightEventRepo{event: &models.Event{ID: eventID, OrganizerID: organizerPersonaID}},
		&storyHighlightPersonaRepo{personas: map[uuid.UUID]*models.Persona{
			ownerPersonaID:     {ID: ownerPersonaID, UserID: ownerUserID, Handle: "owner", DisplayName: "Owner", AvatarURL: "owner.png", Category: models.PatronPersonaCategory},
			organizerPersonaID: {ID: organizerPersonaID, UserID: organizerUserID, Handle: "org", DisplayName: "Organizer", AvatarURL: "org.png", Category: models.OrganizerPersonaCategory},
		}},
		&storyHighlightMediaRepo{},
		&storyHighlightFollowRepo{},
	)

	eventRes := pipe.AddEventHighlightStoryPipe(context.Background(), organizerUserID, eventID, dtos.AddEventHighlightStoryDTO{
		StoryID:          storyID,
		AddedByPersonaID: organizerPersonaID,
	})
	if !eventRes.Success {
		t.Fatalf("expected event highlight success, got %s", eventRes.Message)
	}

	profileRes := pipe.AddProfileStoryHighlightPipe(context.Background(), ownerUserID, storyID, dtos.ProfileStoryHighlightDTO{
		PersonaID: ownerPersonaID,
	})
	if !profileRes.Success {
		t.Fatalf("expected profile highlight success, got %s", profileRes.Message)
	}

	repo := pipe.storyRepo.(*storyHighlightTestRepo)
	if len(repo.eventHighlights) != 1 {
		t.Fatalf("expected 1 event highlight, got %d", len(repo.eventHighlights))
	}
	if len(repo.profileHighlights) != 1 {
		t.Fatalf("expected 1 profile highlight, got %d", len(repo.profileHighlights))
	}
	if repo.eventHighlights[0].StoryID != storyID {
		t.Fatalf("expected event highlight story id %s, got %s", storyID, repo.eventHighlights[0].StoryID)
	}
	if repo.profileHighlights[0].StoryID != storyID {
		t.Fatalf("expected profile highlight story id %s, got %s", storyID, repo.profileHighlights[0].StoryID)
	}
}

// storyHighlightTestRepo stores highlight mutations in memory for story pipe tests.
type storyHighlightTestRepo struct {
	story             *models.Story
	items             []*models.StoryItem
	eventHighlights   []*models.EventHighlightStory
	profileHighlights []*models.ProfileStoryHighlight
	acceptsAdditions  bool
}

// CreateStory is unused in this test stub.
func (r *storyHighlightTestRepo) CreateStory(ctx context.Context, ownerUserID uuid.UUID, dto dtos.CreateStoryDTO) (*models.Story, error) {
	panic("unexpected call to CreateStory")
}

// FindStoryByID returns the active story fixture.
func (r *storyHighlightTestRepo) FindStoryByID(ctx context.Context, storyID uuid.UUID) (*models.Story, error) {
	if r.story == nil || r.story.ID != storyID {
		return nil, story_repo.ErrStoryNotFound
	}
	return r.story, nil
}

// FindStoryByIDAny returns the story fixture regardless of expiry.
func (r *storyHighlightTestRepo) FindStoryByIDAny(ctx context.Context, storyID uuid.UUID) (*models.Story, error) {
	return r.FindStoryByID(ctx, storyID)
}

// FindStoriesByEventID is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoriesByEventID(ctx context.Context, eventID uuid.UUID, limit int, offset int) ([]*models.Story, error) {
	panic("unexpected call to FindStoriesByEventID")
}

// FindStoriesByOwnerPersonaID is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoriesByOwnerPersonaID(ctx context.Context, personaID uuid.UUID, limit int, offset int) ([]*models.Story, error) {
	panic("unexpected call to FindStoriesByOwnerPersonaID")
}

// Returns the configured add-window decision.
func (r *storyHighlightTestRepo) StoryAcceptsAdditions(ctx context.Context, storyID uuid.UUID) (bool, error) {
	return r.acceptsAdditions, nil
}

// DeleteStory is unused in this test stub.
func (r *storyHighlightTestRepo) DeleteStory(ctx context.Context, storyID uuid.UUID) error {
	panic("unexpected call to DeleteStory")
}

// AddStoryItem is unused in this test stub.
func (r *storyHighlightTestRepo) AddStoryItem(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, durationSeconds int, anonymousLabel string, dto dtos.AddStoryItemDTO) (*models.StoryItem, error) {
	panic("unexpected call to AddStoryItem")
}

// FindStoryItemByID is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItemByID(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error) {
	panic("unexpected call to FindStoryItemByID")
}

// FindStoryItems is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItems(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error) {
	items := make([]*models.StoryItem, 0, len(r.items))
	now := time.Now()
	for _, item := range r.items {
		if item.StoryID == storyID && item.ExpiresAt.After(now) {
			items = append(items, item)
		}
	}
	return items, nil
}

// Returns all stored story items for preserved highlight reads.
func (r *storyHighlightTestRepo) FindStoryItemsAny(ctx context.Context, storyID uuid.UUID) ([]*models.StoryItem, error) {
	items := make([]*models.StoryItem, 0, len(r.items))
	for _, item := range r.items {
		if item.StoryID == storyID {
			items = append(items, item)
		}
	}
	return items, nil
}

// DeleteStoryItem is unused in this test stub.
func (r *storyHighlightTestRepo) DeleteStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) (*models.StoryItem, error) {
	panic("unexpected call to DeleteStoryItem")
}

// ReorderStoryItem is unused in this test stub.
func (r *storyHighlightTestRepo) ReorderStoryItem(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, position int) (*models.StoryItem, error) {
	panic("unexpected call to ReorderStoryItem")
}

// CreateStoryContributionRequest is unused in this test stub.
func (r *storyHighlightTestRepo) CreateStoryContributionRequest(ctx context.Context, storyID uuid.UUID, contributorUserID uuid.UUID, dto dtos.CreateStoryContributionRequestDTO) (*models.StoryContributionRequest, error) {
	panic("unexpected call to CreateStoryContributionRequest")
}

// FindStoryContributionRequestByID is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryContributionRequestByID(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID) (*models.StoryContributionRequest, error) {
	panic("unexpected call to FindStoryContributionRequestByID")
}

// FindStoryContributionRequests is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryContributionRequests(ctx context.Context, storyID uuid.UUID, status *models.StoryContributionRequestStatus, limit int, offset int) ([]*models.StoryContributionRequest, error) {
	panic("unexpected call to FindStoryContributionRequests")
}

// AcceptStoryContributionRequest is unused in this test stub.
func (r *storyHighlightTestRepo) AcceptStoryContributionRequest(ctx context.Context, story *models.Story, requestID uuid.UUID, reviewerPersonaID uuid.UUID, durationSeconds int) (*models.StoryContributionRequest, *models.StoryItem, error) {
	panic("unexpected call to AcceptStoryContributionRequest")
}

// RejectStoryContributionRequest is unused in this test stub.
func (r *storyHighlightTestRepo) RejectStoryContributionRequest(ctx context.Context, storyID uuid.UUID, requestID uuid.UUID, reviewerPersonaID uuid.UUID) (*models.StoryContributionRequest, error) {
	panic("unexpected call to RejectStoryContributionRequest")
}

// AddEventHighlightStory stores one event highlight row.
func (r *storyHighlightTestRepo) AddEventHighlightStory(ctx context.Context, eventID uuid.UUID, dto dtos.AddEventHighlightStoryDTO) (*models.EventHighlightStory, error) {
	highlight := &models.EventHighlightStory{
		ID:               uuid.New(),
		EventID:          eventID,
		StoryID:          dto.StoryID,
		AddedByPersonaID: dto.AddedByPersonaID,
		Position:         len(r.eventHighlights) + 1,
		CreatedAt:        time.Now(),
	}
	r.eventHighlights = append(r.eventHighlights, highlight)
	return highlight, nil
}

// FindEventHighlightStories is unused in this test stub.
func (r *storyHighlightTestRepo) FindEventHighlightStories(ctx context.Context, eventID uuid.UUID) ([]*models.EventHighlightStory, error) {
	return r.eventHighlights, nil
}

// ReorderEventHighlightStory is unused in this test stub.
func (r *storyHighlightTestRepo) ReorderEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID, position int) (*models.EventHighlightStory, error) {
	panic("unexpected call to ReorderEventHighlightStory")
}

// RemoveEventHighlightStory is unused in this test stub.
func (r *storyHighlightTestRepo) RemoveEventHighlightStory(ctx context.Context, eventID uuid.UUID, storyID uuid.UUID) error {
	panic("unexpected call to RemoveEventHighlightStory")
}

// UpsertStoryItemView is unused in this test stub.
func (r *storyHighlightTestRepo) UpsertStoryItemView(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, viewerUserID uuid.UUID, viewerPersonaID uuid.UUID) (*models.StoryItemView, bool, error) {
	panic("unexpected call to UpsertStoryItemView")
}

// FindStoryItemViewCounts is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItemViewCounts(ctx context.Context, storyID uuid.UUID) (map[uuid.UUID]int, error) {
	return map[uuid.UUID]int{}, nil
}

// FindStoryItemViewerPersonaIDs is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItemViewerPersonaIDs(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

// FindViewedStoryItemIDs is unused in this test stub.
func (r *storyHighlightTestRepo) FindViewedStoryItemIDs(ctx context.Context, storyID uuid.UUID, viewerPersonaID uuid.UUID) ([]uuid.UUID, error) {
	return []uuid.UUID{}, nil
}

// UpsertStoryItemReaction is unused in this test stub.
func (r *storyHighlightTestRepo) UpsertStoryItemReaction(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, reactorUserID uuid.UUID, reactorPersonaID uuid.UUID, reactionType models.StoryReactionType) (*models.StoryItemReaction, error) {
	panic("unexpected call to UpsertStoryItemReaction")
}

// DeleteStoryItemReaction is unused in this test stub.
func (r *storyHighlightTestRepo) DeleteStoryItemReaction(ctx context.Context, storyID uuid.UUID, itemID uuid.UUID, reactorPersonaID uuid.UUID) error {
	panic("unexpected call to DeleteStoryItemReaction")
}

// FindStoryItemReactionCounts is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItemReactionCounts(ctx context.Context, storyID uuid.UUID) (map[uuid.UUID]map[models.StoryReactionType]int, error) {
	return map[uuid.UUID]map[models.StoryReactionType]int{}, nil
}

// FindStoryItemReactionsByPersona is unused in this test stub.
func (r *storyHighlightTestRepo) FindStoryItemReactionsByPersona(ctx context.Context, storyID uuid.UUID, personaID uuid.UUID) (map[uuid.UUID]models.StoryReactionType, error) {
	return map[uuid.UUID]models.StoryReactionType{}, nil
}

// AddProfileStoryHighlight stores one profile highlight row.
func (r *storyHighlightTestRepo) AddProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) (*models.ProfileStoryHighlight, error) {
	highlight := &models.ProfileStoryHighlight{
		ID:             uuid.New(),
		OwnerPersonaID: ownerPersonaID,
		StoryID:        storyID,
		Position:       len(r.profileHighlights) + 1,
		CreatedAt:      time.Now(),
	}
	r.profileHighlights = append(r.profileHighlights, highlight)
	return highlight, nil
}

// FindProfileStoryHighlights is unused in this test stub.
func (r *storyHighlightTestRepo) FindProfileStoryHighlights(ctx context.Context, ownerPersonaID uuid.UUID) ([]*models.ProfileStoryHighlight, error) {
	return r.profileHighlights, nil
}

// RemoveProfileStoryHighlight is unused in this test stub.
func (r *storyHighlightTestRepo) RemoveProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) error {
	panic("unexpected call to RemoveProfileStoryHighlight")
}

// HasProfileStoryHighlight is unused in this test stub.
func (r *storyHighlightTestRepo) HasProfileStoryHighlight(ctx context.Context, ownerPersonaID uuid.UUID, storyID uuid.UUID) (bool, error) {
	return false, nil
}

// No-ops the auto-rejection batch for this in-memory stub.
func (r *storyHighlightTestRepo) RejectPendingContributionRequestsForClosedStories(ctx context.Context, limit int) (int64, error) {
	return 0, nil
}

// No-ops retained item cleanup for this in-memory stub.
func (r *storyHighlightTestRepo) DeleteExpiredNonHighlightedStoryItems(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return 0, nil
}

// No-ops empty-story cleanup for this in-memory stub.
func (r *storyHighlightTestRepo) DeleteRetainedEmptyStories(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	return 0, nil
}

// storyHighlightEventRepo returns one event fixture.
type storyHighlightEventRepo struct {
	event *models.Event
}

// CreateEvent is unused in this test stub.
func (r *storyHighlightEventRepo) CreateEvent(ctx context.Context, dto event_dtos.CreateEventDTO) (*models.Event, error) {
	panic("unexpected call to CreateEvent")
}

// FindEventByID returns one event fixture.
func (r *storyHighlightEventRepo) FindEventByID(ctx context.Context, eventID uuid.UUID) (*models.Event, error) {
	if r.event == nil || r.event.ID != eventID {
		return nil, event_repo.ErrEventNotFound
	}
	return r.event, nil
}

// FindEvents is unused in this test stub.
func (r *storyHighlightEventRepo) FindEvents(ctx context.Context, limit int) ([]*models.Event, error) {
	panic("unexpected call to FindEvents")
}

// storyHighlightPersonaRepo returns persona fixtures by id.
type storyHighlightPersonaRepo struct {
	personas map[uuid.UUID]*models.Persona
}

// CreatePersona is unused in this test stub.
func (r *storyHighlightPersonaRepo) CreatePersona(ctx context.Context, userID uuid.UUID, dto persona_dtos.CreatePersonaDTO) (*models.Persona, error) {
	panic("unexpected call to CreatePersona")
}

// FindPersonaByID returns one persona fixture.
func (r *storyHighlightPersonaRepo) FindPersonaByID(ctx context.Context, personaID uuid.UUID) (*models.Persona, error) {
	persona := r.personas[personaID]
	if persona == nil {
		return nil, persona_repo.ErrPersonaNotFound
	}
	return persona, nil
}

// FindPersonasByUserID is unused in this test stub.
func (r *storyHighlightPersonaRepo) FindPersonasByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Persona, error) {
	panic("unexpected call to FindPersonasByUserID")
}

// FindPersonaByHandle is unused in this test stub.
func (r *storyHighlightPersonaRepo) FindPersonaByHandle(ctx context.Context, handle string) (*models.Persona, error) {
	panic("unexpected call to FindPersonaByHandle")
}

// UpdatePersona is unused in this test stub.
func (r *storyHighlightPersonaRepo) UpdatePersona(ctx context.Context, personaID uuid.UUID, dto persona_dtos.UpdatePersonaDTO) (*models.Persona, error) {
	panic("unexpected call to UpdatePersona")
}

// storyHighlightMediaRepo satisfies the media dependency for this test.
type storyHighlightMediaRepo struct{}

// CreateMediaAsset is unused in this test stub.
func (r *storyHighlightMediaRepo) CreateMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.CreateMediaAssetDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreateMediaAsset")
}

// CreatePostMediaAsset is unused in this test stub.
func (r *storyHighlightMediaRepo) CreatePostMediaAsset(ctx context.Context, ownerUserID uuid.UUID, dto media_dtos.ConfirmPostMediaUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePostMediaAsset")
}

// CreatePendingMediaAsset is unused in this test stub.
func (r *storyHighlightMediaRepo) CreatePendingMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateSetVideoUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePendingMediaAsset")
}

// CreatePendingStoryMediaAsset is unused in this test stub.
func (r *storyHighlightMediaRepo) CreatePendingStoryMediaAsset(ctx context.Context, ownerUserID uuid.UUID, storageKey string, playbackURL string, dto media_dtos.InitiateStoryVideoUploadDTO) (*models.MediaAsset, error) {
	panic("unexpected call to CreatePendingStoryMediaAsset")
}

// FindMediaAssetByID is unused in this test stub.
func (r *storyHighlightMediaRepo) FindMediaAssetByID(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	panic("unexpected call to FindMediaAssetByID")
}

// MarkMediaAssetReady is unused in this test stub.
func (r *storyHighlightMediaRepo) MarkMediaAssetReady(ctx context.Context, mediaAssetID uuid.UUID, dto media_dtos.CompleteMediaProcessingDTO) (*models.MediaAsset, error) {
	panic("unexpected call to MarkMediaAssetReady")
}

// MarkMediaAssetFailed is unused in this test stub.
func (r *storyHighlightMediaRepo) MarkMediaAssetFailed(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, error) {
	panic("unexpected call to MarkMediaAssetFailed")
}

// DeleteOrphanedMediaAssets is unused in this test stub.
func (r *storyHighlightMediaRepo) DeleteOrphanedMediaAssets(ctx context.Context, olderThan time.Time, limit int) (int64, error) {
	panic("unexpected call to DeleteOrphanedMediaAssets")
}

// storyHighlightFollowRepo satisfies the follow dependency for this test.
type storyHighlightFollowRepo struct{}

// Follow is unused in this test stub.
func (r *storyHighlightFollowRepo) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	panic("unexpected call to Follow")
}

// Unfollow is unused in this test stub.
func (r *storyHighlightFollowRepo) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	panic("unexpected call to Unfollow")
}

// IsFollowing is unused in this test stub.
func (r *storyHighlightFollowRepo) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return false, nil
}

// FindFollowingIDs is unused in this test stub.
func (r *storyHighlightFollowRepo) FindFollowingIDs(ctx context.Context, followerID uuid.UUID, followingIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	return map[uuid.UUID]bool{}, nil
}

// FindFollowers is unused in this test stub.
func (r *storyHighlightFollowRepo) FindFollowers(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) ([]*models.Persona, error) {
	panic("unexpected call to FindFollowers")
}

// FindFollowing is unused in this test stub.
func (r *storyHighlightFollowRepo) FindFollowing(ctx context.Context, personaID uuid.UUID, options follow_repo.ListOptions) ([]*models.Persona, error) {
	panic("unexpected call to FindFollowing")
}

var _ story_repo.StoryRepository = (*storyHighlightTestRepo)(nil)
var _ event_repo.EventRepository = (*storyHighlightEventRepo)(nil)
var _ persona_repo.PersonaRepository = (*storyHighlightPersonaRepo)(nil)
var _ media_repo.MediaRepository = (*storyHighlightMediaRepo)(nil)
var _ follow_repo.FollowRepository = (*storyHighlightFollowRepo)(nil)
