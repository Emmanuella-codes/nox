package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ModerateContent(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	entityID, err := uuid.Parse(ctx.Params("entityID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.ModerateContentDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.ModerateContentPipe(requestContext(ctx), adminUserID, models.ModerationEntityType(ctx.Params("entityType")), entityID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListModerationActions(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListModerationActionsParams{}
	if value := ctx.Query("entity_type"); value != "" {
		entityType := models.ModerationEntityType(value)
		params.EntityType = &entityType
	}
	if value := ctx.Query("status"); value != "" {
		status := models.ModerationStatus(value)
		params.Status = &status
	}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListModerationActionsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}
