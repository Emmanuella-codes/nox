package controllers

import (
	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SendMessage handles message creation for one conversation member.
func (c *MessagingController) SendMessage(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	var dto dtos.SendMessageDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	if dto.IdempotencyKey == "" {
		dto.IdempotencyKey = ctx.Get("Idempotency-Key")
	}
	res := c.pipe.SendMessagePipe(ctx.Context(), userID, conversationID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusCreated, res.Message, res.Data)
}

// ListMessages lists visible messages for one conversation member.
func (c *MessagingController) ListMessages(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	personaID, err := uuid.Parse(ctx.Query("persona_id"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	res := c.pipe.ListMessagesPipe(ctx.Context(), userID, conversationID, personaID, queryLimit(ctx, 30), queryOffset(ctx))
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// MarkRead advances the conversation read cursor for one member.
func (c *MessagingController) MarkRead(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	var dto dtos.MarkReadDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.MarkReadPipe(ctx.Context(), userID, conversationID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// EditMessage updates the body of one sender-owned message.
func (c *MessagingController) EditMessage(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	messageID, err := uuid.Parse(ctx.Params("messageID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_message_id")
	}
	var dto dtos.EditMessageDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.EditMessagePipe(ctx.Context(), userID, messageID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// DeleteMessage deletes one sender-owned message for every member.
func (c *MessagingController) DeleteMessage(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	messageID, err := uuid.Parse(ctx.Params("messageID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_message_id")
	}
	res := c.pipe.DeleteMessagePipe(ctx.Context(), userID, messageID)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
