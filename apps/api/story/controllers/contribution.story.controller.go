package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateStoryContributionRequest submits one contribution request for a story owner to review.
func (c *StoryController) CreateStoryContributionRequest(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	var dto dtos.CreateStoryContributionRequestDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.CreateStoryContributionRequestPipe(ctx.Context(), userID, storyID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

// ListStoryContributionRequests lists contribution requests for one story owner.
func (c *StoryController) ListStoryContributionRequests(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	res := c.pipe.ListStoryContributionRequestsPipe(ctx.Context(), userID, storyID, ctx.Query("status"), queryLimit(ctx, 20), queryOffset(ctx))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// AcceptStoryContributionRequest accepts one pending contribution request and creates the story item.
func (c *StoryController) AcceptStoryContributionRequest(ctx *fiber.Ctx) error {
	return c.reviewStoryContributionRequest(ctx, true)
}

// RejectStoryContributionRequest rejects one pending contribution request.
func (c *StoryController) RejectStoryContributionRequest(ctx *fiber.Ctx) error {
	return c.reviewStoryContributionRequest(ctx, false)
}

// reviewStoryContributionRequest dispatches one accept or reject action for a contribution request.
func (c *StoryController) reviewStoryContributionRequest(ctx *fiber.Ctx, accept bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	requestID, err := uuid.Parse(ctx.Params("requestID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_contribution_request_id")
	}
	if accept {
		pipeRes := c.pipe.AcceptStoryContributionRequestPipe(ctx.Context(), userID, storyID, requestID)
		if !pipeRes.Success {
			return pipeError(ctx, pipeErrorStatus(pipeRes.Message), pipeRes.Message)
		}
		return pipeSuccess(ctx, fiber.StatusOK, pipeRes.Message, pipeRes.Data)
	}
	pipeRes := c.pipe.RejectStoryContributionRequestPipe(ctx.Context(), userID, storyID, requestID)
	if !pipeRes.Success {
		return pipeError(ctx, pipeErrorStatus(pipeRes.Message), pipeRes.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, pipeRes.Message, pipeRes.Data)
}
