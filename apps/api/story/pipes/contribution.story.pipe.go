package pipes

import (
	"context"

	"github.com/emmanuella-codes/nox/models"
	story_repo "github.com/emmanuella-codes/nox/repositories/story"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/emmanuella-codes/nox/story/messages"
	"github.com/google/uuid"
)

func (p *StoryPipe) CreateStoryContributionRequestPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, dto dtos.CreateStoryContributionRequestDTO) *shared.PipeRes[StoryContributionRequestResponse] {
	contributor, message := p.ownedPersona(ctx, userID, dto.ContributorPersonaID)
	if message != "" {
		return shared.PipeError[StoryContributionRequestResponse](message)
	}
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_story")
	}
	if story.OwnerPersonaID == contributor.ID {
		return shared.PipeError[StoryContributionRequestResponse](messages.Forbidden)
	}
	allowed, err := p.canContribute(ctx, story, contributor.ID)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_permission")
	}
	if !allowed {
		return shared.PipeError[StoryContributionRequestResponse](messages.Forbidden)
	}
	acceptsAdditions, err := p.storyRepo.StoryAcceptsAdditions(ctx, storyID)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_window")
	}
	if !acceptsAdditions {
		return shared.PipeError[StoryContributionRequestResponse](messages.Story_Closed_For_Additions)
	}
	asset, message := p.mediaAsset(ctx, dto.MediaAssetID)
	if message != "" {
		return shared.PipeError[StoryContributionRequestResponse](message)
	}
	if asset.OwnerUserID != userID || asset.OwnerPersonaID != contributor.ID || !validStoryMedia(asset) {
		return shared.PipeError[StoryContributionRequestResponse](messages.Invalid_Story)
	}
	request, err := p.storyRepo.CreateStoryContributionRequest(ctx, storyID, userID, dto)
	if err != nil {
		if err == story_repo.ErrStoryMediaInUse {
			return shared.PipeError[StoryContributionRequestResponse](messages.Media_Asset_In_Use)
		}
		if err == story_repo.ErrStoryContributionRequestPending {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Already_Pending)
		}
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_create")
	}
	if err := p.notifyStoryContributionRequest(ctx, story, contributor, request); err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_notify")
	}
	response, err := p.storyContributionRequestResponse(ctx, request)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_response")
	}
	return shared.PipeSuccess(messages.Story_Contribution_Request_Created, response)
}

func (p *StoryPipe) ListStoryContributionRequestsPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, statusValue string, limit int, offset int) *shared.PipeRes[[]StoryContributionRequestResponse] {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return shared.PipeError[[]StoryContributionRequestResponse](messages.Story_Not_Found)
		}
		return pipeInternalError[[]StoryContributionRequestResponse](err, "story.contribution_request_list_story")
	}
	if story.OwnerUserID != userID {
		return shared.PipeError[[]StoryContributionRequestResponse](messages.Forbidden)
	}
	status, message := parseStoryContributionRequestStatus(statusValue)
	if message != "" {
		return shared.PipeError[[]StoryContributionRequestResponse](message)
	}
	requests, err := p.storyRepo.FindStoryContributionRequests(ctx, storyID, status, normalizeLimit(limit), normalizeOffset(offset))
	if err != nil {
		return pipeInternalError[[]StoryContributionRequestResponse](err, "story.contribution_request_list")
	}
	responses, err := p.storyContributionRequestResponses(ctx, requests)
	if err != nil {
		return pipeInternalError[[]StoryContributionRequestResponse](err, "story.contribution_request_list_response")
	}
	return shared.PipeSuccess(messages.Story_Contribution_Requests_Listed, &responses)
}

func (p *StoryPipe) AcceptStoryContributionRequestPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, requestID uuid.UUID) *shared.PipeRes[StoryContributionRequestResponse] {
	request, story, contributor, response := p.loadContributionReviewContext(ctx, userID, storyID, requestID)
	if response != nil {
		return response
	}
	asset, message := p.mediaAsset(ctx, request.MediaAssetID)
	if message != "" {
		return shared.PipeError[StoryContributionRequestResponse](message)
	}
	acceptsAdditions, err := p.storyRepo.StoryAcceptsAdditions(ctx, story.ID)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_accept_window")
	}
	if !acceptsAdditions {
		return shared.PipeError[StoryContributionRequestResponse](messages.Story_Closed_For_Additions)
	}
	durationSeconds := storyItemDuration(asset)
	if !validStoryMedia(asset) || durationSeconds <= 0 {
		return shared.PipeError[StoryContributionRequestResponse](messages.Invalid_Story)
	}
	request, _, err = p.storyRepo.AcceptStoryContributionRequest(ctx, story, requestID, story.OwnerPersonaID, durationSeconds)
	if err != nil {
		if err == story_repo.ErrStoryDurationLimitExceeded {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Duration_Limit_Exceeded)
		}
		if err == story_repo.ErrStoryClosedForAdditions {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Closed_For_Additions)
		}
		if err == story_repo.ErrStoryMediaInUse {
			return shared.PipeError[StoryContributionRequestResponse](messages.Media_Asset_In_Use)
		}
		if err == story_repo.ErrStoryContributionRequestReviewed {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Already_Reviewed)
		}
		if err == story_repo.ErrStoryContributionRequestNotFound {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Not_Found)
		}
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_accept")
	}
	if err := p.notifyStoryContributionDecision(ctx, story, contributor, models.StoryContributionAcceptedNotificationType, request); err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_accept_notify")
	}
	result, err := p.storyContributionRequestResponse(ctx, request)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_accept_response")
	}
	return shared.PipeSuccess(messages.Story_Contribution_Request_Accepted, result)
}

func (p *StoryPipe) RejectStoryContributionRequestPipe(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, requestID uuid.UUID) *shared.PipeRes[StoryContributionRequestResponse] {
	request, story, contributor, response := p.loadContributionReviewContext(ctx, userID, storyID, requestID)
	if response != nil {
		return response
	}
	request, err := p.storyRepo.RejectStoryContributionRequest(ctx, storyID, requestID, story.OwnerPersonaID)
	if err != nil {
		if err == story_repo.ErrStoryContributionRequestReviewed {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Already_Reviewed)
		}
		if err == story_repo.ErrStoryContributionRequestNotFound {
			return shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Not_Found)
		}
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_reject")
	}
	if err := p.notifyStoryContributionDecision(ctx, story, contributor, models.StoryContributionRejectedNotificationType, request); err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_reject_notify")
	}
	result, err := p.storyContributionRequestResponse(ctx, request)
	if err != nil {
		return pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_reject_response")
	}
	return shared.PipeSuccess(messages.Story_Contribution_Request_Rejected, result)
}

func (p *StoryPipe) loadContributionReviewContext(ctx context.Context, userID uuid.UUID, storyID uuid.UUID, requestID uuid.UUID) (*models.StoryContributionRequest, *models.Story, *models.Persona, *shared.PipeRes[StoryContributionRequestResponse]) {
	story, err := p.storyRepo.FindStoryByID(ctx, storyID)
	if err != nil {
		if err == story_repo.ErrStoryNotFound {
			return nil, nil, nil, shared.PipeError[StoryContributionRequestResponse](messages.Story_Not_Found)
		}
		return nil, nil, nil, pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_review_story")
	}
	if story.OwnerUserID != userID {
		return nil, nil, nil, shared.PipeError[StoryContributionRequestResponse](messages.Forbidden)
	}
	request, err := p.storyRepo.FindStoryContributionRequestByID(ctx, storyID, requestID)
	if err != nil {
		if err == story_repo.ErrStoryContributionRequestNotFound {
			return nil, nil, nil, shared.PipeError[StoryContributionRequestResponse](messages.Story_Contribution_Request_Not_Found)
		}
		return nil, nil, nil, pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_review_find")
	}
	contributor, err := p.personaRepo.FindPersonaByID(ctx, request.ContributorPersonaID)
	if err != nil {
		return nil, nil, nil, pipeInternalError[StoryContributionRequestResponse](err, "story.contribution_request_review_contributor")
	}
	return request, story, contributor, nil
}

func (p *StoryPipe) storyContributionRequestResponse(ctx context.Context, request *models.StoryContributionRequest) (*StoryContributionRequestResponse, error) {
	asset, err := p.mediaRepo.FindMediaAssetByID(ctx, request.MediaAssetID)
	if err != nil {
		return nil, err
	}
	contributor, err := p.personaRepo.FindPersonaByID(ctx, request.ContributorPersonaID)
	if err != nil {
		return nil, err
	}
	var reviewedByPersonaID *string
	var storyItemID *string
	if request.ReviewedByPersonaID != nil {
		value := request.ReviewedByPersonaID.String()
		reviewedByPersonaID = &value
	}
	if request.StoryItemID != nil {
		value := request.StoryItemID.String()
		storyItemID = &value
	}
	return &StoryContributionRequestResponse{
		ID:                  request.ID.String(),
		StoryID:             request.StoryID.String(),
		MediaAsset:          asset,
		Contributor:         personaResponse(contributor),
		Status:              request.Status,
		ReviewedByPersonaID: reviewedByPersonaID,
		StoryItemID:         storyItemID,
		CreatedAt:           request.CreatedAt,
		ReviewedAt:          request.ReviewedAt,
	}, nil
}

func (p *StoryPipe) storyContributionRequestResponses(ctx context.Context, requests []*models.StoryContributionRequest) ([]StoryContributionRequestResponse, error) {
	responses := make([]StoryContributionRequestResponse, 0, len(requests))
	for _, request := range requests {
		response, err := p.storyContributionRequestResponse(ctx, request)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}
	return responses, nil
}

func parseStoryContributionRequestStatus(value string) (*models.StoryContributionRequestStatus, shared.PipeMessage) {
	if value == "" {
		return nil, ""
	}
	status := models.StoryContributionRequestStatus(value)
	switch status {
	case models.PendingStoryContributionRequestStatus, models.AcceptedStoryContributionRequestStatus, models.RejectedStoryContributionRequestStatus:
		return &status, ""
	default:
		return nil, messages.Invalid_Story
	}
}
