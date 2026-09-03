package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *EventController) GetEvent(ctx *fiber.Ctx) error {
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}

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
		res := c.pipe.GetEventForViewerPipe(ctx.Context(), eventID, persona.ID)
		if !res.Success {
			return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
		}
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	res := c.pipe.GetEventPipe(ctx.Context(), eventID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
