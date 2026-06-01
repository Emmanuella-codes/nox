package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *SetController) ListSets(ctx *fiber.Ctx) error {
	viewerPersonaID, err := optionalViewerPersonaID(ctx)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.ListSetsPipe(ctx.Context(), queryLimit(ctx, 20), queryOffset(ctx), ctx.Query("genre"), ctx.Query("sort"), viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *SetController) ListPersonaSets(ctx *fiber.Ctx) error {
	personaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	viewerPersonaID, err := optionalViewerPersonaID(ctx)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	res := c.pipe.ListPersonaSetsPipe(ctx.Context(), personaID, queryLimit(ctx, 20), queryOffset(ctx), viewerPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
