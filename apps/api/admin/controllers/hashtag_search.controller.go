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

func (ac *AdminController) ListAdminHashtags(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListAdminHashtagsParams{Query: ctx.Query("query")}
	if value := ctx.Query("status"); value != "" {
		status := modelsModerationStatus(value)
		params.Status = &status
	}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListAdminHashtagsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ModerateHashtag(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.ModerateHashtagDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.ModerateHashtagPipe(requestContext(ctx), adminUserID, ctx.Params("tag"), dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) SuppressHashtag(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.SuppressHashtagDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.SuppressHashtagPipe(requestContext(ctx), adminUserID, ctx.Params("tag"), dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) UnsuppressHashtag(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	res := ac.pipe.UnsuppressHashtagPipe(requestContext(ctx), adminUserID, ctx.Params("tag"))
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ListSearchSuppressions(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params := adminrepo.ListSearchSuppressionsParams{Query: ctx.Query("query")}
	params.Limit, _ = strconv.Atoi(ctx.Query("limit", "50"))
	params.Offset, _ = strconv.Atoi(ctx.Query("offset", "0"))
	res := ac.pipe.ListSearchSuppressionsPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) CreateSearchSuppression(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	var dto dtos.CreateSearchSuppressionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.CreateSearchSuppressionPipe(requestContext(ctx), adminUserID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) DeleteSearchSuppression(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	id, err := uuid.Parse(ctx.Params("suppressionID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.DeleteSearchSuppressionPipe(requestContext(ctx), adminUserID, id)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ReindexSearch(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	res := ac.pipe.ReindexSearchPipe(requestContext(ctx), adminUserID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func modelsModerationStatus(value string) models.ModerationStatus {
	return models.ModerationStatus(value)
}
