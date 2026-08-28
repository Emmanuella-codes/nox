package controllers

import (
	"github.com/emmanuella-codes/nox/messaging/dtos"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AddMembers adds members to one group conversation.
func (c *MessagingController) AddMembers(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	var dto dtos.AddMembersDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.AddMembersPipe(ctx.Context(), userID, conversationID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}

// RemoveMember removes one member from a group conversation.
func (c *MessagingController) RemoveMember(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	targetPersonaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	var dto dtos.RemoveMemberDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.RemoveMemberPipe(ctx.Context(), userID, conversationID, targetPersonaID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

// LeaveConversation removes the current member from a group conversation.
func (c *MessagingController) LeaveConversation(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	var dto dtos.LeaveConversationDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.LeaveConversationPipe(ctx.Context(), userID, conversationID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

// UpdateMemberRole changes one group member role through an admin action.
func (c *MessagingController) UpdateMemberRole(ctx *fiber.Ctx) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}
	conversationID, err := uuid.Parse(ctx.Params("conversationID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_conversation_id")
	}
	targetPersonaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}
	var dto dtos.UpdateMemberRoleDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}
	res := c.pipe.UpdateMemberRolePipe(ctx.Context(), userID, conversationID, targetPersonaID, dto)
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}
	return pipeSuccess(ctx, fiber.StatusOK, res.Message, res.Data)
}
