package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/search/pipes"
	"github.com/emmanuella-codes/nox/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *SearchController) Search(ctx *fiber.Ctx) error {
	var res *shared.PipeRes[pipes.SearchResponse]
	options := pipes.SearchOptions{Limit: queryLimit(ctx, 10), Offset: queryOffset(ctx)}
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
		res = c.pipe.SearchForViewer(ctx.Context(), ctx.Query("q"), options, persona.ID)
	} else {
		res = c.pipe.Search(ctx.Context(), ctx.Query("q"), options)
	}

	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
