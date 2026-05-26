package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/follow/dtos"
	"github.com/emmanuella-codes/nox/follow/messages"
	"github.com/emmanuella-codes/nox/follow/pipes"
	"github.com/emmanuella-codes/nox/middleware"
	"github.com/emmanuella-codes/nox/shared"
	shared_api "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FollowController struct {
	pipe *pipes.FollowPipe
}

func NewFollowController(pipe *pipes.FollowPipe) *FollowController {
	return &FollowController{pipe: pipe}
}

func (c *FollowController) followAction(ctx *fiber.Ctx, follow bool) error {
	userID, ok := middleware.CurrentUserID(ctx)
	if !ok {
		return pipeError(ctx, fiber.StatusUnauthorized, "invalid_token")
	}

	targetPersonaID, err := uuid.Parse(ctx.Params("personaID"))
	if err != nil {
		return pipeError(ctx, fiber.StatusBadRequest, "invalid_persona_id")
	}

	var dto dtos.FollowDTO
	if err := parseAndValidate(ctx, &dto); err != nil {
		return validationError(ctx, err)
	}

	var res *shared.PipeRes[any]
	if follow {
		res = c.pipe.FollowPersonaPipe(ctx.Context(), userID, dto.PersonaID, targetPersonaID)
	} else {
		res = c.pipe.UnfollowPersonaPipe(ctx.Context(), userID, dto.PersonaID, targetPersonaID)
	}
	if !res.Success {
		return pipeError(ctx, pipeErrorStatus(res.Message), res.Message)
	}

	return pipeSuccess[any](ctx, fiber.StatusOK, res.Message, nil)
}

func parseAndValidate(ctx *fiber.Ctx, dto any) error {
	if err := ctx.BodyParser(dto); err != nil {
		return err
	}
	success, err := shared_api.ValidateAPIData(dto)
	if !success {
		return err
	}
	return nil
}

func validationError(ctx *fiber.Ctx, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": messages.Invalid_Payload, "error": err.Error()})
}

func pipeSuccess[T any](ctx *fiber.Ctx, status int, message shared.PipeMessage, data *T) error {
	return ctx.Status(status).JSON(shared.PipeRes[T]{Success: true, Message: message, Data: data})
}

func pipeError(ctx *fiber.Ctx, status int, message shared.PipeMessage) error {
	return ctx.Status(status).JSON(shared.PipeRes[any]{Success: false, Message: message})
}

func pipeErrorStatus(message shared.PipeMessage) int {
	switch message {
	case messages.Invalid_Payload, messages.Already_Following, messages.Not_Following, messages.Self_Follow_Not_Allowed:
		return fiber.StatusBadRequest
	case messages.Persona_Not_Found:
		return fiber.StatusNotFound
	case messages.Forbidden:
		return fiber.StatusForbidden
	case messages.Internal_Error:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusBadRequest
	}
}

func queryLimit(ctx *fiber.Ctx, fallback int) int {
	value := ctx.Query("limit")
	if value == "" {
		return fallback
	}
	limit, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return limit
}
