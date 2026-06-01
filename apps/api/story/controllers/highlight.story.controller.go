package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *StoryController) AddEventHighlightStory(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}

	var dto dtos.AddEventHighlightStoryDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.AddEventHighlightStoryPipe(ctx.Context(), userID, eventID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *StoryController) ListEventHighlightStories(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}

	res := c.pipe.ListEventHighlightStoriesPipe(ctx.Context(), eventID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *StoryController) RemoveEventHighlightStory(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	addedByPersonaID, err := requiredPersonaQuery(ctx, "added_by_persona_id")
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.RemoveEventHighlightStoryPipe(ctx.Context(), userID, eventID, storyID, addedByPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
