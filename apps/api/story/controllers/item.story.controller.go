package controllers

import (
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (c *StoryController) AddStoryItem(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}

	var dto dtos.AddStoryItemDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.AddStoryItemPipe(ctx.Context(), userID, storyID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

func (c *StoryController) DeleteStoryItem(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	itemID, err := uuid.Parse(ctx.Params("itemID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_item_id")
	}
	var moderatorPersonaID *uuid.UUID
	if value := ctx.Query("moderator_persona_id"); value != "" {
		parsed, err := uuid.Parse(value)
		if err != nil {
			return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
		}
		moderatorPersonaID = &parsed
	}

	res := c.pipe.DeleteStoryItemPipe(ctx.Context(), userID, storyID, itemID, moderatorPersonaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func (c *StoryController) ReorderStoryItem(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	storyID, err := uuid.Parse(ctx.Params("storyID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_id")
	}
	itemID, err := uuid.Parse(ctx.Params("itemID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_story_item_id")
	}

	var dto dtos.ReorderStoryItemDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	res := c.pipe.ReorderStoryItemPipe(ctx.Context(), userID, storyID, itemID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
