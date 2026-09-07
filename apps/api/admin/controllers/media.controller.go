package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/emmanuella-codes/nox/admin/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/models"
	adminrepo "github.com/emmanuella-codes/nox/repositories/admin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (ac *AdminController) ListAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	params, err := adminMediaParams(ctx)
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.ListAdminMediaPipe(requestContext(ctx), adminUserID, params)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) GetAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	mediaID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.GetAdminMediaPipe(requestContext(ctx), adminUserID, mediaID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) RetryAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	mediaID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.RetryMediaPipe(requestContext(ctx), adminUserID, mediaID)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) CorrectAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	mediaID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.CorrectMediaDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.CorrectMediaPipe(requestContext(ctx), adminUserID, mediaID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) ModerateAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	mediaID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	var dto dtos.ModerateMediaDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := ac.pipe.ModerateMediaPipe(requestContext(ctx), adminUserID, mediaID, dto)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func (ac *AdminController) CleanupAdminMedia(ctx *fiber.Ctx) error {
	adminUserID, ok := middleware.CurrentAdminUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	olderThanHours, err := strconv.Atoi(ctx.Query("older_than_hours", "24"))
	if err != nil || olderThanHours <= 0 {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "100"))
	if err != nil || limit <= 0 {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_payload")
	}
	res := ac.pipe.CleanupAdminMediaPipe(requestContext(ctx), adminUserID, time.Duration(olderThanHours)*time.Hour, limit)
	if res.Success {
		return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
	}
	return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
}

func adminMediaParams(ctx *fiber.Ctx) (adminrepo.ListAdminMediaParams, error) {
	params := adminrepo.ListAdminMediaParams{}
	if value := strings.TrimSpace(ctx.Query("status")); value != "" {
		status := models.MediaProcessingStatus(value)
		params.Status = &status
	}
	if value := strings.TrimSpace(ctx.Query("kind")); value != "" {
		kind := models.MediaKind(value)
		params.Kind = &kind
	}
	if value := strings.TrimSpace(ctx.Query("moderation_status")); value != "" {
		status := models.ModerationStatus(value)
		params.Moderation = &status
	}
	var err error
	if value := ctx.Query("owner_id"); value != "" {
		parsed, parseErr := uuid.Parse(value)
		if parseErr != nil {
			return params, parseErr
		}
		params.OwnerID = &parsed
	}
	if value := ctx.Query("orphaned"); value != "" {
		parsed, parseErr := strconv.ParseBool(value)
		if parseErr != nil {
			return params, parseErr
		}
		params.Orphaned = &parsed
	}
	if value := ctx.Query("older_than"); value != "" {
		parsed, parseErr := time.Parse(time.RFC3339, value)
		if parseErr != nil {
			return params, parseErr
		}
		params.OlderThan = &parsed
	}
	params.Limit, err = strconv.Atoi(ctx.Query("limit", "50"))
	if err != nil {
		return params, err
	}
	params.Offset, err = strconv.Atoi(ctx.Query("offset", "0"))
	return params, err
}
