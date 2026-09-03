package controllers

import (
	"github.com/emmanuella-codes/nox/messaging/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/story/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ViewStoryItem records one authenticated viewer impression for one story item.
func (c *StoryController) ViewStoryItem(ctx *fiber.Ctx) error {
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
	var dto dtos.StoryItemViewDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.ViewStoryItemPipe(ctx.Context(), userID, storyID, itemID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// ListStoryItemViewers lists one story item's viewers for the story owner.
func (c *StoryController) ListStoryItemViewers(ctx *fiber.Ctx) error {
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
	res := c.pipe.ListStoryItemViewersPipe(ctx.Context(), userID, storyID, itemID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// ReactToStoryItem upserts one reaction from the current viewer onto one story item.
func (c *StoryController) ReactToStoryItem(ctx *fiber.Ctx) error {
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
	var dto dtos.StoryItemReactionDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.ReactToStoryItemPipe(ctx.Context(), userID, storyID, itemID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// RemoveStoryItemReaction removes the current viewer's reaction from one story item.
func (c *StoryController) RemoveStoryItemReaction(ctx *fiber.Ctx) error {
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
	personaID, err := requiredPersonaQuery(ctx, "persona_id")
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.RemoveStoryItemReactionPipe(ctx.Context(), userID, storyID, itemID, personaID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// ReplyToStoryItem sends one story reply into a direct conversation with the story owner.
func (c *StoryController) ReplyToStoryItem(ctx *fiber.Ctx) error {
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
	var dto dtos.StoryItemReplyDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.ReplyToStoryItemPipe(ctx.Context(), userID, storyID, itemID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[pipes.MessageResponse](ctx, fiber.StatusCreated, res.Message, res.Data)
}
