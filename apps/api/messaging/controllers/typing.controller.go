package controllers

import (
	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UpdateTyping updates realtime typing state for one conversation member.
func (c *MessagingController) UpdateTyping(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	var dto dtos.TypingDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateTypingPipe(ctx.Context(), userID, conversationID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}
