package controllers

import (
	"strconv"
	"time"

	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListContent(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	entityType := contentEntity(ctx)
	if entityType == "" {
		entityType = models.ModerationEntityType(ctx.Query("entity_type"))
	}
	params, err := contentParams(ctx, entityType)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.ListContentPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListContentType(entityType models.ModerationEntityType) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("admin_content_entity", entityType)
		return ac.ListContent(ctx)
	}
}

func (ac *AdminController) GetContent(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	entityID, err := uuid.Parse(ctx.Params("entityID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.GetContentPipe(requestContext(ctx), adminUserID, contentEntity(ctx), entityID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetContentType(entityType models.ModerationEntityType) func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("admin_content_entity", entityType)
		return ac.GetContent(ctx)
	}
}

func contentEntity(ctx *fiber.Ctx) models.ModerationEntityType {
	if value, ok := ctx.Locals("admin_content_entity").(models.ModerationEntityType); ok {
		return value
	}
	return models.ModerationEntityType(ctx.Params("entityType"))
}

func contentParams(ctx *fiber.Ctx, entityType models.ModerationEntityType) (adminrepo.ListAdminContentParams, error) {
	params := adminrepo.ListAdminContentParams{EntityType: entityType}
	if value := ctx.Query("status"); value != "" {
		status := models.ModerationStatus(value)
		params.Status = &status
	}
	var err error
	params.OwnerID, err = optionalUUID(ctx.Query("owner_id"))
	if err != nil {
		return params, err
	}
	params.ParentID, err = optionalUUID(ctx.Query("parent_id"))
	if err != nil {
		return params, err
	}
	params.CreatedFrom, err = optionalTime(ctx.Query("created_from"))
	if err != nil {
		return params, err
	}
	params.CreatedTo, err = optionalTime(ctx.Query("created_to"))
	if err != nil {
		return params, err
	}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	return params, nil
}

func optionalUUID(value string) (*uuid.UUID, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := uuid.Parse(value)
	return &parsed, err
}

func optionalTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	return &parsed, err
}
