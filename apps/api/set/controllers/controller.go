package controllers

import (
	"strconv"

	"github.com/emmanuella-codes/nox/set/messages"
	"github.com/emmanuella-codes/nox/set/pipes"
	"github.com/emmanuella-codes/nox/shared"
	sharedapi "github.com/emmanuella-codes/nox/shared/api"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SetController struct {
	pipe *pipes.SetPipe
}

func NewSetController(pipe *pipes.SetPipe) *SetController {
	return &SetController{pipe: pipe}
}

func parseAndValidate(ctx *fiber.Ctx, dto any) error {
	if err := ctx.BodyParser(dto); err != nil {
		return err
	}
	success, err := sharedapi.ValidateAPIData(dto)
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
	case messages.Invalid_Payload, messages.Invalid_Set, messages.Media_In_Use:
		return fiber.StatusBadRequest
	case messages.Persona_Not_Found, messages.Media_Not_Found, messages.Set_Not_Found, messages.Comment_Not_Found:
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
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		return fallback
	}
	return limit
}

func queryOffset(ctx *fiber.Ctx) int {
	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		return 0
	}
	return offset
}

func optionalViewerPersonaID(ctx *fiber.Ctx) (*uuid.UUID, error) {
	value := ctx.Query("viewer_persona_id")
	if value == "" {
		return nil, nil
	}
	id, err := uuid.Parse(value)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
