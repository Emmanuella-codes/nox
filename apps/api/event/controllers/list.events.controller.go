package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *EventController) ListEvents(ctx *fiber.Ctx) error {
	viewerPersonaID := ctx.Query("viewer_persona_id")
	if viewerPersonaID != "" {
		userID, ok := middleware.CurrentUserID(ctx)
		if !ok {
			return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
		}
		parsedViewerPersonaID, err := uuid.Parse(viewerPersonaID)
		if err != nil {
			return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
		}
		persona, message := c.pipe.FindViewerPersona(ctx.Context(), userID, parsedViewerPersonaID)
		if message != "" {
			return pipeError(ctx, pipeErrorStatus(message), message)
		}
		res := c.pipe.ListEventsForViewerPipe(ctx.Context(), persona.ID, queryLimit(ctx, 30))
		if !res.Success {
			return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
		}
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	res := c.pipe.ListEventsPipe(ctx.Context(), queryLimit(ctx, 30))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
