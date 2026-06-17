package controllers

import (
	"strconv"
	"time"

	"github.com/emmanuella-codes/nox/media/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *MediaController) InitiateSetVideoUpload(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.InitiateSetVideoUploadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.InitiateSetVideoUploadPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *MediaController) InitiateStoryVideoUpload(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.InitiateStoryVideoUploadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.InitiateStoryVideoUploadPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *MediaController) InitiatePostMediaUpload(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.InitiatePostMediaUploadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.InitiatePostMediaUploadPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *MediaController) ConfirmPostMediaUpload(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	var dto dtos.ConfirmPostMediaUploadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.ConfirmPostMediaUploadPipe(ctx.Context(), userID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *MediaController) CompleteMediaProcessing(ctx *fiber.Ctx) error {
	if !c.validProcessingSecret(ctx) {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_processing_secret")
	}
	mediaAssetID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_media_asset_id")
	}

	var dto dtos.CompleteMediaProcessingDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.CompleteMediaProcessingPipe(ctx.Context(), mediaAssetID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *MediaController) CompleteStoryMediaProcessing(ctx *fiber.Ctx) error {
	if !c.validProcessingSecret(ctx) {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_processing_secret")
	}
	mediaAssetID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_media_asset_id")
	}

	var dto dtos.CompleteMediaProcessingDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.CompleteStoryMediaProcessingPipe(ctx.Context(), mediaAssetID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *MediaController) FailMediaProcessing(ctx *fiber.Ctx) error {
	if !c.validProcessingSecret(ctx) {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_processing_secret")
	}
	mediaAssetID, err := uuid.Parse(ctx.Params("mediaAssetID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_media_asset_id")
	}

	res := c.pipe.FailMediaProcessingPipe(ctx.Context(), mediaAssetID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

func (c *MediaController) CleanupOrphanedMediaAssets(ctx *fiber.Ctx) error {
	if !c.validProcessingSecret(ctx) {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_processing_secret")
	}
	olderThanHours, err := strconv.Atoi(ctx.Query("older_than_hours", "24"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_older_than_hours")
	}
	limit, err := strconv.Atoi(ctx.Query("limit", "100"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_limit")
	}

	res := c.pipe.CleanupOrphanedMediaAssetsPipe(ctx.Context(), time.Duration(olderThanHours)*time.Hour, limit)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
