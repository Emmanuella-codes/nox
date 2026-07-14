package controllers

import (
	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *CrewController) UpdateLocationSharing(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_crew_id")
	}
	var dto dtos.UpdateSharingDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateLocationSharingPipe(ctx.Context(), userID, crewID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *CrewController) UpdateLocation(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_crew_id")
	}
	var dto dtos.UpdateLocationDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateLocationPipe(ctx.Context(), userID, crewID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *CrewController) ListLocations(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_crew_id")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListLocationsPipe(ctx.Context(), userID, crewID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
