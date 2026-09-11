package controllers

import (
	"strconv"
	"strings"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/admin/messages"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListCrews(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	limit, _ := strconv.Atoi(ctx.Query("limit", "50"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListCrewsPipe(requestContext(ctx), adminID, strings.TrimSpace(ctx.Query("status")), limit, offset)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetCrew(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	res := ac.pipe.GetCrewPipe(requestContext(ctx), adminID, crewID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) EndCrew(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	var dto dtos.CrewSafetyActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.EndCrewSafetyPipe(requestContext(ctx), adminID, crewID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) DisableCrewSharing(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	var dto dtos.CrewSafetyActionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.DisableCrewSharingPipe(requestContext(ctx), adminID, crewID, dto)
	if res.Success {
		return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListCrewLocations(ctx *fiber.Ctx) error {
	adminID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	crewID, err := uuid.Parse(ctx.Params("crewID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	dto := dtos.CrewLocationAccessDTO{Reason: ctx.Query("reason"), ReportID: ctx.Query("report_id")}
	if strings.TrimSpace(dto.Reason) == "" {
		return pipeError(ctx, fiber.StatusBadRequest, messages.Invalid_Payload)
	}
	res := ac.pipe.ListCrewLocationsPipe(requestContext(ctx), adminID, crewID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
