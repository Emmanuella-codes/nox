package controllers

import (
	"github.com/emmanuella-codes/nox/crew/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *CrewController) CreateCrew(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}
	var dto dtos.CreateCrewDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.CreateCrewPipe(ctx.Context(), userID, eventID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *CrewController) JoinCrew(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.JoinCrewDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.JoinCrewPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *CrewController) ListMyEventCrews(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	eventID, err := uuid.Parse(ctx.Params("eventID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_event_id")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListMyEventCrewsPipe(ctx.Context(), userID, eventID, personaID, queryLimit(ctx, 20), queryOffset(ctx))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *CrewController) GetCrew(ctx *fiber.Ctx) error {
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
	res := c.pipe.GetCrewPipe(ctx.Context(), userID, crewID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *CrewController) LeaveCrew(ctx *fiber.Ctx) error {
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
	res := c.pipe.LeaveCrewPipe(ctx.Context(), userID, crewID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *CrewController) EndCrew(ctx *fiber.Ctx) error {
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
	res := c.pipe.EndCrewPipe(ctx.Context(), userID, crewID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
