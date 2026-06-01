package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	event_repo "github.com/emmanuella-codes/nox/repositories/event"
	follow_repo "github.com/emmanuella-codes/nox/repositories/follow"
	media_repo "github.com/emmanuella-codes/nox/repositories/media"
	persona_repo "github.com/emmanuella-codes/nox/repositories/persona"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type StoryPipe struct {
	storyRepo   story_repo.StoryRepository
	eventRepo   event_repo.EventRepository
	personaRepo persona_repo.PersonaRepository
	mediaRepo   media_repo.MediaRepository
	followRepo  follow_repo.FollowRepository
}

func NewStoryPipe(storyRepo story_repo.StoryRepository, eventRepo event_repo.EventRepository, personaRepo persona_repo.PersonaRepository, mediaRepo media_repo.MediaRepository, followRepo follow_repo.FollowRepository) *StoryPipe {
	return &StoryPipe{storyRepo: storyRepo, eventRepo: eventRepo, personaRepo: personaRepo, mediaRepo: mediaRepo, followRepo: followRepo}
}

func pipeInternalError[T any](err error, operation string) *shared.PipeRes[T] {
	return shared.PipeInternalError[T](err, "story", operation, messages.Internal_Error)
}

func validContributionMode(mode models.StoryContributionMode) bool {
	return mode == models.PublicStoryContributionMode || mode == models.FollowersStoryContributionMode
}

func validPostingMode(mode models.PostingMode) bool {
	return mode == models.PublicPostingMode || mode == models.AnonymousPostingMode
}

func validStoryVideo(asset *models.MediaAsset) bool {
	return asset.MediaKind == models.VideoMediaKind &&
		asset.ProcessingStatus == models.ReadyMediaStatus &&
		asset.DurationSeconds > 0 &&
		asset.DurationSeconds <= 300
}

func (p *StoryPipe) ownedPersona(ctx context.Context, userID uuid.UUID, personaID uuid.UUID) (*models.Persona, shared.PipeMessage) {
	persona, err := p.personaRepo.FindPersonaByID(ctx, personaID)
	if err != nil {
		if err == persona_repo.ErrPersonaNotFound {
			return nil, messages.Persona_Not_Found
		}
		return nil, messages.Internal_Error
	}
	if persona.UserID != userID {
		return nil, messages.Forbidden
	}
	return persona, ""
}

func (p *StoryPipe) mediaAsset(ctx context.Context, mediaAssetID uuid.UUID) (*models.MediaAsset, shared.PipeMessage) {
	asset, err := p.mediaRepo.FindMediaAssetByID(ctx, mediaAssetID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, messages.Media_Asset_Not_Found
		}
		return nil, messages.Internal_Error
	}
	return asset, ""
}

func (p *StoryPipe) canContribute(ctx context.Context, story *models.Story, contributorPersonaID uuid.UUID) (bool, error) {
	if story.OwnerPersonaID == contributorPersonaID {
		return true, nil
	}
	if story.ContributionMode == models.PublicStoryContributionMode {
		return true, nil
	}
	return p.followRepo.IsFollowing(ctx, contributorPersonaID, story.OwnerPersonaID)
}

func (p *StoryPipe) storyResponse(ctx context.Context, story *models.Story, viewerPersonaID *uuid.UUID, includeItems bool) (*StoryResponse, error) {
	owner, err := p.personaRepo.FindPersonaByID(ctx, story.OwnerPersonaID)
	if err != nil {
		return nil, err
	}
	canContribute := false
	if viewerPersonaID != nil {
		allowed, err := p.canContribute(ctx, story, *viewerPersonaID)
		if err != nil {
			return nil, err
		}
		canContribute = allowed
	}
	response := &StoryResponse{
		ID:                   uuidString(story.ID),
		EventID:              uuidString(story.EventID),
		Owner:                personaResponse(owner),
		Title:                story.Title,
		ContributionMode:     story.ContributionMode,
		TotalDurationSeconds: story.TotalDurationSeconds,
		CanContribute:        canContribute,
		Items:                []StoryItemResponse{},
		CreatedAt:            story.CreatedAt,
		UpdatedAt:            story.UpdatedAt,
	}
	if includeItems {
		items, err := p.storyItemResponses(ctx, story.ID)
		if err != nil {
			return nil, err
		}
		response.Items = items
	}
	return response, nil
}

func (p *StoryPipe) storyItemResponses(ctx context.Context, storyID uuid.UUID) ([]StoryItemResponse, error) {
	items, err := p.storyRepo.FindStoryItems(ctx, storyID)
	if err != nil {
		return nil, err
	}
	responses := make([]StoryItemResponse, 0, len(items))
	for _, item := range items {
		asset, err := p.mediaRepo.FindMediaAssetByID(ctx, item.MediaAssetID)
		if err != nil {
			return nil, err
		}
		response := StoryItemResponse{
			ID:              uuidString(item.ID),
			StoryID:         uuidString(item.StoryID),
			MediaAsset:      asset,
			PostingMode:     item.PostingMode,
			DurationSeconds: item.DurationSeconds,
			Position:        item.Position,
			CreatedAt:       item.CreatedAt,
		}
		if item.PostingMode == models.AnonymousPostingMode {
			response.AnonymousLabel = item.AnonymousLabel
		} else {
			persona, err := p.personaRepo.FindPersonaByID(ctx, item.ContributorPersonaID)
			if err != nil {
				return nil, err
			}
			contributor := personaResponse(persona)
			response.Contributor = &contributor
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 50 {
		return 50
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
