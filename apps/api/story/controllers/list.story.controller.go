package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *StoryController) ListEventStories(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}

	viewerUserID, viewerPersonaID, err := optionalViewerContext(ctx)
	if err != nil {
		if err == fiber.ErrUnauthorized {
			return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
		}
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.ListEventStoriesPipe(ctx.Context(), eventID, queryLimit(ctx, 20), queryOffset(ctx), viewerUserID, viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *StoryController) ListPersonaStories(ctx *fiber.Ctx) error {
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

	res := c.pipe.ListPersonaStoriesPipe(ctx.Context(), personaID, queryLimit(ctx, 20), queryOffset(ctx), viewerUserID, viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
