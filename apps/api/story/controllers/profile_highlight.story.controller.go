package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AddProfileStoryHighlight adds one story to the current persona's profile highlights.
func (c *StoryController) AddProfileStoryHighlight(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	var dto dtos.ProfileStoryHighlightDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.AddProfileStoryHighlightPipe(ctx.Context(), userID, storyID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

// ListProfileStoryHighlights lists one persona's saved story highlights.
func (c *StoryController) ListProfileStoryHighlights(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	viewerUserID, viewerPersonaID, err := optionalViewerContext(ctx)
	if err != nil {
		if err == fiber.ErrUnauthorized {
			return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
		}
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListProfileStoryHighlightsPipe(ctx.Context(), personaID, viewerUserID, viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// RemoveProfileStoryHighlight removes one story from the current persona's profile highlights.
func (c *StoryController) RemoveProfileStoryHighlight(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	personaID, err := requiredPersonaQuery(ctx, "persona_id")
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.RemoveProfileStoryHighlightPipe(ctx.Context(), userID, storyID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
