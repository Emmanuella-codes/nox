package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *StoryController) GetStory(ctx *fiber.Ctx) error {
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}

	viewerPersonaID, err := optionalViewerPersonaID(ctx)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.GetStoryPipe(ctx.Context(), storyID, viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *StoryController) ListStoryItems(ctx *fiber.Ctx) error {
	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}

	res := c.pipe.ListStoryItemsPipe(ctx.Context(), storyID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
